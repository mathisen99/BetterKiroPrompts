package generation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"better-kiro-prompts/internal/logger"
	"better-kiro-prompts/internal/openai"
	"better-kiro-prompts/internal/prompts"
	"better-kiro-prompts/internal/queue"
	"better-kiro-prompts/internal/storage"
)

const (
	maxProjectIdeaLength = 2000
	maxAnswerLength      = 1000
	minQuestions         = 5
	maxQuestions         = 10
)

var (
	ErrEmptyProjectIdea   = errors.New("project idea is required")
	ErrProjectIdeaTooLong = errors.New("project idea exceeds maximum length")
	ErrAnswerTooLong      = errors.New("answer exceeds maximum length")
	ErrInvalidResponse    = errors.New("invalid response from AI")
	ErrNoQuestions        = errors.New("no questions generated")
	ErrNoFiles            = errors.New("no files generated")
)

// Question represents a follow-up question for the user.
type Question struct {
	ID       int      `json:"id"`
	Text     string   `json:"text"`
	Hint     string   `json:"hint,omitempty"`
	Examples []string `json:"examples"` // 3 clickable example answers
}

// Answer represents a user's answer to a question.
type Answer struct {
	QuestionID int    `json:"questionId"`
	Answer     string `json:"answer"`
}

// GeneratedFile represents a generated output file.
type GeneratedFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Type    string `json:"type"` // "kickoff", "steering", "hook"
}

// QuestionsResponse is the expected JSON structure from the AI for questions.
type QuestionsResponse struct {
	Questions []Question `json:"questions"`
}

// OutputsResponse is the expected JSON structure from the AI for outputs.
type OutputsResponse struct {
	Files []GeneratedFile `json:"files"`
}

// GenerationResult contains the generated files and the stored generation ID.
type GenerationResult struct {
	Files        []GeneratedFile `json:"files"`
	GenerationID string          `json:"generationId,omitempty"`
}

// Service handles AI-driven generation of questions and outputs.
type Service struct {
	openaiClient *openai.Client
	requestQueue *queue.RequestQueue
	repository   storage.Repository
	log          *slog.Logger
}

// NewService creates a new generation service.
func NewService(client *openai.Client) *Service {
	return &Service{
		openaiClient: client,
		requestQueue: nil, // Optional queue
		repository:   nil, // Optional repository
		log:          slog.Default(),
	}
}

// NewServiceWithQueue creates a new generation service with a request queue.
func NewServiceWithQueue(client *openai.Client, q *queue.RequestQueue) *Service {
	return &Service{
		openaiClient: client,
		requestQueue: q,
		repository:   nil,
		log:          slog.Default(),
	}
}

// NewServiceWithDeps creates a new generation service with all dependencies.
func NewServiceWithDeps(client *openai.Client, q *queue.RequestQueue, repo storage.Repository) *Service {
	return &Service{
		openaiClient: client,
		requestQueue: q,
		repository:   repo,
		log:          slog.Default(),
	}
}

// NewServiceWithLogger creates a new generation service with all dependencies including logger.
func NewServiceWithLogger(client *openai.Client, q *queue.RequestQueue, repo storage.Repository, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	return &Service{
		openaiClient: client,
		requestQueue: q,
		repository:   repo,
		log:          log,
	}
}

// SetLogger sets the logger for the service.
func (s *Service) SetLogger(log *slog.Logger) {
	if log != nil {
		s.log = log
	}
}

// SetRequestQueue sets the request queue for the service.
func (s *Service) SetRequestQueue(q *queue.RequestQueue) {
	s.requestQueue = q
}

// SetRepository sets the storage repository for the service.
func (s *Service) SetRepository(repo storage.Repository) {
	s.repository = repo
}

// ValidateProjectIdea validates the project idea input.
func ValidateProjectIdea(idea string) error {
	trimmed := strings.TrimSpace(idea)
	if trimmed == "" {
		return ErrEmptyProjectIdea
	}
	if len(trimmed) > maxProjectIdeaLength {
		return ErrProjectIdeaTooLong
	}
	return nil
}

// ValidateAnswers validates the answers input.
func ValidateAnswers(answers []Answer) error {
	for _, a := range answers {
		if len(a.Answer) > maxAnswerLength {
			return ErrAnswerTooLong
		}
	}
	return nil
}

// GenerateQuestions generates follow-up questions based on the project idea.
func (s *Service) GenerateQuestions(ctx context.Context, projectIdea string, experienceLevel string) ([]Question, error) {
	requestID := logger.GetRequestID(ctx)
	start := time.Now()

	s.log.Info("generate_questions_start",
		slog.String("request_id", requestID),
		slog.String("experience_level", experienceLevel),
		slog.Int("idea_length", len(projectIdea)),
	)

	if err := ValidateProjectIdea(projectIdea); err != nil {
		s.log.Warn("generate_questions_validation_failed",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	// Acquire queue slot if queue is configured
	if s.requestQueue != nil {
		s.log.Debug("queue_acquire_start", slog.String("request_id", requestID))
		if err := s.requestQueue.Acquire(ctx); err != nil {
			s.log.Error("queue_acquire_failed",
				slog.String("request_id", requestID),
				slog.String("error", err.Error()),
			)
			return nil, fmt.Errorf("failed to acquire queue slot: %w", err)
		}
		defer s.requestQueue.Release()
		s.log.Debug("queue_acquire_success", slog.String("request_id", requestID))
	}

	// Validate experience level
	if !prompts.IsValidExperienceLevel(experienceLevel) {
		experienceLevel = prompts.ExperienceNovice // Default to novice
	}

	// Use experience-level-aware system prompt
	systemPrompt := prompts.GetQuestionsSystemPrompt(experienceLevel)
	userPrompt := prompts.GetQuestionsUserPrompt(strings.TrimSpace(projectIdea), experienceLevel)

	messages := []openai.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	s.log.Debug("openai_call_start",
		slog.String("request_id", requestID),
		slog.String("operation", "generate_questions"),
	)

	response, err := s.openaiClient.ChatCompletion(ctx, messages)
	if err != nil {
		s.log.Error("generate_questions_openai_failed",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return nil, fmt.Errorf("failed to generate questions: %w", err)
	}

	s.log.Debug("openai_call_complete",
		slog.String("request_id", requestID),
		slog.String("operation", "generate_questions"),
	)

	questions, err := parseQuestionsResponse(response)
	if err != nil {
		s.log.Error("generate_questions_parse_failed",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	s.log.Info("generate_questions_complete",
		slog.String("request_id", requestID),
		slog.Int("question_count", len(questions)),
		slog.Duration("duration", time.Since(start)),
	)

	return questions, nil
}

// maxRetries is the number of retry attempts for validation failures
const maxRetries = 1

// GenerateOutputs generates kickoff prompt, steering files, hooks, and AGENTS.md.
func (s *Service) GenerateOutputs(ctx context.Context, projectIdea string, answers []Answer, experienceLevel string, hookPreset string) ([]GeneratedFile, error) {
	requestID := logger.GetRequestID(ctx)
	start := time.Now()

	s.log.Info("generate_outputs_start",
		slog.String("request_id", requestID),
		slog.String("experience_level", experienceLevel),
		slog.String("hook_preset", hookPreset),
		slog.Int("answer_count", len(answers)),
	)

	if err := ValidateProjectIdea(projectIdea); err != nil {
		s.log.Warn("generate_outputs_validation_failed",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
			slog.String("validation_type", "project_idea"),
		)
		return nil, err
	}
	if err := ValidateAnswers(answers); err != nil {
		s.log.Warn("generate_outputs_validation_failed",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
			slog.String("validation_type", "answers"),
		)
		return nil, err
	}

	// Acquire queue slot if queue is configured
	if s.requestQueue != nil {
		s.log.Debug("queue_acquire_start", slog.String("request_id", requestID))
		if err := s.requestQueue.Acquire(ctx); err != nil {
			s.log.Error("queue_acquire_failed",
				slog.String("request_id", requestID),
				slog.String("error", err.Error()),
			)
			return nil, fmt.Errorf("failed to acquire queue slot: %w", err)
		}
		defer s.requestQueue.Release()
		s.log.Debug("queue_acquire_success", slog.String("request_id", requestID))
	}

	// Validate experience level and hook preset
	if !prompts.IsValidExperienceLevel(experienceLevel) {
		experienceLevel = prompts.ExperienceNovice
	}
	if !prompts.IsValidHookPreset(hookPreset) {
		hookPreset = prompts.HookPresetDefault
	}

	// Convert answers to prompts.Answer type
	promptAnswers := make([]prompts.Answer, len(answers))
	for i, a := range answers {
		promptAnswers[i] = prompts.Answer{
			QuestionID: a.QuestionID,
			Answer:     a.Answer,
		}
	}

	// Use comprehensive system and user prompts
	systemPrompt := prompts.GetOutputsSystemPrompt(experienceLevel, hookPreset)
	userPrompt := prompts.GetOutputsUserPrompt(strings.TrimSpace(projectIdea), promptAnswers, experienceLevel, hookPreset)

	messages := []openai.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		s.log.Debug("generate_outputs_attempt",
			slog.String("request_id", requestID),
			slog.Int("attempt", attempt+1),
			slog.Int("max_attempts", maxRetries+1),
		)

		response, err := s.openaiClient.ChatCompletion(ctx, messages)
		if err != nil {
			s.log.Error("generate_outputs_openai_failed",
				slog.String("request_id", requestID),
				slog.Int("attempt", attempt+1),
				slog.String("error", err.Error()),
			)
			return nil, fmt.Errorf("failed to generate outputs: %w", err)
		}

		files, err := parseOutputsResponse(response)
		if err != nil {
			lastErr = err
			s.log.Warn("generate_outputs_parse_failed",
				slog.String("request_id", requestID),
				slog.Int("attempt", attempt+1),
				slog.String("error", err.Error()),
			)
			if attempt < maxRetries {
				// Add retry context to messages for the next attempt
				messages = append(messages,
					openai.Message{Role: "assistant", Content: response},
					openai.Message{Role: "user", Content: buildRetryPrompt(err)},
				)
				continue
			}
			return nil, FormatValidationError(err)
		}

		// Validate generated files
		if err := ValidateGeneratedFiles(files); err != nil {
			lastErr = fmt.Errorf("%w: %v", ErrInvalidResponse, err)
			s.log.Warn("generate_outputs_validation_failed",
				slog.String("request_id", requestID),
				slog.Int("attempt", attempt+1),
				slog.String("error", err.Error()),
				slog.String("validation_type", "generated_files"),
			)
			if attempt < maxRetries {
				// Add retry context to messages for the next attempt
				messages = append(messages,
					openai.Message{Role: "assistant", Content: response},
					openai.Message{Role: "user", Content: buildRetryPrompt(err)},
				)
				continue
			}
			return nil, FormatValidationError(lastErr)
		}

		s.log.Info("generate_outputs_complete",
			slog.String("request_id", requestID),
			slog.Int("file_count", len(files)),
			slog.Int("attempts_used", attempt+1),
			slog.Duration("duration", time.Since(start)),
		)

		return files, nil
	}

	// Should not reach here, but return last error if we do
	return nil, FormatValidationError(lastErr)
}

// GenerateAndStoreOutputs generates outputs and stores them in the database.
// Returns the generated files and the generation ID if storage is configured.
func (s *Service) GenerateAndStoreOutputs(ctx context.Context, projectIdea string, answers []Answer, experienceLevel string, hookPreset string) (*GenerationResult, error) {
	requestID := logger.GetRequestID(ctx)

	// Generate the outputs
	files, err := s.GenerateOutputs(ctx, projectIdea, answers, experienceLevel, hookPreset)
	if err != nil {
		return nil, err
	}

	result := &GenerationResult{
		Files: files,
	}

	// Store in database if repository is configured
	if s.repository != nil {
		s.log.Debug("storage_attempt_start",
			slog.String("request_id", requestID),
			slog.Int("file_count", len(files)),
		)

		// Convert files to JSON
		filesJSON, err := json.Marshal(files)
		if err != nil {
			s.log.Error("storage_marshal_failed",
				slog.String("request_id", requestID),
				slog.String("error", err.Error()),
			)
			// Log error but don't fail the request - user still gets their files
			return result, nil
		}

		// Get category based on project idea
		s.log.Debug("category_lookup_start",
			slog.String("request_id", requestID),
		)
		categoryID, err := s.repository.GetCategoryByKeywords(ctx, projectIdea)
		if err != nil {
			s.log.Warn("category_lookup_failed",
				slog.String("request_id", requestID),
				slog.String("error", err.Error()),
				slog.Int("default_category_id", 5),
			)
			// Default to "Other" category if lookup fails
			categoryID = 5
		} else {
			s.log.Debug("category_lookup_complete",
				slog.String("request_id", requestID),
				slog.Int("category_id", categoryID),
			)
		}

		// Create generation record
		gen := &storage.Generation{
			ProjectIdea:     strings.TrimSpace(projectIdea),
			ExperienceLevel: experienceLevel,
			HookPreset:      hookPreset,
			Files:           filesJSON,
			CategoryID:      categoryID,
		}

		if err := s.repository.CreateGeneration(ctx, gen); err != nil {
			s.log.Error("storage_create_failed",
				slog.String("request_id", requestID),
				slog.String("error", err.Error()),
			)
			// Log error but don't fail the request - user still gets their files
			return result, nil
		}

		s.log.Info("storage_complete",
			slog.String("request_id", requestID),
			slog.String("generation_id", gen.ID),
			slog.Int("category_id", categoryID),
		)

		result.GenerationID = gen.ID
	}

	return result, nil
}

// buildRetryPrompt creates a prompt explaining the validation error for retry
func buildRetryPrompt(err error) string {
	return fmt.Sprintf(`The previous response had validation errors. Please fix the following issues and regenerate the complete JSON response:

Error: %v

Remember:
- All steering files must have valid YAML frontmatter with 'inclusion' field
- fileMatch mode requires 'fileMatchPattern' field
- All hook files must have valid JSON with required fields: name, description, version, enabled, when, then
- when.type must be one of: fileEdited, fileCreated, fileDeleted, promptSubmit, agentStop, userTriggered
- then.type must be 'askAgent' or 'runCommand'
- runCommand can only be used with promptSubmit or agentStop triggers
- File-based triggers (fileEdited, fileCreated, fileDeleted) require patterns array

Please provide the corrected JSON response.`, err)
}

func parseQuestionsResponse(response string) ([]Question, error) {
	// Try to extract JSON from response (handle potential markdown code blocks)
	jsonStr := extractJSON(response)

	var qr QuestionsResponse
	if err := json.Unmarshal([]byte(jsonStr), &qr); err != nil {
		return nil, fmt.Errorf("%w: failed to parse questions JSON: %v", ErrInvalidResponse, err)
	}

	if len(qr.Questions) == 0 {
		return nil, ErrNoQuestions
	}

	// Validate question count
	if len(qr.Questions) < minQuestions || len(qr.Questions) > maxQuestions {
		// Truncate or pad if needed, but still return what we have
		if len(qr.Questions) > maxQuestions {
			qr.Questions = qr.Questions[:maxQuestions]
		}
	}

	// Validate each question has required fields
	for i, q := range qr.Questions {
		if q.Text == "" {
			return nil, fmt.Errorf("%w: question %d has empty text", ErrInvalidResponse, i+1)
		}
		// Ensure IDs are set
		if q.ID == 0 {
			qr.Questions[i].ID = i + 1
		}
	}

	return qr.Questions, nil
}

func parseOutputsResponse(response string) ([]GeneratedFile, error) {
	// Try to extract JSON from response (handle potential markdown code blocks)
	jsonStr := extractJSON(response)

	var or OutputsResponse
	if err := json.Unmarshal([]byte(jsonStr), &or); err != nil {
		return nil, fmt.Errorf("%w: failed to parse outputs JSON: %v", ErrInvalidResponse, err)
	}

	if len(or.Files) == 0 {
		return nil, ErrNoFiles
	}

	// Validate required file types
	hasKickoff := false
	hasSteering := false
	hasHook := false
	hasAgents := false

	for _, f := range or.Files {
		if f.Path == "" || f.Content == "" {
			return nil, fmt.Errorf("%w: file has empty path or content", ErrInvalidResponse)
		}
		switch f.Type {
		case "kickoff":
			hasKickoff = true
		case "steering":
			hasSteering = true
		case "hook":
			hasHook = true
		case "agents":
			hasAgents = true
		}
	}

	if !hasKickoff {
		return nil, fmt.Errorf("%w: missing kickoff file", ErrInvalidResponse)
	}
	if !hasSteering {
		return nil, fmt.Errorf("%w: missing steering file", ErrInvalidResponse)
	}
	if !hasHook {
		return nil, fmt.Errorf("%w: missing hook file", ErrInvalidResponse)
	}
	if !hasAgents {
		return nil, fmt.Errorf("%w: missing AGENTS.md file", ErrInvalidResponse)
	}

	return or.Files, nil
}

// extractJSON attempts to extract JSON from a response that might contain markdown code blocks.
func extractJSON(response string) string {
	response = strings.TrimSpace(response)

	// If it starts with {, assume it's already JSON
	if strings.HasPrefix(response, "{") {
		return response
	}

	// Try to extract from markdown code block
	if idx := strings.Index(response, "```json"); idx != -1 {
		start := idx + 7
		if end := strings.Index(response[start:], "```"); end != -1 {
			return strings.TrimSpace(response[start : start+end])
		}
	}

	// Try generic code block
	if idx := strings.Index(response, "```"); idx != -1 {
		start := idx + 3
		// Skip language identifier if present
		if newline := strings.Index(response[start:], "\n"); newline != -1 {
			start += newline + 1
		}
		if end := strings.Index(response[start:], "```"); end != -1 {
			return strings.TrimSpace(response[start : start+end])
		}
	}

	return response
}

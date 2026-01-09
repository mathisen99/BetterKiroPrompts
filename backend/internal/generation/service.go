package generation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"better-kiro-prompts/internal/openai"
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
	ID   int    `json:"id"`
	Text string `json:"text"`
	Hint string `json:"hint,omitempty"`
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

// Service handles AI-driven generation of questions and outputs.
type Service struct {
	openaiClient *openai.Client
}

// NewService creates a new generation service.
func NewService(client *openai.Client) *Service {
	return &Service{
		openaiClient: client,
	}
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
func (s *Service) GenerateQuestions(ctx context.Context, projectIdea string) ([]Question, error) {
	if err := ValidateProjectIdea(projectIdea); err != nil {
		return nil, err
	}

	prompt := buildQuestionsPrompt(strings.TrimSpace(projectIdea))

	messages := []openai.Message{
		{Role: "system", Content: questionsSystemPrompt},
		{Role: "user", Content: prompt},
	}

	response, err := s.openaiClient.ChatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate questions: %w", err)
	}

	questions, err := parseQuestionsResponse(response)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

// GenerateOutputs generates kickoff prompt, steering files, and hooks.
func (s *Service) GenerateOutputs(ctx context.Context, projectIdea string, answers []Answer) ([]GeneratedFile, error) {
	if err := ValidateProjectIdea(projectIdea); err != nil {
		return nil, err
	}
	if err := ValidateAnswers(answers); err != nil {
		return nil, err
	}

	prompt := buildOutputsPrompt(strings.TrimSpace(projectIdea), answers)

	messages := []openai.Message{
		{Role: "system", Content: outputsSystemPrompt},
		{Role: "user", Content: prompt},
	}

	response, err := s.openaiClient.ChatCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate outputs: %w", err)
	}

	files, err := parseOutputsResponse(response)
	if err != nil {
		return nil, err
	}

	return files, nil
}

const questionsSystemPrompt = `You are helping a developer plan their project. Based on their project idea, generate 5-10 follow-up questions to understand their requirements better.

Rules:
- Adapt question complexity to the project sophistication
- For simple projects (games, basic apps): focus on scope, platform, basic features
- For complex projects (distributed systems, APIs): include architecture, scalability, data consistency
- Questions should help clarify: users, data, auth, tech stack, constraints
- Return ONLY valid JSON, no markdown code blocks

Response format:
{"questions": [{"id": 1, "text": "...", "hint": "..."}]}`

const outputsSystemPrompt = `You are generating Kiro project files for a developer. Based on their project idea and answers, generate:

1. A kickoff prompt (markdown) that summarizes the project requirements
2. Steering files for .kiro/steering/ with appropriate frontmatter
3. Hook files for .kiro/hooks/ in valid Kiro hook JSON format

Rules:
- Kickoff prompt should enforce "answer questions before coding" principle
- Steering files should be concise and actionable
- Include product.md, tech.md, structure.md at minimum
- Add security/quality steering if project warrants it
- Hooks should match project tech stack (Go, TypeScript, React, etc.)
- Use valid Kiro hook schema with name, description, version, enabled, when, then
- Return ONLY valid JSON, no markdown code blocks

Response format:
{
  "files": [
    {"path": "kickoff-prompt.md", "content": "...", "type": "kickoff"},
    {"path": ".kiro/steering/product.md", "content": "...", "type": "steering"},
    {"path": ".kiro/hooks/format.kiro.hook", "content": "...", "type": "hook"}
  ]
}`

func buildQuestionsPrompt(projectIdea string) string {
	return fmt.Sprintf("Project idea: %s", projectIdea)
}

func buildOutputsPrompt(projectIdea string, answers []Answer) string {
	answersJSON, _ := json.Marshal(answers)
	return fmt.Sprintf("Project idea: %s\nAnswers: %s", projectIdea, string(answersJSON))
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

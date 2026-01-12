package api

import (
	"better-kiro-prompts/internal/generation"
	"better-kiro-prompts/internal/ratelimit"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// ExperienceLevel represents the user's programming experience level.
type ExperienceLevel string

const (
	ExperienceLevelBeginner ExperienceLevel = "beginner"
	ExperienceLevelNovice   ExperienceLevel = "novice"
	ExperienceLevelExpert   ExperienceLevel = "expert"
)

// ValidExperienceLevels contains all valid experience level values.
var ValidExperienceLevels = map[ExperienceLevel]bool{
	ExperienceLevelBeginner: true,
	ExperienceLevelNovice:   true,
	ExperienceLevelExpert:   true,
}

// HookPreset represents the hook configuration preset.
type HookPreset string

const (
	HookPresetLight   HookPreset = "light"
	HookPresetBasic   HookPreset = "basic"
	HookPresetDefault HookPreset = "default"
	HookPresetStrict  HookPreset = "strict"
)

// ValidHookPresets contains all valid hook preset values.
var ValidHookPresets = map[HookPreset]bool{
	HookPresetLight:   true,
	HookPresetBasic:   true,
	HookPresetDefault: true,
	HookPresetStrict:  true,
}

// GenerateQuestionsRequest is the request body for generating questions.
type GenerateQuestionsRequest struct {
	ProjectIdea     string          `json:"projectIdea"`
	ExperienceLevel ExperienceLevel `json:"experienceLevel"`
}

// GenerateQuestionsResponse is the response body for generated questions.
type GenerateQuestionsResponse struct {
	Questions []generation.Question `json:"questions"`
}

// GenerateOutputsRequest is the request body for generating outputs.
type GenerateOutputsRequest struct {
	ProjectIdea     string              `json:"projectIdea"`
	Answers         []generation.Answer `json:"answers"`
	ExperienceLevel ExperienceLevel     `json:"experienceLevel"`
	HookPreset      HookPreset          `json:"hookPreset"`
}

// GenerateOutputsResponse is the response body for generated outputs.
type GenerateOutputsResponse struct {
	Files []generation.GeneratedFile `json:"files"`
}

// ErrorResponse is the standard error response format.
type ErrorResponse struct {
	Error      string `json:"error"`
	RetryAfter int    `json:"retryAfter,omitempty"`
}

// GenerateHandler holds dependencies for generation endpoints.
type GenerateHandler struct {
	service     *generation.Service
	rateLimiter *ratelimit.Limiter
}

// NewGenerateHandler creates a new handler with the given dependencies.
func NewGenerateHandler(service *generation.Service, limiter *ratelimit.Limiter) *GenerateHandler {
	return &GenerateHandler{
		service:     service,
		rateLimiter: limiter,
	}
}

// HandleGenerateQuestions handles POST /api/generate/questions.
func (h *GenerateHandler) HandleGenerateQuestions(w http.ResponseWriter, r *http.Request) {
	// Check rate limit
	ip := getClientIP(r)
	allowed, retryAfter := h.rateLimiter.Allow(ip)
	if !allowed {
		writeError(w, http.StatusTooManyRequests, "Rate limit exceeded", int(retryAfter.Seconds()))
		return
	}

	// Parse request body
	var req GenerateQuestionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", 0)
		return
	}

	// Validate input
	if err := generation.ValidateProjectIdea(req.ProjectIdea); err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), 0)
		return
	}

	// Validate experience level
	if err := validateExperienceLevel(req.ExperienceLevel); err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), 0)
		return
	}

	// Generate questions
	questions, err := h.service.GenerateQuestions(r.Context(), req.ProjectIdea, string(req.ExperienceLevel))
	if err != nil {
		handleGenerationError(w, err)
		return
	}

	// Return response
	writeJSON(w, http.StatusOK, GenerateQuestionsResponse{Questions: questions})
}

// HandleGenerateOutputs handles POST /api/generate/outputs.
func (h *GenerateHandler) HandleGenerateOutputs(w http.ResponseWriter, r *http.Request) {
	// Check rate limit
	ip := getClientIP(r)
	allowed, retryAfter := h.rateLimiter.Allow(ip)
	if !allowed {
		writeError(w, http.StatusTooManyRequests, "Rate limit exceeded", int(retryAfter.Seconds()))
		return
	}

	// Parse request body
	var req GenerateOutputsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", 0)
		return
	}

	// Validate input
	if err := generation.ValidateProjectIdea(req.ProjectIdea); err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), 0)
		return
	}
	if err := generation.ValidateAnswers(req.Answers); err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), 0)
		return
	}

	// Validate experience level
	if err := validateExperienceLevel(req.ExperienceLevel); err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), 0)
		return
	}

	// Validate hook preset
	if err := validateHookPreset(req.HookPreset); err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), 0)
		return
	}

	// Generate outputs
	files, err := h.service.GenerateOutputs(r.Context(), req.ProjectIdea, req.Answers, string(req.ExperienceLevel), string(req.HookPreset))
	if err != nil {
		handleGenerationError(w, err)
		return
	}

	// Return response
	writeJSON(w, http.StatusOK, GenerateOutputsResponse{Files: files})
}

// getClientIP extracts the client IP from the request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxied requests)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	addr := r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}

// handleGenerationError converts generation errors to appropriate HTTP responses.
func handleGenerationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, generation.ErrEmptyProjectIdea),
		errors.Is(err, generation.ErrProjectIdeaTooLong),
		errors.Is(err, generation.ErrAnswerTooLong):
		writeError(w, http.StatusBadRequest, err.Error(), 0)
	case errors.Is(err, generation.ErrInvalidResponse),
		errors.Is(err, generation.ErrNoQuestions),
		errors.Is(err, generation.ErrNoFiles):
		writeError(w, http.StatusInternalServerError, "Generation failed. Please try again later.", 0)
	default:
		// Check for timeout
		if strings.Contains(err.Error(), "timed out") {
			writeError(w, http.StatusGatewayTimeout, "Request timed out", 0)
			return
		}
		writeError(w, http.StatusInternalServerError, "Generation failed. Please try again later.", 0)
	}
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError writes an error response.
func writeError(w http.ResponseWriter, status int, message string, retryAfter int) {
	resp := ErrorResponse{
		Error: message,
	}
	if retryAfter > 0 {
		resp.RetryAfter = retryAfter
		w.Header().Set("Retry-After", string(rune(retryAfter)))
	}
	writeJSON(w, status, resp)
}

// validateExperienceLevel validates the experience level value.
func validateExperienceLevel(level ExperienceLevel) error {
	if level == "" {
		return errors.New("experience level is required")
	}
	if !ValidExperienceLevels[level] {
		return errors.New("invalid experience level: must be 'beginner', 'novice', or 'expert'")
	}
	return nil
}

// validateHookPreset validates the hook preset value.
func validateHookPreset(preset HookPreset) error {
	if preset == "" {
		return errors.New("hook preset is required")
	}
	if !ValidHookPresets[preset] {
		return errors.New("invalid hook preset: must be 'light', 'basic', 'default', or 'strict'")
	}
	return nil
}

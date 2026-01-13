package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"better-kiro-prompts/internal/ratelimit"
	"better-kiro-prompts/internal/scanner"
)

// ScanRequest is the request body for starting a scan.
type ScanRequest struct {
	RepoURL string `json:"repo_url"`
}

// ScanConfigResponse is the response for scan configuration.
type ScanConfigResponse struct {
	PrivateRepoEnabled bool `json:"private_repo_enabled"`
	AIReviewEnabled    bool `json:"ai_review_enabled,omitempty"`
	MaxFilesToReview   int  `json:"max_files_to_review,omitempty"`
}

// ScanHandler holds dependencies for scan endpoints.
type ScanHandler struct {
	service     *scanner.Service
	rateLimiter *ratelimit.Limiter
}

// NewScanHandler creates a new handler with the given dependencies.
func NewScanHandler(service *scanner.Service, limiter *ratelimit.Limiter) *ScanHandler {
	return &ScanHandler{
		service:     service,
		rateLimiter: limiter,
	}
}

// HandleStartScan handles POST /api/scan - Start a new security scan.
func (h *ScanHandler) HandleStartScan(w http.ResponseWriter, r *http.Request) {
	// Check rate limit
	ip := getClientIP(r)
	allowed, retryAfter := h.rateLimiter.Allow(ip)
	if !allowed {
		WriteRateLimited(w, r, int(retryAfter.Seconds()))
		return
	}

	// Parse request body
	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteBadRequest(w, r, "Invalid request body")
		return
	}

	// Validate URL format first
	if validationErr := scanner.ValidateGitHubURL(req.RepoURL); validationErr != nil {
		WriteValidationError(w, r, validationErr.Message)
		return
	}

	// Start the scan
	job, err := h.service.StartScan(r.Context(), scanner.ScanRequest{
		RepoURL: req.RepoURL,
	})
	if err != nil {
		handleScanError(w, r, err)
		return
	}

	// Return 202 Accepted with job info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(job)
}

// HandleGetScan handles GET /api/scan/{id} - Get scan status and results.
func (h *ScanHandler) HandleGetScan(w http.ResponseWriter, r *http.Request) {
	// Extract job ID from path
	jobID := r.PathValue("id")
	if jobID == "" {
		WriteBadRequest(w, r, "Scan job ID is required")
		return
	}

	// Get the job
	job, err := h.service.GetJob(r.Context(), jobID)
	if err != nil {
		if errors.Is(err, scanner.ErrJobNotFound) {
			WriteNotFound(w, r, "Scan job not found")
			return
		}
		WriteInternalError(w, r, "Failed to retrieve scan job")
		return
	}

	// Return job info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(job)
}

// HandleGetScanConfig handles GET /api/scan/config - Get scan configuration.
func (h *ScanHandler) HandleGetScanConfig(w http.ResponseWriter, r *http.Request) {
	config := h.service.GetConfig()

	resp := ScanConfigResponse{
		PrivateRepoEnabled: config["private_repo_enabled"].(bool),
	}

	// Include optional fields if available
	if aiEnabled, ok := config["ai_review_enabled"].(bool); ok {
		resp.AIReviewEnabled = aiEnabled
	}
	if maxFiles, ok := config["max_files_to_review"].(int); ok {
		resp.MaxFilesToReview = maxFiles
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// handleScanError converts scan errors to appropriate HTTP responses.
func handleScanError(w http.ResponseWriter, r *http.Request, err error) {
	// Check for validation errors
	var validationErr *scanner.ValidationError
	if errors.As(err, &validationErr) {
		WriteValidationError(w, r, validationErr.Message)
		return
	}

	// Check for specific error types
	if errors.Is(err, scanner.ErrJobNotFound) {
		WriteNotFound(w, r, "Scan job not found")
		return
	}

	if errors.Is(err, scanner.ErrScanFailed) {
		WriteInternalError(w, r, "Scan failed. Please try again later.")
		return
	}

	// Default to internal error
	WriteInternalError(w, r, "Failed to start scan. Please try again later.")
}

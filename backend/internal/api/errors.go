package api

import (
	"encoding/json"
	"net/http"
)

// Error codes for structured error responses.
// Client errors (4xx) start with "CLIENT_", server errors (5xx) start with "SERVER_".
const (
	// Client errors (4xx)
	ErrCodeValidation   = "CLIENT_VALIDATION"
	ErrCodeRateLimited  = "CLIENT_RATE_LIMITED"
	ErrCodeNotFound     = "CLIENT_NOT_FOUND"
	ErrCodeBadRequest   = "CLIENT_BAD_REQUEST"
	ErrCodeUnauthorized = "CLIENT_UNAUTHORIZED"

	// Server errors (5xx)
	ErrCodeInternal    = "SERVER_INTERNAL"
	ErrCodeTimeout     = "SERVER_TIMEOUT"
	ErrCodeUnavailable = "SERVER_UNAVAILABLE"
)

// ErrorResponse represents a structured error response.
type ErrorResponse struct {
	Error      string `json:"error"`                // Human-readable error message
	Code       string `json:"code"`                 // Machine-readable error code
	RequestID  string `json:"requestId,omitempty"`  // Request ID for tracking
	RetryAfter int    `json:"retryAfter,omitempty"` // Seconds until retry (for rate limiting)
}

// WriteError writes a structured error response to the response writer.
func WriteError(w http.ResponseWriter, r *http.Request, statusCode int, code string, message string) {
	resp := ErrorResponse{
		Error:     message,
		Code:      code,
		RequestID: GetRequestID(r.Context()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

// WriteErrorWithRetry writes a structured error response with retry-after header.
func WriteErrorWithRetry(w http.ResponseWriter, r *http.Request, statusCode int, code string, message string, retryAfterSeconds int) {
	resp := ErrorResponse{
		Error:      message,
		Code:       code,
		RequestID:  GetRequestID(r.Context()),
		RetryAfter: retryAfterSeconds,
	}

	w.Header().Set("Content-Type", "application/json")
	if retryAfterSeconds > 0 {
		w.Header().Set("Retry-After", string(rune(retryAfterSeconds)))
	}
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

// Common error response helpers

// WriteBadRequest writes a 400 Bad Request error.
func WriteBadRequest(w http.ResponseWriter, r *http.Request, message string) {
	WriteError(w, r, http.StatusBadRequest, ErrCodeBadRequest, message)
}

// WriteValidationError writes a 400 validation error.
func WriteValidationError(w http.ResponseWriter, r *http.Request, message string) {
	WriteError(w, r, http.StatusBadRequest, ErrCodeValidation, message)
}

// WriteNotFound writes a 404 Not Found error.
func WriteNotFound(w http.ResponseWriter, r *http.Request, message string) {
	WriteError(w, r, http.StatusNotFound, ErrCodeNotFound, message)
}

// WriteRateLimited writes a 429 Too Many Requests error.
func WriteRateLimited(w http.ResponseWriter, r *http.Request, retryAfterSeconds int) {
	WriteErrorWithRetry(w, r, http.StatusTooManyRequests, ErrCodeRateLimited,
		"Too many requests. Please try again later.", retryAfterSeconds)
}

// WriteInternalError writes a 500 Internal Server Error.
func WriteInternalError(w http.ResponseWriter, r *http.Request, message string) {
	// Don't expose internal error details to clients
	if message == "" {
		message = "An internal error occurred. Please try again later."
	}
	WriteError(w, r, http.StatusInternalServerError, ErrCodeInternal, message)
}

// WriteTimeout writes a 504 Gateway Timeout error.
func WriteTimeout(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusGatewayTimeout, ErrCodeTimeout,
		"Request timed out. Please try again.")
}

// WriteServiceUnavailable writes a 503 Service Unavailable error.
func WriteServiceUnavailable(w http.ResponseWriter, r *http.Request, retryAfterSeconds int) {
	WriteErrorWithRetry(w, r, http.StatusServiceUnavailable, ErrCodeUnavailable,
		"Service temporarily unavailable. Please try again later.", retryAfterSeconds)
}

// IsClientError returns true if the error code indicates a client error.
func IsClientError(code string) bool {
	return len(code) >= 7 && code[:7] == "CLIENT_"
}

// IsServerError returns true if the error code indicates a server error.
func IsServerError(code string) bool {
	return len(code) >= 7 && code[:7] == "SERVER_"
}

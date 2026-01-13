package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/quick"
)

// Property 14: Structured Error Responses
// For any API error, the response SHALL include error message, error code, and request_id fields.
// Client errors (4xx) SHALL have codes starting with "CLIENT_", and server errors (5xx) SHALL
// have codes starting with "SERVER_".
// Validates: Requirements 10.5, 10.6
//
// Feature: final-polish, Property 14: Structured Error Responses

// TestWriteError_Property_StructuredResponse tests that all error responses have required fields.
// Property: For any error response, it SHALL contain error, code, and requestId fields.
func TestWriteError_Property_StructuredResponse(t *testing.T) {
	property := func(message string, code string) bool {
		// Skip empty inputs
		if message == "" || code == "" {
			return true
		}

		// Create a request with a request ID
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req = req.WithContext(setRequestID(req.Context(), "test-request-id"))

		// Create response recorder
		w := httptest.NewRecorder()

		// Write error
		WriteError(w, req, http.StatusBadRequest, code, message)

		// Parse response
		var resp ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			return false
		}

		// Verify required fields are present
		if resp.Error == "" {
			return false
		}
		if resp.Code == "" {
			return false
		}
		if resp.RequestID == "" {
			return false
		}

		// Verify values match
		if resp.Error != message {
			return false
		}
		if resp.Code != code {
			return false
		}
		if resp.RequestID != "test-request-id" {
			return false
		}

		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: error responses should have required fields: %v", err)
	}
}

// TestClientErrorCodes_Property tests that client error codes start with "CLIENT_".
// Property: For any client error (4xx), the error code SHALL start with "CLIENT_".
func TestClientErrorCodes_Property(t *testing.T) {
	clientErrorCodes := []string{
		ErrCodeValidation,
		ErrCodeRateLimited,
		ErrCodeNotFound,
		ErrCodeBadRequest,
		ErrCodeUnauthorized,
	}

	for _, code := range clientErrorCodes {
		if !strings.HasPrefix(code, "CLIENT_") {
			t.Errorf("Client error code %q should start with 'CLIENT_'", code)
		}
		if !IsClientError(code) {
			t.Errorf("IsClientError(%q) should return true", code)
		}
		if IsServerError(code) {
			t.Errorf("IsServerError(%q) should return false for client error", code)
		}
	}
}

// TestServerErrorCodes_Property tests that server error codes start with "SERVER_".
// Property: For any server error (5xx), the error code SHALL start with "SERVER_".
func TestServerErrorCodes_Property(t *testing.T) {
	serverErrorCodes := []string{
		ErrCodeInternal,
		ErrCodeTimeout,
		ErrCodeUnavailable,
	}

	for _, code := range serverErrorCodes {
		if !strings.HasPrefix(code, "SERVER_") {
			t.Errorf("Server error code %q should start with 'SERVER_'", code)
		}
		if !IsServerError(code) {
			t.Errorf("IsServerError(%q) should return true", code)
		}
		if IsClientError(code) {
			t.Errorf("IsClientError(%q) should return false for server error", code)
		}
	}
}

// TestWriteRateLimited_Property tests rate limited responses include retry-after.
// Property: For any rate limited response, it SHALL include retryAfter field.
func TestWriteRateLimited_Property(t *testing.T) {
	property := func(retryAfter uint16) bool {
		// Use uint16 to keep values reasonable
		retrySeconds := int(retryAfter % 3600) // Max 1 hour

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req = req.WithContext(setRequestID(req.Context(), "test-id"))
		w := httptest.NewRecorder()

		WriteRateLimited(w, req, retrySeconds)

		// Verify status code
		if w.Code != http.StatusTooManyRequests {
			return false
		}

		// Parse response
		var resp ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			return false
		}

		// Verify error code
		if resp.Code != ErrCodeRateLimited {
			return false
		}

		// Verify retryAfter matches
		if resp.RetryAfter != retrySeconds {
			return false
		}

		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: rate limited responses should include retryAfter: %v", err)
	}
}

// TestErrorHelpers_Property tests that error helper functions produce correct status codes.
// Property: Each error helper function SHALL produce the correct HTTP status code.
func TestErrorHelpers_Property(t *testing.T) {
	testCases := []struct {
		name           string
		writeFunc      func(http.ResponseWriter, *http.Request)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "WriteBadRequest",
			writeFunc: func(w http.ResponseWriter, r *http.Request) {
				WriteBadRequest(w, r, "test error")
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrCodeBadRequest,
		},
		{
			name: "WriteValidationError",
			writeFunc: func(w http.ResponseWriter, r *http.Request) {
				WriteValidationError(w, r, "test error")
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrCodeValidation,
		},
		{
			name: "WriteNotFound",
			writeFunc: func(w http.ResponseWriter, r *http.Request) {
				WriteNotFound(w, r, "test error")
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   ErrCodeNotFound,
		},
		{
			name: "WriteInternalError",
			writeFunc: func(w http.ResponseWriter, r *http.Request) {
				WriteInternalError(w, r, "test error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   ErrCodeInternal,
		},
		{
			name: "WriteTimeout",
			writeFunc: func(w http.ResponseWriter, r *http.Request) {
				WriteTimeout(w, r)
			},
			expectedStatus: http.StatusGatewayTimeout,
			expectedCode:   ErrCodeTimeout,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req = req.WithContext(setRequestID(req.Context(), "test-id"))
			w := httptest.NewRecorder()

			tc.writeFunc(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			var resp ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if resp.Code != tc.expectedCode {
				t.Errorf("Expected code %q, got %q", tc.expectedCode, resp.Code)
			}
		})
	}
}

// TestContentType_Property tests that all error responses have JSON content type.
// Property: For any error response, the Content-Type header SHALL be "application/json".
func TestContentType_Property(t *testing.T) {
	property := func(message string) bool {
		if message == "" {
			return true
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		WriteError(w, req, http.StatusBadRequest, ErrCodeBadRequest, message)

		contentType := w.Header().Get("Content-Type")
		return contentType == "application/json"
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property failed: error responses should have JSON content type: %v", err)
	}
}

// setRequestID is a helper to set request ID in context for testing.
func setRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

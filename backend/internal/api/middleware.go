// Package api provides HTTP handlers and middleware for the API.
package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// RequestIDKey is the context key for the request ID.
	RequestIDKey contextKey = "requestID"
	// RequestIDHeader is the HTTP header name for the request ID.
	RequestIDHeader = "X-Request-ID"
)

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// generateRequestID creates a new unique request ID.
func generateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based ID if random fails
		return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000")))
	}
	return hex.EncodeToString(b)
}

// RequestIDMiddleware adds a unique request ID to each request.
// The ID is stored in the request context and added to the response header.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request already has an ID (from upstream proxy)
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Add to response header
		w.Header().Set(RequestIDHeader, requestID)

		// Add to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}

// LoggingMiddleware logs requests with timing and status.
// It logs security-relevant events without logging sensitive data.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := GetRequestID(r.Context())

		// Wrap response writer to capture status code
		rw := newResponseWriter(w)

		// Process request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Log the request (without sensitive data like auth headers or body)
		log.Printf("[%s] %s %s %d %v %s",
			requestID,
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
			r.RemoteAddr,
		)

		// Log security-relevant events
		if rw.statusCode == http.StatusTooManyRequests {
			log.Printf("[%s] SECURITY: Rate limit hit from %s on %s",
				requestID, r.RemoteAddr, r.URL.Path)
		}
		if rw.statusCode == http.StatusBadRequest {
			log.Printf("[%s] SECURITY: Validation failure from %s on %s",
				requestID, r.RemoteAddr, r.URL.Path)
		}
	})
}

// RecoveryMiddleware recovers from panics and returns 500.
// It logs the panic with stack trace for debugging.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(r.Context())

				// Log the panic with stack trace
				log.Printf("[%s] PANIC: %v\n%s", requestID, err, debug.Stack())

				// Return 500 error
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"Internal server error","code":"SERVER_INTERNAL"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Chain applies middleware in order (first middleware wraps outermost).
// Usage: Chain(handler, middleware1, middleware2, middleware3)
// Results in: middleware1(middleware2(middleware3(handler)))
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

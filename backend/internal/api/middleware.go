// Package api provides HTTP handlers and middleware for the API.
package api

import (
	"better-kiro-prompts/internal/logger"
	"context"
	"log/slog"
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
	// Fallback to logger's context key
	return logger.GetRequestID(ctx)
}

// RequestIDMiddleware adds a unique request ID to each request.
// The ID is stored in the request context and added to the response header.
// It uses the logger context helpers to ensure request ID propagates to all downstream operations.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request already has an ID (from upstream proxy)
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = logger.GenerateRequestID()
		}

		// Add to response header
		w.Header().Set(RequestIDHeader, requestID)

		// Add to context using both the API package key and logger package key
		// This ensures compatibility with existing code and enables logger context propagation
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		ctx = logger.WithRequestID(ctx, requestID)
		ctx = logger.WithUserIP(ctx, r.RemoteAddr)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// responseWriter wraps http.ResponseWriter to capture status code and bytes written.
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
	written      bool
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
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += int64(n)
	return n, err
}

// LoggingMiddleware logs requests with timing and status.
// It logs security-relevant events without logging sensitive data.
func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := GetRequestID(r.Context())

			// Log request start with all required fields
			log.HTTP().Info("request_start",
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.Int64("content_length", r.ContentLength),
			)

			// Wrap response writer to capture status code and bytes written
			rw := newResponseWriter(w)

			// Process request
			next.ServeHTTP(rw, r)

			// Calculate duration
			duration := time.Since(start)

			// Log request completion
			log.HTTP().Info("request_complete",
				slog.String("request_id", requestID),
				slog.Int("status", rw.statusCode),
				slog.Duration("duration", duration),
				slog.Int64("bytes_written", rw.bytesWritten),
			)

			// Log security-relevant events
			if rw.statusCode == http.StatusTooManyRequests {
				log.HTTP().Warn("security_rate_limit",
					slog.String("request_id", requestID),
					slog.String("remote_addr", r.RemoteAddr),
					slog.String("path", r.URL.Path),
				)
			}
			if rw.statusCode == http.StatusBadRequest {
				log.HTTP().Warn("security_validation_failure",
					slog.String("request_id", requestID),
					slog.String("remote_addr", r.RemoteAddr),
					slog.String("path", r.URL.Path),
				)
			}
		})
	}
}

// RecoveryMiddleware recovers from panics and returns 500.
// It logs the panic with stack trace for debugging.
func RecoveryMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestID := GetRequestID(r.Context())

					// Log the panic with stack trace
					log.HTTP().Error("panic_recovered",
						slog.String("request_id", requestID),
						slog.Any("error", err),
						slog.String("stack", string(debug.Stack())),
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
					)

					// Return 500 error
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"Internal server error","code":"SERVER_INTERNAL"}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
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

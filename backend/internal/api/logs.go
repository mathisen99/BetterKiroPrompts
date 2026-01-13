package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"better-kiro-prompts/internal/logger"
)

// ClientLogRequest represents a batch of client logs
type ClientLogRequest struct {
	Logs []ClientLogEntry `json:"logs"`
}

// ClientLogEntry represents a single log entry from the frontend
type ClientLogEntry struct {
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Stack     string    `json:"stack,omitempty"`
	URL       string    `json:"url"`
	Component string    `json:"component,omitempty"`
	UserAgent string    `json:"user_agent"`
	Timestamp time.Time `json:"timestamp"`
}

// HandleClientLogs receives and logs frontend errors
func HandleClientLogs(log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ClientLogRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// Process each log entry
		for _, entry := range req.Logs {
			level := parseLogLevel(entry.Level)

			// Build attributes for the log entry
			attrs := []any{
				slog.String("url", entry.URL),
				slog.String("user_agent", entry.UserAgent),
				slog.Time("client_time", entry.Timestamp),
			}

			// Add optional fields if present
			if entry.Stack != "" {
				attrs = append(attrs, slog.String("stack", entry.Stack))
			}
			if entry.Component != "" {
				attrs = append(attrs, slog.String("client_component", entry.Component))
			}

			log.Client().Log(r.Context(), level, entry.Message, attrs...)
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

// parseLogLevel converts a string level to slog.Level
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug", "DEBUG":
		return slog.LevelDebug
	case "info", "INFO":
		return slog.LevelInfo
	case "warn", "WARN", "warning", "WARNING":
		return slog.LevelWarn
	case "error", "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"better-kiro-prompts/internal/logger"
)

// LogLevelRequest represents a request to change the log level
type LogLevelRequest struct {
	Level string `json:"level"`
}

// LogLevelResponse represents the current log level
type LogLevelResponse struct {
	Level string `json:"level"`
}

// HandleGetLogLevel returns the current log level
func HandleGetLogLevel(log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		level := log.GetLevel()
		resp := LogLevelResponse{
			Level: levelToString(level),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.App().Error("failed to encode log level response",
				slog.String("error", err.Error()),
			)
		}
	}
}

// HandleSetLogLevel changes the log level at runtime
func HandleSetLogLevel(log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LogLevelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteBadRequest(w, r, "invalid request body")
			return
		}

		// Validate the level
		newLevel := logger.ParseLevel(req.Level)
		oldLevel := log.GetLevel()

		// Set the new level
		log.SetLevel(newLevel)

		log.App().Info("log_level_changed",
			slog.String("old_level", levelToString(oldLevel)),
			slog.String("new_level", levelToString(newLevel)),
		)

		resp := LogLevelResponse{
			Level: levelToString(newLevel),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.App().Error("failed to encode log level response",
				slog.String("error", err.Error()),
			)
		}
	}
}

// levelToString converts a slog.Level to its string representation
func levelToString(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARN"
	case slog.LevelError:
		return "ERROR"
	default:
		return "INFO"
	}
}

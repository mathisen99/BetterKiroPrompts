package api

import (
	"net/http"

	"better-kiro-prompts/internal/generation"
	"better-kiro-prompts/internal/ratelimit"
)

// RouterConfig holds dependencies for the router.
type RouterConfig struct {
	GenerationService *generation.Service
	RateLimiter       *ratelimit.Limiter
}

// NewRouter creates a new HTTP router with all API routes.
func NewRouter(cfg *RouterConfig) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /api/health", HandleHealth)

	// Generation endpoints (if service is configured)
	if cfg != nil && cfg.GenerationService != nil && cfg.RateLimiter != nil {
		genHandler := NewGenerateHandler(cfg.GenerationService, cfg.RateLimiter)
		mux.HandleFunc("POST /api/generate/questions", genHandler.HandleGenerateQuestions)
		mux.HandleFunc("POST /api/generate/outputs", genHandler.HandleGenerateOutputs)
	}

	return mux
}

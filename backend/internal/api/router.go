package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	// Serve static files from ./static directory (SPA with fallback to index.html)
	staticDir := "./static"
	if _, err := os.Stat(staticDir); err == nil {
		mux.HandleFunc("/", spaHandler(staticDir))
	}

	return mux
}

// spaHandler serves static files and falls back to index.html for SPA routing.
func spaHandler(staticDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Don't serve static for API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		path := filepath.Join(staticDir, r.URL.Path)

		// Check if file exists
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			// Serve index.html for SPA routing
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		// Serve the actual file
		http.ServeFile(w, r, path)
	}
}

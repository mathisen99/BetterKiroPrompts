package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"better-kiro-prompts/internal/gallery"
	"better-kiro-prompts/internal/generation"
	"better-kiro-prompts/internal/logger"
	"better-kiro-prompts/internal/ratelimit"
	"better-kiro-prompts/internal/scanner"
)

// RouterConfig holds dependencies for the router.
type RouterConfig struct {
	GenerationService *generation.Service
	RateLimiter       *ratelimit.Limiter
	GalleryService    *gallery.Service
	RatingLimiter     *ratelimit.Limiter
	ScannerService    *scanner.Service
	ScanRateLimiter   *ratelimit.Limiter
	Logger            *logger.Logger
}

// NewRouter creates a new HTTP router with all API routes.
func NewRouter(cfg *RouterConfig) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /api/health", HandleHealth)

	// Generation endpoints (if service is configured)
	if cfg != nil && cfg.GenerationService != nil && cfg.RateLimiter != nil {
		genHandler := NewGenerateHandler(cfg.GenerationService, cfg.RateLimiter)
		mux.HandleFunc("POST /api/generate/questions", genHandler.HandleGenerateQuestions)
		mux.HandleFunc("POST /api/generate/outputs", genHandler.HandleGenerateOutputs)
	}

	// Gallery endpoints (if service is configured)
	if cfg != nil && cfg.GalleryService != nil {
		galleryHandler := NewGalleryHandler(cfg.GalleryService, cfg.RatingLimiter)
		mux.HandleFunc("GET /api/gallery", galleryHandler.HandleListGallery)
		mux.HandleFunc("GET /api/gallery/{id}", galleryHandler.HandleGetGalleryItem)
		mux.HandleFunc("POST /api/gallery/{id}/rate", galleryHandler.HandleRateGalleryItem)
	}

	// Scanner endpoints (if service is configured)
	if cfg != nil && cfg.ScannerService != nil && cfg.ScanRateLimiter != nil {
		scanHandler := NewScanHandler(cfg.ScannerService, cfg.ScanRateLimiter)
		mux.HandleFunc("POST /api/scan", scanHandler.HandleStartScan)
		mux.HandleFunc("GET /api/scan/config", scanHandler.HandleGetScanConfig)
		mux.HandleFunc("GET /api/scan/{id}", scanHandler.HandleGetScan)
	}

	// Serve static files from ./static directory (SPA with fallback to index.html)
	staticDir := "./static"
	if _, err := os.Stat(staticDir); err == nil {
		mux.HandleFunc("/", spaHandler(staticDir))
	}

	// Apply middleware chain: Recovery -> RequestID -> Logging
	// Order matters: Recovery is outermost to catch panics from all handlers
	// Logger is required for Recovery and Logging middleware
	if cfg != nil && cfg.Logger != nil {
		return Chain(mux,
			RecoveryMiddleware(cfg.Logger),
			RequestIDMiddleware,
			LoggingMiddleware(cfg.Logger),
		)
	}

	// Fallback without logging (for testing or when logger is not configured)
	return Chain(mux,
		RequestIDMiddleware,
	)
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

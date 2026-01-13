package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"better-kiro-prompts/internal/api"
	"better-kiro-prompts/internal/db"
	"better-kiro-prompts/internal/gallery"
	"better-kiro-prompts/internal/generation"
	"better-kiro-prompts/internal/openai"
	"better-kiro-prompts/internal/ratelimit"
	"better-kiro-prompts/internal/storage"
)

func main() {
	ctx := context.Background()

	if err := db.Connect(ctx); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize dependencies
	routerCfg := &api.RouterConfig{}

	// Initialize storage repository for gallery (only if DB is connected)
	if db.DB != nil {
		repo := storage.NewPostgresRepository(db.DB)

		// Initialize gallery service with rating limiter (20 ratings/hour per IP)
		ratingLimiter := ratelimit.NewLimiterWithConfig(20, time.Hour)
		galleryService := gallery.NewService(repo, ratingLimiter)
		routerCfg.GalleryService = galleryService
		routerCfg.RatingLimiter = ratingLimiter
		log.Printf("Gallery service initialized")
	} else {
		log.Printf("Warning: Database not connected, gallery endpoints will not be available")
	}

	// Try to create OpenAI client (optional - may not have API key in dev)
	openaiClient, err := openai.NewClient()
	if err != nil {
		log.Printf("Warning: OpenAI client not initialized: %v", err)
		log.Printf("Generation endpoints will not be available")
	} else {
		// Create generation service with repository for gallery storage
		var repo storage.Repository
		if db.DB != nil {
			repo = storage.NewPostgresRepository(db.DB)
		}
		genService := generation.NewServiceWithDeps(openaiClient, nil, repo)
		rateLimiter := ratelimit.NewLimiter()
		routerCfg.GenerationService = genService
		routerCfg.RateLimiter = rateLimiter
		log.Printf("Generation service initialized")
	}

	router := api.NewRouter(routerCfg)

	// Create HTTP server with explicit configuration
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Channel to listen for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	sig := <-shutdown
	log.Printf("Received signal %v, initiating graceful shutdown...", sig)

	// Create context with timeout for graceful shutdown
	// Allow up to 30 seconds for in-flight requests to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown error: %v", err)
	} else {
		log.Printf("Server gracefully stopped")
	}

	// Close database connection
	if err := db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	} else {
		log.Printf("Database connection closed")
	}
}

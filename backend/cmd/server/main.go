package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"better-kiro-prompts/internal/api"
	"better-kiro-prompts/internal/db"
	"better-kiro-prompts/internal/gallery"
	"better-kiro-prompts/internal/generation"
	"better-kiro-prompts/internal/logger"
	"better-kiro-prompts/internal/openai"
	"better-kiro-prompts/internal/ratelimit"
	"better-kiro-prompts/internal/scanner"
	"better-kiro-prompts/internal/storage"
)

const version = "1.0.0"

func main() {
	ctx := context.Background()

	// Initialize logger first
	logLevel := logger.ParseLevel(os.Getenv("LOG_LEVEL"))
	logCfg := logger.Config{
		Level:       logLevel,
		LogDir:      "./logs",
		MaxSizeMB:   100,
		MaxAgeDays:  7,
		EnableColor: true,
	}

	appLog, err := logger.New(logCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := appLog.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing logger: %v\n", err)
		}
	}()

	appLog.App().Info("application_starting",
		slog.String("version", version),
		slog.String("log_level", logLevel.String()),
		slog.String("log_dir", logCfg.LogDir),
	)

	// Database connection
	appLog.App().Info("database_connecting")
	if err := db.Connect(ctx); err != nil {
		appLog.App().Error("database_connection_failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	appLog.App().Info("database_connected")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize dependencies
	routerCfg := &api.RouterConfig{
		Logger: appLog,
	}

	// Initialize storage repository for gallery (only if DB is connected)
	if db.DB != nil {
		repo := storage.NewPostgresRepository(db.DB)

		// Initialize gallery service with rating limiter (20 ratings/hour per IP)
		ratingLimiter := ratelimit.NewLimiterWithConfig(20, time.Hour)
		galleryService := gallery.NewService(repo, ratingLimiter, appLog)
		routerCfg.GalleryService = galleryService
		routerCfg.RatingLimiter = ratingLimiter
		appLog.App().Info("gallery_service_initialized")
	} else {
		appLog.App().Warn("gallery_service_unavailable",
			slog.String("reason", "database not connected"))
	}

	// Try to create OpenAI client (optional - may not have API key in dev)
	openaiClient, err := openai.NewClient(appLog.App())
	if err != nil {
		appLog.App().Warn("openai_client_unavailable",
			slog.String("error", err.Error()),
			slog.String("impact", "generation endpoints will not be available"))
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
		appLog.App().Info("generation_service_initialized")
	}

	// Initialize scanner service (requires DB, OpenAI client is optional for AI review)
	if db.DB != nil {
		githubToken := os.Getenv("GITHUB_TOKEN")

		// Set scanner container name from environment
		scannerContainer := os.Getenv("SCANNER_CONTAINER")
		if scannerContainer != "" {
			scanner.SetScannerContainer(scannerContainer)
		}

		scannerService := scanner.NewService(db.DB, openaiClient, githubToken,
			scanner.WithServiceLogger(appLog.Scanner()))
		// Scanner rate limiter: 10 scans per hour per IP (scans are resource-intensive)
		scanRateLimiter := ratelimit.NewLimiterWithConfig(10, time.Hour)
		routerCfg.ScannerService = scannerService
		routerCfg.ScanRateLimiter = scanRateLimiter

		appLog.App().Info("scanner_service_initialized",
			slog.Bool("private_repo_support", githubToken != ""))
	} else {
		appLog.App().Warn("scanner_service_unavailable",
			slog.String("reason", "database not connected"))
	}

	appLog.App().Info("services_initialized",
		slog.Bool("generation_enabled", routerCfg.GenerationService != nil),
		slog.Bool("gallery_enabled", routerCfg.GalleryService != nil),
		slog.Bool("scanner_enabled", routerCfg.ScannerService != nil),
	)

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
		appLog.App().Info("server_starting", slog.String("port", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLog.App().Error("server_error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	sig := <-shutdown
	appLog.App().Info("shutdown_signal_received", slog.String("signal", sig.String()))

	// Create context with timeout for graceful shutdown
	// Allow up to 30 seconds for in-flight requests to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		appLog.App().Error("shutdown_error", slog.String("error", err.Error()))
	} else {
		appLog.App().Info("server_stopped_gracefully")
	}

	// Close database connection
	if err := db.Close(); err != nil {
		appLog.App().Error("database_close_error", slog.String("error", err.Error()))
	} else {
		appLog.App().Info("database_connection_closed")
	}
}

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
	"better-kiro-prompts/internal/config"
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

	// Load configuration first (before logger, as logger config comes from here)
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger using config values directly
	appLog, err := logger.NewFromLoggingConfig(cfg.Logging)
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
		slog.String("log_level", cfg.Logging.Level),
		slog.String("log_dir", cfg.Logging.Directory),
	)

	// Log loaded configuration (with sensitive values redacted)
	cfg.LogConfig(appLog.App())

	// Database connection
	appLog.App().Info("database_connecting")
	db.SetLogger(appLog.DB()) // Set logger for database operations
	if err := db.Connect(ctx); err != nil {
		appLog.App().Error("database_connection_failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	appLog.App().Info("database_connected")

	// Use port from config (already includes env var override)
	port := fmt.Sprintf("%d", cfg.Server.Port)

	// Initialize dependencies
	routerCfg := &api.RouterConfig{
		Logger: appLog,
	}

	// Initialize storage repository for gallery (only if DB is connected)
	var loggingDB *db.LoggingDB
	if db.DB != nil {
		loggingDB = db.NewLoggingDB(db.DB, appLog.DB())
		repo := storage.NewPostgresRepositoryWithLogging(loggingDB)

		// Initialize gallery service with rating limiter using config values
		ratingLimiter := ratelimit.NewLimiterWithConfigAndLogger(cfg.RateLimit.RatingLimitPerHour, time.Hour, appLog.App())
		galleryService := gallery.NewServiceWithConfig(repo, ratingLimiter, appLog, cfg.Gallery)
		routerCfg.GalleryService = galleryService
		routerCfg.RatingLimiter = ratingLimiter
		appLog.App().Info("gallery_service_initialized",
			slog.Int("page_size", cfg.Gallery.PageSize),
			slog.String("default_sort", cfg.Gallery.DefaultSort),
		)
	} else {
		appLog.App().Warn("gallery_service_unavailable",
			slog.String("reason", "database not connected"))
	}

	// Try to create OpenAI client (optional - may not have API key in dev)
	// Use config values for model, timeout, reasoning effort, and verbosity
	openaiClient, err := openai.NewClientWithConfig(openai.ClientConfig{
		APIKey:          os.Getenv("OPENAI_API_KEY"),
		BaseURL:         cfg.OpenAI.BaseURL,
		Model:           cfg.OpenAI.Model,
		Timeout:         cfg.OpenAI.Timeout.Duration(),
		ReasoningEffort: openai.ReasoningEffort(cfg.OpenAI.ReasoningEffort),
		Verbosity:       openai.Verbosity(cfg.OpenAI.Verbosity),
		Logger:          appLog.App(),
	})
	if err != nil {
		appLog.App().Warn("openai_client_unavailable",
			slog.String("error", err.Error()),
			slog.String("impact", "generation endpoints will not be available"))
	} else {
		// Create generation service with repository for gallery storage and config
		var repo storage.Repository
		if loggingDB != nil {
			repo = storage.NewPostgresRepositoryWithLogging(loggingDB)
		}
		genService := generation.NewServiceWithConfig(openaiClient, nil, repo, appLog.App(), cfg.Generation)
		// Use generation rate limit from config
		rateLimiter := ratelimit.NewLimiterWithConfigAndLogger(cfg.RateLimit.GenerationLimitPerHour, time.Hour, appLog.App())
		routerCfg.GenerationService = genService
		routerCfg.RateLimiter = rateLimiter
		appLog.App().Info("generation_service_initialized",
			slog.Int("max_project_idea_length", cfg.Generation.MaxProjectIdeaLength),
			slog.Int("max_answer_length", cfg.Generation.MaxAnswerLength),
			slog.Int("min_questions", cfg.Generation.MinQuestions),
			slog.Int("max_questions", cfg.Generation.MaxQuestions),
			slog.Int("max_retries", cfg.Generation.MaxRetries),
		)
	}

	// Initialize scanner service (requires DB, OpenAI client is optional for AI review)
	if db.DB != nil {
		githubToken := os.Getenv("GITHUB_TOKEN")

		// Set scanner container name from environment
		scannerContainer := os.Getenv("SCANNER_CONTAINER")
		if scannerContainer != "" {
			scanner.SetScannerContainer(scannerContainer)
		}

		// Use NewServiceWithConfig to pass scanner configuration
		scannerService := scanner.NewServiceWithConfig(db.DB, openaiClient, githubToken, cfg.Scanner, cfg.OpenAI.CodeReviewModel,
			scanner.WithServiceLogger(appLog.Scanner()))
		// Scanner rate limiter using config values
		scanRateLimiter := ratelimit.NewLimiterWithConfigAndLogger(cfg.RateLimit.ScanLimitPerHour, time.Hour, appLog.App())
		routerCfg.ScannerService = scannerService
		routerCfg.ScanRateLimiter = scanRateLimiter

		appLog.App().Info("scanner_service_initialized",
			slog.Bool("private_repo_support", githubToken != ""),
			slog.Int("max_repo_size_mb", cfg.Scanner.MaxRepoSizeMB),
			slog.Int("max_review_files", cfg.Scanner.MaxReviewFiles),
			slog.Int("retention_days", cfg.Scanner.RetentionDays),
			slog.Int("tool_timeout_seconds", cfg.Scanner.ToolTimeoutSeconds),
		)
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

	// Create context with timeout for graceful shutdown using config value
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout.Duration())
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

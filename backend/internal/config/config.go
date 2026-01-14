// Package config provides centralized configuration loading and validation.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// Config holds all application configuration.
type Config struct {
	Server     ServerConfig     `toml:"server"`
	OpenAI     OpenAIConfig     `toml:"openai"`
	RateLimit  RateLimitConfig  `toml:"rate_limit"`
	Logging    LoggingConfig    `toml:"logging"`
	Scanner    ScannerConfig    `toml:"scanner"`
	Generation GenerationConfig `toml:"generation"`
	Gallery    GalleryConfig    `toml:"gallery"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port            int      `toml:"port"`
	Host            string   `toml:"host"`
	ShutdownTimeout Duration `toml:"shutdown_timeout"`
}

// OpenAIConfig holds OpenAI API settings.
type OpenAIConfig struct {
	Model           string   `toml:"model"`
	CodeReviewModel string   `toml:"code_review_model"`
	BaseURL         string   `toml:"base_url"`
	Timeout         Duration `toml:"timeout"`
	ReasoningEffort string   `toml:"reasoning_effort"`
	Verbosity       string   `toml:"verbosity"`
}

// RateLimitConfig holds rate limiting settings.
type RateLimitConfig struct {
	GenerationLimitPerHour int `toml:"generation_limit_per_hour"`
	RatingLimitPerHour     int `toml:"rating_limit_per_hour"`
	ScanLimitPerHour       int `toml:"scan_limit_per_hour"`
}

// LoggingConfig holds logging settings.
type LoggingConfig struct {
	Level       string `toml:"level"`
	Directory   string `toml:"directory"`
	MaxSizeMB   int    `toml:"max_size_mb"`
	MaxAgeDays  int    `toml:"max_age_days"`
	EnableColor bool   `toml:"enable_color"`
}

// ScannerConfig holds security scanner settings.
type ScannerConfig struct {
	MaxRepoSizeMB      int      `toml:"max_repo_size_mb"`
	MaxReviewFiles     int      `toml:"max_review_files"`
	ToolTimeoutSeconds int      `toml:"tool_timeout_seconds"`
	RetentionDays      int      `toml:"retention_days"`
	CloneTimeout       Duration `toml:"clone_timeout"`
}

// GenerationConfig holds AI generation settings.
type GenerationConfig struct {
	MaxProjectIdeaLength int `toml:"max_project_idea_length"`
	MaxAnswerLength      int `toml:"max_answer_length"`
	MinQuestions         int `toml:"min_questions"`
	MaxQuestions         int `toml:"max_questions"`
	MaxRetries           int `toml:"max_retries"`
}

// GalleryConfig holds gallery settings.
type GalleryConfig struct {
	PageSize    int    `toml:"page_size"`
	DefaultSort string `toml:"default_sort"`
}

// Duration is a wrapper around time.Duration that supports TOML unmarshaling.
type Duration time.Duration

// UnmarshalText implements encoding.TextUnmarshaler for Duration.
func (d *Duration) UnmarshalText(text []byte) error {
	duration, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}

// MarshalText implements encoding.TextMarshaler for Duration.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

// Duration returns the underlying time.Duration.
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

// DefaultConfig returns configuration with sensible defaults matching current hardcoded values.
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            8080,
			Host:            "0.0.0.0",
			ShutdownTimeout: Duration(30 * time.Second),
		},
		OpenAI: OpenAIConfig{
			Model:           "gpt-5.2",
			CodeReviewModel: "gpt-5.1-codex-max",
			BaseURL:         "https://api.openai.com/v1",
			Timeout:         Duration(180 * time.Second),
			ReasoningEffort: "medium",
			Verbosity:       "medium",
		},
		RateLimit: RateLimitConfig{
			GenerationLimitPerHour: 10,
			RatingLimitPerHour:     20,
			ScanLimitPerHour:       10,
		},
		Logging: LoggingConfig{
			Level:       "INFO",
			Directory:   "./logs",
			MaxSizeMB:   100,
			MaxAgeDays:  7,
			EnableColor: true,
		},
		Scanner: ScannerConfig{
			MaxRepoSizeMB:      500,
			MaxReviewFiles:     10,
			ToolTimeoutSeconds: 300,
			RetentionDays:      7,
			CloneTimeout:       Duration(5 * time.Minute),
		},
		Generation: GenerationConfig{
			MaxProjectIdeaLength: 2000,
			MaxAnswerLength:      1000,
			MinQuestions:         5,
			MaxQuestions:         10,
			MaxRetries:           1,
		},
		Gallery: GalleryConfig{
			PageSize:    20,
			DefaultSort: "newest",
		},
	}
}

// Load reads configuration with the following precedence:
// 1. Default values
// 2. config.toml file values
// 3. Environment variable overrides
func Load() (*Config, error) {
	// Determine config path
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.toml"
	}

	return LoadFromPath(configPath)
}

// LoadFromPath reads configuration from a specific path.
func LoadFromPath(path string) (*Config, error) {
	cfg := DefaultConfig()

	// Load from file if exists
	if _, err := os.Stat(path); err == nil {
		if _, err := toml.DecodeFile(path, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Apply environment overrides
	cfg.ApplyEnvironmentOverrides()

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// ApplyEnvironmentOverrides reads env vars and overrides config values.
// Environment variables take precedence for secrets and backward compatibility.
func (c *Config) ApplyEnvironmentOverrides() {
	// Server
	if v := os.Getenv("PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.Server.Port = port
		}
	}

	// Logging
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		c.Logging.Level = v
	}

	// OpenAI model override
	if v := os.Getenv("OPENAI_MODEL"); v != "" {
		c.OpenAI.Model = v
	}

	// Scanner overrides (existing env vars for backward compatibility)
	if v := os.Getenv("SCANNER_MAX_REPO_SIZE_MB"); v != "" {
		if size, err := strconv.Atoi(v); err == nil {
			c.Scanner.MaxRepoSizeMB = size
		}
	}
	if v := os.Getenv("SCANNER_MAX_REVIEW_FILES"); v != "" {
		if files, err := strconv.Atoi(v); err == nil {
			c.Scanner.MaxReviewFiles = files
		}
	}
	if v := os.Getenv("SCANNER_TOOL_TIMEOUT_SECONDS"); v != "" {
		if timeout, err := strconv.Atoi(v); err == nil {
			c.Scanner.ToolTimeoutSeconds = timeout
		}
	}
	if v := os.Getenv("SCANNER_RESULT_RETENTION_DAYS"); v != "" {
		if days, err := strconv.Atoi(v); err == nil {
			c.Scanner.RetentionDays = days
		}
	}

	// Rate limit overrides
	if v := os.Getenv("RATE_LIMIT_GENERATION"); v != "" {
		if limit, err := strconv.Atoi(v); err == nil {
			c.RateLimit.GenerationLimitPerHour = limit
		}
	}
	if v := os.Getenv("RATE_LIMIT_RATING"); v != "" {
		if limit, err := strconv.Atoi(v); err == nil {
			c.RateLimit.RatingLimitPerHour = limit
		}
	}
	if v := os.Getenv("RATE_LIMIT_SCAN"); v != "" {
		if limit, err := strconv.Atoi(v); err == nil {
			c.RateLimit.ScanLimitPerHour = limit
		}
	}
}

// Valid values for enum fields
var (
	validReasoningEfforts = map[string]bool{
		"none": true, "low": true, "medium": true, "high": true, "xhigh": true,
	}
	validVerbosities = map[string]bool{
		"low": true, "medium": true, "high": true,
	}
	validLogLevels = map[string]bool{
		"DEBUG": true, "INFO": true, "WARN": true, "ERROR": true,
		"debug": true, "info": true, "warn": true, "error": true,
	}
	validSortOptions = map[string]bool{
		"newest": true, "highest_rated": true, "most_viewed": true,
	}
)

// Validate checks all configuration values are within acceptable ranges.
func (c *Config) Validate() error {
	var errs []string

	// Server validation
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		errs = append(errs, fmt.Sprintf("server.port must be 1-65535, got %d", c.Server.Port))
	}
	if c.Server.ShutdownTimeout.Duration() < time.Second {
		errs = append(errs, "server.shutdown_timeout must be at least 1s")
	}

	// OpenAI validation
	if c.OpenAI.Model == "" {
		errs = append(errs, "openai.model is required")
	}
	if !validReasoningEfforts[c.OpenAI.ReasoningEffort] {
		errs = append(errs, fmt.Sprintf("openai.reasoning_effort must be one of: none, low, medium, high, xhigh; got %s", c.OpenAI.ReasoningEffort))
	}
	if !validVerbosities[c.OpenAI.Verbosity] {
		errs = append(errs, fmt.Sprintf("openai.verbosity must be one of: low, medium, high; got %s", c.OpenAI.Verbosity))
	}
	if c.OpenAI.Timeout.Duration() < 10*time.Second {
		errs = append(errs, "openai.timeout must be at least 10s")
	}

	// Rate limit validation
	if c.RateLimit.GenerationLimitPerHour < 1 {
		errs = append(errs, "rate_limit.generation_limit_per_hour must be at least 1")
	}
	if c.RateLimit.RatingLimitPerHour < 1 {
		errs = append(errs, "rate_limit.rating_limit_per_hour must be at least 1")
	}
	if c.RateLimit.ScanLimitPerHour < 1 {
		errs = append(errs, "rate_limit.scan_limit_per_hour must be at least 1")
	}

	// Logging validation
	if !validLogLevels[c.Logging.Level] {
		errs = append(errs, fmt.Sprintf("logging.level must be one of: DEBUG, INFO, WARN, ERROR; got %s", c.Logging.Level))
	}
	if c.Logging.MaxSizeMB < 1 {
		errs = append(errs, "logging.max_size_mb must be at least 1")
	}
	if c.Logging.MaxAgeDays < 1 {
		errs = append(errs, "logging.max_age_days must be at least 1")
	}

	// Scanner validation
	if c.Scanner.MaxRepoSizeMB < 1 {
		errs = append(errs, "scanner.max_repo_size_mb must be at least 1")
	}
	if c.Scanner.MaxReviewFiles < 1 {
		errs = append(errs, "scanner.max_review_files must be at least 1")
	}
	if c.Scanner.ToolTimeoutSeconds < 10 {
		errs = append(errs, "scanner.tool_timeout_seconds must be at least 10")
	}
	if c.Scanner.RetentionDays < 1 {
		errs = append(errs, "scanner.retention_days must be at least 1")
	}
	if c.Scanner.CloneTimeout.Duration() < 10*time.Second {
		errs = append(errs, "scanner.clone_timeout must be at least 10s")
	}

	// Generation validation
	if c.Generation.MaxProjectIdeaLength < 100 {
		errs = append(errs, "generation.max_project_idea_length must be at least 100")
	}
	if c.Generation.MaxAnswerLength < 100 {
		errs = append(errs, "generation.max_answer_length must be at least 100")
	}
	if c.Generation.MinQuestions < 1 {
		errs = append(errs, "generation.min_questions must be at least 1")
	}
	if c.Generation.MaxQuestions < c.Generation.MinQuestions {
		errs = append(errs, "generation.max_questions must be >= min_questions")
	}
	if c.Generation.MaxRetries < 0 {
		errs = append(errs, "generation.max_retries must be at least 0")
	}

	// Gallery validation
	if c.Gallery.PageSize < 1 || c.Gallery.PageSize > 100 {
		errs = append(errs, "gallery.page_size must be 1-100")
	}
	if !validSortOptions[c.Gallery.DefaultSort] {
		errs = append(errs, fmt.Sprintf("gallery.default_sort must be one of: newest, highest_rated, most_viewed; got %s", c.Gallery.DefaultSort))
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation errors:\n  - %s", strings.Join(errs, "\n  - "))
	}

	return nil
}

// LogConfig logs the loaded configuration with sensitive values redacted.
func (c *Config) LogConfig(log *slog.Logger) {
	log.Info("configuration_loaded",
		slog.Group("server",
			slog.Int("port", c.Server.Port),
			slog.String("host", c.Server.Host),
			slog.Duration("shutdown_timeout", c.Server.ShutdownTimeout.Duration()),
		),
		slog.Group("openai",
			slog.String("model", c.OpenAI.Model),
			slog.String("code_review_model", c.OpenAI.CodeReviewModel),
			slog.String("base_url", c.OpenAI.BaseURL),
			slog.Duration("timeout", c.OpenAI.Timeout.Duration()),
			slog.String("reasoning_effort", c.OpenAI.ReasoningEffort),
			slog.String("verbosity", c.OpenAI.Verbosity),
		),
		slog.Group("rate_limit",
			slog.Int("generation_per_hour", c.RateLimit.GenerationLimitPerHour),
			slog.Int("rating_per_hour", c.RateLimit.RatingLimitPerHour),
			slog.Int("scan_per_hour", c.RateLimit.ScanLimitPerHour),
		),
		slog.Group("logging",
			slog.String("level", c.Logging.Level),
			slog.String("directory", c.Logging.Directory),
			slog.Int("max_size_mb", c.Logging.MaxSizeMB),
			slog.Int("max_age_days", c.Logging.MaxAgeDays),
			slog.Bool("enable_color", c.Logging.EnableColor),
		),
		slog.Group("scanner",
			slog.Int("max_repo_size_mb", c.Scanner.MaxRepoSizeMB),
			slog.Int("max_review_files", c.Scanner.MaxReviewFiles),
			slog.Int("tool_timeout_seconds", c.Scanner.ToolTimeoutSeconds),
			slog.Int("retention_days", c.Scanner.RetentionDays),
			slog.Duration("clone_timeout", c.Scanner.CloneTimeout.Duration()),
		),
		slog.Group("generation",
			slog.Int("max_project_idea_length", c.Generation.MaxProjectIdeaLength),
			slog.Int("max_answer_length", c.Generation.MaxAnswerLength),
			slog.Int("min_questions", c.Generation.MinQuestions),
			slog.Int("max_questions", c.Generation.MaxQuestions),
			slog.Int("max_retries", c.Generation.MaxRetries),
		),
		slog.Group("gallery",
			slog.Int("page_size", c.Gallery.PageSize),
			slog.String("default_sort", c.Gallery.DefaultSort),
		),
	)
}

// String returns a string representation of the config for debugging.
// Sensitive values are redacted.
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{Server: {Port: %d, Host: %s}, OpenAI: {Model: %s}, RateLimit: {Gen: %d/h, Rating: %d/h, Scan: %d/h}, Logging: {Level: %s}}",
		c.Server.Port, c.Server.Host, c.OpenAI.Model,
		c.RateLimit.GenerationLimitPerHour, c.RateLimit.RatingLimitPerHour, c.RateLimit.ScanLimitPerHour,
		c.Logging.Level,
	)
}

// Package logger provides structured logging with file rotation and colored console output.
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Level represents log severity
type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// Config holds logger configuration
type Config struct {
	Level       Level
	LogDir      string
	MaxSizeMB   int
	MaxAgeDays  int
	EnableColor bool
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() Config {
	return Config{
		Level:       LevelInfo,
		LogDir:      "./logs",
		MaxSizeMB:   100,
		MaxAgeDays:  7,
		EnableColor: true,
	}
}

// Logger wraps slog with file rotation and colored output
type Logger struct {
	config   Config
	handlers map[string]*slog.Logger
	files    map[string]*RotatingFile
	mu       sync.RWMutex
	levelVar *slog.LevelVar
}

// New creates a new logger instance
func New(cfg Config) (*Logger, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	l := &Logger{
		config:   cfg,
		handlers: make(map[string]*slog.Logger),
		files:    make(map[string]*RotatingFile),
		levelVar: new(slog.LevelVar),
	}
	l.levelVar.Set(cfg.Level)

	// Initialize category-specific loggers
	categories := []string{"app", "http", "db", "scanner", "client"}
	for _, cat := range categories {
		if err := l.initCategory(cat); err != nil {
			_ = l.Close()
			return nil, fmt.Errorf("failed to initialize %s logger: %w", cat, err)
		}
	}

	return l, nil
}

// initCategory initializes a logger for a specific category
func (l *Logger) initCategory(category string) error {
	// Create rotating file for this category
	filename := fmt.Sprintf("%s-%s.log", time.Now().Format("2006-01-02"), category)
	filePath := filepath.Join(l.config.LogDir, filename)

	rf, err := NewRotatingFile(filePath, int64(l.config.MaxSizeMB)*1024*1024, l.config.MaxAgeDays)
	if err != nil {
		return err
	}
	l.files[category] = rf

	// Create multi-writer for file and console
	var writers []io.Writer
	writers = append(writers, rf)

	// Add console output with optional colors
	if l.config.EnableColor {
		writers = append(writers, NewColorWriter(os.Stdout, l.levelVar))
	} else {
		writers = append(writers, os.Stdout)
	}

	multiWriter := io.MultiWriter(writers...)

	// Create JSON handler for file output
	opts := &slog.HandlerOptions{
		Level: l.levelVar,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Redact sensitive data
			return redactAttr(a)
		},
	}

	var handler slog.Handler
	if l.config.EnableColor && isTerminal(os.Stdout) {
		// Use color handler for console, JSON for file
		handler = NewColorHandler(multiWriter, opts, category)
	} else {
		handler = slog.NewJSONHandler(multiWriter, opts)
	}

	l.handlers[category] = slog.New(handler).With(slog.String("component", category))

	return nil
}

// App returns the application logger
func (l *Logger) App() *slog.Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.handlers["app"]
}

// HTTP returns the HTTP request logger
func (l *Logger) HTTP() *slog.Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.handlers["http"]
}

// DB returns the database logger
func (l *Logger) DB() *slog.Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.handlers["db"]
}

// Scanner returns the scanner logger
func (l *Logger) Scanner() *slog.Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.handlers["scanner"]
}

// Client returns the client error logger
func (l *Logger) Client() *slog.Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.handlers["client"]
}

// SetLevel changes the log level at runtime
func (l *Logger) SetLevel(level Level) {
	l.levelVar.Set(level)
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() Level {
	return l.levelVar.Level()
}

// Close closes all log files
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var errs []error
	for name, rf := range l.files {
		if err := rf.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close %s log: %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing log files: %v", errs)
	}
	return nil
}

// ParseLevel parses a string log level
func ParseLevel(s string) Level {
	switch s {
	case "DEBUG", "debug":
		return LevelDebug
	case "INFO", "info", "":
		return LevelInfo
	case "WARN", "warn", "WARNING", "warning":
		return LevelWarn
	case "ERROR", "error":
		return LevelError
	default:
		return LevelInfo
	}
}

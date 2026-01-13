package logger

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	// Create temp directory for logs
	tmpDir := t.TempDir()

	cfg := Config{
		Level:       LevelInfo,
		LogDir:      tmpDir,
		MaxSizeMB:   1,
		MaxAgeDays:  1,
		EnableColor: false,
	}

	log, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer func() { _ = log.Close() }()

	// Verify log directory was created
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("log directory was not created")
	}

	// Verify all category loggers exist
	if log.App() == nil {
		t.Error("App logger is nil")
	}
	if log.HTTP() == nil {
		t.Error("HTTP logger is nil")
	}
	if log.DB() == nil {
		t.Error("DB logger is nil")
	}
	if log.Scanner() == nil {
		t.Error("Scanner logger is nil")
	}
	if log.Client() == nil {
		t.Error("Client logger is nil")
	}
}

func TestSetLevel(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := Config{
		Level:       LevelInfo,
		LogDir:      tmpDir,
		MaxSizeMB:   1,
		MaxAgeDays:  1,
		EnableColor: false,
	}

	log, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer func() { _ = log.Close() }()

	// Verify initial level
	if log.GetLevel() != LevelInfo {
		t.Errorf("expected initial level INFO, got %v", log.GetLevel())
	}

	// Change level
	log.SetLevel(LevelDebug)
	if log.GetLevel() != LevelDebug {
		t.Errorf("expected level DEBUG after SetLevel, got %v", log.GetLevel())
	}

	log.SetLevel(LevelError)
	if log.GetLevel() != LevelError {
		t.Errorf("expected level ERROR after SetLevel, got %v", log.GetLevel())
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"DEBUG", LevelDebug},
		{"debug", LevelDebug},
		{"INFO", LevelInfo},
		{"info", LevelInfo},
		{"", LevelInfo}, // default
		{"WARN", LevelWarn},
		{"warn", LevelWarn},
		{"WARNING", LevelWarn},
		{"ERROR", LevelError},
		{"error", LevelError},
		{"unknown", LevelInfo}, // default for unknown
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLogWritesToFile(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := Config{
		Level:       LevelInfo,
		LogDir:      tmpDir,
		MaxSizeMB:   1,
		MaxAgeDays:  1,
		EnableColor: false,
	}

	log, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Write a log entry
	log.App().Info("test message", "key", "value")

	// Close to flush
	if err := log.Close(); err != nil {
		t.Fatalf("failed to close logger: %v", err)
	}

	// Check that log files were created
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read log directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("no log files were created")
	}

	// Find the app log file
	var appLogFound bool
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".log" {
			appLogFound = true
			break
		}
	}

	if !appLogFound {
		t.Error("app log file was not created")
	}
}

func TestContextHelpers(t *testing.T) {
	ctx := context.Background()

	// Test WithRequestID and GetRequestID
	requestID := "test-request-123"
	ctx = WithRequestID(ctx, requestID)
	if got := GetRequestID(ctx); got != requestID {
		t.Errorf("GetRequestID() = %q, want %q", got, requestID)
	}

	// Test WithComponent and GetComponent
	component := "test-component"
	ctx = WithComponent(ctx, component)
	if got := GetComponent(ctx); got != component {
		t.Errorf("GetComponent() = %q, want %q", got, component)
	}

	// Test WithUserIP and GetUserIP
	userIP := "192.168.1.1"
	ctx = WithUserIP(ctx, userIP)
	if got := GetUserIP(ctx); got != userIP {
		t.Errorf("GetUserIP() = %q, want %q", got, userIP)
	}

	// Test with empty context (no values set)
	emptyCtx := context.Background()
	if got := GetRequestID(emptyCtx); got != "" {
		t.Errorf("GetRequestID(emptyCtx) = %q, want empty string", got)
	}
}

func TestGenerateRequestID(t *testing.T) {
	id1 := GenerateRequestID()
	id2 := GenerateRequestID()

	if id1 == "" {
		t.Error("GenerateRequestID() returned empty string")
	}

	if id1 == id2 {
		t.Error("GenerateRequestID() returned duplicate IDs")
	}

	// Should be 16 hex characters (8 bytes)
	if len(id1) != 16 {
		t.Errorf("GenerateRequestID() returned ID of length %d, want 16", len(id1))
	}
}

func TestRedactSensitive(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected map[string]string
	}{
		{
			name:     "password redacted",
			input:    map[string]string{"password": "secret123"},
			expected: map[string]string{"password": RedactedValue},
		},
		{
			name:     "token redacted",
			input:    map[string]string{"token": "abc123"},
			expected: map[string]string{"token": RedactedValue},
		},
		{
			name:     "api_key redacted",
			input:    map[string]string{"api_key": "key123"},
			expected: map[string]string{"api_key": RedactedValue},
		},
		{
			name:     "secret redacted",
			input:    map[string]string{"secret": "mysecret"},
			expected: map[string]string{"secret": RedactedValue},
		},
		{
			name:     "authorization redacted",
			input:    map[string]string{"authorization": "Bearer xyz"},
			expected: map[string]string{"authorization": RedactedValue},
		},
		{
			name:     "non-sensitive not redacted",
			input:    map[string]string{"username": "john", "email": "john@example.com"},
			expected: map[string]string{"username": "john", "email": "john@example.com"},
		},
		{
			name:     "mixed keys",
			input:    map[string]string{"username": "john", "password": "secret"},
			expected: map[string]string{"username": "john", "password": RedactedValue},
		},
		{
			name:     "case insensitive",
			input:    map[string]string{"PASSWORD": "secret", "Token": "abc"},
			expected: map[string]string{"PASSWORD": RedactedValue, "Token": RedactedValue},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RedactSensitive(tt.input)
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("RedactSensitive()[%q] = %q, want %q", k, result[k], v)
				}
			}
		})
	}
}

func TestIsSensitiveKey(t *testing.T) {
	sensitiveKeys := []string{
		"password",
		"PASSWORD",
		"user_password",
		"token",
		"access_token",
		"api_key",
		"apikey",
		"secret",
		"client_secret",
		"authorization",
		"auth_token",
	}

	for _, key := range sensitiveKeys {
		if !isSensitiveKey(key) {
			t.Errorf("isSensitiveKey(%q) = false, want true", key)
		}
	}

	nonSensitiveKeys := []string{
		"username",
		"email",
		"name",
		"id",
		"status",
		"message",
	}

	for _, key := range nonSensitiveKeys {
		if isSensitiveKey(key) {
			t.Errorf("isSensitiveKey(%q) = true, want false", key)
		}
	}
}

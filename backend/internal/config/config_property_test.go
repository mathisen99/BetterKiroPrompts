package config

import (
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/BurntSushi/toml"
)

// Property 1: Config structure completeness
// For any valid config.toml file, parsing it SHALL produce a Config struct
// with all required sections populated.
// **Validates: Requirements 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7**
func TestProperty1_ConfigStructureCompleteness(t *testing.T) {
	f := func(seed int64) bool {
		rng := rand.New(rand.NewSource(seed))
		cfg := generateValidConfig(rng)

		// Serialize to TOML
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.toml")
		file, err := os.Create(configPath)
		if err != nil {
			t.Logf("Failed to create temp file: %v", err)
			return false
		}
		encoder := toml.NewEncoder(file)
		if err := encoder.Encode(cfg); err != nil {
			_ = file.Close()
			t.Logf("Failed to encode config: %v", err)
			return false
		}
		_ = file.Close()

		// Parse back
		loaded, err := LoadFromPath(configPath)
		if err != nil {
			t.Logf("Failed to load config: %v", err)
			return false
		}

		// Verify all sections are populated (non-zero values)
		if loaded.Server.Port == 0 {
			t.Log("Server.Port is zero")
			return false
		}
		if loaded.OpenAI.Model == "" {
			t.Log("OpenAI.Model is empty")
			return false
		}
		if loaded.RateLimit.GenerationLimitPerHour == 0 {
			t.Log("RateLimit.GenerationLimitPerHour is zero")
			return false
		}
		if loaded.Logging.Level == "" {
			t.Log("Logging.Level is empty")
			return false
		}
		if loaded.Scanner.MaxRepoSizeMB == 0 {
			t.Log("Scanner.MaxRepoSizeMB is zero")
			return false
		}
		if loaded.Generation.MaxProjectIdeaLength == 0 {
			t.Log("Generation.MaxProjectIdeaLength is zero")
			return false
		}
		if loaded.Gallery.PageSize == 0 {
			t.Log("Gallery.PageSize is zero")
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 1 failed: %v", err)
	}
}

// Property 2: Default value fallback
// For any missing config.toml file or missing section, the Config_Loader
// SHALL return a Config struct with default values.
// **Validates: Requirements 1.8**
func TestProperty2_DefaultValueFallback(t *testing.T) {
	f := func(seed int64) bool {
		// Use a non-existent path
		tmpDir := t.TempDir()
		nonExistentPath := filepath.Join(tmpDir, "nonexistent.toml")

		loaded, err := LoadFromPath(nonExistentPath)
		if err != nil {
			t.Logf("Unexpected error loading non-existent config: %v", err)
			return false
		}

		defaults := DefaultConfig()

		// Verify defaults are applied
		if loaded.Server.Port != defaults.Server.Port {
			t.Logf("Server.Port mismatch: got %d, want %d", loaded.Server.Port, defaults.Server.Port)
			return false
		}
		if loaded.OpenAI.Model != defaults.OpenAI.Model {
			t.Logf("OpenAI.Model mismatch: got %s, want %s", loaded.OpenAI.Model, defaults.OpenAI.Model)
			return false
		}
		if loaded.RateLimit.GenerationLimitPerHour != defaults.RateLimit.GenerationLimitPerHour {
			t.Logf("RateLimit.GenerationLimitPerHour mismatch")
			return false
		}
		if loaded.Logging.Level != defaults.Logging.Level {
			t.Logf("Logging.Level mismatch")
			return false
		}
		if loaded.Scanner.MaxRepoSizeMB != defaults.Scanner.MaxRepoSizeMB {
			t.Logf("Scanner.MaxRepoSizeMB mismatch")
			return false
		}
		if loaded.Generation.MaxProjectIdeaLength != defaults.Generation.MaxProjectIdeaLength {
			t.Logf("Generation.MaxProjectIdeaLength mismatch")
			return false
		}
		if loaded.Gallery.PageSize != defaults.Gallery.PageSize {
			t.Logf("Gallery.PageSize mismatch")
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 2 failed: %v", err)
	}
}

// Property 3: Invalid config rejection
// For any config.toml with values outside valid ranges, the Config_Loader
// SHALL return an error describing the specific invalid field.
// **Validates: Requirements 1.9, 7.1, 7.2, 7.3**
func TestProperty3_InvalidConfigRejection(t *testing.T) {
	f := func(seed int64) bool {
		rng := rand.New(rand.NewSource(seed))
		cfg := generateInvalidConfig(rng)

		err := cfg.Validate()
		if err == nil {
			t.Log("Expected validation error but got none")
			return false
		}

		// Error message should contain useful information
		errStr := err.Error()
		if len(errStr) < 10 {
			t.Logf("Error message too short: %s", errStr)
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 3 failed: %v", err)
	}
}

// Property 4: Environment variable override
// For any configuration value that has both a config.toml value and an
// environment variable set, the environment variable value SHALL take precedence.
// **Validates: Requirements 2.2**
func TestProperty4_EnvironmentVariableOverride(t *testing.T) {
	f := func(seed int64) bool {
		rng := rand.New(rand.NewSource(seed))

		// Generate a valid config
		cfg := generateValidConfig(rng)

		// Write to temp file
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.toml")
		file, err := os.Create(configPath)
		if err != nil {
			t.Logf("Failed to create temp file: %v", err)
			return false
		}
		encoder := toml.NewEncoder(file)
		if err := encoder.Encode(cfg); err != nil {
			_ = file.Close()
			t.Logf("Failed to encode config: %v", err)
			return false
		}
		_ = file.Close()

		// Set environment variable to override
		envPort := 9999
		if err := os.Setenv("PORT", "9999"); err != nil {
			t.Logf("Failed to set env var: %v", err)
			return false
		}
		defer func() { _ = os.Unsetenv("PORT") }()

		// Load config
		loaded, err := LoadFromPath(configPath)
		if err != nil {
			t.Logf("Failed to load config: %v", err)
			return false
		}

		// Verify env var took precedence
		if loaded.Server.Port != envPort {
			t.Logf("Port not overridden: got %d, want %d", loaded.Server.Port, envPort)
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 4 failed: %v", err)
	}
}

// Property 6: Config round-trip
// For any valid Config struct, serializing to TOML and parsing back
// SHALL produce an equivalent Config struct.
// **Validates: Requirements 1.1-1.7**
func TestProperty6_ConfigRoundTrip(t *testing.T) {
	f := func(seed int64) bool {
		rng := rand.New(rand.NewSource(seed))
		original := generateValidConfig(rng)

		// Serialize to TOML
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.toml")
		file, err := os.Create(configPath)
		if err != nil {
			t.Logf("Failed to create temp file: %v", err)
			return false
		}
		encoder := toml.NewEncoder(file)
		if err := encoder.Encode(original); err != nil {
			_ = file.Close()
			t.Logf("Failed to encode config: %v", err)
			return false
		}
		_ = file.Close()

		// Parse back
		loaded, err := LoadFromPath(configPath)
		if err != nil {
			t.Logf("Failed to load config: %v", err)
			return false
		}

		// Compare key fields (deep equality)
		if !configsEqual(original, loaded) {
			t.Log("Configs not equal after round-trip")
			return false
		}

		return true
	}

	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 6 failed: %v", err)
	}
}

// Helper functions for generating test data

func generateValidConfig(rng *rand.Rand) *Config {
	reasoningEfforts := []string{"none", "low", "medium", "high", "xhigh"}
	verbosities := []string{"low", "medium", "high"}
	logLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	sortOptions := []string{"newest", "highest_rated", "most_viewed"}

	return &Config{
		Server: ServerConfig{
			Port:            1 + rng.Intn(65534),
			Host:            "0.0.0.0",
			ShutdownTimeout: Duration(time.Duration(1+rng.Intn(60)) * time.Second),
		},
		OpenAI: OpenAIConfig{
			Model:           "gpt-" + randomString(rng, 5),
			CodeReviewModel: "gpt-" + randomString(rng, 5),
			BaseURL:         "https://api.openai.com/v1",
			Timeout:         Duration(time.Duration(10+rng.Intn(300)) * time.Second),
			ReasoningEffort: reasoningEfforts[rng.Intn(len(reasoningEfforts))],
			Verbosity:       verbosities[rng.Intn(len(verbosities))],
		},
		RateLimit: RateLimitConfig{
			GenerationLimitPerHour: 1 + rng.Intn(100),
			RatingLimitPerHour:     1 + rng.Intn(100),
			ScanLimitPerHour:       1 + rng.Intn(100),
		},
		Logging: LoggingConfig{
			Level:       logLevels[rng.Intn(len(logLevels))],
			Directory:   "./logs",
			MaxSizeMB:   1 + rng.Intn(1000),
			MaxAgeDays:  1 + rng.Intn(365),
			EnableColor: rng.Intn(2) == 1,
		},
		Scanner: ScannerConfig{
			MaxRepoSizeMB:      1 + rng.Intn(1000),
			MaxReviewFiles:     1 + rng.Intn(100),
			ToolTimeoutSeconds: 10 + rng.Intn(600),
			RetentionDays:      1 + rng.Intn(365),
			CloneTimeout:       Duration(time.Duration(10+rng.Intn(600)) * time.Second),
		},
		Generation: GenerationConfig{
			MaxProjectIdeaLength: 100 + rng.Intn(10000),
			MaxAnswerLength:      100 + rng.Intn(10000),
			MinQuestions:         1 + rng.Intn(5),
			MaxQuestions:         6 + rng.Intn(15),
			MaxRetries:           rng.Intn(5),
		},
		Gallery: GalleryConfig{
			PageSize:    1 + rng.Intn(100),
			DefaultSort: sortOptions[rng.Intn(len(sortOptions))],
		},
	}
}

func generateInvalidConfig(rng *rand.Rand) *Config {
	cfg := generateValidConfig(rng)

	// Randomly invalidate one field
	invalidationType := rng.Intn(10)
	switch invalidationType {
	case 0:
		cfg.Server.Port = -1 // Invalid port
	case 1:
		cfg.Server.Port = 70000 // Port too high
	case 2:
		cfg.OpenAI.Model = "" // Empty model
	case 3:
		cfg.OpenAI.ReasoningEffort = "invalid" // Invalid enum
	case 4:
		cfg.OpenAI.Verbosity = "invalid" // Invalid enum
	case 5:
		cfg.Logging.Level = "INVALID" // Invalid log level
	case 6:
		cfg.RateLimit.GenerationLimitPerHour = 0 // Zero rate limit
	case 7:
		cfg.Scanner.ToolTimeoutSeconds = 5 // Too low
	case 8:
		cfg.Generation.MaxQuestions = 0 // Less than min
	case 9:
		cfg.Gallery.DefaultSort = "invalid" // Invalid sort
	}

	return cfg
}

func randomString(rng *rand.Rand, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

func configsEqual(a, b *Config) bool {
	// Compare using reflect.DeepEqual
	return reflect.DeepEqual(a.Server, b.Server) &&
		reflect.DeepEqual(a.OpenAI, b.OpenAI) &&
		reflect.DeepEqual(a.RateLimit, b.RateLimit) &&
		reflect.DeepEqual(a.Logging, b.Logging) &&
		reflect.DeepEqual(a.Scanner, b.Scanner) &&
		reflect.DeepEqual(a.Generation, b.Generation) &&
		reflect.DeepEqual(a.Gallery, b.Gallery)
}

package config

import (
	"testing"
)

func TestExampleConfigLoads(t *testing.T) {
	cfg, err := LoadFromPath("../../../config.example.toml")
	if err != nil {
		t.Fatalf("Failed to load config.example.toml: %v", err)
	}

	// Verify default values are set correctly
	if cfg.Server.Port != 8090 {
		t.Errorf("Expected server.port = 8090, got %d", cfg.Server.Port)
	}
	if cfg.OpenAI.Model != "gpt-5.2" {
		t.Errorf("Expected openai.model = gpt-5.2, got %s", cfg.OpenAI.Model)
	}
	if cfg.Logging.Level != "INFO" {
		t.Errorf("Expected logging.level = INFO, got %s", cfg.Logging.Level)
	}
	if cfg.RateLimit.GenerationLimitPerHour != 10 {
		t.Errorf("Expected rate_limit.generation_limit_per_hour = 10, got %d", cfg.RateLimit.GenerationLimitPerHour)
	}
	if cfg.Scanner.MaxRepoSizeMB != 500 {
		t.Errorf("Expected scanner.max_repo_size_mb = 500, got %d", cfg.Scanner.MaxRepoSizeMB)
	}
	if cfg.Generation.MaxProjectIdeaLength != 2000 {
		t.Errorf("Expected generation.max_project_idea_length = 2000, got %d", cfg.Generation.MaxProjectIdeaLength)
	}
	if cfg.Gallery.PageSize != 20 {
		t.Errorf("Expected gallery.page_size = 20, got %d", cfg.Gallery.PageSize)
	}
}

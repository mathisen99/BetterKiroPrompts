package generator

import (
	"encoding/json"
	"testing"
)

func TestGenerateHooks_LightPreset(t *testing.T) {
	config := HooksConfig{
		Preset:    "light",
		TechStack: HooksTechStack{HasGo: true, HasTypeScript: true},
	}

	files, err := GenerateHooks(config)
	if err != nil {
		t.Fatalf("GenerateHooks failed: %v", err)
	}

	// Light: only formatters (2 files)
	if len(files) != 2 {
		t.Errorf("Expected 2 files for light preset, got %d", len(files))
	}
}

func TestGenerateHooks_BasicPreset(t *testing.T) {
	config := HooksConfig{
		Preset:    "basic",
		TechStack: HooksTechStack{HasGo: true, HasTypeScript: true},
	}

	files, err := GenerateHooks(config)
	if err != nil {
		t.Fatalf("GenerateHooks failed: %v", err)
	}

	// Basic: formatters + linters + tests (5 files)
	if len(files) < 4 {
		t.Errorf("Expected at least 4 files for basic preset, got %d", len(files))
	}
}

func TestGenerateHooks_DefaultPreset(t *testing.T) {
	config := HooksConfig{
		Preset:    "default",
		TechStack: HooksTechStack{HasGo: true, HasTypeScript: true},
	}

	files, err := GenerateHooks(config)
	if err != nil {
		t.Fatalf("GenerateHooks failed: %v", err)
	}

	// Default: basic + secret scan + prompt guard
	if len(files) < 6 {
		t.Errorf("Expected at least 6 files for default preset, got %d", len(files))
	}
}

func TestGenerateHooks_StrictPreset(t *testing.T) {
	config := HooksConfig{
		Preset:    "strict",
		TechStack: HooksTechStack{HasGo: true, HasTypeScript: true},
	}

	files, err := GenerateHooks(config)
	if err != nil {
		t.Fatalf("GenerateHooks failed: %v", err)
	}

	// Strict: default + static analysis + vuln scan
	if len(files) < 8 {
		t.Errorf("Expected at least 8 files for strict preset, got %d", len(files))
	}
}

func TestGenerateHooks_ValidJSON(t *testing.T) {
	config := HooksConfig{
		Preset:    "light",
		TechStack: HooksTechStack{HasGo: true},
	}

	files, err := GenerateHooks(config)
	if err != nil {
		t.Fatalf("GenerateHooks failed: %v", err)
	}

	for _, f := range files {
		var hook map[string]interface{}
		if err := json.Unmarshal([]byte(f.Content), &hook); err != nil {
			t.Errorf("Invalid JSON in %s: %v", f.Path, err)
		}
		if hook["name"] == nil {
			t.Errorf("Hook %s missing 'name' field", f.Path)
		}
	}
}

package prompts

import (
	"testing"
)

// TestExperienceLevelValidation tests that experience level validation works correctly.
func TestExperienceLevelValidation(t *testing.T) {
	tests := []struct {
		level string
		valid bool
	}{
		{"beginner", true},
		{"novice", true},
		{"expert", true},
		{"invalid", false},
		{"", false},
		{"BEGINNER", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			got := IsValidExperienceLevel(tt.level)
			if got != tt.valid {
				t.Errorf("IsValidExperienceLevel(%q) = %v, want %v", tt.level, got, tt.valid)
			}
		})
	}
}

// TestHookPresetValidation tests that hook preset validation works correctly.
func TestHookPresetValidation(t *testing.T) {
	tests := []struct {
		preset string
		valid  bool
	}{
		{"light", true},
		{"basic", true},
		{"default", true},
		{"strict", true},
		{"invalid", false},
		{"", false},
		{"LIGHT", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.preset, func(t *testing.T) {
			got := IsValidHookPreset(tt.preset)
			if got != tt.valid {
				t.Errorf("IsValidHookPreset(%q) = %v, want %v", tt.preset, got, tt.valid)
			}
		})
	}
}

// TestQuestionsSystemPromptGeneration tests that question prompts are generated for each level.
func TestQuestionsSystemPromptGeneration(t *testing.T) {
	levels := []string{ExperienceBeginner, ExperienceNovice, ExperienceExpert}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			prompt := GetQuestionsSystemPrompt(level)
			if prompt == "" {
				t.Errorf("GetQuestionsSystemPrompt(%q) returned empty string", level)
			}
			if len(prompt) < 100 {
				t.Errorf("GetQuestionsSystemPrompt(%q) returned suspiciously short prompt: %d chars", level, len(prompt))
			}
		})
	}
}

// TestOutputsSystemPromptGeneration tests that output prompts are generated for each combination.
func TestOutputsSystemPromptGeneration(t *testing.T) {
	levels := []string{ExperienceBeginner, ExperienceNovice, ExperienceExpert}
	presets := []string{HookPresetLight, HookPresetBasic, HookPresetDefault, HookPresetStrict}

	for _, level := range levels {
		for _, preset := range presets {
			t.Run(level+"_"+preset, func(t *testing.T) {
				prompt := GetOutputsSystemPrompt(level, preset)
				if prompt == "" {
					t.Errorf("GetOutputsSystemPrompt(%q, %q) returned empty string", level, preset)
				}
				if len(prompt) < 500 {
					t.Errorf("GetOutputsSystemPrompt(%q, %q) returned suspiciously short prompt: %d chars", level, preset, len(prompt))
				}
			})
		}
	}
}

// TestHookPresetDescriptions tests that all presets have descriptions.
func TestHookPresetDescriptions(t *testing.T) {
	presets := []string{HookPresetLight, HookPresetBasic, HookPresetDefault, HookPresetStrict}

	for _, preset := range presets {
		t.Run(preset, func(t *testing.T) {
			info, ok := HookPresetDescriptions[preset]
			if !ok {
				t.Errorf("HookPresetDescriptions missing entry for %q", preset)
				return
			}
			if info.Title == "" {
				t.Errorf("HookPresetDescriptions[%q].Title is empty", preset)
			}
			if info.Description == "" {
				t.Errorf("HookPresetDescriptions[%q].Description is empty", preset)
			}
			if len(info.Hooks) == 0 {
				t.Errorf("HookPresetDescriptions[%q].Hooks is empty", preset)
			}
		})
	}
}

// TestValidExperienceLevels tests that ValidExperienceLevels returns all levels.
func TestValidExperienceLevels(t *testing.T) {
	levels := ValidExperienceLevels()
	if len(levels) != 3 {
		t.Errorf("ValidExperienceLevels() returned %d levels, want 3", len(levels))
	}

	expected := map[string]bool{
		ExperienceBeginner: true,
		ExperienceNovice:   true,
		ExperienceExpert:   true,
	}

	for _, level := range levels {
		if !expected[level] {
			t.Errorf("ValidExperienceLevels() contains unexpected level %q", level)
		}
	}
}

// TestValidHookPresets tests that ValidHookPresets returns all presets.
func TestValidHookPresets(t *testing.T) {
	presets := ValidHookPresets()
	if len(presets) != 4 {
		t.Errorf("ValidHookPresets() returned %d presets, want 4", len(presets))
	}

	expected := map[string]bool{
		HookPresetLight:   true,
		HookPresetBasic:   true,
		HookPresetDefault: true,
		HookPresetStrict:  true,
	}

	for _, preset := range presets {
		if !expected[preset] {
			t.Errorf("ValidHookPresets() contains unexpected preset %q", preset)
		}
	}
}

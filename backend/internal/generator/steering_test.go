package generator

import (
	"strings"
	"testing"
)

func TestGenerateSteering_FoundationFiles(t *testing.T) {
	config := SteeringConfig{
		ProjectName:        "Test Project",
		ProjectDescription: "A test project",
	}

	files, err := GenerateSteering(config)
	if err != nil {
		t.Fatalf("GenerateSteering failed: %v", err)
	}

	// Should have 4 foundation files: product, tech, structure, AGENTS.md
	if len(files) != 4 {
		t.Errorf("Expected 4 files, got %d", len(files))
	}

	// Verify expected paths
	paths := make(map[string]bool)
	for _, f := range files {
		paths[f.Path] = true
	}

	expected := []string{
		".kiro/steering/product.md",
		".kiro/steering/tech.md",
		".kiro/steering/structure.md",
		"AGENTS.md",
	}
	for _, p := range expected {
		if !paths[p] {
			t.Errorf("Expected file %s not found", p)
		}
	}
}

func TestGenerateSteering_ConditionalFiles(t *testing.T) {
	config := SteeringConfig{
		ProjectName:        "Test Project",
		IncludeConditional: true,
	}

	files, err := GenerateSteering(config)
	if err != nil {
		t.Fatalf("GenerateSteering failed: %v", err)
	}

	// 4 foundation + 4 conditional = 8
	if len(files) != 8 {
		t.Errorf("Expected 8 files with conditional, got %d", len(files))
	}
}

func TestGenerateSteering_ManualFiles(t *testing.T) {
	config := SteeringConfig{
		ProjectName:   "Test Project",
		IncludeManual: true,
	}

	files, err := GenerateSteering(config)
	if err != nil {
		t.Fatalf("GenerateSteering failed: %v", err)
	}

	// 4 foundation + 1 manual = 5
	if len(files) != 5 {
		t.Errorf("Expected 5 files with manual, got %d", len(files))
	}

	// Verify manual file has correct frontmatter
	for _, f := range files {
		if f.Path == ".kiro/steering/manual-example.md" {
			if !strings.Contains(f.Content, "inclusion: manual") {
				t.Error("Manual file should contain 'inclusion: manual' frontmatter")
			}
			return
		}
	}
	t.Error("Manual file not found")
}

func TestGenerateSteering_Frontmatter(t *testing.T) {
	config := SteeringConfig{
		ProjectName: "Test Project",
	}

	files, err := GenerateSteering(config)
	if err != nil {
		t.Fatalf("GenerateSteering failed: %v", err)
	}

	for _, f := range files {
		if f.Path == ".kiro/steering/product.md" {
			if !strings.Contains(f.Content, "inclusion: always") {
				t.Error("Product file should contain 'inclusion: always' frontmatter")
			}
			return
		}
	}
}

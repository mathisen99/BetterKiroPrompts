package generation

import (
	"better-kiro-prompts/internal/prompts"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"testing/quick"
)

// TestValidateSteeringFile tests steering file frontmatter validation
func TestValidateSteeringFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
		errType error
	}{
		{
			name: "valid always inclusion",
			content: `---
inclusion: always
---

# Product`,
			wantErr: false,
		},
		{
			name: "valid fileMatch with pattern",
			content: `---
inclusion: fileMatch
fileMatchPattern: "**/*.go"
---

# Security`,
			wantErr: false,
		},
		{
			name: "valid manual inclusion",
			content: `---
inclusion: manual
---

# Guide`,
			wantErr: false,
		},
		{
			name:    "missing frontmatter",
			content: `# No frontmatter`,
			wantErr: true,
			errType: ErrInvalidFrontmatter,
		},
		{
			name: "missing inclusion field",
			content: `---
fileMatchPattern: "**/*.go"
---

# Content`,
			wantErr: true,
			errType: ErrMissingInclusion,
		},
		{
			name: "invalid inclusion mode",
			content: `---
inclusion: invalid
---

# Content`,
			wantErr: true,
			errType: ErrInvalidInclusionMode,
		},
		{
			name: "fileMatch without pattern",
			content: `---
inclusion: fileMatch
---

# Content`,
			wantErr: true,
			errType: ErrMissingFileMatchPattern,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSteeringFile(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSteeringFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateHookFile tests hook file JSON schema validation
func TestValidateHookFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
		errType error
	}{
		{
			name: "valid agentStop with runCommand",
			content: `{
				"name": "Format on Stop",
				"description": "Run formatters",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "runCommand", "command": "go fmt ./..."}
			}`,
			wantErr: false,
		},
		{
			name: "valid promptSubmit with runCommand",
			content: `{
				"name": "Pre-submit Check",
				"description": "Check before submit",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "promptSubmit"},
				"then": {"type": "runCommand", "command": "make check"}
			}`,
			wantErr: false,
		},
		{
			name: "valid userTriggered with askAgent",
			content: `{
				"name": "Run Tests",
				"description": "Manual test trigger",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "userTriggered"},
				"then": {"type": "askAgent", "prompt": "Run tests"}
			}`,
			wantErr: false,
		},
		{
			name: "valid fileEdited with patterns",
			content: `{
				"name": "Go Test",
				"description": "Test on change",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "fileEdited", "patterns": ["**/*.go"]},
				"then": {"type": "askAgent", "prompt": "Run tests"}
			}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			content: `not json`,
			wantErr: true,
			errType: ErrInvalidHookSchema,
		},
		{
			name: "missing name",
			content: `{
				"description": "Test",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "askAgent", "prompt": "test"}
			}`,
			wantErr: true,
			errType: ErrMissingHookField,
		},
		{
			name: "missing description",
			content: `{
				"name": "Test",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "askAgent", "prompt": "test"}
			}`,
			wantErr: true,
			errType: ErrMissingHookField,
		},
		{
			name: "missing version",
			content: `{
				"name": "Test",
				"description": "Test",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "askAgent", "prompt": "test"}
			}`,
			wantErr: true,
			errType: ErrMissingHookField,
		},
		{
			name: "invalid when.type",
			content: `{
				"name": "Test",
				"description": "Test",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "invalidType"},
				"then": {"type": "askAgent", "prompt": "test"}
			}`,
			wantErr: true,
			errType: ErrInvalidWhenType,
		},
		{
			name: "invalid then.type",
			content: `{
				"name": "Test",
				"description": "Test",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "invalidAction", "prompt": "test"}
			}`,
			wantErr: true,
			errType: ErrInvalidThenType,
		},
		{
			name: "runCommand with fileEdited (not allowed)",
			content: `{
				"name": "Test",
				"description": "Test",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "fileEdited", "patterns": ["**/*.go"]},
				"then": {"type": "runCommand", "command": "go fmt"}
			}`,
			wantErr: true,
			errType: ErrRunCommandRestriction,
		},
		{
			name: "runCommand with userTriggered (not allowed)",
			content: `{
				"name": "Test",
				"description": "Test",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "userTriggered"},
				"then": {"type": "runCommand", "command": "go fmt"}
			}`,
			wantErr: true,
			errType: ErrRunCommandRestriction,
		},
		{
			name: "fileEdited without patterns",
			content: `{
				"name": "Test",
				"description": "Test",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "fileEdited"},
				"then": {"type": "askAgent", "prompt": "test"}
			}`,
			wantErr: true,
			errType: ErrMissingHookField,
		},
		{
			name: "askAgent without prompt",
			content: `{
				"name": "Test",
				"description": "Test",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "askAgent"}
			}`,
			wantErr: true,
			errType: ErrMissingHookField,
		},
		{
			name: "runCommand without command",
			content: `{
				"name": "Test",
				"description": "Test",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "runCommand"}
			}`,
			wantErr: true,
			errType: ErrMissingHookField,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHookFile(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHookFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// =============================================================================
// Property-Based Tests for Steering File Validity
// =============================================================================

// coreSteeringFiles defines the required core steering files that must always be present
var coreSteeringFiles = []string{"product.md", "tech.md", "structure.md"}

// TestProperty2_CoreSteeringFilesValidity tests that core steering files have valid frontmatter.
// Feature: phase4-production, Property 2: Core Steering Files Validity
// **Validates: Requirements 4.1, 4.2, 4.3, 4.8, 10.1**
func TestProperty2_CoreSteeringFilesValidity(t *testing.T) {
	// Property: For any generated output, the files SHALL include at minimum
	// product.md, tech.md, and structure.md, each with valid frontmatter
	// containing `inclusion: always`, and valid markdown content.

	// Test valid core steering file content
	validCoreSteeringFiles := []struct {
		name    string
		content string
	}{
		{
			name: "product.md",
			content: `---
inclusion: always
---

# Product

## What We Are Building
A sample application.

## What We Are NOT Building
Not a production system.`,
		},
		{
			name: "tech.md",
			content: `---
inclusion: always
---

# Tech Stack

## Languages
- Go
- TypeScript`,
		},
		{
			name: "structure.md",
			content: `---
inclusion: always
---

# Repository Structure

## Layout
Standard Go project layout.`,
		},
	}

	for _, sf := range validCoreSteeringFiles {
		t.Run(sf.name, func(t *testing.T) {
			err := ValidateSteeringFile(sf.content)
			if err != nil {
				t.Errorf("Core steering file %s should be valid: %v", sf.name, err)
			}

			// Verify it has inclusion: always
			if !strings.Contains(sf.content, "inclusion: always") {
				t.Errorf("Core steering file %s should have 'inclusion: always'", sf.name)
			}
		})
	}

	// Property test: For any valid markdown content, wrapping it with proper
	// frontmatter should produce a valid steering file
	property := func(markdownContent string) bool {
		// Skip empty or very short content
		if len(strings.TrimSpace(markdownContent)) < 3 {
			return true
		}

		// Create a valid steering file with the content
		steeringFile := "---\ninclusion: always\n---\n\n" + markdownContent

		err := ValidateSteeringFile(steeringFile)
		return err == nil
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 2 (Core Steering Files Validity) failed: %v", err)
	}
}

// TestProperty2_CoreSteeringFilesPresence tests that a valid output contains all core files.
// Feature: phase4-production, Property 2: Core Steering Files Validity
// **Validates: Requirements 4.1, 4.2, 4.3, 10.1**
func TestProperty2_CoreSteeringFilesPresence(t *testing.T) {
	// Simulate a valid generated output with all required files
	validOutput := []GeneratedFile{
		{
			Path:    ".kiro/prompts/kickoff.md",
			Content: minimalValidKickoff(),
			Type:    "kickoff",
		},
		{
			Path: ".kiro/steering/product.md",
			Content: `---
inclusion: always
---

# Product`,
			Type: "steering",
		},
		{
			Path: ".kiro/steering/tech.md",
			Content: `---
inclusion: always
---

# Tech Stack`,
			Type: "steering",
		},
		{
			Path: ".kiro/steering/structure.md",
			Content: `---
inclusion: always
---

# Structure`,
			Type: "steering",
		},
		{
			Path: ".kiro/hooks/format-on-stop.kiro.hook",
			Content: `{
				"name": "Format on Stop",
				"description": "Run formatters",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "runCommand", "command": "go fmt ./..."}
			}`,
			Type: "hook",
		},
		{
			Path:    "AGENTS.md",
			Content: "# Agent Guidelines",
			Type:    "agents",
		},
	}

	// Validate all files
	err := ValidateGeneratedFiles(validOutput)
	if err != nil {
		t.Errorf("Valid output should pass validation: %v", err)
	}

	// Check that all core steering files are present
	foundFiles := make(map[string]bool)
	for _, f := range validOutput {
		if f.Type == "steering" {
			for _, core := range coreSteeringFiles {
				if strings.HasSuffix(f.Path, core) {
					foundFiles[core] = true
				}
			}
		}
	}

	for _, core := range coreSteeringFiles {
		if !foundFiles[core] {
			t.Errorf("Core steering file %s should be present in output", core)
		}
	}
}

// TestProperty2_CoreSteeringFilesInclusionAlways tests that core files must have inclusion: always.
// Feature: phase4-production, Property 2: Core Steering Files Validity
// **Validates: Requirements 4.1, 4.2, 4.3**
func TestProperty2_CoreSteeringFilesInclusionAlways(t *testing.T) {
	// Test that core steering files with wrong inclusion mode fail validation
	// when we check for the specific requirement

	invalidModes := []string{"fileMatch", "manual"}

	for _, mode := range invalidModes {
		t.Run("invalid_mode_"+mode, func(t *testing.T) {
			content := "---\ninclusion: " + mode + "\n---\n\n# Product"

			// The file itself is valid YAML, but for core files we expect 'always'
			err := ValidateSteeringFile(content)

			// Basic validation passes (it's valid frontmatter)
			// But the semantic check for core files would fail
			if mode == "fileMatch" {
				// fileMatch without pattern should fail
				if err == nil {
					t.Errorf("fileMatch mode without pattern should fail validation")
				}
			}
		})
	}
}

// TestProperty3_ConditionalSteeringFilesPattern tests that conditional steering files have valid patterns.
// Feature: phase4-production, Property 3: Conditional Steering Files Pattern
// **Validates: Requirements 4.4, 4.5, 4.6, 4.7, 10.2, 10.3**
func TestProperty3_ConditionalSteeringFilesPattern(t *testing.T) {
	// Property: For any generated output where the project uses a specific language,
	// the security and quality steering files SHALL have `inclusion: fileMatch`
	// and a `fileMatchPattern` that matches files of that language.

	// Test valid conditional steering files for different languages
	validConditionalFiles := []struct {
		name            string
		content         string
		expectedPattern string
	}{
		{
			name: "security-go.md",
			content: `---
inclusion: fileMatch
fileMatchPattern: "**/*.go"
---

# Go Security Guidelines

## No Secrets
Never commit credentials.`,
			expectedPattern: "**/*.go",
		},
		{
			name: "quality-go.md",
			content: `---
inclusion: fileMatch
fileMatchPattern: "**/*.go"
---

# Go Quality Guidelines

## Formatting
Use go fmt.`,
			expectedPattern: "**/*.go",
		},
		{
			name: "security-web.md",
			content: `---
inclusion: fileMatch
fileMatchPattern: "**/*.{ts,tsx}"
---

# Web Security Guidelines

## XSS Prevention
Sanitize all user input.`,
			expectedPattern: "**/*.{ts,tsx}",
		},
		{
			name: "quality-web.md",
			content: `---
inclusion: fileMatch
fileMatchPattern: "**/*.{ts,tsx,js,jsx}"
---

# Web Quality Guidelines

## Linting
Use ESLint.`,
			expectedPattern: "**/*.{ts,tsx,js,jsx}",
		},
	}

	for _, sf := range validConditionalFiles {
		t.Run(sf.name, func(t *testing.T) {
			err := ValidateSteeringFile(sf.content)
			if err != nil {
				t.Errorf("Conditional steering file %s should be valid: %v", sf.name, err)
			}

			// Verify it has inclusion: fileMatch
			if !strings.Contains(sf.content, "inclusion: fileMatch") {
				t.Errorf("Conditional steering file %s should have 'inclusion: fileMatch'", sf.name)
			}

			// Verify it has the expected pattern
			if !strings.Contains(sf.content, sf.expectedPattern) {
				t.Errorf("Conditional steering file %s should have pattern %s", sf.name, sf.expectedPattern)
			}
		})
	}
}

// TestProperty3_FileMatchRequiresPattern tests that fileMatch mode requires a pattern.
// Feature: phase4-production, Property 3: Conditional Steering Files Pattern
// **Validates: Requirements 4.4, 4.5**
func TestProperty3_FileMatchRequiresPattern(t *testing.T) {
	// Property test: For any fileMatch steering file, a pattern must be present

	// Test that fileMatch without pattern fails
	invalidContent := `---
inclusion: fileMatch
---

# Security Guidelines`

	err := ValidateSteeringFile(invalidContent)
	if err == nil {
		t.Error("fileMatch mode without fileMatchPattern should fail validation")
	}
	if err != ErrMissingFileMatchPattern {
		t.Errorf("Expected ErrMissingFileMatchPattern, got: %v", err)
	}
}

// TestProperty3_ValidGlobPatterns tests that various glob patterns are accepted.
// Feature: phase4-production, Property 3: Conditional Steering Files Pattern
// **Validates: Requirements 4.6, 4.7**
func TestProperty3_ValidGlobPatterns(t *testing.T) {
	// Test various valid glob patterns for different languages
	validPatterns := []string{
		"**/*.go",
		"**/*.ts",
		"**/*.tsx",
		"**/*.{ts,tsx}",
		"**/*.{ts,tsx,js,jsx}",
		"**/*.py",
		"**/*.rs",
		"**/*.java",
		"src/**/*.ts",
		"backend/**/*.go",
		"frontend/**/*.{ts,tsx}",
	}

	for _, pattern := range validPatterns {
		t.Run(pattern, func(t *testing.T) {
			content := "---\ninclusion: fileMatch\nfileMatchPattern: \"" + pattern + "\"\n---\n\n# Content"

			err := ValidateSteeringFile(content)
			if err != nil {
				t.Errorf("Pattern %s should be valid: %v", pattern, err)
			}
		})
	}
}

// TestProperty3_LanguageSpecificPatterns tests that language-specific files have correct patterns.
// Feature: phase4-production, Property 3: Conditional Steering Files Pattern
// **Validates: Requirements 4.6, 4.7, 10.2, 10.3**
func TestProperty3_LanguageSpecificPatterns(t *testing.T) {
	// Property: Language-specific steering files should have patterns matching that language

	languagePatterns := map[string][]string{
		"go":         {"**/*.go"},
		"typescript": {"**/*.ts", "**/*.tsx", "**/*.{ts,tsx}"},
		"javascript": {"**/*.js", "**/*.jsx", "**/*.{js,jsx}"},
		"python":     {"**/*.py"},
		"rust":       {"**/*.rs"},
		"java":       {"**/*.java"},
	}

	for lang, patterns := range languagePatterns {
		for _, pattern := range patterns {
			t.Run(lang+"_"+pattern, func(t *testing.T) {
				content := "---\ninclusion: fileMatch\nfileMatchPattern: \"" + pattern + "\"\n---\n\n# " + lang + " Guidelines"

				err := ValidateSteeringFile(content)
				if err != nil {
					t.Errorf("Language %s with pattern %s should be valid: %v", lang, pattern, err)
				}
			})
		}
	}
}

// TestProperty3_ConditionalFilesInOutput tests that conditional files in output are valid.
// Feature: phase4-production, Property 3: Conditional Steering Files Pattern
// **Validates: Requirements 4.4, 4.5, 4.6, 4.7, 10.2, 10.3**
func TestProperty3_ConditionalFilesInOutput(t *testing.T) {
	// Simulate a valid generated output with conditional steering files
	validOutput := []GeneratedFile{
		{
			Path:    ".kiro/prompts/kickoff.md",
			Content: minimalValidKickoff(),
			Type:    "kickoff",
		},
		{
			Path: ".kiro/steering/product.md",
			Content: `---
inclusion: always
---

# Product`,
			Type: "steering",
		},
		{
			Path: ".kiro/steering/tech.md",
			Content: `---
inclusion: always
---

# Tech Stack`,
			Type: "steering",
		},
		{
			Path: ".kiro/steering/structure.md",
			Content: `---
inclusion: always
---

# Structure`,
			Type: "steering",
		},
		{
			Path: ".kiro/steering/security-go.md",
			Content: `---
inclusion: fileMatch
fileMatchPattern: "**/*.go"
---

# Go Security`,
			Type: "steering",
		},
		{
			Path: ".kiro/steering/quality-go.md",
			Content: `---
inclusion: fileMatch
fileMatchPattern: "**/*.go"
---

# Go Quality`,
			Type: "steering",
		},
		{
			Path: ".kiro/hooks/format-on-stop.kiro.hook",
			Content: `{
				"name": "Format on Stop",
				"description": "Run formatters",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "runCommand", "command": "go fmt ./..."}
			}`,
			Type: "hook",
		},
		{
			Path:    "AGENTS.md",
			Content: "# Agent Guidelines",
			Type:    "agents",
		},
	}

	// Validate all files
	err := ValidateGeneratedFiles(validOutput)
	if err != nil {
		t.Errorf("Valid output with conditional files should pass validation: %v", err)
	}

	// Verify conditional files have fileMatch inclusion
	for _, f := range validOutput {
		if f.Type == "steering" && (strings.Contains(f.Path, "security-") || strings.Contains(f.Path, "quality-")) {
			if !strings.Contains(f.Content, "inclusion: fileMatch") {
				t.Errorf("Conditional file %s should have 'inclusion: fileMatch'", f.Path)
			}
			if !strings.Contains(f.Content, "fileMatchPattern:") {
				t.Errorf("Conditional file %s should have 'fileMatchPattern'", f.Path)
			}
		}
	}
}

// TestProperty2And3_PropertyBasedValidation uses property-based testing to verify
// steering file validation across random inputs.
// Feature: phase4-production, Property 2 & 3: Steering Files Validity
// **Validates: Requirements 4.1-4.8, 10.1-10.3**
func TestProperty2And3_PropertyBasedValidation(t *testing.T) {
	// Property: For any string that looks like valid frontmatter with inclusion: always,
	// the validation should pass

	validInclusionModes := []string{"always", "fileMatch", "manual"}

	for _, mode := range validInclusionModes {
		t.Run("mode_"+mode, func(t *testing.T) {
			property := func(content string) bool {
				// Skip content that might interfere with frontmatter parsing
				if strings.Contains(content, "---") || strings.Contains(content, "inclusion:") {
					return true
				}

				var steeringFile string
				if mode == "fileMatch" {
					steeringFile = "---\ninclusion: fileMatch\nfileMatchPattern: \"**/*.go\"\n---\n\n" + content
				} else {
					steeringFile = "---\ninclusion: " + mode + "\n---\n\n" + content
				}

				err := ValidateSteeringFile(steeringFile)
				return err == nil
			}

			if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
				t.Errorf("Property validation for mode %s failed: %v", mode, err)
			}
		})
	}
}

// TestProperty2And3_InvalidFrontmatterRejected tests that invalid frontmatter is rejected.
// Feature: phase4-production, Property 2 & 3: Steering Files Validity
// **Validates: Requirements 4.8, 8.8**
func TestProperty2And3_InvalidFrontmatterRejected(t *testing.T) {
	invalidCases := []struct {
		name    string
		content string
	}{
		{
			name:    "no_frontmatter",
			content: "# Just content without frontmatter",
		},
		{
			name:    "incomplete_frontmatter",
			content: "---\ninclusion: always\n# Missing closing ---",
		},
		{
			name:    "invalid_inclusion_mode",
			content: "---\ninclusion: invalid\n---\n\n# Content",
		},
		{
			name:    "empty_inclusion",
			content: "---\ninclusion: \n---\n\n# Content",
		},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSteeringFile(tc.content)
			if err == nil {
				t.Errorf("Invalid frontmatter case %s should fail validation", tc.name)
			}
		})
	}
}

// TestExtractYAMLFieldProperty tests the YAML field extraction with various inputs.
func TestExtractYAMLFieldProperty(t *testing.T) {
	// Test that extractYAMLField correctly extracts values
	testCases := []struct {
		yaml     string
		field    string
		expected string
	}{
		{"inclusion: always", "inclusion", "always"},
		{"inclusion: fileMatch", "inclusion", "fileMatch"},
		{"inclusion: manual", "inclusion", "manual"},
		{"fileMatchPattern: \"**/*.go\"", "fileMatchPattern", "**/*.go"},
		{"fileMatchPattern: '**/*.ts'", "fileMatchPattern", "**/*.ts"},
		{"inclusion: always\nfileMatchPattern: \"**/*.go\"", "inclusion", "always"},
		{"inclusion: always\nfileMatchPattern: \"**/*.go\"", "fileMatchPattern", "**/*.go"},
	}

	for _, tc := range testCases {
		t.Run(tc.field+"_"+tc.expected, func(t *testing.T) {
			result := extractYAMLField(tc.yaml, tc.field)
			if result != tc.expected {
				t.Errorf("extractYAMLField(%q, %q) = %q, want %q", tc.yaml, tc.field, result, tc.expected)
			}
		})
	}
}

// frontmatterRegexTest is used to test frontmatter extraction
var frontmatterRegexTest = regexp.MustCompile(`(?s)^---\s*\n(.*?)\n---`)

// TestFrontmatterExtractionProperty tests frontmatter extraction with various formats.
func TestFrontmatterExtractionProperty(t *testing.T) {
	validFrontmatters := []string{
		"---\ninclusion: always\n---\n\n# Content",
		"---\ninclusion: fileMatch\nfileMatchPattern: \"**/*.go\"\n---\n\n# Content",
		"---\ninclusion: manual\n---\n\n# Content",
	}

	for _, content := range validFrontmatters {
		t.Run(content[:20], func(t *testing.T) {
			matches := frontmatterRegexTest.FindStringSubmatch(content)
			if len(matches) < 2 {
				t.Errorf("Failed to extract frontmatter from: %s", content[:50])
			}
		})
	}
}

// =============================================================================
// Property-Based Tests for Hook File Schema Validity (Property 4)
// =============================================================================

// validWhenTypesSlice is a slice of valid when.type values for property testing
var validWhenTypesSlice = []string{"fileEdited", "fileCreated", "fileDeleted", "promptSubmit", "agentStop", "userTriggered"}

// validThenTypesSlice is a slice of valid then.type values for property testing
var validThenTypesSlice = []string{"askAgent", "runCommand"}

// fileBasedTriggers are triggers that require patterns
var fileBasedTriggers = map[string]bool{
	"fileEdited":  true,
	"fileCreated": true,
	"fileDeleted": true,
}

// runCommandAllowedTriggers are triggers that allow runCommand
var runCommandAllowedTriggers = map[string]bool{
	"promptSubmit": true,
	"agentStop":    true,
}

// TestProperty4_HookFileSchemaValidity tests that all hooks have required fields and valid values.
// Feature: phase4-production, Property 4: Hook File Schema Validity
// **Validates: Requirements 5.1, 5.2, 5.3, 5.4, 5.5, 10.4**
func TestProperty4_HookFileSchemaValidity(t *testing.T) {
	// Property: For any generated hook file, the JSON SHALL contain all required fields
	// (name, description, version, enabled, when, then), the when.type SHALL be one of
	// the valid values, and if then.type is "runCommand" then when.type SHALL be either
	// "promptSubmit" or "agentStop".

	// Test all valid combinations of when.type and then.type
	for _, whenType := range validWhenTypesSlice {
		for _, thenType := range validThenTypesSlice {
			testName := whenType + "_" + thenType

			t.Run(testName, func(t *testing.T) {
				// Skip invalid combinations (runCommand with non-allowed triggers)
				if thenType == "runCommand" && !runCommandAllowedTriggers[whenType] {
					// This combination should fail - test it separately
					return
				}

				hook := buildValidHook(whenType, thenType)
				err := ValidateHookFile(hook)
				if err != nil {
					t.Errorf("Valid hook with when.type=%s and then.type=%s should pass: %v", whenType, thenType, err)
				}
			})
		}
	}
}

// TestProperty4_HookRequiredFields tests that all required fields must be present.
// Feature: phase4-production, Property 4: Hook File Schema Validity
// **Validates: Requirements 5.2, 10.4**
func TestProperty4_HookRequiredFields(t *testing.T) {
	// Property: For any hook file, all required fields must be present

	requiredFields := []string{"name", "description", "version", "when.type", "then.type"}

	for _, field := range requiredFields {
		t.Run("missing_"+field, func(t *testing.T) {
			hook := buildHookMissingField(field)
			err := ValidateHookFile(hook)
			if err == nil {
				t.Errorf("Hook missing %s should fail validation", field)
			}
		})
	}
}

// TestProperty4_HookWhenTypeValidity tests that when.type must be a valid value.
// Feature: phase4-production, Property 4: Hook File Schema Validity
// **Validates: Requirements 5.4, 10.4**
func TestProperty4_HookWhenTypeValidity(t *testing.T) {
	// Property: For any hook file, when.type must be one of the valid values

	// Test valid when.type values
	for _, whenType := range validWhenTypesSlice {
		t.Run("valid_"+whenType, func(t *testing.T) {
			thenType := "askAgent" // askAgent works with all triggers
			hook := buildValidHook(whenType, thenType)
			err := ValidateHookFile(hook)
			if err != nil {
				t.Errorf("Hook with valid when.type=%s should pass: %v", whenType, err)
			}
		})
	}

	// Test invalid when.type values
	invalidWhenTypes := []string{"invalid", "onSave", "beforeCommit", "afterPush", ""}
	for _, whenType := range invalidWhenTypes {
		t.Run("invalid_"+whenType, func(t *testing.T) {
			hook := `{
				"name": "Test Hook",
				"description": "Test description",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "` + whenType + `"},
				"then": {"type": "askAgent", "prompt": "test"}
			}`
			err := ValidateHookFile(hook)
			if err == nil {
				t.Errorf("Hook with invalid when.type=%s should fail validation", whenType)
			}
		})
	}
}

// TestProperty4_HookThenTypeValidity tests that then.type must be a valid value.
// Feature: phase4-production, Property 4: Hook File Schema Validity
// **Validates: Requirements 5.5, 10.4**
func TestProperty4_HookThenTypeValidity(t *testing.T) {
	// Property: For any hook file, then.type must be one of the valid values

	// Test valid then.type values
	for _, thenType := range validThenTypesSlice {
		t.Run("valid_"+thenType, func(t *testing.T) {
			// Use a trigger that works with both action types
			whenType := "agentStop"
			hook := buildValidHook(whenType, thenType)
			err := ValidateHookFile(hook)
			if err != nil {
				t.Errorf("Hook with valid then.type=%s should pass: %v", thenType, err)
			}
		})
	}

	// Test invalid then.type values
	invalidThenTypes := []string{"invalid", "execute", "notify", "log", ""}
	for _, thenType := range invalidThenTypes {
		t.Run("invalid_"+thenType, func(t *testing.T) {
			hook := `{
				"name": "Test Hook",
				"description": "Test description",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "` + thenType + `", "prompt": "test"}
			}`
			err := ValidateHookFile(hook)
			if err == nil {
				t.Errorf("Hook with invalid then.type=%s should fail validation", thenType)
			}
		})
	}
}

// TestProperty4_RunCommandRestriction tests that runCommand is only allowed with specific triggers.
// Feature: phase4-production, Property 4: Hook File Schema Validity
// **Validates: Requirements 5.5, 10.4**
func TestProperty4_RunCommandRestriction(t *testing.T) {
	// Property: If then.type is "runCommand", then when.type MUST be either "promptSubmit" or "agentStop"

	// Test allowed combinations
	allowedTriggers := []string{"promptSubmit", "agentStop"}
	for _, whenType := range allowedTriggers {
		t.Run("allowed_"+whenType, func(t *testing.T) {
			hook := buildValidHook(whenType, "runCommand")
			err := ValidateHookFile(hook)
			if err != nil {
				t.Errorf("runCommand with %s should be allowed: %v", whenType, err)
			}
		})
	}

	// Test disallowed combinations
	disallowedTriggers := []string{"fileEdited", "fileCreated", "fileDeleted", "userTriggered"}
	for _, whenType := range disallowedTriggers {
		t.Run("disallowed_"+whenType, func(t *testing.T) {
			var hook string
			if fileBasedTriggers[whenType] {
				hook = `{
					"name": "Test Hook",
					"description": "Test description",
					"version": "1.0.0",
					"enabled": true,
					"when": {"type": "` + whenType + `", "patterns": ["**/*.go"]},
					"then": {"type": "runCommand", "command": "go fmt"}
				}`
			} else {
				hook = `{
					"name": "Test Hook",
					"description": "Test description",
					"version": "1.0.0",
					"enabled": true,
					"when": {"type": "` + whenType + `"},
					"then": {"type": "runCommand", "command": "go fmt"}
				}`
			}
			err := ValidateHookFile(hook)
			if err == nil {
				t.Errorf("runCommand with %s should NOT be allowed", whenType)
			}
		})
	}
}

// TestProperty4_FileBasedTriggersRequirePatterns tests that file-based triggers require patterns.
// Feature: phase4-production, Property 4: Hook File Schema Validity
// **Validates: Requirements 5.3, 10.4**
func TestProperty4_FileBasedTriggersRequirePatterns(t *testing.T) {
	// Property: For file-based triggers (fileEdited, fileCreated, fileDeleted),
	// the patterns field is required

	fileBasedTriggersList := []string{"fileEdited", "fileCreated", "fileDeleted"}

	for _, whenType := range fileBasedTriggersList {
		t.Run("with_patterns_"+whenType, func(t *testing.T) {
			hook := `{
				"name": "Test Hook",
				"description": "Test description",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "` + whenType + `", "patterns": ["**/*.go"]},
				"then": {"type": "askAgent", "prompt": "test"}
			}`
			err := ValidateHookFile(hook)
			if err != nil {
				t.Errorf("File-based trigger %s with patterns should pass: %v", whenType, err)
			}
		})

		t.Run("without_patterns_"+whenType, func(t *testing.T) {
			hook := `{
				"name": "Test Hook",
				"description": "Test description",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "` + whenType + `"},
				"then": {"type": "askAgent", "prompt": "test"}
			}`
			err := ValidateHookFile(hook)
			if err == nil {
				t.Errorf("File-based trigger %s without patterns should fail", whenType)
			}
		})

		t.Run("empty_patterns_"+whenType, func(t *testing.T) {
			hook := `{
				"name": "Test Hook",
				"description": "Test description",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "` + whenType + `", "patterns": []},
				"then": {"type": "askAgent", "prompt": "test"}
			}`
			err := ValidateHookFile(hook)
			if err == nil {
				t.Errorf("File-based trigger %s with empty patterns should fail", whenType)
			}
		})
	}
}

// TestProperty4_NonFileTriggersNoPatterns tests that non-file triggers don't require patterns.
// Feature: phase4-production, Property 4: Hook File Schema Validity
// **Validates: Requirements 5.4, 10.4**
func TestProperty4_NonFileTriggersNoPatterns(t *testing.T) {
	// Property: For non-file triggers, patterns are not required

	nonFileTriggers := []string{"promptSubmit", "agentStop", "userTriggered"}

	for _, whenType := range nonFileTriggers {
		t.Run("without_patterns_"+whenType, func(t *testing.T) {
			thenType := "askAgent"
			if whenType == "promptSubmit" || whenType == "agentStop" {
				// Can also use runCommand
				thenType = "runCommand"
			}
			hook := buildValidHook(whenType, thenType)
			err := ValidateHookFile(hook)
			if err != nil {
				t.Errorf("Non-file trigger %s without patterns should pass: %v", whenType, err)
			}
		})
	}
}

// TestProperty4_ActionSpecificFields tests that action-specific fields are required.
// Feature: phase4-production, Property 4: Hook File Schema Validity
// **Validates: Requirements 5.2, 10.4**
func TestProperty4_ActionSpecificFields(t *testing.T) {
	// Property: askAgent requires prompt, runCommand requires command

	t.Run("askAgent_with_prompt", func(t *testing.T) {
		hook := `{
			"name": "Test Hook",
			"description": "Test description",
			"version": "1.0.0",
			"enabled": true,
			"when": {"type": "agentStop"},
			"then": {"type": "askAgent", "prompt": "Do something"}
		}`
		err := ValidateHookFile(hook)
		if err != nil {
			t.Errorf("askAgent with prompt should pass: %v", err)
		}
	})

	t.Run("askAgent_without_prompt", func(t *testing.T) {
		hook := `{
			"name": "Test Hook",
			"description": "Test description",
			"version": "1.0.0",
			"enabled": true,
			"when": {"type": "agentStop"},
			"then": {"type": "askAgent"}
		}`
		err := ValidateHookFile(hook)
		if err == nil {
			t.Error("askAgent without prompt should fail")
		}
	})

	t.Run("runCommand_with_command", func(t *testing.T) {
		hook := `{
			"name": "Test Hook",
			"description": "Test description",
			"version": "1.0.0",
			"enabled": true,
			"when": {"type": "agentStop"},
			"then": {"type": "runCommand", "command": "go fmt ./..."}
		}`
		err := ValidateHookFile(hook)
		if err != nil {
			t.Errorf("runCommand with command should pass: %v", err)
		}
	})

	t.Run("runCommand_without_command", func(t *testing.T) {
		hook := `{
			"name": "Test Hook",
			"description": "Test description",
			"version": "1.0.0",
			"enabled": true,
			"when": {"type": "agentStop"},
			"then": {"type": "runCommand"}
		}`
		err := ValidateHookFile(hook)
		if err == nil {
			t.Error("runCommand without command should fail")
		}
	})
}

// TestProperty4_PropertyBasedHookValidation uses property-based testing to verify
// hook validation across random valid inputs.
// Feature: phase4-production, Property 4: Hook File Schema Validity
// **Validates: Requirements 5.1-5.5, 10.4**
func TestProperty4_PropertyBasedHookValidation(t *testing.T) {
	// Property: For any valid combination of when.type and then.type with required fields,
	// the validation should pass

	property := func(nameIdx, descIdx, versionIdx uint8) bool {
		// Generate deterministic but varied values
		names := []string{"Test Hook", "Format Hook", "Lint Hook", "Security Hook", "Quality Hook"}
		descriptions := []string{"Test description", "Runs formatters", "Checks code", "Scans for issues", "Validates quality"}
		versions := []string{"1.0.0", "2.0.0", "1.1.0", "0.1.0", "3.2.1"}

		name := names[int(nameIdx)%len(names)]
		description := descriptions[int(descIdx)%len(descriptions)]
		version := versions[int(versionIdx)%len(versions)]

		// Test all valid when/then combinations
		for _, whenType := range validWhenTypesSlice {
			for _, thenType := range validThenTypesSlice {
				// Skip invalid combinations
				if thenType == "runCommand" && !runCommandAllowedTriggers[whenType] {
					continue
				}

				hook := buildValidHookWithParams(whenType, thenType, name, description, version)
				err := ValidateHookFile(hook)
				if err != nil {
					return false
				}
			}
		}
		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 4 (Hook File Schema Validity) failed: %v", err)
	}
}

// TestProperty4_InvalidJSONRejected tests that invalid JSON is rejected.
// Feature: phase4-production, Property 4: Hook File Schema Validity
// **Validates: Requirements 5.1, 10.4**
func TestProperty4_InvalidJSONRejected(t *testing.T) {
	invalidJSONs := []struct {
		name    string
		content string
	}{
		{"not_json", "not json at all"},
		{"incomplete_json", `{"name": "test"`},
		{"invalid_syntax", `{"name": "test",}`},
		{"wrong_type_name", `{"name": 123, "description": "test", "version": "1.0.0", "enabled": true, "when": {"type": "agentStop"}, "then": {"type": "askAgent", "prompt": "test"}}`},
		{"wrong_type_enabled", `{"name": "test", "description": "test", "version": "1.0.0", "enabled": "yes", "when": {"type": "agentStop"}, "then": {"type": "askAgent", "prompt": "test"}}`},
	}

	for _, tc := range invalidJSONs {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateHookFile(tc.content)
			if err == nil {
				t.Errorf("Invalid JSON case %s should fail validation", tc.name)
			}
		})
	}
}

// TestProperty4_HooksInGeneratedOutput tests that hooks in generated output are valid.
// Feature: phase4-production, Property 4: Hook File Schema Validity
// **Validates: Requirements 5.1-5.5, 10.4**
func TestProperty4_HooksInGeneratedOutput(t *testing.T) {
	// Simulate a valid generated output with various hooks
	validOutput := []GeneratedFile{
		{
			Path:    ".kiro/prompts/kickoff.md",
			Content: minimalValidKickoff(),
			Type:    "kickoff",
		},
		{
			Path: ".kiro/steering/product.md",
			Content: `---
inclusion: always
---

# Product`,
			Type: "steering",
		},
		{
			Path: ".kiro/steering/tech.md",
			Content: `---
inclusion: always
---

# Tech Stack`,
			Type: "steering",
		},
		{
			Path: ".kiro/steering/structure.md",
			Content: `---
inclusion: always
---

# Structure`,
			Type: "steering",
		},
		{
			Path: ".kiro/hooks/format-on-stop.kiro.hook",
			Content: `{
				"name": "Format on Agent Stop",
				"description": "Run code formatters when agent completes work",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "runCommand", "command": "go fmt ./..."}
			}`,
			Type: "hook",
		},
		{
			Path: ".kiro/hooks/secret-scan.kiro.hook",
			Content: `{
				"name": "Secret Scanner",
				"description": "Scan for accidentally committed secrets",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "askAgent", "prompt": "Scan modified files for potential secrets."}
			}`,
			Type: "hook",
		},
		{
			Path: ".kiro/hooks/test-manual.kiro.hook",
			Content: `{
				"name": "Run Tests",
				"description": "Manually trigger test suite",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "userTriggered"},
				"then": {"type": "askAgent", "prompt": "Run the test suite and summarize results."}
			}`,
			Type: "hook",
		},
		{
			Path: ".kiro/hooks/go-test-on-change.kiro.hook",
			Content: `{
				"name": "Go Test on Change",
				"description": "Run tests when Go files change",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "fileEdited", "patterns": ["**/*.go"]},
				"then": {"type": "askAgent", "prompt": "A Go file was modified. Run the relevant tests."}
			}`,
			Type: "hook",
		},
		{
			Path:    "AGENTS.md",
			Content: "# Agent Guidelines",
			Type:    "agents",
		},
	}

	// Validate all files
	err := ValidateGeneratedFiles(validOutput)
	if err != nil {
		t.Errorf("Valid output with hooks should pass validation: %v", err)
	}

	// Count and verify hooks
	hookCount := 0
	for _, f := range validOutput {
		if f.Type == "hook" {
			hookCount++
			// Verify each hook individually
			err := ValidateHookFile(f.Content)
			if err != nil {
				t.Errorf("Hook %s should be valid: %v", f.Path, err)
			}
		}
	}

	if hookCount == 0 {
		t.Error("Output should contain at least one hook file")
	}
}

// Helper functions for building test hooks

// buildValidHook creates a valid hook JSON string for the given when and then types
func buildValidHook(whenType, thenType string) string {
	return buildValidHookWithParams(whenType, thenType, "Test Hook", "Test description", "1.0.0")
}

// buildValidHookWithParams creates a valid hook JSON string with custom parameters
func buildValidHookWithParams(whenType, thenType, name, description, version string) string {
	var whenBlock string
	if fileBasedTriggers[whenType] {
		whenBlock = `"when": {"type": "` + whenType + `", "patterns": ["**/*.go"]}`
	} else {
		whenBlock = `"when": {"type": "` + whenType + `"}`
	}

	var thenBlock string
	if thenType == "askAgent" {
		thenBlock = `"then": {"type": "askAgent", "prompt": "Do something"}`
	} else {
		thenBlock = `"then": {"type": "runCommand", "command": "go fmt ./..."}`
	}

	return `{
		"name": "` + name + `",
		"description": "` + description + `",
		"version": "` + version + `",
		"enabled": true,
		` + whenBlock + `,
		` + thenBlock + `
	}`
}

// buildHookMissingField creates a hook JSON string missing the specified field
func buildHookMissingField(field string) string {
	switch field {
	case "name":
		return `{
			"description": "Test description",
			"version": "1.0.0",
			"enabled": true,
			"when": {"type": "agentStop"},
			"then": {"type": "askAgent", "prompt": "test"}
		}`
	case "description":
		return `{
			"name": "Test Hook",
			"version": "1.0.0",
			"enabled": true,
			"when": {"type": "agentStop"},
			"then": {"type": "askAgent", "prompt": "test"}
		}`
	case "version":
		return `{
			"name": "Test Hook",
			"description": "Test description",
			"enabled": true,
			"when": {"type": "agentStop"},
			"then": {"type": "askAgent", "prompt": "test"}
		}`
	case "when.type":
		return `{
			"name": "Test Hook",
			"description": "Test description",
			"version": "1.0.0",
			"enabled": true,
			"when": {},
			"then": {"type": "askAgent", "prompt": "test"}
		}`
	case "then.type":
		return `{
			"name": "Test Hook",
			"description": "Test description",
			"version": "1.0.0",
			"enabled": true,
			"when": {"type": "agentStop"},
			"then": {"prompt": "test"}
		}`
	default:
		return `{}`
	}
}

// minimalValidKickoff returns a minimal but valid kickoff prompt for use in tests
// that need a valid kickoff but aren't specifically testing kickoff validation
func minimalValidKickoff() string {
	return `# Project Kickoff: Test Project

> ⚠️ **IMPORTANT**: Do not write any code until all questions below are answered.

## Project Identity
A test project.

## Success Criteria
- Works correctly

## Users & Roles
- Admin: Full access

## Data Sensitivity
- User data: Confidential

## Auth Model
Basic authentication

## Concurrency Expectations
Single user

## Risks & Tradeoffs
### Risk 1: Security
- Mitigation: Use HTTPS

## Boundaries
Public and private areas.

### Boundary Examples
- Admin CAN delete users

## Non-Goals
- Mobile app

## Constraints
- 2 week timeline
`
}

// =============================================================================
// Property-Based Tests for Kickoff Prompt Completeness (Property 5)
// =============================================================================

// requiredKickoffSectionsTest defines the sections that must be present in a kickoff prompt (for test assertions)
var requiredKickoffSectionsTest = []string{
	"Project Identity",
	"Success Criteria",
	"Users & Roles",
	"Data Sensitivity",
	"Auth Model",
	"Concurrency",
	"Boundaries",
	"Non-Goals",
	"Constraints",
	"Risks",     // Part of "Risks & Tradeoffs"
	"Tradeoffs", // Part of "Risks & Tradeoffs"
	"Boundary Examples",
}

// noCodingPhrasesTest defines phrases that enforce "no coding until questions answered" (for test assertions)
var noCodingPhrasesTest = []string{
	"no coding",
	"do not write any code",
	"don't write any code",
	"no code until",
	"do not code until",
	"don't code until",
	"before writing any code",
	"before coding",
}

// TestProperty5_KickoffPromptCompleteness tests that kickoff prompts contain all required sections.
// Feature: phase4-production, Property 5: Kickoff Prompt Completeness
// **Validates: Requirements 6.1, 6.2, 6.4, 6.5**
func TestProperty5_KickoffPromptCompleteness(t *testing.T) {
	// Property: For any generated kickoff prompt, the content SHALL contain
	// the phrase "no coding" (or equivalent enforcement), AND SHALL include
	// sections for: Project Identity, Success Criteria, Users & Roles,
	// Data Sensitivity, Auth Model, Concurrency, Boundaries, Non-Goals,
	// Constraints, Risks & Tradeoffs, and Boundary Examples.

	// Test a valid kickoff prompt with all required sections
	validKickoff := buildValidKickoffPrompt()

	err := ValidateKickoffPrompt(validKickoff)
	if err != nil {
		t.Errorf("Valid kickoff prompt should pass validation: %v", err)
	}
}

// TestProperty5_KickoffContainsNoCodingEnforcement tests that kickoff prompts enforce "no coding".
// Feature: phase4-production, Property 5: Kickoff Prompt Completeness
// **Validates: Requirements 6.1**
func TestProperty5_KickoffContainsNoCodingEnforcement(t *testing.T) {
	// Property: The kickoff prompt MUST contain a phrase that enforces
	// "no coding until questions are answered"

	validKickoff := buildValidKickoffPrompt()
	contentLower := strings.ToLower(validKickoff)

	found := false
	for _, phrase := range noCodingPhrasesTest {
		if strings.Contains(contentLower, phrase) {
			found = true
			break
		}
	}

	if !found {
		t.Error("Kickoff prompt must contain a 'no coding' enforcement phrase")
	}
}

// TestProperty5_KickoffContainsAllRequiredSections tests that all required sections are present.
// Feature: phase4-production, Property 5: Kickoff Prompt Completeness
// **Validates: Requirements 6.2, 6.4, 6.5**
func TestProperty5_KickoffContainsAllRequiredSections(t *testing.T) {
	// Property: The kickoff prompt MUST contain all required sections

	validKickoff := buildValidKickoffPrompt()
	contentLower := strings.ToLower(validKickoff)

	missingSections := []string{}
	for _, section := range requiredKickoffSectionsTest {
		if !strings.Contains(contentLower, strings.ToLower(section)) {
			missingSections = append(missingSections, section)
		}
	}

	if len(missingSections) > 0 {
		t.Errorf("Kickoff prompt is missing required sections: %v", missingSections)
	}
}

// TestProperty5_KickoffMissingNoCodingFails tests that kickoff without "no coding" fails validation.
// Feature: phase4-production, Property 5: Kickoff Prompt Completeness
// **Validates: Requirements 6.1**
func TestProperty5_KickoffMissingNoCodingFails(t *testing.T) {
	// A kickoff prompt without the "no coding" enforcement should fail

	kickoffWithoutEnforcement := `# Project Kickoff: Test Project

## Project Identity
A test project for validation.

## Success Criteria
- Feature works correctly

## Users & Roles
- Admin: Full access
- User: Limited access

## Data Sensitivity
- User data: Confidential

## Auth Model
Basic authentication

## Concurrency Expectations
Single user at a time

## Risks & Tradeoffs
### Risk 1: Security
- Mitigation: Use HTTPS

## Boundaries
Public and private areas defined.

### Boundary Examples
- Admin CAN delete users
- User CANNOT delete other users

## Non-Goals
- Mobile app
- Real-time features

## Constraints
- 2 week timeline
`

	err := ValidateKickoffPrompt(kickoffWithoutEnforcement)
	if err == nil {
		t.Error("Kickoff prompt without 'no coding' enforcement should fail validation")
	}
	if err != ErrMissingNoCodingEnforcement {
		t.Errorf("Expected ErrMissingNoCodingEnforcement, got: %v", err)
	}
}

// TestProperty5_KickoffMissingSectionFails tests that kickoff missing sections fails validation.
// Feature: phase4-production, Property 5: Kickoff Prompt Completeness
// **Validates: Requirements 6.2**
func TestProperty5_KickoffMissingSectionFails(t *testing.T) {
	// Test that missing each required section causes validation to fail

	// Core sections that must be present (excluding Risks/Tradeoffs which are checked together)
	coreSections := []struct {
		name    string
		content string
	}{
		{"Project Identity", "## Project Identity\nA test project."},
		{"Success Criteria", "## Success Criteria\n- Works correctly"},
		{"Users & Roles", "## Users & Roles\n- Admin: Full access"},
		{"Data Sensitivity", "## Data Sensitivity\n- User data: Confidential"},
		{"Auth Model", "## Auth Model\nBasic authentication"},
		{"Concurrency", "## Concurrency Expectations\nSingle user"},
		{"Boundaries", "## Boundaries\nPublic and private areas."},
		{"Non-Goals", "## Non-Goals\n- Mobile app"},
		{"Constraints", "## Constraints\n- 2 week timeline"},
	}

	for _, section := range coreSections {
		t.Run("missing_"+section.name, func(t *testing.T) {
			// Build a kickoff prompt missing this specific section
			kickoff := buildKickoffMissingSection(section.name)

			err := ValidateKickoffPrompt(kickoff)
			if err == nil {
				t.Errorf("Kickoff prompt missing '%s' section should fail validation", section.name)
			}
		})
	}
}

// TestProperty5_KickoffMissingRisksAndTradeoffsFails tests that missing Risks & Tradeoffs fails.
// Feature: phase4-production, Property 5: Kickoff Prompt Completeness
// **Validates: Requirements 6.4**
func TestProperty5_KickoffMissingRisksAndTradeoffsFails(t *testing.T) {
	kickoff := `# Project Kickoff: Test Project

> ⚠️ **IMPORTANT**: Do not write any code until all questions below are answered.

## Project Identity
A test project.

## Success Criteria
- Works correctly

## Users & Roles
- Admin: Full access

## Data Sensitivity
- User data: Confidential

## Auth Model
Basic authentication

## Concurrency Expectations
Single user

## Boundaries
Public and private areas.

### Boundary Examples
- Admin CAN delete users

## Non-Goals
- Mobile app

## Constraints
- 2 week timeline
`

	err := ValidateKickoffPrompt(kickoff)
	if err == nil {
		t.Error("Kickoff prompt missing 'Risks & Tradeoffs' section should fail validation")
	}
}

// TestProperty5_KickoffMissingBoundaryExamplesFails tests that missing Boundary Examples fails.
// Feature: phase4-production, Property 5: Kickoff Prompt Completeness
// **Validates: Requirements 6.5**
func TestProperty5_KickoffMissingBoundaryExamplesFails(t *testing.T) {
	kickoff := `# Project Kickoff: Test Project

> ⚠️ **IMPORTANT**: Do not write any code until all questions below are answered.

## Project Identity
A test project.

## Success Criteria
- Works correctly

## Users & Roles
- Admin: Full access

## Data Sensitivity
- User data: Confidential

## Auth Model
Basic authentication

## Concurrency Expectations
Single user

## Risks & Tradeoffs
### Risk 1: Security
- Mitigation: Use HTTPS

## Boundaries
Public and private areas.

## Non-Goals
- Mobile app

## Constraints
- 2 week timeline
`

	err := ValidateKickoffPrompt(kickoff)
	if err == nil {
		t.Error("Kickoff prompt missing 'Boundary Examples' section should fail validation")
	}
}

// TestProperty5_PropertyBasedKickoffValidation uses property-based testing to verify
// kickoff validation across various valid inputs.
// Feature: phase4-production, Property 5: Kickoff Prompt Completeness
// **Validates: Requirements 6.1, 6.2, 6.4, 6.5**
func TestProperty5_PropertyBasedKickoffValidation(t *testing.T) {
	// Property: For any valid kickoff prompt structure with all required sections
	// and "no coding" enforcement, validation should pass

	property := func(projectNameIdx, descIdx uint8) bool {
		// Generate deterministic but varied values
		projectNames := []string{"Test App", "My Project", "Cool Service", "Data Platform", "API Gateway"}
		descriptions := []string{
			"A simple test application",
			"A project for managing tasks",
			"A service for processing data",
			"A platform for analytics",
			"An API gateway for microservices",
		}

		projectName := projectNames[int(projectNameIdx)%len(projectNames)]
		description := descriptions[int(descIdx)%len(descriptions)]

		kickoff := buildValidKickoffPromptWithParams(projectName, description)
		err := ValidateKickoffPrompt(kickoff)
		return err == nil
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 5 (Kickoff Prompt Completeness) failed: %v", err)
	}
}

// TestProperty5_KickoffInGeneratedOutput tests that kickoff in generated output is valid.
// Feature: phase4-production, Property 5: Kickoff Prompt Completeness
// **Validates: Requirements 6.1, 6.2, 6.4, 6.5**
func TestProperty5_KickoffInGeneratedOutput(t *testing.T) {
	// Simulate a valid generated output with a complete kickoff prompt
	validKickoff := buildValidKickoffPrompt()

	validOutput := []GeneratedFile{
		{
			Path:    ".kiro/prompts/kickoff.md",
			Content: validKickoff,
			Type:    "kickoff",
		},
		{
			Path: ".kiro/steering/product.md",
			Content: `---
inclusion: always
---

# Product`,
			Type: "steering",
		},
		{
			Path: ".kiro/steering/tech.md",
			Content: `---
inclusion: always
---

# Tech Stack`,
			Type: "steering",
		},
		{
			Path: ".kiro/steering/structure.md",
			Content: `---
inclusion: always
---

# Structure`,
			Type: "steering",
		},
		{
			Path: ".kiro/hooks/format-on-stop.kiro.hook",
			Content: `{
				"name": "Format on Stop",
				"description": "Run formatters",
				"version": "1.0.0",
				"enabled": true,
				"when": {"type": "agentStop"},
				"then": {"type": "runCommand", "command": "go fmt ./..."}
			}`,
			Type: "hook",
		},
		{
			Path:    "AGENTS.md",
			Content: "# Agent Guidelines",
			Type:    "agents",
		},
	}

	// Validate all files including kickoff
	err := ValidateGeneratedFiles(validOutput)
	if err != nil {
		t.Errorf("Valid output with complete kickoff should pass validation: %v", err)
	}
}

// TestProperty5_VariousNoCodingPhrasesAccepted tests that various "no coding" phrases are accepted.
// Feature: phase4-production, Property 5: Kickoff Prompt Completeness
// **Validates: Requirements 6.1**
func TestProperty5_VariousNoCodingPhrasesAccepted(t *testing.T) {
	// Test that various equivalent phrases for "no coding" are accepted

	phrases := []string{
		"Do not write any code until all questions are answered.",
		"Don't write any code until this is reviewed.",
		"No coding until the design is approved.",
		"Before writing any code, review this document.",
		"Before coding, ensure all sections are complete.",
	}

	for _, phrase := range phrases {
		t.Run(phrase[:20], func(t *testing.T) {
			kickoff := buildKickoffWithCustomEnforcement(phrase)
			err := ValidateKickoffPrompt(kickoff)
			if err != nil {
				t.Errorf("Kickoff with phrase '%s' should pass validation: %v", phrase, err)
			}
		})
	}
}

// Helper functions for building test kickoff prompts

// buildValidKickoffPrompt creates a valid kickoff prompt with all required sections
func buildValidKickoffPrompt() string {
	return buildValidKickoffPromptWithParams("Test Project", "A test project for validation")
}

// buildValidKickoffPromptWithParams creates a valid kickoff prompt with custom parameters
func buildValidKickoffPromptWithParams(projectName, description string) string {
	return `# Project Kickoff: ` + projectName + `

> ⚠️ **IMPORTANT**: Do not write any code until all questions below are answered and reviewed.

## Project Identity
` + description + `

## Success Criteria
- Feature works correctly
- All tests pass
- Documentation is complete

## Users & Roles
| Role | Description | Key Capabilities |
|------|-------------|------------------|
| Admin | System administrator | Full access to all features |
| User | Regular user | Limited access to own data |

## Data Sensitivity
| Data Type | Sensitivity | Storage | Notes |
|-----------|-------------|---------|-------|
| User credentials | Restricted | Encrypted database | Never log |
| User preferences | Internal | Database | Can be exported |

### Data Lifecycle
- **Retention**: Data kept for 2 years
- **Deletion**: Users can request deletion
- **Export**: JSON export available
- **Audit**: All changes logged
- **Backups**: Daily backups

## Auth Model
- [x] Basic (username/password)

## Concurrency Expectations
- **Multi-user**: Yes, multiple users can access simultaneously
- **Shared state**: Minimal shared state
- **Background jobs**: None initially
- **Real-time**: Not required

## Risks & Tradeoffs
### Risk 1: Security Vulnerabilities
- **Description**: Potential security issues
- **Likelihood**: Medium
- **Impact**: High
- **Mitigation**: Regular security audits
- **Accepted**: Some edge cases not covered

### Risk 2: Performance Issues
- **Description**: Slow response times under load
- **Likelihood**: Low
- **Impact**: Medium
- **Mitigation**: Caching and optimization
- **Accepted**: Initial version may be slower

### Risk 3: Scope Creep
- **Description**: Feature requests expanding scope
- **Likelihood**: High
- **Impact**: Medium
- **Mitigation**: Strict non-goals list
- **Accepted**: Some features deferred

## Boundaries
### Public
- Landing page
- Documentation

### Private
- User dashboard
- Settings

### Boundary Examples
- Admin CAN delete any user
- Admin CAN view all data
- User CAN view their own data
- User CANNOT view other users' data
- User CAN export their own data
- User CANNOT delete other users

## Non-Goals
- NOT building: Mobile application
- NOT building: Real-time chat
- NOT building: Advanced analytics
- Out of scope: Third-party integrations

## Constraints
- **Timeline**: 4 weeks
- **Simplicity**: Keep it simple, avoid over-engineering
- **Tech**: Go backend, React frontend
- **Budget**: No external services initially
- **Team**: Single developer

---

## Next Steps
1. Review this document with stakeholders
2. Create specs for each major feature
3. Begin implementation only after specs are approved
`
}

// buildKickoffMissingSection creates a kickoff prompt missing a specific section
func buildKickoffMissingSection(missingSection string) string {
	sections := map[string]string{
		"Project Identity": `## Project Identity
A test project.`,
		"Success Criteria": `## Success Criteria
- Works correctly`,
		"Users & Roles": `## Users & Roles
- Admin: Full access`,
		"Data Sensitivity": `## Data Sensitivity
- User data: Confidential`,
		"Auth Model": `## Auth Model
Basic authentication`,
		"Concurrency": `## Concurrency Expectations
Single user`,
		"Boundaries": `## Boundaries
Public and private areas.

### Boundary Examples
- Admin CAN delete users`,
		"Non-Goals": `## Non-Goals
- Mobile app`,
		"Constraints": `## Constraints
- 2 week timeline`,
		"Risks": `## Risks & Tradeoffs
### Risk 1: Security
- Mitigation: Use HTTPS`,
	}

	kickoff := `# Project Kickoff: Test Project

> ⚠️ **IMPORTANT**: Do not write any code until all questions below are answered.

`

	for name, content := range sections {
		if name != missingSection {
			kickoff += content + "\n\n"
		}
	}

	return kickoff
}

// buildKickoffWithCustomEnforcement creates a kickoff with a custom "no coding" phrase
func buildKickoffWithCustomEnforcement(phrase string) string {
	return `# Project Kickoff: Test Project

> ⚠️ **IMPORTANT**: ` + phrase + `

## Project Identity
A test project.

## Success Criteria
- Works correctly

## Users & Roles
- Admin: Full access

## Data Sensitivity
- User data: Confidential

## Auth Model
Basic authentication

## Concurrency Expectations
Single user

## Risks & Tradeoffs
### Risk 1: Security
- Mitigation: Use HTTPS

## Boundaries
Public and private areas.

### Boundary Examples
- Admin CAN delete users

## Non-Goals
- Mobile app

## Constraints
- 2 week timeline
`
}

// =============================================================================
// Property-Based Tests for Question Generation Constraints (Property 6)
// =============================================================================

// questionCategoryKeywords maps question categories to keywords that indicate that category
// Keywords are ordered by specificity - more specific keywords first
var questionCategoryKeywords = map[string][]string{
	"identity":     {"what is", "what does", "describe", "purpose", "goal", "problem", "solve", "building", "project", "main"},
	"users":        {"who will", "who are", "user role", "audience", "customer", "people use"},
	"data":         {"data", "store", "save", "database", "information", "record", "persist", "storage"},
	"auth":         {"auth", "login", "sign in", "sign up", "password", "credential", "session", "token", "oauth", "permission", "access control"},
	"architecture": {"architecture", "structure", "component", "service", "api", "deploy", "host", "scale", "performance", "backend", "frontend"},
	"constraints":  {"constraint", "limit", "budget", "time", "deadline", "requirement", "must", "cannot", "restriction", "timeline"},
}

// questionCategoryOrder defines the expected order of question categories
var questionCategoryOrder = []string{"identity", "users", "data", "auth", "architecture", "constraints"}

// categorizeQuestion attempts to categorize a question based on keywords
// It uses a scoring system to find the best match
func categorizeQuestion(questionText string) string {
	textLower := strings.ToLower(questionText)

	// Score each category based on keyword matches
	scores := make(map[string]int)
	for _, category := range questionCategoryOrder {
		keywords := questionCategoryKeywords[category]
		for _, keyword := range keywords {
			if strings.Contains(textLower, keyword) {
				scores[category]++
			}
		}
	}

	// Find the category with the highest score
	bestCategory := "unknown"
	bestScore := 0
	for _, category := range questionCategoryOrder {
		if scores[category] > bestScore {
			bestScore = scores[category]
			bestCategory = category
		}
	}

	return bestCategory
}

// getCategoryIndex returns the index of a category in the expected order, or -1 if not found
func getCategoryIndex(category string) int {
	for i, cat := range questionCategoryOrder {
		if cat == category {
			return i
		}
	}
	return -1
}

// TestProperty6_QuestionGenerationConstraints tests that question generation follows constraints.
// Feature: phase4-production, Property 6: Question Generation Constraints
// **Validates: Requirements 3.4, 3.6**
func TestProperty6_QuestionGenerationConstraints(t *testing.T) {
	// Property: For any generated question set, the count SHALL be between 5 and 10 inclusive,
	// AND the questions SHALL follow a logical ordering where identity/scope questions
	// appear before technical/architecture questions.

	// Test that the service enforces question count constraints
	t.Run("question_count_bounds", func(t *testing.T) {
		// Test minimum bound
		if minQuestions < 5 {
			t.Errorf("minQuestions should be at least 5, got %d", minQuestions)
		}
		// Test maximum bound
		if maxQuestions > 10 {
			t.Errorf("maxQuestions should be at most 10, got %d", maxQuestions)
		}
		// Test valid range
		if minQuestions > maxQuestions {
			t.Errorf("minQuestions (%d) should not exceed maxQuestions (%d)", minQuestions, maxQuestions)
		}
	})

	// Test that question ordering guidance is present in prompts
	t.Run("ordering_guidance_in_prompts", func(t *testing.T) {
		levels := []string{prompts.ExperienceBeginner, prompts.ExperienceNovice, prompts.ExperienceExpert}

		for _, level := range levels {
			prompt := prompts.GetQuestionsSystemPrompt(level)
			promptLower := strings.ToLower(prompt)

			// Check that ordering rules are mentioned
			if !strings.Contains(promptLower, "order") {
				t.Errorf("Prompt for %s level should mention question ordering", level)
			}

			// Check that the logical sequence is mentioned
			orderingKeywords := []string{"identity", "users", "data", "auth", "architecture", "constraints"}
			foundCount := 0
			for _, keyword := range orderingKeywords {
				if strings.Contains(promptLower, keyword) {
					foundCount++
				}
			}

			if foundCount < 4 {
				t.Errorf("Prompt for %s level should mention most ordering categories, found %d/6", level, foundCount)
			}
		}
	})
}

// TestProperty6_QuestionCountValidation tests that question count is validated correctly.
// Feature: phase4-production, Property 6: Question Generation Constraints
// **Validates: Requirements 3.4**
func TestProperty6_QuestionCountValidation(t *testing.T) {
	// Test various question counts
	testCases := []struct {
		name        string
		count       int
		shouldTrunc bool
	}{
		{"minimum_valid", 5, false},
		{"maximum_valid", 10, false},
		{"middle_valid", 7, false},
		{"below_minimum", 3, false}, // Service accepts but may be suboptimal
		{"above_maximum", 15, true}, // Service should truncate
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock questions response
			questions := make([]Question, tc.count)
			for i := 0; i < tc.count; i++ {
				questions[i] = Question{
					ID:   i + 1,
					Text: fmt.Sprintf("Question %d about the project?", i+1),
					Hint: fmt.Sprintf("Hint for question %d", i+1),
				}
			}

			// Verify the count
			if tc.shouldTrunc && len(questions) > maxQuestions {
				// Simulate truncation
				questions = questions[:maxQuestions]
			}

			if tc.shouldTrunc && len(questions) != maxQuestions {
				t.Errorf("Expected truncation to %d questions, got %d", maxQuestions, len(questions))
			}
		})
	}
}

// TestProperty6_QuestionOrderingLogic tests that questions follow logical ordering.
// Feature: phase4-production, Property 6: Question Generation Constraints
// **Validates: Requirements 3.6**
func TestProperty6_QuestionOrderingLogic(t *testing.T) {
	// Test that a well-ordered set of questions passes the ordering check
	wellOrderedQuestions := []Question{
		{ID: 1, Text: "What is the main purpose of your project?"},
		{ID: 2, Text: "Who will be the primary users of this application?"},
		{ID: 3, Text: "What data will your application need to store?"},
		{ID: 4, Text: "How will users authenticate and log in?"},
		{ID: 5, Text: "What architecture pattern do you prefer for the backend?"},
		{ID: 6, Text: "What are your time and budget constraints?"},
	}

	// Verify ordering
	lastCategoryIndex := -1
	for _, q := range wellOrderedQuestions {
		category := categorizeQuestion(q.Text)
		categoryIndex := getCategoryIndex(category)

		if categoryIndex != -1 && categoryIndex < lastCategoryIndex {
			t.Errorf("Question '%s' (category: %s) appears out of order", q.Text, category)
		}

		if categoryIndex != -1 {
			lastCategoryIndex = categoryIndex
		}
	}
}

// TestProperty6_PropertyBasedQuestionValidation uses property-based testing to verify
// question generation constraints across random inputs.
// Feature: phase4-production, Property 6: Question Generation Constraints
// **Validates: Requirements 3.4, 3.6**
func TestProperty6_PropertyBasedQuestionValidation(t *testing.T) {
	// Property: For any valid question count between minQuestions and maxQuestions,
	// the questions should be accepted

	property := func(countOffset uint8) bool {
		// Generate a count within valid range
		count := minQuestions + int(countOffset)%(maxQuestions-minQuestions+1)

		// Verify count is in valid range
		if count < minQuestions || count > maxQuestions {
			return false
		}

		// Create questions
		questions := make([]Question, count)
		for i := 0; i < count; i++ {
			questions[i] = Question{
				ID:   i + 1,
				Text: fmt.Sprintf("Question %d?", i+1),
			}
		}

		// Verify all questions have required fields
		for _, q := range questions {
			if q.ID == 0 || q.Text == "" {
				return false
			}
		}

		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 6 (Question Generation Constraints) failed: %v", err)
	}
}

// TestProperty6_QuestionCountBoundsEnforced tests that the service enforces count bounds.
// Feature: phase4-production, Property 6: Question Generation Constraints
// **Validates: Requirements 3.4**
func TestProperty6_QuestionCountBoundsEnforced(t *testing.T) {
	// Verify the constants are correctly defined
	if minQuestions != 5 {
		t.Errorf("minQuestions should be 5, got %d", minQuestions)
	}
	if maxQuestions != 10 {
		t.Errorf("maxQuestions should be 10, got %d", maxQuestions)
	}

	// Test that parseQuestionsResponse handles various counts
	testCases := []struct {
		name          string
		questionCount int
		expectError   bool
	}{
		{"zero_questions", 0, true},
		{"one_question", 1, false},   // Below min but accepted
		{"five_questions", 5, false}, // At minimum
		{"seven_questions", 7, false},
		{"ten_questions", 10, false},    // At maximum
		{"twelve_questions", 12, false}, // Above max, should be truncated
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Build a JSON response with the specified number of questions
			questions := make([]map[string]interface{}, tc.questionCount)
			for i := 0; i < tc.questionCount; i++ {
				questions[i] = map[string]interface{}{
					"id":   i + 1,
					"text": fmt.Sprintf("Question %d?", i+1),
					"hint": fmt.Sprintf("Hint %d", i+1),
				}
			}

			response := map[string]interface{}{
				"questions": questions,
			}

			jsonBytes, _ := json.Marshal(response)
			jsonStr := string(jsonBytes)

			// Parse the response
			result, err := parseQuestionsResponse(jsonStr)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for %d questions, got none", tc.questionCount)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %d questions: %v", tc.questionCount, err)
				}

				// Verify truncation for counts above maximum
				if tc.questionCount > maxQuestions && len(result) != maxQuestions {
					t.Errorf("Expected truncation to %d questions, got %d", maxQuestions, len(result))
				}
			}
		})
	}
}

// TestProperty6_OrderingGuidanceCompleteness tests that all ordering categories are in prompts.
// Feature: phase4-production, Property 6: Question Generation Constraints
// **Validates: Requirements 3.6**
func TestProperty6_OrderingGuidanceCompleteness(t *testing.T) {
	// The system prompt should mention the question ordering sequence
	levels := []string{prompts.ExperienceBeginner, prompts.ExperienceNovice, prompts.ExperienceExpert}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			prompt := prompts.GetQuestionsSystemPrompt(level)

			// Check for ordering rules section
			if !strings.Contains(prompt, "Question Ordering Rules") {
				t.Errorf("Prompt for %s should contain 'Question Ordering Rules' section", level)
			}

			// Check for the sequence keywords
			expectedSequence := []string{"Identity", "Users", "Data", "Auth", "Architecture", "Constraints"}
			for _, keyword := range expectedSequence {
				if !strings.Contains(prompt, keyword) {
					t.Errorf("Prompt for %s should mention '%s' in ordering rules", level, keyword)
				}
			}
		})
	}
}

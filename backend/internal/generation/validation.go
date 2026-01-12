package generation

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Validation errors
var (
	ErrInvalidFrontmatter      = errors.New("invalid frontmatter")
	ErrMissingInclusion        = errors.New("missing inclusion field in frontmatter")
	ErrInvalidInclusionMode    = errors.New("invalid inclusion mode")
	ErrMissingFileMatchPattern = errors.New("fileMatch mode requires fileMatchPattern")
	ErrInvalidHookSchema       = errors.New("invalid hook schema")
	ErrMissingHookField        = errors.New("missing required hook field")
	ErrInvalidWhenType         = errors.New("invalid when.type value")
	ErrInvalidThenType         = errors.New("invalid then.type value")
	ErrRunCommandRestriction   = errors.New("runCommand can only be used with promptSubmit or agentStop triggers")
)

// Valid inclusion modes for steering files
var validInclusionModes = map[string]bool{
	"always":    true,
	"fileMatch": true,
	"manual":    true,
}

// Valid when.type values for hooks
var validWhenTypes = map[string]bool{
	"fileEdited":    true,
	"fileCreated":   true,
	"fileDeleted":   true,
	"promptSubmit":  true,
	"agentStop":     true,
	"userTriggered": true,
}

// Valid then.type values for hooks
var validThenTypes = map[string]bool{
	"askAgent":   true,
	"runCommand": true,
}

// When types that allow runCommand
var runCommandAllowedWhenTypes = map[string]bool{
	"promptSubmit": true,
	"agentStop":    true,
}

// frontmatterRegex matches YAML frontmatter at the start of a file
var frontmatterRegex = regexp.MustCompile(`(?s)^---\s*\n(.*?)\n---`)

// SteeringFrontmatter represents the parsed frontmatter of a steering file
type SteeringFrontmatter struct {
	Inclusion        string `yaml:"inclusion"`
	FileMatchPattern string `yaml:"fileMatchPattern"`
}

// ValidateSteeringFile validates a steering file's frontmatter
func ValidateSteeringFile(content string) error {
	// Extract frontmatter
	matches := frontmatterRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return fmt.Errorf("%w: no frontmatter found", ErrInvalidFrontmatter)
	}

	frontmatter := matches[1]

	// Parse inclusion mode
	inclusion := extractYAMLField(frontmatter, "inclusion")
	if inclusion == "" {
		return ErrMissingInclusion
	}

	if !validInclusionModes[inclusion] {
		return fmt.Errorf("%w: got '%s', expected 'always', 'fileMatch', or 'manual'", ErrInvalidInclusionMode, inclusion)
	}

	// If fileMatch mode, require fileMatchPattern
	if inclusion == "fileMatch" {
		pattern := extractYAMLField(frontmatter, "fileMatchPattern")
		if pattern == "" {
			return ErrMissingFileMatchPattern
		}
	}

	return nil
}

// extractYAMLField extracts a simple string field from YAML content
func extractYAMLField(yaml, field string) string {
	// Simple regex-based extraction for single-line string values
	pattern := regexp.MustCompile(fmt.Sprintf(`(?m)^%s:\s*["']?([^"'\n]+)["']?\s*$`, regexp.QuoteMeta(field)))
	matches := pattern.FindStringSubmatch(yaml)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// HookFile represents the structure of a Kiro hook file
type HookFile struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Enabled     bool     `json:"enabled"`
	When        HookWhen `json:"when"`
	Then        HookThen `json:"then"`
}

// HookWhen represents the trigger configuration
type HookWhen struct {
	Type     string   `json:"type"`
	Patterns []string `json:"patterns,omitempty"`
}

// HookThen represents the action configuration
type HookThen struct {
	Type    string `json:"type"`
	Prompt  string `json:"prompt,omitempty"`
	Command string `json:"command,omitempty"`
}

// ValidateHookFile validates a hook file's JSON schema
func ValidateHookFile(content string) error {
	var hook HookFile
	if err := json.Unmarshal([]byte(content), &hook); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidHookSchema, err)
	}

	// Validate required fields
	if hook.Name == "" {
		return fmt.Errorf("%w: name", ErrMissingHookField)
	}
	if hook.Description == "" {
		return fmt.Errorf("%w: description", ErrMissingHookField)
	}
	if hook.Version == "" {
		return fmt.Errorf("%w: version", ErrMissingHookField)
	}

	// Validate when.type
	if hook.When.Type == "" {
		return fmt.Errorf("%w: when.type", ErrMissingHookField)
	}
	if !validWhenTypes[hook.When.Type] {
		return fmt.Errorf("%w: got '%s'", ErrInvalidWhenType, hook.When.Type)
	}

	// File-based triggers require patterns
	if isFileBasedTrigger(hook.When.Type) && len(hook.When.Patterns) == 0 {
		return fmt.Errorf("%w: patterns required for %s trigger", ErrMissingHookField, hook.When.Type)
	}

	// Validate then.type
	if hook.Then.Type == "" {
		return fmt.Errorf("%w: then.type", ErrMissingHookField)
	}
	if !validThenTypes[hook.Then.Type] {
		return fmt.Errorf("%w: got '%s'", ErrInvalidThenType, hook.Then.Type)
	}

	// Validate runCommand restriction
	if hook.Then.Type == "runCommand" && !runCommandAllowedWhenTypes[hook.When.Type] {
		return ErrRunCommandRestriction
	}

	// Validate action-specific fields
	if hook.Then.Type == "askAgent" && hook.Then.Prompt == "" {
		return fmt.Errorf("%w: prompt required for askAgent action", ErrMissingHookField)
	}
	if hook.Then.Type == "runCommand" && hook.Then.Command == "" {
		return fmt.Errorf("%w: command required for runCommand action", ErrMissingHookField)
	}

	return nil
}

// isFileBasedTrigger returns true if the trigger type requires file patterns
func isFileBasedTrigger(triggerType string) bool {
	return triggerType == "fileEdited" || triggerType == "fileCreated" || triggerType == "fileDeleted"
}

// ValidateGeneratedFiles validates all generated files
func ValidateGeneratedFiles(files []GeneratedFile) error {
	for _, f := range files {
		switch f.Type {
		case "steering":
			if err := ValidateSteeringFile(f.Content); err != nil {
				return fmt.Errorf("invalid steering file %s: %w", f.Path, err)
			}
		case "hook":
			if err := ValidateHookFile(f.Content); err != nil {
				return fmt.Errorf("invalid hook file %s: %w", f.Path, err)
			}
		}
	}
	return nil
}

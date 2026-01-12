package generation

import (
	"testing"
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

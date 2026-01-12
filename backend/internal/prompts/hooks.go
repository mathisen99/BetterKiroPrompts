package prompts

// Hook preset constants
const (
	HookPresetLight   = "light"
	HookPresetBasic   = "basic"
	HookPresetDefault = "default"
	HookPresetStrict  = "strict"
)

// HookSchemaSpec contains the complete Kiro hook file schema specification.
const HookSchemaSpec = `# Kiro Hook File Schema Specification

## Overview
Hook files automate agent actions based on events. They are JSON files in '.kiro/hooks/' with the '.kiro.hook' extension.

## File Location
All hook files must be placed in: .kiro/hooks/{name}.kiro.hook

## Complete Schema
` + "```json" + `
{
  "name": "string (REQUIRED) - Human-readable hook name",
  "description": "string (REQUIRED) - What this hook does",
  "version": "string (REQUIRED) - Semantic version, e.g., '1.0.0'",
  "enabled": "boolean (REQUIRED) - Whether hook is active",
  "when": {
    "type": "string (REQUIRED) - Trigger type (see valid values below)",
    "patterns": ["string[] (REQUIRED for file events) - Glob patterns"]
  },
  "then": {
    "type": "string (REQUIRED) - Action type (see valid values and restrictions)",
    "prompt": "string (for askAgent) - Instructions for the agent",
    "command": "string (for runCommand) - Shell command to execute"
  }
}
` + "```" + `

## Valid when.type Values
| Type | Description | Requires patterns? |
|------|-------------|-------------------|
| fileEdited | Triggered when a file is edited | YES |
| fileCreated | Triggered when a file is created | YES |
| fileDeleted | Triggered when a file is deleted | YES |
| promptSubmit | Triggered when user submits a prompt | NO |
| agentStop | Triggered when agent completes work | NO |
| userTriggered | Manually triggered by user | NO |

## Valid then.type Values and Restrictions
| Type | Description | Valid with when.type |
|------|-------------|---------------------|
| askAgent | Send instructions to the agent | ALL trigger types |
| runCommand | Execute a shell command | ONLY promptSubmit, agentStop |

### CRITICAL RESTRICTION
If then.type is "runCommand", then when.type MUST be either "promptSubmit" or "agentStop".
This prevents arbitrary command execution on file changes.

## Pattern Examples
- "**/*.go" - All Go files
- "**/*.{ts,tsx}" - All TypeScript files
- "src/**/*" - All files in src directory
- "*.md" - Markdown files in root only
- "**/*.test.{ts,js}" - All test files
`

// HookExamples contains example hooks for each preset level.
const HookExamples = `## Hook Examples by Preset

### Light Preset Hooks
Minimum friction - just formatters on agent stop.

#### format-on-stop.kiro.hook
` + "```json" + `
{
  "name": "Format on Agent Stop",
  "description": "Run code formatters when agent completes work",
  "version": "1.0.0",
  "enabled": true,
  "when": {
    "type": "agentStop"
  },
  "then": {
    "type": "runCommand",
    "command": "go fmt ./... && pnpm --prefix frontend format"
  }
}
` + "```" + `

### Basic Preset Hooks
Daily discipline - formatters, linters, manual test runner.

#### lint-on-stop.kiro.hook
` + "```json" + `
{
  "name": "Lint on Agent Stop",
  "description": "Run linters when agent completes work",
  "version": "1.0.0",
  "enabled": true,
  "when": {
    "type": "agentStop"
  },
  "then": {
    "type": "runCommand",
    "command": "golangci-lint run ./... && pnpm --prefix frontend lint"
  }
}
` + "```" + `

#### test-manual.kiro.hook
` + "```json" + `
{
  "name": "Run Tests",
  "description": "Manually trigger test suite",
  "version": "1.0.0",
  "enabled": true,
  "when": {
    "type": "userTriggered"
  },
  "then": {
    "type": "askAgent",
    "prompt": "Run the test suite and summarize results. Report any failures with file locations and suggested fixes."
  }
}
` + "```" + `

### Default Preset Hooks
Balanced safety - adds secret scanning and prompt guardrails.

#### secret-scan.kiro.hook
` + "```json" + `
{
  "name": "Secret Scanner",
  "description": "Scan for accidentally committed secrets",
  "version": "1.0.0",
  "enabled": true,
  "when": {
    "type": "agentStop"
  },
  "then": {
    "type": "askAgent",
    "prompt": "Scan all modified files for potential secrets (API keys, passwords, tokens, private keys). Look for patterns like 'sk-', 'api_key', 'password=', 'secret', base64-encoded strings that look like credentials. Report any findings with file locations."
  }
}
` + "```" + `

#### prompt-guardrails.kiro.hook
` + "```json" + `
{
  "name": "Prompt Guardrails",
  "description": "Validate prompts before execution",
  "version": "1.0.0",
  "enabled": true,
  "when": {
    "type": "promptSubmit"
  },
  "then": {
    "type": "askAgent",
    "prompt": "Before proceeding with this request, verify: 1) This request doesn't skip required planning or spec steps, 2) Security implications have been considered, 3) The scope is reasonable for a single change. If any concerns, ask clarifying questions before proceeding."
  }
}
` + "```" + `

### Strict Preset Hooks
Maximum enforcement - adds static analysis and dependency scanning.

#### static-analysis.kiro.hook
` + "```json" + `
{
  "name": "Static Analysis",
  "description": "Run static analysis on code changes",
  "version": "1.0.0",
  "enabled": true,
  "when": {
    "type": "agentStop"
  },
  "then": {
    "type": "askAgent",
    "prompt": "Perform static analysis on modified files: 1) Check for common security vulnerabilities (SQL injection, XSS, path traversal), 2) Identify potential null pointer issues, 3) Look for resource leaks, 4) Check error handling completeness. Report findings with severity levels."
  }
}
` + "```" + `

#### dep-scan.kiro.hook
` + "```json" + `
{
  "name": "Dependency Scanner",
  "description": "Check for vulnerable dependencies",
  "version": "1.0.0",
  "enabled": true,
  "when": {
    "type": "agentStop"
  },
  "then": {
    "type": "askAgent",
    "prompt": "If any dependency files were modified (go.mod, package.json, requirements.txt, etc.), check for: 1) Known vulnerable versions, 2) Outdated major versions, 3) Unnecessary dependencies that could be removed. Suggest updates if needed."
  }
}
` + "```" + `

### File-Based Hook Examples

#### go-test-on-change.kiro.hook
` + "```json" + `
{
  "name": "Go Test on Change",
  "description": "Run tests when Go files change",
  "version": "1.0.0",
  "enabled": true,
  "when": {
    "type": "fileEdited",
    "patterns": ["**/*.go"]
  },
  "then": {
    "type": "askAgent",
    "prompt": "A Go file was modified. Run the relevant tests for this file and report results."
  }
}
` + "```" + `

#### typescript-typecheck.kiro.hook
` + "```json" + `
{
  "name": "TypeScript Type Check",
  "description": "Type check TypeScript files on change",
  "version": "1.0.0",
  "enabled": true,
  "when": {
    "type": "fileEdited",
    "patterns": ["**/*.{ts,tsx}"]
  },
  "then": {
    "type": "askAgent",
    "prompt": "A TypeScript file was modified. Run type checking and report any type errors."
  }
}
` + "```" + `
`

// HookPresetDescriptions describes what each preset includes.
var HookPresetDescriptions = map[string]struct {
	Title       string
	Description string
	Hooks       []string
}{
	HookPresetLight: {
		Title:       "Light",
		Description: "Minimum friction - just formatters on agent stop",
		Hooks:       []string{"format-on-stop"},
	},
	HookPresetBasic: {
		Title:       "Basic",
		Description: "Daily discipline - formatters, linters, manual test runner",
		Hooks:       []string{"format-on-stop", "lint-on-stop", "test-manual"},
	},
	HookPresetDefault: {
		Title:       "Default (Recommended)",
		Description: "Balanced safety - adds secret scanning and prompt guardrails",
		Hooks:       []string{"format-on-stop", "lint-on-stop", "test-manual", "secret-scan", "prompt-guardrails"},
	},
	HookPresetStrict: {
		Title:       "Strict",
		Description: "Maximum enforcement - adds static analysis and dependency scanning",
		Hooks:       []string{"format-on-stop", "lint-on-stop", "test-manual", "secret-scan", "prompt-guardrails", "static-analysis", "dep-scan"},
	},
}

// HooksSystemPrompt returns the complete system prompt for hook file generation.
func HooksSystemPrompt() string {
	return HookSchemaSpec + "\n\n" + HookExamples
}

package prompts

// SteeringFormatSpec contains the complete Kiro steering file format specification
// for inclusion in AI system prompts.
const SteeringFormatSpec = `# Kiro Steering File Format Specification

## Overview
Steering files are markdown files in '.kiro/steering/' that provide persistent context to Kiro.
They guide the AI assistant with project-specific rules, conventions, and constraints.

## File Structure
Every steering file is a markdown file with YAML frontmatter at the top.

### Frontmatter Format (Required)
` + "```yaml" + `
---
inclusion: always | fileMatch | manual
fileMatchPattern: "glob-pattern"  # Required ONLY for fileMatch mode
---
` + "```" + `

## Inclusion Modes

### 1. always (Default)
- Loaded into EVERY Kiro interaction automatically
- Use for: project-wide standards, tech stack decisions, core conventions
- Typical files: product.md, tech.md, structure.md

### 2. fileMatch
- Loaded ONLY when working with files matching the glob pattern
- REQUIRES: fileMatchPattern field with valid glob
- Use for: language-specific rules, domain-specific standards
- Examples:
  - fileMatchPattern: "**/*.go" → Loaded for Go files
  - fileMatchPattern: "**/*.{ts,tsx}" → Loaded for TypeScript/React files
  - fileMatchPattern: "**/*.py" → Loaded for Python files
  - fileMatchPattern: "src/api/**/*" → Loaded for API directory files
  - fileMatchPattern: "**/*.{js,jsx,ts,tsx}" → Loaded for all JavaScript/TypeScript

### 3. manual
- Loaded on-demand via #steering-file-name in chat
- Use for: specialized workflows, troubleshooting guides, rarely-needed context
- User must explicitly reference the file to include it

## Content Best Practices
1. Keep files focused: ONE domain per file
2. Be concise: short, strict, practical guidance
3. Provide examples: code snippets, before/after comparisons
4. Explain why: rationale for decisions, not just rules
5. NEVER include: API keys, passwords, secrets, credentials
6. Use file references: link to live project files when helpful
7. Keep under 500 lines: shorter is better for AI context

## Required Steering Files

Every project should have these three core files with 'inclusion: always':
1. product.md - What we're building and NOT building
2. tech.md - Technology choices and architecture rules
3. structure.md - Repository layout and conventions
`

// SteeringTemplates contains templates for each steering file type.
const SteeringTemplates = `## Steering File Templates

### product.md Template
` + "```markdown" + `
---
inclusion: always
---

# Product

## What We Are Building
{One to two sentence description of the product's core purpose}

## What We Are NOT Building
{Explicit list of out-of-scope features and anti-goals}
- Not a {thing}
- Not a replacement for {thing}
- No {feature} in this version

## Definition of Done
{Measurable completion criteria}
1. User can {action}
2. {Feature} works end-to-end
3. Tests pass with {coverage}% coverage

## Absolute Rules
{Non-negotiable constraints - things that must NEVER happen}
1. Never {action}
2. Always {action} before {action}
3. Do not {action} without {condition}
` + "```" + `

### tech.md Template
` + "```markdown" + `
---
inclusion: always
---

# Tech Stack

## Pinned Versions
- {Language}: {version}
- {Framework}: {version}
- {Database}: {version}

## Architecture Rules
- {Pattern}: {description}
- {Constraint}: {reason}
- {Boundary}: {what goes where}

## Simplicity Rules
- No {anti-pattern}
- Prefer {approach} over {alternative}
- Standard library over dependencies when possible

## Build & Run
{Commands to build, test, and run the project}
` + "```bash" + `
# Development
{dev command}

# Testing
{test command}

# Production
{prod command}
` + "```" + `
` + "```" + `

### structure.md Template
` + "```markdown" + `
---
inclusion: always
---

# Repository Structure

## Layout
` + "```" + `
/
├── {folder}/
│   ├── {subfolder}/
│   └── {file}
├── {folder}/
│   └── {file}
└── {config files}
` + "```" + `

## Conventions
- {Naming convention}: {pattern}
- {Organization rule}: {description}
- {File placement}: {where things go}

## Rules
- No {anti-pattern folders}
- One {thing} per {scope}
- Keep {structure} until {condition}
` + "```" + `

### security-{lang}.md Template (fileMatch)
` + "```markdown" + `
---
inclusion: fileMatch
fileMatchPattern: "**/*.{ext}"
---

# Security Guidelines for {Language}

## No Secrets
- NEVER commit credentials, API keys, or tokens
- Use environment variables for all secrets
- Add secret patterns to .gitignore

## Input Validation
- Validate ALL user input before processing
- Sanitize data before database queries
- Use parameterized queries, never string concatenation

## Auth Boundaries
- Check permissions before every sensitive operation
- Never trust client-side validation alone
- Log security-relevant events

## Least Privilege
- Request minimum necessary permissions
- Scope database access appropriately
- Use read-only connections where possible
` + "```" + `

### quality-{lang}.md Template (fileMatch)
` + "```markdown" + `
---
inclusion: fileMatch
fileMatchPattern: "**/*.{ext}"
---

# Code Quality for {Language}

## Formatting
- Use {formatter} for consistent style
- Run formatter before every commit
- {Specific formatting rules}

## Linting
- Use {linter} with {config}
- Zero warnings policy
- {Specific lint rules}

## Testing
- Minimum {X}% coverage for new code
- Unit tests for all business logic
- Integration tests for API endpoints

## Documentation
- Document public APIs
- Explain non-obvious code with comments
- Keep README up to date
` + "```" + `
`

// LanguagePatterns maps languages to their file match patterns.
var LanguagePatterns = map[string]string{
	"go":         "**/*.go",
	"typescript": "**/*.{ts,tsx}",
	"javascript": "**/*.{js,jsx}",
	"python":     "**/*.py",
	"rust":       "**/*.rs",
	"java":       "**/*.java",
	"csharp":     "**/*.cs",
	"ruby":       "**/*.rb",
	"php":        "**/*.php",
	"web":        "**/*.{ts,tsx,js,jsx,html,css}",
}

// SteeringSystemPrompt returns the complete system prompt for steering file generation.
func SteeringSystemPrompt() string {
	return SteeringFormatSpec + "\n\n" + SteeringTemplates
}

# Design Document: Phase 4 - Final Overhaul

## Overview

This design delivers a production-ready Kiro Prompting & Guardrails Generator with three major improvements:

1. **Professional UI** - Dark theme with blue accents matching shadcn's blue theme
2. **Smarter AI Generation** - Experience-level adaptation, comprehensive file format knowledge
3. **Complete File Coverage** - All steering files, valid hooks, AGENTS.md

The architecture maintains the existing Go backend + React frontend pattern while significantly enhancing the AI prompts and UI components.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (React)                         │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │ LandingPage  │  │ QuestionFlow │  │    OutputEditor      │  │
│  │              │  │              │  │                      │  │
│  │ - Hero       │  │ - Progress   │  │ - Tabs (Kickoff/     │  │
│  │ - Level      │  │ - Q&A Cards  │  │   Steering/Hooks)    │  │
│  │   Selector   │  │ - Navigation │  │ - File Editor        │  │
│  │ - Project    │  │              │  │ - Download Actions   │  │
│  │   Input      │  │              │  │                      │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Backend (Go)                                │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                   Generation Service                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐   │   │
│  │  │  Question   │  │   Output    │  │    Validator    │   │   │
│  │  │  Generator  │  │  Generator  │  │                 │   │   │
│  │  └─────────────┘  └─────────────┘  └─────────────────┘   │   │
│  │         │                │                  │            │   │
│  │         ▼                ▼                  ▼            │   │
│  │  ┌─────────────────────────────────────────────────┐     │   │
│  │  │           AI Knowledge Base (Prompts)           │     │   │
│  │  │  - Steering Format Specs                        │     │   │
│  │  │  - Hook Schema & Examples                       │     │   │
│  │  │  - Kickoff Prompt Template                      │     │   │
│  │  │  - Experience Level Adaptations                 │     │   │
│  │  └─────────────────────────────────────────────────┘     │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Frontend Components

#### 1. ExperienceLevelSelector
```typescript
interface ExperienceLevelSelectorProps {
  onSelect: (level: ExperienceLevel) => void
  selected?: ExperienceLevel
}

type ExperienceLevel = 'beginner' | 'novice' | 'expert'

interface LevelOption {
  id: ExperienceLevel
  title: string
  description: string
  icon: ReactNode
}
```

#### 2. Enhanced LandingPage
```typescript
interface LandingPageState {
  phase: 'level-select' | 'input' | 'questions' | 'generating' | 'output' | 'error'
  experienceLevel: ExperienceLevel | null
  projectIdea: string
  hookPreset: HookPreset
  questions: Question[]
  answers: Map<number, string>
  currentQuestionIndex: number
  generatedFiles: GeneratedFile[]
  editedFiles: Map<string, string>
  error: string | null
  retryAfter: number | null
}

type HookPreset = 'light' | 'basic' | 'default' | 'strict'
```

#### 3. HookPresetSelector
```typescript
interface HookPresetSelectorProps {
  onSelect: (preset: HookPreset) => void
  selected: HookPreset
}

interface PresetOption {
  id: HookPreset
  title: string
  description: string
  hooks: string[] // List of hooks included
}
```

### Backend Interfaces

#### 1. Enhanced API Types
```go
type GenerateQuestionsRequest struct {
    ProjectIdea     string `json:"projectIdea"`
    ExperienceLevel string `json:"experienceLevel"` // "beginner", "novice", "expert"
}

type GenerateOutputsRequest struct {
    ProjectIdea     string   `json:"projectIdea"`
    Answers         []Answer `json:"answers"`
    ExperienceLevel string   `json:"experienceLevel"`
    HookPreset      string   `json:"hookPreset"` // "light", "basic", "default", "strict"
}
```

#### 2. File Validator
```go
type FileValidator interface {
    ValidateSteeringFile(content string) error
    ValidateHookFile(content string) error
    ValidateKickoffPrompt(content string) error
}
```

## Data Models

### Experience Level Definitions

```typescript
const EXPERIENCE_LEVELS = {
  beginner: {
    title: 'Beginner',
    description: 'New to programming. I need guidance on basics.',
    questionStyle: 'simple',
    avoidTerms: ['microservices', 'distributed', 'scalability', 'concurrency', 'middleware'],
    suggestSimpler: true
  },
  novice: {
    title: 'Novice', 
    description: 'Some experience. I understand basic concepts.',
    questionStyle: 'moderate',
    includeHints: true,
    suggestSimpler: false
  },
  expert: {
    title: 'Expert',
    description: 'Experienced developer. Give me the technical details.',
    questionStyle: 'technical',
    includeArchitecture: true,
    suggestSimpler: false
  }
}
```

### Hook Preset Definitions

```typescript
const HOOK_PRESETS = {
  light: {
    title: 'Light',
    description: 'Minimum friction - just formatters on agent stop',
    hooks: ['format-on-stop']
  },
  basic: {
    title: 'Basic',
    description: 'Daily discipline - formatters, linters, manual test runner',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual']
  },
  default: {
    title: 'Default (Recommended)',
    description: 'Balanced safety - adds secret scanning and prompt guardrails',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual', 'secret-scan', 'prompt-guardrails']
  },
  strict: {
    title: 'Strict',
    description: 'Maximum enforcement - adds static analysis and dependency scanning',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual', 'secret-scan', 'prompt-guardrails', 'static-analysis', 'dep-scan']
  }
}
```

### Generated File Types

```typescript
interface GeneratedFile {
  path: string
  content: string
  type: 'kickoff' | 'steering' | 'hook' | 'agents'
}

// Steering file structure
interface SteeringFile {
  frontmatter: {
    inclusion: 'always' | 'fileMatch' | 'manual'
    fileMatchPattern?: string
  }
  content: string
}

// Hook file structure (Kiro schema)
interface KiroHook {
  name: string
  description: string
  version: string
  enabled: boolean
  when: {
    type: 'fileEdited' | 'fileCreated' | 'fileDeleted' | 'promptSubmit' | 'agentStop' | 'userTriggered'
    patterns?: string[]
  }
  then: {
    type: 'askAgent' | 'runCommand'
    prompt?: string
    command?: string
  }
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Experience Level Adaptation

*For any* generated question set, when the experience level is "beginner", the questions SHALL NOT contain technical jargon terms (microservices, distributed, scalability, concurrency, middleware) AND SHALL include more explanatory hints.

**Validates: Requirements 1.2, 1.3, 1.4, 3.1, 6.6**

### Property 2: Core Steering Files Validity

*For any* generated output, the files SHALL include at minimum product.md, tech.md, and structure.md, each with valid frontmatter containing `inclusion: always`, and valid markdown content.

**Validates: Requirements 4.1, 4.2, 4.3, 4.8, 10.1**

### Property 3: Conditional Steering Files Pattern

*For any* generated output where the project uses a specific language (Go, TypeScript, etc.), the security and quality steering files SHALL have `inclusion: fileMatch` and a `fileMatchPattern` that matches files of that language.

**Validates: Requirements 4.4, 4.5, 4.6, 4.7, 10.2, 10.3**

### Property 4: Hook File Schema Validity

*For any* generated hook file, the JSON SHALL contain all required fields (name, description, version, enabled, when, then), the `when.type` SHALL be one of the valid values, and if `then.type` is "runCommand" then `when.type` SHALL be either "promptSubmit" or "agentStop".

**Validates: Requirements 5.1, 5.2, 5.3, 5.4, 5.5, 10.4**

### Property 5: Kickoff Prompt Completeness

*For any* generated kickoff prompt, the content SHALL contain the phrase "no coding" (or equivalent enforcement), AND SHALL include sections for: Project Identity, Success Criteria, Users & Roles, Data Sensitivity, Auth Model, Concurrency, Boundaries, Non-Goals, Constraints, Risks & Tradeoffs, and Boundary Examples.

**Validates: Requirements 6.1, 6.2, 6.4, 6.5**

### Property 6: Question Generation Constraints

*For any* generated question set, the count SHALL be between 5 and 10 inclusive, AND the questions SHALL follow a logical ordering where identity/scope questions appear before technical/architecture questions.

**Validates: Requirements 3.4, 3.6**

## AI Knowledge Base (System Prompts)

This is the critical component that enables the AI to generate high-quality, valid Kiro files. The system prompts must contain comprehensive knowledge about file formats, schemas, and best practices.

### Steering File Format Specification

```markdown
# Kiro Steering File Format

## Structure
Every steering file is a markdown file with YAML frontmatter.

## Frontmatter (Required)
---
inclusion: always | fileMatch | manual
fileMatchPattern: "glob-pattern"  # Required only for fileMatch
---

## Inclusion Modes

### always (Default)
- Loaded into every Kiro interaction
- Use for: project-wide standards, tech stack, core conventions
- Files: product.md, tech.md, structure.md

### fileMatch
- Loaded only when working with files matching the pattern
- Use for: language-specific rules, domain-specific standards
- Requires: fileMatchPattern with valid glob
- Examples:
  - fileMatchPattern: "**/*.go" → Go files
  - fileMatchPattern: "**/*.{ts,tsx}" → TypeScript/React files
  - fileMatchPattern: "src/api/**/*" → API directory

### manual
- Loaded on-demand via #steering-file-name in chat
- Use for: specialized workflows, troubleshooting guides

## Content Best Practices
- Keep files focused: one domain per file
- Be concise: short, strict, practical guidance
- Provide examples: code snippets, before/after comparisons
- Explain why: rationale for decisions, not just rules
- Never include: API keys, passwords, secrets
- Use file references: link to live project files when helpful
```

### Required Steering Files

```markdown
# product.md (inclusion: always)
Required sections:
- What We Are Building: 1-2 sentence product description
- What We Are NOT Building: explicit non-goals
- Definition of Done: measurable completion criteria
- Absolute Rules: non-negotiable constraints

# tech.md (inclusion: always)
Required sections:
- Stack Choices: languages, frameworks, databases
- Architecture Rules: patterns, constraints, boundaries
- Simplicity Rules: what to avoid, preferences
- Build & Run: commands to build/test/run

# structure.md (inclusion: always)
Required sections:
- Repository Layout: folder structure diagram
- Conventions: naming, organization patterns
- Rules: what goes where, anti-patterns to avoid

# security-{lang}.md (inclusion: fileMatch)
Required sections:
- No Secrets: never commit credentials
- Input Validation: sanitization requirements
- Auth Boundaries: explicit access control
- Least Privilege: minimal permissions

# quality-{lang}.md (inclusion: fileMatch)
Required sections:
- Formatting: tools and standards
- Linting: rules and exceptions
- Testing: coverage expectations
- Documentation: what to document
```

### Hook File Schema

```json
{
  "name": "string (required) - Human-readable hook name",
  "description": "string (required) - What this hook does",
  "version": "string (required) - Semantic version, e.g., '1.0.0'",
  "enabled": "boolean (required) - Whether hook is active",
  "when": {
    "type": "string (required) - One of: fileEdited, fileCreated, fileDeleted, promptSubmit, agentStop, userTriggered",
    "patterns": ["string[] (required for file events) - Glob patterns, e.g., '**/*.go'"]
  },
  "then": {
    "type": "string (required) - 'askAgent' (all triggers) or 'runCommand' (only promptSubmit/agentStop)",
    "prompt": "string (for askAgent) - Instructions for the agent",
    "command": "string (for runCommand) - Shell command to execute"
  }
}
```

### Hook Examples by Preset

```json
// Light Preset - format-on-stop.kiro.hook
{
  "name": "Format on Agent Stop",
  "description": "Run formatters when agent completes work",
  "version": "1.0.0",
  "enabled": true,
  "when": { "type": "agentStop" },
  "then": {
    "type": "runCommand",
    "command": "go fmt ./... && pnpm --prefix frontend format"
  }
}

// Basic Preset - test-manual.kiro.hook
{
  "name": "Run Tests",
  "description": "Manually trigger test suite",
  "version": "1.0.0",
  "enabled": true,
  "when": { "type": "userTriggered" },
  "then": {
    "type": "askAgent",
    "prompt": "Run the test suite and summarize results. Report any failures with file locations."
  }
}

// Default Preset - secret-scan.kiro.hook
{
  "name": "Secret Scanner",
  "description": "Scan for accidentally committed secrets",
  "version": "1.0.0",
  "enabled": true,
  "when": { "type": "agentStop" },
  "then": {
    "type": "askAgent",
    "prompt": "Scan modified files for potential secrets (API keys, passwords, tokens). Report any findings."
  }
}

// Default Preset - prompt-guardrails.kiro.hook
{
  "name": "Prompt Guardrails",
  "description": "Validate prompts before execution",
  "version": "1.0.0",
  "enabled": true,
  "when": { "type": "promptSubmit" },
  "then": {
    "type": "askAgent",
    "prompt": "Before proceeding, verify: 1) This request doesn't skip required planning steps, 2) Security implications are considered, 3) The scope is reasonable for one change."
  }
}
```

### Kickoff Prompt Template

```markdown
# Project Kickoff: {Project Name}

> ⚠️ **IMPORTANT**: Do not write any code until all questions below are answered and reviewed.

## Project Identity
{One sentence description of what this project is}

## Success Criteria
{What does "done" look like? Measurable outcomes}

## Users & Roles
{Who uses this? Anonymous, authenticated, admin, etc.}

## Data Sensitivity
{What data is stored? Label sensitive data explicitly}

### Data Lifecycle
- Retention: {how long is data kept}
- Deletion: {how users delete their data}
- Export: {can users export their data}
- Audit: {what is logged}
- Backups: {backup strategy}

## Auth Model
{none / basic / external provider / custom}

## Concurrency Expectations
{Multi-user? Background jobs? Shared state?}

## Risks & Tradeoffs
1. **Risk**: {description}
   - Mitigation: {how to address}
   - Not Handled: {what we're accepting}
2. **Risk**: {description}
   - Mitigation: {how to address}
   - Not Handled: {what we're accepting}
3. **Risk**: {description}
   - Mitigation: {how to address}
   - Not Handled: {what we're accepting}

## Boundaries
{Public vs private data boundaries}

### Boundary Examples
- {Role} CAN {action} on {resource}
- {Role} CANNOT {action} on {resource}
- {Role} CAN {action} on {resource} IF {condition}

## Non-Goals
{What will NOT be built}

## Constraints
{Time, simplicity, tech limits}

---

## Next Steps
1. Review this document with stakeholders
2. Create specs for each major feature
3. Begin implementation only after specs are approved
```

### AGENTS.md Template

```markdown
# Agent Guidelines

## Core Principles
1. Always follow steering files in `.kiro/steering/`
2. Never invent requirements - ask if unclear
3. Prefer small, reviewable changes
4. Update docs when behavior changes

## Before Coding
- Ensure requirements are clear
- Check for existing patterns in codebase
- Consider security implications
- Plan for testing

## Commit Standards
- Atomic commits (one concern per commit)
- Prefix: feat: | fix: | docs: | chore:
- One-sentence summary
- No mixed or vague commits

## When Stuck
- Ask for clarification rather than guessing
- Reference relevant steering files
- Suggest alternatives if blocked
```

## Error Handling

### Frontend Error States

1. **Rate Limit Exceeded (429)**
   - Display countdown timer with `retryAfter` seconds
   - Disable submit buttons until timer expires
   - Show friendly message: "Too many requests. Please wait X seconds."

2. **Generation Failed (500)**
   - Display error message with retry option
   - Offer to start over from beginning
   - Log error details for debugging

3. **Timeout (504)**
   - Display timeout message
   - Suggest refreshing and trying again
   - Consider breaking into smaller requests

4. **Invalid Response**
   - Backend retries once automatically
   - If still invalid, show generic error
   - Log malformed response for analysis

### Backend Validation

1. **Input Validation**
   - Project idea: max 2000 chars, non-empty
   - Answers: max 1000 chars each
   - Experience level: must be valid enum
   - Hook preset: must be valid enum

2. **Output Validation**
   - Steering files: valid frontmatter YAML
   - Hook files: valid JSON matching schema
   - Kickoff prompt: contains required sections

## Testing Strategy

### Unit Tests

- Component rendering tests for new UI components
- State management tests for experience level persistence
- Validation function tests for file format checking

### Property-Based Tests

Using a property-based testing library (e.g., fast-check for TypeScript tests):

1. **Experience Level Adaptation Test**
   - Generate random project ideas
   - For each experience level, verify question content matches level expectations
   - Minimum 100 iterations

2. **Steering File Validity Test**
   - Generate outputs for random project configurations
   - Verify all steering files have valid frontmatter
   - Verify conditional files have correct patterns
   - Minimum 100 iterations

3. **Hook Schema Validity Test**
   - Generate hooks for all presets
   - Parse JSON and verify schema compliance
   - Verify runCommand restrictions
   - Minimum 100 iterations

4. **Kickoff Prompt Completeness Test**
   - Generate kickoff prompts for various projects
   - Verify all required sections present
   - Verify "no coding" enforcement phrase
   - Minimum 100 iterations

5. **Question Constraints Test**
   - Generate questions for random projects
   - Verify count is 5-10
   - Verify ordering follows logical sequence
   - Minimum 100 iterations

### Integration Tests

- End-to-end flow from level selection to file download
- API endpoint tests with various inputs
- Error handling scenarios
# Design Document: Phase 3 - Final Polish

## Overview

This design delivers a production-ready Kiro Prompting & Guardrails Generator with three major improvements:

1. **Professional UI** - Dark theme with blue accents matching shadcn's blue theme
2. **Smarter AI Generation** - Experience-level adaptation, comprehensive file format knowledge
3. **Complete File Coverage** - All steering files, valid hooks, AGENTS.md

The architecture maintains the existing Go backend + React frontend pattern while significantly enhancing the AI prompts and UI components.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (React)                         │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │ LandingPage  │  │ QuestionFlow │  │    OutputEditor      │  │
│  │              │  │              │  │                      │  │
│  │ - Hero       │  │ - Progress   │  │ - Tabs (Kickoff/     │  │
│  │ - Level      │  │ - Q&A Cards  │  │   Steering/Hooks)    │  │
│  │   Selector   │  │ - Navigation │  │ - File Editor        │  │
│  │ - Project    │  │              │  │ - Download Actions   │  │
│  │   Input      │  │              │  │                      │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Backend (Go)                                │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                   Generation Service                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐   │   │
│  │  │  Question   │  │   Output    │  │    Validator    │   │   │
│  │  │  Generator  │  │  Generator  │  │                 │   │   │
│  │  └─────────────┘  └─────────────┘  └─────────────────┘   │   │
│  │         │                │                  │            │   │
│  │         ▼                ▼                  ▼            │   │
│  │  ┌─────────────────────────────────────────────────┐     │   │
│  │  │           AI Knowledge Base (Prompts)           │     │   │
│  │  │  - Steering Format Specs                        │     │   │
│  │  │  - Hook Schema & Examples                       │     │   │
│  │  │  - Kickoff Prompt Template                      │     │   │
│  │  │  - Experience Level Adaptations                 │     │   │
│  │  └─────────────────────────────────────────────────┘     │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Frontend Components

#### 1. ExperienceLevelSelector
```typescript
interface ExperienceLevelSelectorProps {
  onSelect: (level: ExperienceLevel) => void
  selected?: ExperienceLevel
}

type ExperienceLevel = 'beginner' | 'novice' | 'expert'

interface LevelOption {
  id: ExperienceLevel
  title: string
  description: string
  icon: ReactNode
}
```

#### 2. Enhanced LandingPage
```typescript
interface LandingPageState {
  phase: 'level-select' | 'input' | 'questions' | 'generating' | 'output' | 'error'
  experienceLevel: ExperienceLevel | null
  projectIdea: string
  hookPreset: HookPreset
  questions: Question[]
  answers: Map<number, string>
  currentQuestionIndex: number
  generatedFiles: GeneratedFile[]
  editedFiles: Map<string, string>
  error: string | null
  retryAfter: number | null
}

type HookPreset = 'light' | 'basic' | 'default' | 'strict'
```

#### 3. HookPresetSelector
```typescript
interface HookPresetSelectorProps {
  onSelect: (preset: HookPreset) => void
  selected: HookPreset
}

interface PresetOption {
  id: HookPreset
  title: string
  description: string
  hooks: string[] // List of hooks included
}
```

### Backend Interfaces

#### 1. Enhanced API Types
```go
type GenerateQuestionsRequest struct {
    ProjectIdea     string `json:"projectIdea"`
    ExperienceLevel string `json:"experienceLevel"` // "beginner", "novice", "expert"
}

type GenerateOutputsRequest struct {
    ProjectIdea     string   `json:"projectIdea"`
    Answers         []Answer `json:"answers"`
    ExperienceLevel string   `json:"experienceLevel"`
    HookPreset      string   `json:"hookPreset"` // "light", "basic", "default", "strict"
}
```

#### 2. File Validator
```go
type FileValidator interface {
    ValidateSteeringFile(content string) error
    ValidateHookFile(content string) error
    ValidateKickoffPrompt(content string) error
}
```

## Data Models

### Experience Level Definitions

```typescript
const EXPERIENCE_LEVELS = {
  beginner: {
    title: 'Beginner',
    description: 'New to programming. I need guidance on basics.',
    questionStyle: 'simple',
    avoidTerms: ['microservices', 'distributed', 'scalability', 'concurrency', 'middleware'],
    suggestSimpler: true
  },
  novice: {
    title: 'Novice', 
    description: 'Some experience. I understand basic concepts.',
    questionStyle: 'moderate',
    includeHints: true,
    suggestSimpler: false
  },
  expert: {
    title: 'Expert',
    description: 'Experienced developer. Give me the technical details.',
    questionStyle: 'technical',
    includeArchitecture: true,
    suggestSimpler: false
  }
}
```

### Hook Preset Definitions

```typescript
const HOOK_PRESETS = {
  light: {
    title: 'Light',
    description: 'Minimum friction - just formatters on agent stop',
    hooks: ['format-on-stop']
  },
  basic: {
    title: 'Basic',
    description: 'Daily discipline - formatters, linters, manual test runner',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual']
  },
  default: {
    title: 'Default (Recommended)',
    description: 'Balanced safety - adds secret scanning and prompt guardrails',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual', 'secret-scan', 'prompt-guardrails']
  },
  strict: {
    title: 'Strict',
    description: 'Maximum enforcement - adds static analysis and dependency scanning',
    hooks: ['format-on-stop', 'lint-on-stop', 'test-manual', 'secret-scan', 'prompt-guardrails', 'static-analysis', 'dep-scan']
  }
}
```

### Generated File Types

```typescript
interface GeneratedFile {
  path: string
  content: string
  type: 'kickoff' | 'steering' | 'hook' | 'agents'
}

// Steering file structure
interface SteeringFile {
  frontmatter: {
    inclusion: 'always' | 'fileMatch' | 'manual'
    fileMatchPattern?: string
  }
  content: string
}

// Hook file structure (Kiro schema)
interface KiroHook {
  name: string
  description: string
  version: string
  enabled: boolean
  when: {
    type: 'fileEdited' | 'fileCreated' | 'fileDeleted' | 'promptSubmit' | 'agentStop' | 'userTriggered'
    patterns?: string[]
  }
  then: {
    type: 'askAgent' | 'runCommand'
    prompt?: string
    command?: string
  }
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Experience Level Adaptation

*For any* generated question set, when the experience level is "beginner", the questions SHALL NOT contain technical jargon terms (microservices, distributed, scalability, concurrency, middleware) AND SHALL include more explanatory hints.

**Validates: Requirements 1.2, 1.3, 1.4, 3.1, 6.6**

### Property 2: Core Steering Files Validity

*For any* generated output, the files SHALL include at minimum product.md, tech.md, and structure.md, each with valid frontmatter containing `inclusion: always`, and valid markdown content.

**Validates: Requirements 4.1, 4.2, 4.3, 4.8, 10.1**

### Property 3: Conditional Steering Files Pattern

*For any* generated output where the project uses a specific language (Go, TypeScript, etc.), the security and quality steering files SHALL have `inclusion: fileMatch` and a `fileMatchPattern` that matches files of that language.

**Validates: Requirements 4.4, 4.5, 4.6, 4.7, 10.2, 10.3**

### Property 4: Hook File Schema Validity

*For any* generated hook file, the JSON SHALL contain all required fields (name, description, version, enabled, when, then), the `when.type` SHALL be one of the valid values, and if `then.type` is "runCommand" then `when.type` SHALL be either "promptSubmit" or "agentStop".

**Validates: Requirements 5.1, 5.2, 5.3, 5.4, 5.5, 10.4**

### Property 5: Kickoff Prompt Completeness

*For any* generated kickoff prompt, the content SHALL contain the phrase "no coding" (or equivalent enforcement), AND SHALL include sections for: Project Identity, Success Criteria, Users & Roles, Data Sensitivity, Auth Model, Concurrency, Boundaries, Non-Goals, Constraints, Risks & Tradeoffs, and Boundary Examples.

**Validates: Requirements 6.1, 6.2, 6.4, 6.5**

### Property 6: Question Generation Constraints

*For any* generated question set, the count SHALL be between 5 and 10 inclusive, AND the questions SHALL follow a logical ordering where identity/scope questions appear before technical/architecture questions.

**Validates: Requirements 3.4, 3.6**

## AI Knowledge Base (System Prompts)

This is the critical component that enables the AI to generate high-quality, valid Kiro files. The system prompts must contain comprehensive knowledge about file formats, schemas, and best practices.

### Steering File Format Specification

```markdown
# Kiro Steering File Format

## Structure
Every steering file is a markdown file with YAML frontmatter.

## Frontmatter (Required)
---
inclusion: always | fileMatch | manual
fileMatchPattern: "glob-pattern"  # Required only for fileMatch
---

## Inclusion Modes

### always (Default)
- Loaded into every Kiro interaction
- Use for: project-wide standards, tech stack, core conventions
- Files: product.md, tech.md, structure.md

### fileMatch
- Loaded only when working with files matching the pattern
- Use for: language-specific rules, domain-specific standards
- Requires: fileMatchPattern with valid glob
- Examples:
  - fileMatchPattern: "**/*.go" → Go files
  - fileMatchPattern: "**/*.{ts,tsx}" → TypeScript/React files
  - fileMatchPattern: "src/api/**/*" → API directory

### manual
- Loaded on-demand via #steering-file-name in chat
- Use for: specialized workflows, troubleshooting guides

## Content Best Practices
- Keep files focused: one domain per file
- Be concise: short, strict, practical guidance
- Provide examples: code snippets, before/after comparisons
- Explain why: rationale for decisions, not just rules
- Never include: API keys, passwords, secrets
- Use file references: link to live project files when helpful
```

### Required Steering Files

```markdown
# product.md (inclusion: always)
Required sections:
- What We Are Building: 1-2 sentence product description
- What We Are NOT Building: explicit non-goals
- Definition of Done: measurable completion criteria
- Absolute Rules: non-negotiable constraints

# tech.md (inclusion: always)
Required sections:
- Stack Choices: languages, frameworks, databases
- Architecture Rules: patterns, constraints, boundaries
- Simplicity Rules: what to avoid, preferences
- Build & Run: commands to build/test/run

# structure.md (inclusion: always)
Required sections:
- Repository Layout: folder structure diagram
- Conventions: naming, organization patterns
- Rules: what goes where, anti-patterns to avoid

# security-{lang}.md (inclusion: fileMatch)
Required sections:
- No Secrets: never commit credentials
- Input Validation: sanitization requirements
- Auth Boundaries: explicit access control
- Least Privilege: minimal permissions

# quality-{lang}.md (inclusion: fileMatch)
Required sections:
- Formatting: tools and standards
- Linting: rules and exceptions
- Testing: coverage expectations
- Documentation: what to document
```

### Hook File Schema

```json
{
  "name": "string (required) - Human-readable hook name",
  "description": "string (required) - What this hook does",
  "version": "string (required) - Semantic version, e.g., '1.0.0'",
  "enabled": "boolean (required) - Whether hook is active",
  "when": {
    "type": "string (required) - One of: fileEdited, fileCreated, fileDeleted, promptSubmit, agentStop, userTriggered",
    "patterns": ["string[] (required for file events) - Glob patterns, e.g., '**/*.go'"]
  },
  "then": {
    "type": "string (required) - 'askAgent' (all triggers) or 'runCommand' (only promptSubmit/agentStop)",
    "prompt": "string (for askAgent) - Instructions for the agent",
    "command": "string (for runCommand) - Shell command to execute"
  }
}
```

### Hook Examples by Preset

```json
// Light Preset - format-on-stop.kiro.hook
{
  "name": "Format on Agent Stop",
  "description": "Run formatters when agent completes work",
  "version": "1.0.0",
  "enabled": true,
  "when": { "type": "agentStop" },
  "then": {
    "type": "runCommand",
    "command": "go fmt ./... && pnpm --prefix frontend format"
  }
}

// Basic Preset - test-manual.kiro.hook
{
  "name": "Run Tests",
  "description": "Manually trigger test suite",
  "version": "1.0.0",
  "enabled": true,
  "when": { "type": "userTriggered" },
  "then": {
    "type": "askAgent",
    "prompt": "Run the test suite and summarize results. Report any failures with file locations."
  }
}

// Default Preset - secret-scan.kiro.hook
{
  "name": "Secret Scanner",
  "description": "Scan for accidentally committed secrets",
  "version": "1.0.0",
  "enabled": true,
  "when": { "type": "agentStop" },
  "then": {
    "type": "askAgent",
    "prompt": "Scan modified files for potential secrets (API keys, passwords, tokens). Report any findings."
  }
}

// Default Preset - prompt-guardrails.kiro.hook
{
  "name": "Prompt Guardrails",
  "description": "Validate prompts before execution",
  "version": "1.0.0",
  "enabled": true,
  "when": { "type": "promptSubmit" },
  "then": {
    "type": "askAgent",
    "prompt": "Before proceeding, verify: 1) This request doesn't skip required planning steps, 2) Security implications are considered, 3) The scope is reasonable for one change."
  }
}
```

### Kickoff Prompt Template

```markdown
# Project Kickoff: {Project Name}

> ⚠️ **IMPORTANT**: Do not write any code until all questions below are answered and reviewed.

## Project Identity
{One sentence description of what this project is}

## Success Criteria
{What does "done" look like? Measurable outcomes}

## Users & Roles
{Who uses this? Anonymous, authenticated, admin, etc.}

## Data Sensitivity
{What data is stored? Label sensitive data explicitly}

### Data Lifecycle
- Retention: {how long is data kept}
- Deletion: {how users delete their data}
- Export: {can users export their data}
- Audit: {what is logged}
- Backups: {backup strategy}

## Auth Model
{none / basic / external provider / custom}

## Concurrency Expectations
{Multi-user? Background jobs? Shared state?}

## Risks & Tradeoffs
1. **Risk**: {description}
   - Mitigation: {how to address}
   - Not Handled: {what we're accepting}
2. **Risk**: {description}
   - Mitigation: {how to address}
   - Not Handled: {what we're accepting}
3. **Risk**: {description}
   - Mitigation: {how to address}
   - Not Handled: {what we're accepting}

## Boundaries
{Public vs private data boundaries}

### Boundary Examples
- {Role} CAN {action} on {resource}
- {Role} CANNOT {action} on {resource}
- {Role} CAN {action} on {resource} IF {condition}

## Non-Goals
{What will NOT be built}

## Constraints
{Time, simplicity, tech limits}

---

## Next Steps
1. Review this document with stakeholders
2. Create specs for each major feature
3. Begin implementation only after specs are approved
```

### AGENTS.md Template

```markdown
# Agent Guidelines

## Core Principles
1. Always follow steering files in `.kiro/steering/`
2. Never invent requirements - ask if unclear
3. Prefer small, reviewable changes
4. Update docs when behavior changes

## Before Coding
- Ensure requirements are clear
- Check for existing patterns in codebase
- Consider security implications
- Plan for testing

## Commit Standards
- Atomic commits (one concern per commit)
- Prefix: feat: | fix: | docs: | chore:
- One-sentence summary
- No mixed or vague commits

## When Stuck
- Ask for clarification rather than guessing
- Reference relevant steering files
- Suggest alternatives if blocked
```

## Error Handling

### Frontend Error States

1. **Rate Limit Exceeded (429)**
   - Display countdown timer with `retryAfter` seconds
   - Disable submit buttons until timer expires
   - Show friendly message: "Too many requests. Please wait X seconds."

2. **Generation Failed (500)**
   - Display error message with retry option
   - Offer to start over from beginning
   - Log error details for debugging

3. **Timeout (504)**
   - Display timeout message
   - Suggest refreshing and trying again
   - Consider breaking into smaller requests

4. **Invalid Response**
   - Backend retries once automatically
   - If still invalid, show generic error
   - Log malformed response for analysis

### Backend Validation

1. **Input Validation**
   - Project idea: max 2000 chars, non-empty
   - Answers: max 1000 chars each
   - Experience level: must be valid enum
   - Hook preset: must be valid enum

2. **Output Validation**
   - Steering files: valid frontmatter YAML
   - Hook files: valid JSON matching schema
   - Kickoff prompt: contains required sections

## Testing Strategy

### Unit Tests

- Component rendering tests for new UI components
- State management tests for experience level persistence
- Validation function tests for file format checking

### Property-Based Tests

Using a property-based testing library (e.g., fast-check for TypeScript tests):

1. **Experience Level Adaptation Test**
   - Generate random project ideas
   - For each experience level, verify question content matches level expectations
   - Minimum 100 iterations

2. **Steering File Validity Test**
   - Generate outputs for random project configurations
   - Verify all steering files have valid frontmatter
   - Verify conditional files have correct patterns
   - Minimum 100 iterations

3. **Hook Schema Validity Test**
   - Generate hooks for all presets
   - Parse JSON and verify schema compliance
   - Verify runCommand restrictions
   - Minimum 100 iterations

4. **Kickoff Prompt Completeness Test**
   - Generate kickoff prompts for various projects
   - Verify all required sections present
   - Verify "no coding" enforcement phrase
   - Minimum 100 iterations

5. **Question Constraints Test**
   - Generate questions for random projects
   - Verify count is 5-10
   - Verify ordering follows logical sequence
   - Minimum 100 iterations

### Integration Tests

- End-to-end flow from level selection to file download
- API endpoint tests with various inputs
- Error handling scenarios

# Phase 2: Feature Implementation — Requirements

## Introduction

Phase 2 implements the three core generators of the Kiro Prompting & Guardrails Generator:

- Kickoff Prompt Generator: Generates a single prompt artifact enforcing answer-first, no-coding-first thinking
- Steering Document Generator: Generates `.kiro/steering/` files with proper frontmatter and inclusion modes
- Hooks Generator: Generates `.kiro/hooks/*.kiro.hook` files with preset-based configurations
- All generators produce downloadable/copyable output artifacts
- UI guides users through configuration before generation
- No repo scanning (future module, out of scope)
- No AI integration for generation (templates only for MVP)

## User Stories

### US-1: Kickoff Prompt Generation
As a beginner developer, I want to generate a kickoff prompt that forces me to answer critical questions before coding so that I think through my project properly.

### US-2: Steering Document Generation
As a developer setting up a new project, I want to generate steering files with correct Kiro frontmatter so that my AI assistant follows project conventions.

### US-3: Hooks Generation
As a developer, I want to generate hook files from presets so that I have automated quality checks without manual configuration.

### US-4: Output Artifacts
As a user, I want to preview, copy, and download generated artifacts so that I can use them in my project.

## Acceptance Criteria (EARS Format)

### Kickoff Prompt Generator

#### AC-1: Question Flow UI
WHEN a user starts the kickoff prompt generator
THE SYSTEM SHALL display questions in strict order:
1. Project Identity
2. Success Criteria
3. Users & Roles
4. Data Sensitivity (including Data Lifecycle: retention, deletion, export, audit logging, backups)
5. Auth Model
6. Concurrency Expectations
7. Risks & Tradeoffs (Top 3 risks, simplest mitigations, what is explicitly not handled)
8. Boundaries (including 2-3 concrete access examples)
9. Non-Goals
10. Constraints

#### AC-2: Answer Enforcement
WHEN a user attempts to proceed without answering a required question
THE SYSTEM SHALL block progression and indicate the missing answer.

#### AC-3: Prompt Generation
WHEN a user completes all questions and requests generation
THE SYSTEM SHALL produce a single cohesive prompt artifact that:
- Lists all questions with user answers
- Instructs that no coding is allowed until questions are answered
- Enforces answer-first, no-coding-first

### Steering Document Generator

#### AC-4: Foundation Steering Files
WHEN a user generates steering documents
THE SYSTEM SHALL create files with correct frontmatter:
- `product.md` with `inclusion: always`
- `tech.md` with `inclusion: always`
- `structure.md` with `inclusion: always`

#### AC-5: Conditional Steering Files
WHEN a user enables conditional steering
THE SYSTEM SHALL create files with fileMatch inclusion:
- `security-go.md` with `inclusion: fileMatch`, `fileMatchPattern: "**/*.go"`
- `security-web.md` with `inclusion: fileMatch`, `fileMatchPattern: "**/*.ts"`
- `quality-go.md` with `inclusion: fileMatch`, `fileMatchPattern: "**/*.go"`
- `quality-web.md` with `inclusion: fileMatch`, `fileMatchPattern: "**/*.tsx"`

#### AC-6: AGENTS.md Generation
WHEN a user generates steering documents
THE SYSTEM SHALL create `AGENTS.md` at repo root containing:
- Always follow steering
- Never invent requirements
- Prefer small, reviewable changes
- Update docs when behavior changes

#### AC-7: Steering Content
WHEN steering files are generated
THE SYSTEM SHALL include content that is short, strict, and practical.

### Hooks Generator

#### AC-8: Hook File Format
WHEN a user generates hooks
THE SYSTEM SHALL create files with `.kiro.hook` extension containing valid schema:
- `name`, `description`, `version`, `enabled`
- `when` with valid `type` (fileEdited, fileCreated, fileDeleted, promptSubmit, agentStop, userTriggered)
- `then` with valid `type` (askAgent everywhere, runCommand only for promptSubmit/agentStop)

#### AC-9: Preset Light
WHEN a user selects Light preset
THE SYSTEM SHALL generate hooks for:
- `agentStop`: runCommand formatters (Go + frontend)

#### AC-10: Preset Basic
WHEN a user selects Basic preset
THE SYSTEM SHALL generate hooks for:
- `agentStop`: formatters + basic linters
- `userTriggered`: run unit tests and summarize results

#### AC-11: Preset Default
WHEN a user selects Default preset
THE SYSTEM SHALL generate hooks for:
- Everything in Basic
- `agentStop`: quick secret scan
- `promptSubmit`: prompt guardrails (block unsafe prompts, require confirmation)

#### AC-12: Preset Strict
WHEN a user selects Strict preset
THE SYSTEM SHALL generate hooks for:
- Everything in Default
- `agentStop`: static analysis
- Manual or `agentStop`: dependency vulnerability scan

### Output Artifacts

#### AC-13: Preview Mode
WHEN a user requests preview
THE SYSTEM SHALL display generated content before finalizing.

#### AC-14: Copy to Clipboard
WHEN a user clicks copy
THE SYSTEM SHALL copy the artifact content to clipboard.

#### AC-15: Download Files
WHEN a user clicks download
THE SYSTEM SHALL download artifacts as files with correct names and paths.

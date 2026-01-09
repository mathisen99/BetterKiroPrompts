# Phase 3: Polish & Testing — Requirements

## Introduction

Phase 3 completes the Kiro Prompting & Guardrails Generator:

- UI/UX polish: consistent dark theme with blue base, error handling, loading states
- Accessibility compliance (WCAG 2.1 AA)
- Missing features: manual steering (`inclusion: manual`), file references (`#[[file:<path>]]`)
- Commit message contract enforcement (atomic, prefixed, one-sentence)
- Comprehensive testing: unit, integration, E2E
- Documentation: README, API docs, user guide
- Definition of Done verification
- (OPTIONAL) Repo scanning module

## User Stories

### US-1: Polished User Experience
As a user, I want consistent visual feedback and error handling so that I understand what the application is doing.

### US-2: Accessible Interface
As a user with accessibility needs, I want the application to be keyboard navigable and screen reader compatible.

### US-3: Manual Steering
As a developer, I want to generate manual steering files that I can reference with `#steering-file-name`.

### US-4: Verified Quality
As a maintainer, I want comprehensive test coverage so that I can refactor with confidence.

### US-5: Clear Documentation
As a new user, I want setup instructions and usage guides so that I can use the tool effectively.

## Acceptance Criteria (EARS Format)

### UI/UX Polish

#### AC-1: Error Handling
WHEN an API request fails
THE SYSTEM SHALL display a user-friendly error message with retry option.

#### AC-2: Loading States
WHEN an API request is in progress
THE SYSTEM SHALL display loading indicators on affected components.

#### AC-3: Toast Notifications
WHEN a user action succeeds (copy, download)
THE SYSTEM SHALL display a brief success notification.

### Accessibility

#### AC-4: Keyboard Navigation
WHEN a user navigates using keyboard only
THE SYSTEM SHALL allow complete wizard flow without mouse.

#### AC-5: Screen Reader Support
WHEN a screen reader is active
THE SYSTEM SHALL announce form labels, errors, and state changes.

#### AC-6: Color Contrast
WHEN displaying text and interactive elements
THE SYSTEM SHALL meet WCAG 2.1 AA contrast ratios.

### Missing Features

#### AC-7: Manual Steering Generation
WHEN a user enables manual steering
THE SYSTEM SHALL generate files with `inclusion: manual` frontmatter.

#### AC-8: File References
WHEN a user adds file references to steering
THE SYSTEM SHALL include `#[[file:<relative_path>]]` syntax in output.

#### AC-9: Commit Contract Display
WHEN generating any output
THE SYSTEM SHALL display the commit message contract:
- Atomic (one concern per commit)
- Prefixed (feat:, fix:, docs:, chore:)
- One-sentence summary

### Testing

#### AC-10: Unit Test Coverage
WHEN running unit tests
THE SYSTEM SHALL have tests for all generator functions in `backend/internal/generator/`.

#### AC-11: Integration Test Coverage
WHEN running integration tests
THE SYSTEM SHALL have tests for all API endpoints.

#### AC-12: E2E Test Coverage
WHEN running E2E tests
THE SYSTEM SHALL have tests for:
- Complete kickoff wizard flow
- Complete steering generation flow
- Complete hooks generation flow

### Documentation

#### AC-13: README
WHEN a developer clones the repo
THE SYSTEM SHALL have a README with setup instructions and quick start.

#### AC-14: API Documentation
WHEN a developer needs API details
THE SYSTEM SHALL have documentation for all endpoints with request/response examples.

#### AC-15: User Guide
WHEN a user wants to understand generated outputs
THE SYSTEM SHALL have a guide explaining kickoff prompts, steering files, and hooks.

### Definition of Done (from plan)

#### AC-16: Kickoff Prompt Generation
WHEN a user completes the kickoff wizard
THE SYSTEM SHALL generate a full Kiro kickoff prompt.

#### AC-17: Steering Files
WHEN a user generates steering
THE SYSTEM SHALL produce usable, correctly scoped, concise files.

#### AC-18: Hooks Files
WHEN a user generates hooks
THE SYSTEM SHALL produce valid, usable hook files.

#### AC-19: (OPTIONAL) Repo Scanning
WHEN a user triggers repo scanning
THE SYSTEM SHALL run tools end-to-end and produce structured output.

## OPTIONAL: Repo Scanning Requirements

### AC-20: Read-Only Clone
WHEN scanning a repository
THE SYSTEM SHALL clone read-only without write access.

### AC-21: Tool Execution
WHEN scanning
THE SYSTEM SHALL run: TruffleHog, Gitleaks, osv-scanner, govulncheck.

### AC-22: Hard Timeouts
WHEN a scan tool exceeds timeout
THE SYSTEM SHALL terminate and return partial results.

### AC-23: No Outbound Network
WHEN scanning
THE SYSTEM SHALL block outbound network access during tool execution.

### AC-24: Structured Output
WHEN scan completes
THE SYSTEM SHALL return severity-ranked findings with file/line references.

### AC-25: AI Summary
WHEN presenting results
THE SYSTEM SHALL use AI to summarize tool output only (not invent vulnerabilities).

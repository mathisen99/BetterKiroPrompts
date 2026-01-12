# Requirements Document

## Introduction

This phase delivers a complete overhaul of the Kiro Prompting & Guardrails Generator to make it production-ready. The focus is on three areas: professional UI design with shadcn blue dark theme, significantly smarter AI generation that adapts to user experience levels, and complete coverage of all Kiro file types with proper formats.

## Glossary

- **Generator**: The web application that generates Kiro kickoff prompts, steering files, and hooks
- **Experience_Level**: User's programming proficiency (beginner, novice, expert)
- **Steering_File**: Markdown files in `.kiro/steering/` that provide persistent context to Kiro
- **Hook_File**: JSON files in `.kiro/hooks/` that automate agent actions on events
- **Kickoff_Prompt**: A structured prompt that enforces thinking before coding
- **Frontmatter**: YAML metadata at the top of steering files defining inclusion mode
- **Inclusion_Mode**: How steering files are loaded (always, fileMatch, manual)

## Requirements

### Requirement 1: Experience Level Selection

**User Story:** As a user, I want to select my programming experience level, so that the questions and suggestions are appropriate for my skill level.

#### Acceptance Criteria

1. WHEN a user visits the landing page, THE Generator SHALL display three experience level options: Beginner, Novice, and Expert
2. WHEN a user selects "Beginner", THE Generator SHALL adapt questions to avoid technical jargon and provide more guidance
3. WHEN a user selects "Novice", THE Generator SHALL use moderate technical language with helpful hints
4. WHEN a user selects "Expert", THE Generator SHALL use technical terminology and assume familiarity with architecture concepts
5. THE Generator SHALL persist the selected experience level throughout the session
6. WHEN a beginner describes an overly complex project, THE Generator SHALL suggest simpler alternatives or phased approaches

### Requirement 2: Professional UI Design

**User Story:** As a user, I want a clean, professional dark interface with blue accents, so that the tool feels trustworthy and modern.

#### Acceptance Criteria

1. THE Generator SHALL use a dark background theme matching shadcn's blue theme palette
2. THE Generator SHALL use blue (#3B82F6 primary) as the accent color for interactive elements
3. THE Generator SHALL display cards with subtle borders and proper spacing
4. THE Generator SHALL use consistent typography with clear hierarchy
5. WHEN displaying forms, THE Generator SHALL use properly styled inputs with focus states
6. THE Generator SHALL include a header with branding and navigation
7. THE Generator SHALL be fully responsive on mobile, tablet, and desktop

### Requirement 3: Smarter Question Generation

**User Story:** As a user, I want questions that are relevant to my project and experience level, so that I can provide meaningful answers.

#### Acceptance Criteria

1. WHEN generating questions, THE Generator SHALL consider the user's experience level
2. WHEN a beginner describes wanting to build a SaaS, THE Generator SHALL ask about core features before technical choices
3. WHEN an expert describes a project, THE Generator SHALL ask about architecture, scalability, and data consistency
4. THE Generator SHALL generate between 5-10 questions based on project complexity
5. WHEN a question has multiple valid approaches, THE Generator SHALL provide example suggestions
6. THE Generator SHALL ask questions in a logical order: identity → users → data → auth → architecture → constraints

### Requirement 4: Complete Steering File Generation

**User Story:** As a user, I want all necessary steering files generated with proper Kiro format, so that I can use them directly in my project.

#### Acceptance Criteria

1. THE Generator SHALL generate `product.md` with frontmatter `inclusion: always`
2. THE Generator SHALL generate `tech.md` with frontmatter `inclusion: always`
3. THE Generator SHALL generate `structure.md` with frontmatter `inclusion: always`
4. THE Generator SHALL generate `security-{lang}.md` with frontmatter `inclusion: fileMatch` and appropriate `fileMatchPattern`
5. THE Generator SHALL generate `quality-{lang}.md` with frontmatter `inclusion: fileMatch` and appropriate `fileMatchPattern`
6. WHEN the project uses Go, THE Generator SHALL generate `security-go.md` with `fileMatchPattern: "**/*.go"`
7. WHEN the project uses TypeScript/React, THE Generator SHALL generate `security-web.md` with `fileMatchPattern: "**/*.{ts,tsx}"`
8. THE Generator SHALL ensure all steering files follow Kiro's documented format

### Requirement 5: Valid Hook File Generation

**User Story:** As a user, I want valid Kiro hook files that work immediately, so that I can automate my development workflow.

#### Acceptance Criteria

1. THE Generator SHALL generate hook files with `.kiro.hook` extension
2. THE Generator SHALL include required fields: name, description, version, enabled, when, then
3. WHEN generating file-based hooks, THE Generator SHALL include `patterns` in the `when` block
4. THE Generator SHALL use valid `when.type` values: fileEdited, fileCreated, fileDeleted, promptSubmit, agentStop, userTriggered
5. THE Generator SHALL use valid `then.type` values: askAgent (for all triggers) or runCommand (only for promptSubmit/agentStop)
6. THE Generator SHALL generate hooks appropriate to the project's tech stack
7. THE Generator SHALL offer hook presets: Light, Basic, Default, Strict

### Requirement 6: Kickoff Prompt Quality

**User Story:** As a user, I want a comprehensive kickoff prompt that enforces thinking before coding, so that I start my project with clarity.

#### Acceptance Criteria

1. THE Generator SHALL generate a kickoff prompt that explicitly states "no coding until questions are answered"
2. THE Generator SHALL include all required sections: Project Identity, Success Criteria, Users & Roles, Data Sensitivity, Auth Model, Concurrency, Boundaries, Non-Goals, Constraints
3. WHEN the project involves sensitive data, THE Generator SHALL include Data Lifecycle section
4. THE Generator SHALL include Risks & Tradeoffs section with top 3 risks and mitigations
5. THE Generator SHALL include concrete Boundary Examples showing who can access what
6. THE Generator SHALL adapt language complexity to the user's experience level

### Requirement 7: Output Preview and Download

**User Story:** As a user, I want to preview, edit, and download all generated files, so that I can customize them before use.

#### Acceptance Criteria

1. THE Generator SHALL display generated files in organized tabs by type (Kickoff, Steering, Hooks)
2. THE Generator SHALL allow inline editing of any generated file
3. THE Generator SHALL show a "Modified" indicator when a file has been edited
4. THE Generator SHALL provide a "Reset" button to restore original content
5. THE Generator SHALL allow downloading individual files
6. THE Generator SHALL allow downloading all files as a ZIP archive
7. THE Generator SHALL preserve the correct directory structure in the ZIP

### Requirement 8: AI Knowledge Base for File Formats

**User Story:** As a system, I need comprehensive knowledge about Kiro file formats and best practices, so that I can generate high-quality, valid files.

#### Acceptance Criteria

1. THE Generator SHALL include detailed system prompts with Kiro steering file format specifications
2. THE Generator SHALL include the complete Kiro hook schema with all valid field values
3. THE Generator SHALL include examples of well-structured steering files for each type (product, tech, structure, security, quality)
4. THE Generator SHALL include frontmatter format rules: `inclusion: always | fileMatch | manual` and `fileMatchPattern` for conditional files
5. THE Generator SHALL include hook trigger types and their valid combinations with action types
6. THE Generator SHALL include best practices for each file type (concise, actionable, no secrets, etc.)
7. THE Generator SHALL include kickoff prompt structure with all required sections and their purposes
8. THE Generator SHALL validate generated files against known schemas before returning them

### Requirement 9: Error Handling and Feedback

**User Story:** As a user, I want clear feedback when something goes wrong, so that I know how to proceed.

#### Acceptance Criteria

1. WHEN a rate limit is exceeded, THE Generator SHALL display a countdown timer
2. WHEN generation fails, THE Generator SHALL display a user-friendly error message
3. WHEN the AI returns invalid output, THE Generator SHALL retry once before showing an error
4. THE Generator SHALL show loading states with estimated time during generation
5. IF an error occurs, THEN THE Generator SHALL provide a clear path to retry or start over

### Requirement 10: Comprehensive File Coverage

**User Story:** As a user, I want all the steering and hook files I need for a complete Kiro setup, so that I don't have to create them manually.

#### Acceptance Criteria

1. THE Generator SHALL generate a minimum of 3 steering files (product.md, tech.md, structure.md)
2. WHEN the project involves backend code, THE Generator SHALL generate security and quality steering for that language
3. WHEN the project involves frontend code, THE Generator SHALL generate security and quality steering for web technologies
4. THE Generator SHALL generate at least one hook file appropriate to the project
5. THE Generator SHALL generate an AGENTS.md file at repo root with agent behavior guidelines
6. THE Generator SHALL offer hook preset selection (Light, Basic, Default, Strict) with clear descriptions of each

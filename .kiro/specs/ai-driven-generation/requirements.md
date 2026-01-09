# Requirements Document

## Introduction

This feature transforms the BetterKiroPrompts application from a multi-page hardcoded wizard into a single-page AI-driven experience. Users enter their project idea, the AI generates contextual follow-up questions, and then produces tailored kickoff prompts, steering files, and hooks—all in one flow. The old multi-page navigation and template-based generation will be removed.

## Glossary

- **Landing_Page**: The single page where users interact with the entire generation flow
- **Project_Idea**: The initial free-text description of what the user wants to build
- **Question_Plan**: AI-generated list of follow-up questions tailored to the project idea
- **Generation_Service**: Backend service that communicates with GPT-5.2 API
- **Output_Editor**: UI component allowing manual editing of generated files before download
- **Rate_Limiter**: Backend component that limits API calls per client
- **Session**: A single generation flow from project idea to final output (no regeneration)

## Requirements

### Requirement 1: Single Landing Page

**User Story:** As a user, I want to access all functionality from one page, so that I can complete the entire flow without navigating between pages.

#### Acceptance Criteria

1. WHEN a user visits the application, THE Landing_Page SHALL display a single input field asking "What project do you want to make?"
2. THE Landing_Page SHALL display 3-5 example project ideas as clickable suggestions below the input field
3. WHEN the user submits a project idea, THE Landing_Page SHALL transition to the question phase without page navigation
4. THE Landing_Page SHALL remove the old navigation menu (Kickoff, Steering, Hooks tabs)

### Requirement 2: AI-Generated Question Plan

**User Story:** As a user, I want the system to ask me relevant questions based on my project idea, so that I don't have to answer irrelevant hardcoded questions.

#### Acceptance Criteria

1. WHEN a project idea is submitted, THE Generation_Service SHALL call GPT-5.2 to generate a Question_Plan
2. THE Question_Plan SHALL contain 5-10 questions tailored to the specific project type
3. THE Question_Plan SHALL adapt question complexity based on project sophistication (beginner vs advanced indicators)
4. WHEN generating questions for a simple project (e.g., "snake game"), THE Generation_Service SHALL focus on scope, platform, and basic features
5. WHEN generating questions for a complex project (e.g., "distributed event-sourcing system"), THE Generation_Service SHALL include architecture, scalability, and data consistency questions
6. THE Landing_Page SHALL display questions one at a time with previous answers visible
7. THE Landing_Page SHALL allow users to go back and edit previous answers

### Requirement 3: AI-Generated Outputs

**User Story:** As a user, I want all outputs (kickoff prompt, steering files, hooks) generated based on my specific project, so that they are relevant and useful.

#### Acceptance Criteria

1. WHEN all questions are answered, THE Generation_Service SHALL call GPT-5.2 to generate all outputs in a single API call
2. THE Generation_Service SHALL generate a kickoff prompt tailored to the project
3. THE Generation_Service SHALL generate steering files appropriate for the project's tech stack and complexity
4. THE Generation_Service SHALL generate hooks appropriate for the project type
5. THE Landing_Page SHALL display all generated files in a tabbed or accordion interface
6. IF the AI fails to generate valid output, THEN THE Generation_Service SHALL return a structured error message

### Requirement 4: Manual Editing Before Download

**User Story:** As a user, I want to edit generated files before downloading, so that I can make adjustments without regenerating.

#### Acceptance Criteria

1. THE Output_Editor SHALL display each generated file with syntax highlighting
2. THE Output_Editor SHALL allow inline text editing of any generated file
3. WHEN a user edits a file, THE Output_Editor SHALL preserve the changes in browser state
4. THE Output_Editor SHALL provide a "Reset to Original" button per file
5. THE Landing_Page SHALL NOT provide a regenerate button (one-time flow for cost efficiency)

### Requirement 5: Download Functionality

**User Story:** As a user, I want to download all generated files, so that I can use them in my project.

#### Acceptance Criteria

1. THE Landing_Page SHALL provide a "Download All" button that creates a ZIP file
2. THE Landing_Page SHALL provide individual "Copy" and "Download" buttons per file
3. WHEN downloading, THE system SHALL use the edited content if modifications were made
4. THE ZIP file SHALL preserve the correct directory structure (.kiro/steering/, .kiro/hooks/)

### Requirement 6: API Security and Rate Limiting

**User Story:** As a system operator, I want to protect the API from abuse, so that costs are controlled and service remains available.

#### Acceptance Criteria

1. THE Generation_Service SHALL store the OpenAI API key in environment variables only
2. THE Generation_Service SHALL never log API keys or full prompts containing sensitive data
3. THE Rate_Limiter SHALL limit requests to 10 generations per IP per hour
4. WHEN rate limit is exceeded, THE Generation_Service SHALL return HTTP 429 with retry-after header
5. THE Generation_Service SHALL validate all input before sending to OpenAI API
6. THE Generation_Service SHALL set a timeout of 60 seconds for OpenAI API calls
7. IF OpenAI API returns an error, THEN THE Generation_Service SHALL return a user-friendly error without exposing internal details

### Requirement 7: Cleanup of Old Code

**User Story:** As a developer, I want unused code removed, so that the codebase remains maintainable.

#### Acceptance Criteria

1. THE system SHALL remove the KickoffPage, SteeringPage, and HooksPage components
2. THE system SHALL remove the KickoffWizard and all kickoff-specific components
3. THE system SHALL remove the SteeringConfigurator and all steering-specific components
4. THE system SHALL remove the HooksPresetSelector and all hooks-specific components
5. THE system SHALL remove the Navigation component
6. THE system SHALL remove the old /api/kickoff/generate, /api/steering/generate, /api/hooks/generate endpoints
7. THE system SHALL remove all .tmpl template files from backend/internal/templates/
8. THE system SHALL remove the old generator package functions (GenerateKickoff, GenerateSteering, GenerateHooks)
9. THE system SHALL update the API client (lib/api.ts) to remove old types and functions

### Requirement 8: Error Handling and Loading States

**User Story:** As a user, I want clear feedback during generation, so that I know the system is working.

#### Acceptance Criteria

1. WHILE waiting for AI response, THE Landing_Page SHALL display a loading indicator with estimated wait time
2. WHEN an error occurs, THE Landing_Page SHALL display a clear error message
3. THE Landing_Page SHALL NOT provide retry functionality (one-time flow)
4. IF the session times out, THEN THE Landing_Page SHALL instruct the user to refresh and start over

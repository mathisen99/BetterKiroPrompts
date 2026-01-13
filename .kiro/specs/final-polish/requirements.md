# Requirements Document

## Introduction

This specification covers the final polish phase for BetterKiroPrompts, focusing on UI/UX improvements, syntax highlighting, error recovery, generation persistence, and a public gallery with ratings. The goal is to transform the MVP into a production-ready application that handles multiple concurrent users securely.

## Glossary

- **Generation**: A complete set of generated files (kickoff prompt, steering files, hooks, AGENTS.md) created from a user's project idea and answers
- **Gallery**: A public browsable collection of past generations that users can view and rate
- **Session_State**: The current user's progress through the generation flow, including project idea, answers, and experience level
- **Syntax_Highlighter**: A component that renders code with language-appropriate color highlighting
- **Rating**: A 1-5 star score that users can assign to gallery items
- **LocalStorage_Manager**: A service that persists and retrieves session state from browser localStorage

## Requirements

### Requirement 1: Improved UI Flow

**User Story:** As a user, I want a cleaner, more focused interface during the question and generation phases, so that I can concentrate on providing quality answers without visual distractions.

#### Acceptance Criteria

1. WHEN the user progresses past the level-select phase, THE App SHALL hide the large centered logo with a smooth fade-out animation
2. WHEN the user is in the questions phase, THE App SHALL display a compact header with a small logo that links back to start
3. WHEN the user navigates between phases, THE App SHALL apply smooth transition animations to content changes
4. THE QuestionFlow component SHALL maintain consistent card styling and spacing throughout the question sequence
5. WHEN the user completes generation successfully, THE App SHALL show a success celebration animation before displaying results

### Requirement 2: Syntax Highlighting for Output

**User Story:** As a user, I want to see generated JSON and Markdown files with proper syntax highlighting, so that I can easily read and understand the content.

#### Acceptance Criteria

1. THE OutputEditor SHALL render JSON files with syntax highlighting for keys, strings, numbers, and brackets
2. THE OutputEditor SHALL render Markdown files with syntax highlighting for headers, links, code blocks, and emphasis
3. THE Syntax_Highlighter SHALL support a dark theme that matches the application's color scheme
4. WHEN a user edits a file, THE OutputEditor SHALL update syntax highlighting in real-time
5. THE OutputEditor SHALL provide a toggle between highlighted view and plain text edit mode
6. THE Syntax_Highlighter SHALL handle malformed JSON/Markdown gracefully without crashing

### Requirement 3: Extended Timeouts

**User Story:** As a user, I want the application to wait longer for AI generation, so that complex projects don't fail due to timeout.

#### Acceptance Criteria

1. THE OpenAI_Client SHALL use a 120-second timeout for API requests
2. THE Frontend SHALL display a progress indicator that updates during long-running requests
3. WHEN a request exceeds 90 seconds, THE LoadingState SHALL display a message indicating the request is taking longer than usual
4. THE Frontend SHALL use AbortController to properly cancel requests on user navigation or timeout
5. IF a timeout occurs, THEN THE System SHALL provide a clear error message with retry option

### Requirement 4: State Persistence and Error Recovery

**User Story:** As a user, I want my progress saved automatically, so that I don't lose my work if generation fails or I accidentally close the browser.

#### Acceptance Criteria

1. THE LocalStorage_Manager SHALL save session state after each user action (level selection, project idea submission, each answer)
2. WHEN the application loads, THE System SHALL check for existing session state and offer to restore it
3. WHEN generation fails, THE System SHALL preserve all user inputs and allow retry without re-entering data
4. WHEN generation succeeds, THE LocalStorage_Manager SHALL clear the saved session state
5. THE LocalStorage_Manager SHALL store state with a timestamp and expire entries older than 24 hours
6. THE System SHALL provide a "Start Over" button that clears saved state and resets the flow
7. IF localStorage is unavailable, THEN THE System SHALL continue functioning without persistence

### Requirement 5: Generation Storage in Database

**User Story:** As a user, I want my successful generations saved to the database, so that they can be displayed in the public gallery.

#### Acceptance Criteria

1. WHEN generation completes successfully, THE Backend SHALL store the generation in PostgreSQL with all metadata
2. THE Generation record SHALL include: unique ID, project idea, experience level, hook preset, generated files, creation timestamp, and category tags
3. THE Backend SHALL automatically categorize generations based on project idea keywords (e.g., "API", "CLI", "Web App", "Mobile")
4. THE Backend SHALL use parameterized queries for all database operations
5. THE Backend SHALL handle concurrent generation storage without data corruption
6. THE Generation storage SHALL NOT include any user-identifying information (anonymous by default)

### Requirement 6: Public Gallery

**User Story:** As a user, I want to browse a gallery of past generations, so that I can find inspiration and see examples of good prompts.

#### Acceptance Criteria

1. THE Gallery page SHALL display a paginated list of past generations with project idea preview and category tags
2. THE Gallery SHALL support filtering by category (API, CLI, Web App, Mobile, Other)
3. THE Gallery SHALL support sorting by: newest, highest rated, most viewed
4. WHEN a user clicks a gallery item, THE System SHALL display the full generation details in a modal or detail page
5. THE Gallery SHALL load efficiently with pagination (20 items per page)
6. THE Gallery detail view SHALL allow users to copy or download the generated files

### Requirement 7: Rating System

**User Story:** As a user, I want to rate generations in the gallery, so that the community can identify the most useful examples.

#### Acceptance Criteria

1. THE Rating component SHALL display a 1-5 star rating interface
2. WHEN a user submits a rating, THE Backend SHALL store it associated with the generation ID
3. THE Backend SHALL calculate and cache average ratings for efficient retrieval
4. THE System SHALL use browser fingerprinting or localStorage to prevent duplicate ratings from the same user
5. THE Gallery SHALL display average rating and total vote count for each item
6. THE Rating submission SHALL be rate-limited to prevent abuse (max 20 ratings per hour per IP)

### Requirement 8: Concurrency and Performance

**User Story:** As a system operator, I want the application to handle multiple concurrent users, so that the service remains responsive under load.

#### Acceptance Criteria

1. THE Backend SHALL handle multiple concurrent generation requests without blocking
2. THE PostgreSQL connection SHALL use connection pooling with appropriate limits
3. THE Backend SHALL implement request queuing for OpenAI API calls to respect rate limits
4. THE Backend SHALL use goroutines appropriately for concurrent operations
5. THE System SHALL return appropriate HTTP status codes when under heavy load (503 with Retry-After)
6. THE Backend SHALL implement graceful shutdown to complete in-flight requests

### Requirement 9: Security Hardening

**User Story:** As a system operator, I want the application to be secure, so that user data and API keys are protected.

#### Acceptance Criteria

1. THE Backend SHALL never expose API keys in responses or logs
2. THE Backend SHALL sanitize all user input before storage to prevent XSS
3. THE Backend SHALL validate all input at API boundaries with appropriate length limits
4. THE Gallery display SHALL escape all user-generated content when rendering
5. THE Backend SHALL implement request ID tracking for security event logging
6. THE Backend SHALL log security-relevant events (rate limit hits, validation failures) without logging sensitive data
7. THE System SHALL use HTTPS for all external resources

### Requirement 10: Error Handling Improvements

**User Story:** As a user, I want clear and helpful error messages, so that I understand what went wrong and how to fix it.

#### Acceptance Criteria

1. WHEN an error occurs, THE System SHALL display a user-friendly message with suggested actions
2. THE Error display SHALL include a "Try Again" button that retries the failed operation
3. THE Error display SHALL include a "Start Over" button that resets the flow
4. IF the error is recoverable, THEN THE System SHALL automatically retry once before showing the error
5. THE Backend SHALL return structured error responses with error codes for frontend handling
6. THE System SHALL distinguish between client errors (bad input) and server errors (try again later)

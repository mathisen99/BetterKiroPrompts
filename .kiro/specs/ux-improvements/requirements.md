# Requirements Document

## Introduction

This spec addresses critical UX and functionality improvements identified during user testing. The main issues are: questions not being differentiated enough by experience level, lack of loading feedback during long generations, missing clickable answer examples, view/vote abuse prevention, timeout issues, and navigation clarity problems.

## Glossary

- **Experience_Level**: User's programming skill level (beginner, novice, expert)
- **Question_Generator**: The AI system that generates follow-up questions based on project idea
- **Answer_Example**: A pre-generated clickable option that users can select instead of typing
- **View_Count**: Number of unique views a gallery item has received
- **Vote**: A rating (1-5 stars) given to a gallery generation
- **IP_Fingerprint**: A hash of the user's IP address used for abuse prevention

## Requirements

### Requirement 1: Experience-Level-Appropriate Questions

**User Story:** As a beginner user, I want questions that use simple language and avoid technical jargon, so that I can understand and answer them without prior programming knowledge.

#### Acceptance Criteria

1. WHEN a beginner selects their experience level, THE Question_Generator SHALL produce questions using only non-technical language
2. WHEN a beginner receives questions, THE Question_Generator SHALL avoid terms like "API", "database schema", "authentication flow", "microservices", "CI/CD", "containerization"
3. WHEN a novice selects their experience level, THE Question_Generator SHALL use moderate technical terms with brief explanations
4. WHEN an expert selects their experience level, THE Question_Generator SHALL use full technical terminology and ask about architecture patterns
5. THE Question_Generator SHALL produce distinctly different question sets for each experience level given the same project idea

### Requirement 2: Clickable Answer Examples

**User Story:** As a user, I want to see example answers I can click to select, so that I can quickly answer questions without typing if an example fits my needs.

#### Acceptance Criteria

1. WHEN questions are displayed, THE System SHALL show 3 clickable example answers for each question
2. WHEN a user clicks an example answer, THE System SHALL populate that answer for the question
3. WHEN a user has selected an example, THE System SHALL allow them to edit the selected text
4. WHEN a user prefers to type, THE System SHALL still allow free-text input in the answer field
5. THE example answers SHALL be appropriate to the user's experience level

### Requirement 3: Loading Feedback During Generation

**User Story:** As a user, I want to see clear feedback that generation is in progress, so that I know the system is working during long waits.

#### Acceptance Criteria

1. WHEN question generation starts, THE System SHALL display a loading indicator with message "Generating questions... This may take up to 2 minutes"
2. WHEN output generation starts, THE System SHALL display a loading indicator with message "Generating your files... This may take up to 3 minutes"
3. WHILE generation is in progress, THE System SHALL show an animated spinner or progress indicator
4. IF generation takes longer than 30 seconds, THE System SHALL display an encouraging message like "Still working..."

### Requirement 4: Increased Timeout Configuration

**User Story:** As a system administrator, I want longer timeouts for API requests, so that slow generations don't fail unnecessarily.

#### Acceptance Criteria

1. THE Backend SHALL set API request timeout to 180 seconds (3 minutes)
2. THE Frontend SHALL set fetch timeout to 180 seconds (3 minutes)
3. THE OpenAI_Client SHALL set HTTP client timeout to 180 seconds (3 minutes)
4. WHEN a timeout occurs, THE System SHALL display a user-friendly error message suggesting retry

### Requirement 5: IP-Based View and Vote Abuse Prevention

**User Story:** As a gallery curator, I want views and votes to be limited per IP address, so that metrics reflect genuine engagement.

#### Acceptance Criteria

1. WHEN a user views a gallery item, THE System SHALL record the view only once per IP address per item
2. WHEN a user votes on a gallery item, THE System SHALL allow only one vote per IP address per item
3. WHEN the same IP attempts to view an item again, THE System SHALL NOT increment the view count
4. WHEN the same IP attempts to vote again, THE System SHALL update their existing vote instead of creating a new one
5. THE System SHALL use IP address hashing for privacy

### Requirement 6: Improved Navigation Visibility

**User Story:** As a user, I want clear and visible navigation elements, so that I can easily move between pages and close modals.

#### Acceptance Criteria

1. WHEN on the landing page, THE System SHALL display a prominent "Browse Gallery" button with high contrast
2. WHEN on the gallery page, THE System SHALL display a clearly visible "Back to Home" or logo link
3. WHEN a gallery item modal is open, THE System SHALL display a prominent close button (X) in the top-right corner
4. THE close button SHALL have sufficient size (minimum 44x44px touch target) and contrast
5. WHEN user clicks outside the modal, THE System SHALL close the modal
6. WHEN user presses Escape key, THE System SHALL close the modal

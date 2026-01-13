# Implementation Plan: UX Improvements

## Overview

This plan implements six improvement areas: experience-level-appropriate questions, clickable answer examples, loading feedback, increased timeouts, IP-based abuse prevention, and navigation visibility. Tasks are ordered to build foundational changes first, then layer on features.

## Tasks

- [x] 1. Update timeout configuration
  - [x] 1.1 Update backend OpenAI client timeout to 180 seconds
    - Modify `defaultTimeout` constant in `backend/internal/openai/client.go`
    - _Requirements: 4.3_
  - [x] 1.2 Update frontend API timeout to 180 seconds
    - Modify `DEFAULT_TIMEOUT_MS` constant in `frontend/src/lib/api.ts`
    - _Requirements: 4.2_

- [x] 2. Create views tracking table and IP deduplication
  - [x] 2.1 Create database migration for views table
    - Add migration file with views table schema
    - Include unique constraint on (generation_id, ip_hash)
    - _Requirements: 5.1, 5.3_
  - [x] 2.2 Implement view tracking in storage repository
    - Add `RecordView(ctx, generationID, ipHash)` method
    - Return whether view was new or duplicate
    - Update view_count only for new views
    - _Requirements: 5.1, 5.3_
  - [x] 2.3 Write property test for idempotent view counting
    - **Property 4: Idempotent View Counting**
    - **Validates: Requirements 5.1, 5.3**
  - [x] 2.4 Update gallery handler to track views by IP
    - Hash IP address using SHA-256
    - Call RecordView on gallery item fetch
    - _Requirements: 5.1, 5.5_
  - [x] 2.5 Write property test for IP hashing
    - **Property 6: IP Addresses Are Hashed**
    - **Validates: Requirements 5.5**

- [x] 3. Fix vote deduplication by IP
  - [x] 3.1 Update rating storage to use IP hash
    - Ensure CreateOrUpdateRating uses IP hash for voter identification
    - Verify upsert behavior works correctly
    - _Requirements: 5.2, 5.4_
  - [x] 3.2 Write property test for vote upsert behavior
    - **Property 5: Vote Upsert Behavior**
    - **Validates: Requirements 5.2, 5.4**

- [x] 4. Checkpoint - Database and timeout changes
  - Ensure all tests pass, ask the user if questions arise.

- [x] 5. Differentiate experience-level prompts
  - [x] 5.1 Create beginner-specific question prompt
    - Add forbidden terms list
    - Use everyday language only
    - Include examples using real-world analogies
    - _Requirements: 1.1, 1.2_
  - [x] 5.2 Create novice-specific question prompt
    - Allow basic technical terms with explanations
    - Moderate complexity questions
    - _Requirements: 1.3_
  - [x] 5.3 Create expert-specific question prompt
    - Full technical terminology
    - Architecture and scaling questions
    - _Requirements: 1.4_
  - [x] 5.4 Write property test for beginner jargon avoidance
    - **Property 1: Beginner Questions Avoid Technical Jargon**
    - **Validates: Requirements 1.1, 1.2**
  - [x] 5.5 Write property test for experience level differentiation
    - **Property 2: Experience Levels Produce Different Questions**
    - **Validates: Requirements 1.5**

- [x] 6. Add example answers to questions
  - [x] 6.1 Update Question struct to include examples
    - Add `Examples []string` field to Question type
    - Update JSON serialization
    - _Requirements: 2.1_
  - [x] 6.2 Update question prompts to generate 3 examples per question
    - Modify system prompt to request examples
    - Examples should match experience level
    - _Requirements: 2.1, 2.5_
  - [x] 6.3 Update frontend Question type
    - Add `examples: string[]` to Question interface
    - _Requirements: 2.1_
  - [x] 6.4 Write property test for example count
    - **Property 3: Questions Include Exactly Three Examples**
    - **Validates: Requirements 2.1**

- [ ] 7. Checkpoint - Backend prompt changes complete
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 8. Implement frontend loading feedback
  - [ ] 8.1 Add loading messages to QuestionFlow component
    - Show "Generating questions... This may take up to 2 minutes"
    - Show "Still working..." after 30 seconds
    - _Requirements: 3.1, 3.4_
  - [ ] 8.2 Add loading messages to output generation
    - Show "Generating your files... This may take up to 3 minutes"
    - Show encouraging message after 30 seconds
    - _Requirements: 3.2, 3.4_
  - [ ] 8.3 Add animated spinner during generation
    - Use existing UI spinner component or add one
    - _Requirements: 3.3_

- [ ] 9. Implement clickable example answers UI
  - [ ] 9.1 Create ExampleAnswers component
    - Display 3 clickable example buttons per question
    - On click, populate answer field
    - _Requirements: 2.2_
  - [ ] 9.2 Integrate ExampleAnswers into QuestionFlow
    - Show examples below each question
    - Allow editing after selection
    - Preserve free-text input option
    - _Requirements: 2.2, 2.3, 2.4_

- [ ] 10. Improve navigation visibility
  - [ ] 10.1 Make "Browse Gallery" button more prominent on landing page
    - Increase size and contrast
    - Use primary button styling
    - _Requirements: 6.1_
  - [ ] 10.2 Add visible "Back to Home" link on gallery page
    - Add clear navigation link or logo that returns to landing
    - _Requirements: 6.2_
  - [ ] 10.3 Improve modal close button visibility
    - Increase close button size to 44x44px minimum
    - Add high contrast styling
    - _Requirements: 6.3, 6.4_
  - [ ] 10.4 Ensure modal closes on outside click and Escape key
    - Verify click-outside behavior works
    - Add Escape key handler if missing
    - _Requirements: 6.5, 6.6_

- [ ] 11. Final checkpoint
  - Ensure all tests pass, ask the user if questions arise.
  - Run frontend type check and lint
  - Run backend lint and vet

## Notes

- All tasks including property-based tests are required
- Timeout changes (Task 1) are quick wins that should be done first
- Database migration (Task 2) must be done before view tracking code
- Prompt changes (Tasks 5-6) are the core improvement for question quality
- Frontend changes (Tasks 8-10) can be done in parallel after backend is ready

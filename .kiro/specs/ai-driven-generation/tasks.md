# Implementation Plan: AI-Driven Generation

## Overview

This plan transforms BetterKiroPrompts from a multi-page template-based app to a single-page AI-driven experience. We'll first clean up old code, then build the new backend services, and finally implement the new frontend.

## Tasks

- [x] 1. Backend cleanup - Remove old endpoints and generators
  - [x] 1.1 Remove old API handlers (kickoff.go, steering.go, hooks.go)
    - Delete `internal/api/kickoff.go`
    - Delete `internal/api/steering.go`
    - Delete `internal/api/hooks.go`
    - Update router.go to remove old routes
    - _Requirements: 7.6_

  - [x] 1.2 Remove old generator package
    - Delete `internal/generator/kickoff.go`
    - Delete `internal/generator/steering.go`
    - Delete `internal/generator/hooks.go`
    - _Requirements: 7.8_

  - [x] 1.3 Remove template files
    - Delete `internal/templates/kickoff.tmpl`
    - Delete `internal/templates/steering/` directory
    - Delete `internal/templates/hooks/` directory
    - Update `internal/templates/templates.go` or remove if empty
    - _Requirements: 7.7_

- [x] 2. Backend - OpenAI client implementation
  - [x] 2.1 Create OpenAI client package
    - Create `internal/openai/client.go`
    - Implement ChatCompletion method with timeout
    - Load API key from environment variable
    - _Requirements: 6.1, 6.6_

  - [x] 2.2 Write property test for input validation
    - **Property 7: Input Validation**
    - **Validates: Requirements 6.5**

- [x] 3. Backend - Rate limiter implementation
  - [x] 3.1 Create rate limiter package
    - Create `internal/ratelimit/limiter.go`
    - Implement in-memory sliding window rate limiter
    - 10 requests per IP per hour
    - Return remaining time until reset
    - _Requirements: 6.3, 6.4_

  - [x] 3.2 Write property test for rate limiting
    - **Property 6: Rate Limiting Enforcement**
    - **Validates: Requirements 6.3**

- [x] 4. Backend - Generation service implementation
  - [x] 4.1 Create generation service
    - Create `internal/generation/service.go`
    - Implement GenerateQuestions method
    - Implement GenerateOutputs method
    - Include AI prompts from design
    - _Requirements: 2.1, 3.1_

  - [x] 4.2 Create API handlers for new endpoints
    - Create `internal/api/generate.go`
    - Implement POST /api/generate/questions handler
    - Implement POST /api/generate/outputs handler
    - Wire up rate limiter middleware
    - _Requirements: 2.1, 3.1, 6.3_

  - [x] 4.3 Write property test for question plan structure
    - **Property 1: Question Plan Structure**
    - **Validates: Requirements 2.2**

  - [x] 4.4 Write property test for generation response completeness
    - **Property 2: Generation Response Completeness**
    - **Validates: Requirements 3.2, 3.3, 3.4**

- [x] 5. Checkpoint - Backend complete
  - Ensure all backend tests pass
  - Verify endpoints work with curl/httpie
  - Ask the user if questions arise

- [ ] 6. Frontend cleanup - Remove old pages and components
  - [ ] 6.1 Remove old page components
    - Delete `src/pages/KickoffPage.tsx`
    - Delete `src/pages/SteeringPage.tsx`
    - Delete `src/pages/HooksPage.tsx`
    - _Requirements: 7.1_

  - [ ] 6.2 Remove old feature components
    - Delete `src/components/kickoff/` directory
    - Delete `src/components/steering/` directory
    - Delete `src/components/hooks/` directory
    - _Requirements: 7.2, 7.3, 7.4_

  - [ ] 6.3 Remove navigation component
    - Delete `src/components/shared/Navigation.tsx`
    - Delete `src/components/shared/StepIndicator.tsx`
    - _Requirements: 7.5, 1.4_

  - [ ] 6.4 Update API client
    - Remove old types (KickoffAnswers, SteeringConfig, HooksConfig)
    - Remove old functions (generateKickoff, generateSteering, generateHooks)
    - Add new types (Question, GeneratedFile, etc.)
    - Add new functions (generateQuestions, generateOutputs)
    - _Requirements: 7.9_

- [ ] 7. Frontend - Landing page implementation
  - [ ] 7.1 Create LandingPage component with state machine
    - Create `src/pages/LandingPage.tsx`
    - Implement phase state: input → questions → generating → output → error
    - Manage all state (projectIdea, questions, answers, files, editedFiles)
    - _Requirements: 1.1, 1.3_

  - [ ] 7.2 Create ProjectInput component
    - Create `src/components/ProjectInput.tsx`
    - Input field with placeholder "What project do you want to make?"
    - Display 3-5 example project ideas as clickable chips
    - Submit button with loading state
    - _Requirements: 1.1, 1.2_

  - [ ] 7.3 Create QuestionFlow component
    - Create `src/components/QuestionFlow.tsx`
    - Display current question with input
    - Show previous Q&A above current question
    - Allow clicking previous answers to edit
    - Next/Back navigation
    - _Requirements: 2.6, 2.7_

  - [ ] 7.4 Create OutputEditor component
    - Create `src/components/OutputEditor.tsx`
    - Tabbed interface for file types (kickoff, steering, hooks)
    - Syntax highlighting for markdown/JSON
    - Inline editing with textarea
    - Reset to original button per file
    - _Requirements: 4.1, 4.2, 4.4_

  - [ ] 7.5 Write property test for edit state preservation
    - **Property 3: Edit State Preservation**
    - **Validates: Requirements 4.3**

- [ ] 8. Frontend - Download functionality
  - [ ] 8.1 Implement download features
    - Individual file copy/download buttons
    - Download All as ZIP with correct directory structure
    - Use edited content if modified
    - _Requirements: 5.1, 5.2, 5.3, 5.4_

  - [ ] 8.2 Write property test for download content integrity
    - **Property 4: Download Content Integrity**
    - **Validates: Requirements 5.3**

  - [ ] 8.3 Write property test for ZIP structure
    - **Property 5: ZIP Directory Structure**
    - **Validates: Requirements 5.4**

- [ ] 9. Frontend - Error handling and loading states
  - [ ] 9.1 Implement loading and error states
    - Loading indicator during AI generation
    - Error message display (no retry button)
    - Timeout handling with refresh instruction
    - Rate limit display with countdown
    - _Requirements: 8.1, 8.2, 8.3, 8.4_

- [ ] 10. Integration - Wire everything together
  - [ ] 10.1 Update App.tsx
    - Remove Navigation import
    - Remove page state management
    - Render only LandingPage
    - _Requirements: 1.4_

  - [ ] 10.2 Update router.go with new routes
    - Add /api/generate/questions route
    - Add /api/generate/outputs route
    - Remove old routes
    - _Requirements: 2.1, 3.1_

- [ ] 11. Final checkpoint
  - Ensure all tests pass
  - Verify full flow works end-to-end
  - Check no old code remains
  - Ask the user if questions arise

## Notes

- All tasks including property-based tests are required
- Backend cleanup happens first to avoid import errors
- Frontend cleanup happens before new implementation
- Each task references specific requirements for traceability
- Property tests validate universal correctness properties

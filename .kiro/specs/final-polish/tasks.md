# Implementation Plan: Final Polish

## Overview

This implementation plan covers the final polish phase, organized into logical groups: database setup, backend services, frontend UI improvements, and integration. Tasks are ordered to build incrementally with early validation.

## Tasks

- [x] 1. Database Schema and Migrations
  - [x] 1.1 Create migration for categories table
    - Create `backend/migrations/20260113000001_create_categories.sql`
    - Define categories table with id, name, keywords array
    - Insert default categories (API, CLI, Web App, Mobile, Other)
    - _Requirements: 5.2, 5.3_
  - [x] 1.2 Create migration for generations table
    - Create `backend/migrations/20260113000002_create_generations.sql`
    - Define generations table with all required fields
    - Add foreign key to categories, indexes for sorting
    - _Requirements: 5.1, 5.2_
  - [x] 1.3 Create migration for ratings table
    - Create `backend/migrations/20260113000003_create_ratings.sql`
    - Define ratings table with unique constraint on (generation_id, voter_hash)
    - Add index on generation_id
    - _Requirements: 7.2_
  - [x] 1.4 Update database connection with pooling
    - Modify `backend/internal/db/db.go` to configure connection pool
    - Set max open connections, max idle, connection lifetime
    - _Requirements: 8.2_

- [x] 2. Backend Storage Layer
  - [x] 2.1 Create storage repository interface and implementation
    - Create `backend/internal/storage/repository.go` with Repository interface
    - Implement PostgreSQL repository with parameterized queries
    - _Requirements: 5.1, 5.4, 5.5_
  - [x] 2.2 Write property tests for storage repository
    - **Property 3: Generation Record Completeness**
    - **Validates: Requirements 5.2, 5.6**
  - [x] 2.3 Implement category matching logic
    - Create `backend/internal/storage/category.go`
    - Implement keyword-based category detection
    - _Requirements: 5.3_
  - [x] 2.4 Write property tests for category matching
    - **Property 4: Category Assignment Correctness**
    - **Validates: Requirements 5.3**

- [x] 3. Backend Input Sanitization
  - [x] 3.1 Create input sanitizer package
    - Create `backend/internal/sanitize/sanitize.go`
    - Implement HTML tag removal, special character escaping
    - Implement length validation
    - _Requirements: 9.2, 9.3_
  - [x] 3.2 Write property tests for sanitizer
    - **Property 13: Input Sanitization**
    - **Validates: Requirements 9.2, 9.3**

- [x] 4. Backend Request Infrastructure
  - [x] 4.1 Create request queue for OpenAI calls
    - Create `backend/internal/queue/queue.go`
    - Implement semaphore-based concurrency limiter
    - _Requirements: 8.1, 8.3_
  - [x] 4.2 Write property tests for request queue
    - **Property 11: Concurrent Request Handling**
    - **Property 12: Request Queue Fairness**
    - **Validates: Requirements 8.1, 8.3**
  - [x] 4.3 Create middleware for request ID and logging
    - Create `backend/internal/api/middleware.go`
    - Implement RequestIDMiddleware, LoggingMiddleware, RecoveryMiddleware
    - _Requirements: 9.5, 9.6_
  - [x] 4.4 Update router to use middleware chain
    - Modify `backend/internal/api/router.go`
    - Apply middleware to all routes
    - _Requirements: 9.5_
  - [x] 4.5 Implement structured error responses
    - Create `backend/internal/api/errors.go`
    - Define error codes and response format
    - Update handlers to use structured errors
    - _Requirements: 10.5, 10.6_
  - [x] 4.6 Write property tests for error responses
    - **Property 14: Structured Error Responses**
    - **Validates: Requirements 10.5, 10.6**

- [x] 5. Backend Timeout and OpenAI Updates
  - [x] 5.1 Increase OpenAI client timeout to 120 seconds
    - Modify `backend/internal/openai/client.go`
    - Update defaultTimeout constant
    - _Requirements: 3.1_
  - [x] 5.2 Integrate request queue with generation service
    - Modify `backend/internal/generation/service.go`
    - Acquire queue slot before OpenAI calls
    - _Requirements: 8.3_
  - [x] 5.3 Update generation service to store results
    - Modify `backend/internal/generation/service.go`
    - Save successful generations to database
    - Return generation ID in response
    - _Requirements: 5.1_

- [ ] 6. Checkpoint - Backend Core Complete
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 7. Backend Gallery Service
  - [ ] 7.1 Create gallery service
    - Create `backend/internal/gallery/service.go`
    - Implement ListGenerations, GetGeneration, RateGeneration
    - _Requirements: 6.1, 6.2, 6.3, 7.2_
  - [ ] 7.2 Write property tests for gallery filtering
    - **Property 5: Gallery Filtering Correctness**
    - **Validates: Requirements 6.2**
  - [ ] 7.3 Write property tests for gallery sorting
    - **Property 6: Gallery Sorting Correctness**
    - **Validates: Requirements 6.3**
  - [ ] 7.4 Write property tests for pagination
    - **Property 7: Pagination Bounds**
    - **Validates: Requirements 6.5**

- [ ] 8. Backend Rating System
  - [ ] 8.1 Implement rating endpoints
    - Add POST /api/gallery/:id/rate endpoint
    - Add GET /api/gallery/:id endpoint with user rating
    - _Requirements: 7.1, 7.2_
  - [ ] 8.2 Implement rating rate limiter
    - Create separate rate limiter for ratings (20/hour)
    - _Requirements: 7.6_
  - [ ] 8.3 Write property tests for rating calculation
    - **Property 8: Rating Storage and Calculation**
    - **Validates: Requirements 7.2, 7.3**
  - [ ] 8.4 Write property tests for duplicate prevention
    - **Property 9: Duplicate Rating Prevention**
    - **Validates: Requirements 7.4**
  - [ ] 8.5 Write property tests for rating rate limit
    - **Property 10: Rate Limit Enforcement**
    - **Validates: Requirements 7.6**

- [ ] 9. Backend Gallery API Routes
  - [ ] 9.1 Create gallery API handlers
    - Create `backend/internal/api/gallery.go`
    - Implement handlers for list, detail, rate endpoints
    - _Requirements: 6.1, 6.4, 7.2_
  - [ ] 9.2 Register gallery routes in router
    - Add GET /api/gallery, GET /api/gallery/:id, POST /api/gallery/:id/rate
    - _Requirements: 6.1_

- [ ] 10. Checkpoint - Backend Complete
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 11. Frontend LocalStorage Manager
  - [ ] 11.1 Create storage manager utility
    - Create `frontend/src/lib/storage.ts`
    - Implement save, load, clear, isExpired functions
    - _Requirements: 4.1, 4.5_
  - [ ] 11.2 Write property tests for storage manager
    - **Property 1: State Persistence Consistency**
    - **Validates: Requirements 4.1, 4.3**
  - [ ] 11.3 Integrate storage manager with LandingPage
    - Modify `frontend/src/pages/LandingPage.tsx`
    - Save state after each action, restore on load
    - Clear on successful generation
    - _Requirements: 4.1, 4.2, 4.3, 4.4_
  - [ ] 11.4 Add restore prompt dialog
    - Create restore confirmation UI
    - Handle localStorage unavailable gracefully
    - _Requirements: 4.2, 4.7_

- [ ] 12. Frontend Syntax Highlighting
  - [ ] 12.1 Install react-syntax-highlighter
    - Add dependency to package.json
    - _Requirements: 2.1_
  - [ ] 12.2 Create SyntaxHighlighter component
    - Create `frontend/src/components/SyntaxHighlighter.tsx`
    - Support JSON, Markdown, YAML languages
    - Implement dark theme matching app colors
    - _Requirements: 2.1, 2.2, 2.3_
  - [ ] 12.3 Write property tests for syntax highlighter
    - **Property 2: Syntax Highlighter Robustness**
    - **Validates: Requirements 2.6**
  - [ ] 12.4 Update OutputEditor to use syntax highlighting
    - Modify `frontend/src/components/OutputEditor.tsx`
    - Add toggle between highlighted view and edit mode
    - _Requirements: 2.4, 2.5_

- [ ] 13. Frontend UI Flow Improvements
  - [ ] 13.1 Update App.tsx for logo visibility
    - Add phase-based logo visibility logic
    - Implement fade-out animation
    - _Requirements: 1.1_
  - [ ] 13.2 Create CompactHeader component
    - Create `frontend/src/components/shared/CompactHeader.tsx`
    - Small logo with link to start, Start Over button
    - _Requirements: 1.2, 4.6_
  - [ ] 13.3 Add phase transition animations
    - Add CSS transitions for content changes
    - _Requirements: 1.3_
  - [ ] 13.4 Add success celebration animation
    - Create celebration component (confetti or similar)
    - Show briefly before displaying results
    - _Requirements: 1.5_

- [ ] 14. Frontend Error Recovery
  - [ ] 14.1 Update error handling with retry support
    - Modify `frontend/src/pages/LandingPage.tsx`
    - Add canRetry state, implement retry logic
    - _Requirements: 4.3, 10.2_
  - [ ] 14.2 Implement automatic retry for recoverable errors
    - Add retry logic for timeouts and 503 errors
    - _Requirements: 10.4_
  - [ ] 14.3 Write property tests for retry behavior
    - **Property 15: Automatic Retry Behavior**
    - **Validates: Requirements 10.4**
  - [ ] 14.4 Update ErrorMessage component
    - Add Try Again and Start Over buttons
    - Show appropriate actions based on error type
    - _Requirements: 10.1, 10.2, 10.3_

- [ ] 15. Frontend Timeout Handling
  - [ ] 15.1 Add AbortController to API calls
    - Modify `frontend/src/lib/api.ts`
    - Implement timeout with AbortController
    - _Requirements: 3.4_
  - [ ] 15.2 Update LoadingState with progress tracking
    - Track loading start time
    - Show "taking longer than usual" after 90 seconds
    - _Requirements: 3.2, 3.3_

- [ ] 16. Checkpoint - Frontend Core Complete
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 17. Frontend Gallery
  - [ ] 17.1 Create gallery API functions
    - Add to `frontend/src/lib/api.ts`
    - Implement listGallery, getGalleryItem, rateGalleryItem
    - _Requirements: 6.1, 6.4, 7.2_
  - [ ] 17.2 Create GalleryList component
    - Create `frontend/src/components/Gallery/GalleryList.tsx`
    - Implement paginated list with filters
    - _Requirements: 6.1, 6.2, 6.3_
  - [ ] 17.3 Create GalleryDetail component
    - Create `frontend/src/components/Gallery/GalleryDetail.tsx`
    - Show full generation with copy/download
    - _Requirements: 6.4, 6.6_
  - [ ] 17.4 Create Rating component
    - Create `frontend/src/components/Gallery/Rating.tsx`
    - Implement 1-5 star interface
    - _Requirements: 7.1, 7.5_
  - [ ] 17.5 Create GalleryPage
    - Create `frontend/src/pages/GalleryPage.tsx`
    - Integrate list, detail modal, and rating
    - _Requirements: 6.1_
  - [ ] 17.6 Add gallery navigation
    - Add Gallery link to header/navigation
    - Update App.tsx with routing if needed
    - _Requirements: 6.1_

- [ ] 18. Frontend Voter Hash
  - [ ] 18.1 Implement voter hash generation
    - Create `frontend/src/lib/voter.ts`
    - Generate consistent hash from localStorage + fingerprint
    - _Requirements: 7.4_

- [ ] 19. Final Integration
  - [ ] 19.1 Update API response types
    - Add generationId to GenerateOutputsResponse
    - Add gallery types
    - _Requirements: 5.1_
  - [ ] 19.2 Connect generation success to gallery
    - Show "View in Gallery" link after successful generation
    - _Requirements: 6.4_

- [ ] 20. Graceful Shutdown
  - [ ] 20.1 Implement graceful shutdown in main.go
    - Handle SIGTERM/SIGINT signals
    - Wait for in-flight requests to complete
    - _Requirements: 8.6_

- [ ] 21. Final Checkpoint
  - Run full test suite
  - Run `pnpm typecheck && pnpm lint` in frontend
  - Run `golangci-lint run && go vet` in backend
  - Build production artifacts
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- All property-based tests are required for comprehensive coverage
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases

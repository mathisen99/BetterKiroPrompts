# Implementation Plan: Comprehensive Logging

## Overview

This plan implements comprehensive file-based logging across the entire application stack. We'll create a centralized logger package, integrate it into all backend components, add frontend error collection, and configure Docker for log file mounting.

## Tasks

- [x] 1. Create Logger Package Infrastructure
  - [x] 1.1 Create `backend/internal/logger/logger.go` with core Logger struct
    - Config struct with Level, LogDir, MaxSizeMB, MaxAgeDays, EnableColor
    - New() constructor that creates log directory and initializes handlers
    - Category-specific loggers: App(), HTTP(), DB(), Scanner(), Client()
    - SetLevel() for runtime level changes
    - Close() for cleanup
    - _Requirements: 1.1, 1.2, 1.3, 9.1, 9.2_

  - [x] 1.2 Create `backend/internal/logger/rotating.go` with RotatingFile
    - Write() that tracks size and triggers rotation at MaxSizeMB
    - Rotate() that renames current file with timestamp suffix
    - Cleanup() that removes files older than MaxAgeDays
    - _Requirements: 1.6, 1.7_

  - [x] 1.3 Create `backend/internal/logger/color.go` with ColorHandler
    - Custom slog.Handler that wraps output with ANSI colors
    - Level-to-color mapping: ERROR→RED, WARN→YELLOW, INFO→GREEN, DEBUG→CYAN
    - Request ID highlighting in MAGENTA, timestamps in GRAY
    - TTY detection to auto-disable colors for non-terminals
    - NO_COLOR environment variable support
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5, 10.6, 10.7, 10.8, 10.9_

  - [x] 1.4 Create `backend/internal/logger/context.go` with context helpers
    - Context keys: RequestIDKey, ComponentKey, UserIPKey
    - WithRequestID(), WithComponent(), GetRequestID() functions
    - _Requirements: 2.2, 2.4_

  - [x] 1.5 Create `backend/internal/logger/redact.go` with sensitive data redaction
    - RedactSensitive() that replaces password, token, api_key, secret, authorization values
    - Apply redaction before any log write
    - _Requirements: 2.5, 3.5_

- [x] 2. Integrate Logger into Application Startup
  - [x] 2.1 Update `backend/cmd/server/main.go` to initialize logger
    - Parse LOG_LEVEL from environment (default INFO)
    - Create logger with ./logs directory
    - Pass logger to all services and handlers
    - Add startup/shutdown logging
    - _Requirements: 9.1, 9.2_

  - [x] 2.2 Update `docker-compose.yml` to mount logs directory
    - Add volume mount: `./logs:/app/logs`
    - Ensure directory exists on host
    - _Requirements: 8.3_

  - [x] 2.3 Update `.gitignore` to exclude logs directory
    - Add `logs/` entry
    - _Requirements: 8.4_

  - [x] 2.4 Create `logs/README.md` explaining log structure
    - Document file naming convention
    - Document retention policy
    - Document log levels
    - _Requirements: 8.5_

- [x] 3. Update HTTP Middleware with Logging
  - [x] 3.1 Update `backend/internal/api/middleware.go` LoggingMiddleware
    - Accept logger as parameter
    - Log request_start with method, path, query, remote_addr, user_agent, content_length
    - Log request_complete with status, duration, bytes_written
    - Log security events (rate limits, validation failures)
    - _Requirements: 2.1, 2.3_

  - [x] 3.2 Update `backend/internal/api/middleware.go` RequestIDMiddleware
    - Use logger context helpers
    - Ensure request ID propagates to all downstream operations
    - _Requirements: 2.2, 2.4_

  - [x] 3.3 Update `backend/internal/api/router.go` to pass logger to middleware
    - Update NewRouter to accept logger
    - Pass logger to LoggingMiddleware
    - _Requirements: 2.1_

- [ ] 4. Add Logging to Generation Service
  - [ ] 4.1 Update `backend/internal/generation/service.go` with logger
    - Add logger field to Service struct
    - Update constructors to accept logger
    - _Requirements: 4.1_

  - [ ] 4.2 Add logging to GenerateQuestions
    - Log start with experience_level, idea_length
    - Log validation failures
    - Log queue acquisition
    - Log OpenAI call start/complete
    - Log parse failures
    - Log completion with question_count, duration
    - _Requirements: 4.1, 4.2, 4.3_

  - [ ] 4.3 Add logging to GenerateOutputs
    - Log start with experience_level, hook_preset, answer_count
    - Log each retry attempt
    - Log parse/validation failures
    - Log completion with file_count, attempts_used, duration
    - _Requirements: 4.1, 4.2, 4.3_

  - [ ] 4.4 Add logging to GenerateAndStoreOutputs
    - Log storage attempt
    - Log category lookup
    - Log storage success/failure
    - _Requirements: 4.5_

- [ ] 5. Add Logging to Gallery Service
  - [ ] 5.1 Update `backend/internal/gallery/service.go` with logger
    - Add logger field to Service struct
    - Update constructor to accept logger
    - _Requirements: 4.1_

  - [ ] 5.2 Add logging to ListGenerations
    - Log start with sort_by, page, page_size, category_id
    - Log completion with item_count, total, duration
    - _Requirements: 4.5_

  - [ ] 5.3 Add logging to GetGenerationWithView
    - Log start with generation_id
    - Log view recording (new_view boolean)
    - Log completion
    - _Requirements: 4.5_

  - [ ] 5.4 Add logging to RateGeneration
    - Log start with generation_id, score
    - Log rate limit checks
    - Log completion or failure
    - _Requirements: 4.5_

- [ ] 6. Add Logging to Scanner Service
  - [ ] 6.1 Update `backend/internal/scanner/service.go` with logger
    - Add logger field to Service struct
    - Update constructor to accept logger
    - _Requirements: 6.1_

  - [ ] 6.2 Add logging to runScan pipeline
    - Log scan_pipeline_start with job_id
    - Log each phase start/complete with timing
    - Log scan_pipeline_complete with total_findings, total_duration
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

  - [ ] 6.3 Add logging to clone phase
    - Log repo_url at start
    - Log path, duration at complete
    - Log errors with details
    - _Requirements: 6.2_

  - [ ] 6.4 Add logging to language detection phase
    - Log detected languages and counts
    - _Requirements: 6.3_

  - [ ] 6.5 Add logging to tool execution phase
    - Log each tool start/complete
    - Log finding_count, timed_out, duration per tool
    - Log tool errors
    - _Requirements: 6.4_

  - [ ] 6.6 Add logging to aggregation phase
    - Log result_count at start
    - Log total_findings with severity breakdown at complete
    - _Requirements: 6.5_

  - [ ] 6.7 Add logging to AI review phase
    - Log findings_to_review at start
    - Log reviewed_findings, duration at complete
    - Log skip reason if skipped
    - _Requirements: 6.4_

- [ ] 7. Add Logging to OpenAI Client
  - [ ] 7.1 Update `backend/internal/openai/client.go` with logger
    - Add logger field to Client struct
    - Update constructors to accept logger
    - _Requirements: 5.1_

  - [ ] 7.2 Add logging to ChatCompletion
    - Log request_start with model, prompt_length, message_count, reasoning_effort
    - Log truncated prompt preview at DEBUG level (first 500 chars)
    - Log response_received with status_code, response_length, latency
    - Log truncated response preview at DEBUG level
    - Log errors with details
    - _Requirements: 5.1, 5.2, 5.4, 5.5_

- [ ] 8. Add Logging to Database Operations
  - [ ] 8.1 Create `backend/internal/db/logging.go` with LoggingDB wrapper
    - Wrap sql.DB with logging
    - QueryContext logging with type, duration, success
    - ExecContext logging with type, duration, rows_affected
    - _Requirements: 3.1, 3.2_

  - [ ] 8.2 Update `backend/internal/db/db.go` with connection logging
    - Log connection attempts and retries
    - Log successful connection with pool config
    - Log migration execution
    - Log connection close
    - _Requirements: 3.3, 3.4_

  - [ ] 8.3 Update repository to use LoggingDB
    - Update storage.NewPostgresRepository to accept LoggingDB
    - Ensure all queries go through logging wrapper
    - _Requirements: 3.1, 3.2_

- [ ] 9. Add Logging to Queue and Rate Limiter
  - [ ] 9.1 Update `backend/internal/queue/queue.go` with logger
    - Add logger field
    - Log acquire start/success/timeout
    - Log release with stats
    - _Requirements: 4.1_

  - [ ] 9.2 Update `backend/internal/ratelimit/limiter.go` with logger
    - Add logger field
    - Log allow/deny decisions with IP hash
    - Log remaining count at DEBUG level
    - _Requirements: 4.1_

- [ ] 10. Add Client Logging Endpoint
  - [ ] 10.1 Create `backend/internal/api/logs.go` with HandleClientLogs
    - Parse ClientLogRequest with array of log entries
    - Write each entry to client log file
    - Return 202 Accepted
    - _Requirements: 7.5_

  - [ ] 10.2 Register `/api/logs/client` endpoint in router
    - Add POST handler
    - No rate limiting (logs are important)
    - _Requirements: 7.5_

- [ ] 11. Add Frontend Log Collector
  - [ ] 11.1 Create `frontend/src/lib/logger.ts` with LogCollector class
    - Buffer for log entries
    - setupErrorHandlers for window.onerror and unhandledrejection
    - log(), debug(), info(), warn(), error() methods
    - logApiCall() for API timing
    - flush() to send batched logs to backend
    - Colored console output
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

  - [ ] 11.2 Update `frontend/src/lib/api.ts` to use logger
    - Import logger
    - Add logApiCall after each fetch
    - Log errors on failures
    - _Requirements: 7.3_

  - [ ] 11.3 Update `frontend/src/components/shared/ErrorBoundary.tsx` to use logger
    - Import logger
    - Log React errors in componentDidCatch
    - _Requirements: 7.2_

  - [ ] 11.4 Add logging to key frontend components
    - QuestionFlow: log phase transitions
    - OutputEditor: log file operations
    - ScanProgress: log scan status changes
    - _Requirements: 7.1_

- [ ] 12. Add Log Level Admin Endpoint
  - [ ] 12.1 Create `/api/admin/log-level` endpoint
    - GET to retrieve current level
    - POST to change level at runtime
    - _Requirements: 9.3_

- [ ] 13. Update Environment Configuration
  - [ ] 13.1 Update `.env.example` with logging variables
    - LOG_LEVEL (DEBUG, INFO, WARN, ERROR)
    - NO_COLOR (optional)
    - _Requirements: 9.1, 10.8_

- [ ] 14. Checkpoint - Verify Logging Works
  - Run the application and verify:
    - Log files are created in ./logs directory
    - Console output has colors
    - Request IDs propagate through operations
    - All services log their operations
  - Ensure all tests pass, ask the user if questions arise.

- [ ]* 15. Write Property Tests for Logger
  - [ ]* 15.1 Write property test for JSON structure
    - **Property 1: Log Entry JSON Structure**
    - **Validates: Requirements 1.4, 1.5**

  - [ ]* 15.2 Write property test for request ID uniqueness
    - **Property 3: Request ID Uniqueness**
    - **Validates: Requirements 2.2**

  - [ ]* 15.3 Write property test for sensitive data redaction
    - **Property 5: Sensitive Data Redaction**
    - **Validates: Requirements 2.5, 3.5**

  - [ ]* 15.4 Write property test for level filtering
    - **Property 14: Level Filtering**
    - **Validates: Requirements 9.4, 9.5**

  - [ ]* 15.5 Write property test for color mapping
    - **Property 13: Color Mapping by Level**
    - **Validates: Requirements 10.1, 10.2, 10.3, 10.4, 10.5**

- [ ] 16. Final Checkpoint
  - Run full test suite
  - Verify log output in all scenarios
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases

# Implementation Plan: Info Page and Security Scanning

## Overview

This implementation plan covers the Info Page and Security Scanning features for the BetterKiroPrompts hackathon submission. The implementation is organized to build incrementally, starting with the simpler Info Page, then the Security Container, backend scanner service, and finally the frontend Security Scan Page.

## Tasks

- [x] 1. Create Info Page
  - [x] 1.1 Create InfoPage component with content sections
    - Create `frontend/src/pages/InfoPage.tsx`
    - Add hero section explaining site purpose
    - Add problem statement section (vibe-coding without thinking)
    - Add feature cards for kickoff prompts, steering, hooks, security scanning
    - Add self-hosting explanation
    - Add "Get Started" CTA button
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7_

  - [x] 1.2 Add navigation to Info Page
    - Update `App.tsx` to add 'info' view state
    - Add "About" link to landing page header area
    - Add "About" link to gallery page
    - Add navigation back to home and gallery from Info Page
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

- [x] 2. Checkpoint - Info Page Complete
  - Ensure Info Page renders correctly
  - Verify navigation works between all pages
  - Run `pnpm typecheck && pnpm lint` in frontend

- [x] 3. Create Security Container
  - [x] 3.1 Create Dockerfile.scanner
    - Create `backend/Dockerfile.scanner`
    - Multi-stage build with Go, Python, Rust builders
    - Install universal tools: Trivy, Semgrep, TruffleHog, Gitleaks
    - Install Go tools: govulncheck
    - Install Python tools: bandit, pip-audit, safety
    - Install Node.js and npm for npm audit
    - Install Rust tools: cargo audit
    - Install Ruby tools: bundler-audit, brakeman
    - Configure non-root user
    - _Requirements: 12.1, 12.2, 12.3, 12.4, 12.5, 12.6, 12.7_

  - [x] 3.2 Update docker-compose.yml
    - Add scanner service with profile "scan"
    - Configure shared volume for repos
    - Set resource limits (CPU, memory)
    - Configure network access
    - _Requirements: 12.8, 12.9, 12.10_

  - [x] 3.3 Update .env.example
    - Add GITHUB_TOKEN documentation
    - Add scanner-related configuration variables
    - _Requirements: 13.4_

- [x] 4. Checkpoint - Security Container Ready
  - [x] Build scanner container: `docker compose --profile scan build scanner`
  - [x] Verify all tools are installed in container
  - [x] Test container starts correctly

- [x] 5. Create Database Schema for Scans
  - [x] 5.1 Create migration for scan tables
    - Create `backend/migrations/005_create_scan_tables.sql`
    - Create scan_jobs table with status, languages, error fields
    - Create scan_findings table with severity, tool, file_path, remediation
    - Add indexes for efficient queries
    - _Requirements: 11.1, 11.3_

- [x] 6. Implement Scanner Backend Service
  - [x] 6.1 Create URL validator
    - Create `backend/internal/scanner/validator.go`
    - Validate GitHub URL format
    - Support both https://github.com/owner/repo and .git suffix
    - Return structured validation errors
    - _Requirements: 4.1, 4.4, 4.6_

  - [x] 6.2 Write property test for URL validation
    - **Property 1: URL Validation**
    - **Validates: Requirements 4.1, 4.4, 4.6**

  - [x] 6.3 Create repository cloner
    - Create `backend/internal/scanner/cloner.go`
    - Clone to temporary directory
    - Support GitHub token for private repos
    - Implement size limit check
    - Implement cleanup function
    - Never log token values
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

  - [x] 6.4 Write property test for token security
    - **Property 3: Token Security**
    - **Validates: Requirements 5.3**

  - [x] 6.5 Create language detector
    - Create `backend/internal/scanner/language.go`
    - Detect languages by file extension
    - Support: Go, JS, TS, Python, Java, Ruby, PHP, C, C++, Rust
    - Rank by file count
    - _Requirements: 6.1, 6.2, 6.3_

  - [x] 6.6 Write property test for language detection
    - **Property 6: Language Detection Accuracy**
    - **Validates: Requirements 6.1, 6.2, 6.3**

  - [x] 6.7 Create tool runner
    - Create `backend/internal/scanner/tools.go`
    - Implement RunTrivy, RunSemgrep, RunTruffleHog, RunGitleaks
    - Implement RunGovulncheck (Go)
    - Implement RunBandit, RunPipAudit, RunSafety (Python)
    - Implement RunNpmAudit (JavaScript/TypeScript)
    - Implement RunCargoAudit (Rust)
    - Implement RunBundlerAudit, RunBrakeman (Ruby)
    - Implement timeout enforcement
    - Implement GetToolsForLanguages selector
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 7.8, 7.9, 7.10, 7.11_

  - [x] 6.8 Write property test for tool timeout
    - **Property 7: Tool Timeout Enforcement**
    - **Validates: Requirements 7.10, 7.11**

  - [x] 6.9 Create finding aggregator
    - Create `backend/internal/scanner/aggregator.go`
    - Parse tool outputs into unified Finding format
    - Deduplicate findings (same file, line, description)
    - Rank by severity (critical > high > medium > low > info)
    - _Requirements: 8.1, 8.2, 8.3, 8.4_

  - [x] 6.10 Write property test for finding aggregation
    - **Property 8: Finding Aggregation Completeness**
    - **Validates: Requirements 8.1, 8.2, 8.3, 8.4**

  - [x] 6.11 Create code reviewer
    - Create `backend/internal/scanner/reviewer.go`
    - Build GPT-5.1-Codex-Max request with findings and file contents
    - Enforce max files limit (default 10)
    - Parse JSON response for remediation and code examples
    - Handle API errors gracefully
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5, 9.6, 9.7_

  - [x] 6.12 Write property test for AI review scope
    - **Property 9: AI Review Scope Limitation**
    - **Validates: Requirements 9.2, 9.3, 9.7**

  - [x] 6.13 Create scanner service
    - Create `backend/internal/scanner/service.go`
    - Orchestrate: validate → clone → detect → scan → aggregate → review → cleanup
    - Manage scan job status transitions
    - Persist results to database
    - _Requirements: 4.5, 11.1, 11.2, 11.4_

  - [x] 6.14 Write property test for job creation round-trip
    - **Property 2: Job Creation Round-Trip**
    - **Validates: Requirements 4.5, 11.1, 11.2**

- [ ] 7. Checkpoint - Scanner Service Complete
  - Run `golangci-lint run && go fmt ./... && go vet ./...` in backend
  - Run all property tests
  - Ensure all tests pass

- [ ] 8. Create Scanner API Endpoints
  - [ ] 8.1 Create scan handler
    - Create `backend/internal/api/scan.go`
    - POST /api/scan - Start scan job
    - GET /api/scan/{id} - Get scan status/results
    - GET /api/scan/config - Get configuration (private repo enabled)
    - _Requirements: 4.1, 4.2, 4.3, 4.5, 11.2, 13.1, 13.2, 13.3_

  - [ ] 8.2 Register scan routes
    - Update `backend/internal/api/router.go`
    - Add scan endpoints with rate limiting
    - _Requirements: 4.1_

- [ ] 9. Create Security Scan Page
  - [ ] 9.1 Create SecurityScanPage component
    - Create `frontend/src/pages/SecurityScanPage.tsx`
    - Add repo URL input field
    - Add scan explanation text
    - Show private repo availability indicator
    - Add "Start Scan" button
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

  - [ ] 9.2 Create scan API client functions
    - Update `frontend/src/lib/api.ts`
    - Add startScan, getScanStatus, getScanConfig functions
    - Add ScanJob, Finding, ScanConfig types
    - _Requirements: 4.5, 11.2_

  - [ ] 9.3 Create scan results display
    - Create `frontend/src/components/ScanResults.tsx`
    - Group findings by severity
    - Show file path, line number, description, tool
    - Display remediation with syntax highlighting
    - _Requirements: 10.1, 10.2, 10.3, 10.4_

  - [ ] 9.4 Write property test for finding display
    - **Property 11: Finding Display Completeness**
    - **Validates: Requirements 10.2, 10.3, 10.4**

  - [ ] 9.5 Create scan progress indicator
    - Create `frontend/src/components/ScanProgress.tsx`
    - Show current status (cloning, scanning, reviewing)
    - Poll for status updates
    - _Requirements: 10.5_

  - [ ] 9.6 Add navigation to Security Scan Page
    - Update `App.tsx` to add 'scan' view state
    - Add "Security Scan" link to Info Page
    - Add "Security Scan" link to navigation
    - _Requirements: 3.1_

- [ ] 10. Checkpoint - Security Scan Page Complete
  - Run `pnpm typecheck && pnpm lint && pnpm build` in frontend
  - Verify scan flow works end-to-end
  - Test with public repository

- [ ] 11. Final Integration and Testing
  - [ ] 11.1 End-to-end integration test
    - Test complete scan flow with real repository
    - Verify findings display correctly
    - Verify AI remediation appears when findings exist
    - _Requirements: All_

  - [ ] 11.2 Error handling verification
    - Test invalid URL handling
    - Test private repo without token
    - Test timeout handling
    - Test API error handling
    - _Requirements: 4.3, 4.4, 7.11, 11.4_

- [ ] 12. Final Checkpoint
  - Ensure all tests pass
  - Run full build: `./build.sh --prod --build up`
  - Verify production build works
  - Ask the user if questions arise

## Notes

- All tasks including property-based tests are required for comprehensive coverage
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases

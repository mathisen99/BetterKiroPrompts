# Phase 3: Polish & Testing — Tasks

## Task List

### UI/UX Polish

- [x] 1. Add ErrorBoundary component at app level
  - Refs: AC-1
  - Outcome: Graceful error display with retry

- [x] 2. Add error handling to all API calls with user-friendly messages
  - Refs: AC-1
  - Outcome: Errors show toast with retry option

- [x] 3. Add loading skeletons to wizard steps
  - Refs: AC-2
  - Outcome: Skeleton UI during data fetch

- [x] 4. Add loading states to all buttons during submission
  - Refs: AC-2
  - Outcome: Spinner and disabled state on submit

- [x] 5. Integrate shadcn/ui toast for success notifications
  - Refs: AC-3
  - Outcome: Toast on copy/download success

### Accessibility

- [x] 6. Add labels to all form inputs
  - Refs: AC-5
  - Outcome: Screen readers announce field names

- [x] 7. Implement focus management in wizard
  - Refs: AC-4
  - Outcome: Auto-focus on step change

- [x] 8. Add ARIA live regions for dynamic content
  - Refs: AC-5
  - Outcome: State changes announced

- [x] 9. Verify color contrast meets WCAG 2.1 AA
  - Refs: AC-6
  - Outcome: axe-core audit passes

- [x] 10. Test complete keyboard navigation flow
  - Refs: AC-4
  - Outcome: All flows completable without mouse

### Missing Features

- [x] 11. Add manual steering option to SteeringConfigurator
  - Refs: AC-7
  - Outcome: Checkbox generates `inclusion: manual` files

- [x] 12. Update steering templates for manual inclusion
  - Refs: AC-7
  - Outcome: Backend generates correct frontmatter

- [x] 13. Add file reference input to steering UI
  - Refs: AC-8
  - Outcome: Users can add `#[[file:<path>]]` references

- [x] 14. Update steering templates with file reference syntax
  - Refs: AC-8
  - Outcome: Generated files include references

- [x] 15. Create CommitContract component
  - Refs: AC-9
  - Outcome: Displays atomic/prefixed/one-sentence rules

- [x] 16. Add CommitContract to all output panels
  - Refs: AC-9
  - Outcome: Contract visible on every generation

### Unit Tests

- [x] 17. Create `backend/internal/generator/kickoff_test.go`
  - Refs: AC-10
  - Outcome: Tests for prompt generation logic

- [x] 18. Create `backend/internal/generator/steering_test.go`
  - Refs: AC-10
  - Outcome: Tests for steering file generation

- [x] 19. Create `backend/internal/generator/hooks_test.go`
  - Refs: AC-10
  - Outcome: Tests for hook generation and presets

### Integration Tests

- [x] 20. Create `backend/internal/api/kickoff_test.go`
  - Refs: AC-11
  - Outcome: Tests for /api/kickoff/generate endpoint

- [ ] 21. Create `backend/internal/api/steering_test.go`
  - Refs: AC-11
  - Outcome: Tests for /api/steering/generate endpoint

- [ ] 22. Create `backend/internal/api/hooks_test.go`
  - Refs: AC-11
  - Outcome: Tests for /api/hooks/generate endpoint

### E2E Tests

- [ ] 23. Set up Playwright in frontend
  - Refs: AC-12
  - Outcome: E2E test infrastructure ready

- [ ] 24. Create kickoff wizard E2E test
  - Refs: AC-12
  - Outcome: Full wizard flow tested

- [ ] 25. Create steering generation E2E test
  - Refs: AC-12
  - Outcome: Config to download tested

- [ ] 26. Create hooks generation E2E test
  - Refs: AC-12
  - Outcome: Preset selection to download tested

### Documentation

- [ ] 27. Write README.md with setup and quick start
  - Refs: AC-13
  - Outcome: New developers can start in <5 minutes

- [ ] 28. Create docs/api.md with endpoint documentation
  - Refs: AC-14
  - Outcome: All endpoints documented with examples

- [ ] 29. Create docs/user-guide.md
  - Refs: AC-15
  - Outcome: Users understand generated outputs

### Definition of Done Verification

- [ ] 30. Verify: User can generate full kickoff prompt
  - Refs: AC-16
  - Outcome: Manual test passes

- [ ] 31. Verify: Steering files are usable and correctly scoped
  - Refs: AC-17
  - Outcome: Files work in Kiro IDE

- [ ] 32. Verify: Hooks are valid and usable
  - Refs: AC-18
  - Outcome: Hooks recognized by Kiro IDE

- [ ] 33. Final accessibility audit
  - Refs: AC-4, AC-5, AC-6
  - Outcome: No critical accessibility issues

### OPTIONAL: Repo Scanning

- [ ] 34. (OPTIONAL) Create scan worker Dockerfile
  - Refs: AC-20, AC-23
  - Outcome: Isolated container with no outbound network

- [ ] 35. (OPTIONAL) Implement TruffleHog integration
  - Refs: AC-21
  - Outcome: Secret scanning works

- [ ] 36. (OPTIONAL) Implement Gitleaks integration
  - Refs: AC-21
  - Outcome: Secret scanning works

- [ ] 37. (OPTIONAL) Implement osv-scanner integration
  - Refs: AC-21
  - Outcome: Dependency scanning works

- [ ] 38. (OPTIONAL) Implement govulncheck integration
  - Refs: AC-21
  - Outcome: Go vulnerability scanning works

- [ ] 39. (OPTIONAL) Implement hard timeout handling
  - Refs: AC-22
  - Outcome: Scans terminate and return partial results

- [ ] 40. (OPTIONAL) Create scan results API endpoints
  - Refs: AC-24
  - Outcome: Structured findings returned

- [ ] 41. (OPTIONAL) Implement AI summary of scan results
  - Refs: AC-25
  - Outcome: Prioritized summary without invented findings

- [ ] 42. (OPTIONAL) Create scan UI page
  - Refs: AC-19
  - Outcome: Users can trigger and view scans

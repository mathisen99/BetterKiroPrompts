# Phase 2: Feature Implementation — Tasks

## Task List

### Shared Infrastructure

- [x] 1. Create `backend/internal/generator/` package structure
  - Refs: Design
  - Outcome: Generator package ready for implementation

- [x] 2. Create `backend/internal/templates/` with embed directive
  - Refs: Design
  - Outcome: Templates embedded in binary

- [x] 3. Create `frontend/src/lib/api.ts` with API client functions
  - Refs: Design
  - Outcome: Typed API client for all endpoints

- [x] 4. Create `frontend/src/components/shared/OutputPanel.tsx`
  - Refs: AC-13, AC-14, AC-15
  - Outcome: Preview/Copy/Download component

- [x] 5. Create `frontend/src/components/shared/StepIndicator.tsx`
  - Refs: Design
  - Outcome: Wizard step indicator component

### Kickoff Prompt Generator — Backend

- [x] 6. Create `backend/internal/api/kickoff.go` handler
  - Refs: AC-3
  - Deps: 1
  - Outcome: POST /api/kickoff/generate endpoint

- [x] 7. Create `backend/internal/generator/kickoff.go` logic
  - Refs: AC-3
  - Deps: 1
  - Outcome: Prompt generation from answers

- [x] 8. Create `backend/internal/templates/kickoff.tmpl`
  - Refs: AC-3
  - Deps: 2
  - Outcome: Prompt template with all questions

- [x] 9. Register kickoff route in router.go
  - Refs: AC-3
  - Deps: 6
  - Outcome: Route accessible at /api/kickoff/generate

- [x] 10. Test: Kickoff endpoint returns valid prompt
  - Refs: AC-3
  - Deps: 9
  - Outcome: curl test passes

### Kickoff Prompt Generator — Frontend

- [x] 11. Create `frontend/src/pages/KickoffPage.tsx`
  - Refs: AC-1
  - Deps: 3
  - Outcome: Page container for kickoff wizard

- [x] 12. Create `frontend/src/components/kickoff/KickoffWizard.tsx`
  - Refs: AC-1, AC-2
  - Deps: 5
  - Outcome: Multi-step wizard with state management

- [x] 13. Create `frontend/src/components/kickoff/QuestionStep.tsx`
  - Refs: AC-1
  - Outcome: Reusable question step component

- [x] 14. Implement Step 1-3: Project Identity, Success Criteria, Users & Roles
  - Refs: AC-1
  - Deps: 12, 13
  - Outcome: First three questions working

- [x] 15. Implement Step 4: Data Sensitivity with Data Lifecycle sub-fields
  - Refs: AC-1
  - Deps: 14
  - Outcome: Data sensitivity + lifecycle fields

- [x] 16. Implement Step 5-6: Auth Model, Concurrency
  - Refs: AC-1
  - Deps: 15
  - Outcome: Auth and concurrency questions

- [ ] 17. Implement Step 7: Risks & Tradeoffs with sub-fields
  - Refs: AC-1
  - Deps: 16
  - Outcome: Risks with top 3, mitigations, not handled

- [ ] 18. Implement Step 8: Boundaries with concrete examples
  - Refs: AC-1
  - Deps: 17
  - Outcome: Boundaries + 2-3 examples

- [ ] 19. Implement Step 9-10: Non-Goals, Constraints
  - Refs: AC-1
  - Deps: 18
  - Outcome: Final questions complete

- [ ] 20. Create `frontend/src/components/kickoff/PromptPreview.tsx`
  - Refs: AC-13
  - Deps: 4, 19
  - Outcome: Preview generated prompt

- [ ] 21. Integrate kickoff wizard with API
  - Refs: AC-3
  - Deps: 10, 20
  - Outcome: End-to-end kickoff generation

- [ ] 22. Add navigation to KickoffPage from App.tsx
  - Refs: AC-1
  - Deps: 11
  - Outcome: Kickoff accessible from main app

### Steering Document Generator — Backend

- [ ] 23. Create `backend/internal/api/steering.go` handler
  - Refs: AC-4, AC-5, AC-6
  - Deps: 1
  - Outcome: POST /api/steering/generate endpoint

- [ ] 24. Create `backend/internal/generator/steering.go` logic
  - Refs: AC-4, AC-5, AC-6, AC-7
  - Deps: 1
  - Outcome: Steering file generation

- [ ] 25. Create foundation steering templates (product, tech, structure)
  - Refs: AC-4
  - Deps: 2
  - Outcome: Templates with inclusion: always

- [ ] 26. Create conditional steering templates (security-go, security-web, quality-go, quality-web)
  - Refs: AC-5
  - Deps: 2
  - Outcome: Templates with fileMatch inclusion

- [ ] 27. Create AGENTS.md template
  - Refs: AC-6
  - Deps: 2
  - Outcome: AGENTS.md template

- [ ] 28. Register steering route in router.go
  - Refs: AC-4
  - Deps: 23
  - Outcome: Route accessible at /api/steering/generate

- [ ] 29. Test: Steering endpoint returns valid files
  - Refs: AC-4, AC-5, AC-6
  - Deps: 28
  - Outcome: curl test passes with correct frontmatter

### Steering Document Generator — Frontend

- [ ] 30. Create `frontend/src/pages/SteeringPage.tsx`
  - Refs: AC-4
  - Deps: 3
  - Outcome: Page container for steering config

- [ ] 31. Create `frontend/src/components/steering/SteeringConfigurator.tsx`
  - Refs: AC-4, AC-5
  - Outcome: Configuration form

- [ ] 32. Create `frontend/src/components/steering/SteeringOptions.tsx`
  - Refs: AC-5
  - Outcome: Checkboxes for conditional files

- [ ] 33. Create `frontend/src/components/steering/FilePreview.tsx`
  - Refs: AC-13
  - Deps: 4
  - Outcome: Multi-file preview with tabs

- [ ] 34. Integrate steering configurator with API
  - Refs: AC-4, AC-5, AC-6
  - Deps: 29, 33
  - Outcome: End-to-end steering generation

- [ ] 35. Add navigation to SteeringPage from App.tsx
  - Refs: AC-4
  - Deps: 30
  - Outcome: Steering accessible from main app

### Hooks Generator — Backend

- [ ] 36. Create `backend/internal/api/hooks.go` handler
  - Refs: AC-8
  - Deps: 1
  - Outcome: POST /api/hooks/generate endpoint

- [ ] 37. Create `backend/internal/generator/hooks.go` logic
  - Refs: AC-8, AC-9, AC-10, AC-11, AC-12
  - Deps: 1
  - Outcome: Hook file generation with presets

- [ ] 38. Create hook templates for Light preset
  - Refs: AC-9
  - Deps: 2
  - Outcome: Formatter hooks

- [ ] 39. Create hook templates for Basic preset additions
  - Refs: AC-10
  - Deps: 38
  - Outcome: Linter + test hooks

- [ ] 40. Create hook templates for Default preset additions
  - Refs: AC-11
  - Deps: 39
  - Outcome: Secret scan + prompt guard hooks

- [ ] 41. Create hook templates for Strict preset additions
  - Refs: AC-12
  - Deps: 40
  - Outcome: Static analysis + vuln scan hooks

- [ ] 42. Register hooks route in router.go
  - Refs: AC-8
  - Deps: 36
  - Outcome: Route accessible at /api/hooks/generate

- [ ] 43. Test: Hooks endpoint returns valid hook files
  - Refs: AC-8, AC-9, AC-10, AC-11, AC-12
  - Deps: 42
  - Outcome: curl test passes with valid JSON schema

### Hooks Generator — Frontend

- [ ] 44. Create `frontend/src/pages/HooksPage.tsx`
  - Refs: AC-8
  - Deps: 3
  - Outcome: Page container for hooks config

- [ ] 45. Create `frontend/src/components/hooks/HooksPresetSelector.tsx`
  - Refs: AC-9, AC-10, AC-11, AC-12
  - Outcome: Preset selection UI

- [ ] 46. Create `frontend/src/components/hooks/PresetCard.tsx`
  - Refs: AC-9, AC-10, AC-11, AC-12
  - Outcome: Card showing preset details

- [ ] 47. Create `frontend/src/components/hooks/HookFilePreview.tsx`
  - Refs: AC-13
  - Deps: 4
  - Outcome: Hook file preview

- [ ] 48. Integrate hooks selector with API
  - Refs: AC-8
  - Deps: 43, 47
  - Outcome: End-to-end hooks generation

- [ ] 49. Add navigation to HooksPage from App.tsx
  - Refs: AC-8
  - Deps: 44
  - Outcome: Hooks accessible from main app

### Integration & Polish

- [ ] 50. Create main navigation component
  - Refs: Design
  - Outcome: Nav between Kickoff, Steering, Hooks pages

- [ ] 51. Update App.tsx with routing
  - Refs: Design
  - Deps: 22, 35, 49, 50
  - Outcome: All pages accessible via routes

- [ ] 52. Implement download as zip for multi-file outputs
  - Refs: AC-15
  - Deps: 4
  - Outcome: Zip download for steering and hooks

- [ ] 53. Test: Full kickoff flow end-to-end
  - Refs: AC-1, AC-2, AC-3
  - Deps: 21
  - Outcome: Complete wizard generates valid prompt

- [ ] 54. Test: Full steering flow end-to-end
  - Refs: AC-4, AC-5, AC-6, AC-7
  - Deps: 34
  - Outcome: Config generates valid steering files

- [ ] 55. Test: Full hooks flow end-to-end
  - Refs: AC-8, AC-9, AC-10, AC-11, AC-12
  - Deps: 48
  - Outcome: Preset generates valid hook files

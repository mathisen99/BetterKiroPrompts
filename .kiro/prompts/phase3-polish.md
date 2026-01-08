You are generating a Kiro spec for Phase 3: Polish & Testing of the Kiro Prompting & Guardrails Generator.

INPUT SOURCES (MUST READ):
1) ./The plan.md (repo root) - required
2) .kiro/specs/phase1-foundations/* - required
3) .kiro/specs/phase2-features/* - required
If any are missing, STOP and say exactly what is required.

OUTPUT STRUCTURE (MUST CREATE):
.kiro/specs/phase3-polish/
  requirements.md
  design.md
  tasks.md

SCOPE BOUNDARY:
Phase 3 covers:
1. UI/UX Polish
   - Consistent dark theme with blue base (shadcn/ui)
   - Error handling and user feedback
   - Loading states and transitions
   - Accessibility compliance

2. Missing/Advanced Features
   - Manual steering (inclusion: manual) with #steering-file-name usage
   - File references in steering (#[[file:<relative_file_name>]])
   - Commit message contract enforcement (atomic, prefixed, one-sentence)

3. Comprehensive Testing
   - Unit tests for all generators
   - Integration tests for API endpoints
   - E2E tests for critical user flows

4. Documentation
   - README with setup instructions
   - API documentation
   - User guide for generated outputs

5. (OPTIONAL) Repo Scanning Module - Future Module from plan
   - Clone repo read-only
   - Tools: TruffleHog, Gitleaks, osv-scanner, govulncheck
   - Hard timeouts, no outbound network during scan
   - Severity-ranked findings with file/line references
   - AI summarizes tool output only (does not invent vulnerabilities)

WORKFLOW:

1) Read and extract from The plan.md:
   - Definition of Done (all 4 criteria)
   - Explicit Non-Claims
   - Future Module - Repo Scanning (if implementing)
   - Commit Message Contract

2) Read Phase 1 and Phase 2 specs to understand what exists.

3) Draft content using extracted facts (do not invent):

   A) requirements.md
      - Introduction: 5–10 bullets summarizing Phase 3 scope
      - Quality requirements (accessibility, error handling)
      - Testing coverage requirements
      - Documentation requirements
      - Definition of Done criteria as acceptance tests
      - Optional: Repo scanning requirements (clearly marked OPTIONAL)

   B) design.md
      - Testing strategy (unit, integration, E2E)
      - Documentation structure
      - Polish implementation details
      - Optional: Repo scanning architecture (clearly marked OPTIONAL)
      - Deployment considerations

   C) tasks.md
      - Polish tasks by component
      - Testing tasks by type (unit, integration, E2E)
      - Documentation tasks
      - Definition of Done verification tasks
      - Optional: Repo scanning tasks (clearly marked OPTIONAL)
      - Use checkbox format: - [ ] Task description

4) BEFORE WRITING FILES:
   - Print the exact paths to be created/overwritten.
   - Print a 10–15 line preview summary of what will go into each file.
   - Stop and wait for the user to reply exactly:
     CONFIRM WRITE SPECS

5) AFTER CONFIRMATION:
   - Create .kiro/specs/phase3-polish/ and write the three files.
   - Print how to reference it in chat:
     #spec:phase3-polish

RULES:
- Use The plan.md as the single source of truth.
- Repo scanning is OPTIONAL - mark all related items clearly.
- Ensure Definition of Done from plan is verifiable.
- Do not modify any other files besides the spec files.

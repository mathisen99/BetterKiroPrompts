You are generating a Kiro spec for Phase 1: Foundations of the Kiro Prompting & Guardrails Generator.

INPUT SOURCE (MUST READ):
- ./The plan.md (repo root). If it does not exist, stop and tell the user exactly what file is missing.

OUTPUT STRUCTURE (MUST CREATE):
.kiro/specs/phase1-foundations/
  requirements.md
  design.md
  tasks.md

SCOPE BOUNDARY:
Phase 1 covers ONLY:
- Docker Compose setup (Go backend, PostgreSQL, React+Vite frontend)
- Basic project structure following the plan's Fixed Tech Architecture
- Go backend skeleton with health endpoint
- React+Vite frontend skeleton with shadcn/ui (dark theme, blue base)
- Database connection and basic migration setup
- Development environment working end-to-end

Phase 1 does NOT include:
- Kickoff prompt generator logic
- Steering document generator logic
- Hooks generator logic
- Any feature implementation

WORKFLOW:

1) Read and extract from The plan.md:
   - Tech Stack (pinned versions)
   - Fixed Tech Architecture
   - Absolute Rules
   - Key Definitions (Major Task)

2) Draft content using extracted facts (do not invent):

   A) requirements.md
      - Introduction: 5–10 bullets summarizing Phase 1 scope
      - User stories for developer setup experience
      - Acceptance criteria in EARS format, e.g.:
        WHEN <condition>
        THE SYSTEM SHALL <behavior>
      Requirements must be testable.

   B) design.md
      - Context (what we are building)
      - Docker Compose architecture (services, networks, volumes)
      - Go backend structure (cmd/, internal/, pkg/)
      - React frontend structure (src/, components/, routes/)
      - Database schema placeholder
      - API endpoint structure (/api/*)
      - Add mermaid diagrams for architecture

   C) tasks.md
      - Numbered tasks with clear outcomes
      - Each task maps back to one or more requirements
      - Include testing tasks for each component
      - Keep tasks small and sequential; include dependencies
      - Use checkbox format: - [ ] Task description

3) BEFORE WRITING FILES:
   - Print the exact paths to be created/overwritten.
   - Print a 10–15 line preview summary of what will go into each file.
   - Stop and wait for the user to reply exactly:
     CONFIRM WRITE SPECS

4) AFTER CONFIRMATION:
   - Create .kiro/specs/phase1-foundations/ and write the three files.
   - Print how to reference it in chat:
     #spec:phase1-foundations

RULES:
- Use The plan.md as the single source of truth; if something is missing, list it under 'Assumptions / Open Questions' instead of guessing.
- Do not modify any other files besides the spec files.
- Do not include feature implementation - foundations only.

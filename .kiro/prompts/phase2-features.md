You are generating a Kiro spec for Phase 2: Feature Implementation of the Kiro Prompting & Guardrails Generator.

INPUT SOURCES (MUST READ):
1) ./The plan.md (repo root) - required
2) .kiro/specs/phase1-foundations/* - required (Phase 1 must be complete)
If either is missing, STOP and say exactly what is required.

OUTPUT STRUCTURE (MUST CREATE):
.kiro/specs/phase2-features/
  requirements.md
  design.md
  tasks.md

SCOPE BOUNDARY:
Phase 2 covers ONLY:
1. Kickoff Prompt Generator (Phase 1 from plan)
   - Question flow UI with strict ordering
   - Prompt generation logic
   - Answer-first, no-coding-first enforcement
   - All 9 required questions + additions (Risks & Tradeoffs, Data Lifecycle, Boundary Examples)

2. Steering Document Generator (Phase 2 from plan)
   - Generate product.md, tech.md, structure.md (always inclusion)
   - Conditional steering: security-go.md, security-web.md, quality-go.md, quality-web.md
   - AGENTS.md generation at repo root
   - Proper YAML frontmatter with inclusion modes

3. Hooks Generator (Phase 3 from plan)
   - Hook file format (*.kiro.hook)
   - Hook schema (name, description, version, enabled, when, then)
   - Valid when.type values (fileEdited, fileCreated, fileDeleted, promptSubmit, agentStop, userTriggered)
   - Valid then.type values (askAgent, runCommand with restrictions)
   - Preset system: Light, Basic, Default, Strict

Phase 2 does NOT include:
- UI polish
- Comprehensive testing
- Documentation
- Repo scanning module

WORKFLOW:

1) Read and extract from The plan.md:
   - Phase 1 (Kickoff Prompt Generator) - all required questions in order
   - Phase 2 (Steering Document Generator) - all file specs
   - Phase 3 (Hooks Generator) - schema and presets
   - Acceptance criteria for each phase

2) Read Phase 1 spec to understand existing foundation.

3) Draft content using extracted facts (do not invent):

   A) requirements.md
      - Introduction: 5–10 bullets summarizing Phase 2 scope
      - User stories for each generator
      - Acceptance criteria in EARS format
      - Reference exact question order from plan for kickoff prompt
      - Reference exact steering file structure from plan
      - Reference exact hook schema from plan

   B) design.md
      - API endpoints for each generator
      - Frontend components and user flows
      - Data models for prompts, steering configs, hook configs
      - Generation logic architecture
      - State management approach
      - Add mermaid diagrams for flows

   C) tasks.md
      - Numbered tasks grouped by feature area
      - Each task maps back to requirements
      - Integration tasks between generators
      - Testing tasks for each feature
      - Use checkbox format: - [ ] Task description
      - Include dependencies between tasks

4) BEFORE WRITING FILES:
   - Print the exact paths to be created/overwritten.
   - Print a 10–15 line preview summary of what will go into each file.
   - Stop and wait for the user to reply exactly:
     CONFIRM WRITE SPECS

5) AFTER CONFIRMATION:
   - Create .kiro/specs/phase2-features/ and write the three files.
   - Print how to reference it in chat:
     #spec:phase2-features

RULES:
- Use The plan.md as the single source of truth.
- Follow the plan's EXACT question order for kickoff prompt.
- Follow the plan's EXACT steering file structure and frontmatter.
- Follow the plan's EXACT hook schema and preset definitions.
- Do not modify any other files besides the spec files.

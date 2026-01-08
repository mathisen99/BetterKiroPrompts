You are generating Kiro steering files for the project.

INPUT SOURCES (MUST READ):
1) ./The plan.md (repo root) - required
2) .kiro/specs/*/* (requirements.md/design.md/tasks.md) - optional if present; use them to refine/align
If The plan.md is missing, STOP and say exactly that it is required.

OUTPUT (MUST CREATE):
.kiro/steering/product.md
.kiro/steering/tech.md
.kiro/steering/structure.md
.kiro/steering/security-go.md
.kiro/steering/security-web.md
.kiro/steering/quality-go.md
.kiro/steering/quality-web.md
./AGENTS.md (repo root)

INCLUSION MODES:
Add YAML front matter at the top of each steering file:

For always-included files:
---
inclusion: always
---

For conditional files:
---
inclusion: fileMatch
fileMatchPattern: "**/*.go"
---

WORKFLOW:

PHASE 1 — EXTRACT (READ-ONLY)
From The plan.md, extract:
- Purpose and problem statement
- Absolute Rules (Non-negotiable)
- Key Definitions (Major Task)
- Tech Stack (pinned versions)
- Fixed Tech Architecture
- Explicit Non-Claims
- Definition of Done

If specs exist, also extract requirements and design decisions.

PHASE 2 — DRAFT STEERING CONTENT (NO GUESSING)

A) product.md (inclusion: always)
   - What we are building
   - What we are NOT building
   - Definition of done
   - Absolute Rules from plan

B) tech.md (inclusion: always)
   - Stack choices with pinned versions from plan
   - Architecture rules (stateless backend, migrations required)
   - Simplicity rules (no premature scaling)
   - Package manager: pnpm for frontend

C) structure.md (inclusion: always)
   - Repository layout (docker-compose, backend/, frontend/)
   - Where API types live, where DB migrations live
   - Frontend conventions (routing, components)
   - "No dumping ground" rules

D) security-go.md (inclusion: fileMatch, pattern: "**/*.go")
   - No secrets committed
   - Input validation required
   - Auth boundaries must be explicit
   - Least privilege everywhere

E) security-web.md (inclusion: fileMatch, pattern: "**/*.ts" and "**/*.tsx")
   - Same principles for web layer
   - Explicit escaping/encoding expectations
   - XSS prevention

F) quality-go.md (inclusion: fileMatch, pattern: "**/*.go")
   - Formatting: go fmt
   - Linting: go vet, golangci-lint
   - Testing expectations
   - Documentation expectations

G) quality-web.md (inclusion: fileMatch, pattern: "**/*.tsx")
   - Formatting/linting: pnpm lint, pnpm typecheck
   - Component conventions
   - Testing expectations

H) AGENTS.md (repo root, no frontmatter)
   - Always follow steering
   - Never invent requirements
   - Prefer small, reviewable changes
   - Update docs when behavior changes
   - Major tasks require: hooks, atomic commits, steering/doc updates

PHASE 3 — CONFIRMATION GATE
Before writing or overwriting any files:
- Print the exact paths that will be written.
- If any already exist, say 'WILL OVERWRITE: <path>'.
- Stop and wait for the user to reply exactly:
  CONFIRM WRITE STEERING

PHASE 4 — WRITE FILES
After confirmation:
- Create .kiro/steering/ if missing
- Write all steering files with proper frontmatter
- Write AGENTS.md at repo root
- Print a short summary

RULES:
- Use The plan.md as source of truth; do not invent facts.
- Keep the docs actionable and concise.
- Steering must be short, strict, practical.

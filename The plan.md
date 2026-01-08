# Kiro Prompting & Guardrails Generator — Complete Roadmap

## Purpose

This project solves one problem:

Beginners often “vibe-code” without understanding architecture, security, data, or concurrency. We improve their thinking by generating **better prompts**, **steering documents**, and **Kiro hooks**—not by writing the application for them.

This project **does not replace Kiro specs**. Kiro already does specs well. Our value is improving **inputs** to Kiro (kickoff prompt + steering + hooks), not outputs.

---

## Absolute Rules (Non‑negotiable)

1. **Do not generate application code** unless explicitly instructed.
2. **Do not invent requirements.** If required input is missing or ambiguous, stop and ask.
3. **Do not skip the prompt flow.** Execute phases in order.
4. **Do not overbuild.** Prefer the simplest working solution.
5. **Do not claim security guarantees.** Provide safer defaults and guardrails only.
6. **Major tasks require discipline:** hooks, atomic commits, and steering/doc updates.

---

## Key Definitions

### Major Task

A “major task” is any change that affects one or more of:

* Application behavior
* API surface
* Authentication/authorization
* Database schema or data model
* Security posture
* Deployment/runtime configuration

Major tasks must:

* Trigger relevant hooks
* Result in **clean, atomic commits**
* Include **documentation/steering updates** when behavior changes

---

## Roadmap Overview

### Phase 0 — Foundations

Goal: Ensure the project is aligned with **how Kiro actually works**.

Deliverables:

* A reference section describing:

  * Kiro steering inclusion modes: `always`, `fileMatch`, `manual`
  * Steering frontmatter keys (notably `inclusion` and `fileMatchPattern`)
  * Hooks file format and schema (`.kiro.hook`, `when/then`, and valid enum types)
  * Operational constraints: one hook at a time, workspace trust, manual triggers

Acceptance criteria:

* The plan’s language and file formats match Kiro’s documented behavior.

---

## Phase 1 — Base Project Creation Prompt (Kickoff Prompt Generator)

Goal: Generate **one** “project kickoff prompt” that enforces the thinking sequence and blocks coding until answered.

### Output

* A **single** kickoff prompt text artifact.
* It must:

  * Ask questions in a strict order
  * Instruct that **no coding is allowed** until the questions are answered
  * Stop after asking questions (wait for answers)

### Required Questions (in order)

1. **Project Identity**: Restate the project in one sentence.
2. **Success Criteria**: What does “done” mean?
3. **Users & Roles**: Who uses it? (anonymous/auth/admin/etc.)
4. **Data Sensitivity**: What data is stored? Label sensitive data explicitly.
5. **Auth Model**: none / basic / external provider.
6. **Concurrency Expectations**: multi-user? background jobs? shared state?
7. **Boundaries**: public vs private data boundaries.
8. **Non-Goals**: what will not be built.
9. **Constraints**: time, simplicity, tech limits.

### Additions to strengthen beginner thinking (still same major idea)

* **Risks & Tradeoffs** (immediately after Concurrency):

  * Top 3 risks
  * Simplest mitigations
  * What is explicitly not handled
* **Data Lifecycle** (within Data Sensitivity):

  * retention, deletion, export, audit logging, backups
* **Boundary Examples** (within Boundaries):

  * 2–3 concrete access examples (who can read/write what)

### Acceptance Criteria

* Kickoff prompt is generated as **one** cohesive artifact.
* It enforces: answer-first, no-coding-first.

---

## Phase 2 — Steering Document Generator

Goal: Generate concise, opinionated steering files under `.kiro/steering/` using Kiro inclusion modes.

### Steering File Scopes

* **Workspace steering:** lives in the repo under `.kiro/steering/`.
* **Global steering (optional):** can live under `~/.kiro/steering/` and applies across workspaces.

### Steering Files (Foundation Set)

Create these files:

1. `.kiro/steering/product.md`

* **Frontmatter:** `inclusion: always`
* Content:

  * What we are building
  * What we are not building
  * Definition of done

2. `.kiro/steering/tech.md`

* **Frontmatter:** `inclusion: always`
* Content:

  * Stack choices (Go backend, React+Vite frontend, PostgreSQL)
  * Architecture rules (stateless backend, migrations required)
  * Simplicity rules (no premature scaling)

3. `.kiro/steering/structure.md`

* **Frontmatter:** `inclusion: always`
* Content:

  * Repository layout (folders, boundaries)
  * Where API types live, where DB migrations live
  * Frontend conventions (routing, components)
  * “No dumping ground” rules

### Conditional Steering Files

4. `.kiro/steering/security-go.md`

* **Frontmatter:** `inclusion: fileMatch`, `fileMatchPattern: "**/*.go"`
* Content:

  * No secrets committed
  * Input validation required
  * Auth boundaries must be explicit
  * Least privilege everywhere

5. `.kiro/steering/security-web.md`

* **Frontmatter:** `inclusion: fileMatch`, `fileMatchPattern: "**/*.ts"` (and a separate file for `**/*.tsx` if needed)
* Content:

  * Same principles for web layer
  * Explicit escaping/encoding expectations

6. `.kiro/steering/quality-go.md`

* **Frontmatter:** `inclusion: fileMatch`, `fileMatchPattern: "**/*.go"`
* Content:

  * Formatting
  * Linting
  * Testing
  * Documentation expectations

7. `.kiro/steering/quality-web.md`

* **Frontmatter:** `inclusion: fileMatch`, `fileMatchPattern: "**/*.tsx"`
* Content:

  * Formatting/linting
  * Component conventions
  * Testing expectations

### Manual Steering (Optional)

* Use `inclusion: manual` for specialized docs users pull in via `#steering-file-name`.

### AGENTS.md

Create `AGENTS.md` at repo root.

* Content:

  * Always follow steering
  * Never invent requirements
  * Prefer small, reviewable changes
  * Update docs when behavior changes

### Constraints

* Steering must be **short, strict, practical**.
* Use file references when helpful (e.g., reference OpenAPI or `.env.example`).

### Acceptance Criteria

* Files exist under `.kiro/steering/` with valid frontmatter.
* Guidance is minimal, enforceable, and consistent.

---

## Phase 3 — Hooks and Agents Configuration Generator

Goal: Generate valid **Kiro IDE hooks** located in `.kiro/hooks/`.

### File Format

* Use the proper hook file extension: `*.kiro.hook`.

### Hook Schema (IDE)

Each hook includes (at minimum):

* `name`
* `description`
* `version`
* `enabled`
* `when` (with `type` and optional fields)
* `then` (with `type` and config)

### Valid `when.type` values

* File hooks: `fileEdited`, `fileCreated`, `fileDeleted`
* Contextual hooks: `promptSubmit`, `agentStop`
* Manual hooks: `userTriggered`

### File hook matching

* For file hooks, `when.patterns` is **required** and uses workspace-relative glob patterns.
* Keep patterns narrow to avoid performance issues and accidental triggering.

### Valid `then.type` values

* `askAgent` (allowed everywhere)
* `runCommand` (allowed only for `promptSubmit` and `agentStop` triggers)

### Design Rules

* **No `runCommand` on file event hooks.** Use `askAgent` to instruct checks or request a follow-up action.
* **Prefer frictionless defaults.** Heavy scans should be manual (`userTriggered`).
* Remember Kiro operational constraints:

  * Only one hook runs at a time (loop prevention)
  * Hooks are disabled in untrusted workspaces

### Presets

Support three presets.

#### Preset A — Light

Goal: minimum friction.

* On `agentStop`:

  * `runCommand` formatters (Go + frontend)

#### Preset B — Basic

Goal: daily discipline.

* On `agentStop`:

  * formatters + basic linters
* Manual `userTriggered`:

  * run unit tests and summarize results

#### Preset C — Default (Recommended)

Goal: balanced safety.

* Everything in Basic
* Add:

  * quick secret scan on `agentStop`
  * prompt guardrails on `promptSubmit` (block obviously unsafe prompts, require confirmation steps)

#### Preset D — Strict (Optional)

Goal: maximum enforcement.

* Everything in Default
* Add:

  * static analysis on `agentStop`
  * dependency vulnerability scan as manual or `agentStop` depending on performance

### Commit Message Contract

All commits must:

* Be atomic (one concern per commit)
* Use a prefix: `feat:`, `fix:`, `docs:`, `chore:`
* Have a one-sentence summary
* Avoid mixed or vague commits

### Acceptance Criteria

* Hook files are valid and recognized by Kiro.
* Presets generate correct hook sets.

---

## Future Module — Repo Scanning (Separate Site Feature)

This is **not part of the core prompt/steering/hooks generation flow**. It is a separate feature area of the site with its own UX, permissions, and execution model.

### Scope boundary

* Core flow (Phases 1–3): generate kickoff prompt + steering + hooks.
* Repo scanning: a separate module that can be added later without affecting the core flow.

### Behavior Rules

* Clone repo **read-only**.
* Run local tools with **hard timeouts**.
* No outbound network access during scan.
* Partial results are acceptable.

### Tools (MVP)

* TruffleHog
* Gitleaks
* `osv-scanner`
* `govulncheck` (Go)

### Output Expectations

* Severity-ranked findings
* File and line references
* “Why it matters” context
* Minimal fix guidance

### AI Usage Policy

* AI only summarizes and prioritizes **tool output**.
* AI must not invent vulnerabilities.

### Trigger model

* Expose as an explicit user action (site button / job run).
* If integrated into Kiro later, prefer a **manual (`userTriggered`) hook** for scans.

### Acceptance Criteria

* Scan runs end-to-end.
* Output is structured and traceable to tool output.

## Tech Stack (Pinned Versions)

Pinned to current stable releases (update policy: review quarterly).

* **Go:** 1.25.5
* **PostgreSQL:** 18.1
* **Node.js (LTS):** 24.12.0 (meets Vite’s Node requirements)
* **React:** 19.2.3
* **Vite:** 7.3.1
* **shadcn/ui:** use the official shadcn/ui docs + CLI; `shadcn-ui` npm package 0.9.5
* **Docker Engine:** 29.1.3
* **Docker Desktop:** 4.55.0 (optional, for dev convenience)
* **Docker Compose plugin:** 5.0.1

## AI Models

* **Prompt Generation (kickoff prompts, steering, hooks):** GPT-5.2
  * Input: $1.75/1M tokens
  * Cached: $0.175/1M tokens
  * Output: $14.00/1M tokens

* **Code Review & Repo Scanning:** GPT-5.1-Codex-Max
  * Input: $1.25/1M tokens
  * Cached: $0.125/1M tokens
  * Output: $10.00/1M tokens
  * Guide: https://cookbook.openai.com/examples/gpt-5/gpt-5-1-codex-max_prompting_guide

## Fixed Tech Architecture (For MVP)

* Containerization: Docker Compose–based project.
* Backend: Go serving `/api/*` JSON API + serving built static frontend.
* Frontend: React + Vite + shadcn/ui (dark theme, blue base).
* Database: PostgreSQL in same compose stack.
* Optional Worker: Repo scanning only (future module).

---

## Explicit Non‑Claims

This system:

* Does not guarantee security
* Does not replace human review
* Does not teach programming
* Does not enforce correctness

It provides guardrails and better thinking.

---

## Definition of Done

The project is complete when:

1. A user can generate a full Kiro kickoff prompt.
2. Steering files are generated, usable, correctly scoped, and concise.
3. Hooks are generated, valid, and usable.
4. (Optional) Repo scanning works end-to-end.

Do not proceed beyond this unless explicitly instructed.

# User Guide

## Kickoff Prompts

A kickoff prompt is a structured set of questions that forces you to think through your project before writing code. It covers:

1. **Project Identity** — What are you building?
2. **Success Criteria** — What does "done" mean?
3. **Users & Roles** — Who uses it?
4. **Data Sensitivity** — What data is stored? What's sensitive?
5. **Auth Model** — How do users authenticate?
6. **Concurrency** — Multi-user? Background jobs?
7. **Risks & Tradeoffs** — What could go wrong?
8. **Boundaries** — What's public vs private?
9. **Non-Goals** — What won't you build?
10. **Constraints** — Time, tech, simplicity limits?

### How to Use

1. Complete the wizard, answering each question
2. Generate the prompt
3. Copy or download the output
4. Paste into Kiro as your project kickoff

The prompt instructs Kiro to ask clarifying questions before writing any code.

---

## Steering Files

Steering files live in `.kiro/steering/` and guide Kiro's behavior. They use YAML frontmatter to control when they're included.

### Inclusion Modes

| Mode | Frontmatter | When Included |
|------|-------------|---------------|
| `always` | `inclusion: always` | Every conversation |
| `fileMatch` | `inclusion: fileMatch`<br>`fileMatchPattern: "**/*.go"` | When editing matching files |
| `manual` | `inclusion: manual` | Only when referenced via `#steering-file-name` |

### Generated Files

**Foundation (always included):**
- `product.md` — What you're building, not building, definition of done
- `tech.md` — Stack choices, architecture rules, simplicity rules
- `structure.md` — Repository layout, conventions

**Conditional (fileMatch):**
- `security-go.md` — Go security rules (no secrets, input validation)
- `security-web.md` — Web security rules
- `quality-go.md` — Go quality rules (formatting, testing)
- `quality-web.md` — Web quality rules

**Manual:**
- Referenced explicitly when needed via `#steering-file-name`

### File References

Steering files can reference other files using:
```
#[[file:.env.example]]
#[[file:backend/migrations/README.md]]
```

---

## Hooks

Hooks are automated actions that run at specific points. They live in `.kiro/hooks/` as `*.kiro.hook` files.

### Presets

| Preset | What It Does |
|--------|--------------|
| **Light** | Formatters only (Go + frontend) |
| **Basic** | Formatters + linters + manual test trigger |
| **Default** | Basic + secret scan + prompt guardrails |
| **Strict** | Default + static analysis + vulnerability scan |

### When Hooks Run

- `agentStop` — After Kiro finishes a task
- `promptSubmit` — Before Kiro processes your prompt
- `userTriggered` — Manual trigger only

### Tech Stack Options

Hooks are generated based on your tech stack:
- **Go** — `go fmt`, `go vet`, `govulncheck`
- **TypeScript** — `prettier`, `eslint`, `tsc`
- **React** — Additional React-specific linting

---

## Commit Message Contract

All commits should follow this contract:

1. **Atomic** — One concern per commit
2. **Prefixed** — Use `feat:`, `fix:`, `docs:`, `chore:`, `test:`, `refactor:`
3. **One-sentence summary** — Under 72 characters

**Examples:**
```
feat: add user authentication
fix: resolve database connection timeout
docs: update API documentation
chore: upgrade dependencies
```

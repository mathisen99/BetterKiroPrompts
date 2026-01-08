# Kiro CLI Workflow Cheatsheet

## Quick Start
```
/steering              # Generate steering files (once)
/phase1-foundations    # Generate Phase 1 spec
/next                  # Work loop - repeat until done
```

## The Loop
```
/next → finds task → implements → reviews → commits → marks done
```
Just keep running `/next` until all tasks are complete.

## Phase Transitions
```
Phase 1 done → /readme → /phase2-features → /next ...
Phase 2 done → /readme → /phase3-polish   → /next ...
Phase 3 done → /readme
```

---

## All Prompts

| Prompt | What it does |
|--------|--------------|
| `/steering` | Generates steering files + AGENTS.md |
| `/phase1-foundations` | Generates Phase 1 spec (setup, docker, skeleton) |
| `/phase2-features` | Generates Phase 2 spec (generators) |
| `/phase3-polish` | Generates Phase 3 spec (polish, testing, docs) |
| `/next` | **Full automation** - find → implement → review → commit |
| `/task` | Work on specific task (manual pick) |
| `/task-list` | View open/done tasks |
| `/task-complete` | Mark a task done manually |
| `/commit` | Quality gates + commit manually |
| `/review` | Code review anytime |
| `/readme` | Generate/update README after phase completion |

---

## Quality Gates (run before commits)

**Backend (Go)**
```
go fmt ./...
go vet ./...
golangci-lint run ./...
```

**Frontend (React/TS)**
```
pnpm typecheck
pnpm lint
```

---

## Confirmation Keywords

| Prompt | Keyword |
|--------|---------|
| Spec generators | `CONFIRM WRITE SPECS` |
| Steering | `CONFIRM WRITE STEERING` |
| `/next` implement | `PROCEED` or `SKIP` |
| `/next` commit | `YES`, `NO`, or `FIX <issue>` |
| `/task-complete` | `CONFIRM MARK DONE` |
| `/commit` | `CONFIRM COMMIT` |
| `/readme` | `CONFIRM WRITE README` |

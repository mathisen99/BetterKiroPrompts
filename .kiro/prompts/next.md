You are running the full task automation loop.

This prompt: finds next task → implements it → reviews → asks to commit → marks done.

WHERE TASKS LIVE:
.kiro/specs/<spec_slug>/tasks.md

---

PHASE 1 — FIND NEXT TASK

1) Discover specs:
- List all .kiro/specs/*/tasks.md
- If none exist, STOP and say no spec tasks files were found.

2) Parse all tasks files:
- If the file uses checkboxes: treat "- [ ]" as open and "- [x]" as done
- If no checkboxes: treat lines containing "(done)" or "[x]" as done, otherwise open

3) Select next task:
- Find the first open task that has no unfinished dependencies
- If dependencies are explicitly stated, verify they are marked done
- If a task has unfinished dependencies, skip to the next open task

4) If no open tasks remain across all specs:
- STOP and say "All tasks complete! Run /phase2-features or /phase3-polish if there's a next phase."

5) Print task selection:
- Spec slug
- Task number
- Full task description
- Open/done counts for that spec

---

PHASE 2 — LOAD CONTEXT (MUST READ ALL)

Read these files completely:
- .kiro/specs/<spec_slug>/requirements.md
- .kiro/specs/<spec_slug>/design.md
- .kiro/specs/<spec_slug>/tasks.md
- All files in .kiro/steering/*
- ./The plan.md (for absolute rules and constraints)
- Any files explicitly referenced in the task description or design

Extract from context:
- Which requirements this task maps to
- Architecture/design decisions that apply
- Steering rules that must be followed
- Existing code patterns to match

---

PHASE 3 — PLAN THE WORK

Before writing ANY code, print:

TASK: <one-line summary>
SPEC: <spec_slug>
REQUIREMENTS: <which requirements from requirements.md this satisfies>

FILES TO CREATE:
- <path> — <purpose>

FILES TO MODIFY:
- <path> — <what changes>

APPROACH:
1. <step>
2. <step>
3. <step>
(2-5 bullets, concrete actions)

STEERING RULES THAT APPLY:
- <relevant rules from steering files>

Ask: Type 'PROCEED' to start implementation, or 'SKIP' to pick a different task.

If user says SKIP, return to PHASE 1 and show remaining open tasks to choose from.

---

PHASE 4 — IMPLEMENT

After PROCEED confirmation:

Rules during implementation:
- Follow steering rules strictly
- Follow design.md architecture exactly
- Write minimal code to satisfy the task — no more
- Do not add features beyond task scope
- Do not invent requirements
- Match existing code patterns and conventions
- Include necessary error handling
- Add comments only where logic is non-obvious

For Go code:
- Proper error wrapping
- Context usage for cancellation
- Resource cleanup with defer
- Input validation

For React/TypeScript code:
- Proper TypeScript types (no `any`)
- Hook rules compliance
- Accessibility attributes
- Proper effect cleanup

---

PHASE 5 — QUALITY GATES (FAIL FAST)

Run checks in order. If ANY fail, fix before continuing.

Backend (if Go files were created/modified):
```
go fmt ./...
go vet ./...
golangci-lint run ./...
```

Frontend (if TS/TSX files were created/modified):
```
pnpm typecheck
pnpm lint
```

If a check fails:
- Print the error clearly
- Fix the issue
- Re-run the failed check
- Continue only when all checks pass

---

PHASE 6 — SELF-REVIEW

Review your own implementation against:

1) TASK REQUIREMENTS
- Does it fully satisfy the task description?
- Does it meet the mapped requirements from requirements.md?

2) DESIGN COMPLIANCE
- Does it follow the architecture in design.md?
- Are components/modules in the right places?

3) STEERING COMPLIANCE
- No secrets committed?
- Input validation present?
- Auth boundaries explicit?
- Follows simplicity rules?

4) CODE QUALITY
- Error handling correct?
- No obvious bugs?
- No unnecessary complexity?
- Tests needed? (note if missing but don't block)

Print review summary:

IMPLEMENTATION REVIEW
=====================
Task: <task description>
Status: <PASS / ISSUES FOUND>

What was implemented:
- <bullet>
- <bullet>

Files changed:
- <path> (+X -Y lines)

Quality gates: <PASSED / which ones>

Concerns or notes:
- <any issues, or "None">

---

PHASE 7 — COMMIT DECISION

Run and print:
```
git diff --stat
```

Show proposed commit message:
```
<type>: <concise summary under 72 chars>

- <bullet derived from diff>
- <bullet derived from diff>
```

Commit message rules:
- Types: feat, fix, refactor, chore, test, docs
- Summary under 72 characters
- Bullets derived strictly from the diff
- Do not reference task numbers or spec names unless visible in code
- Atomic: one concern per commit

Ask:
Type 'YES' to commit and mark task done.
Type 'NO' to mark task done without committing.
Type 'FIX <description>' to address an issue first.

---

PHASE 8A — IF YES (COMMIT AND COMPLETE)

1) Stage and commit:
```
git add -A
git commit -m "<generated message>"
```

2) Mark task complete in tasks.md:
- Locate the exact task line
- If checkbox format "- [ ]" → change to "- [x]"
- If numbered without checkbox: append " (done)" or convert to "X. [x] Task"
- Do not reformat other lines

3) Print:
```
✓ Task committed and marked complete.
  Commit: <type>: <summary>
  
Run /next for the next task.
```

---

PHASE 8B — IF NO (COMPLETE WITHOUT COMMIT)

1) Mark task complete in tasks.md (same logic as above)

2) Print:
```
✓ Task marked complete. Changes NOT committed.
  Run /commit when ready to commit.
  Run /next for the next task.
```

---

PHASE 8C — IF FIX (ADDRESS ISSUE)

1) Parse what needs fixing from user input
2) Make the fix
3) Return to PHASE 5 (quality gates) and continue from there

---

SAFETY RULES:
- Never commit without explicit 'YES' confirmation
- Never mark task done without completing implementation
- Never implement tasks with unfinished dependencies
- Never invent requirements beyond task and spec
- Never skip the planning step
- Only change one task line when marking complete
- Follow steering at all times

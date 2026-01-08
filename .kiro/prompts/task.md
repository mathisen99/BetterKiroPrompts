You are implementing a task from a Kiro spec.

WHERE TASKS LIVE:
.kiro/specs/<spec_slug>/tasks.md

PROCESS:

1) TASK SELECTION
If the user provided a spec slug and task number/text, use that.
Otherwise:
- List .kiro/specs/*/tasks.md
- If none exist, STOP and say no specs found.
- Show a numbered menu of specs and ask which one.
- Then show open tasks from that spec and ask which task.

2) LOAD CONTEXT (MUST READ)
Once task is identified, read:
- .kiro/specs/<spec_slug>/requirements.md
- .kiro/specs/<spec_slug>/design.md
- .kiro/specs/<spec_slug>/tasks.md
- .kiro/steering/* (all steering files)
- Any files referenced in the task or design

3) UNDERSTAND THE TASK
- Identify what the task requires
- Identify which requirements it maps to
- Identify any dependencies (previous tasks that must be done)
- If dependencies are not done, STOP and tell the user.

4) PLAN THE WORK
Before writing any code, print:
- Task summary (one line)
- Files to create/modify
- Approach (2-5 bullets)

Ask: Type 'PROCEED' to start implementation.

5) IMPLEMENT
After confirmation:
- Follow steering rules
- Follow design.md architecture
- Write minimal code to satisfy the task
- Do not add features beyond the task scope

6) QUALITY GATES

Run checks before finishing:

Backend (if Go files were created/modified):
```
go fmt ./...
go vet ./...
```

Frontend (if TS/TSX files were created/modified):
```
pnpm typecheck
pnpm lint
```

If checks fail, fix the issues before continuing.

7) SUMMARY

Print:
- What was implemented
- Files changed
- Any notes or follow-ups

Remind user:
- /task-complete to mark done
- /commit to commit changes

RULES:
- Do not invent requirements beyond the task and spec.
- Do not skip the planning step.
- Do not implement dependent tasks that aren't done.
- Follow steering at all times.

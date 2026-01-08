You are finalizing work and preparing a git commit.

SOURCE OF TRUTH:
- Use `git diff --stat` and `git diff` to understand what changed.
- Do NOT invent intent. Derive it from the diff only.

NOTE: Run quality gates before committing:
- Backend: go fmt ./... && go vet ./...
- Frontend: pnpm typecheck && pnpm lint
If there are errors, fix them before proceeding.

PROCESS:

1) CHANGE ANALYSIS
Run:
- git diff --stat
- git diff

From the diff, determine:
- Primary intent of the change (e.g. refactor, new feature, bug fix)
- Key areas/modules affected
- Scope (small/medium/large)

2) COMMIT MESSAGE GENERATION
Generate a commit message in this format:

<type>: <concise summary>

<optional body>
- bullet points describing notable changes
- derived strictly from the diff

Rules:
- Use conventional types: feat, fix, refactor, chore, test, docs
- Keep the summary under 72 characters
- Do NOT reference tasks, specs, or assumptions unless visible in the diff
- Commits must be atomic (one concern per commit)

3) CONFIRMATION
Print:
- git diff --stat
- the proposed commit message

Then ask:
Type EXACTLY 'CONFIRM COMMIT' to proceed.

4) COMMIT
Only after confirmation:
- git add -A
- git commit -m "<generated message>" (include body if present)

SAFETY RULES:
- Never commit without explicit confirmation.
- Never modify files other than via git add/commit.
- If the working tree is clean, STOP and say so.

Perform a complete code review of the current codebase.

SCOPE:
- Review the entire repo, focusing on Go backend, React frontend, tests, build config, and tooling.
- Do not modify files unless the user explicitly asks for fixes.

PROCESS:

1) REPO INVENTORY
- Identify project structure (docker-compose, backend/, frontend/)
- Identify entry points and main binaries
- Identify CI/lint config if present

2) CORRECTNESS AND RELIABILITY
- Error handling patterns
- Concurrency safety (goroutines, channels, mutexes, context)
- Resource management (files, connections, defer correctness)
- Input validation and edge cases
- React: proper hook usage, effect cleanup, state management

3) BEST PRACTICES
- Go: package boundaries, import hygiene, naming, doc comments, simplicity
- React: component structure, prop types, accessibility
- Avoiding anti-patterns (global state, init complexity, excessive interfaces)

4) SECURITY
- No secrets committed
- Input validation required
- Auth boundaries explicit
- Least privilege everywhere
- XSS/injection prevention in frontend

5) PERFORMANCE (only where relevant)
- Obvious hot paths, unnecessary allocations, inefficient IO
- React: unnecessary re-renders, missing memoization
- Avoid premature micro-optimizations; focus on clear wins

6) TESTING AND QUALITY
- Test coverage of critical paths
- Table-driven tests, test helpers
- Determinism, race potential, missing negative tests
- Frontend: component tests, integration tests

7) TOOLING ALIGNMENT
- Ensure formatting is compliant (gofmt, prettier)
- Ensure linter findings are anticipated
- If hooks are configured, verify they match repo reality

OUTPUT FORMAT:
- Start with a short overall assessment (5–10 lines).
- Then provide findings grouped by severity:
  - Blockers (must fix)
  - Major
  - Minor
  - Suggestions
For each finding:
- Quote the file path and symbol name (function/type/component) where applicable.
- Explain why it matters.
- Provide an actionable recommendation.
If you cannot locate the relevant code, say so explicitly.

RULES:
- Do not invent project requirements.
- Reference steering/spec files for intent: .kiro/steering/*, .kiro/specs/*
- Do not modify files.

You are generating or updating the project README.md after a phase is complete.

INPUT SOURCES (MUST READ):
1) ./The plan.md - project purpose and scope
2) .kiro/specs/*/requirements.md - what was built
3) .kiro/specs/*/tasks.md - verify phase completion
4) Existing ./README.md (if present) - preserve and extend
5) Current codebase - actual implementation details

PROCESS:

1) DETERMINE PHASE STATUS
- Check .kiro/specs/*/tasks.md
- Count open vs done tasks per spec
- Identify which phases are complete (all tasks marked [x] or (done))
- Print: "Phases complete: <list>"

2) GATHER IMPLEMENTATION DETAILS
From the actual codebase, extract:
- Project structure (folders, key files)
- Tech stack in use (versions from package.json, go.mod, docker-compose)
- Available scripts/commands
- API endpoints (if any)
- Environment variables needed
- How to run locally

3) DRAFT README CONTENT

Structure:
```
# <Project Name>

<One paragraph description from The plan.md>

## Status

<Current phase, what's implemented, what's next>

## Tech Stack

<Actual versions from codebase>

## Prerequisites

<What needs to be installed>

## Getting Started

<Step-by-step setup instructions>

## Development

<How to run, test, lint>

## Project Structure

<Folder layout with descriptions>

## API (if applicable)

<Endpoints, methods, brief descriptions>

## Environment Variables

<Required env vars with descriptions>

## Contributing

<Reference to AGENTS.md and steering>

## License

<If specified>
```

4) CONFIRMATION GATE
- Print whether README.md exists (will create vs will update)
- Print a 10-line preview of key sections
- Ask: Type 'CONFIRM WRITE README' to proceed.

5) WRITE README
After confirmation:
- Write ./README.md
- Print: "README.md updated for Phase X completion."

RULES:
- Derive content from actual code, not assumptions
- Keep instructions accurate and testable
- Update status section to reflect current phase
- Preserve any custom sections from existing README
- Use pnpm (not npm) for frontend commands

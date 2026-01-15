# AGENTS.md

Instructions for AI coding agents working on this repository.

## Project Overview

BetterKiroPrompts — A tool that generates Kiro configurations (prompts, steering files, hooks) to help beginners think better, not write code for them.

Stack: Go 1.25.5 backend + React 19 frontend + PostgreSQL 18.1, all in Docker Compose.

## Setup Commands

```bash
# Start everything (stop, rebuild, start)
./build.sh

# Restart without rebuilding
./build.sh --restart

# Stop all containers
./build.sh --stop

# Backend runs on :8090, serves API and built frontend
```

## Code Style

### Go (backend/)
- Standard library preferred over dependencies
- Entry point: `backend/cmd/server/main.go`
- Handlers: `backend/internal/api/`
- Migrations: `backend/migrations/`
- No `pkg/` folder unless shared externally
- Run before commit: `golangci-lint run`, `go fmt ./...`, `go vet ./...`

### TypeScript/React (frontend/)
- Strict TypeScript mode
- shadcn/ui components in `frontend/src/components/ui/`
- App components in `frontend/src/components/`
- Package manager: pnpm
- Run before commit: `pnpm typecheck` and `pnpm lint`

## Testing

```bash
# Backend tests
cd backend && go test ./...

# Frontend tests
cd frontend && pnpm test
```

## Build & Verify

```bash
# Backend
cd backend && go build -o server ./cmd/server

# Frontend
cd frontend && pnpm build
```

## Commit Guidelines

- Use conventional commits: `feat:`, `fix:`, `docs:`, `chore:`, `test:`, `refactor:`
- One concern per commit
- Run quality gates before committing
- Update docs when behavior changes

## Project Structure

```
/
├── build.sh                 # Stack operations script
├── docker-compose.yml       # Dev environment
├── config.toml              # App settings (committed)
├── .env                     # Secrets (not committed)
├── backend/
│   ├── cmd/server/main.go   # Entry point
│   ├── internal/            # All packages (api, db, gallery, generation, etc.)
│   └── migrations/          # SQL migrations
├── frontend/
│   ├── src/components/      # React components
│   └── src/components/ui/   # shadcn/ui components
└── docs/                    # API, developer, user guides
```

## Architecture Rules

- Backend is stateless, serves `/api/*` JSON endpoints + built frontend
- All schema changes require migrations
- No premature scaling or unnecessary abstractions
- Secrets in `.env`, settings in `config.toml`

## Security Considerations

- Never hardcode secrets — use environment variables
- Validate and sanitize all user input
- Scanner container runs isolated with no outbound network during scans
- Rate limiting enabled per IP for public endpoints

## What NOT To Do

- Don't generate application code unless explicitly asked
- Don't invent requirements — ask if input is missing
- Don't skip quality gates
- Don't create dumping ground folders (`utils/`, `helpers/`, `misc/`)

---
inclusion: always
---

# Tech Stack

## Pinned Versions

- Go: 1.25.5
- PostgreSQL: 18.1
- Node.js: 24.12.0
- React: 19.2.3
- Vite: 7.3.1
- shadcn/ui: 0.9.5
- Docker Engine: 29.1.3
- Docker Compose: 5.0.1

## Architecture Rules

- Backend: Go serving `/api/*` JSON API + built static frontend
- Frontend: React + Vite + shadcn/ui (dark theme, blue base)
- Database: PostgreSQL in Docker Compose stack
- Backend must be stateless
- All schema changes require migrations

## Simplicity Rules

- No premature scaling
- No unnecessary abstractions
- Standard library preferred over dependencies
- Package manager: pnpm for frontend

## Code quality

- Frontend always run pnpm typecheck && pnpm lint and build the project after each major task.
- Backend always run golangci-lint && go fmt && go vet and build the project after each major task

## Build & Run

Use `./build.sh` for all stack operations:

```bash
./build.sh up                    # Start dev stack (hot reload)
./build.sh --dev -d up           # Start dev in background
./build.sh --prod --build -d up  # Build and start production
./build.sh stop                  # Stop stack
./build.sh down                  # Stop and remove containers
./build.sh logs                  # Follow logs
./build.sh status                # Show container status
./build.sh clean                 # Remove build artifacts
```

## MCP Servers

Available MCP servers for development assistance:

- `godoc` - Go documentation lookup (use for Go stdlib and package docs)
- `shadcn` - shadcn/ui component search, examples, and installation commands

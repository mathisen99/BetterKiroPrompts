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
./build.sh           # Stop, rebuild everything, and start
./build.sh --restart # Restart without rebuilding
./build.sh --stop    # Stop all containers
```

## MCP Servers

Available MCP servers for development assistance:

- `godoc` - Go documentation lookup (use for Go stdlib and package docs)
- `shadcn` - shadcn/ui component search, examples, and installation commands

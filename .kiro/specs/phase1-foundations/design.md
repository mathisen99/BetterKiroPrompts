# Phase 1: Foundations — Design

## Context

This phase creates the skeleton for the Kiro Prompting & Guardrails Generator. No feature logic—only infrastructure and project structure.

## Architecture Overview

```mermaid
graph TB
    subgraph Docker Compose
        FE[Frontend<br/>React+Vite<br/>:5173]
        BE[Backend<br/>Go<br/>:8080]
        DB[(PostgreSQL<br/>:5432)]
    end
    
    Browser --> FE
    FE -->|/api/*| BE
    BE --> DB
```

## Docker Compose Services

| Service | Image/Build | Ports | Purpose |
|---------|-------------|-------|---------|
| backend | `./backend` | 8080:8080 | Go API server |
| frontend | `./frontend` | 5173:5173 | Vite dev server |
| postgres | `postgres:18.1` | 5432:5432 | Database |

### Volumes
- `postgres_data`: persistent DB storage
- `./backend:/app`: Go source mount for rebuild
- `./frontend:/app`: React source mount for hot reload

### Networks
- `app-network`: internal bridge for service communication

## Go Backend Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   └── api/
│       ├── router.go        # HTTP router setup
│       └── health.go        # Health handler
├── migrations/              # SQL migrations (placeholder)
├── go.mod
├── go.sum
└── Dockerfile
```

### Pinned Versions
- Go: 1.25.5
- Router: standard library `net/http`

## React Frontend Structure

```
frontend/
├── src/
│   ├── App.tsx              # Root component
│   ├── main.tsx             # Entry point
│   ├── index.css            # Global styles + shadcn
│   └── components/
│       └── ui/              # shadcn/ui components
├── index.html
├── vite.config.ts
├── tailwind.config.js
├── package.json
├── pnpm-lock.yaml
└── Dockerfile
```

### Pinned Versions
- Node.js: 24.12.0
- React: 19.2.3
- Vite: 7.3.1
- shadcn/ui: 0.9.5

### Theme
- Mode: dark
- Base color: blue

## Database

- PostgreSQL 18.1
- Connection via environment variables
- Migrations in `backend/migrations/` (applied on startup)

## API Structure

| Endpoint | Method | Response |
|----------|--------|----------|
| `/api/health` | GET | `{"status":"ok"}` |

## Environment Variables

```
# Backend (set in .env)
DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable
POSTGRES_USER=bkp_user
POSTGRES_PASSWORD=CHANGE_ME_IN_PRODUCTION
POSTGRES_DB=betterkiro
PORT=8080

# Frontend
VITE_API_URL=http://localhost:8080
```

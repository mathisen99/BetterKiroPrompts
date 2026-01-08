---
inclusion: always
---

# Repository Structure

## Layout

```
/
├── build.sh
├── docker-compose.yml
├── docker-compose.prod.yml
├── .env.example
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/api/
│   ├── migrations/
│   ├── Dockerfile
│   └── Dockerfile.prod
├── frontend/
│   ├── src/
│   │   ├── components/ui/
│   │   └── App.tsx
│   └── Dockerfile
└── .kiro/
    ├── steering/
    ├── specs/
    └── hooks/
```

## Backend Conventions

- Entry point: `backend/cmd/server/main.go`
- API handlers: `backend/internal/api/`
- DB migrations: `backend/migrations/`
- No `pkg/` unless shared externally

## Frontend Conventions

- shadcn/ui components: `frontend/src/components/ui/`
- App components: `frontend/src/components/`
- Routes: `frontend/src/routes/` (when needed)

## Rules

- No dumping ground folders (`utils/`, `helpers/`, `misc/`)
- One concern per file
- Keep flat until nesting is necessary
- Use MCP server for ShadCN UI components
- Use MCP server for Godoc for go syntax when needed 
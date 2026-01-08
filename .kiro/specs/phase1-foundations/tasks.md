# Phase 1: Foundations — Tasks

## Task List

### Docker Compose Setup

- [x] 1. Create `docker-compose.yml` with backend, frontend, postgres services
  - Refs: AC-1
  - Outcome: `docker compose up` starts all services

- [x] 2. Create `.env.example` with required environment variables
  - Refs: AC-4
  - Outcome: Documented env vars for local dev

### Go Backend

- [x] 3. Initialize Go module in `backend/`
  - Refs: AC-2
  - Outcome: `go.mod` with Go 1.25.5

- [x] 4. Create `backend/cmd/server/main.go` entry point
  - Refs: AC-2
  - Outcome: Server starts on PORT from env

- [x] 5. Create `backend/internal/api/router.go` with route setup
  - Refs: AC-2
  - Outcome: Router mounts `/api/*` routes

- [x] 6. Create `backend/internal/api/health.go` handler
  - Refs: AC-2
  - Outcome: GET `/api/health` returns `{"status":"ok"}`

- [x] 7. Create `backend/Dockerfile`
  - Refs: AC-1, AC-7
  - Outcome: Multi-stage build, hot reload support

- [x] 8. Add database connection in `backend/internal/db/`
  - Refs: AC-4
  - Outcome: Connection pool, logs success on startup

- [x] 9. Create `backend/migrations/` directory with README
  - Refs: AC-5
  - Outcome: Migration structure ready

### React Frontend

- [x] 10. Initialize Vite + React project in `frontend/`
  - Refs: AC-3, AC-6
  - Outcome: `pnpm create vite` with React template

- [x] 11. Install and configure shadcn/ui
  - Refs: AC-3
  - Outcome: Dark theme, blue base color

- [ ] 12. Create `frontend/Dockerfile`
  - Refs: AC-1, AC-6
  - Outcome: Dev server with hot reload

- [ ] 13. Update `frontend/src/App.tsx` with basic layout
  - Refs: AC-3
  - Outcome: Displays app title with shadcn styling

### Verification

- [ ] 14. Test: `docker compose up` starts all services
  - Refs: AC-1
  - Outcome: No errors, all containers healthy

- [ ] 15. Test: Health endpoint returns 200
  - Refs: AC-2
  - Outcome: `curl localhost:8080/api/health` succeeds

- [ ] 16. Test: Frontend loads in browser
  - Refs: AC-3
  - Outcome: Dark theme visible at localhost:5173

- [ ] 17. Test: Database connection logged
  - Refs: AC-4
  - Outcome: Backend logs show DB connected

- [ ] 18. Test: Hot reload works
  - Refs: AC-6, AC-7
  - Outcome: Changes reflect without manual restart

# Phase 1: Foundations — Requirements

## Introduction

Phase 1 establishes the development environment and project skeleton:

- Docker Compose orchestration for all services
- Go backend with health endpoint serving `/api/*`
- React + Vite frontend with shadcn/ui (dark theme, blue base)
- PostgreSQL database with connection and migration structure
- End-to-end development workflow (hot reload, logs, rebuild)
- No feature logic—foundations only

## User Stories

### US-1: Developer Environment Setup
As a developer cloning the repo, I want to run a single command to start all services so that I can begin development immediately.

### US-2: Service Health Verification
As a developer, I want to verify all services are running correctly so that I can trust the environment before writing code.

### US-3: Frontend Development
As a frontend developer, I want hot reload working so that I see changes without restarting containers.

## Acceptance Criteria (EARS Format)

### AC-1: Docker Compose Startup
WHEN a developer runs `docker compose up`
THE SYSTEM SHALL start backend, frontend, and postgres services without errors.

### AC-2: Health Endpoint
WHEN a client sends GET to `/api/health`
THE SYSTEM SHALL return HTTP 200 with JSON `{"status":"ok"}`.

### AC-3: Frontend Loads
WHEN a developer opens `http://localhost:5173`
THE SYSTEM SHALL display the React app with shadcn/ui dark theme and blue accent.

### AC-4: Database Connection
WHEN the backend starts
THE SYSTEM SHALL connect to PostgreSQL and log successful connection.

### AC-5: Migration Structure
WHEN migrations exist in `backend/migrations/`
THE SYSTEM SHALL apply them in order on startup.

### AC-6: Hot Reload
WHEN frontend source files change
THE SYSTEM SHALL reload the browser automatically.

### AC-7: Backend Rebuild
WHEN backend source files change
THE SYSTEM SHALL rebuild and restart the Go service.

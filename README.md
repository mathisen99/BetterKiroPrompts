# BetterKiroPrompts

A tool that generates better prompts, steering documents, and Kiro hooks to improve beginner thinking—not write applications for them.

## Prerequisites

- Docker Engine 29.1.3+
- Docker Compose 5.0.1+
- Node.js 24.12.0+ (for local frontend development)
- pnpm (for frontend package management)

## Quick Start

```bash
# Clone and enter directory
git clone <repo-url>
cd Better-Kiro-Prompts

# Copy environment file
cp .env.example .env

# Start all services
./build.sh up
```

The app will be available at:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

## Development Commands

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

## Project Structure

```
/
├── backend/
│   ├── cmd/server/main.go       # Entry point
│   ├── internal/api/            # HTTP handlers
│   ├── internal/generator/      # Prompt/steering/hooks generation
│   ├── internal/templates/      # Embedded templates
│   └── migrations/              # Database migrations
├── frontend/
│   ├── src/components/          # React components
│   ├── src/pages/               # Page components
│   └── src/lib/                 # API client, utilities
├── docs/
│   ├── api.md                   # API documentation
│   └── user-guide.md            # User guide
└── .kiro/
    ├── steering/                # Project steering files
    ├── specs/                   # Feature specifications
    └── hooks/                   # Kiro hooks
```

## Tech Stack

- Backend: Go 1.25.5
- Frontend: React 19.2.3 + Vite 7.3.1 + shadcn/ui
- Database: PostgreSQL 18.1
- Containerization: Docker Compose

## Documentation

- [API Documentation](docs/api.md)
- [User Guide](docs/user-guide.md)

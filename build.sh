#!/bin/bash
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="$ROOT_DIR/frontend"
BACKEND_DIR="$ROOT_DIR/backend"
STATIC_DIR="$BACKEND_DIR/static"

log() { echo -e "${GREEN}[INFO]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

build_all() {
    log "Building frontend..."
    cd "$FRONTEND_DIR"
    pnpm install --frozen-lockfile
    pnpm build
    rm -rf "$STATIC_DIR"
    cp -r "$FRONTEND_DIR/dist" "$STATIC_DIR"
    
    log "Building backend..."
    cd "$BACKEND_DIR"
    CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
    
    # Copy config.example.toml to backend for Docker build
    cp "$ROOT_DIR/config.example.toml" "$BACKEND_DIR/config.example.toml"
    
    log "Building Docker images..."
    cd "$ROOT_DIR"
    docker compose -f docker-compose.prod.yml --profile scan build --no-cache
    
    log "Build complete!"
}

case "${1:-}" in
    "")
        # ./build.sh - stop, rebuild, start
        log "Stopping containers..."
        docker compose -f docker-compose.prod.yml --profile scan down 2>/dev/null || true
        build_all
        log "Starting containers..."
        docker compose -f docker-compose.prod.yml --profile scan up -d
        log "Stack running at http://localhost:8090"
        ;;
    --restart)
        # ./build.sh --restart - just restart without rebuild
        log "Restarting containers..."
        docker compose -f docker-compose.prod.yml --profile scan down 2>/dev/null || true
        docker compose -f docker-compose.prod.yml --profile scan up -d
        log "Stack running at http://localhost:8090"
        ;;
    --stop)
        # ./build.sh --stop - stop everything
        log "Stopping containers..."
        docker compose -f docker-compose.prod.yml --profile scan down
        log "Stopped!"
        ;;
    *)
        echo "Usage: ./build.sh [OPTION]"
        echo ""
        echo "Options:"
        echo "  (none)      Stop, rebuild everything, and start"
        echo "  --restart   Restart without rebuilding"
        echo "  --stop      Stop all containers"
        ;;
esac

#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Directories
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="$ROOT_DIR/frontend"
BACKEND_DIR="$ROOT_DIR/backend"
STATIC_DIR="$BACKEND_DIR/static"

# Default values
ACTION="up"
BUILD_FRONTEND=false
BUILD_BACKEND=false
REBUILD=false
DETACH=false
PROD=false
DEV=true
NO_ARGS=true

usage() {
    echo "Usage: $0 [OPTIONS] [ACTION]"
    echo ""
    echo "Actions:"
    echo "  up          Start the stack (default)"
    echo "  stop        Stop the stack"
    echo "  restart     Restart the stack"
    echo "  down        Stop and remove containers"
    echo "  logs        Show logs (follow mode)"
    echo "  status      Show container status"
    echo "  clean       Remove build artifacts"
    echo ""
    echo "Options:"
    echo "  --dev               Use development compose (default)"
    echo "  --prod              Use production compose file"
    echo "  --build-frontend    Build frontend for production"
    echo "  --build-backend     Build backend binary"
    echo "  --build             Build both frontend and backend"
    echo "  --rebuild           Force rebuild Docker images"
    echo "  -d, --detach        Run containers in background"
    echo "  -h, --help          Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 up                        Start dev stack (hot reload)"
    echo "  $0 --dev -d up               Start dev in background"
    echo "  $0 --prod --build -d up      Build and start production"
    echo "  $0 --rebuild up              Rebuild Docker images and start"
    echo "  $0 stop                      Stop the stack"
    echo "  $0 --prod down               Tear down production stack"
    echo "  $0 logs                      Follow logs"
}

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

build_frontend() {
    log_info "Building frontend for production..."
    cd "$FRONTEND_DIR"
    
    if ! command -v pnpm &> /dev/null; then
        log_error "pnpm not found. Please install pnpm first."
        exit 1
    fi
    
    pnpm install --frozen-lockfile
    pnpm build
    
    # Copy built assets to backend static folder
    log_info "Copying frontend build to backend/static..."
    rm -rf "$STATIC_DIR"
    cp -r "$FRONTEND_DIR/dist" "$STATIC_DIR"
    
    log_info "Frontend build complete."
    cd "$ROOT_DIR"
}

build_backend() {
    log_info "Building backend binary..."
    cd "$BACKEND_DIR"
    
    if ! command -v go &> /dev/null; then
        log_error "Go not found. Please install Go first."
        exit 1
    fi
    
    CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
    
    log_info "Backend build complete."
    cd "$ROOT_DIR"
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    NO_ARGS=false
    case $1 in
        --build-frontend)
            BUILD_FRONTEND=true
            shift
            ;;
        --build-backend)
            BUILD_BACKEND=true
            shift
            ;;
        --build)
            BUILD_FRONTEND=true
            BUILD_BACKEND=true
            shift
            ;;
        --rebuild)
            REBUILD=true
            shift
            ;;
        --prod)
            PROD=true
            DEV=false
            shift
            ;;
        --dev)
            DEV=true
            PROD=false
            shift
            ;;
        -d|--detach)
            DETACH=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        up|stop|restart|down|logs|status|clean)
            ACTION=$1
            shift
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# If no arguments provided, build everything and start production
if [ "$NO_ARGS" = true ]; then
    log_info "No arguments provided - building and starting production stack..."
    BUILD_FRONTEND=true
    BUILD_BACKEND=true
    PROD=true
    DEV=false
fi

# Execute builds if requested
if [ "$BUILD_FRONTEND" = true ]; then
    build_frontend
fi

if [ "$BUILD_BACKEND" = true ]; then
    build_backend
fi

# Docker compose commands
if [ "$PROD" = true ]; then
    COMPOSE_FILE="docker-compose.prod.yml"
    log_info "Using production compose file"
else
    COMPOSE_FILE="docker-compose.yml"
fi
COMPOSE_CMD="docker compose -f $COMPOSE_FILE"

case $ACTION in
    up)
        log_info "Starting stack..."
        ARGS="--build"
        if [ "$DETACH" = true ]; then
            ARGS="$ARGS -d"
        fi
        if [ "$REBUILD" = true ]; then
            log_info "Forcing fresh rebuild (removing old images)..."
            $COMPOSE_CMD build --no-cache
        fi
        $COMPOSE_CMD up $ARGS
        ;;
    stop)
        log_info "Stopping stack..."
        $COMPOSE_CMD stop
        ;;
    restart)
        log_info "Restarting stack..."
        $COMPOSE_CMD restart
        ;;
    down)
        log_info "Stopping and removing containers..."
        $COMPOSE_CMD down
        ;;
    logs)
        log_info "Following logs..."
        $COMPOSE_CMD logs -f
        ;;
    status)
        $COMPOSE_CMD ps
        ;;
    clean)
        log_info "Cleaning build artifacts..."
        rm -rf "$STATIC_DIR"
        rm -rf "$FRONTEND_DIR/dist"
        rm -f "$BACKEND_DIR/server"
        log_info "Clean complete."
        ;;
    *)
        log_error "Unknown action: $ACTION"
        usage
        exit 1
        ;;
esac

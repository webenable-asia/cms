#!/bin/bash

# WebEnable CMS Production Helper Script (Podman)

set -e

COMPOSE_FILE="podman-compose.yml"
PROJECT_NAME="webenable-cms"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}======================================${NC}"
    echo -e "${BLUE} WebEnable CMS Production Environment${NC}"
    echo -e "${BLUE}======================================${NC}"
}

# Check if Podman is running
check_podman() {
    if ! podman info >/dev/null 2>&1; then
        print_error "Podman is not running or not properly configured."
        print_status "Please ensure Podman is installed and running."
        print_status "On macOS: brew install podman && podman machine init && podman machine start"
        print_status "On Linux: sudo apt install podman (or equivalent package manager)"
        exit 1
    fi
}

# Load environment variables
load_env() {
    if [ -f .env ]; then
        print_status "Loading production environment variables..."
        set -a
        source .env
        set +a
    else
        print_warning ".env file not found. Using default values."
    fi
}

# Show help
show_help() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  start         Start all services"
    echo "  stop          Stop all services"
    echo "  restart       Restart all services"
    echo "  logs          Show logs for all services"
    echo "  logs <service> Show logs for specific service"
    echo "  build         Build all services"
    echo "  build <service> Build specific service"
    echo "  status        Show status of all services"
    echo "  clean         Clean up containers and volumes"
    echo "  open          Open application in browser"
    echo "  shell <service> Open shell in service container"
    echo "  help          Show this help message"
    echo ""
    echo "Services:"
    echo "  frontend      Next.js frontend (via nginx)"
    echo "  backend       Go backend API (via nginx)"
    echo "  db            CouchDB database"
    echo "  cache         Valkey cache"
    echo "  nginx         Nginx reverse proxy (port 80)"
}

# Start services
start_services() {
    load_env
    print_status "Starting WebEnable CMS production services..."
    podman compose -p $PROJECT_NAME up -d
    print_status "Services started successfully!"
    print_status "Application: http://localhost"
    print_status "API: http://localhost/api"
}

# Stop services
stop_services() {
    print_status "Stopping WebEnable CMS services..."
    podman compose -p $PROJECT_NAME down
    print_status "Services stopped successfully!"
}

# Restart services
restart_services() {
    load_env
    print_status "Restarting WebEnable CMS services..."
    podman compose -p $PROJECT_NAME restart
    print_status "Services restarted successfully!"
}

# Show logs
show_logs() {
    if [ -z "$1" ]; then
        podman compose -p $PROJECT_NAME logs -f
    else
        podman compose -p $PROJECT_NAME logs -f "$1"
    fi
}

# Build services
build_services() {
    load_env
    if [ -z "$1" ]; then
        print_status "Building all services..."
        podman compose -p $PROJECT_NAME build --no-cache
    else
        print_status "Building $1 service..."
        podman compose -p $PROJECT_NAME build --no-cache "$1"
    fi
    print_status "Build completed successfully!"
}

# Show status
show_status() {
    print_status "WebEnable CMS Services Status:"
    podman compose -p $PROJECT_NAME ps
}

# Clean up
clean_up() {
    print_warning "This will remove all containers and volumes. Are you sure? (y/N)"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        print_status "Cleaning up containers and volumes..."
        podman compose -p $PROJECT_NAME down -v --remove-orphans
        podman system prune -f
        print_status "Cleanup completed!"
    else
        print_status "Cleanup cancelled."
    fi
}

# Open in browser
open_app() {
    print_status "Opening application in browser..."
    case "$(uname -s)" in
        Darwin) open http://localhost ;;
        Linux) xdg-open http://localhost ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*) start http://localhost ;;
        *) print_error "Unable to open browser automatically" ;;
    esac
}

# Open shell in container
open_shell() {
    if [ -z "$1" ]; then
        print_error "Please specify a service name"
        exit 1
    fi
    print_status "Opening shell in $1 container..."
    podman compose -p $PROJECT_NAME exec "$1" /bin/sh
}

# Main script logic
print_header

check_podman

case "${1:-help}" in
    start)
        start_services
        ;;
    stop)
        stop_services
        ;;
    restart)
        restart_services
        ;;
    logs)
        show_logs "$2"
        ;;
    build)
        build_services "$2"
        ;;
    status)
        show_status
        ;;
    clean)
        clean_up
        ;;
    open)
        open_app
        ;;
    shell)
        open_shell "$2"
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac

#!/bin/bash

# WebEnable CMS Development Helper Script

set -e

COMPOSE_FILE="docker-compose.yml"
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
    echo -e "${BLUE} WebEnable CMS Development Environment${NC}"
    echo -e "${BLUE}======================================${NC}"
}

# Check if Docker is running
check_docker() {
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker Desktop."
        exit 1
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
    echo "  frontend      Open frontend in browser"
    echo "  backend       Open backend API in browser"
    echo "  db            Open CouchDB admin in browser"
    echo "  shell <service> Open shell in service container"
    echo "  help          Show this help message"
    echo ""
    echo "Services:"
    echo "  frontend      Next.js frontend (port 3000)"
    echo "  backend       Node.js backend (port 8080)"
    echo "  db            CouchDB database (port 5984)"
    echo "  cache         Valkey cache (port 6379)"
}

# Start services
start_services() {
    print_status "Starting WebEnable CMS services..."
    docker-compose -p $PROJECT_NAME up -d
    print_status "Services started successfully!"
    print_status "Frontend: http://localhost:3000"
    print_status "Backend API: http://localhost:8080"
    print_status "CouchDB Admin: http://localhost:5984/_utils"
}

# Stop services
stop_services() {
    print_status "Stopping WebEnable CMS services..."
    docker-compose -p $PROJECT_NAME down
    print_status "Services stopped successfully!"
}

# Restart services
restart_services() {
    print_status "Restarting WebEnable CMS services..."
    docker-compose -p $PROJECT_NAME restart
    print_status "Services restarted successfully!"
}

# Show logs
show_logs() {
    if [ -z "$1" ]; then
        docker-compose -p $PROJECT_NAME logs -f
    else
        docker-compose -p $PROJECT_NAME logs -f "$1"
    fi
}

# Build services
build_services() {
    if [ -z "$1" ]; then
        print_status "Building all services..."
        docker-compose -p $PROJECT_NAME build
    else
        print_status "Building $1 service..."
        docker-compose -p $PROJECT_NAME build "$1"
    fi
    print_status "Build completed successfully!"
}

# Show status
show_status() {
    print_status "WebEnable CMS Services Status:"
    docker-compose -p $PROJECT_NAME ps
}

# Clean up
clean_up() {
    print_warning "This will remove all containers and volumes. Are you sure? (y/N)"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        print_status "Cleaning up containers and volumes..."
        docker-compose -p $PROJECT_NAME down -v --remove-orphans
        docker system prune -f
        print_status "Cleanup completed!"
    else
        print_status "Cleanup cancelled."
    fi
}

# Open in browser
open_frontend() {
    print_status "Opening frontend in browser..."
    case "$(uname -s)" in
        Darwin) open http://localhost:3000 ;;
        Linux) xdg-open http://localhost:3000 ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*) start http://localhost:3000 ;;
        *) print_error "Unable to open browser automatically" ;;
    esac
}

open_backend() {
    print_status "Opening backend API in browser..."
    case "$(uname -s)" in
        Darwin) open http://localhost:8080 ;;
        Linux) xdg-open http://localhost:8080 ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*) start http://localhost:8080 ;;
        *) print_error "Unable to open browser automatically" ;;
    esac
}

open_db() {
    print_status "Opening CouchDB admin in browser..."
    case "$(uname -s)" in
        Darwin) open http://localhost:5984/_utils ;;
        Linux) xdg-open http://localhost:5984/_utils ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*) start http://localhost:5984/_utils ;;
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
    docker-compose -p $PROJECT_NAME exec "$1" /bin/sh
}

# Main script logic
print_header

check_docker

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
    frontend)
        open_frontend
        ;;
    backend)
        open_backend
        ;;
    db)
        open_db
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

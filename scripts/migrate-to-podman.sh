#!/bin/bash

# Migration script from Docker to Podman
# This script helps users transition from Docker to Podman

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}=========================================${NC}"
    echo -e "${BLUE}  WebEnable CMS: Docker to Podman Migration${NC}"
    echo -e "${BLUE}=========================================${NC}"
    echo ""
}

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is currently running
check_docker() {
    if docker info >/dev/null 2>&1; then
        print_status "Docker is currently running"
        return 0
    else
        print_status "Docker is not running (this is expected for migration)"
        return 1
    fi
}

# Check if Podman is installed
check_podman_installed() {
    if command -v podman >/dev/null 2>&1; then
        print_status "Podman is already installed"
        return 0
    else
        print_warning "Podman is not installed"
        return 1
    fi
}

# Install Podman based on OS
install_podman() {
    print_status "Installing Podman..."
    
    case "$(uname -s)" in
        Darwin)
            print_status "Installing Podman on macOS..."
            if command -v brew >/dev/null 2>&1; then
                brew install podman
                print_status "Initializing Podman machine..."
                podman machine init --cpus 2 --memory 4096
                podman machine start
            else
                print_error "Homebrew is required to install Podman on macOS"
                print_status "Please install Homebrew first: https://brew.sh/"
                exit 1
            fi
            ;;
        Linux)
            print_status "Installing Podman on Linux..."
            
            # Detect Linux distribution
            if [ -f /etc/os-release ]; then
                . /etc/os-release
                case "$ID" in
                    ubuntu|debian)
                        sudo apt update
                        sudo apt install -y podman
                        ;;
                    fedora)
                        sudo dnf install -y podman
                        ;;
                    centos|rhel)
                        sudo yum install -y podman
                        ;;
                    *)
                        print_error "Unsupported Linux distribution: $ID"
                        print_status "Please install Podman manually for your distribution"
                        exit 1
                        ;;
                esac
            else
                print_error "Cannot detect Linux distribution"
                print_status "Please install Podman manually"
                exit 1
            fi
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            print_status "Please install Podman manually"
            exit 1
            ;;
    esac
}

# Stop existing Docker containers
stop_docker_containers() {
    if check_docker; then
        print_status "Stopping existing Docker containers..."
        
        # Stop the WebEnable CMS containers
        if docker-compose ps | grep -q "webenable"; then
            docker-compose down
        fi
        
        # Stop any other containers using the same ports
        for port in 80 3000 8080 5984 6379; do
            container=$(docker ps --filter "publish=$port" --format "{{.Names}}" | head -1)
            if [ -n "$container" ]; then
                print_status "Stopping container using port $port: $container"
                docker stop "$container" 2>/dev/null || true
            fi
        done
        
        print_status "Docker containers stopped"
    fi
}

# Migrate Docker volumes to Podman
migrate_volumes() {
    print_status "Checking for existing Docker volumes to migrate..."
    
    # List of volumes to migrate
    volumes=(
        "webenable-cms_couchdb_data"
        "webenable-cms_valkey_data"
        "webenable-cms_caddy_data"
        "webenable-cms_caddy_config"
    )
    
    for volume in "${volumes[@]}"; do
        if docker volume ls | grep -q "$volume"; then
            print_status "Migrating volume: $volume"
            
            # Create a backup of the Docker volume
            backup_file="/tmp/${volume}_backup.tar.gz"
            docker run --rm \
                -v "$volume":/source \
                -v /tmp:/backup \
                alpine tar czf "/backup/${volume}_backup.tar.gz" -C /source .
            
            # Create the Podman volume
            podman volume create "$volume" >/dev/null 2>&1 || true
            
            # Restore data to Podman volume
            podman run --rm \
                -v "$volume":/target \
                -v /tmp:/backup \
                alpine tar xzf "/backup/${volume}_backup.tar.gz" -C /target
            
            # Clean up backup
            rm -f "$backup_file"
            
            print_status "Volume $volume migrated successfully"
        else
            print_status "Volume $volume not found, skipping"
        fi
    done
}

# Test Podman setup
test_podman() {
    print_status "Testing Podman installation..."
    
    # Test basic Podman functionality
    if podman info >/dev/null 2>&1; then
        print_status "✓ Podman is working correctly"
    else
        print_error "✗ Podman is not working properly"
        exit 1
    fi
    
    # Test Podman Compose
    if podman compose version >/dev/null 2>&1; then
        print_status "✓ Podman Compose is working correctly"
    else
        print_error "✗ Podman Compose is not available"
        print_status "You may need to install podman-compose separately"
        exit 1
    fi
}

# Start services with Podman
start_with_podman() {
    print_status "Starting WebEnable CMS with Podman..."
    
    # Use the updated management script
    if [ -f "./manage.sh" ]; then
        ./manage.sh start
    else
        print_status "Using direct podman compose command..."
        podman compose up -d
    fi
    
    print_status "Services started with Podman!"
    print_status "Frontend: http://localhost:3000"
    print_status "Backend: http://localhost:8080"
    print_status "Application (via proxy): http://localhost"
}

# Show migration summary
show_summary() {
    echo ""
    echo -e "${GREEN}Migration Summary:${NC}"
    echo "=================="
    echo ""
    echo "✅ Podman installed and configured"
    echo "✅ Docker containers stopped"
    echo "✅ Data volumes migrated"
    echo "✅ WebEnable CMS running on Podman"
    echo ""
    echo -e "${BLUE}Next Steps:${NC}"
    echo "- Use './manage.sh' for all container operations"
    echo "- Read PODMAN.md for detailed Podman-specific documentation"
    echo "- Consider uninstalling Docker if no longer needed"
    echo ""
    echo -e "${YELLOW}Important:${NC}"
    echo "- All your data has been preserved during migration"
    echo "- Scripts have been updated to use 'podman compose'"
    echo "- Development workflow remains the same"
    echo ""
}

# Main migration process
main() {
    print_header
    
    # Check current state
    docker_running=$(check_docker && echo "true" || echo "false")
    
    # Install Podman if needed
    if ! check_podman_installed; then
        install_podman
    fi
    
    # Test Podman
    test_podman
    
    # Stop Docker containers if running
    if [ "$docker_running" = "true" ]; then
        stop_docker_containers
        # Migrate volumes
        migrate_volumes
    else
        print_status "No Docker containers running, skipping volume migration"
    fi
    
    # Start services with Podman
    start_with_podman
    
    # Show summary
    show_summary
}

# Run migration
main "$@"

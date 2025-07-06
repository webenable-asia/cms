#!/bin/bash

# Rate Limit Reset Script for WebEnable CMS
# This script provides easy commands to reset rate limits using the API

API_BASE="http://localhost:8080/api"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to get JWT token (you need to implement login)
get_jwt_token() {
    local username="${1:-admin}"
    local password="${2:-admin123}"
    
    echo "Getting JWT token for user: $username" >&2
    
    local response=$(curl -s -X POST "$API_BASE/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$username\",\"password\":\"$password\"}")
    
    local token=$(echo "$response" | jq -r '.token // empty')
    
    if [ -z "$token" ] || [ "$token" = "null" ]; then
        echo -e "${RED}Failed to get JWT token. Please check credentials.${NC}" >&2
        echo "Response: $response" >&2
        exit 1
    fi
    
    echo "$token"
}

# Function to make authenticated API calls
api_call() {
    local method="$1"
    local endpoint="$2"
    local token="$3"
    
    curl -s -X "$method" "$API_BASE$endpoint" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json"
}

# Function to reset rate limit for IP
reset_ip() {
    local ip="$1"
    local token="$2"
    
    if [ -z "$ip" ]; then
        echo -e "${RED}Error: IP address is required${NC}"
        echo "Usage: $0 reset-ip <ip_address>"
        exit 1
    fi
    
    echo -e "${BLUE}Resetting rate limit for IP: $ip${NC}"
    
    local response=$(api_call "POST" "/admin/rate-limit/reset?type=ip&target=$ip" "$token")
    local message=$(echo "$response" | jq -r '.message // .error // empty')
    
    if echo "$response" | jq -e '.message' >/dev/null 2>&1; then
        echo -e "${GREEN}✓ $message${NC}"
    else
        echo -e "${RED}✗ Failed: $message${NC}"
        exit 1
    fi
}

# Function to reset rate limit for user
reset_user() {
    local user_id="$1"
    local token="$2"
    
    if [ -z "$user_id" ]; then
        echo -e "${RED}Error: User ID is required${NC}"
        echo "Usage: $0 reset-user <user_id>"
        exit 1
    fi
    
    echo -e "${BLUE}Resetting rate limit for user: $user_id${NC}"
    
    local response=$(api_call "POST" "/admin/rate-limit/reset?type=user&target=$user_id" "$token")
    local message=$(echo "$response" | jq -r '.message // .error // empty')
    
    if echo "$response" | jq -e '.message' >/dev/null 2>&1; then
        echo -e "${GREEN}✓ $message${NC}"
    else
        echo -e "${RED}✗ Failed: $message${NC}"
        exit 1
    fi
}

# Function to reset all rate limits
reset_all() {
    local type="$1"
    local token="$2"
    
    echo -e "${BLUE}Resetting $type rate limits...${NC}"
    
    local response=$(api_call "POST" "/admin/rate-limit/reset?type=$type" "$token")
    local message=$(echo "$response" | jq -r '.message // .error // empty')
    
    if echo "$response" | jq -e '.message' >/dev/null 2>&1; then
        echo -e "${GREEN}✓ $message${NC}"
    else
        echo -e "${RED}✗ Failed: $message${NC}"
        exit 1
    fi
}

# Function to check rate limit status
check_status() {
    local type="$1"
    local target="$2"
    local token="$3"
    
    if [ -z "$type" ] || [ -z "$target" ]; then
        echo -e "${RED}Error: Type and target are required${NC}"
        echo "Usage: $0 status <ip|user> <target>"
        exit 1
    fi
    
    echo -e "${BLUE}Checking rate limit status for $type: $target${NC}"
    
    local response=$(api_call "GET" "/admin/rate-limit/status?type=$type&target=$target" "$token")
    
    if echo "$response" | jq -e '.api_limit or .auth_limit or .user_limit' >/dev/null 2>&1; then
        echo "$response" | jq .
    else
        local message=$(echo "$response" | jq -r '.message // .error // "No data found"')
        echo -e "${YELLOW}$message${NC}"
    fi
}

# Function to show help
show_help() {
    cat << EOF
${BLUE}WebEnable CMS Rate Limit Management${NC}

${YELLOW}USAGE:${NC}
    $0 <command> [options]

${YELLOW}COMMANDS:${NC}
    reset-ip <ip_address>           Reset rate limit for specific IP
    reset-user <user_id>            Reset rate limit for specific user
    reset-api                       Reset all API rate limits
    reset-auth                      Reset all authentication rate limits
    reset-users                     Reset all user-specific rate limits
    reset-all                       Reset ALL rate limits
    status <ip|user> <target>       Check rate limit status
    help                            Show this help message

${YELLOW}EXAMPLES:${NC}
    $0 reset-ip 192.168.1.100
    $0 reset-user user123
    $0 reset-all
    $0 status ip 192.168.1.100
    $0 status user user123

${YELLOW}AUTHENTICATION:${NC}
    Set environment variables:
    export CMS_USERNAME=admin
    export CMS_PASSWORD=your_password

    Or the script will prompt for credentials.

${YELLOW}REQUIREMENTS:${NC}
    - curl
    - jq
    - Backend server running on localhost:8080
EOF
}

# Main script logic
main() {
    local command="$1"
    
    # Check dependencies
    if ! command -v curl >/dev/null 2>&1; then
        echo -e "${RED}Error: curl is required but not installed${NC}"
        exit 1
    fi
    
    if ! command -v jq >/dev/null 2>&1; then
        echo -e "${RED}Error: jq is required but not installed${NC}"
        exit 1
    fi
    
    case "$command" in
        "reset-ip")
            local token=$(get_jwt_token "${CMS_USERNAME}" "${CMS_PASSWORD}")
            reset_ip "$2" "$token"
            ;;
        "reset-user")
            local token=$(get_jwt_token "${CMS_USERNAME}" "${CMS_PASSWORD}")
            reset_user "$2" "$token"
            ;;
        "reset-api")
            local token=$(get_jwt_token "${CMS_USERNAME}" "${CMS_PASSWORD}")
            reset_all "api" "$token"
            ;;
        "reset-auth")
            local token=$(get_jwt_token "${CMS_USERNAME}" "${CMS_PASSWORD}")
            reset_all "auth" "$token"
            ;;
        "reset-users")
            local token=$(get_jwt_token "${CMS_USERNAME}" "${CMS_PASSWORD}")
            reset_all "users" "$token"
            ;;
        "reset-all")
            local token=$(get_jwt_token "${CMS_USERNAME}" "${CMS_PASSWORD}")
            reset_all "all" "$token"
            ;;
        "status")
            local token=$(get_jwt_token "${CMS_USERNAME}" "${CMS_PASSWORD}")
            check_status "$2" "$3" "$token"
            ;;
        "help"|"--help"|"-h"|"")
            show_help
            ;;
        *)
            echo -e "${RED}Unknown command: $command${NC}"
            echo "Use '$0 help' for usage information."
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"

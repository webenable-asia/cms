#!/bin/bash

# Direct Valkey Rate Limit Reset Script
# This script connects directly to Valkey cache to reset rate limits
# Bypasses API authentication for emergency situations

VALKEY_HOST="${VALKEY_HOST:-localhost}"
VALKEY_PORT="${VALKEY_PORT:-6379}"
VALKEY_PASSWORD="${VALKEY_PASSWORD:-valkeypassword}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to execute Redis/Valkey commands
valkey_exec() {
    local command="$1"
    
    if [ -n "$VALKEY_PASSWORD" ]; then
        echo "$command" | redis-cli -h "$VALKEY_HOST" -p "$VALKEY_PORT" -a "$VALKEY_PASSWORD" --no-auth-warning
    else
        echo "$command" | redis-cli -h "$VALKEY_HOST" -p "$VALKEY_PORT"
    fi
}

# Function to reset rate limits by pattern
reset_by_pattern() {
    local pattern="$1"
    local description="$2"
    
    echo -e "${BLUE}Resetting $description...${NC}"
    
    # Get all keys matching the pattern
    local keys=$(valkey_exec "KEYS rate_limit:$pattern")
    
    if [ -z "$keys" ] || [ "$keys" = "(empty array)" ]; then
        echo -e "${YELLOW}No rate limit keys found for pattern: $pattern${NC}"
        return 0
    fi
    
    # Count keys
    local key_count=$(echo "$keys" | wc -l)
    
    # Delete all matching keys
    echo "$keys" | while read -r key; do
        if [ -n "$key" ]; then
            valkey_exec "DEL $key" >/dev/null
        fi
    done
    
    echo -e "${GREEN}✓ Reset $key_count rate limit entries for $description${NC}"
}

# Function to reset specific IP
reset_ip() {
    local ip="$1"
    
    if [ -z "$ip" ]; then
        echo -e "${RED}Error: IP address is required${NC}"
        echo "Usage: $0 ip <ip_address>"
        exit 1
    fi
    
    echo -e "${BLUE}Resetting rate limits for IP: $ip${NC}"
    
    local api_key="rate_limit:api:$ip"
    local auth_key="rate_limit:auth:$ip"
    
    local api_result=$(valkey_exec "DEL $api_key")
    local auth_result=$(valkey_exec "DEL $auth_key")
    
    local total_reset=$((api_result + auth_result))
    
    if [ "$total_reset" -gt 0 ]; then
        echo -e "${GREEN}✓ Reset $total_reset rate limit entries for IP: $ip${NC}"
    else
        echo -e "${YELLOW}No rate limit entries found for IP: $ip${NC}"
    fi
}

# Function to reset specific user
reset_user() {
    local user_id="$1"
    
    if [ -z "$user_id" ]; then
        echo -e "${RED}Error: User ID is required${NC}"
        echo "Usage: $0 user <user_id>"
        exit 1
    fi
    
    echo -e "${BLUE}Resetting rate limits for user: $user_id${NC}"
    
    local user_key="rate_limit:user:$user_id"
    local result=$(valkey_exec "DEL $user_key")
    
    if [ "$result" -gt 0 ]; then
        echo -e "${GREEN}✓ Reset rate limit for user: $user_id${NC}"
    else
        echo -e "${YELLOW}No rate limit entries found for user: $user_id${NC}"
    fi
}

# Function to show current rate limit status
show_status() {
    echo -e "${BLUE}Current Rate Limit Status:${NC}"
    echo ""
    
    # Count different types of rate limits
    local api_count=$(valkey_exec "KEYS rate_limit:api:*" | wc -l)
    local auth_count=$(valkey_exec "KEYS rate_limit:auth:*" | wc -l)
    local user_count=$(valkey_exec "KEYS rate_limit:user:*" | wc -l)
    local total_count=$(valkey_exec "KEYS rate_limit:*" | wc -l)
    
    echo -e "API Rate Limits:    ${YELLOW}$api_count${NC}"
    echo -e "Auth Rate Limits:   ${YELLOW}$auth_count${NC}"
    echo -e "User Rate Limits:   ${YELLOW}$user_count${NC}"
    echo -e "Total Rate Limits:  ${YELLOW}$total_count${NC}"
    echo ""
    
    # Show some example entries
    if [ "$total_count" -gt 0 ]; then
        echo -e "${BLUE}Recent Rate Limit Entries:${NC}"
        valkey_exec "KEYS rate_limit:*" | head -10 | while read -r key; do
            if [ -n "$key" ]; then
                local value=$(valkey_exec "GET $key")
                local ttl=$(valkey_exec "TTL $key")
                echo -e "  $key: ${YELLOW}$value${NC} (TTL: ${ttl}s)"
            fi
        done
    fi
}

# Function to show help
show_help() {
    cat << EOF
${BLUE}WebEnable CMS Direct Valkey Rate Limit Reset${NC}

${YELLOW}USAGE:${NC}
    $0 <command> [options]

${YELLOW}COMMANDS:${NC}
    ip <ip_address>                 Reset rate limits for specific IP
    user <user_id>                  Reset rate limits for specific user
    api                             Reset all API rate limits
    auth                            Reset all authentication rate limits
    users                           Reset all user-specific rate limits
    all                             Reset ALL rate limits
    status                          Show current rate limit status
    help                            Show this help message

${YELLOW}EXAMPLES:${NC}
    $0 ip 192.168.1.100
    $0 user user123
    $0 api
    $0 all
    $0 status

${YELLOW}ENVIRONMENT VARIABLES:${NC}
    VALKEY_HOST=localhost           Valkey server host
    VALKEY_PORT=6379                Valkey server port  
    VALKEY_PASSWORD=valkeypassword  Valkey password

${YELLOW}REQUIREMENTS:${NC}
    - redis-cli
    - Access to Valkey server

${YELLOW}NOTE:${NC}
    This script connects directly to Valkey and bypasses API authentication.
    Use for emergency situations or when API is unavailable.
EOF
}

# Main script logic
main() {
    local command="$1"
    
    # Check dependencies
    if ! command -v redis-cli >/dev/null 2>&1; then
        echo -e "${RED}Error: redis-cli is required but not installed${NC}"
        exit 1
    fi
    
    # Test connection
    local ping_result=$(valkey_exec "PING" 2>/dev/null)
    if [ "$ping_result" != "PONG" ]; then
        echo -e "${RED}Error: Cannot connect to Valkey server at $VALKEY_HOST:$VALKEY_PORT${NC}"
        echo "Check VALKEY_HOST, VALKEY_PORT, and VALKEY_PASSWORD environment variables."
        exit 1
    fi
    
    case "$command" in
        "ip")
            reset_ip "$2"
            ;;
        "user")
            reset_user "$2"
            ;;
        "api")
            reset_by_pattern "api:*" "API rate limits"
            ;;
        "auth")
            reset_by_pattern "auth:*" "authentication rate limits"
            ;;
        "users")
            reset_by_pattern "user:*" "user-specific rate limits"
            ;;
        "all")
            reset_by_pattern "*" "ALL rate limits"
            ;;
        "status")
            show_status
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

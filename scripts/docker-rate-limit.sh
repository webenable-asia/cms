#!/bin/bash

# Docker-based Valkey Rate Limit Reset Script
# Uses the Valkey container to execute commands

COMPOSE_FILE="docker-compose.yml"
VALKEY_SERVICE="cache"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to execute Valkey commands via Docker
valkey_exec() {
    local command="$1"
    docker-compose exec -T "$VALKEY_SERVICE" redis-cli -a valkeypassword --no-auth-warning $command
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
    
    # Delete all matching keys using DEL command
    if [ "$key_count" -gt 0 ]; then
        echo "$keys" | xargs -I {} docker-compose exec -T "$VALKEY_SERVICE" redis-cli -a valkeypassword --no-auth-warning DEL {}
        echo -e "${GREEN}✓ Reset $key_count rate limit entries for $description${NC}"
    fi
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
    
    local api_result=$(valkey_exec "DEL rate_limit:api:$ip")
    local auth_result=$(valkey_exec "DEL rate_limit:auth:$ip")
    
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
    
    local result=$(valkey_exec "DEL rate_limit:user:$user_id")
    
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
    local api_count=$(docker-compose exec -T "$VALKEY_SERVICE" redis-cli -a valkeypassword --no-auth-warning KEYS "rate_limit:api:*" | wc -l)
    local auth_count=$(docker-compose exec -T "$VALKEY_SERVICE" redis-cli -a valkeypassword --no-auth-warning KEYS "rate_limit:auth:*" | wc -l) 
    local user_count=$(docker-compose exec -T "$VALKEY_SERVICE" redis-cli -a valkeypassword --no-auth-warning KEYS "rate_limit:user:*" | wc -l)
    local total_count=$(docker-compose exec -T "$VALKEY_SERVICE" redis-cli -a valkeypassword --no-auth-warning KEYS "rate_limit:*" | wc -l)
    
    echo -e "API Rate Limits:    ${YELLOW}$api_count${NC}"
    echo -e "Auth Rate Limits:   ${YELLOW}$auth_count${NC}"
    echo -e "User Rate Limits:   ${YELLOW}$user_count${NC}"
    echo -e "Total Rate Limits:  ${YELLOW}$total_count${NC}"
    echo ""
    
    # Show some example entries
    if [ "$total_count" -gt 0 ]; then
        echo -e "${BLUE}Recent Rate Limit Entries:${NC}"
        docker-compose exec -T "$VALKEY_SERVICE" redis-cli -a valkeypassword --no-auth-warning KEYS "rate_limit:*" | head -5 | while read -r key; do
            if [ -n "$key" ] && [ "$key" != "(empty array)" ]; then
                local value=$(docker-compose exec -T "$VALKEY_SERVICE" redis-cli -a valkeypassword --no-auth-warning GET "$key")
                local ttl=$(docker-compose exec -T "$VALKEY_SERVICE" redis-cli -a valkeypassword --no-auth-warning TTL "$key")
                echo -e "  $key: ${YELLOW}$value${NC} (TTL: ${ttl}s)"
            fi
        done
    fi
}

# Function to show help
show_help() {
    cat << EOF
${BLUE}WebEnable CMS Docker Valkey Rate Limit Reset${NC}

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

${YELLOW}REQUIREMENTS:${NC}
    - Docker and docker-compose
    - Valkey service running in Docker

${YELLOW}NOTE:${NC}
    This script uses docker-compose to connect to the Valkey container.
    Make sure you're running this from the project root directory.
EOF
}

# Main script logic
main() {
    local command="$1"
    
    # Check if docker-compose is available
    if ! command -v docker-compose >/dev/null 2>&1; then
        echo -e "${RED}Error: docker-compose is required but not installed${NC}"
        exit 1
    fi
    
    # Check if Valkey service is running
    if ! docker-compose ps "$VALKEY_SERVICE" | grep -q "Up"; then
        echo -e "${RED}Error: Valkey service ($VALKEY_SERVICE) is not running${NC}"
        echo "Start it with: docker-compose up -d $VALKEY_SERVICE"
        exit 1
    fi
    
    # Test connection
    local ping_result=$(docker-compose exec -T "$VALKEY_SERVICE" redis-cli -a valkeypassword --no-auth-warning PING 2>/dev/null)
    if [ "$ping_result" != "PONG" ]; then
        echo -e "${RED}Error: Cannot connect to Valkey service${NC}"
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

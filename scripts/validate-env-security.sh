#!/bin/bash

# Environment Security Validation Script
# Checks production environment configuration for security best practices

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

ENV_FILE="${1:-.env}"
SCORE=0
MAX_SCORE=0

print_header() {
    echo -e "${BLUE}=================================================${NC}"
    echo -e "${BLUE}  WebEnable CMS Environment Security Check${NC}"
    echo -e "${BLUE}=================================================${NC}"
    echo ""
    echo "Checking: $ENV_FILE"
    echo ""
}

check_passed() {
    echo -e "${GREEN}✓${NC} $1"
    SCORE=$((SCORE + 1))
}

check_failed() {
    echo -e "${RED}✗${NC} $1"
}

check_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

increment_max() {
    MAX_SCORE=$((MAX_SCORE + 1))
}

# Check if environment file exists
check_env_file() {
    increment_max
    if [ -f "$ENV_FILE" ]; then
        check_passed "Environment file exists"
    else
        check_failed "Environment file not found: $ENV_FILE"
        exit 1
    fi
}

# Check file permissions
check_permissions() {
    increment_max
    local perms=$(stat -c "%a" "$ENV_FILE" 2>/dev/null || stat -f "%Lp" "$ENV_FILE" 2>/dev/null || echo "unknown")
    
    if [ "$perms" = "600" ]; then
        check_passed "Secure file permissions (600)"
    elif [ "$perms" = "unknown" ]; then
        check_warning "Could not determine file permissions"
    else
        check_failed "Insecure file permissions ($perms). Should be 600. Run: chmod 600 $ENV_FILE"
    fi
}

# Check JWT secret strength
check_jwt_secret() {
    increment_max
    local jwt_secret=$(grep "^JWT_SECRET=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if [ ${#jwt_secret} -ge 32 ]; then
        check_passed "JWT secret has sufficient length (${#jwt_secret} characters)"
    else
        check_failed "JWT secret too short (${#jwt_secret} characters). Should be at least 32 characters."
    fi
    
    increment_max
    if echo "$jwt_secret" | grep -q "your-.*-secret\|change.*this\|example"; then
        check_failed "JWT secret appears to be a placeholder. Generate new: openssl rand -base64 32"
    else
        check_passed "JWT secret is not a default placeholder"
    fi
}

# Check database password
check_database_password() {
    increment_max
    local db_pass=$(grep "^COUCHDB_PASSWORD=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if [ ${#db_pass} -ge 16 ]; then
        check_passed "Database password has sufficient length (${#db_pass} characters)"
    else
        check_failed "Database password too short (${#db_pass} characters). Should be at least 16 characters."
    fi
    
    increment_max
    if echo "$db_pass" | grep -q "password\|admin\|123\|qwerty"; then
        check_failed "Database password appears to be weak or default"
    else
        check_passed "Database password is not a common weak password"
    fi
}

# Check cache password
check_cache_password() {
    increment_max
    local cache_pass=$(grep "^VALKEY_PASSWORD=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if [ ${#cache_pass} -ge 16 ]; then
        check_passed "Cache password has sufficient length (${#cache_pass} characters)"
    else
        check_failed "Cache password too short (${#cache_pass} characters). Should be at least 16 characters."
    fi
}

# Check admin credentials
check_admin_credentials() {
    increment_max
    local admin_user=$(grep "^ADMIN_USERNAME=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if [ "$admin_user" != "admin" ] && [ "$admin_user" != "root" ] && [ "$admin_user" != "administrator" ]; then
        check_passed "Admin username is not a default value"
    else
        check_failed "Admin username is a default/common value. Change to something unique."
    fi
    
    increment_max
    local admin_pass=$(grep "^ADMIN_PASSWORD=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if [ ${#admin_pass} -ge 12 ]; then
        check_passed "Admin password has sufficient length (${#admin_pass} characters)"
    else
        check_failed "Admin password too short (${#admin_pass} characters). Should be at least 12 characters."
    fi
}

# Check domain configuration
check_domain_config() {
    increment_max
    local cors_origins=$(grep "^CORS_ORIGINS=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if echo "$cors_origins" | grep -q "localhost\|yourdomain.com\|example.com"; then
        check_failed "CORS origins contains placeholder domains. Update with your actual domain."
    else
        check_passed "CORS origins configured with actual domains"
    fi
    
    increment_max
    local api_url=$(grep "^NEXT_PUBLIC_API_URL=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if echo "$api_url" | grep -q "https://"; then
        check_passed "API URL uses HTTPS"
    else
        check_failed "API URL should use HTTPS in production"
    fi
}

# Check session security
check_session_security() {
    increment_max
    local session_secure=$(grep "^SESSION_SECURE=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if [ "$session_secure" = "true" ]; then
        check_passed "Session security is enabled"
    else
        check_failed "SESSION_SECURE should be 'true' in production"
    fi
    
    increment_max
    local session_domain=$(grep "^SESSION_DOMAIN=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if [ "$session_domain" != "localhost" ] && [ -n "$session_domain" ]; then
        check_passed "Session domain is configured for production"
    else
        check_failed "SESSION_DOMAIN should be set to your production domain"
    fi
}

# Check environment type
check_environment_type() {
    increment_max
    local node_env=$(grep "^NODE_ENV=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if [ "$node_env" = "production" ]; then
        check_passed "NODE_ENV is set to production"
    else
        check_failed "NODE_ENV should be 'production' for production deployments"
    fi
    
    increment_max
    local go_env=$(grep "^GO_ENV=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if [ "$go_env" = "production" ]; then
        check_passed "GO_ENV is set to production"
    else
        check_warning "GO_ENV should be 'production' for production deployments"
    fi
}

# Check email configuration
check_email_config() {
    increment_max
    local smtp_host=$(grep "^SMTP_HOST=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"')
    
    if [ -n "$smtp_host" ] && [ "$smtp_host" != "localhost" ] && [ "$smtp_host" != "smtp.example.com" ]; then
        check_passed "SMTP host is configured"
    else
        check_warning "SMTP host not configured or using placeholder"
    fi
}

# Check for sensitive information exposure
check_sensitive_exposure() {
    increment_max
    if git check-ignore "$ENV_FILE" >/dev/null 2>&1; then
        check_passed "Environment file is in .gitignore"
    else
        check_failed "Environment file should be added to .gitignore to prevent accidental commits"
    fi
}

# Generate security report
generate_report() {
    echo ""
    echo -e "${BLUE}Security Assessment Results${NC}"
    echo "=================================="
    
    local percentage=$((SCORE * 100 / MAX_SCORE))
    
    echo "Score: $SCORE/$MAX_SCORE ($percentage%)"
    
    if [ $percentage -ge 90 ]; then
        echo -e "${GREEN}Excellent! Your environment is well-secured.${NC}"
    elif [ $percentage -ge 75 ]; then
        echo -e "${YELLOW}Good security posture. Address the failed checks above.${NC}"
    elif [ $percentage -ge 60 ]; then
        echo -e "${YELLOW}Moderate security. Several issues need attention.${NC}"
    else
        echo -e "${RED}Poor security! Critical issues must be fixed before production.${NC}"
    fi
    
    echo ""
    echo "Recommendations:"
    echo "- Generate strong, unique passwords for all services"
    echo "- Use HTTPS for all production URLs"
    echo "- Set secure file permissions (chmod 600)"
    echo "- Never commit .env files to version control"
    echo "- Rotate secrets regularly"
    echo "- Monitor for security updates"
}

# Main validation process
main() {
    print_header
    
    check_env_file
    check_permissions
    check_jwt_secret
    check_database_password
    check_cache_password
    check_admin_credentials
    check_domain_config
    check_session_security
    check_environment_type
    check_email_config
    check_sensitive_exposure
    
    generate_report
}

# Show usage if no arguments
if [ $# -eq 0 ] && [ ! -f ".env" ]; then
    echo "Usage: $0 [environment-file]"
    echo "Example: $0 .env.production"
    echo ""
    echo "If no file is specified, will check .env in current directory"
    exit 1
fi

# Run validation
main "$@"

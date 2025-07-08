#!/bin/bash

# WebEnable CMS Production Environment Setup Script
# This script helps configure secure production environment variables

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}=================================================${NC}"
    echo -e "${BLUE}  WebEnable CMS Production Environment Setup${NC}"
    echo -e "${BLUE}=================================================${NC}"
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

# Generate secure random string
generate_secret() {
    local length=${1:-32}
    openssl rand -base64 $length | tr -d "\n"
}

# Prompt for user input with default
prompt_with_default() {
    local prompt="$1"
    local default="$2"
    local result
    
    if [ -n "$default" ]; then
        read -p "$prompt [$default]: " result
        echo "${result:-$default}"
    else
        read -p "$prompt: " result
        echo "$result"
    fi
}

# Prompt for secure input (hidden)
prompt_secure() {
    local prompt="$1"
    local result
    
    read -s -p "$prompt: " result
    echo ""
    echo "$result"
}

# Check if running as root
check_root() {
    if [ "$EUID" -eq 0 ]; then
        print_error "Do not run this script as root for security reasons"
        exit 1
    fi
}

# Generate production environment file
generate_production_env() {
    print_status "Setting up production environment configuration..."
    echo ""
    
    # Domain configuration
    echo -e "${BLUE}Domain Configuration${NC}"
    DOMAIN=$(prompt_with_default "Enter your production domain (without https://)" "yourdomain.com")
    API_SUBDOMAIN=$(prompt_with_default "Enter API subdomain (optional)" "")
    
    # Database configuration
    echo ""
    echo -e "${BLUE}Database Configuration${NC}"
    DB_PASSWORD=$(generate_secret 24)
    print_status "Generated secure database password"
    
    # Cache configuration
    echo ""
    echo -e "${BLUE}Cache Configuration${NC}"
    CACHE_PASSWORD=$(generate_secret 24)
    print_status "Generated secure cache password"
    
    # Admin user configuration
    echo ""
    echo -e "${BLUE}Admin User Configuration${NC}"
    ADMIN_USERNAME=$(prompt_with_default "Admin username" "admin")
    ADMIN_PASSWORD=$(generate_secret 16)
    print_status "Generated secure admin password: $ADMIN_PASSWORD"
    print_warning "Save this password securely - you'll need it to log in!"
    
    # JWT Secret
    JWT_SECRET=$(generate_secret 32)
    print_status "Generated secure JWT secret"
    
    # Email configuration
    echo ""
    echo -e "${BLUE}Email Configuration (Optional)${NC}"
    SMTP_HOST=$(prompt_with_default "SMTP host" "smtp.${DOMAIN}")
    SMTP_PORT=$(prompt_with_default "SMTP port" "587")
    SMTP_USER=$(prompt_with_default "SMTP username" "noreply@${DOMAIN}")
    SMTP_PASS=$(prompt_secure "SMTP password (leave empty to skip)")
    
    # Build CORS origins
    if [ -n "$API_SUBDOMAIN" ]; then
        CORS_ORIGINS="https://${DOMAIN},https://www.${DOMAIN},https://${API_SUBDOMAIN}.${DOMAIN}"
        API_URL="https://${API_SUBDOMAIN}.${DOMAIN}/api"
    else
        CORS_ORIGINS="https://${DOMAIN},https://www.${DOMAIN}"
        API_URL="https://${DOMAIN}/api"
    fi
    
    # Create production .env file
    cat > .env.production << EOF
# WebEnable CMS Production Environment Configuration
# Generated on: $(date)
# WARNING: This file contains sensitive information. Do not commit to version control.

# =============================================================================
# SECURITY CONFIGURATION
# =============================================================================

JWT_SECRET=${JWT_SECRET}
ADMIN_USERNAME=${ADMIN_USERNAME}
ADMIN_PASSWORD=${ADMIN_PASSWORD}

# =============================================================================
# DATABASE CONFIGURATION
# =============================================================================

COUCHDB_USER=admin
COUCHDB_PASSWORD=${DB_PASSWORD}
COUCHDB_URL=http://admin:${DB_PASSWORD}@db:5984/

# =============================================================================
# CACHE CONFIGURATION
# =============================================================================

VALKEY_PASSWORD=${CACHE_PASSWORD}
VALKEY_URL=redis://:${CACHE_PASSWORD}@cache:6379

# =============================================================================
# SERVER CONFIGURATION
# =============================================================================

PORT=8080
CORS_ORIGINS=${CORS_ORIGINS}
SESSION_DOMAIN=${DOMAIN}
SESSION_SECURE=true
LOG_LEVEL=info

# =============================================================================
# FRONTEND CONFIGURATION
# =============================================================================

NEXT_PUBLIC_API_URL=${API_URL}
BACKEND_URL=http://backend:8080
NODE_ENV=production

# =============================================================================
# EMAIL CONFIGURATION
# =============================================================================

SMTP_HOST=${SMTP_HOST}
SMTP_PORT=${SMTP_PORT}
SMTP_USER=${SMTP_USER}
SMTP_PASS=${SMTP_PASS}
FROM_EMAIL=${SMTP_USER}
FROM_NAME=WebEnable CMS

# =============================================================================
# PRODUCTION OPTIMIZATIONS
# =============================================================================

GO_ENV=production
GIN_MODE=release
CACHE_TTL_POSTS=3600
CACHE_TTL_PAGES=1800
CACHE_TTL_API=300
RATE_LIMIT_API=100
RATE_LIMIT_AUTH=20
RATE_LIMIT_USER=150

# =============================================================================
# MONITORING & SECURITY
# =============================================================================

ENABLE_ACCESS_LOG=true
ENABLE_ERROR_LOG=true
ENABLE_SECURITY_LOG=true
HEALTH_CHECK_INTERVAL=30
HEALTH_CHECK_TIMEOUT=10
HSTS_MAX_AGE=31536000
FRAME_OPTIONS=DENY
CONTENT_TYPE_OPTIONS=nosniff
EOF

    # Set secure file permissions
    chmod 600 .env.production
    
    print_status "Production environment file created: .env.production"
    
    # Create summary
    echo ""
    echo -e "${GREEN}Production Configuration Summary:${NC}"
    echo "================================="
    echo "Domain: https://${DOMAIN}"
    echo "API URL: ${API_URL}"
    echo "Admin Username: ${ADMIN_USERNAME}"
    echo "Admin Password: ${ADMIN_PASSWORD}"
    echo ""
    print_warning "IMPORTANT: Save the admin password securely!"
    print_warning "File permissions set to 600 (owner read/write only)"
}

# Update Caddyfile for production
update_caddyfile() {
    echo ""
    print_status "Updating Caddyfile for production domain..."
    
    if [ -f "caddy/Caddyfile" ]; then
        # Backup original
        cp caddy/Caddyfile caddy/Caddyfile.backup
        
        # Replace localhost with actual domain
        sed -i.bak "s/localhost:80/${DOMAIN}:80/g" caddy/Caddyfile
        sed -i.bak "s/localhost:443/${DOMAIN}:443/g" caddy/Caddyfile
        sed -i.bak "s/auto_https off/auto_https on/g" caddy/Caddyfile
        
        print_status "Caddyfile updated for domain: ${DOMAIN}"
        print_status "Backup saved as: caddy/Caddyfile.backup"
    else
        print_warning "Caddyfile not found, skipping update"
    fi
}

# Create deployment checklist
create_checklist() {
    cat > DEPLOYMENT_CHECKLIST.md << EOF
# Production Deployment Checklist

## Pre-Deployment
- [ ] Domain DNS configured to point to server
- [ ] Server firewall configured (ports 80, 443, 5984)
- [ ] SSL certificates will be automatically generated by Caddy
- [ ] .env.production file is secure (permissions 600)

## Deployment
- [ ] Copy .env.production to server as .env
- [ ] Update Caddyfile with production domain
- [ ] Run: podman compose up -d
- [ ] Initialize admin user: podman compose exec backend ./main init-admin

## Post-Deployment
- [ ] Test HTTPS access: https://${DOMAIN}
- [ ] Test API: ${API_URL}/health
- [ ] Login with admin credentials
- [ ] Configure email settings if needed
- [ ] Set up backup schedule
- [ ] Configure monitoring

## Admin Credentials
- Username: ${ADMIN_USERNAME}
- Password: ${ADMIN_PASSWORD}

## Important URLs
- Website: https://${DOMAIN}
- API: ${API_URL}
- Database Admin: https://${DOMAIN}:5984/_utils

Remember to:
1. Change default passwords after first login
2. Set up regular backups
3. Configure monitoring and alerts
4. Review security settings regularly
EOF

    print_status "Deployment checklist created: DEPLOYMENT_CHECKLIST.md"
}

# Security recommendations
show_security_tips() {
    echo ""
    echo -e "${YELLOW}Security Recommendations:${NC}"
    echo "=========================="
    echo "1. 🔒 Change admin password after first login"
    echo "2. 📁 Never commit .env files to version control"
    echo "3. 🔄 Rotate secrets regularly (quarterly)"
    echo "4. 🔐 Use strong, unique passwords for all services"
    echo "5. 🛡️  Enable firewall on production server"
    echo "6. 📊 Set up monitoring and log aggregation"
    echo "7. 💾 Configure automated backups"
    echo "8. 🔍 Regular security audits and updates"
    echo ""
}

# Main setup process
main() {
    print_header
    check_root
    
    echo "This script will help you create a secure production environment configuration."
    echo ""
    
    generate_production_env
    update_caddyfile
    create_checklist
    show_security_tips
    
    echo -e "${GREEN}Production environment setup complete!${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Review .env.production file"
    echo "2. Copy to your production server as .env"
    echo "3. Follow DEPLOYMENT_CHECKLIST.md"
    echo "4. Deploy with: podman compose up -d"
}

# Run setup
main "$@"

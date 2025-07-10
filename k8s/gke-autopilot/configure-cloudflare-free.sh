#!/bin/bash

# Cloudflare Free Tier Configuration Script for WebEnable CMS
# This script configures Cloudflare DNS and security settings for free tier

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CLOUDFLARE_API_TOKEN=${CLOUDFLARE_API_TOKEN:-""}
CLOUDFLARE_ZONE_ID=${CLOUDFLARE_ZONE_ID:-""}
DOMAIN=${DOMAIN:-"webenable-cms.com"}
STATIC_IP=${STATIC_IP:-""}

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    if [ -z "$CLOUDFLARE_API_TOKEN" ]; then
        print_error "CLOUDFLARE_API_TOKEN is required"
        exit 1
    fi
    
    if [ -z "$CLOUDFLARE_ZONE_ID" ]; then
        print_error "CLOUDFLARE_ZONE_ID is required"
        exit 1
    fi
    
    if [ -z "$STATIC_IP" ]; then
        print_error "STATIC_IP is required"
        exit 1
    fi
    
    print_success "Prerequisites check completed"
}

# Function to create DNS records
create_dns_records() {
    print_status "Creating DNS records for free tier..."
    
    # List of hosts to create
    hosts=(
        "$DOMAIN"
        "www.$DOMAIN"
        "api.$DOMAIN"
        "admin.$DOMAIN"
    )
    
    for host in "${hosts[@]}"; do
        print_status "Creating A record for $host..."
        
        # Check if record already exists
        existing_record=$(curl -s -X GET "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/dns_records?name=$host" \
            -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
            -H "Content-Type: application/json" | jq -r '.result[0].id // empty')
        
        if [ -n "$existing_record" ]; then
            # Update existing record
            curl -s -X PUT "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/dns_records/$existing_record" \
                -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
                -H "Content-Type: application/json" \
                --data "{
                    \"type\": \"A\",
                    \"name\": \"$host\",
                    \"content\": \"$STATIC_IP\",
                    \"proxied\": true,
                    \"ttl\": 1
                }" > /dev/null
            
            print_success "Updated A record for $host"
        else
            # Create new record
            curl -s -X POST "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/dns_records" \
                -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
                -H "Content-Type: application/json" \
                --data "{
                    \"type\": \"A\",
                    \"name\": \"$host\",
                    \"content\": \"$STATIC_IP\",
                    \"proxied\": true,
                    \"ttl\": 1
                }" > /dev/null
            
            print_success "Created A record for $host"
        fi
    done
    
    print_success "DNS records created for free tier"
}

# Function to configure SSL/TLS settings (free tier)
configure_ssl_tls() {
    print_status "Configuring SSL/TLS settings for free tier..."
    
    # Set SSL/TLS encryption mode to Full (strict) - Free
    curl -s -X PATCH "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/ssl" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json" \
        --data '{"value":"full_strict"}' > /dev/null
    
    # Enable Always Use HTTPS - Free
    curl -s -X PATCH "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/always_use_https" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json" \
        --data '{"value":"on"}' > /dev/null
    
    # Enable HSTS - Free
    curl -s -X PATCH "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/security_header" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json" \
        --data '{
            "value": {
                "strict_transport_security": {
                    "enabled": true,
                    "max_age": 31536000,
                    "include_subdomains": true,
                    "preload": true
                }
            }
        }' > /dev/null
    
    print_success "SSL/TLS settings configured for free tier"
}

# Function to configure security settings (free tier)
configure_security() {
    print_status "Configuring security settings for free tier..."
    
    # Set security level to Medium - Free
    curl -s -X PATCH "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/security_level" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json" \
        --data '{"value":"medium"}' > /dev/null
    
    # Enable Browser Integrity Check - Free
    curl -s -X PATCH "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/browser_check" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json" \
        --data '{"value":"on"}' > /dev/null
    
    # Enable Challenge Passage - Free
    curl -s -X PATCH "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/challenge_ttl" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json" \
        --data '{"value":1800}' > /dev/null
    
    print_success "Security settings configured for free tier"
}

# Function to configure performance settings (free tier)
configure_performance() {
    print_status "Configuring performance settings for free tier..."
    
    # Enable Auto Minify for JS, CSS, and HTML - Free
    curl -s -X PATCH "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/minify" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json" \
        --data '{
            "value": {
                "css": "on",
                "html": "on",
                "js": "on"
            }
        }' > /dev/null
    
    # Enable Brotli compression - Free
    curl -s -X PATCH "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/brotli" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json" \
        --data '{"value":"on"}' > /dev/null
    
    # Enable Rocket Loader - Free
    curl -s -X PATCH "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/rocket_loader" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json" \
        --data '{"value":"on"}' > /dev/null
    
    print_success "Performance settings configured for free tier"
}

# Function to create basic rate limiting (free tier: 10,000 requests/day)
create_rate_limiting() {
    print_status "Creating basic rate limiting for free tier..."
    
    # Note: Free tier has limited rate limiting capabilities
    # This creates a basic rule for API endpoints
    curl -s -X POST "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/rulesets" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json" \
        --data '{
            "name": "API Rate Limiting (Free Tier)",
            "description": "Basic rate limiting for API endpoints (10,000 requests/day limit)",
            "kind": "zone",
            "phase": "http_ratelimit",
            "rules": [
                {
                    "expression": "(http.request.uri.path contains \"/api/\")",
                    "action": "block",
                    "ratelimit": {
                        "requests_per_period": 50,
                        "period": 60
                    }
                }
            ]
        }' > /dev/null
    
    print_success "Basic rate limiting created for free tier"
}

# Function to verify configuration
verify_configuration() {
    print_status "Verifying configuration..."
    
    # Check DNS records
    dns_records=$(curl -s -X GET "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/dns_records" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json")
    
    echo "DNS Records:"
    echo "$dns_records" | jq -r '.result[] | "  \(.name) -> \(.content)"'
    
    # Check SSL/TLS settings
    ssl_setting=$(curl -s -X GET "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/ssl" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json")
    
    echo ""
    echo "SSL/TLS Mode: $(echo "$ssl_setting" | jq -r '.result.value')"
    
    # Check security level
    security_level=$(curl -s -X GET "https://api.cloudflare.com/client/v4/zones/$CLOUDFLARE_ZONE_ID/settings/security_level" \
        -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
        -H "Content-Type: application/json")
    
    echo "Security Level: $(echo "$security_level" | jq -r '.result.value')"
    
    print_success "Configuration verification completed"
}

# Function to show free tier limitations
show_free_tier_info() {
    print_status "Cloudflare Free Tier Information:"
    echo ""
    echo "✅ Included Features:"
    echo "  - SSL/TLS encryption"
    echo "  - Global CDN (200+ data centers)"
    echo "  - Basic DDoS protection"
    echo "  - DNS management"
    echo "  - Auto minify (JS, CSS, HTML)"
    echo "  - Brotli compression"
    echo "  - Rocket Loader"
    echo "  - Browser integrity check"
    echo "  - Basic WAF rules"
    echo ""
    echo "⚠️  Free Tier Limitations:"
    echo "  - Rate limiting: 10,000 requests/day"
    echo "  - Bot Fight Mode: Not available"
    echo "  - Advanced WAF rules: Limited"
    echo "  - Image optimization: Not available"
    echo "  - Load balancing: Not available"
    echo "  - Priority support: Not available"
    echo ""
    echo "💰 Cost: $0/month"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --dns-only        Only configure DNS records"
    echo "  --ssl-only        Only configure SSL/TLS settings"
    echo "  --security-only   Only configure security settings"
    echo "  --performance-only Only configure performance settings"
    echo "  --verify          Verify current configuration"
    echo "  --info            Show free tier information"
    echo "  --help            Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  CLOUDFLARE_API_TOKEN  Cloudflare API token (required)"
    echo "  CLOUDFLARE_ZONE_ID    Cloudflare zone ID (required)"
    echo "  DOMAIN                Domain name (default: webenable-cms.com)"
    echo "  STATIC_IP             Static IP address (required)"
}

# Main function
main() {
    case "${1:-all}" in
        "all")
            check_prerequisites
            create_dns_records
            configure_ssl_tls
            configure_security
            configure_performance
            create_rate_limiting
            verify_configuration
            show_free_tier_info
            ;;
        "dns-only")
            check_prerequisites
            create_dns_records
            ;;
        "ssl-only")
            check_prerequisites
            configure_ssl_tls
            ;;
        "security-only")
            check_prerequisites
            configure_security
            ;;
        "performance-only")
            check_prerequisites
            configure_performance
            ;;
        "verify")
            check_prerequisites
            verify_configuration
            ;;
        "info")
            show_free_tier_info
            ;;
        "help"|"--help"|"-h")
            show_usage
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@" 
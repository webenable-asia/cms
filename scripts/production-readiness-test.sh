#!/bin/bash

# Production Readiness Test Script
# Tests all aspects of the CMS for production deployment

# Load environment variables if available
if [ -f ".env.production" ]; then
    set -a
    source .env.production
    set +a
elif [ -f ".env" ]; then
    set -a
    source .env
    set +a
fi

echo "🚀 Production Readiness Test Suite"
echo "=================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results
PASSED=0
FAILED=0
WARNINGS=0

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ $2${NC}"
        ((FAILED++))
    fi
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
    ((WARNINGS++))
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Test 1: Docker Services Health
echo "📋 Test 1: Docker Services Health"
echo "================================="

# Check if docker-compose is running
if docker-compose ps | grep -q "Up"; then
    print_result 0 "Docker Compose services are running"
else
    print_result 1 "Docker Compose services are not running"
    echo "Run: docker-compose up -d"
    exit 1
fi

# Check individual service health
services=("cms-frontend-1" "cms-backend-1" "cms-db-1" "cms-cache-1" "cms-caddy-1")
for service in "${services[@]}"; do
    status=$(docker ps --format "{{.Names}}\t{{.Status}}" | grep "$service" || echo "not found")
    if echo "$status" | grep -q "Up"; then
        print_result 0 "$service is healthy"
    else
        print_result 1 "$service is not healthy: $status"
    fi
done

echo ""

# Test 2: Application Connectivity
echo "📋 Test 2: Application Connectivity"
echo "==================================="

# Test main website
if curl -s -o /dev/null -w "%{http_code}" http://localhost | grep -q "200"; then
    print_result 0 "Frontend is accessible (HTTP 200)"
else
    print_result 1 "Frontend is not accessible"
fi

# Test API endpoints
if curl -s -o /dev/null -w "%{http_code}" http://localhost/api/posts | grep -q "200"; then
    print_result 0 "Backend API is accessible (HTTP 200)"
else
    print_result 1 "Backend API is not accessible"
fi

# Test database proxy
if curl -s -o /dev/null -w "%{http_code}" http://localhost:5984 | grep -q "200"; then
    print_result 0 "Database proxy is accessible (HTTP 200)"
else
    print_result 1 "Database proxy is not accessible"
fi

echo ""

# Test 3: Security Headers
echo "📋 Test 3: Security Headers"
echo "==========================="

headers_response=$(curl -s -I http://localhost)

# Check for essential security headers
if echo "$headers_response" | grep -q "X-Frame-Options"; then
    print_result 0 "X-Frame-Options header present"
else
    print_result 1 "X-Frame-Options header missing"
fi

if echo "$headers_response" | grep -q "X-Content-Type-Options"; then
    print_result 0 "X-Content-Type-Options header present"
else
    print_result 1 "X-Content-Type-Options header missing"
fi

if echo "$headers_response" | grep -q "Referrer-Policy"; then
    print_result 0 "Referrer-Policy header present"
else
    print_result 1 "Referrer-Policy header missing"
fi

# Check if server header is hidden
if echo "$headers_response" | grep -q "Server:"; then
    print_warning "Server header is exposed (consider hiding for security)"
else
    print_result 0 "Server header is hidden"
fi

echo ""

# Test 4: Performance Features
echo "📋 Test 4: Performance Features"
echo "==============================="

# Test compression
compression_response=$(curl -s -H "Accept-Encoding: gzip" -I http://localhost)
if echo "$compression_response" | grep -q "Content-Encoding: gzip"; then
    print_result 0 "Gzip compression is enabled"
else
    print_result 1 "Gzip compression is not enabled"
fi

# Test response times
frontend_time=$(curl -s -o /dev/null -w "%{time_total}" http://localhost)
api_time=$(curl -s -o /dev/null -w "%{time_total}" http://localhost/api/posts)

if awk "BEGIN {exit !($frontend_time < 2.0)}"; then
    print_result 0 "Frontend response time acceptable (${frontend_time}s)"
else
    print_warning "Frontend response time slow (${frontend_time}s)"
fi

if awk "BEGIN {exit !($api_time < 1.0)}"; then
    print_result 0 "API response time acceptable (${api_time}s)"
else
    print_warning "API response time slow (${api_time}s)"
fi

echo ""

# Test 5: Database Functionality
echo "📋 Test 5: Database Functionality"
echo "================================="

# Test CouchDB welcome endpoint
db_response=$(curl -s http://localhost:5984)
if echo "$db_response" | grep -q "Welcome"; then
    print_result 0 "CouchDB is responding correctly"
else
    print_result 1 "CouchDB is not responding correctly"
fi

# Test authenticated database access (if credentials are available)
if [ -n "${COUCHDB_USER}" ] && [ -n "${COUCHDB_PASSWORD}" ]; then
    auth_response=$(curl -s -u "${COUCHDB_USER}:${COUCHDB_PASSWORD}" http://localhost:5984/_all_dbs)
    if echo "$auth_response" | grep -q "\["; then
        print_result 0 "Database authentication working"
    else
        print_result 1 "Database authentication failed"
    fi
else
    print_warning "Database credentials not set for authentication test"
fi

echo ""

# Test 6: Resource Usage
echo "📋 Test 6: Resource Usage"
echo "========================="

# Get container stats with names
docker stats --no-stream --format "{{.Name}},{{.CPUPerc}},{{.MemUsage}}" > /tmp/cms_stats.txt

while IFS=',' read -r container cpu memory; do
    # Only process CMS containers
    if [[ "$container" == cms-* ]]; then
        # Remove % from CPU
        cpu_clean=$(echo "$cpu" | sed 's/%//')
        
        # Extract memory usage (before the /)
        mem_usage=$(echo "$memory" | cut -d'/' -f1 | sed 's/MiB//' | sed 's/GiB//' | tr -d ' ')
        
        # Simple comparison using shell arithmetic for CPU (convert to integer)
        cpu_int=$(echo "$cpu_clean" | cut -d'.' -f1)
        if [ "$cpu_int" -lt 50 ] 2>/dev/null; then
            print_result 0 "$container CPU usage acceptable (${cpu})"
        else
            print_warning "$container CPU usage high (${cpu})"
        fi
        
        # Simple comparison for memory (should be less than 1000 MiB)
        mem_int=$(echo "$mem_usage" | cut -d'.' -f1)
        if [ "$mem_int" -lt 1000 ] 2>/dev/null; then
            print_result 0 "$container memory usage acceptable (${memory})"
        else
            print_warning "$container memory usage high (${memory})"
        fi
    fi
done < /tmp/cms_stats.txt

rm -f /tmp/cms_stats.txt

echo ""

# Test 7: Environment Configuration
echo "📋 Test 7: Environment Configuration"
echo "===================================="

# Check for production environment file
if [ -f ".env.production" ]; then
    print_result 0 "Production environment file exists"
    
    # Check for required variables
    required_vars=("COUCHDB_USER" "COUCHDB_PASSWORD" "JWT_SECRET" "VALKEY_PASSWORD")
    for var in "${required_vars[@]}"; do
        if grep -q "^${var}=" .env.production; then
            print_result 0 "$var is configured"
        else
            print_result 1 "$var is missing from production config"
        fi
    done
else
    print_result 1 "Production environment file missing"
fi

# Check for sensitive data in version control
if git ls-files | grep -q "\.env$"; then
    print_result 1 "Sensitive .env file is tracked in git"
else
    print_result 0 "Environment files are properly ignored"
fi

echo ""

# Test 8: SSL/TLS Readiness
echo "📋 Test 8: SSL/TLS Readiness"
echo "============================"

# Check Caddy configuration for HTTPS readiness
if grep -q "auto_https off" caddy/Caddyfile; then
    print_warning "HTTPS is disabled (set auto_https on for production with domain)"
else
    print_result 0 "HTTPS is enabled in Caddy configuration"
fi

# Check if localhost is used (not production ready)
if grep -q "localhost" caddy/Caddyfile; then
    print_warning "Using localhost in Caddyfile (replace with domain for production)"
else
    print_result 0 "Domain configuration ready"
fi

echo ""

# Test 9: Log Configuration
echo "📋 Test 9: Log Configuration"
echo "============================"

# Check if logs are properly configured
log_output=$(docker logs cms-caddy-1 2>&1 | tail -n 5)
if [ -n "$log_output" ]; then
    print_result 0 "Caddy logging is working"
else
    print_warning "No recent Caddy logs found"
fi

# Check for structured logging
if echo "$log_output" | grep -q "{"; then
    print_result 0 "Structured JSON logging is enabled"
else
    print_warning "Structured logging not detected"
fi

echo ""

# Test 10: Backup and Data Persistence
echo "📋 Test 10: Data Persistence"
echo "============================"

# Check Docker volumes
volumes=$(docker volume ls --format "{{.Name}}" | grep cms)
if [ -n "$volumes" ]; then
    print_result 0 "Docker volumes are configured for data persistence"
    echo "   Volumes: $volumes"
else
    print_result 1 "No Docker volumes found for data persistence"
fi

echo ""

# Final Summary
echo "🎯 Test Summary"
echo "==============="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    if [ $WARNINGS -eq 0 ]; then
        echo -e "${GREEN}🎉 Production Ready! All tests passed.${NC}"
        exit 0
    else
        echo -e "${YELLOW}⚠️  Production Ready with warnings. Address warnings for optimal security.${NC}"
        exit 0
    fi
else
    echo -e "${RED}❌ Not Production Ready. Please fix the failed tests.${NC}"
    exit 1
fi

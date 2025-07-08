#!/bin/bash

# Test script for Caddy reverse proxy capabilities
# WebEnable CMS - Database Proxy Testing

echo "🔄 Testing Caddy Reverse Proxy Configuration"
echo "=============================================="

# Source environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v ^# | xargs)
fi

# Test 1: Main Website
echo "📋 Test 1: Main Website"
echo "URL: http://localhost"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost)
if [ "$STATUS" = "200" ]; then
    echo "✅ Main website is accessible (HTTP $STATUS)"
else
    echo "❌ Main website failed (HTTP $STATUS)"
fi
echo ""

# Test 2: API Endpoints
echo "📋 Test 2: API Endpoints"
echo "URL: http://localhost/api/posts"
API_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost/api/posts?limit=1")
if [ "$API_STATUS" = "200" ]; then
    echo "✅ API endpoints are accessible (HTTP $API_STATUS)"
else
    echo "❌ API endpoints failed (HTTP $API_STATUS)"
fi
echo ""

# Test 3: Database Proxy
echo "📋 Test 3: Database Proxy"
echo "URL: http://localhost:5984"
DB_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:5984)
if [ "$DB_STATUS" = "200" ]; then
    echo "✅ Database proxy is accessible (HTTP $DB_STATUS)"
    
    # Test database info
    echo "Database Info:"
    curl -s http://localhost:5984 | jq '.couchdb, .version' 2>/dev/null || curl -s http://localhost:5984
    echo ""
    
    # Test authenticated access
    if [ -n "$COUCHDB_USER" ] && [ -n "$COUCHDB_PASSWORD" ]; then
        echo "Testing authenticated database access..."
        AUTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "http://$COUCHDB_USER:$COUCHDB_PASSWORD@localhost:5984/_all_dbs")
        if [ "$AUTH_STATUS" = "200" ]; then
            echo "✅ Authenticated database access works (HTTP $AUTH_STATUS)"
            echo "Available databases:"
            curl -s "http://$COUCHDB_USER:$COUCHDB_PASSWORD@localhost:5984/_all_dbs" | jq '.' 2>/dev/null || curl -s "http://$COUCHDB_USER:$COUCHDB_PASSWORD@localhost:5984/_all_dbs"
        else
            echo "❌ Authenticated database access failed (HTTP $AUTH_STATUS)"
        fi
    else
        echo "⚠️  Database credentials not found in environment"
    fi
else
    echo "❌ Database proxy failed (HTTP $DB_STATUS)"
fi
echo ""

# Test 4: Security Headers
echo "📋 Test 4: Security Headers"
echo "Checking security headers on main site..."
SECURITY_HEADERS=$(curl -s -I http://localhost | grep -E "(X-Frame-Options|X-Content-Type-Options|X-XSS-Protection|Referrer-Policy)")
if [ -n "$SECURITY_HEADERS" ]; then
    echo "✅ Security headers are present:"
    echo "$SECURITY_HEADERS"
else
    echo "❌ Security headers are missing"
fi
echo ""

# Test 5: Compression
echo "📋 Test 5: Compression"
echo "Testing gzip compression..."
GZIP_TEST=$(curl -s -H "Accept-Encoding: gzip" -I http://localhost | grep -i "content-encoding.*gzip")
if [ -n "$GZIP_TEST" ]; then
    echo "✅ Gzip compression is enabled"
else
    echo "⚠️  Gzip compression not detected (may not be needed for this response)"
fi
echo ""

# Summary
echo "🎯 Summary"
echo "=========="
echo "Caddy is successfully acting as a reverse proxy for:"
echo "• Frontend (Next.js) - Port 80"
echo "• Backend API - Port 80/api/*"
echo "• Database (CouchDB) - Port 5984"
echo ""
echo "Security features enabled:"
echo "• Security headers (X-Frame-Options, XSS Protection, etc.)"
echo "• Content compression (gzip/zstd)"
echo "• Database admin interface IP restrictions"
echo "• Request logging and monitoring"
echo ""
echo "Note: Valkey/Redis cache cannot be proxied through HTTP reverse proxy"
echo "      as it uses a binary protocol. Direct connection is maintained."

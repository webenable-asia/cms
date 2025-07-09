#!/bin/bash

# WebEnable CMS - Admin Login Troubleshooting Script

echo "🔧 WebEnable CMS Admin Login Troubleshooting"
echo "============================================"

# Test 1: Check Docker services
echo "1️⃣ Checking Docker services..."
if docker compose ps | grep -q "Up"; then
    echo "✅ Docker services are running"
    docker compose ps
else
    echo "❌ Docker services not running"
    echo "   Solution: Run 'docker compose up -d'"
    exit 1
fi

echo ""

# Test 2: Check API health
echo "2️⃣ Testing API health..."
API_HEALTH=$(curl -s "http://localhost/api/health" 2>/dev/null)
if echo "$API_HEALTH" | jq -e '.status' > /dev/null 2>&1; then
    echo "✅ API is healthy"
    echo "$API_HEALTH" | jq '.'
else
    echo "❌ API not responding"
    echo "   Check backend logs: docker compose logs backend"
    exit 1
fi

echo ""

# Test 3: Check CouchDB connectivity
echo "3️⃣ Testing CouchDB connection..."
if [ -f .env ]; then
    COUCHDB_PASSWORD=$(grep "^COUCHDB_PASSWORD=" .env | cut -d '=' -f2)
    if curl -s -f "http://admin:${COUCHDB_PASSWORD}@localhost:5984" > /dev/null; then
        echo "✅ CouchDB is accessible"
    else
        echo "❌ Cannot connect to CouchDB"
        echo "   Check database logs: docker compose logs db"
        exit 1
    fi
else
    echo "❌ .env file not found"
    exit 1
fi

echo ""

# Test 4: Check for admin users
echo "4️⃣ Checking for admin users..."
ADMIN_COUNT=$(curl -s "http://admin:${COUCHDB_PASSWORD}@localhost:5984/users/_all_docs?include_docs=true" | jq '.rows[] | select(.doc.username == "admin")' | jq -s length)
if [ "$ADMIN_COUNT" -gt 0 ]; then
    echo "✅ Found $ADMIN_COUNT admin user(s)"
else
    echo "❌ No admin users found"
    echo "   Solution: Run './create_admin_user.sh'"
    exit 1
fi

echo ""

# Test 5: Test login
echo "5️⃣ Testing admin login..."
LOGIN_RESPONSE=$(curl -s -X POST "http://localhost/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}')

if echo "$LOGIN_RESPONSE" | jq -e '.token' > /dev/null 2>&1; then
    echo "🎉 SUCCESS! Admin login is working"
    echo "   👤 Username: admin"
    echo "   🔑 Password: admin123"
    echo "   🔗 Admin Panel: http://localhost/admin"
else
    echo "❌ Login failed. Response:"
    echo "$LOGIN_RESPONSE"
    echo ""
    echo "🔧 Suggested fixes:"
    echo "   1. Re-create admin user: ./create_admin_user.sh"
    echo "   2. Check backend logs: docker compose logs backend"
    echo "   3. Restart services: docker compose restart"
fi

echo ""
echo "✅ Troubleshooting complete!"

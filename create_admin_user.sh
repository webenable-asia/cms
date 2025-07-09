#!/bin/bash

# WebEnable CMS - Create Admin User Script
# This script creates a fresh admin user with proper authentication

set -e  # Exit on any error

echo "🚀 Creating fresh admin user for WebEnable CMS..."

# Load environment variables
if [ -f .env ]; then
    COUCHDB_PASSWORD=$(grep "^COUCHDB_PASSWORD=" .env | cut -d '=' -f2)
    if [ -z "$COUCHDB_PASSWORD" ]; then
        echo "❌ Error: COUCHDB_PASSWORD not found in .env file"
        exit 1
    fi
    echo "✅ Loaded CouchDB password from .env"
else
    echo "❌ Error: .env file not found"
    exit 1
fi

# Check CouchDB connectivity
echo "🔍 Testing CouchDB connection..."
if ! curl -s -f "http://admin:${COUCHDB_PASSWORD}@localhost:5984" > /dev/null; then
    echo "❌ Error: Cannot connect to CouchDB"
    echo "   - Check if Docker services are running: docker compose ps"
    echo "   - Verify CouchDB password in .env file"
    exit 1
fi
echo "✅ CouchDB is accessible"

# Check if backend directory exists
if [ ! -d "backend" ]; then
    echo "❌ Error: backend directory not found. Run this script from project root."
    exit 1
fi

echo "🧹 Cleaning up existing admin users..."

# Get all existing admin users
EXISTING_ADMINS=$(curl -s "http://admin:${COUCHDB_PASSWORD}@localhost:5984/users/_all_docs?include_docs=true" | jq -r '.rows[] | select(.doc.username == "admin") | "\(.id),\(.doc._rev)"' 2>/dev/null || echo "")

if [ -n "$EXISTING_ADMINS" ]; then
    echo "📋 Found existing admin users, removing them..."
    echo "$EXISTING_ADMINS" | while IFS=',' read -r id rev; do
        if [ -n "$id" ] && [ -n "$rev" ]; then
            echo "   Deleting admin user: $id"
            curl -s -X DELETE "http://admin:${COUCHDB_PASSWORD}@localhost:5984/users/$id?rev=$rev" > /dev/null || true
        fi
    done
    echo "✅ Cleaned up existing admin users"
else
    echo "✅ No existing admin users found"
fi

echo "🔐 Generating secure password hash..."

# Create temporary hash generator
cat > temp_hash_generator.go << 'EOF'
package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run temp_hash_generator.go <password>")
		os.Exit(1)
	}
	
	password := os.Args[1]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Printf("Error generating hash: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Print(string(hash))
}
EOF

# Generate bcrypt hash
BCRYPT_HASH=$(cd backend && go run ../temp_hash_generator.go "admin123" 2>/dev/null)
if [ $? -ne 0 ] || [ -z "$BCRYPT_HASH" ]; then
    echo "❌ Error: Failed to generate password hash"
    rm -f temp_hash_generator.go
    exit 1
fi

# Clean up temp file
rm -f temp_hash_generator.go

echo "✅ Password hash generated successfully"

# Generate new admin user ID
ADMIN_USER_ID=$(uuidgen | tr '[:upper:]' '[:lower:]')
echo "📝 Creating admin user with ID: $ADMIN_USER_ID"

# Create the admin user
CREATE_RESPONSE=$(curl -s -X POST "http://admin:${COUCHDB_PASSWORD}@localhost:5984/users" \
  -H "Content-Type: application/json" \
  -d "{
    \"_id\": \"$ADMIN_USER_ID\",
    \"username\": \"admin\",
    \"email\": \"admin@webenable.asia\",
    \"password_hash\": \"$BCRYPT_HASH\",
    \"role\": \"admin\",
    \"active\": true,
    \"created_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"updated_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
  }")

# Check if user creation was successful
if echo "$CREATE_RESPONSE" | jq -e '.ok' > /dev/null 2>&1; then
    echo "✅ Admin user created successfully in database"
else
    echo "❌ Error creating admin user. Response: $CREATE_RESPONSE"
    exit 1
fi

# Wait a moment for the database to sync
echo "⏳ Waiting for database sync..."
sleep 3

# Verify the user was created
echo "🔍 Verifying admin user creation..."
USER_CHECK=$(curl -s "http://admin:${COUCHDB_PASSWORD}@localhost:5984/users/$ADMIN_USER_ID" | jq -r '.username' 2>/dev/null)
if [ "$USER_CHECK" = "admin" ]; then
    echo "✅ Admin user verified in database"
else
    echo "❌ Error: Admin user not found in database after creation"
    exit 1
fi

# Test the authentication
echo "🧪 Testing admin login via API..."
LOGIN_RESPONSE=$(curl -s -X POST "http://localhost/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}')

# Check if login was successful
if echo "$LOGIN_RESPONSE" | jq -e '.token' > /dev/null 2>&1; then
    echo "🎉 SUCCESS! Admin login test passed"
    echo ""
    echo "📋 Admin User Details:"
    echo "   👤 Username: admin"
    echo "   🔑 Password: admin123"
    echo "   📧 Email: admin@webenable.asia"
    echo "   🔗 Admin Panel: http://localhost/admin"
    echo ""
    echo "⚠️  IMPORTANT: Change the password after first login!"
    echo ""
    echo "🎯 You can now login to the admin panel!"
else
    echo "❌ Login test failed. API Response:"
    echo "$LOGIN_RESPONSE"
    echo ""
    echo "🔧 Troubleshooting steps:"
    echo "   1. Check if backend service is running: docker compose ps"
    echo "   2. Check backend logs: docker compose logs backend"
    echo "   3. Verify API health: curl http://localhost/api/health"
    exit 1
fi

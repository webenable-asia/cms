#!/bin/bash

# WebEnable CMS - Populate Contacts Script
# This script populates the CouchDB database with sample contact form submissions

set -e  # Exit on any error

echo "🚀 Populating CouchDB with sample contacts..."

# Load CouchDB password from .env file
if [ -f .env ]; then
    COUCHDB_PASSWORD=$(grep "^COUCHDB_PASSWORD=" .env | cut -d '=' -f2)
    if [ -z "$COUCHDB_PASSWORD" ]; then
        echo "❌ Error: COUCHDB_PASSWORD not found in .env file"
        exit 1
    fi
    echo "✅ Loaded CouchDB password from .env"
else
    echo "❌ Error: .env file not found. Please ensure you're in the project root directory."
    exit 1
fi

# Verify CouchDB connection
echo "🔍 Verifying CouchDB connection..."
if curl -s -f "http://admin:${COUCHDB_PASSWORD}@localhost:5984" > /dev/null; then
    echo "✅ CouchDB is accessible"
else
    echo "❌ Error: Cannot connect to CouchDB. Please ensure:"
    echo "   - Docker services are running (docker compose ps)"
    echo "   - CouchDB password in .env is correct"
    echo "   - CouchDB is accessible on localhost:5984"
    exit 1
fi

# Check if Go is available and we're in the right directory
if [ ! -f "backend/go.mod" ]; then
    echo "❌ Error: backend/go.mod not found. Please run this script from the project root directory."
    exit 1
fi

# Run the populate_contacts.go script from the backend directory
echo "📝 Running populate_contacts.go script..."
cd backend

# Set the COUCHDB_URL environment variable for the Go script
export COUCHDB_URL="http://admin:${COUCHDB_PASSWORD}@localhost:5984"

# Run the Go script
if go run ../scripts/populate_contacts.go; then
    echo "✅ Contacts populated successfully!"
    echo ""
    echo "📊 Database Summary:"
    
    # Count total contacts
    TOTAL_CONTACTS=$(curl -s "http://admin:${COUCHDB_PASSWORD}@localhost:5984/contacts/_all_docs" | jq -r '.total_rows')
    echo "   📞 Total contacts in database: ${TOTAL_CONTACTS}"
    
    # Show contact status breakdown
    echo "   📋 Contact status breakdown:"
    NEW_COUNT=$(curl -s "http://admin:${COUCHDB_PASSWORD}@localhost:5984/contacts/_design/contacts/_view/by_status?key=\"new\"" 2>/dev/null | jq -r '.rows | length' || echo "0")
    READ_COUNT=$(curl -s "http://admin:${COUCHDB_PASSWORD}@localhost:5984/contacts/_design/contacts/_view/by_status?key=\"read\"" 2>/dev/null | jq -r '.rows | length' || echo "0")
    REPLIED_COUNT=$(curl -s "http://admin:${COUCHDB_PASSWORD}@localhost:5984/contacts/_design/contacts/_view/by_status?key=\"replied\"" 2>/dev/null | jq -r '.rows | length' || echo "0")
    
    echo "      - New: ${NEW_COUNT}"
    echo "      - Read: ${READ_COUNT}" 
    echo "      - Replied: ${REPLIED_COUNT}"
    
    echo ""
    echo "🔗 Access contacts via:"
    echo "   📱 Admin Panel: http://localhost/admin/contacts"
    echo "   🔧 API Endpoint: http://localhost/api/contacts (requires authentication)"
    echo "   💾 Direct DB: http://localhost:5984/_utils/#database/contacts/_all_docs"
    
else
    echo "❌ Error: Failed to populate contacts. Check the error messages above."
    exit 1
fi

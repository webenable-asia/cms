#!/bin/bash

# Development environment
echo "🚀 Starting development environment..."

export NODE_ENV=development
export GO_ENV=development

# Build images if they don't exist
docker-compose build --parallel

# Start services
docker-compose up -d

echo "✅ Development environment started!"
echo "📱 Frontend: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"
echo "🗄️  CouchDB: http://localhost:5984"
echo "💾 Valkey: localhost:6379"

# Show logs
docker-compose logs -f

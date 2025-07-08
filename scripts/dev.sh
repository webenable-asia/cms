#!/bin/bash

# Development environment (Podman)
echo "🚀 Starting development environment with Podman..."

export NODE_ENV=development
export GO_ENV=development

# Build images if they don't exist
podman compose build --parallel

# Start services
podman compose up -d

echo "✅ Development environment started!"
echo "📱 Frontend: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"
echo "🗄️  CouchDB: http://localhost:5984"
echo "💾 Valkey: localhost:6379"

# Show logs
podman compose logs -f

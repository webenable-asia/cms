#!/bin/bash

# Production environment (Docker)
echo "🚀 Building and starting production environment with Docker..."

export NODE_ENV=production
export GO_ENV=production

# Build optimized images
docker compose -f docker-compose.yml build --parallel --no-cache

# Start services
docker compose -f docker-compose.yml up -d

echo "✅ Production environment started!"
echo "🌐 Application: http://localhost"
echo "📊 Health check: http://localhost/health"

# Show logs
docker compose -f docker-compose.yml logs -f

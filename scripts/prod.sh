#!/bin/bash

# Production environment
echo "🚀 Building and starting production environment..."

export NODE_ENV=production
export GO_ENV=production

# Build optimized images
docker-compose -f docker-compose.prod.yml build --parallel --no-cache

# Start services
docker-compose -f docker-compose.prod.yml up -d

echo "✅ Production environment started!"
echo "🌐 Application: http://localhost"
echo "📊 Health check: http://localhost/health"

# Show logs
docker-compose -f docker-compose.prod.yml logs -f

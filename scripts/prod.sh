#!/bin/bash

# Production environment (Podman)
echo "🚀 Building and starting production environment with Podman..."

export NODE_ENV=production
export GO_ENV=production

# Build optimized images
podman compose -f docker-compose.prod.yml build --parallel --no-cache

# Start services
podman compose -f docker-compose.prod.yml up -d

echo "✅ Production environment started!"
echo "🌐 Application: http://localhost"
echo "📊 Health check: http://localhost/health"

# Show logs
podman compose -f docker-compose.prod.yml logs -f

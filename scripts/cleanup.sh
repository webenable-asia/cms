#!/bin/bash

echo "🧹 Cleaning up Docker resources..."

# Stop all containers
docker-compose down
docker-compose -f docker-compose.prod.yml down 2>/dev/null || true

# Remove unused images
docker image prune -f

# Remove unused volumes (careful!)
echo "⚠️  Do you want to remove unused volumes? This will delete database data! (y/N)"
read -r response
if [[ "$response" =~ ^[Yy]$ ]]; then
    docker volume prune -f
    echo "🗑️  Volumes removed"
else
    echo "📦 Volumes preserved"
fi

# Remove unused networks
docker network prune -f

# Remove build cache
docker builder prune -f

echo "✅ Cleanup completed!"
docker system df

#!/bin/bash

echo "🧹 Cleaning up Podman resources..."

# Stop all containers
podman compose down
podman compose -f podman-compose.yml down 2>/dev/null || true

# Remove unused images
podman image prune -f

# Remove unused volumes (careful!)
echo "⚠️  Do you want to remove unused volumes? This will delete database data! (y/N)"
read -r response
if [[ "$response" =~ ^[Yy]$ ]]; then
    podman volume prune -f
    echo "🗑️  Volumes removed"
else
    echo "📦 Volumes preserved"
fi

# Remove unused networks
podman network prune -f

# Remove build cache (if supported)
podman builder prune -f 2>/dev/null || echo "⚠️  Builder prune not supported, skipping"

echo "✅ Cleanup completed!"
podman system df

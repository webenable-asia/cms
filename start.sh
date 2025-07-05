#!/bin/bash

echo "🚀 Starting WebEnable CMS Development Environment..."
echo ""
echo "This will start:"
echo "  - CouchDB at http://localhost:5984"
echo "  - Go Backend API at http://localhost:8080"
echo "  - Next.js Frontend at http://localhost:3000"
echo ""
echo "Default admin credentials: admin / password"
echo ""

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Start the services
echo "📦 Building and starting containers..."
docker-compose up --build

echo ""
echo "🎉 Development environment started successfully!"
echo ""
echo "Access your applications:"
echo "  Frontend: http://localhost:3000"
echo "  Backend API: http://localhost:8080"
echo "  CouchDB Admin: http://localhost:5984/_utils (admin/password)"
echo ""
echo "To stop the environment, press Ctrl+C"

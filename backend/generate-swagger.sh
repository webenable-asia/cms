#!/bin/bash

echo "Generating Swagger documentation..."

# Generate swagger docs
/Users/tsaa/go/bin/swag init

if [ $? -eq 0 ]; then
    echo "✅ Swagger documentation generated successfully!"
    echo "📄 Files created:"
    echo "  - docs/docs.go"
    echo "  - docs/swagger.json" 
    echo "  - docs/swagger.yaml"
    echo ""
    echo "🌐 Access Swagger UI at: http://localhost:8080/swagger/index.html"
    echo "📊 Access JSON spec at: http://localhost:8080/swagger/doc.json"
else
    echo "❌ Failed to generate Swagger documentation"
    exit 1
fi

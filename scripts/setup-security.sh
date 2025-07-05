#!/bin/bash

echo "🔒 Setting up WebEnable CMS Security..."

# Generate JWT secret
JWT_SECRET=$(openssl rand -base64 32)
echo "✅ Generated JWT Secret"

# Generate admin password
ADMIN_PASS=$(openssl rand -base64 12)
echo "✅ Generated Admin Password: $ADMIN_PASS"

# Create production env file
cat > backend/.env.production <<EOF
JWT_SECRET=$JWT_SECRET
COUCHDB_URL=http://admin:password@db:5984/
PORT=8080
ALLOWED_ORIGINS=https://your-domain.com
SMTP_HOST=your-smtp-host
SMTP_PORT=587
SMTP_USER=your-email
SMTP_PASS=your-smtp-password
EOF

echo "✅ Created .env.production"
echo ""
echo "📝 Next steps:"
echo "1. Update ALLOWED_ORIGINS with your domain"
echo "2. Configure SMTP settings"
echo "3. Run: ADMIN_PASSWORD=$ADMIN_PASS go run scripts/init_admin.go"
echo ""
echo "⚠️  Save the admin password securely!"
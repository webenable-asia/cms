#!/bin/bash

# Sample data population script for WebEnable CMS
# This script populates the database with sample posts and contacts

set -e

DB_URL="http://admin:dGR5FbtwrkTCEbl1xFQZikKw7rIazzA6@localhost:5984"
POSTS_DB="$DB_URL/posts"
CONTACTS_DB="$DB_URL/contacts"

echo "🚀 Populating WebEnable CMS with sample data..."

# Function to generate UUID
generate_uuid() {
    uuidgen | tr '[:upper:]' '[:lower:]'
}

# Function to get current timestamp
current_timestamp() {
    date -u +%Y-%m-%dT%H:%M:%SZ
}

# Create sample posts
echo "📝 Creating sample posts..."

# Post 1: Welcome post
POST_ID_1=$(generate_uuid)
curl -X POST "$POSTS_DB" -H "Content-Type: application/json" -d "{
  \"_id\": \"$POST_ID_1\",
  \"title\": \"Welcome to WebEnable CMS\",
  \"content\": \"<h1>Welcome to WebEnable CMS</h1><p>This is your first blog post! WebEnable CMS is a modern, fast, and secure content management system built with Go and Next.js.</p><p>Features include:</p><ul><li>Modern responsive design</li><li>SEO optimization</li><li>Fast performance</li><li>Secure architecture</li><li>Easy content management</li></ul><p>Start creating amazing content today!</p>\",
  \"excerpt\": \"Welcome to WebEnable CMS - a modern, fast, and secure content management system. Start creating amazing content today!\",
  \"author\": \"admin\",
  \"status\": \"published\",
  \"tags\": [\"welcome\", \"getting-started\", \"cms\"],
  \"categories\": [\"General\", \"Announcements\"],
  \"featured_image\": \"/images/welcome-banner.jpg\",
  \"image_alt\": \"WebEnable CMS Welcome Banner\",
  \"meta_title\": \"Welcome to WebEnable CMS - Modern Content Management\",
  \"meta_description\": \"Get started with WebEnable CMS, a modern and secure content management system built for performance and ease of use.\",
  \"reading_time\": 2,
  \"is_featured\": true,
  \"view_count\": 0,
  \"created_at\": \"$(current_timestamp)\",
  \"updated_at\": \"$(current_timestamp)\",
  \"published_at\": \"$(current_timestamp)\"
}"

# Post 2: Getting Started Guide
POST_ID_2=$(generate_uuid)
curl -X POST "$POSTS_DB" -H "Content-Type: application/json" -d "{
  \"_id\": \"$POST_ID_2\",
  \"title\": \"Getting Started with Your New Website\",
  \"content\": \"<h1>Getting Started Guide</h1><p>Congratulations on setting up your new website! Here's a quick guide to get you started:</p><h2>Admin Panel</h2><p>Access your admin panel at <code>/admin</code> to manage your content, users, and settings.</p><h2>Creating Content</h2><p>Use the intuitive editor to create blog posts, pages, and manage your content easily.</p><h2>Customization</h2><p>Customize your site's appearance, SEO settings, and more through the admin interface.</p><h2>Support</h2><p>Need help? Check out our documentation or contact our support team.</p>\",
  \"excerpt\": \"Learn how to get started with your new WebEnable CMS website. A comprehensive guide for beginners.\",
  \"author\": \"admin\",
  \"status\": \"published\",
  \"tags\": [\"guide\", \"tutorial\", \"getting-started\"],
  \"categories\": [\"Tutorials\", \"Help\"],
  \"featured_image\": \"/images/getting-started.jpg\",
  \"image_alt\": \"Getting Started with WebEnable CMS\",
  \"meta_title\": \"Getting Started Guide - WebEnable CMS\",
  \"meta_description\": \"Complete getting started guide for your new WebEnable CMS website. Learn the basics and start creating content.\",
  \"reading_time\": 3,
  \"is_featured\": false,
  \"view_count\": 0,
  \"created_at\": \"$(current_timestamp)\",
  \"updated_at\": \"$(current_timestamp)\",
  \"published_at\": \"$(current_timestamp)\"
}"

# Post 3: Features Overview
POST_ID_3=$(generate_uuid)
curl -X POST "$POSTS_DB" -H "Content-Type: application/json" -d "{
  \"_id\": \"$POST_ID_3\",
  \"title\": \"WebEnable CMS Features Overview\",
  \"content\": \"<h1>Feature Overview</h1><p>WebEnable CMS comes packed with powerful features:</p><h2>🚀 Performance</h2><p>Built with Go and Next.js for lightning-fast performance and scalability.</p><h2>🔒 Security</h2><p>Enterprise-grade security with JWT authentication, rate limiting, and secure headers.</p><h2>📱 Responsive Design</h2><p>Mobile-first design that looks great on all devices.</p><h2>🔍 SEO Optimized</h2><p>Built-in SEO features including meta tags, schema markup, and sitemap generation.</p><h2>📊 Analytics Ready</h2><p>Easy integration with Google Analytics and other tracking services.</p><h2>🛠 Easy Management</h2><p>Intuitive admin interface for content management, user administration, and site settings.</p>\",
  \"excerpt\": \"Discover the powerful features that make WebEnable CMS the perfect choice for your website.\",
  \"author\": \"admin\",
  \"status\": \"published\",
  \"tags\": [\"features\", \"overview\", \"capabilities\"],
  \"categories\": [\"Features\", \"Information\"],
  \"featured_image\": \"/images/features-overview.jpg\",
  \"image_alt\": \"WebEnable CMS Features\",
  \"meta_title\": \"Features Overview - WebEnable CMS Capabilities\",
  \"meta_description\": \"Explore the comprehensive features of WebEnable CMS including performance, security, SEO optimization, and more.\",
  \"reading_time\": 4,
  \"is_featured\": true,
  \"view_count\": 0,
  \"created_at\": \"$(current_timestamp)\",
  \"updated_at\": \"$(current_timestamp)\",
  \"published_at\": \"$(current_timestamp)\"
}"

# Create sample contacts
echo "📞 Creating sample contacts..."

# Contact 1: Welcome inquiry
CONTACT_ID_1=$(generate_uuid)
curl -X POST "$CONTACTS_DB" -H "Content-Type: application/json" -d "{
  \"_id\": \"$CONTACT_ID_1\",
  \"name\": \"John Smith\",
  \"email\": \"john.smith@example.com\",
  \"subject\": \"Welcome and Questions\",
  \"message\": \"Hi! I just discovered your website and I'm impressed with the design and functionality. I have a few questions about your services. Could you please get back to me when you have a chance? Thank you!\",
  \"status\": \"new\",
  \"created_at\": \"$(current_timestamp)\",
  \"updated_at\": \"$(current_timestamp)\"
}"

# Contact 2: Business inquiry
CONTACT_ID_2=$(generate_uuid)
curl -X POST "$CONTACTS_DB" -H "Content-Type: application/json" -d "{
  \"_id\": \"$CONTACT_ID_2\",
  \"name\": \"Sarah Johnson\",
  \"email\": \"sarah.johnson@business.com\",
  \"subject\": \"Business Partnership Inquiry\",
  \"message\": \"Hello, I represent a growing technology company and we're interested in exploring potential partnership opportunities. Would you be available for a brief call next week to discuss this further?\",
  \"status\": \"new\",
  \"created_at\": \"$(current_timestamp)\",
  \"updated_at\": \"$(current_timestamp)\"
}"

echo "✅ Sample data populated successfully!"
echo ""
echo "📊 Summary:"
echo "   • 3 sample blog posts created"
echo "   • 2 sample contacts created"
echo ""
echo "🌐 You can now view your website with sample content!"
echo "   • Frontend: http://localhost"
echo "   • Admin Panel: http://localhost/admin"
echo ""

# Contact Population Script

## Overview

The `populate_contacts.sh` script automatically populates your WebEnable CMS database with realistic sample contact form submissions for testing and demonstration purposes.

## Features

- ✅ **Environment Validation**: Verifies .env configuration and CouchDB connectivity
- 📝 **Sample Contacts**: Creates 15+ realistic contact form submissions
- 🏢 **Business Scenarios**: Includes various business inquiries (web development, e-commerce, healthcare, fintech, etc.)
- 📊 **Status Variety**: Mixes contact statuses (new, read, replied) with realistic timestamps
- 🔍 **Duplicate Detection**: Skips existing contacts to prevent duplicates
- 📈 **Progress Reporting**: Shows detailed summary of created contacts

## Usage

### Quick Start
```bash
# From project root directory
./populate_contacts.sh
```

### Manual Execution
```bash
# Make executable (first time only)
chmod +x populate_contacts.sh

# Run the script
./populate_contacts.sh
```

## Sample Contact Types

The script creates contacts from various business scenarios:

### **Business Inquiries**
- Custom web application development
- E-commerce platform setup
- Patient portal development (healthcare)
- IoT dashboard development
- Learning management systems
- Financial analytics dashboards
- Inventory management systems
- Freelancer marketplace platforms
- Travel booking platforms
- Fitness tracking applications

### **Contact Statuses**
- **New**: Unread contact submissions
- **Read**: Contacts that have been viewed by admin
- **Replied**: Contacts that have received responses

## Prerequisites

1. **Running Services**: Ensure Docker services are running
   ```bash
   docker compose ps
   ```

2. **Environment File**: Valid `.env` file with `COUCHDB_PASSWORD`

3. **CouchDB Access**: Database accessible on `localhost:5984`

4. **Go Runtime**: Available for script execution (via backend container context)

## Output Example

```bash
🚀 Populating CouchDB with sample contacts...
✅ Loaded CouchDB password from .env
🔍 Verifying CouchDB connection...
✅ CouchDB is accessible
📝 Running populate_contacts.go script...

Added contact: Sarah Johnson (TechCorp Solutions) - Custom Web Application Development
Added contact: Michael Chen (InnovateStart) - E-commerce Platform Setup
Added contact: Emma Rodriguez (HealthPlus Clinic) - Patient Portal Development
...

✅ Contacts populated successfully!

📊 Database Summary:
   📞 Total contacts in database: 22
   📋 Contact status breakdown:
      - New: 15
      - Read: 5
      - Replied: 2

🔗 Access contacts via:
   📱 Admin Panel: http://localhost/admin/contacts
   🔧 API Endpoint: http://localhost/api/contacts (requires authentication)
   💾 Direct DB: http://localhost:5984/_utils/#database/contacts/_all_docs
```

## Accessing Populated Contacts

### **Admin Panel**
Visit: `http://localhost/admin/contacts`
- View all contacts with filtering and sorting
- Mark contacts as read/replied
- Respond to contact inquiries

### **API Endpoint**
Access: `http://localhost/api/contacts`
- Requires admin authentication
- Returns JSON format with pagination
- Supports filtering by status

### **Direct Database**
CouchDB Interface: `http://localhost:5984/_utils/#database/contacts/_all_docs`
- Direct database access (admin credentials required)
- Raw document view
- Full database management capabilities

## Troubleshooting

### **Script Fails to Start**
```bash
# Check if .env file exists
ls -la .env

# Verify CouchDB password
grep COUCHDB_PASSWORD .env
```

### **CouchDB Connection Issues**
```bash
# Test CouchDB connectivity
curl http://admin:YOUR_PASSWORD@localhost:5984

# Check Docker services
docker compose ps
docker compose logs db
```

### **Go Module Issues**
```bash
# Ensure backend/go.mod exists
ls -la backend/go.mod

# Run from correct directory
pwd  # Should be in project root
```

### **Permission Issues**
```bash
# Make script executable
chmod +x populate_contacts.sh

# Check file permissions
ls -la populate_contacts.sh
```

## Integration with Production Setup

This script is automatically included in the production setup guide and can be run as part of the initial deployment process:

```bash
# During production setup (Step 5.2)
./populate_contacts.sh
```

The script is safe to run multiple times - it will skip existing contacts and only add new ones.

## Sample Contact Data

The script includes realistic business inquiries such as:

- **TechCorp Solutions**: Custom web application development
- **InnovateStart**: E-commerce platform setup  
- **HealthPlus Clinic**: HIPAA-compliant patient portal
- **GreenTech Solutions**: IoT environmental monitoring dashboard
- **EduLearn Platform**: Scalable learning management system
- **FinancePlus**: Secure financial analytics dashboard
- **RetailChain Inc**: Multi-location inventory management
- **FreelanceHub**: Marketplace platform with payment processing
- **Travel Explore**: International travel booking platform
- **SportsTracker**: Fitness tracking with wearable integration

Each contact includes realistic company information, phone numbers, detailed project descriptions, and appropriate timestamps.

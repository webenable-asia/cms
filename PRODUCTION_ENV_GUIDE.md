# Production Environment Setup Guide

## 🚀 Quick Production Setup

### Automated Setup (Recommended)
```bash
# Run the interactive production environment setup
./scripts/setup-production-env.sh
```

This script will:
- ✅ Generate secure secrets and passwords
- ✅ Configure domain settings
- ✅ Create .env.production file
- ✅ Update Caddyfile for your domain
- ✅ Set secure file permissions
- ✅ Create deployment checklist

### Manual Setup
If you prefer manual configuration, copy and customize the `.env` file:

```bash
# Copy the generated production environment
cp .env .env.production

# Edit with your production values
nano .env.production

# Set secure permissions
chmod 600 .env.production
```

## 🔐 Security Validation

Before deploying, validate your environment security:

```bash
# Check your production environment file
./scripts/validate-env-security.sh .env.production
```

This will check for:
- Strong passwords and secrets
- Secure file permissions
- Production domain configuration
- HTTPS configuration
- Session security settings

## 📋 Production Environment Variables

### Critical Security Settings
```bash
# Generate secure values with:
JWT_SECRET=$(openssl rand -base64 32)
ADMIN_PASSWORD=$(openssl rand -base64 16)
COUCHDB_PASSWORD=$(openssl rand -base64 24)
VALKEY_PASSWORD=$(openssl rand -base64 24)
```

### Domain Configuration
Replace placeholder domains with your actual domains:
```bash
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
NEXT_PUBLIC_API_URL=https://yourdomain.com/api
SESSION_DOMAIN=yourdomain.com
```

### Production Optimizations
```bash
NODE_ENV=production
GO_ENV=production
SESSION_SECURE=true
LOG_LEVEL=info
```

## 🛠️ Deployment Steps

### 1. Server Preparation
```bash
# Install Podman
sudo apt update
sudo apt install podman git

# Create application user
sudo useradd -m -s /bin/bash cms
```

### 2. Application Deployment
```bash
# Clone repository
git clone <your-repository-url> /opt/cms
cd /opt/cms

# Copy production environment
cp .env.production .env

# Set ownership and permissions
sudo chown -R cms:cms /opt/cms
chmod 600 .env
```

### 3. Start Services
```bash
# Switch to cms user
sudo -u cms bash

# Start production services
podman compose up -d

# Initialize admin user
podman compose exec backend go run scripts/init_admin.go
```

### 4. Verify Deployment
```bash
# Check service status
podman compose ps

# Check logs
podman compose logs

# Test health endpoint
curl https://yourdomain.com/api/health
```

## 🔧 Configuration Files

### Environment File Structure
The production `.env` file contains:

- **Security Settings**: JWT secrets, passwords, admin credentials
- **Database Config**: CouchDB connection and credentials
- **Cache Config**: Valkey (Redis) settings
- **Server Settings**: CORS, session, logging
- **Frontend Config**: API URLs, Node.js settings
- **Email Config**: SMTP settings for notifications
- **Performance**: Cache TTL, rate limiting
- **Monitoring**: Logging and health check settings

### Caddyfile Updates
The setup script automatically updates your Caddyfile:
- Replaces `localhost` with your production domain
- Enables automatic HTTPS
- Configures SSL certificate generation

## 🔒 Security Best Practices

### Secrets Management
- ✅ Use generated secure passwords (minimum 16 characters)
- ✅ Set file permissions to 600 (owner read/write only)
- ✅ Never commit `.env` files to version control
- ✅ Rotate secrets quarterly
- ✅ Use different passwords for each service

### Domain Security
- ✅ Use HTTPS for all production URLs
- ✅ Configure proper CORS origins
- ✅ Set secure session configuration
- ✅ Enable HSTS headers

### Server Security
- ✅ Run services as non-root user
- ✅ Configure firewall (ports 80, 443, 5984 only)
- ✅ Keep server and packages updated
- ✅ Enable fail2ban or similar intrusion detection

## 📊 Monitoring & Maintenance

### Health Checks
```bash
# API health
curl https://yourdomain.com/api/health

# Service status
podman compose ps

# Resource usage
podman stats
```

### Log Management
```bash
# View application logs
podman compose logs -f

# View specific service logs
podman compose logs backend
podman compose logs frontend
```

### Backup Strategy
```bash
# Backup database
podman exec cms-db-1 curl -X GET http://admin:password@localhost:5984/_all_dbs

# Backup volumes
podman run --rm -v cms_couchdb_data:/data -v $(pwd):/backup alpine \
  tar czf /backup/couchdb-$(date +%Y%m%d).tar.gz /data
```

## 🚨 Troubleshooting

### Common Issues

#### SSL Certificate Issues
```bash
# Check Caddy logs
podman compose logs caddy

# Verify domain DNS
nslookup yourdomain.com

# Test port accessibility
telnet yourdomain.com 80
```

#### Database Connection Issues
```bash
# Check database status
podman compose logs db

# Test database connection
curl http://admin:password@localhost:5984/_all_dbs
```

#### Performance Issues
```bash
# Monitor resource usage
podman stats

# Check cache status
podman compose exec cache redis-cli ping
```

### Emergency Procedures

#### Service Recovery
```bash
# Restart all services
podman compose restart

# Rebuild if needed
podman compose up -d --force-recreate
```

#### Security Incident Response
1. Immediately rotate all secrets
2. Check logs for suspicious activity
3. Update environment with new secrets
4. Restart all services
5. Monitor for continued issues

## 📚 Additional Resources

### Scripts Available
- `setup-production-env.sh` - Interactive production setup
- `validate-env-security.sh` - Security validation
- `manage.sh` - Service management
- `scripts/cleanup.sh` - Cleanup and maintenance

### Documentation
- `PODMAN.md` - Podman development guide
- `PRODUCTION_DEPLOYMENT.md` - Detailed deployment guide
- `SECURITY_CHECKLIST.md` - Security implementation checklist

---

**Note**: Always test your production configuration in a staging environment first before deploying to production.

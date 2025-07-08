# Production Deployment Checklist

## Pre-Deployment Checklist

### 🔧 Infrastructure Setup
- [ ] **Domain Name**: Registered and DNS configured to point to server IP
- [ ] **Server Requirements**: Minimum 4GB RAM, 2 vCPU, 40GB storage
- [ ] **Podman Installation**: Podman and Podman Compose installed on server
- [ ] **Firewall Configuration**: Ports 80, 443, and 5984 open
- [ ] **SSH Access**: Secure SSH key-based access configured

### 🔐 Security Configuration
- [ ] **Environment Variables**: All secrets configured in `.env.production`
- [ ] **Strong Passwords**: Database and cache passwords are complex and unique
- [ ] **JWT Secret**: Cryptographically secure JWT secret generated
- [ ] **Git Security**: No sensitive files tracked in version control
- [ ] **File Permissions**: Proper file permissions set on server

### 🌐 SSL/TLS Setup
- [ ] **Domain Configuration**: Replace `localhost` with actual domain in `caddy/Caddyfile`
- [ ] **HTTPS Enablement**: Set `auto_https on` in Caddyfile (or remove `auto_https off`)
- [ ] **Certificate Validation**: Verify automatic certificate generation works

### 📊 Performance Optimization
- [ ] **Resource Limits**: Podman container resource limits configured
- [ ] **Database Indexing**: CouchDB views and indexes optimized
- [ ] **Cache Configuration**: Valkey cache properly configured with TTL
- [ ] **Compression**: Gzip/Zstd compression enabled
- [ ] **Static Assets**: CDN or efficient static file serving configured

## Deployment Process

### 1. Server Preparation
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install podman git -y

# Create application user
sudo useradd -m -s /bin/bash cms
```

### 2. Application Deployment
```bash
# Clone repository
git clone <your-repository-url> /opt/cms
cd /opt/cms

# Set ownership
sudo chown -R cms:cms /opt/cms

# Switch to cms user
sudo -u cms bash

# Configure environment
cp .env.example .env.production
nano .env.production  # Edit with production values

# Update Caddyfile for production
nano caddy/Caddyfile  # Replace localhost with your domain

# Start services
podman compose up -d

# Initialize admin user
podman compose exec backend ./main init-admin
```

### 3. Post-Deployment Verification
```bash
# Run production readiness test
./scripts/production-readiness-test.sh

# Check service health
podman compose ps
podman compose logs

# Test application functionality
curl https://yourdomain.com
curl https://yourdomain.com/api/posts
curl https://yourdomain.com:5984
```

## Production Environment Variables

### Required Variables
```bash
# Database Configuration
COUCHDB_USER=admin
COUCHDB_PASSWORD=<strong-password>
COUCHDB_DATABASE=cms

# Cache Configuration  
VALKEY_PASSWORD=<strong-password>

# Authentication
JWT_SECRET=<cryptographically-secure-secret>

# Frontend Configuration
NEXT_PUBLIC_API_URL=https://yourdomain.com/api

# Email Configuration (if used)
SMTP_HOST=smtp.yourdomain.com
SMTP_PORT=587
SMTP_USER=noreply@yourdomain.com
SMTP_PASS=<smtp-password>
```

### Generating Secure Values
```bash
# Generate JWT secret (64 characters)
openssl rand -hex 32

# Generate strong passwords
openssl rand -base64 32
```

## Monitoring & Maintenance

### 📊 Health Monitoring
- [ ] **Service Health**: Monitor container health status
- [ ] **Resource Usage**: Track CPU, memory, and disk usage
- [ ] **Response Times**: Monitor application response times
- [ ] **Error Rates**: Track application error rates and logs

### 🔄 Backup Strategy
- [ ] **Database Backups**: Automated CouchDB backups configured
- [ ] **Configuration Backups**: Environment and configuration files backed up
- [ ] **Volume Backups**: Podman volumes regularly backed up
- [ ] **Recovery Testing**: Backup restoration process tested

### 🛡️ Security Monitoring
- [ ] **Access Logs**: Monitor access patterns and suspicious activity
- [ ] **Security Updates**: Regular system and container image updates
- [ ] **Certificate Renewal**: Monitor SSL certificate expiration
- [ ] **Vulnerability Scanning**: Regular security vulnerability scans

### 📈 Performance Optimization
- [ ] **Database Performance**: Monitor query performance and optimize indexes
- [ ] **Cache Hit Rates**: Monitor and optimize cache performance
- [ ] **Resource Scaling**: Plan for resource scaling based on usage
- [ ] **CDN Configuration**: Configure CDN for static assets if needed

## Troubleshooting

### Common Issues

#### 1. SSL Certificate Issues
```bash
# Check Caddy logs
podman logs cms-caddy-1

# Verify domain DNS
nslookup yourdomain.com

# Test HTTP challenge
curl -I http://yourdomain.com/.well-known/acme-challenge/test
```

#### 2. Database Connection Issues
```bash
# Check database health
podman logs cms-db-1

# Test database connectivity
curl http://localhost:5984

# Check database authentication
curl -u admin:password http://localhost:5984/_all_dbs
```

#### 3. Performance Issues
```bash
# Check resource usage
podman stats

# Check application logs
podman compose logs backend
podman compose logs frontend

# Monitor database performance
curl -u admin:password http://localhost:5984/_utils/
```

## Security Best Practices

### 🔐 Access Control
- Use strong, unique passwords for all services
- Implement IP-based restrictions for admin interfaces
- Regular access review and user management
- Enable two-factor authentication where possible

### 🛡️ Network Security
- Configure firewall rules (UFW recommended)
- Use VPN for administrative access
- Implement rate limiting for public endpoints
- Monitor and log all access attempts

### 📋 Regular Maintenance
- Keep system packages updated
- Update Docker images regularly
- Review and rotate secrets periodically
- Monitor security advisories for dependencies

## Support & Documentation

- **Production Readiness Test**: `./scripts/production-readiness-test.sh`
- **Caddy Proxy Test**: `./scripts/test-caddy-proxy.sh`
- **Management Helper**: `./manage.sh help`
- **Architecture Documentation**: `docs/REVERSE_PROXY.md`
- **Security Checklist**: `SECURITY_CHECKLIST.md`

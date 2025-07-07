# VM Deployment Guide

This guide provides step-by-step instructions for deploying the CMS application on a VM instance in production.

## Prerequisites

- VM instance with Docker and Docker Compose installed
- Domain name configured (optional but recommended)
- SSL certificates (for HTTPS)
- Environment variables configured

## VM Requirements

### Minimum Specifications
- **CPU**: 2 vCPUs
- **RAM**: 4GB
- **Storage**: 40GB SSD
- **Network**: 10 Mbps bandwidth
- **OS**: Ubuntu 22.04 LTS or similar

### Recommended Specifications
- **CPU**: 4 vCPUs
- **RAM**: 8GB
- **Storage**: 80GB SSD
- **Network**: 100 Mbps bandwidth
- **OS**: Ubuntu 22.04 LTS

## Deployment Steps

### 1. Server Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install docker-compose-plugin -y

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Install additional tools
sudo apt install git nginx-utils htop -y
```

### 2. Application Deployment

```bash
# Clone repository
git clone <repository-url> cms-app
cd cms-app

# Create environment file
cp .env.example .env.production
nano .env.production

# Create SSL directory
mkdir -p nginx/ssl

# Generate self-signed SSL (for testing)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/privkey.pem \
  -out nginx/ssl/fullchain.pem

# Build and start services
docker compose -f docker-compose.prod.yml up -d --build
```

### 3. Environment Configuration

Create `.env.production` with the following variables:

```env
# Database Configuration
COUCHDB_USER=admin
COUCHDB_PASSWORD=your_secure_password_here

# Cache Configuration
VALKEY_PASSWORD=your_cache_password_here

# API Configuration
NEXT_PUBLIC_API_URL=https://yourdomain.com/api
BACKEND_URL=http://backend:8080

# Security
JWT_SECRET=your_jwt_secret_here
ENCRYPTION_KEY=your_32_character_encryption_key

# Email Configuration (optional)
SMTP_HOST=smtp.yourdomain.com
SMTP_PORT=587
SMTP_USER=noreply@yourdomain.com
SMTP_PASS=your_smtp_password
```

### 4. SSL Certificate Setup

For production, replace self-signed certificates with proper SSL certificates:

```bash
# Using Let's Encrypt with Certbot
sudo apt install certbot -y

# Generate certificates
sudo certbot certonly --standalone -d yourdomain.com

# Copy certificates
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem nginx/ssl/
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem nginx/ssl/

# Set permissions
sudo chown -R $USER:$USER nginx/ssl/
chmod 600 nginx/ssl/privkey.pem
chmod 644 nginx/ssl/fullchain.pem

# Restart nginx service
docker compose -f docker-compose.prod.yml restart nginx
```

### 5. Firewall Configuration

```bash
# Enable UFW
sudo ufw enable

# Allow SSH
sudo ufw allow ssh

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Check status
sudo ufw status
```

## Monitoring and Maintenance

### Health Checks

```bash
# Check service status
docker compose -f docker-compose.prod.yml ps

# Check logs
docker compose -f docker-compose.prod.yml logs -f

# Check resource usage
docker stats

# Check system resources
htop
```

### Backup Strategy

```bash
# Create backup script
cat > backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup CouchDB
docker exec cms-db-1 couchdb-backup > $BACKUP_DIR/couchdb_$DATE.backup

# Backup Valkey
docker exec cms-cache-1 valkey-cli --rdb /data/dump_$DATE.rdb BGSAVE

# Backup application files
tar -czf $BACKUP_DIR/app_$DATE.tar.gz \
  --exclude='node_modules' \
  --exclude='.git' \
  /path/to/cms-app/

# Clean old backups (keep last 7 days)
find $BACKUP_DIR -name "*.backup" -mtime +7 -delete
find $BACKUP_DIR -name "*.rdb" -mtime +7 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
EOF

chmod +x backup.sh

# Add to crontab (daily backup at 2 AM)
(crontab -l 2>/dev/null; echo "0 2 * * * /path/to/backup.sh") | crontab -
```

### Log Rotation

```bash
# Configure Docker log rotation
sudo nano /etc/docker/daemon.json

# Add log configuration
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}

# Restart Docker
sudo systemctl restart docker
```

### Auto-Update Script

```bash
# Create update script
cat > update.sh << 'EOF'
#!/bin/bash
cd /path/to/cms-app

# Pull latest changes
git pull origin main

# Rebuild and restart services
docker compose -f docker-compose.prod.yml up -d --build

# Clean unused images
docker image prune -f
EOF

chmod +x update.sh
```

## Scaling Considerations

### Horizontal Scaling

For high traffic scenarios, consider:

1. **Load Balancer**: Use external load balancer (AWS ALB, Google Cloud Load Balancer)
2. **Database Clustering**: CouchDB cluster setup
3. **Cache Clustering**: Valkey cluster configuration
4. **Multiple VM Instances**: Deploy across multiple regions

### Vertical Scaling

Resource upgrade path:
- **Light Traffic**: 2 vCPU, 4GB RAM
- **Medium Traffic**: 4 vCPU, 8GB RAM  
- **High Traffic**: 8 vCPU, 16GB RAM
- **Enterprise**: 16+ vCPU, 32+ GB RAM

## Troubleshooting

### Common Issues

1. **Service Won't Start**
   ```bash
   # Check logs
   docker compose -f docker-compose.prod.yml logs service_name
   
   # Check resource usage
   docker stats
   ```

2. **Out of Memory**
   ```bash
   # Check memory usage
   free -h
   docker stats
   
   # Increase swap
   sudo fallocate -l 2G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```

3. **SSL Certificate Issues**
   ```bash
   # Check certificate validity
   openssl x509 -in nginx/ssl/fullchain.pem -text -noout
   
   # Test SSL connection
   openssl s_client -connect yourdomain.com:443
   ```

4. **Database Connection Issues**
   ```bash
   # Check CouchDB status
   docker exec cms-db-1 curl -X GET http://admin:password@localhost:5984/
   
   # Check network connectivity
   docker exec cms-backend-1 curl -X GET http://db:5984/
   ```

### Performance Optimization

1. **Database Optimization**
   ```bash
   # Compact database
   docker exec cms-db-1 curl -X POST http://admin:password@localhost:5984/your_db/_compact
   
   # Create indexes
   docker exec cms-db-1 curl -X POST http://admin:password@localhost:5984/your_db/_index \
     -H "Content-Type: application/json" \
     -d '{"index": {"fields": ["created_at"]}, "name": "created-index"}'
   ```

2. **Cache Optimization**
   ```bash
   # Monitor cache hit ratio
   docker exec cms-cache-1 valkey-cli info stats | grep keyspace
   
   # Optimize memory usage
   docker exec cms-cache-1 valkey-cli config set maxmemory-policy allkeys-lru
   ```

## Security Best Practices

1. **Regular Updates**
   - Keep OS updated
   - Update Docker images regularly
   - Update application dependencies

2. **Access Control**
   - Use strong passwords
   - Enable 2FA where possible
   - Limit SSH access
   - Use VPN for admin access

3. **Network Security**
   - Configure firewall rules
   - Use private networks
   - Enable SSL/TLS
   - Regular security audits

4. **Data Protection**
   - Encrypt data at rest
   - Regular backups
   - Test restore procedures
   - Monitor access logs

## Support

For issues and questions:
1. Check application logs
2. Review this deployment guide
3. Consult Docker documentation
4. Contact system administrator

---

**Last Updated**: $(date +%Y-%m-%d)
**Version**: 1.0

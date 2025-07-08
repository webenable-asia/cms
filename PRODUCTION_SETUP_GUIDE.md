# WebEnable CMS Production Setup Guide
## Complete Step-by-Step Deployment Instructions

*Updated: July 9, 2025 | Version: 2.0*

---

## 🎯 **Overview**

This guide provides complete step-by-step instructions for deploying WebEnable CMS to production using Podman (recommended) or Docker. The setup includes security hardening, performance optimization, and monitoring.

---

## 📋 **Prerequisites**

### **System Requirements:**
- **OS**: Ubuntu 22.04 LTS / CentOS 8+ / RHEL 8+ (recommended)
- **CPU**: 2+ cores (4+ recommended)
- **RAM**: 4GB minimum (8GB+ recommended)
- **Storage**: 20GB+ available space
- **Network**: Static IP address and domain name

### **Software Prerequisites:**
- Root or sudo access
- Git installed
- Text editor (nano/vim)
- Basic Linux command line knowledge

---

## 🔧 **Step 1: System Preparation**

### **1.1 Update System**
```bash
# Ubuntu/Debian
sudo apt update && sudo apt upgrade -y

# CentOS/RHEL/Fedora
sudo dnf update -y
```

### **1.2 Install Required Packages**
```bash
# Ubuntu/Debian
sudo apt install -y curl wget git unzip firewalld fail2ban

# CentOS/RHEL/Fedora
sudo dnf install -y curl wget git unzip firewalld fail2ban
```

### **1.3 Configure Firewall**
```bash
# Enable firewall
sudo systemctl enable --now firewalld

# Allow required ports
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --permanent --add-port=22/tcp
sudo firewall-cmd --reload
```

### **1.4 Setup Fail2Ban**
```bash
# Enable fail2ban
sudo systemctl enable --now fail2ban

# Create jail configuration
sudo tee /etc/fail2ban/jail.d/webenable.conf << EOF
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5

[sshd]
enabled = true
port = ssh
logpath = %(sshd_log)s
backend = %(sshd_backend)s
EOF

sudo systemctl restart fail2ban
```

---

## 🐳 **Step 2: Container Runtime Installation**

### **Option A: Podman Installation (Recommended)**

#### **2A.1 Install Podman**
```bash
# Ubuntu 22.04+
sudo apt install -y podman podman-compose

# CentOS/RHEL 8+
sudo dnf install -y podman podman-compose

# Verify installation
podman --version
```

#### **2A.2 Configure Podman**
```bash
# Enable lingering for current user
sudo loginctl enable-linger $USER

# Configure registries
sudo tee /etc/containers/registries.conf << EOF
[registries.search]
registries = ['docker.io', 'registry.redhat.io']

[registries.insecure]
registries = []

[registries.block]
registries = []
EOF
```

### **Option B: Docker Installation (Alternative)**

#### **2B.1 Install Docker**
```bash
# Ubuntu
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Start and enable Docker
sudo systemctl enable --now docker
```

---

## 📁 **Step 3: Application Deployment**

### **3.1 Create Application User**
```bash
# Create dedicated user for security
sudo useradd -m -s /bin/bash webenable
sudo usermod -aG wheel webenable  # CentOS/RHEL
sudo usermod -aG sudo webenable   # Ubuntu

# Switch to application user
sudo su - webenable
```

### **3.2 Clone Repository**
```bash
# Clone the WebEnable CMS repository
cd /home/webenable
git clone https://github.com/your-org/webenable-asia.git
cd webenable-asia

# Verify project structure
ls -la
```

### **3.3 Configure Environment**
```bash
# Copy environment template
cp .env.example .env

# Generate secure environment variables
nano .env
```

### **3.4 Environment Configuration**
Edit `.env` file with production values:

```bash
# ===== PRODUCTION ENVIRONMENT CONFIGURATION =====

# Basic Configuration
NODE_ENV=production
PORT=8080
ENVIRONMENT=production

# Domain Configuration
DOMAIN=yourdomain.com
NEXT_PUBLIC_API_URL=https://yourdomain.com/api
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Security Configuration (Generate new secrets!)
JWT_SECRET=$(openssl rand -base64 32)
ADMIN_PASSWORD=$(openssl rand -base64 16)
SESSION_SECRET=$(openssl rand -base64 32)

# Database Configuration
COUCHDB_USER=admin
COUCHDB_PASSWORD=$(openssl rand -base64 16)
COUCHDB_URL=http://admin:YOUR_DB_PASSWORD@db:5984/

# Cache Configuration
VALKEY_PASSWORD=$(openssl rand -base64 16)
VALKEY_URL=redis://:YOUR_VALKEY_PASSWORD@cache:6379

# Email Configuration (Configure with your SMTP provider)
SMTP_HOST=smtp.yourdomain.com
SMTP_PORT=587
SMTP_USER=noreply@yourdomain.com
SMTP_PASSWORD=your_smtp_password
SMTP_FROM=noreply@yourdomain.com

# Security Headers
SESSION_SECURE=true
SESSION_DOMAIN=yourdomain.com

# Rate Limiting
RATE_LIMIT_WINDOW=900000
RATE_LIMIT_MAX=100

# Logging
LOG_LEVEL=info
LOG_FILE=/var/log/webenable/app.log
```

---

## 🚀 **Step 4: SSL/TLS Certificate Setup**

### **4.1 Install Certbot**
```bash
# Ubuntu
sudo apt install -y certbot

# CentOS/RHEL
sudo dnf install -y certbot
```

### **4.2 Obtain SSL Certificate**
```bash
# Stop any running web servers
sudo systemctl stop nginx apache2 2>/dev/null || true

# Obtain certificate
sudo certbot certonly --standalone \
  --email admin@yourdomain.com \
  --agree-tos \
  --no-eff-email \
  -d yourdomain.com \
  -d www.yourdomain.com

# Setup auto-renewal
sudo crontab -e
# Add this line:
# 0 12 * * * /usr/bin/certbot renew --quiet
```

### **4.3 Configure Caddy for Production**
Update `caddy/Caddyfile`:

```bash
# Production Caddyfile
{
    admin off
    auto_https on
    
    # Global rate limiting
    servers {
        metrics
    }
}

# Main site configuration
yourdomain.com, www.yourdomain.com {
    # Security headers
    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        X-Frame-Options DENY
        X-Content-Type-Options nosniff
        X-XSS-Protection "1; mode=block"
        Referrer-Policy strict-origin-when-cross-origin
        Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'"
        Permissions-Policy "geolocation=(), microphone=(), camera=()"
        -Server
    }

    # API routes
    handle /api/* {
        reverse_proxy backend:8080 {
            header_up Host {upstream_hostport}
            header_up X-Real-IP {remote_ip}
            header_up X-Forwarded-For {remote_ip}
            header_up X-Forwarded-Proto {scheme}
            health_uri /api/health
            health_interval 30s
        }
    }

    # Frontend routes
    handle {
        reverse_proxy frontend:8080 {
            header_up Host {upstream_hostport}
            header_up X-Real-IP {remote_ip}
            header_up X-Forwarded-For {remote_ip}
            header_up X-Forwarded-Proto {scheme}
        }
    }

    # Enable compression
    encode gzip zstd

    # Rate limiting
    rate_limit {
        zone static_ip {
            key {remote_ip}
            events 100
            window 1m
        }
    }

    # Logging
    log {
        output file /var/log/caddy/access.log {
            roll_size 100mb
            roll_keep 5
        }
        format json
    }
}

# Database proxy (restrict access)
api.yourdomain.com:5984 {
    # Restrict to internal networks only
    @internal {
        remote_ip 10.0.0.0/8 172.16.0.0/12 192.168.0.0/16
    }
    
    handle @internal {
        reverse_proxy db:5984
    }
    
    handle {
        respond "Access denied" 403
    }
}
```

---

## 🔒 **Step 5: Production Deployment**

### **5.1 Build and Start Services**
```bash
# Using Podman (Recommended)
cd /home/webenable/webenable-asia

# Build images
podman compose build

# Start services
podman compose up -d

# Verify all services are running
podman compose ps
```

### **5.2 Initialize Database**
```bash
# Wait for services to be ready (30 seconds)
sleep 30

# Initialize admin user
cd backend
go run ./scripts/init_admin.go

# Populate sample data (optional)
go run ../scripts/populate_database.go
```

### **5.3 Verify Deployment**
```bash
# Check service health
curl -f http://localhost/api/health

# Check frontend
curl -f http://localhost/

# Check logs
podman compose logs --tail=50
```

---

## 📊 **Step 6: Monitoring and Logging**

### **6.1 Setup Log Directories**
```bash
# Create log directories
sudo mkdir -p /var/log/webenable
sudo mkdir -p /var/log/caddy
sudo chown webenable:webenable /var/log/webenable
sudo chown webenable:webenable /var/log/caddy
```

### **6.2 Configure Logrotate**
```bash
sudo tee /etc/logrotate.d/webenable << EOF
/var/log/webenable/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    copytruncate
}

/var/log/caddy/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    copytruncate
}
EOF
```

### **6.3 Setup Health Monitoring**
```bash
# Create health check script
sudo tee /usr/local/bin/webenable-health << 'EOF'
#!/bin/bash

HEALTH_URL="http://localhost/api/health"
LOG_FILE="/var/log/webenable/health.log"

if curl -f -s "$HEALTH_URL" > /dev/null; then
    echo "$(date): Health check PASSED" >> "$LOG_FILE"
    exit 0
else
    echo "$(date): Health check FAILED" >> "$LOG_FILE"
    # Send alert (implement your notification method)
    exit 1
fi
EOF

sudo chmod +x /usr/local/bin/webenable-health

# Add to crontab for monitoring every 5 minutes
echo "*/5 * * * * /usr/local/bin/webenable-health" | sudo crontab -
```

---

## 🔧 **Step 7: System Service Integration**

### **7.1 Create Systemd Service**
```bash
sudo tee /etc/systemd/system/webenable-cms.service << EOF
[Unit]
Description=WebEnable CMS
After=network.target
Requires=network.target

[Service]
Type=forking
User=webenable
Group=webenable
WorkingDirectory=/home/webenable/webenable-asia
ExecStart=/usr/bin/podman compose up -d
ExecStop=/usr/bin/podman compose down
ExecReload=/usr/bin/podman compose restart
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable webenable-cms
sudo systemctl start webenable-cms
```

---

## 🛡️ **Step 8: Security Hardening**

### **8.1 Configure UFW (Alternative to firewalld)**
```bash
# If using UFW instead of firewalld
sudo ufw --force enable
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
```

### **8.2 Secure SSH**
```bash
sudo tee -a /etc/ssh/sshd_config << EOF
# Security hardening
Protocol 2
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
X11Forwarding no
MaxAuthTries 3
ClientAliveInterval 300
ClientAliveCountMax 2
EOF

sudo systemctl restart sshd
```

### **8.3 Setup Automated Security Updates**
```bash
# Ubuntu
sudo apt install -y unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades

# CentOS/RHEL
sudo dnf install -y dnf-automatic
sudo systemctl enable --now dnf-automatic.timer
```

---

## 📈 **Step 9: Performance Optimization**

### **9.1 Optimize Container Resources**
Update `docker-compose.yml` with production resource limits:

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 512M
        reservations:
          cpus: '1.0'
          memory: 256M
      replicas: 2
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3

  frontend:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
      replicas: 2

  # Add resource limits to other services...
```

### **9.2 Configure System Limits**
```bash
# Increase file descriptor limits
sudo tee -a /etc/security/limits.conf << EOF
webenable soft nofile 65536
webenable hard nofile 65536
webenable soft nproc 32768
webenable hard nproc 32768
EOF

# Configure kernel parameters
sudo tee -a /etc/sysctl.conf << EOF
# Network optimizations
net.core.rmem_max = 134217728
net.core.wmem_max = 134217728
net.ipv4.tcp_rmem = 4096 65536 134217728
net.ipv4.tcp_wmem = 4096 65536 134217728
net.core.netdev_max_backlog = 30000
net.ipv4.tcp_congestion_control = bbr

# Security
net.ipv4.conf.all.accept_redirects = 0
net.ipv4.conf.all.send_redirects = 0
net.ipv4.conf.all.accept_source_route = 0
EOF

sudo sysctl -p
```

---

## 🔄 **Step 10: Backup and Recovery**

### **10.1 Setup Database Backup**
```bash
# Create backup script
sudo tee /usr/local/bin/backup-webenable << 'EOF'
#!/bin/bash

BACKUP_DIR="/var/backups/webenable"
DATE=$(date +%Y%m%d_%H%M%S)
CONTAINER_NAME="webenable-asia-db-1"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Backup CouchDB
podman exec "$CONTAINER_NAME" curl -X GET http://admin:YOUR_DB_PASSWORD@localhost:5984/_all_dbs | \
    jq -r '.[]' | while read db; do
    echo "Backing up database: $db"
    podman exec "$CONTAINER_NAME" curl -X GET \
        "http://admin:YOUR_DB_PASSWORD@localhost:5984/$db/_all_docs?include_docs=true" \
        > "$BACKUP_DIR/${db}_${DATE}.json"
done

# Compress backups older than 1 day
find "$BACKUP_DIR" -name "*.json" -mtime +1 -exec gzip {} \;

# Remove backups older than 30 days
find "$BACKUP_DIR" -name "*.gz" -mtime +30 -delete

echo "Backup completed: $DATE"
EOF

sudo chmod +x /usr/local/bin/backup-webenable

# Schedule daily backups at 2 AM
echo "0 2 * * * /usr/local/bin/backup-webenable" | sudo crontab -
```

### **10.2 Setup Application Backup**
```bash
# Create application backup script
sudo tee /usr/local/bin/backup-app << 'EOF'
#!/bin/bash

BACKUP_DIR="/var/backups/webenable-app"
DATE=$(date +%Y%m%d_%H%M%S)
APP_DIR="/home/webenable/webenable-asia"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Backup application files (excluding node_modules, logs, etc.)
tar -czf "$BACKUP_DIR/webenable-app_${DATE}.tar.gz" \
    --exclude="node_modules" \
    --exclude="*.log" \
    --exclude=".git" \
    --exclude="backend/main" \
    -C "$(dirname $APP_DIR)" \
    "$(basename $APP_DIR)"

# Remove backups older than 7 days
find "$BACKUP_DIR" -name "*.tar.gz" -mtime +7 -delete

echo "Application backup completed: $DATE"
EOF

sudo chmod +x /usr/local/bin/backup-app
```

---

## 🧪 **Step 11: Testing and Validation**

### **11.1 Performance Testing**
```bash
# Install Apache Bench for testing
sudo apt install -y apache2-utils  # Ubuntu
sudo dnf install -y httpd-tools    # CentOS/RHEL

# Test homepage performance
ab -n 1000 -c 10 https://yourdomain.com/

# Test API performance
ab -n 1000 -c 10 https://yourdomain.com/api/posts

# Test with SSL Labs
echo "Test SSL configuration at: https://www.ssllabs.com/ssltest/analyze.html?d=yourdomain.com"
```

### **11.2 Security Testing**
```bash
# Install security tools
sudo apt install -y nmap nikto

# Basic security scan
nmap -sS -O yourdomain.com
nikto -h https://yourdomain.com
```

### **11.3 Health Check Validation**
```bash
# Test all endpoints
curl -f https://yourdomain.com/api/health
curl -f https://yourdomain.com/
curl -f https://yourdomain.com/blog
curl -f https://yourdomain.com/admin
```

---

## 📋 **Step 12: Post-Deployment Checklist**

### **✅ Security Checklist:**
- [ ] SSL/TLS certificate configured and working
- [ ] Firewall rules configured
- [ ] SSH hardened (no root login, key-based auth)
- [ ] Fail2ban configured
- [ ] Security headers implemented
- [ ] Database access restricted
- [ ] Strong passwords generated
- [ ] Auto-updates configured

### **✅ Performance Checklist:**
- [ ] Resource limits configured
- [ ] Container health checks enabled
- [ ] Compression enabled
- [ ] Caching configured
- [ ] CDN setup (if applicable)
- [ ] Database optimization applied
- [ ] Monitoring implemented

### **✅ Backup Checklist:**
- [ ] Database backup script configured
- [ ] Application backup script configured
- [ ] Backup retention policy implemented
- [ ] Recovery procedure documented
- [ ] Backup restoration tested

### **✅ Monitoring Checklist:**
- [ ] Health checks configured
- [ ] Log rotation configured
- [ ] Performance monitoring enabled
- [ ] Alert notifications setup
- [ ] Service monitoring enabled

---

## 🆘 **Troubleshooting Guide**

### **Common Issues and Solutions:**

#### **Service Won't Start**
```bash
# Check service status
systemctl status webenable-cms

# Check container logs
podman compose logs

# Check system resources
free -h
df -h
```

#### **Database Connection Issues**
```bash
# Test database connectivity
podman exec webenable-asia-db-1 curl http://localhost:5984

# Check database logs
podman logs webenable-asia-db-1
```

#### **SSL Certificate Issues**
```bash
# Check certificate status
sudo certbot certificates

# Renew certificate manually
sudo certbot renew --force-renewal
```

#### **Performance Issues**
```bash
# Check container resources
podman stats

# Check system load
top
htop
iotop
```

---

## 📞 **Support and Maintenance**

### **Regular Maintenance Tasks:**
- **Daily**: Check health status and logs
- **Weekly**: Review security logs and updates
- **Monthly**: Performance review and optimization
- **Quarterly**: Security audit and backup testing

### **Update Procedure:**
```bash
# 1. Create backup
/usr/local/bin/backup-webenable
/usr/local/bin/backup-app

# 2. Pull latest changes
cd /home/webenable/webenable-asia
git pull origin main

# 3. Rebuild and restart
podman compose build
podman compose up -d

# 4. Verify deployment
curl -f https://yourdomain.com/api/health
```

---

## 📚 **Additional Resources**

- **Podman Documentation**: https://docs.podman.io/
- **Caddy Documentation**: https://caddyserver.com/docs/
- **CouchDB Documentation**: https://docs.couchdb.org/
- **Security Best Practices**: https://owasp.org/www-project-top-ten/

---

*This production setup guide ensures a secure, performant, and maintainable WebEnable CMS deployment. Follow each step carefully and customize as needed for your specific environment.*

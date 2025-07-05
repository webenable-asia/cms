# WebEnable CMS - Security Fix Checklist

## 🚨 Critical Security Fixes (Complete ASAP)

### 1. Environment Variables Setup

#### Backend (.env file):
```bash
# Create backend/.env.production
JWT_SECRET=<generate-strong-random-secret>
ADMIN_USERNAME=<change-from-admin>
ADMIN_PASSWORD_HASH=<bcrypt-hash>
COUCHDB_URL=http://<user>:<password>@db:5984/
SMTP_HOST=<your-smtp-host>
SMTP_PORT=587
SMTP_USER=<smtp-username>
SMTP_PASS=<smtp-password>
ALLOWED_ORIGINS=https://your-domain.com
```

#### Frontend (.env.local):
```bash
NEXT_PUBLIC_API_URL=https://api.your-domain.com
```

### 2. Fix Authentication System

- [ ] Install bcrypt: `go get golang.org/x/crypto/bcrypt`
- [ ] Create user initialization script
- [ ] Update auth.go to use database authentication
- [ ] Implement password hashing
- [ ] Add user management endpoints
- [ ] Implement refresh token mechanism

### 3. Add Input Validation

- [ ] Install validator: `go get github.com/go-playground/validator/v10`
- [ ] Add validation middleware
- [ ] Sanitize HTML content to prevent XSS
- [ ] Validate all user inputs
- [ ] Add request size limits

### 4. Secure Database Connection

- [ ] Use TLS for CouchDB connections
- [ ] Implement connection pooling
- [ ] Add database query timeouts
- [ ] Create separate read/write database users

### 5. API Security

- [ ] Implement rate limiting (use `github.com/ulule/limiter`)
- [ ] Add API key authentication for public endpoints
- [ ] Implement request logging
- [ ] Add CSRF protection
- [ ] Set security headers (HSTS, CSP, etc.)

### 6. Security Headers Implementation

Add these headers to all responses:
```go
// backend/middleware/security.go
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        next.ServeHTTP(w, r)
    })
}
```

## 📝 Implementation Examples

### Secure JWT Implementation:
```go
// backend/config/config.go
package config

import (
    "log"
    "os"
    "strings"
)

type Config struct {
    JWTSecret       []byte
    AdminUsername   string
    AdminPassHash   string
    DatabaseURL     string
    AllowedOrigins  []string
}

func Load() *Config {
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET is required")
    }
    
    return &Config{
        JWTSecret:     []byte(jwtSecret),
        AdminUsername: os.Getenv("ADMIN_USERNAME"),
        AdminPassHash: os.Getenv("ADMIN_PASSWORD_HASH"),
        DatabaseURL:   os.Getenv("COUCHDB_URL"),
        AllowedOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
    }
}
```

### Password Hashing:
```go
// backend/utils/auth.go
package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### Input Validation Middleware:
```go
// backend/middleware/validation.go
package middleware

import (
    "encoding/json"
    "net/http"
    "github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate request body size
        r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB
        
        next.ServeHTTP(w, r)
    })
}
```

## 🔒 Production Deployment Checklist

- [ ] Use HTTPS everywhere (Let's Encrypt)
- [ ] Enable CORS only for your domain
- [ ] Set up firewall rules
- [ ] Use secrets management (AWS Secrets Manager, Vault)
- [ ] Enable audit logging
- [ ] Set up intrusion detection
- [ ] Regular security updates
- [ ] Implement backup strategy
- [ ] Set up monitoring and alerts
- [ ] Configure secure headers (Helmet.js for Node/Express)
- [ ] Disable debug mode in production
- [ ] Remove or secure development endpoints
- [ ] Set up DDoS protection (Cloudflare, AWS Shield)
- [ ] Implement session timeout
- [ ] Add IP-based rate limiting for login attempts

## 📊 Security Testing

- [ ] Run OWASP ZAP security scan
- [ ] Test for SQL injection (though using NoSQL)
- [ ] Test for XSS vulnerabilities
- [ ] Test authentication bypass attempts
- [ ] Load test with rate limiting
- [ ] Penetration testing

## 🚀 Quick Start Security Script

Create `scripts/secure-setup.sh`:
```bash
#!/bin/bash

# Generate secure secrets
echo "Generating secure secrets..."
JWT_SECRET=$(openssl rand -base64 32)
ADMIN_PASS=$(openssl rand -base64 16)

# Note: htpasswd might not be available on all systems
# Alternative: Use Go script to generate bcrypt hash
echo "Generated admin password: $ADMIN_PASS"

# Create .env files
cat > backend/.env.production <<EOF
JWT_SECRET=$JWT_SECRET
ADMIN_USERNAME=cms_admin
# Generate ADMIN_PASSWORD_HASH using Go bcrypt
COUCHDB_URL=http://admin:password@db:5984/
ALLOWED_ORIGINS=https://your-domain.com
EOF

echo "Admin password: $ADMIN_PASS"
echo "Save this password securely!"
echo "Use the init_admin.go script to create the admin user with proper bcrypt hashing"
```

## 📅 Timeline

1. **Day 1-2**: Environment variables and configuration
2. **Day 3-4**: Authentication system overhaul
3. **Day 5-6**: Input validation and sanitization
4. **Day 7**: Security headers and rate limiting
5. **Week 2**: Testing and documentation

## ⚠️ Common Security Mistakes to Avoid

1. **Never commit `.env` files** to version control
2. **Don't log sensitive data** (passwords, tokens, personal info)
3. **Avoid using `fmt.Sprintf` for SQL/NoSQL queries** - use parameterized queries
4. **Don't trust client-side validation** - always validate on server
5. **Never store passwords in plain text** - always use bcrypt or similar
6. **Don't expose stack traces** in production error messages
7. **Avoid using outdated dependencies** - regularly update packages
8. **Don't use predictable IDs** - use UUIDs instead of sequential IDs
9. **Never disable TLS certificate verification** in production
10. **Don't expose internal service ports** to the internet

Remember: Security is not a one-time task but an ongoing process!
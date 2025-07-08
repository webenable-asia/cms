# Caddy Reverse Proxy Architecture

## Overview

This CMS uses Caddy as a comprehensive reverse proxy for all HTTP-based services, providing centralized access control, security headers, and performance optimization.

## Architecture

```
[Client] → [Caddy:80/443/5984] → [Services]
                    ↓
    ┌─── Frontend (Next.js:3000)
    ├─── Backend API (Go:8080) 
    └─── Database (CouchDB:5984)
```

## Service Mapping

| Service | Access URL | Upstream | Purpose |
|---------|------------|----------|---------|
| Frontend | `http://localhost` | `frontend:3000` | Main website |
| API | `http://localhost/api/*` | `backend:8080` | REST API endpoints |
| Database | `http://localhost:5984` | `couchdb:5984` | CouchDB direct access |

## Security Features

### 1. Security Headers
- `X-Frame-Options: DENY` - Prevents clickjacking
- `X-Content-Type-Options: nosniff` - Prevents MIME sniffing
- `Referrer-Policy: strict-origin-when-cross-origin` - Controls referrer info

### 2. Database Access Control
- Admin interface restricted to localhost (127.0.0.1)
- Public database access allowed for application usage
- Authentication required for write operations

### 3. Performance Optimization
- Gzip/Zstd compression enabled
- Connection pooling for database connections
- Request/response logging

## Configuration Files

### Caddy Configuration
- **File**: `caddy/Caddyfile`
- **Purpose**: Main reverse proxy configuration
- **Key Features**: Multi-service routing, security headers, compression

### Docker Compose
- **File**: `docker-compose.yml`
- **Ports Exposed**: 80 (HTTP), 443 (HTTPS), 5984 (Database)
- **Dependencies**: Frontend, Backend, Database services

## Testing

Run the comprehensive test suite:
```bash
./scripts/test-caddy-proxy.sh
```

The test script validates:
- ✅ Frontend accessibility
- ✅ API endpoint functionality
- ✅ Database proxy access
- ✅ Security headers presence
- ✅ Compression enablement

## Cache Access

**Note**: Valkey/Redis cache cannot be proxied through HTTP reverse proxy due to its binary protocol. Services connect directly to the cache at `cache:6379`.

## Monitoring

Caddy logs all requests with:
- Request method and path
- Response status codes
- Response times
- Client IP addresses

View logs:
```bash
docker logs cms-caddy-1
```

## Troubleshooting

### Common Issues

1. **Database Connection Refused**
   - Check if CouchDB container is running
   - Verify port 5984 is exposed in docker-compose.yml

2. **Security Headers Missing**
   - Restart Caddy container to reload configuration
   - Check Caddyfile syntax with `caddy validate`

3. **API Routes Not Working**
   - Verify backend container is healthy
   - Check API path matching in Caddyfile

### Health Checks

```bash
# Check all services
docker ps

# Test specific endpoints
curl http://localhost                    # Frontend
curl http://localhost/api/posts         # API
curl http://localhost:5984              # Database
```

## Security Considerations

1. **Production Deployment**
   - Enable HTTPS with proper certificates
   - Restrict database admin access to VPN/internal networks
   - Implement rate limiting for public endpoints

2. **Database Security**
   - Use strong authentication credentials
   - Enable database audit logging
   - Regular security updates for CouchDB

3. **API Security**
   - JWT token validation
   - Input sanitization
   - CORS configuration

## Performance Tuning

- **Connection Pooling**: Configured for database connections
- **Compression**: Enabled for all text-based responses
- **Caching**: Consider adding Caddy cache module for static assets
- **Load Balancing**: Can be extended for multiple backend instances

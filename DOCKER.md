# Docker Development Guide

## 🐳 Docker-First Development

WebEnable CMS is designed for **Docker Compose development** to ensure consistency across all environments and simplify the development workflow.

## Quick Start

### Start Development Environment

```bash
# Navigate to project root
cd /Users/tsaa/workspace/projects/webenable/cms

# Start all services
./dev.sh start

# Or manually
docker-compose up --build
```

### Access Applications

- **Frontend**: http://localhost:3000 (Next.js 15.3.5)
- **Backend API**: http://localhost:8080 (Go 1.24)
- **CouchDB Admin**: http://localhost:5984/_utils
- **Valkey Cache**: localhost:6379

### Default Credentials

- **CouchDB**: admin / password
- **Valkey**: password: `valkeypassword`

## Development Workflow

### 1. Daily Development

```bash
./dev.sh start     # Start all services
./dev.sh logs      # Monitor logs
./dev.sh stop      # Stop when done
```

### 2. Code Changes

- **Frontend**: Hot reload enabled, changes reflect immediately
- **Backend**: Air live reload, Go code restarts automatically
- **Database**: Persistent data in Docker volumes

### 3. Rebuilding

```bash
./dev.sh build frontend  # Rebuild specific service
./dev.sh build           # Rebuild all services
```

### 4. Debugging

```bash
./dev.sh logs frontend   # View frontend logs
./dev.sh logs backend    # View backend logs
./dev.sh shell frontend  # Open shell in container
```

## Docker Services Architecture

### Frontend Service (`webenable-cms-frontend`)

- **Base Image**: `node:22-alpine`
- **Port**: 3000
- **Features**: Hot reload, volume mounts, TypeScript
- **Dependencies**: Backend service

### Backend Service (`webenable-cms-backend`)

- **Base Image**: `golang:1.24-alpine`
- **Port**: 8080
- **Features**: Air live reload, Go modules
- **Dependencies**: Database, Cache

### Database Service (`webenable-cms-db`)

- **Base Image**: `couchdb:3`
- **Port**: 5984
- **Features**: Persistent volumes, admin interface
- **Data**: Stored in `webenable-cms_couchdb_data` volume

### Cache Service (`webenable-cms-cache`)

- **Base Image**: `valkey/valkey:alpine3.22`
- **Port**: 6379
- **Features**: Redis-compatible, health checks
- **Data**: Stored in `webenable-cms_valkey_data` volume

## Development Helper Commands

### Service Management

```bash
./dev.sh start      # Start all services
./dev.sh stop       # Stop all services
./dev.sh restart    # Restart all services
./dev.sh status     # Show service status
```

### Monitoring & Debugging

```bash
./dev.sh logs                # All service logs
./dev.sh logs frontend       # Frontend logs only
./dev.sh logs backend        # Backend logs only
./dev.sh shell frontend      # Shell access to frontend
./dev.sh shell backend       # Shell access to backend
```

### Building & Maintenance

```bash
./dev.sh build              # Build all services
./dev.sh build frontend     # Build frontend only
./dev.sh clean              # Remove containers & volumes
```

### Quick Access

```bash
./dev.sh frontend           # Open frontend in browser
./dev.sh backend            # Open backend API in browser
./dev.sh db                 # Open CouchDB admin in browser
```

## Docker Volumes

### Data Persistence

- **couchdb_data**: Database files and configurations
- **valkey_data**: Cache data and snapshots

### Volume Management

```bash
# View volumes
docker volume ls

# Inspect volume
docker volume inspect webenable-cms_couchdb_data

# Backup volume
docker run --rm -v webenable-cms_couchdb_data:/data -v $(pwd):/backup alpine tar czf /backup/couchdb-backup.tar.gz /data
```

## Network Configuration

### Internal Communication

- Services communicate via Docker network (`webenable-cms_default`)
- Internal hostnames: `frontend`, `backend`, `db`, `cache`
- No external dependencies required

### Port Mapping

| Service  | Internal | External | Purpose |
|----------|----------|----------|---------|
| Frontend | 3000     | 3000     | Web interface |
| Backend  | 8080     | 8080     | API endpoints |
| Database | 5984     | 5984     | Admin interface |
| Cache    | 6379     | 6379     | Direct access |

## Environment Variables

### Configuration Files

- **Docker Compose**: `docker-compose.yml`
- **Environment Template**: `.env.example`
- **Frontend Dockerfile**: `frontend/Dockerfile`
- **Backend Dockerfile**: `backend/Dockerfile`

### Key Variables

```bash
# Database
COUCHDB_USER=admin
COUCHDB_PASSWORD=password

# Cache
VALKEY_PASSWORD=valkeypassword

# Security
JWT_SECRET=your-super-secret-jwt-key

# CORS
CORS_ORIGINS=http://localhost:3000,http://frontend:3000
```

## Performance Optimization

### Build Optimization

- **Layer Caching**: Dockerfile layers optimized for caching
- **Multi-stage Builds**: Separate build and runtime stages
- **Ignore Files**: `.dockerignore` excludes unnecessary files

### Development Speed

- **Volume Mounts**: Code changes without rebuilds
- **Health Checks**: Automated service readiness
- **Hot Reload**: Instant feedback on changes

## Troubleshooting

### Common Issues

1. **Port Conflicts**
   ```bash
   # Check running processes
   lsof -i :3000
   lsof -i :8080
   ```

2. **Permission Issues**
   ```bash
   # Fix volume permissions
   sudo chown -R $USER:$USER .
   ```

3. **Build Failures**
   ```bash
   # Clean build
   ./dev.sh clean
   docker system prune -f
   ./dev.sh start
   ```

### Debug Mode

```bash
# Verbose logging
docker-compose --verbose up

# Individual service debugging
docker-compose up frontend --build
```

## Production Deployment

### Build for Production

```bash
# Create production images
docker-compose -f docker-compose.prod.yml build

# Deploy to registry
docker-compose -f docker-compose.prod.yml push
```

### Environment-Specific Configs

- **Development**: `docker-compose.yml`
- **Production**: `docker-compose.prod.yml`
- **Testing**: `docker-compose.test.yml`

---

Remember: **Always use Docker Compose** for WebEnable CMS development to ensure consistency and avoid environment-specific issues.

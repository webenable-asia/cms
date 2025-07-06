# Rate Limit Reset Documentation

This document explains how to reset rate limits in the WebEnable CMS using Valkey cache.

## Overview

The WebEnable CMS implements rate limiting to protect against abuse and ensure fair usage. Rate limits are stored in Valkey cache with automatic expiration. Sometimes you may need to reset these limits for legitimate users or during testing.

## Rate Limit Types

### 1. API Rate Limits
- **Pattern**: `rate_limit:api:{ip_address}`
- **Default Limit**: 100 requests per minute per IP
- **Applied To**: General API endpoints

### 2. Authentication Rate Limits  
- **Pattern**: `rate_limit:auth:{ip_address}`
- **Default Limit**: 10 attempts per hour per IP
- **Applied To**: Login and authentication endpoints

### 3. User Rate Limits
- **Pattern**: `rate_limit:user:{user_id}`
- **Default Limit**: 120 requests per minute per authenticated user
- **Applied To**: Protected endpoints for authenticated users

## Reset Methods

### Method 1: API Endpoints (Recommended)

Use the built-in admin API endpoints to reset rate limits.

#### Prerequisites
- Admin user account
- JWT authentication token

#### Endpoints

**Reset Rate Limit**
```bash
POST /api/admin/rate-limit/reset
Authorization: Bearer <jwt_token>

Query Parameters:
- type: ip|user|api|auth|users|all
- target: IP address or user ID (required for ip/user types)
```

**Get Rate Limit Status**
```bash
GET /api/admin/rate-limit/status
Authorization: Bearer <jwt_token>

Query Parameters:
- type: ip|user (required)
- target: IP address or user ID (required)
```

#### Examples

Reset rate limit for specific IP:
```bash
curl -X POST "http://localhost:8080/api/admin/rate-limit/reset?type=ip&target=192.168.1.100" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Reset rate limit for specific user:
```bash
curl -X POST "http://localhost:8080/api/admin/rate-limit/reset?type=user&target=user123" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Reset all API rate limits:
```bash
curl -X POST "http://localhost:8080/api/admin/rate-limit/reset?type=api" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Reset all rate limits:
```bash
curl -X POST "http://localhost:8080/api/admin/rate-limit/reset?type=all" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Check rate limit status for IP:
```bash
curl "http://localhost:8080/api/admin/rate-limit/status?type=ip&target=192.168.1.100" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Method 2: Shell Scripts

Use the provided shell scripts for easy command-line access.

#### Script 1: API-based Reset Script

**File**: `scripts/reset-rate-limit.sh`

```bash
# Reset rate limit for IP
./scripts/reset-rate-limit.sh reset-ip 192.168.1.100

# Reset rate limit for user  
./scripts/reset-rate-limit.sh reset-user user123

# Reset all API rate limits
./scripts/reset-rate-limit.sh reset-api

# Reset all authentication rate limits
./scripts/reset-rate-limit.sh reset-auth

# Reset all user rate limits
./scripts/reset-rate-limit.sh reset-users

# Reset ALL rate limits
./scripts/reset-rate-limit.sh reset-all

# Check rate limit status
./scripts/reset-rate-limit.sh status ip 192.168.1.100
./scripts/reset-rate-limit.sh status user user123
```

**Authentication**: Set environment variables:
```bash
export CMS_USERNAME=admin
export CMS_PASSWORD=your_password
```

#### Script 2: Direct Docker Access Script

**File**: `scripts/docker-rate-limit.sh`

This script connects directly to the Valkey container, bypassing API authentication.

```bash
# Reset rate limit for IP
./scripts/docker-rate-limit.sh ip 192.168.1.100

# Reset rate limit for user
./scripts/docker-rate-limit.sh user user123

# Reset all API rate limits
./scripts/docker-rate-limit.sh api

# Reset all authentication rate limits  
./scripts/docker-rate-limit.sh auth

# Reset all user rate limits
./scripts/docker-rate-limit.sh users

# Reset ALL rate limits
./scripts/docker-rate-limit.sh all

# Show current status
./scripts/docker-rate-limit.sh status
```

**Requirements**: Docker and docker-compose

### Method 3: Direct Valkey Commands

For manual debugging or custom scripts, you can connect directly to Valkey:

```bash
# Connect to Valkey via Docker
docker-compose exec cache redis-cli -a valkeypassword

# List all rate limit keys
KEYS rate_limit:*

# List API rate limits
KEYS rate_limit:api:*

# List auth rate limits  
KEYS rate_limit:auth:*

# List user rate limits
KEYS rate_limit:user:*

# Get current value and TTL
GET rate_limit:api:192.168.1.100
TTL rate_limit:api:192.168.1.100

# Delete specific key
DEL rate_limit:api:192.168.1.100

# Delete multiple keys (be careful!)
DEL rate_limit:api:192.168.1.100 rate_limit:auth:192.168.1.100

# Delete all rate limit keys (nuclear option)
EVAL "return redis.call('del', unpack(redis.call('keys', 'rate_limit:*')))" 0
```

## Use Cases

### 1. Legitimate User Locked Out
When a legitimate user is rate limited due to network issues or automated tools:

```bash
# Reset all rate limits for the user's IP
./scripts/docker-rate-limit.sh ip 192.168.1.100
```

### 2. Testing and Development
During testing when you need to reset limits frequently:

```bash
# Reset all rate limits
./scripts/docker-rate-limit.sh all
```

### 3. Authentication Issues
When users can't log in due to auth rate limits:

```bash
# Reset authentication rate limits
./scripts/docker-rate-limit.sh auth
```

### 4. User-Specific Issues
When a specific authenticated user is having problems:

```bash
# Reset rate limits for specific user
./scripts/docker-rate-limit.sh user user123
```

## Monitoring and Alerting

### Check Current Status
```bash
./scripts/docker-rate-limit.sh status
```

### Monitor Rate Limit Usage
```bash
# View all current rate limits
docker-compose exec cache redis-cli -a valkeypassword KEYS "rate_limit:*"

# Count active rate limits
docker-compose exec cache redis-cli -a valkeypassword KEYS "rate_limit:*" | wc -l

# View rate limits by type
docker-compose exec cache redis-cli -a valkeypassword KEYS "rate_limit:api:*"
docker-compose exec cache redis-cli -a valkeypassword KEYS "rate_limit:auth:*"  
docker-compose exec cache redis-cli -a valkeypassword KEYS "rate_limit:user:*"
```

### Rate Limit Headers
The API includes rate limit information in response headers:
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Requests remaining in current window
- `Retry-After`: Seconds to wait before retrying (when rate limited)

## Security Considerations

1. **Admin Access**: Rate limit reset endpoints require admin authentication
2. **Logging**: Rate limit resets should be logged for security auditing
3. **Direct Access**: Docker scripts bypass authentication - use carefully
4. **Rate Limit Patterns**: Don't accidentally reset other cache keys

## Troubleshooting

### Script Permission Issues
```bash
chmod +x scripts/*.sh
```

### Docker Connection Issues
```bash
# Check if services are running
docker-compose ps

# Check Valkey service logs
docker-compose logs cache
```

### Authentication Issues with API
```bash
# Test login first
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your_password"}'
```

### Valkey Connection Issues
```bash
# Test Valkey directly
docker-compose exec cache redis-cli -a valkeypassword PING
```

## Best Practices

1. **Use API endpoints** for production environments
2. **Use Docker scripts** for development and emergency situations
3. **Monitor rate limit usage** to identify patterns
4. **Reset specific targets** rather than all limits when possible
5. **Document rate limit resets** for security auditing
6. **Test scripts** in development before using in production

## Integration Examples

### Automated Monitoring Script
```bash
#!/bin/bash
# Check if rate limits are approaching thresholds
THRESHOLD=50
CURRENT=$(./scripts/docker-rate-limit.sh status | grep "Total Rate Limits" | awk '{print $4}')

if [ "$CURRENT" -gt "$THRESHOLD" ]; then
    echo "High rate limit usage detected: $CURRENT active limits"
    # Send alert or take action
fi
```

### Emergency Reset Script
```bash
#!/bin/bash
# Emergency reset for when users report access issues
echo "Performing emergency rate limit reset..."
./scripts/docker-rate-limit.sh all
echo "All rate limits have been reset"
```

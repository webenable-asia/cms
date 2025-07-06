# Swagger API Documentation

This document explains the Swagger/OpenAPI documentation setup for the WebEnable CMS API.

## Overview

The API documentation is automatically generated using `swaggo/swag` from code annotations and is served via Swagger UI.

## Access Points

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON Specification**: http://localhost:8080/swagger/doc.json
- **YAML Specification**: Static file at `docs/swagger.yaml`

## API Information

- **Title**: WebEnable CMS API
- **Version**: 1.0
- **Description**: A Content Management System API with JWT authentication
- **License**: MIT
- **Base URL**: http://localhost:8080/api

## Authentication

The API uses JWT Bearer token authentication:

```
Authorization: Bearer <your-jwt-token>
```

### Getting a Token

1. **POST** `/api/auth/login` with username/password
2. Use the returned token in the `Authorization` header
3. Token expires after 24 hours

## API Endpoints

### Authentication
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout  
- `GET /auth/me` - Get current user info

### Posts (Public)
- `GET /posts` - Get all published posts
- `GET /posts/{id}` - Get single post

### Posts (Protected)
- `POST /posts` - Create new post
- `PUT /posts/{id}` - Update post
- `DELETE /posts/{id}` - Delete post

### Users (Admin Only)
- `GET /users` - Get all users
- `POST /users` - Create new user
- `GET /users/{id}` - Get user by ID
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user
- `GET /users/stats` - Get user statistics

### Contact
- `POST /contact` - Submit contact form (public)

### System
- `GET /health` - API health check

## Response Format

### Success Responses
All successful responses return JSON with appropriate HTTP status codes:
- `200` - OK
- `201` - Created

### Error Responses
Error responses follow this format:
```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

Common error codes:
- `400` - Bad Request
- `401` - Unauthorized  
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Models

### User
```json
{
  "id": "string",
  "username": "string",
  "email": "string", 
  "role": "admin|editor|author",
  "active": true,
  "created_at": "2025-07-06T08:00:00Z",
  "updated_at": "2025-07-06T08:00:00Z"
}
```

### Post
```json
{
  "id": "string",
  "title": "string",
  "content": "string",
  "excerpt": "string",
  "author": "string",
  "status": "draft|published|scheduled",
  "tags": ["string"],
  "categories": ["string"],
  "featured_image": "string",
  "reading_time": 0,
  "is_featured": false,
  "view_count": 0,
  "created_at": "2025-07-06T08:00:00Z",
  "updated_at": "2025-07-06T08:00:00Z"
}
```

### Contact
```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "company": "string",
  "phone": "string", 
  "subject": "string",
  "message": "string",
  "status": "new|read|replied",
  "created_at": "2025-07-06T08:00:00Z"
}
```

## Regenerating Documentation

When you make changes to API endpoints or annotations:

1. **Using the script**:
   ```bash
   cd backend
   ./generate-swagger.sh
   ```

2. **Manual generation**:
   ```bash
   cd backend
   /Users/tsaa/go/bin/swag init
   ```

3. **Restart the backend**:
   ```bash
   docker-compose restart backend
   ```

## Development Notes

- Swagger annotations are placed directly above handler functions
- The main API info is defined in `main.go`
- Generated files are in the `docs/` directory
- The docs are automatically served at runtime via `http-swagger`

## Tags

API endpoints are organized by these tags:
- **Authentication** - Login, logout, user info
- **Posts** - Blog post management
- **Users** - User management (admin only)
- **Contact** - Contact form handling
- **System** - Health checks and system info

## Security

- JWT tokens are required for most endpoints
- Admin role required for user management
- Public endpoints: posts (read), contact form, health check
- Tokens should be included as: `Authorization: Bearer <token>`

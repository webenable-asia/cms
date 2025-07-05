# WebEnable CMS - Quick Security Fixes Implementation

## Step 1: Create Secure Configuration

### 1.1 Create `backend/config/config.go`:
```go
package config

import (
    "log"
    "os"
    "strings"
)

type Config struct {
    JWTSecret      []byte
    DatabaseURL    string
    Port           string
    AllowedOrigins []string
    SMTPHost       string
    SMTPPort       string
    SMTPUser       string
    SMTPPass       string
}

var AppConfig *Config

func Init() {
    AppConfig = &Config{
        JWTSecret:      getRequiredEnvBytes("JWT_SECRET"),
        DatabaseURL:    getEnvOrDefault("COUCHDB_URL", "http://admin:password@localhost:5984/"),
        Port:           getEnvOrDefault("PORT", "8080"),
        AllowedOrigins: strings.Split(getEnvOrDefault("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
        SMTPHost:       getEnvOrDefault("SMTP_HOST", "localhost"),
        SMTPPort:       getEnvOrDefault("SMTP_PORT", "1025"),
        SMTPUser:       getEnvOrDefault("SMTP_USER", "hello@webenable.asia"),
        SMTPPass:       os.Getenv("SMTP_PASS"),
    }
}

func getRequiredEnvBytes(key string) []byte {
    value := os.Getenv(key)
    if value == "" {
        log.Fatalf("%s environment variable is required", key)
    }
    return []byte(value)
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### 1.2 Update `backend/main.go`:
```go
// Add at the beginning of main()
config.Init()

// Update CORS configuration
c := cors.New(cors.Options{
    AllowedOrigins:   config.AppConfig.AllowedOrigins,
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"*"},
    AllowCredentials: true,
})
```

## Step 2: Fix Authentication

### 2.1 Create `backend/models/user.go`:
```go
package models

import (
    "golang.org/x/crypto/bcrypt"
    "time"
)

type User struct {
    ID           string    `json:"id,omitempty" db:"_id"`
    Rev          string    `json:"rev,omitempty" db:"_rev"`
    Username     string    `json:"username" validate:"required,min=3,max=20"`
    Email        string    `json:"email" validate:"required,email"`
    PasswordHash string    `json:"password_hash,omitempty"`
    Role         string    `json:"role" validate:"required,oneof=admin editor author"`
    Active       bool      `json:"active"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) SetPassword(password string) error {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return err
    }
    u.PasswordHash = string(bytes)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
    return err == nil
}
```

### 2.2 Create `backend/database/users.go`:
```go
package database

import (
    "context"
    "webenable-cms-backend/models"
)

func GetUserByUsername(username string) (*models.User, error) {
    ctx := context.Background()
    
    // Create a simple view to find users by username
    query := map[string]interface{}{
        "selector": map[string]interface{}{
            "username": username,
        },
        "limit": 1,
    }
    
    rows := Instance.UsersDB.Find(ctx, query)
    defer rows.Close()
    
    if rows.Next() {
        var user models.User
        if err := rows.ScanDoc(&user); err != nil {
            return nil, err
        }
        return &user, nil
    }
    
    return nil, nil
}

func CreateUser(user *models.User) error {
    ctx := context.Background()
    _, err := Instance.UsersDB.Put(ctx, user.ID, user)
    return err
}
```

### 2.3 Update `backend/handlers/auth.go`:
```go
package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    
    "webenable-cms-backend/config"
    "webenable-cms-backend/database"
    "webenable-cms-backend/middleware"
    "webenable-cms-backend/models"
    
    "github.com/golang-jwt/jwt/v5"
)

func Login(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    var loginReq models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Get user from database
    user, err := database.GetUserByUsername(loginReq.Username)
    if err != nil {
        http.Error(w, "Authentication failed", http.StatusUnauthorized)
        return
    }
    
    if user == nil || !user.CheckPassword(loginReq.Password) || !user.Active {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
    
    // Create JWT token
    claims := &middleware.Claims{
        Username: user.Username,
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(config.AppConfig.JWTSecret)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }
    
    response := models.LoginResponse{
        Token: tokenString,
        User: models.User{
            ID:       user.ID,
            Username: user.Username,
            Email:    user.Email,
            Role:     user.Role,
        },
    }
    
    json.NewEncoder(w).Encode(response)
}
```

## Step 3: Add Validation Middleware

### 3.1 Install validator:
```bash
cd backend
go get github.com/go-playground/validator/v10
```

### 3.2 Create `backend/middleware/validation.go`:
```go
package middleware

import (
    "encoding/json"
    "net/http"
    "github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateJSON[T any](next func(http.ResponseWriter, *http.Request, T)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var data T
        
        // Limit request body size
        r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB
        
        if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }
        
        if err := validate.Struct(data); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        next(w, r, data)
    }
}
```

## Step 4: Create User Initialization Script

### 4.1 Create `backend/scripts/init_admin.go`:
```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "webenable-cms-backend/config"
    "webenable-cms-backend/database"
    "webenable-cms-backend/models"
    
    "github.com/google/uuid"
)

func main() {
    config.Init()
    database.Init()
    
    adminPassword := os.Getenv("ADMIN_PASSWORD")
    if adminPassword == "" {
        log.Fatal("ADMIN_PASSWORD environment variable required")
    }
    
    admin := &models.User{
        ID:       uuid.New().String(),
        Username: "admin",
        Email:    "admin@webenable.asia",
        Role:     "admin",
        Active:   true,
    }
    
    if err := admin.SetPassword(adminPassword); err != nil {
        log.Fatal("Failed to hash password:", err)
    }
    
    if err := database.CreateUser(admin); err != nil {
        log.Fatal("Failed to create admin user:", err)
    }
    
    fmt.Println("Admin user created successfully!")
    fmt.Printf("Username: %s\n", admin.Username)
}
```

## Step 5: Environment Setup

### 5.1 Create `backend/.env.example`:
```bash
# Security
JWT_SECRET=your-256-bit-secret-here

# Database
COUCHDB_URL=http://admin:password@localhost:5984/

# Server
PORT=8080
ALLOWED_ORIGINS=http://localhost:3000

# Email
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=noreply@example.com
SMTP_PASS=your-smtp-password
```

### 5.2 Create setup script `scripts/setup-security.sh`:
```bash
#!/bin/bash

echo "🔒 Setting up WebEnable CMS Security..."

# Generate JWT secret
JWT_SECRET=$(openssl rand -base64 32)
echo "✅ Generated JWT Secret"

# Generate admin password
ADMIN_PASS=$(openssl rand -base64 12)
echo "✅ Generated Admin Password: $ADMIN_PASS"

# Create production env file
cat > backend/.env.production <<EOF
JWT_SECRET=$JWT_SECRET
COUCHDB_URL=http://admin:password@db:5984/
PORT=8080
ALLOWED_ORIGINS=https://your-domain.com
SMTP_HOST=your-smtp-host
SMTP_PORT=587
SMTP_USER=your-email
SMTP_PASS=your-smtp-password
EOF

echo "✅ Created .env.production"
echo ""
echo "📝 Next steps:"
echo "1. Update ALLOWED_ORIGINS with your domain"
echo "2. Configure SMTP settings"
echo "3. Run: ADMIN_PASSWORD=$ADMIN_PASS go run scripts/init_admin.go"
echo ""
echo "⚠️  Save the admin password securely!"
```

## Quick Start Commands

```bash
# 1. Make setup script executable
chmod +x scripts/setup-security.sh

# 2. Run security setup
./scripts/setup-security.sh

# 3. Install dependencies
cd backend
go get golang.org/x/crypto/bcrypt
go get github.com/go-playground/validator/v10

# 4. Initialize admin user
ADMIN_PASSWORD=<generated-password> go run scripts/init_admin.go

# 5. Start with new config
docker-compose up --build
```

## Testing the Fixes

```bash
# Test login with real credentials
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"<your-admin-password>"}'

# Test with wrong credentials (should fail)
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrong"}'
```

This implementation provides the foundation for secure authentication. Continue with the other security measures in the checklist!
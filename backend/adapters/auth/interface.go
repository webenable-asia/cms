package auth

import (
	"time"
)

// AuthAdapter defines the interface for authentication operations
type AuthAdapter interface {
	// Token Operations
	GenerateToken(claims AuthClaims) (string, error)
	ValidateToken(token string) (*AuthClaims, error)
	RefreshToken(token string) (string, error)
	RevokeToken(token string) error

	// User Authentication
	AuthenticateUser(credentials AuthCredentials) (*AuthResult, error)

	// Claims Management
	ExtractClaims(token string) (*AuthClaims, error)

	// Configuration
	Configure(config AuthConfig) error

	// Health Check
	Health() error
}

// AuthClaims represents authentication claims
type AuthClaims struct {
	UserID    string                 `json:"user_id"`
	Username  string                 `json:"username"`
	Role      string                 `json:"role"`
	Email     string                 `json:"email"`
	IssuedAt  time.Time              `json:"issued_at"`
	ExpiresAt time.Time              `json:"expires_at"`
	Custom    map[string]interface{} `json:"custom"`
}

// AuthCredentials represents authentication credentials
type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

// AuthResult represents authentication result
type AuthResult struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
	ExpiresAt    time.Time   `json:"expires_at"`
	Claims       *AuthClaims `json:"claims"`
}

// AuthConfig holds configuration for authentication adapters
type AuthConfig struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// AuthType constants for supported authentication types
const (
	AuthTypeJWT    = "jwt"
	AuthTypeOAuth2 = "oauth2"
	AuthTypeSAML   = "saml"
	AuthTypeBasic  = "basic"
)

// Common authentication errors
const (
	ErrInvalidToken      = "invalid_token"
	ErrTokenExpired      = "token_expired"
	ErrInvalidCredentials = "invalid_credentials"
	ErrUserNotFound      = "user_not_found"
	ErrTokenRevoked      = "token_revoked"
)
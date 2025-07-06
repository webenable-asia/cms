package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"webenable-cms-backend/models"
)

// JWTAdapter implements AuthAdapter for JWT authentication
type JWTAdapter struct {
	secret     []byte
	expiration time.Duration
	config     map[string]interface{}
}

// NewJWTAdapter creates a new JWT adapter
func NewJWTAdapter(config map[string]interface{}) (AuthAdapter, error) {
	adapter := &JWTAdapter{
		config: config,
	}

	if err := adapter.Configure(AuthConfig{
		Type:   AuthTypeJWT,
		Config: config,
	}); err != nil {
		return nil, err
	}

	return adapter, nil
}

// Configure configures the JWT adapter
func (j *JWTAdapter) Configure(config AuthConfig) error {
	secret, ok := config.Config["secret"].(string)
	if !ok {
		return fmt.Errorf("jwt secret is required")
	}

	j.secret = []byte(secret)

	// Parse expiration duration
	expirationStr, ok := config.Config["expiration"].(string)
	if !ok {
		expirationStr = "24h" // default
	}

	duration, err := time.ParseDuration(expirationStr)
	if err != nil {
		return fmt.Errorf("invalid expiration duration: %w", err)
	}

	j.expiration = duration
	return nil
}

// GenerateToken generates a JWT token for the given claims
func (j *JWTAdapter) GenerateToken(claims AuthClaims) (string, error) {
	now := time.Now()
	if claims.IssuedAt.IsZero() {
		claims.IssuedAt = now
	}
	if claims.ExpiresAt.IsZero() {
		claims.ExpiresAt = now.Add(j.expiration)
	}

	// Convert to JWT claims
	jwtClaims := &JWTClaims{
		Username: claims.Username,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.UserID,
			IssuedAt:  jwt.NewNumericDate(claims.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(claims.ExpiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTAdapter) ValidateToken(tokenString string) (*AuthClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Convert to AuthClaims
	authClaims := &AuthClaims{
		UserID:    claims.Subject,
		Username:  claims.Username,
		Role:      claims.Role,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	return authClaims, nil
}

// RefreshToken refreshes a JWT token
func (j *JWTAdapter) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid token for refresh: %w", err)
	}

	// Create new claims with extended expiration
	newClaims := AuthClaims{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Role:      claims.Role,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(j.expiration),
	}

	return j.GenerateToken(newClaims)
}

// RevokeToken revokes a JWT token (JWT tokens can't be revoked, so this is a no-op)
func (j *JWTAdapter) RevokeToken(token string) error {
	// JWT tokens are stateless and cannot be revoked
	// In a production system, you might want to maintain a blacklist
	return nil
}

// AuthenticateUser authenticates a user with credentials
func (j *JWTAdapter) AuthenticateUser(credentials AuthCredentials) (*AuthResult, error) {
	// This method should typically call a user service to validate credentials
	// For now, we'll return an error as this requires integration with the user system
	return nil, fmt.Errorf("user authentication should be handled by the calling service")
}

// ExtractClaims extracts claims from a token without full validation
func (j *JWTAdapter) ExtractClaims(tokenString string) (*AuthClaims, error) {
	// Parse without verification for claim extraction
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	// Convert to AuthClaims
	authClaims := &AuthClaims{
		UserID:    claims.Subject,
		Username:  claims.Username,
		Role:      claims.Role,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	return authClaims, nil
}

// Health checks the health of the JWT adapter
func (j *JWTAdapter) Health() error {
	if len(j.secret) == 0 {
		return fmt.Errorf("jwt secret not configured")
	}
	return nil
}

// JWTClaims represents JWT-specific claims structure
type JWTClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthenticateUserWithPassword authenticates a user with username/password
func (j *JWTAdapter) AuthenticateUserWithPassword(username, password string, user *models.User) (*AuthResult, error) {
	// Verify password
	if !user.CheckPassword(password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Create claims
	claims := AuthClaims{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		Email:     user.Email,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(j.expiration),
	}

	// Generate token
	token, err := j.GenerateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResult{
		Token:     token,
		User:      user,
		ExpiresAt: claims.ExpiresAt,
		Claims:    &claims,
	}, nil
}

// CreateMiddlewareClaims creates claims compatible with the existing middleware
func (j *JWTAdapter) CreateMiddlewareClaims(authClaims *AuthClaims) *MiddlewareClaims {
	return &MiddlewareClaims{
		Username: authClaims.Username,
		Role:     authClaims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   authClaims.UserID,
			IssuedAt:  jwt.NewNumericDate(authClaims.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(authClaims.ExpiresAt),
		},
	}
}

// MiddlewareClaims represents the claims structure used by the existing middleware
type MiddlewareClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
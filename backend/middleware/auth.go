package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"webenable-cms-backend/adapters/auth"
	"webenable-cms-backend/config"
	"webenable-cms-backend/container"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Global service container for middleware
var globalServiceContainer *container.Container

// SetServiceContainer sets the global service container for middleware
func SetServiceContainer(container *container.Container) {
	globalServiceContainer = container
}

// AuthMiddleware provides JWT authentication using the auth adapter
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]

		// Use auth adapter if available, otherwise fall back to direct JWT
		if globalServiceContainer != nil {
			authAdapter := globalServiceContainer.GetAuthAdapter()
			claims, err := authAdapter.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Convert auth claims to middleware claims
			middlewareClaims := &Claims{
				Username: claims.Username,
				Role:     claims.Role,
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), "user", middlewareClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// Legacy JWT validation (for backward compatibility)
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return config.AppConfig.JWTSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

// AuthMiddlewareWithAdapter creates auth middleware with a specific auth adapter
func AuthMiddlewareWithAdapter(authAdapter auth.AuthAdapter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := bearerToken[1]
			claims, err := authAdapter.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Convert auth claims to middleware claims
			middlewareClaims := &Claims{
				Username: claims.Username,
				Role:     claims.Role,
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), "user", middlewareClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

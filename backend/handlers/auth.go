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

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user and return JWT token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		models.LoginRequest	true	"User credentials"
//	@Success		200			{object}	models.LoginResponse
//	@Failure		400			{object}	models.ErrorResponse
//	@Failure		401			{object}	models.ErrorResponse
//	@Router			/auth/login [post]
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
			Active:   user.Active,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser godoc
//
//	@Summary		Get current user
//	@Description	Get current authenticated user information
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.LoginResponse
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		403	{object}	models.ErrorResponse
//	@Router			/auth/me [get]
//
// GetCurrentUser returns the current authenticated user from JWT claims
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user claims from JWT middleware
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Verify user still exists and is active
	user, err := database.GetUserByUsername(claims.Username)
	if err != nil || user == nil || !user.Active {
		http.Error(w, "User not found or inactive", http.StatusUnauthorized)
		return
	}

	// Ensure user has admin role
	if user.Role != "admin" {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	response := models.LoginResponse{
		User: models.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Active:   user.Active,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// Logout godoc
//
//	@Summary		User logout
//	@Description	Logout user (client-side token invalidation)
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.SuccessResponse
//	@Router			/auth/logout [post]
//
// Logout for cookieless auth simply returns success
// Token invalidation is handled client-side
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

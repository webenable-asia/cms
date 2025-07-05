package handlers

import (
	"encoding/json"
	"net/http"
	"time"

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

	// In a real application, you would validate credentials against database
	// For demo purposes, we'll use hardcoded credentials
	if loginReq.Username != "admin" || loginReq.Password != "password" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	claims := &middleware.Claims{
		Username: loginReq.Username,
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := models.LoginResponse{
		Token: tokenString,
		User: models.User{
			Username: loginReq.Username,
			Role:     "admin",
		},
	}

	json.NewEncoder(w).Encode(response)
}

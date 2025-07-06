package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"webenable-cms-backend/database"
	"webenable-cms-backend/middleware"
	"webenable-cms-backend/models"

	"github.com/gorilla/mux"
)

// GetUsers godoc
//
//	@Summary		Get all users
//	@Description	Get all users (admin only)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		models.User
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		403	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users [get]
//
// GetUsers returns all users (admin only)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user claims from JWT middleware
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Only admin can access user management
	if claims.Role != "admin" {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	users, err := database.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// GetUser godoc
//
//	@Summary		Get user by ID
//	@Description	Get a single user by ID (admin only)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	models.User
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		403	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users/{id} [get]
//
// GetUser returns a single user by ID (admin only)
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user claims from JWT middleware
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Only admin can access user management
	if claims.Role != "admin" {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := database.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Don't return password hash and revision in API response
	user.PasswordHash = ""
	user.Rev = ""

	json.NewEncoder(w).Encode(user)
}

// CreateUser godoc
//
//	@Summary		Create new user
//	@Description	Create a new user (admin only)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			user	body		object{username=string,email=string,password=string,role=string,active=bool}	true	"User data"
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		403		{object}	models.ErrorResponse
//	@Failure		409		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/users [post]
//
// CreateUser creates a new user (admin only)
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user claims from JWT middleware
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Only admin can create users
	if claims.Role != "admin" {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Active   bool   `json:"active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Validate role
	if req.Role != "admin" && req.Role != "editor" && req.Role != "author" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	existingUser, err := database.GetUserByUsername(req.Username)
	if err == nil && existingUser != nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Check if email already exists
	existingUser, err = database.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// Create new user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		Active:   req.Active,
	}

	// Set password
	if err := user.SetPassword(req.Password); err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Save user
	if err := database.CreateUser(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Don't return password hash and revision in API response
	user.PasswordHash = ""
	user.Rev = ""

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser godoc
//
//	@Summary		Update user
//	@Description	Update an existing user (admin only)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path	string																									true	"User ID"
//	@Param			user	body	object{username=string,email=string,password=string,role=string,active=bool}	true	"User data"
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		403		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/users/{id} [put]
//
// UpdateUser updates an existing user (admin only)
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user claims from JWT middleware
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Only admin can update users
	if claims.Role != "admin" {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Active   *bool  `json:"active"` // Use pointer to distinguish between false and not set
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user exists
	existingUser, err := database.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Validate role if provided
	if req.Role != "" && req.Role != "admin" && req.Role != "editor" && req.Role != "author" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	// Check if username already exists (if changing)
	if req.Username != "" && req.Username != existingUser.Username {
		existing, err := database.GetUserByUsername(req.Username)
		if err == nil && existing != nil {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
	}

	// Check if email already exists (if changing)
	if req.Email != "" && req.Email != existingUser.Email {
		existing, err := database.GetUserByEmail(req.Email)
		if err == nil && existing != nil {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
	}

	// Prepare updates
	updates := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
	}

	if req.Active != nil {
		updates.Active = *req.Active
	} else {
		updates.Active = existingUser.Active
	}

	// Hash password if provided
	if req.Password != "" {
		if err := updates.SetPassword(req.Password); err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
	}

	// Update user
	updatedUser, err := database.UpdateUser(userID, updates)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser deletes a user (admin only)
// DeleteUser godoc
//
//	@Summary		Delete user
//	@Description	Delete a user (admin only)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	models.SuccessResponse
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		403	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user claims from JWT middleware
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Only admin can delete users
	if claims.Role != "admin" {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	// Prevent admin from deleting themselves
	currentUser, err := database.GetUserByUsername(claims.Username)
	if err == nil && currentUser != nil && currentUser.ID == userID {
		http.Error(w, "Cannot delete your own account", http.StatusBadRequest)
		return
	}

	// Check if user exists
	_, err = database.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Delete user
	if err := database.DeleteUser(userID); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

// GetUserStats godoc
//
//	@Summary		Get user statistics
//	@Description	Get user statistics (admin only)
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.UserStatsResponse
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		403	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users/stats [get]
//
// GetUserStats returns user statistics (admin only)
func GetUserStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user claims from JWT middleware
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Only admin can access stats
	if claims.Role != "admin" {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	users, err := database.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	stats := map[string]interface{}{
		"total_users":  len(users),
		"active_users": 0,
		"admin_users":  0,
		"editor_users": 0,
		"author_users": 0,
		"last_updated": time.Now(),
	}

	for _, user := range users {
		if user.Active {
			stats["active_users"] = stats["active_users"].(int) + 1
		}

		switch user.Role {
		case "admin":
			stats["admin_users"] = stats["admin_users"].(int) + 1
		case "editor":
			stats["editor_users"] = stats["editor_users"].(int) + 1
		case "author":
			stats["author_users"] = stats["author_users"].(int) + 1
		}
	}

	json.NewEncoder(w).Encode(stats)
}

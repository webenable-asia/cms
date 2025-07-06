package database

import (
	"context"
	"time"
	"webenable-cms-backend/models"

	"github.com/google/uuid"
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

func GetUserByEmail(email string) (*models.User, error) {
	ctx := context.Background()

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"email": email,
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

func GetUserByID(userID string) (*models.User, error) {
	ctx := context.Background()

	var user models.User
	err := Instance.UsersDB.Get(ctx, userID).ScanDoc(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	ctx := context.Background()

	// Simplified query without sort for compatibility
	query := map[string]interface{}{
		"selector": map[string]interface{}{},
	}

	rows := Instance.UsersDB.Find(ctx, query)
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.ScanDoc(&user); err != nil {
			continue
		}
		// Don't return password hashes and revisions in lists
		user.PasswordHash = ""
		user.Rev = ""
		users = append(users, user)
	}

	return users, nil
}

func CreateUser(user *models.User) error {
	ctx := context.Background()

	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := Instance.UsersDB.Put(ctx, user.ID, user)
	return err
}

func UpdateUser(userID string, updates *models.User) (*models.User, error) {
	ctx := context.Background()

	// Get existing user first
	existing, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if updates.Username != "" {
		existing.Username = updates.Username
	}
	if updates.Email != "" {
		existing.Email = updates.Email
	}
	if updates.Role != "" {
		existing.Role = updates.Role
	}
	if updates.PasswordHash != "" {
		existing.PasswordHash = updates.PasswordHash
	}
	// Active can be explicitly set to false
	existing.Active = updates.Active
	existing.UpdatedAt = time.Now()

	_, err = Instance.UsersDB.Put(ctx, userID, existing)
	if err != nil {
		return nil, err
	}

	// Don't return password hash and revision
	existing.PasswordHash = ""
	existing.Rev = ""
	return existing, nil
}

func DeleteUser(userID string) error {
	ctx := context.Background()

	// Get the user first to get the revision
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	_, err = Instance.UsersDB.Delete(ctx, userID, user.Rev)
	return err
}

func GetUserCount() (int, error) {
	ctx := context.Background()

	query := map[string]interface{}{
		"selector": map[string]interface{}{},
	}

	rows := Instance.UsersDB.Find(ctx, query)
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	return count, nil
}

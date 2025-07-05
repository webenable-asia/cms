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

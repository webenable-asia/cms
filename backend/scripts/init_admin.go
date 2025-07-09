package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin" // fallback to default
	}

	admin := &models.User{
		ID:        uuid.New().String(),
		Username:  adminUsername,
		Email:     "admin@webenable.asia",
		Role:      "admin",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

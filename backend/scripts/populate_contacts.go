package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
)

type Contact struct {
	ID        string     `json:"_id,omitempty"`
	Rev       string     `json:"_rev,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Company   string     `json:"company,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Subject   string     `json:"subject"`
	Message   string     `json:"message"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	RepliedAt *time.Time `json:"replied_at,omitempty"`
}

func main() {
	// Get CouchDB connection details from environment or use defaults
	couchdbURL := os.Getenv("COUCHDB_URL")
	if couchdbURL == "" {
		couchdbURL = "http://admin:password@localhost:5984"
	}

	// Connect to CouchDB
	client, err := kivik.New("couch", couchdbURL)
	if err != nil {
		log.Fatal("Failed to connect to CouchDB:", err)
	}

	ctx := context.Background()

	// Get or create the contacts database
	dbExists, err := client.DBExists(ctx, "contacts")
	if err != nil {
		log.Fatal("Failed to check if database exists:", err)
	}

	var db *kivik.DB
	if !dbExists {
		err = client.CreateDB(ctx, "contacts")
		if err != nil {
			log.Fatal("Failed to create contacts database:", err)
		}
		fmt.Println("Created contacts database")
	}

	db = client.DB("contacts")

	// Sample contact messages
	contacts := []Contact{
		{
			ID:        "contact-001",
			Name:      "John Smith",
			Email:     "john.smith@example.com",
			Company:   "TechCorp Inc.",
			Phone:     "+1-555-0123",
			Subject:   "Website Development Inquiry",
			Message:   "Hi, I'm interested in having a new website built for my company. We're looking for a modern, responsive design with e-commerce functionality. Could you provide a quote and timeline?",
			Status:    "new",
			CreatedAt: time.Now().Add(-2 * 24 * time.Hour), // 2 days ago
		},
		{
			ID:        "contact-002",
			Name:      "Sarah Johnson",
			Email:     "sarah.j@marketing-pro.com",
			Company:   "Marketing Pro",
			Phone:     "+1-555-0456",
			Subject:   "Digital Marketing Partnership",
			Message:   "Hello! We're a marketing agency and would like to explore partnership opportunities. We have clients who need web development services and think we could work together.",
			Status:    "read",
			CreatedAt: time.Now().Add(-5 * 24 * time.Hour), // 5 days ago
		},
		{
			ID:        "contact-003",
			Name:      "Mike Chen",
			Email:     "m.chen@startup.io",
			Company:   "Startup Inc.",
			Subject:   "MVP Development",
			Message:   "We're a early-stage startup looking to build our MVP. We need a full-stack web application with user authentication, payments, and admin dashboard. What's your availability for Q2?",
			Status:    "new",
			CreatedAt: time.Now().Add(-1 * 24 * time.Hour), // 1 day ago
		},
		{
			ID:        "contact-004",
			Name:      "Emily Rodriguez",
			Email:     "emily@nonprofit.org",
			Company:   "Community Helper Nonprofit",
			Phone:     "+1-555-0789",
			Subject:   "Nonprofit Website Redesign",
			Message:   "Our nonprofit needs a website redesign to better showcase our mission and accept donations online. We have a limited budget but are looking for high-quality work. Do you offer any nonprofit discounts?",
			Status:    "replied",
			CreatedAt: time.Now().Add(-7 * 24 * time.Hour), // 1 week ago
		},
		{
			ID:        "contact-005",
			Name:      "Alex Thompson",
			Email:     "alex.thompson@freelancer.com",
			Subject:   "General Question",
			Message:   "What technologies do you specialize in? I'm looking for help with a React/Node.js project and wondered if that's in your wheelhouse.",
			Status:    "new",
			CreatedAt: time.Now().Add(-3 * time.Hour), // 3 hours ago
		},
	}

	// Insert contacts into database
	for _, contact := range contacts {
		// Check if contact already exists
		exists := true
		row := db.Get(ctx, contact.ID)
		var existingContact Contact
		if err := row.ScanDoc(&existingContact); err != nil {
			exists = false
		}

		if exists {
			fmt.Printf("Contact from '%s' already exists, skipping...\n", contact.Name)
			continue
		}

		// Insert contact
		_, err := db.Put(ctx, contact.ID, contact)
		if err != nil {
			log.Printf("Failed to insert contact from '%s': %v", contact.Name, err)
			continue
		}

		fmt.Printf("Successfully inserted contact: %s - %s\n", contact.Name, contact.Subject)
	}

	fmt.Println("Contact database population completed!")
}

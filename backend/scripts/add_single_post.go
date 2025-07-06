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

type Post struct {
	ID            string     `json:"_id,omitempty"`
	Rev           string     `json:"_rev,omitempty"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	Excerpt       string     `json:"excerpt"`
	Author        string     `json:"author"`
	Status        string     `json:"status"`
	Tags          []string   `json:"tags"`
	Categories    []string   `json:"categories"`
	FeaturedImage string     `json:"featured_image"`
	ImageAlt      string     `json:"image_alt"`
	MetaTitle     string     `json:"meta_title"`
	MetaDesc      string     `json:"meta_description"`
	ReadingTime   int        `json:"reading_time"`
	IsFeatured    bool       `json:"is_featured"`
	ViewCount     int        `json:"view_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
	ScheduledAt   *time.Time `json:"scheduled_at,omitempty"`
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
	db := client.DB("posts")

	// Create a new blog post via curl command
	now := time.Now()
	post := Post{
		ID:            "api-testing-with-curl",
		Title:         "Testing Our CMS API with cURL",
		Content:       "# Testing Our CMS API with cURL\n\n## Introduction\n\nToday we successfully tested our CMS API using cURL commands! This demonstrates that our backend API is working correctly and can handle blog post creation.\n\n## What We Accomplished\n\n- ✅ Set up authentication with admin user\n- ✅ Created blog posts via database population script\n- ✅ Verified API endpoints are responding correctly\n- ✅ Confirmed frontend displays posts properly\n\n## API Testing Process\n\nHere's how we tested the API:\n\n1. **Created Admin User** - Initialized the admin account with credentials\n2. **Authentication** - Used /api/auth/login to get session tokens\n3. **Database Population** - Used Go scripts to add sample content\n4. **Verification** - Checked endpoints with curl commands\n\n## Next Steps\n\nWith the CMS now fully functional, we can:\n\n- Create new posts through the admin interface\n- Edit existing content\n- Manage categories and tags\n- Handle contact form submissions\n\n## Conclusion\n\nOur CMS system is now complete and working perfectly! The combination of Go backend, CouchDB database, and Next.js frontend provides a robust content management solution.",
		Excerpt:       "Learn how we successfully tested our CMS API using cURL commands and database scripts.",
		Author:        "API Testing Team",
		Status:        "published",
		Tags:          []string{"api", "testing", "curl", "cms"},
		Categories:    []string{"Development", "Testing"},
		FeaturedImage: "https://images.unsplash.com/photo-1555949963-aa79dcee981c?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2070&q=80",
		ImageAlt:      "Terminal window showing curl commands and API testing",
		MetaTitle:     "Testing Our CMS API with cURL | WebEnable",
		MetaDesc:      "Learn how we successfully tested our CMS API using cURL commands and database scripts.",
		ReadingTime:   4,
		IsFeatured:    false,
		ViewCount:     0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Set published date
	publishedAt := now
	post.PublishedAt = &publishedAt

	// Check if post already exists
	exists := true
	row := db.Get(ctx, post.ID)
	var existingPost Post
	if err := row.ScanDoc(&existingPost); err != nil {
		exists = false
	}

	if exists {
		fmt.Printf("Post '%s' already exists!\n", post.Title)
		return
	}

	// Insert post
	_, err = db.Put(ctx, post.ID, post)
	if err != nil {
		log.Fatal("Failed to insert post:", err)
	}

	fmt.Printf("Successfully created blog post: %s\n", post.Title)
	fmt.Printf("Post ID: %s\n", post.ID)
	fmt.Printf("View at: http://localhost:3000/blog/%s\n", post.ID)
}

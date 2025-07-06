package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
	"github.com/google/uuid"
	"webenable-cms-backend/models"
)

// CouchDBAdapter implements DatabaseAdapter for CouchDB
type CouchDBAdapter struct {
	client     *kivik.Client
	postsDB    *kivik.DB
	usersDB    *kivik.DB
	contactsDB *kivik.DB
	config     map[string]interface{}
}

// NewCouchDBAdapter creates a new CouchDB adapter
func NewCouchDBAdapter(config map[string]interface{}) (DatabaseAdapter, error) {
	adapter := &CouchDBAdapter{
		config: config,
	}

	if err := adapter.Connect(DatabaseConfig{
		Type:   DatabaseTypeCouchDB,
		Config: config,
	}); err != nil {
		return nil, err
	}

	return adapter, nil
}

// Connect establishes connection to CouchDB
func (c *CouchDBAdapter) Connect(config DatabaseConfig) error {
	url, ok := config.Config["url"].(string)
	if !ok {
		return fmt.Errorf("couchdb url is required")
	}

	client, err := kivik.New("couch", url)
	if err != nil {
		return fmt.Errorf("failed to connect to CouchDB: %w", err)
	}

	c.client = client

	// Create databases if they don't exist
	ctx := context.Background()

	// Create posts database
	if exists, _ := client.DBExists(ctx, "posts"); !exists {
		if err := client.CreateDB(ctx, "posts"); err != nil {
			return fmt.Errorf("failed to create posts database: %w", err)
		}
	}

	// Create users database
	if exists, _ := client.DBExists(ctx, "users"); !exists {
		if err := client.CreateDB(ctx, "users"); err != nil {
			return fmt.Errorf("failed to create users database: %w", err)
		}
	}

	// Create contacts database
	if exists, _ := client.DBExists(ctx, "contacts"); !exists {
		if err := client.CreateDB(ctx, "contacts"); err != nil {
			return fmt.Errorf("failed to create contacts database: %w", err)
		}
	}

	c.postsDB = client.DB("posts")
	c.usersDB = client.DB("users")
	c.contactsDB = client.DB("contacts")

	log.Println("CouchDB adapter connected successfully")
	return nil
}

// Close closes the connection
func (c *CouchDBAdapter) Close() error {
	// CouchDB client doesn't have a close method
	return nil
}

// Health checks the health of the connection
func (c *CouchDBAdapter) Health() error {
	ctx := context.Background()
	_, err := c.client.Ping(ctx)
	return err
}

// Post Operations

// CreatePost creates a new post
func (c *CouchDBAdapter) CreatePost(post *models.Post) error {
	ctx := context.Background()

	if post.ID == "" {
		post.ID = uuid.New().String()
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	if post.Status == "published" && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}

	// Create document
	doc := map[string]interface{}{
		"title":          post.Title,
		"content":        post.Content,
		"excerpt":        post.Excerpt,
		"author":         post.Author,
		"status":         post.Status,
		"tags":           post.Tags,
		"categories":     post.Categories,
		"featured_image": post.FeaturedImage,
		"image_alt":      post.ImageAlt,
		"meta_title":     post.MetaTitle,
		"meta_description": post.MetaDesc,
		"reading_time":   post.ReadingTime,
		"is_featured":    post.IsFeatured,
		"view_count":     post.ViewCount,
		"created_at":     post.CreatedAt,
		"updated_at":     post.UpdatedAt,
		"published_at":   post.PublishedAt,
		"scheduled_at":   post.ScheduledAt,
	}

	rev, err := c.postsDB.Put(ctx, post.ID, doc)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	post.Rev = rev
	return nil
}

// GetPost retrieves a post by ID
func (c *CouchDBAdapter) GetPost(id string) (*models.Post, error) {
	ctx := context.Background()

	row := c.postsDB.Get(ctx, id)
	var post models.Post
	if err := row.ScanDoc(&post); err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

// GetPosts retrieves posts with pagination
func (c *CouchDBAdapter) GetPosts(limit, offset int) ([]models.Post, error) {
	ctx := context.Background()

	// Use AllDocs with include_docs for simplicity
	// In production, you might want to use a view for better performance
	rows := c.postsDB.AllDocs(ctx, kivik.Param("include_docs", true))
	defer rows.Close()

	var posts []models.Post
	count := 0
	skipped := 0

	for rows.Next() {
		if skipped < offset {
			skipped++
			continue
		}

		if count >= limit {
			break
		}

		var post models.Post
		if err := rows.ScanDoc(&post); err != nil {
			continue
		}

		posts = append(posts, post)
		count++
	}

	return posts, nil
}

// UpdatePost updates a post
func (c *CouchDBAdapter) UpdatePost(id string, post *models.Post) error {
	ctx := context.Background()

	// Get existing post first
	existing, err := c.GetPost(id)
	if err != nil {
		return fmt.Errorf("failed to get existing post: %w", err)
	}

	// Update fields
	post.ID = id
	post.Rev = existing.Rev
	post.CreatedAt = existing.CreatedAt
	post.UpdatedAt = time.Now()

	if post.Status == "published" && existing.Status != "published" {
		now := time.Now()
		post.PublishedAt = &now
	}

	// Create document
	doc := map[string]interface{}{
		"_rev":           post.Rev,
		"title":          post.Title,
		"content":        post.Content,
		"excerpt":        post.Excerpt,
		"author":         post.Author,
		"status":         post.Status,
		"tags":           post.Tags,
		"categories":     post.Categories,
		"featured_image": post.FeaturedImage,
		"image_alt":      post.ImageAlt,
		"meta_title":     post.MetaTitle,
		"meta_description": post.MetaDesc,
		"reading_time":   post.ReadingTime,
		"is_featured":    post.IsFeatured,
		"view_count":     post.ViewCount,
		"created_at":     post.CreatedAt,
		"updated_at":     post.UpdatedAt,
		"published_at":   post.PublishedAt,
		"scheduled_at":   post.ScheduledAt,
	}

	rev, err := c.postsDB.Put(ctx, id, doc)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	post.Rev = rev
	return nil
}

// DeletePost deletes a post
func (c *CouchDBAdapter) DeletePost(id string) error {
	ctx := context.Background()

	// Get existing post to get revision
	existing, err := c.GetPost(id)
	if err != nil {
		return fmt.Errorf("failed to get existing post: %w", err)
	}

	_, err = c.postsDB.Delete(ctx, id, existing.Rev)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// User Operations

// CreateUser creates a new user
func (c *CouchDBAdapter) CreateUser(user *models.User) error {
	ctx := context.Background()

	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := c.usersDB.Put(ctx, user.ID, user)
	return err
}

// GetUser retrieves a user by ID
func (c *CouchDBAdapter) GetUser(id string) (*models.User, error) {
	ctx := context.Background()

	var user models.User
	err := c.usersDB.Get(ctx, id).ScanDoc(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (c *CouchDBAdapter) GetUserByUsername(username string) (*models.User, error) {
	ctx := context.Background()

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"username": username,
		},
		"limit": 1,
	}

	rows := c.usersDB.Find(ctx, query)
	defer rows.Close()

	if rows.Next() {
		var user models.User
		if err := rows.ScanDoc(&user); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		return &user, nil
	}

	return nil, fmt.Errorf("user not found")
}

// GetUserByEmail retrieves a user by email
func (c *CouchDBAdapter) GetUserByEmail(email string) (*models.User, error) {
	ctx := context.Background()

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"email": email,
		},
		"limit": 1,
	}

	rows := c.usersDB.Find(ctx, query)
	defer rows.Close()

	if rows.Next() {
		var user models.User
		if err := rows.ScanDoc(&user); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		return &user, nil
	}

	return nil, fmt.Errorf("user not found")
}

// GetUsers retrieves users with pagination
func (c *CouchDBAdapter) GetUsers(limit, offset int) ([]models.User, error) {
	ctx := context.Background()

	query := map[string]interface{}{
		"selector": map[string]interface{}{},
		"limit":    limit,
		"skip":     offset,
	}

	rows := c.usersDB.Find(ctx, query)
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

// UpdateUser updates a user
func (c *CouchDBAdapter) UpdateUser(id string, user *models.User) error {
	ctx := context.Background()

	// Get existing user first
	existing, err := c.GetUser(id)
	if err != nil {
		return fmt.Errorf("failed to get existing user: %w", err)
	}

	// Update fields
	if user.Username != "" {
		existing.Username = user.Username
	}
	if user.Email != "" {
		existing.Email = user.Email
	}
	if user.Role != "" {
		existing.Role = user.Role
	}
	if user.PasswordHash != "" {
		existing.PasswordHash = user.PasswordHash
	}
	// Active can be explicitly set to false
	existing.Active = user.Active
	existing.UpdatedAt = time.Now()

	_, err = c.usersDB.Put(ctx, id, existing)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Update the passed user with the new values
	*user = *existing
	// Don't return password hash and revision
	user.PasswordHash = ""
	user.Rev = ""

	return nil
}

// DeleteUser deletes a user
func (c *CouchDBAdapter) DeleteUser(id string) error {
	ctx := context.Background()

	// Get the user first to get the revision
	user, err := c.GetUser(id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	_, err = c.usersDB.Delete(ctx, id, user.Rev)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// Contact Operations

// CreateContact creates a new contact
func (c *CouchDBAdapter) CreateContact(contact *models.Contact) error {
	ctx := context.Background()

	if contact.ID == "" {
		contact.ID = uuid.New().String()
	}

	contact.CreatedAt = time.Now()
	contact.Status = "new"

	// Create document
	doc := map[string]interface{}{
		"name":       contact.Name,
		"email":      contact.Email,
		"company":    contact.Company,
		"phone":      contact.Phone,
		"subject":    contact.Subject,
		"message":    contact.Message,
		"status":     contact.Status,
		"created_at": contact.CreatedAt,
		"read_at":    contact.ReadAt,
		"replied_at": contact.RepliedAt,
	}

	rev, err := c.contactsDB.Put(ctx, contact.ID, doc)
	if err != nil {
		return fmt.Errorf("failed to create contact: %w", err)
	}

	contact.Rev = rev
	return nil
}

// GetContact retrieves a contact by ID
func (c *CouchDBAdapter) GetContact(id string) (*models.Contact, error) {
	ctx := context.Background()

	row := c.contactsDB.Get(ctx, id)
	var contact models.Contact
	if err := row.ScanDoc(&contact); err != nil {
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	return &contact, nil
}

// GetContacts retrieves contacts with pagination
func (c *CouchDBAdapter) GetContacts(limit, offset int) ([]models.Contact, error) {
	ctx := context.Background()

	rows := c.contactsDB.AllDocs(ctx, kivik.Param("include_docs", true))
	defer rows.Close()

	var contacts []models.Contact
	count := 0
	skipped := 0

	for rows.Next() {
		if skipped < offset {
			skipped++
			continue
		}

		if count >= limit {
			break
		}

		var contact models.Contact
		if err := rows.ScanDoc(&contact); err != nil {
			continue
		}

		contacts = append(contacts, contact)
		count++
	}

	return contacts, nil
}

// UpdateContact updates a contact
func (c *CouchDBAdapter) UpdateContact(id string, contact *models.Contact) error {
	ctx := context.Background()

	// Get existing contact first
	existing, err := c.GetContact(id)
	if err != nil {
		return fmt.Errorf("failed to get existing contact: %w", err)
	}

	// Update fields
	contact.ID = id
	contact.Rev = existing.Rev
	contact.CreatedAt = existing.CreatedAt

	// Create document
	doc := map[string]interface{}{
		"_rev":       contact.Rev,
		"name":       contact.Name,
		"email":      contact.Email,
		"company":    contact.Company,
		"phone":      contact.Phone,
		"subject":    contact.Subject,
		"message":    contact.Message,
		"status":     contact.Status,
		"created_at": contact.CreatedAt,
		"read_at":    contact.ReadAt,
		"replied_at": contact.RepliedAt,
	}

	_, err = c.contactsDB.Put(ctx, id, doc)
	if err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
	}

	return nil
}

// DeleteContact deletes a contact
func (c *CouchDBAdapter) DeleteContact(id string) error {
	ctx := context.Background()

	// Get existing contact to get revision
	existing, err := c.GetContact(id)
	if err != nil {
		return fmt.Errorf("failed to get existing contact: %w", err)
	}

	_, err = c.contactsDB.Delete(ctx, id, existing.Rev)
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}

	return nil
}

// Transaction Support

// BeginTransaction begins a transaction (CouchDB doesn't support transactions)
func (c *CouchDBAdapter) BeginTransaction() (Transaction, error) {
	return &CouchDBTransaction{}, nil
}

// CouchDBTransaction implements Transaction interface for CouchDB
type CouchDBTransaction struct{}

// Commit commits the transaction (no-op for CouchDB)
func (t *CouchDBTransaction) Commit() error {
	return nil
}

// Rollback rolls back the transaction (no-op for CouchDB)
func (t *CouchDBTransaction) Rollback() error {
	return nil
}

// Context returns the transaction context
func (t *CouchDBTransaction) Context() context.Context {
	return context.Background()
}
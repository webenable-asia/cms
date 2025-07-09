package database

import (
	"context"
	"fmt"
	"time"

	"webenable-cms-backend/models"
	"webenable-cms-backend/utils"

	"github.com/sirupsen/logrus"
)

// OptimizedQueries provides optimized database query methods
type OptimizedQueries struct {
	db *OptimizedDB
}

// NewOptimizedQueries creates a new optimized queries instance
func NewOptimizedQueries(db *OptimizedDB) *OptimizedQueries {
	return &OptimizedQueries{db: db}
}

// Posts queries

// GetPostsPaginated retrieves posts with optimized pagination and filtering
func (q *OptimizedQueries) GetPostsPaginated(ctx context.Context, status string, page, limit int) (*models.PaginatedPostsResponse, error) {
	start := time.Now()
	defer func() {
		q.db.trackQuery("GetPostsPaginated", time.Since(start))
	}()

	conn := q.db.getConnection()
	defer q.db.returnConnection(conn)

	postsDB := conn.DB("posts")

	// Build optimized query using indexes
	query := map[string]interface{}{
		"selector": map[string]interface{}{},
		"sort": []map[string]string{
			{"published_at": "desc"},
			{"created_at": "desc"},
		},
		"limit": limit,
		"skip":  (page - 1) * limit,
	}

	// Add status filter if specified
	if status != "" {
		query["selector"].(map[string]interface{})["status"] = status
		// Use the status-published index for better performance
		query["use_index"] = "status-published-index"
	}

	rows := postsDB.Find(ctx, query)
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.ScanDoc(&post); err != nil {
			utils.LogError(err, "Failed to scan post document", logrus.Fields{
				"query": "GetPostsPaginated",
			})
			continue
		}

		// Set document ID and revision
		if id, err := rows.ID(); err == nil && id != "" {
			post.ID = id
		}
		if rev, err := rows.Rev(); err == nil && rev != "" {
			post.Rev = rev
		}

		posts = append(posts, post)
	}

	// Get total count for pagination metadata
	total, err := q.getPostsCount(ctx, status)
	if err != nil {
		utils.LogError(err, "Failed to get posts count", logrus.Fields{
			"status": status,
		})
		total = len(posts) // Fallback to current page count
	}

	// Calculate pagination metadata
	totalPages := (total + limit - 1) / limit
	meta := models.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return &models.PaginatedPostsResponse{
		Data: posts,
		Meta: meta,
	}, nil
}

// GetPostByID retrieves a single post by ID with caching
func (q *OptimizedQueries) GetPostByID(ctx context.Context, postID string) (*models.Post, error) {
	start := time.Now()
	defer func() {
		q.db.trackQuery("GetPostByID", time.Since(start))
	}()

	conn := q.db.getConnection()
	defer q.db.returnConnection(conn)

	postsDB := conn.DB("posts")

	var post models.Post
	row := postsDB.Get(ctx, postID)
	if err := row.ScanDoc(&post); err != nil {
		return nil, fmt.Errorf("failed to get post %s: %w", postID, err)
	}

	// Set document metadata
	post.ID = postID
	if rev, err := row.Rev(); err == nil && rev != "" {
		post.Rev = rev
	}

	return &post, nil
}

// GetFeaturedPosts retrieves featured posts efficiently
func (q *OptimizedQueries) GetFeaturedPosts(ctx context.Context, limit int) ([]models.Post, error) {
	start := time.Now()
	defer func() {
		q.db.trackQuery("GetFeaturedPosts", time.Since(start))
	}()

	conn := q.db.getConnection()
	defer q.db.returnConnection(conn)

	postsDB := conn.DB("posts")

	// Use the featured-status index for optimal performance
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"is_featured": true,
			"status":      "published",
		},
		"sort": []map[string]string{
			{"published_at": "desc"},
		},
		"limit":     limit,
		"use_index": "featured-status-index",
	}

	rows := postsDB.Find(ctx, query)
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.ScanDoc(&post); err != nil {
			continue
		}

		if id, err := rows.ID(); err == nil && id != "" {
			post.ID = id
		}
		if rev, err := rows.Rev(); err == nil && rev != "" {
			post.Rev = rev
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// GetPostsByAuthor retrieves posts by author with pagination
func (q *OptimizedQueries) GetPostsByAuthor(ctx context.Context, author string, page, limit int) ([]models.Post, error) {
	start := time.Now()
	defer func() {
		q.db.trackQuery("GetPostsByAuthor", time.Since(start))
	}()

	conn := q.db.getConnection()
	defer q.db.returnConnection(conn)

	postsDB := conn.DB("posts")

	// Use the author-created index
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"author": author,
			"status": "published",
		},
		"sort": []map[string]string{
			{"created_at": "desc"},
		},
		"limit":     limit,
		"skip":      (page - 1) * limit,
		"use_index": "author-created-index",
	}

	rows := postsDB.Find(ctx, query)
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.ScanDoc(&post); err != nil {
			continue
		}

		if id, err := rows.ID(); err == nil && id != "" {
			post.ID = id
		}
		if rev, err := rows.Rev(); err == nil && rev != "" {
			post.Rev = rev
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// getPostsCount gets the total count of posts for pagination
func (q *OptimizedQueries) getPostsCount(ctx context.Context, status string) (int, error) {
	conn := q.db.getConnection()
	defer q.db.returnConnection(conn)

	postsDB := conn.DB("posts")

	query := map[string]interface{}{
		"selector": map[string]interface{}{},
		"fields":   []string{"_id"}, // Only fetch IDs for counting
	}

	if status != "" {
		query["selector"].(map[string]interface{})["status"] = status
	}

	rows := postsDB.Find(ctx, query)
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	return count, nil
}

// User queries

// GetUserByUsernameOptimized retrieves user by username using index
func (q *OptimizedQueries) GetUserByUsernameOptimized(ctx context.Context, username string) (*models.User, error) {
	start := time.Now()
	defer func() {
		q.db.trackQuery("GetUserByUsername", time.Since(start))
	}()

	conn := q.db.getConnection()
	defer q.db.returnConnection(conn)

	usersDB := conn.DB("users")

	// Use the username index for optimal performance
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"username": username,
		},
		"limit":     1,
		"use_index": "username-index",
	}

	rows := usersDB.Find(ctx, query)
	defer rows.Close()

	if rows.Next() {
		var user models.User
		if err := rows.ScanDoc(&user); err != nil {
			return nil, fmt.Errorf("failed to scan user document: %w", err)
		}

		if id, err := rows.ID(); err == nil && id != "" {
			user.ID = id
		}
		if rev, err := rows.Rev(); err == nil && rev != "" {
			user.Rev = rev
		}

		return &user, nil
	}

	return nil, nil
}

// GetUserByEmailOptimized retrieves user by email using index
func (q *OptimizedQueries) GetUserByEmailOptimized(ctx context.Context, email string) (*models.User, error) {
	start := time.Now()
	defer func() {
		q.db.trackQuery("GetUserByEmail", time.Since(start))
	}()

	conn := q.db.getConnection()
	defer q.db.returnConnection(conn)

	usersDB := conn.DB("users")

	// Use the email index for optimal performance
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"email": email,
		},
		"limit":     1,
		"use_index": "email-index",
	}

	rows := usersDB.Find(ctx, query)
	defer rows.Close()

	if rows.Next() {
		var user models.User
		if err := rows.ScanDoc(&user); err != nil {
			return nil, fmt.Errorf("failed to scan user document: %w", err)
		}

		if id, err := rows.ID(); err == nil && id != "" {
			user.ID = id
		}
		if rev, err := rows.Rev(); err == nil && rev != "" {
			user.Rev = rev
		}

		return &user, nil
	}

	return nil, nil
}

// GetUsersPaginated retrieves users with pagination and role filtering
func (q *OptimizedQueries) GetUsersPaginated(ctx context.Context, role string, active *bool, page, limit int) (*models.PaginatedUsersResponse, error) {
	start := time.Now()
	defer func() {
		q.db.trackQuery("GetUsersPaginated", time.Since(start))
	}()

	conn := q.db.getConnection()
	defer q.db.returnConnection(conn)

	usersDB := conn.DB("users")

	// Build query with filters
	selector := map[string]interface{}{}
	if role != "" {
		selector["role"] = role
	}
	if active != nil {
		selector["active"] = *active
	}

	query := map[string]interface{}{
		"selector": selector,
		"sort": []map[string]string{
			{"created_at": "desc"},
		},
		"limit": limit,
		"skip":  (page - 1) * limit,
	}

	// Use appropriate index based on filters
	if role != "" || active != nil {
		query["use_index"] = "role-active-index"
	} else {
		query["use_index"] = "created-at-index"
	}

	rows := usersDB.Find(ctx, query)
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.ScanDoc(&user); err != nil {
			continue
		}

		// Don't return password hashes in lists
		user.PasswordHash = ""

		if id, err := rows.ID(); err == nil && id != "" {
			user.ID = id
		}
		if rev, err := rows.Rev(); err == nil && rev != "" {
			user.Rev = rev
		}

		users = append(users, user)
	}

	// Get total count
	total, err := q.getUsersCount(ctx, role, active)
	if err != nil {
		total = len(users) // Fallback
	}

	// Calculate pagination metadata
	totalPages := (total + limit - 1) / limit
	meta := models.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return &models.PaginatedUsersResponse{
		Data: users,
		Meta: meta,
	}, nil
}

// getUsersCount gets the total count of users for pagination
func (q *OptimizedQueries) getUsersCount(ctx context.Context, role string, active *bool) (int, error) {
	conn := q.db.getConnection()
	defer q.db.returnConnection(conn)

	usersDB := conn.DB("users")

	selector := map[string]interface{}{}
	if role != "" {
		selector["role"] = role
	}
	if active != nil {
		selector["active"] = *active
	}

	query := map[string]interface{}{
		"selector": selector,
		"fields":   []string{"_id"}, // Only fetch IDs for counting
	}

	rows := usersDB.Find(ctx, query)
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	return count, nil
}

// Contact queries

// GetContactsPaginated retrieves contacts with pagination and status filtering
func (q *OptimizedQueries) GetContactsPaginated(ctx context.Context, status string, page, limit int) ([]models.Contact, error) {
	start := time.Now()
	defer func() {
		q.db.trackQuery("GetContactsPaginated", time.Since(start))
	}()

	conn := q.db.getConnection()
	defer q.db.returnConnection(conn)

	contactsDB := conn.DB("contacts")

	// Build query with status filter
	selector := map[string]interface{}{}
	if status != "" {
		selector["status"] = status
	}

	query := map[string]interface{}{
		"selector": selector,
		"sort": []map[string]string{
			{"created_at": "desc"},
		},
		"limit":     limit,
		"skip":      (page - 1) * limit,
		"use_index": "status-created-index",
	}

	rows := contactsDB.Find(ctx, query)
	defer rows.Close()

	var contacts []models.Contact
	for rows.Next() {
		var contact models.Contact
		if err := rows.ScanDoc(&contact); err != nil {
			continue
		}

		if id, err := rows.ID(); err == nil && id != "" {
			contact.ID = id
		}
		if rev, err := rows.Rev(); err == nil && rev != "" {
			contact.Rev = rev
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// BatchOperations for bulk operations

// BulkUpdatePostStatus updates multiple posts status in a single operation
func (q *OptimizedQueries) BulkUpdatePostStatus(ctx context.Context, postIDs []string, status string) error {
	start := time.Now()
	defer func() {
		q.db.trackQuery("BulkUpdatePostStatus", time.Since(start))
	}()

	conn := q.db.getConnection()
	defer q.db.returnConnection(conn)

	postsDB := conn.DB("posts")

	// Prepare bulk update documents
	var docs []interface{}
	for _, postID := range postIDs {
		// Get current document
		var post models.Post
		row := postsDB.Get(ctx, postID)
		if err := row.ScanDoc(&post); err != nil {
			utils.LogError(err, "Failed to get post for bulk update", logrus.Fields{
				"post_id": postID,
			})
			continue
		}

		// Update status and timestamp
		post.Status = status
		post.UpdatedAt = time.Now()
		if status == "published" && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}

		// Set document metadata
		post.ID = postID
		if rev, err := row.Rev(); err == nil && rev != "" {
			post.Rev = rev
		}

		docs = append(docs, post)
	}

	// Perform bulk update
	if len(docs) > 0 {
		results, err := postsDB.BulkDocs(ctx, docs)
		if err != nil {
			return fmt.Errorf("bulk update failed: %w", err)
		}

		// Check for errors in results
		for _, result := range results {
			if result.Error != nil {
				utils.LogError(nil, "Bulk update document error", logrus.Fields{
					"error": result.Error.Error(),
					"id":    result.ID,
				})
			}
		}
	}

	return nil
}

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"webenable-cms-backend/database"
	"webenable-cms-backend/models"
	"webenable-cms-backend/utils"

	"github.com/go-kivik/kivik/v4"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GetPosts godoc
//
//	@Summary		Get all posts
//	@Description	Get all published posts with optional status filter and pagination
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			status	query		string	false	"Filter by post status (published, draft, scheduled)"
//	@Param			page	query		int		false	"Page number (default: 1)"
//	@Param			limit	query		int		false	"Items per page (default: 10, max: 100)"
//	@Success		200		{object}	models.PaginatedPostsResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/posts [get]
func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get pagination parameters
	page, limit := getPaginationParams(r)

	// Get status filter from query parameters
	statusFilter := r.URL.Query().Get("status")

	// Check if this is an authenticated request (admin access)
	isAuthenticated := r.Context().Value("user") != nil

	// Default to "published" only for public (non-authenticated) access
	if statusFilter == "" && !isAuthenticated {
		statusFilter = "published"
	}

	// Create cache key based on query parameters
	cacheKey := fmt.Sprintf("posts_list_status_%s_page_%d_limit_%d", statusFilter, page, limit)

	// Try to get from cache first
	if globalCache != nil {
		var cachedResponse models.PaginatedPostsResponse
		err := globalCache.GetCachedPostsList(cacheKey, &cachedResponse)
		if err == nil {
			w.Header().Set("X-Cache", "HIT")
			json.NewEncoder(w).Encode(cachedResponse)
			return
		}
	}

	ctx := context.Background()
	rows := database.Instance.PostsDB.AllDocs(ctx, kivik.Param("include_docs", true))
	defer rows.Close()

	var allPosts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.ScanDoc(&post); err != nil {
			continue
		}

		// Ensure the document ID and revision are set properly
		if id, err := rows.ID(); err == nil && id != "" {
			post.ID = id
		}
		if rev, err := rows.Rev(); err == nil && rev != "" {
			post.Rev = rev
		}

		// Filter by status if specified, or if public access (non-authenticated)
		if statusFilter != "" && post.Status != statusFilter {
			continue
		}

		allPosts = append(allPosts, post)
	}

	// Calculate pagination
	total := len(allPosts)
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	offset := (page - 1) * limit

	// Get the slice for current page
	var posts []models.Post
	if offset < total {
		end := offset + limit
		if end > total {
			end = total
		}
		posts = allPosts[offset:end]
	}

	// Create pagination metadata
	meta := models.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	response := models.PaginatedPostsResponse{
		Data: posts,
		Meta: meta,
	}

	// Cache the result for 10 minutes
	if globalCache != nil {
		go func() {
			err := globalCache.CachePostsList(cacheKey, response, 10*time.Minute)
			if err != nil {
				utils.LogError(err, "Failed to cache posts list", logrus.Fields{
					"cache_key": cacheKey,
				})
			}
		}()
	}

	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(response)
}

// GetPost godoc
//
//	@Summary		Get post by ID
//	@Description	Get a single post by ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	models.Post
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/posts/{id} [get]
func GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	// Try to get from cache first
	if globalCache != nil {
		var cachedPost models.Post
		err := globalCache.GetCachedPost(id, &cachedPost)
		if err == nil {
			w.Header().Set("X-Cache", "HIT")
			json.NewEncoder(w).Encode(cachedPost)
			return
		}
	}

	ctx := context.Background()
	row := database.Instance.PostsDB.Get(ctx, id)

	var post models.Post
	if err := row.ScanDoc(&post); err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Check if this is an authenticated request (admin access)
	isAuthenticated := r.Context().Value("user") != nil

	// Only return published posts for public (non-authenticated) access
	if !isAuthenticated && post.Status != "published" {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Ensure the document ID and revision are set properly
	post.ID = id
	if rev, err := row.Rev(); err == nil && rev != "" {
		post.Rev = rev
	}

	// Cache the post for 30 minutes
	if globalCache != nil {
		go func() {
			err := globalCache.CachePost(id, post, 30*time.Minute)
			if err != nil {
				utils.LogError(err, "Failed to cache post", logrus.Fields{
					"post_id": id,
				})
			}
		}()
	}

	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(post)
}

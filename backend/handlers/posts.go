package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"webenable-cms-backend/cache"
	"webenable-cms-backend/database"
	"webenable-cms-backend/models"

	"github.com/go-kivik/kivik/v4"
	"github.com/gorilla/mux"
)

var globalCache *cache.ValkeyClient

// SetGlobalCache sets the global cache client for handlers
func SetGlobalCache(valkeyClient *cache.ValkeyClient) {
	globalCache = valkeyClient
}

// GetPosts godoc
//
//	@Summary		Get all posts
//	@Description	Get all published posts with optional status filter
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			status	query		string	false	"Filter by post status (published, draft, scheduled)"
//	@Success		200		{array}		models.Post
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/posts [get]
func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get status filter from query parameters
	statusFilter := r.URL.Query().Get("status")

	// Create cache key based on query parameters
	cacheKey := "posts_list"
	if statusFilter != "" {
		cacheKey = fmt.Sprintf("posts_list_status_%s", statusFilter)
	}

	// Try to get from cache first
	if globalCache != nil {
		var cachedPosts []models.Post
		err := globalCache.GetCachedPostsList(cacheKey, &cachedPosts)
		if err == nil {
			w.Header().Set("X-Cache", "HIT")
			json.NewEncoder(w).Encode(cachedPosts)
			return
		}
	}

	ctx := context.Background()
	rows := database.Instance.PostsDB.AllDocs(ctx, kivik.Param("include_docs", true))
	defer rows.Close()

	var posts []models.Post
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

		// Filter by status if specified
		if statusFilter != "" && post.Status != statusFilter {
			continue
		}

		posts = append(posts, post)
	}

	// Cache the result for 10 minutes
	if globalCache != nil {
		go func() {
			err := globalCache.CachePostsList(cacheKey, posts, 10*time.Minute)
			if err != nil {
				fmt.Printf("Failed to cache posts list: %v\n", err)
			}
		}()
	}

	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(posts)
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
				fmt.Printf("Failed to cache post %s: %v\n", id, err)
			}
		}()
	}

	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(post)
}

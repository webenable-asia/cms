package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"webenable-cms-backend/database"
	"webenable-cms-backend/middleware"
	"webenable-cms-backend/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreatePost godoc
//
//	@Summary		Create new post
//	@Description	Create a new post (authenticated users only)
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			post	body		models.Post	true	"Post data"
//	@Success		201		{object}	models.Post
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/posts [post]
func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from context (set by auth middleware)
	claims := r.Context().Value("user").(*middleware.Claims)
	post.Author = claims.Username
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	if post.Status == "" {
		post.Status = "draft"
	}

	// Generate a UUID for the document ID
	postID := uuid.New().String()
	post.ID = postID

	// Create a document map with proper CouchDB fields
	doc := map[string]interface{}{
		"_id":        postID,
		"title":      post.Title,
		"content":    post.Content,
		"excerpt":    post.Excerpt,
		"author":     post.Author,
		"status":     post.Status,
		"tags":       post.Tags,
		"created_at": post.CreatedAt,
		"updated_at": post.UpdatedAt,
	}

	if post.PublishedAt != nil {
		doc["published_at"] = post.PublishedAt
	}

	ctx := context.Background()
	rev, err := database.Instance.PostsDB.Put(ctx, postID, doc)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	post.Rev = rev

	// Invalidate posts list cache
	if globalCache != nil {
		go func() {
			err := globalCache.InvalidatePostsListCache()
			if err != nil {
				fmt.Printf("Failed to invalidate posts list cache: %v\n", err)
			}
		}()
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// UpdatePost godoc
//
//	@Summary		Update post
//	@Description	Update an existing post (authenticated users only)
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string		true	"Post ID"
//	@Param			post	body		models.Post	true	"Post data"
//	@Success		200		{object}	models.Post
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/posts/{id} [put]
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var updatedPost models.Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Get existing post
	row := database.Instance.PostsDB.Get(ctx, id)
	var existingPost models.Post
	if err := row.ScanDoc(&existingPost); err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Ensure the document ID and revision are set properly
	existingPost.ID = id
	if rev, err := row.Rev(); err == nil && rev != "" {
		existingPost.Rev = rev
	}

	// Update fields
	existingPost.Title = updatedPost.Title
	existingPost.Content = updatedPost.Content
	existingPost.Excerpt = updatedPost.Excerpt
	existingPost.Status = updatedPost.Status
	existingPost.Tags = updatedPost.Tags
	existingPost.UpdatedAt = time.Now()

	if updatedPost.Status == "published" && existingPost.PublishedAt == nil {
		now := time.Now()
		existingPost.PublishedAt = &now
	}

	// Create a document map with proper CouchDB fields
	doc := map[string]interface{}{
		"_id":        existingPost.ID,
		"_rev":       existingPost.Rev,
		"title":      existingPost.Title,
		"content":    existingPost.Content,
		"excerpt":    existingPost.Excerpt,
		"author":     existingPost.Author,
		"status":     existingPost.Status,
		"tags":       existingPost.Tags,
		"created_at": existingPost.CreatedAt,
		"updated_at": existingPost.UpdatedAt,
	}

	if existingPost.PublishedAt != nil {
		doc["published_at"] = existingPost.PublishedAt
	}

	// Update in database
	rev, err := database.Instance.PostsDB.Put(ctx, id, doc)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	existingPost.Rev = rev

	// Invalidate caches
	if globalCache != nil {
		go func() {
			// Invalidate specific post cache
			err := globalCache.InvalidatePostCache(id)
			if err != nil {
				fmt.Printf("Failed to invalidate post cache for %s: %v\n", id, err)
			}
			// Invalidate posts list cache
			err = globalCache.InvalidatePostsListCache()
			if err != nil {
				fmt.Printf("Failed to invalidate posts list cache: %v\n", err)
			}
		}()
	}

	json.NewEncoder(w).Encode(existingPost)
}

// DeletePost godoc
//
//	@Summary		Delete post
//	@Description	Delete a post (authenticated users only)
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	models.SuccessResponse
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/posts/{id} [delete]
func DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()

	// Get existing post to get revision
	row := database.Instance.PostsDB.Get(ctx, id)
	var post models.Post
	if err := row.ScanDoc(&post); err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Ensure we have the proper revision
	post.ID = id
	if rev, err := row.Rev(); err == nil && rev != "" {
		post.Rev = rev
	}

	// Delete the post
	_, err := database.Instance.PostsDB.Delete(ctx, id, post.Rev)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	// Invalidate caches
	if globalCache != nil {
		go func() {
			// Invalidate specific post cache
			err := globalCache.InvalidatePostCache(id)
			if err != nil {
				fmt.Printf("Failed to invalidate post cache for %s: %v\n", id, err)
			}
			// Invalidate posts list cache
			err = globalCache.InvalidatePostsListCache()
			if err != nil {
				fmt.Printf("Failed to invalidate posts list cache: %v\n", err)
			}
		}()
	}

	response := map[string]string{"message": "Post deleted successfully"}
	json.NewEncoder(w).Encode(response)
}

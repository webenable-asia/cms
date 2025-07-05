package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"webenable-cms-backend/database"
	"webenable-cms-backend/models"

	"github.com/go-kivik/kivik/v4"
	"github.com/gorilla/mux"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get status filter from query parameters
	statusFilter := r.URL.Query().Get("status")
	
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

	json.NewEncoder(w).Encode(posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	id := vars["id"]

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

	json.NewEncoder(w).Encode(post)
}

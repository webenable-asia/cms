package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
)

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

	// Check if the posts database exists
	dbExists, err := client.DBExists(ctx, "posts")
	if err != nil {
		log.Fatal("Failed to check if database exists:", err)
	}

	if !dbExists {
		fmt.Println("Posts database does not exist. Nothing to delete.")
		return
	}

	db := client.DB("posts")
	// Get all documents
	rows := db.AllDocs(ctx, kivik.Param("include_docs", true))
	defer rows.Close()

	var documentsToDelete []struct {
		ID  string `json:"_id"`
		Rev string `json:"_rev"`
	}

	// Collect all document IDs and revisions
	for rows.Next() {
		var doc map[string]interface{}
		if err := rows.ScanDoc(&doc); err != nil {
			log.Printf("Failed to scan document: %v", err)
			continue
		}

		id, ok := doc["_id"].(string)
		if !ok {
			continue
		}

		rev, ok := doc["_rev"].(string)
		if !ok {
			continue
		}

		// Skip design documents (they start with _design/)
		if len(id) > 8 && id[:8] == "_design/" {
			continue
		}

		documentsToDelete = append(documentsToDelete, struct {
			ID  string `json:"_id"`
			Rev string `json:"_rev"`
		}{
			ID:  id,
			Rev: rev,
		})
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating through documents:", err)
	}

	if len(documentsToDelete) == 0 {
		fmt.Println("No blog posts found to delete.")
		return
	}

	fmt.Printf("Found %d blog posts to delete.\n", len(documentsToDelete))

	// Delete each document
	deletedCount := 0
	for _, doc := range documentsToDelete {
		_, err := db.Delete(ctx, doc.ID, doc.Rev)
		if err != nil {
			log.Printf("Failed to delete document %s: %v", doc.ID, err)
			continue
		}

		fmt.Printf("Deleted post: %s\n", doc.ID)
		deletedCount++
	}

	fmt.Printf("\nSuccessfully deleted %d out of %d blog posts.\n", deletedCount, len(documentsToDelete))

	if deletedCount < len(documentsToDelete) {
		fmt.Printf("Failed to delete %d posts. Check the logs above for details.\n", len(documentsToDelete)-deletedCount)
	}
}

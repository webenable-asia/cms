package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	couchdbURL := os.Getenv("COUCHDB_URL")
	if couchdbURL == "" {
		log.Fatal("COUCHDB_URL environment variable is required")
	}

	// Get all users first
	resp, err := http.Get(couchdbURL + "users/_all_docs?include_docs=true")
	if err != nil {
		log.Printf("Error querying users: %v", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Printf("Error decoding response: %v", err)
		return
	}

	rows, ok := result["rows"].([]interface{})
	if !ok {
		log.Printf("No rows found or invalid format")
		return
	}

	fmt.Printf("Found %d users in database\n", len(rows))

	// Find and delete users with username "admin" (keep only "webenable_admin")
	for _, row := range rows {
		rowMap := row.(map[string]interface{})
		doc := rowMap["doc"].(map[string]interface{})

		username, hasUsername := doc["username"].(string)
		id, hasId := doc["_id"].(string)
		rev, hasRev := doc["_rev"].(string)

		if !hasUsername || !hasId || !hasRev {
			continue
		}

		fmt.Printf("Found user: %s (ID: %s)\n", username, id)

		// Delete users with username "admin" but keep "webenable_admin"
		if username == "admin" {
			fmt.Printf("Deleting user: %s\n", username)

			// Delete the document
			client := &http.Client{}
			deleteURL := fmt.Sprintf("%susers/%s?rev=%s", couchdbURL, id, rev)
			req, err := http.NewRequest("DELETE", deleteURL, nil)
			if err != nil {
				log.Printf("Error creating delete request: %v", err)
				continue
			}

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error deleting user %s: %v", username, err)
				continue
			}
			resp.Body.Close()

			if resp.StatusCode == 200 || resp.StatusCode == 202 {
				fmt.Printf("Successfully deleted user: %s\n", username)
			} else {
				fmt.Printf("Failed to delete user %s: %d\n", username, resp.StatusCode)
			}
		}
	}

	fmt.Println("Cleanup completed")
}

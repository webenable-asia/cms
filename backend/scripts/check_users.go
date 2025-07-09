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

	// Query all users to see what exists
	resp, err := http.Get(couchdbURL + "webenable/_design/users/_view/by_username")
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

	fmt.Printf("Existing users: %+v\n", result)

	// Also try to get all documents in the database
	resp2, err := http.Get(couchdbURL + "webenable/_all_docs?include_docs=true")
	if err != nil {
		log.Printf("Error querying all docs: %v", err)
		return
	}
	defer resp2.Body.Close()

	var result2 map[string]interface{}
	err = json.NewDecoder(resp2.Body).Decode(&result2)
	if err != nil {
		log.Printf("Error decoding all docs response: %v", err)
		return
	}

	fmt.Printf("\nAll documents: %+v\n", result2)
}

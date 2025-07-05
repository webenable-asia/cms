package database

import (
	"context"
	"log"
	"os"

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
)

type DB struct {
	Client     *kivik.Client
	PostsDB    *kivik.DB
	UsersDB    *kivik.DB
	ContactsDB *kivik.DB
}

var Instance *DB

func Init() {
	couchdbURL := os.Getenv("COUCHDB_URL")
	if couchdbURL == "" {
		couchdbURL = "http://admin:password@localhost:5984/"
	}

	client, err := kivik.New("couch", couchdbURL)
	if err != nil {
		log.Fatal("Failed to connect to CouchDB:", err)
	}

	// Create databases if they don't exist
	ctx := context.Background()
	
	// Create posts database
	if exists, _ := client.DBExists(ctx, "posts"); !exists {
		if err := client.CreateDB(ctx, "posts"); err != nil {
			log.Fatal("Failed to create posts database:", err)
		}
	}

	// Create users database
	if exists, _ := client.DBExists(ctx, "users"); !exists {
		if err := client.CreateDB(ctx, "users"); err != nil {
			log.Fatal("Failed to create users database:", err)
		}
	}

	// Create contacts database
	if exists, _ := client.DBExists(ctx, "contacts"); !exists {
		if err := client.CreateDB(ctx, "contacts"); err != nil {
			log.Fatal("Failed to create contacts database:", err)
		}
	}

	Instance = &DB{
		Client:     client,
		PostsDB:    client.DB("posts"),
		UsersDB:    client.DB("users"),
		ContactsDB: client.DB("contacts"),
	}

	log.Println("Database initialized successfully")
}

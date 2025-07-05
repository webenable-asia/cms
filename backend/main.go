package main

import (
	"log"
	"net/http"
	"os"

	"webenable-cms-backend/database"
	"webenable-cms-backend/handlers"
	"webenable-cms-backend/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/go-kivik/kivik/v4/couchdb"
)

func main() {
	// Initialize database
	database.Init()

	// Initialize router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Public routes
	api.HandleFunc("/posts", handlers.GetPosts).Methods("GET")
	api.HandleFunc("/posts/{id}", handlers.GetPost).Methods("GET")
	api.HandleFunc("/auth/login", handlers.Login).Methods("POST")
	api.HandleFunc("/contact", handlers.SubmitContact).Methods("POST")

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/posts", handlers.CreatePost).Methods("POST")
	protected.HandleFunc("/posts/{id}", handlers.UpdatePost).Methods("PUT")
	protected.HandleFunc("/posts/{id}", handlers.DeletePost).Methods("DELETE")
	protected.HandleFunc("/contacts", handlers.GetContacts).Methods("GET")
	protected.HandleFunc("/contacts/{id}", handlers.GetContact).Methods("GET")
	protected.HandleFunc("/contacts/{id}", handlers.UpdateContactStatus).Methods("PUT")
	protected.HandleFunc("/contacts/{id}/reply", handlers.ReplyToContact).Methods("POST")
	protected.HandleFunc("/contacts/{id}", handlers.DeleteContact).Methods("DELETE")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

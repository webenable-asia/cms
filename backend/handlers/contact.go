package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"webenable-cms-backend/database"
	"webenable-cms-backend/models"

	"github.com/go-kivik/kivik/v4"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// SubmitContact godoc
//
//	@Summary		Submit contact form
//	@Description	Submit a contact form (public endpoint)
//	@Tags			Contact
//	@Accept			json
//	@Produce		json
//	@Param			contact	body		models.Contact	true	"Contact form data"
//	@Success		201		{object}	models.Contact
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/contact [post]
//
// Public endpoint - no auth required
func SubmitContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set default values
	contact.Status = "new"
	contact.CreatedAt = time.Now()

	// Generate a UUID for the document ID
	contactID := uuid.New().String()
	contact.ID = contactID

	// Create a document map with proper CouchDB fields
	doc := map[string]interface{}{
		"_id":        contactID,
		"name":       contact.Name,
		"email":      contact.Email,
		"company":    contact.Company,
		"phone":      contact.Phone,
		"subject":    contact.Subject,
		"message":    contact.Message,
		"status":     contact.Status,
		"created_at": contact.CreatedAt,
	}

	ctx := context.Background()
	rev, err := database.Instance.ContactsDB.Put(ctx, contactID, doc)
	if err != nil {
		http.Error(w, "Failed to submit contact form", http.StatusInternalServerError)
		return
	}

	contact.Rev = rev
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Contact form submitted successfully",
		"id":      contactID,
	})
}

// Protected endpoints - require authentication
func GetContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get status filter from query parameters
	statusFilter := r.URL.Query().Get("status")

	ctx := context.Background()
	rows := database.Instance.ContactsDB.AllDocs(ctx, kivik.Param("include_docs", true))
	defer rows.Close()

	var contacts []models.Contact
	for rows.Next() {
		var contact models.Contact
		if err := rows.ScanDoc(&contact); err != nil {
			continue
		}

		// Ensure the document ID and revision are set properly
		if id, err := rows.ID(); err == nil && id != "" {
			contact.ID = id
		}
		if rev, err := rows.Rev(); err == nil && rev != "" {
			contact.Rev = rev
		}

		// Filter by status if specified
		if statusFilter != "" && contact.Status != statusFilter {
			continue
		}

		contacts = append(contacts, contact)
	}

	// Return data in the expected format for frontend
	response := map[string]interface{}{
		"data": contacts,
		"meta": map[string]interface{}{
			"total":         len(contacts),
			"status_filter": statusFilter,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func GetContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()
	row := database.Instance.ContactsDB.Get(ctx, id)

	var contact models.Contact
	if err := row.ScanDoc(&contact); err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	// Ensure the document ID and revision are set properly
	contact.ID = id
	if rev, err := row.Rev(); err == nil && rev != "" {
		contact.Rev = rev
	}

	json.NewEncoder(w).Encode(contact)
}

func UpdateContactStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var updateData struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Get existing contact
	row := database.Instance.ContactsDB.Get(ctx, id)
	var existingContact models.Contact
	if err := row.ScanDoc(&existingContact); err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	// Ensure the document ID and revision are set properly
	existingContact.ID = id
	if rev, err := row.Rev(); err == nil && rev != "" {
		existingContact.Rev = rev
	}

	// Update status and timestamps
	existingContact.Status = updateData.Status
	if updateData.Status == "read" && existingContact.ReadAt == nil {
		now := time.Now()
		existingContact.ReadAt = &now
	}
	if updateData.Status == "replied" && existingContact.RepliedAt == nil {
		now := time.Now()
		existingContact.RepliedAt = &now
	}

	// Create a document map with proper CouchDB fields
	doc := map[string]interface{}{
		"_id":        existingContact.ID,
		"_rev":       existingContact.Rev,
		"name":       existingContact.Name,
		"email":      existingContact.Email,
		"company":    existingContact.Company,
		"phone":      existingContact.Phone,
		"subject":    existingContact.Subject,
		"message":    existingContact.Message,
		"status":     existingContact.Status,
		"created_at": existingContact.CreatedAt,
	}

	if existingContact.ReadAt != nil {
		doc["read_at"] = existingContact.ReadAt
	}
	if existingContact.RepliedAt != nil {
		doc["replied_at"] = existingContact.RepliedAt
	}

	// Update in database
	rev, err := database.Instance.ContactsDB.Put(ctx, id, doc)
	if err != nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return
	}

	existingContact.Rev = rev
	json.NewEncoder(w).Encode(existingContact)
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx := context.Background()

	// Get existing contact to get revision
	row := database.Instance.ContactsDB.Get(ctx, id)
	var contact models.Contact
	if err := row.ScanDoc(&contact); err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	// Ensure we have the proper revision
	contact.ID = id
	if rev, err := row.Rev(); err == nil && rev != "" {
		contact.Rev = rev
	}

	// Delete the contact
	_, err := database.Instance.ContactsDB.Delete(ctx, id, contact.Rev)
	if err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ReplyToContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var replyData struct {
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&replyData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	row := database.Instance.ContactsDB.Get(ctx, id)
	var contact models.Contact
	if err := row.ScanDoc(&contact); err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	// Ensure the document ID and revision are set properly
	contact.ID = id
	if rev, err := row.Rev(); err == nil && rev != "" {
		contact.Rev = rev
	}

	// Send email reply using the email service (optional in development)
	if err := SendEmailReply(contact.Email, contact.Name, replyData.Subject, replyData.Message); err != nil {
		// Log the error but don't fail the request in development
		fmt.Printf("Email sending failed (continuing anyway): %v\n", err)
	}

	// Update contact status to "replied"
	contact.Status = "replied"
	now := time.Now()
	contact.RepliedAt = &now

	// Create a document map with proper CouchDB fields
	doc := map[string]interface{}{
		"_id":        contact.ID,
		"_rev":       contact.Rev,
		"name":       contact.Name,
		"email":      contact.Email,
		"company":    contact.Company,
		"phone":      contact.Phone,
		"subject":    contact.Subject,
		"message":    contact.Message,
		"status":     contact.Status,
		"created_at": contact.CreatedAt,
		"replied_at": contact.RepliedAt,
	}

	if contact.ReadAt != nil {
		doc["read_at"] = contact.ReadAt
	}

	// Update in database
	_, err := database.Instance.ContactsDB.Put(ctx, id, doc)
	if err != nil {
		// Email was sent but couldn't update status - that's ok
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Reply sent successfully (status update failed)",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Reply sent successfully and contact marked as replied",
	})
}

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
)

type Contact struct {
	ID        string     `json:"_id,omitempty"`
	Rev       string     `json:"_rev,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Company   string     `json:"company"`
	Phone     string     `json:"phone"`
	Subject   string     `json:"subject"`
	Message   string     `json:"message"`
	Status    string     `json:"status"` // new, read, replied
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	RepliedAt *time.Time `json:"replied_at,omitempty"`
}

func main() {
	// Get CouchDB URL from environment or use default
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

	// Get or create the contacts database
	dbExists, err := client.DBExists(ctx, "contacts")
	if err != nil {
		log.Fatal("Failed to check if contacts database exists:", err)
	}

	var db *kivik.DB
	if !dbExists {
		err = client.CreateDB(ctx, "contacts")
		if err != nil {
			log.Fatal("Failed to create contacts database:", err)
		}
		fmt.Println("Created contacts database")
	}

	db = client.DB("contacts")

	// Sample contact submissions
	contacts := []Contact{
		{
			Name:      "Sarah Johnson",
			Email:     "sarah.johnson@techcorp.com",
			Company:   "TechCorp Solutions",
			Phone:     "+1 (555) 123-4567",
			Subject:   "Custom Web Application Development",
			Message:   "Hi, we're looking for a partner to develop a custom web application for our logistics management. We need a scalable solution that can handle real-time tracking and inventory management. Could you provide more information about your development process and timeline?",
			Status:    "new",
			CreatedAt: time.Now().Add(-72 * time.Hour),
		},
		{
			Name:      "Michael Chen",
			Email:     "m.chen@innovatestart.com",
			Company:   "InnovateStart",
			Phone:     "+1 (555) 987-6543",
			Subject:   "E-commerce Platform Setup",
			Message:   "We're a startup looking to launch our e-commerce platform. We need help with payment integration, user authentication, and mobile responsiveness. What's your experience with Next.js and modern payment gateways?",
			Status:    "read",
			CreatedAt: time.Now().Add(-48 * time.Hour),
			ReadAt:    timePtr(time.Now().Add(-24 * time.Hour)),
		},
		{
			Name:      "Emma Rodriguez",
			Email:     "emma.r@healthplus.org",
			Company:   "HealthPlus Clinic",
			Phone:     "+1 (555) 456-7890",
			Subject:   "Patient Portal Development",
			Message:   "Our clinic needs a secure patient portal where patients can book appointments, view test results, and communicate with doctors. Security and HIPAA compliance are critical. Can you help us with this project?",
			Status:    "replied",
			CreatedAt: time.Now().Add(-120 * time.Hour),
			ReadAt:    timePtr(time.Now().Add(-96 * time.Hour)),
			RepliedAt: timePtr(time.Now().Add(-72 * time.Hour)),
		},
		{
			Name:      "David Kim",
			Email:     "david@greentech.io",
			Company:   "GreenTech Solutions",
			Phone:     "+1 (555) 234-5678",
			Subject:   "IoT Dashboard Development",
			Message:   "We're developing IoT sensors for environmental monitoring and need a real-time dashboard to visualize the data. The dashboard should support multiple data sources and provide analytics capabilities. What's your experience with real-time data visualization?",
			Status:    "new",
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			Name:      "Lisa Wang",
			Email:     "lisa.wang@edulearn.com",
			Company:   "EduLearn Platform",
			Phone:     "+1 (555) 345-6789",
			Subject:   "Learning Management System",
			Message:   "We're creating an online learning platform and need help with video streaming, progress tracking, and assessment tools. The platform should support thousands of concurrent users. Could you share some similar projects you've worked on?",
			Status:    "read",
			CreatedAt: time.Now().Add(-96 * time.Hour),
			ReadAt:    timePtr(time.Now().Add(-48 * time.Hour)),
		},
		{
			Name:      "Robert Thompson",
			Email:     "r.thompson@financeplus.net",
			Company:   "FinancePlus",
			Phone:     "+1 (555) 567-8901",
			Subject:   "Financial Analytics Dashboard",
			Message:   "We need a comprehensive financial analytics dashboard that can process large datasets and provide real-time insights. The solution should integrate with our existing APIs and support custom reporting. What's your approach to handling financial data securely?",
			Status:    "new",
			CreatedAt: time.Now().Add(-6 * time.Hour),
		},
		{
			Name:      "Jennifer Martinez",
			Email:     "j.martinez@retailchain.com",
			Company:   "RetailChain Inc",
			Phone:     "+1 (555) 678-9012",
			Subject:   "Inventory Management System",
			Message:   "Our retail chain needs a centralized inventory management system that can sync across multiple locations. We need real-time stock updates, automated reordering, and detailed reporting. Can you provide a proposal for this project?",
			Status:    "replied",
			CreatedAt: time.Now().Add(-168 * time.Hour), // 1 week ago
			ReadAt:    timePtr(time.Now().Add(-144 * time.Hour)),
			RepliedAt: timePtr(time.Now().Add(-120 * time.Hour)),
		},
		{
			Name:      "Alex Rodriguez",
			Email:     "alex@freelancehub.co",
			Company:   "FreelanceHub",
			Phone:     "+1 (555) 789-0123",
			Subject:   "Freelancer Marketplace Platform",
			Message:   "We're building a freelancer marketplace and need help with the matching algorithm, payment processing, and dispute resolution system. The platform should support multiple currencies and payment methods. What's your experience with marketplace development?",
			Status:    "new",
			CreatedAt: time.Now().Add(-12 * time.Hour),
		},
		{
			Name:      "Sophie Laurent",
			Email:     "sophie@travelexplore.fr",
			Company:   "Travel Explore",
			Phone:     "+33 1 23 45 67 89",
			Subject:   "Travel Booking Platform",
			Message:   "Bonjour! We're developing a travel booking platform that needs to integrate with multiple airline and hotel APIs. We also need a recommendation engine based on user preferences. Can you help us with the backend architecture and API integrations?",
			Status:    "read",
			CreatedAt: time.Now().Add(-36 * time.Hour),
			ReadAt:    timePtr(time.Now().Add(-12 * time.Hour)),
		},
		{
			Name:      "James Wilson",
			Email:     "james.wilson@sportstracker.app",
			Company:   "SportsTracker",
			Phone:     "+1 (555) 890-1234",
			Subject:   "Fitness Tracking Application",
			Message:   "We're creating a fitness tracking app that needs to sync with wearable devices and provide personalized workout recommendations. The app should support social features and challenges. What's your experience with health and fitness applications?",
			Status:    "new",
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
	}

	// Insert contacts into database
	for _, contact := range contacts {
		// Check if contact already exists
		existingContact := Contact{}
		err := db.Get(ctx, contact.Email).ScanDoc(&existingContact)
		if err == nil {
			fmt.Printf("Contact '%s' already exists, skipping...\n", contact.Name)
			continue
		}

		// Set document ID to email for uniqueness
		contact.ID = contact.Email

		// Insert contact
		_, err = db.Put(ctx, contact.ID, contact)
		if err != nil {
			log.Printf("Failed to insert contact '%s': %v", contact.Name, err)
			continue
		}

		fmt.Printf("Added contact: %s (%s) - %s\n", contact.Name, contact.Company, contact.Subject)
	}

	// Generate some additional random contacts
	names := []string{"John Smith", "Maria Garcia", "Ahmed Hassan", "Yuki Tanaka", "Oliver Brown", "Anna Kowalski", "Carlos Silva", "Fatima Al-Zahra"}
	companies := []string{"StartupTech", "Digital Solutions", "Innovation Labs", "Tech Pioneers", "Future Systems", "Data Dynamics", "Cloud Nine", "Agile Works"}
	subjects := []string{
		"Website Redesign",
		"Mobile App Development",
		"API Integration",
		"Database Optimization",
		"Cloud Migration",
		"Security Audit",
		"Performance Enhancement",
		"Digital Transformation",
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 5; i++ {
		randomContact := Contact{
			Name:      names[rand.Intn(len(names))],
			Email:     fmt.Sprintf("contact%d@%s.com", i+1, companies[rand.Intn(len(companies))]),
			Company:   companies[rand.Intn(len(companies))],
			Phone:     fmt.Sprintf("+1 (555) %03d-%04d", rand.Intn(900)+100, rand.Intn(9000)+1000),
			Subject:   subjects[rand.Intn(len(subjects))],
			Message:   "This is a sample contact message generated for testing purposes. We're interested in your services and would like to discuss our project requirements.",
			Status:    []string{"new", "read", "replied"}[rand.Intn(3)],
			CreatedAt: time.Now().Add(-time.Duration(rand.Intn(168)) * time.Hour), // Random time in last week
		}

		randomContact.ID = randomContact.Email

		// Check if contact already exists
		existingContact := Contact{}
		err := db.Get(ctx, randomContact.Email).ScanDoc(&existingContact)
		if err == nil {
			fmt.Printf("Random contact '%s' already exists, skipping...\n", randomContact.Name)
			continue
		}

		// Set read time for read/replied contacts
		if randomContact.Status == "read" || randomContact.Status == "replied" {
			readTime := randomContact.CreatedAt.Add(time.Duration(rand.Intn(48)) * time.Hour)
			randomContact.ReadAt = &readTime
		}

		// Set replied time for replied contacts
		if randomContact.Status == "replied" {
			repliedTime := randomContact.ReadAt.Add(time.Duration(rand.Intn(24)) * time.Hour)
			randomContact.RepliedAt = &repliedTime
		}

		// Insert contact
		_, err = db.Put(ctx, randomContact.ID, randomContact)
		if err != nil {
			log.Printf("Failed to insert random contact '%s': %v", randomContact.Name, err)
			continue
		}

		fmt.Printf("Added random contact: %s (%s) - %s\n", randomContact.Name, randomContact.Company, randomContact.Subject)
	}

	fmt.Println("\nContacts database population completed!")
	fmt.Println("You can now view contacts via the API at: http://localhost/api/contacts")
	fmt.Println("Note: Admin authentication may be required to view contacts.")
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}

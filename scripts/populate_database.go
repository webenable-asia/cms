package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
)

type Post struct {
	ID            string     `json:"_id,omitempty"`
	Rev           string     `json:"_rev,omitempty"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	Excerpt       string     `json:"excerpt"`
	Author        string     `json:"author"`
	Status        string     `json:"status"`
	Tags          []string   `json:"tags"`
	Categories    []string   `json:"categories"`
	FeaturedImage string     `json:"featured_image"`
	ImageAlt      string     `json:"image_alt"`
	MetaTitle     string     `json:"meta_title"`
	MetaDesc      string     `json:"meta_description"`
	ReadingTime   int        `json:"reading_time"`
	IsFeatured    bool       `json:"is_featured"`
	ViewCount     int        `json:"view_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
	ScheduledAt   *time.Time `json:"scheduled_at,omitempty"`
}

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

	// Get or create the posts database
	dbExists, err := client.DBExists(ctx, "posts")
	if err != nil {
		log.Fatal("Failed to check if database exists:", err)
	}

	var db *kivik.DB
	if !dbExists {
		err = client.CreateDB(ctx, "posts")
		if err != nil {
			log.Fatal("Failed to create posts database:", err)
		}
		fmt.Println("Created posts database")
	}

	db = client.DB("posts")

	// Sample blog posts based on webenable.asia content
	posts := []Post{
		{
			ID:    "getting-started-modern-web-development",
			Title: "Getting Started with Modern Web Development",
			Content: `# Getting Started with Modern Web Development

## Introduction

Modern web development has evolved significantly over the past few years. With new frameworks, tools, and best practices emerging constantly, it can be overwhelming for developers to keep up.

## Key Technologies

Here are some of the most important technologies every modern web developer should know:

- **React/Next.js** - For building user interfaces and full-stack applications
- **TypeScript** - For type-safe JavaScript development
- **Tailwind CSS** - For utility-first styling
- **Vercel/Netlify** - For deployment and hosting

## Best Practices

Following these best practices will help you build better web applications:

1. Write clean, maintainable code
2. Implement proper error handling
3. Optimize for performance
4. Ensure accessibility compliance
5. Test your applications thoroughly

## Conclusion

Modern web development is an exciting field with endless possibilities. By staying up-to-date with the latest technologies and best practices, you can build amazing web experiences that delight users.`,
			Excerpt:       "Learn the fundamentals of building modern web applications with the latest technologies and best practices.",
			Author:        "WebEnable Team",
			Status:        "published",
			Tags:          []string{"webdev", "javascript", "react", "typescript"},
			Categories:    []string{"Development"},
			FeaturedImage: "https://images.unsplash.com/photo-1461749280684-dccba630e2f6?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2669&q=80",
			ImageAlt:      "Modern web development setup with code on multiple screens",
			MetaTitle:     "Getting Started with Modern Web Development | WebEnable",
			MetaDesc:      "Learn the fundamentals of building modern web applications with the latest technologies and best practices.",
			ReadingTime:   5,
			IsFeatured:    true,
			ViewCount:     142,
			CreatedAt:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:    "future-digital-marketing",
			Title: "The Future of Digital Marketing",
			Content: `# The Future of Digital Marketing

## Introduction

Digital marketing is rapidly evolving with new technologies and changing consumer behaviors. Understanding these trends is crucial for businesses looking to stay competitive.

## Emerging Trends

### AI and Machine Learning
- Personalized content recommendations
- Automated ad optimization
- Predictive analytics for customer behavior

### Voice Search Optimization
- Growing adoption of smart speakers
- Conversational search queries
- Local business opportunities

### Interactive Content
- Augmented reality experiences
- Interactive videos and polls
- Gamification strategies

## Strategies for Success

1. **Focus on First-Party Data** - Build direct relationships with customers
2. **Embrace Omnichannel Marketing** - Consistent experience across all touchpoints
3. **Invest in Video Content** - Short-form and live video continue to dominate
4. **Prioritize Mobile Experience** - Mobile-first approach is essential

## Conclusion

The future of digital marketing lies in creating authentic, personalized experiences that add value to customers' lives. Businesses that adapt to these changes will thrive in the digital landscape.`,
			Excerpt:       "Explore emerging trends and strategies that will shape the digital marketing landscape in the coming years.",
			Author:        "Marketing Team",
			Status:        "published",
			Tags:          []string{"marketing", "digital", "trends", "strategy"},
			Categories:    []string{"Marketing"},
			FeaturedImage: "https://images.unsplash.com/photo-1460925895917-afdab827c52f?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2015&q=80",
			ImageAlt:      "Digital marketing analytics dashboard showing various metrics",
			MetaTitle:     "The Future of Digital Marketing | WebEnable",
			MetaDesc:      "Explore emerging trends and strategies that will shape the digital marketing landscape in the coming years.",
			ReadingTime:   8,
			IsFeatured:    false,
			ViewCount:     98,
			CreatedAt:     time.Date(2024, 1, 10, 14, 30, 0, 0, time.UTC),
			UpdatedAt:     time.Date(2024, 1, 10, 14, 30, 0, 0, time.UTC),
		},
		{
			ID:    "building-scalable-business-solutions",
			Title: "Building Scalable Business Solutions",
			Content: `# Building Scalable Business Solutions

## Introduction

As businesses grow, their technology needs evolve. Building scalable solutions from the start can save time, money, and headaches down the road.

## Key Principles of Scalability

### Architecture Design
- **Microservices Architecture** - Break applications into smaller, manageable services
- **Cloud-Native Approach** - Leverage cloud platforms for flexibility and scalability
- **API-First Design** - Build with integration and expansion in mind

### Technology Choices
- **Database Selection** - Choose databases that can handle growth
- **Caching Strategies** - Implement effective caching for performance
- **Load Balancing** - Distribute traffic efficiently

## Implementation Strategies

1. **Start with MVP** - Build minimum viable product first
2. **Monitor Performance** - Track key metrics from day one
3. **Plan for Growth** - Design with future expansion in mind
4. **Automate Processes** - Reduce manual work through automation

## Common Pitfalls to Avoid

- Over-engineering from the start
- Ignoring security considerations
- Not planning for data migration
- Underestimating maintenance needs

## Conclusion

Building scalable business solutions requires careful planning, the right technology choices, and a focus on long-term growth. Invest in scalability early to ensure your business can adapt to changing market needs.`,
			Excerpt:       "Discover how to create business solutions that grow with your company and adapt to changing market needs.",
			Author:        "Business Solutions Team",
			Status:        "published",
			Tags:          []string{"business", "scalability", "architecture", "growth"},
			Categories:    []string{"Business"},
			FeaturedImage: "https://images.unsplash.com/photo-1551434678-e076c223a692?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2070&q=80",
			ImageAlt:      "Business growth chart and scalable architecture diagram",
			MetaTitle:     "Building Scalable Business Solutions | WebEnable",
			MetaDesc:      "Discover how to create business solutions that grow with your company and adapt to changing market needs.",
			ReadingTime:   6,
			IsFeatured:    false,
			ViewCount:     76,
			CreatedAt:     time.Date(2024, 1, 5, 9, 15, 0, 0, time.UTC),
			UpdatedAt:     time.Date(2024, 1, 5, 9, 15, 0, 0, time.UTC),
		},
		{
			ID:    "featured-1",
			Title: "Welcome to WebEnable Blog",
			Content: `# Welcome to WebEnable Blog

## About Our Blog

Welcome to the WebEnable blog, where we share insights, tips, and stories from our team of experts in web development, digital marketing, and business solutions.

## What You'll Find Here

Our blog covers a wide range of topics including:

- **Web Development** - Latest frameworks, tools, and best practices
- **Digital Marketing** - Strategies to grow your online presence
- **Business Solutions** - Scalable approaches to common business challenges
- **Industry Insights** - Trends and predictions for the digital world

## Our Mission

At WebEnable, we believe in building exceptional digital experiences that help businesses grow and succeed. Through our blog, we aim to share knowledge and insights that can help you achieve your digital goals.

## Stay Connected

Follow us for regular updates and insights:
- Subscribe to our newsletter
- Connect with us on social media
- Join our community discussions

Thank you for visiting our blog. We hope you find our content valuable and inspiring!`,
			Excerpt:       "Welcome to the WebEnable blog, where we share insights, tips, and stories from our team of digital experts.",
			Author:        "WebEnable Team",
			Status:        "published",
			Tags:          []string{"welcome", "blog", "introduction"},
			Categories:    []string{"General"},
			FeaturedImage: "https://images.unsplash.com/photo-1522202176988-66273c2fd55f?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2071&q=80",
			ImageAlt:      "Team collaboration and digital innovation",
			MetaTitle:     "Welcome to WebEnable Blog",
			MetaDesc:      "Welcome to the WebEnable blog, where we share insights, tips, and stories from our team of digital experts.",
			ReadingTime:   3,
			IsFeatured:    true,
			ViewCount:     250,
			CreatedAt:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	// Insert posts into database
	for _, post := range posts {
		// Check if post already exists
		exists := true
		row := db.Get(ctx, post.ID)
		var existingPost Post
		if err := row.ScanDoc(&existingPost); err != nil {
			exists = false
		}

		if exists {
			fmt.Printf("Post '%s' already exists, skipping...\n", post.Title)
			continue
		}

		// Set published date
		publishedAt := post.CreatedAt
		post.PublishedAt = &publishedAt

		// Insert post
		_, err := db.Put(ctx, post.ID, post)
		if err != nil {
			log.Printf("Failed to insert post '%s': %v", post.Title, err)
			continue
		}

		fmt.Printf("Successfully inserted post: %s\n", post.Title)
	}

	fmt.Println("Database population completed!")
}

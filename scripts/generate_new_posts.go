package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type BlogPost struct {
	Title           string   `json:"title"`
	Content         string   `json:"content"`
	Excerpt         string   `json:"excerpt"`
	Author          string   `json:"author"`
	Status          string   `json:"status"`
	Tags            []string `json:"tags"`
	Categories      []string `json:"categories"`
	FeaturedImage   string   `json:"featured_image,omitempty"`
	ImageAlt        string   `json:"image_alt,omitempty"`
	MetaTitle       string   `json:"meta_title,omitempty"`
	MetaDescription string   `json:"meta_description,omitempty"`
	IsFeatured      bool     `json:"is_featured"`
	ReadingTime     int      `json:"reading_time"`
	ViewCount       int      `json:"view_count"`
	PublishedAt     string   `json:"published_at,omitempty"`
}

func main() {
	// Sample posts to create
	posts := []BlogPost{
		{
			Title:           "The Future of Web Development in 2025",
			Content:         "# The Future of Web Development in 2025\n\nWeb development continues to evolve at a rapid pace. In 2025, we're seeing exciting trends that are reshaping how we build and interact with websites.\n\n## Key Trends\n\n### 1. AI-Powered Development Tools\nArtificial Intelligence is revolutionizing how developers write code. From intelligent code completion to automated testing, AI tools are becoming indispensable.\n\n### 2. Progressive Web Apps (PWAs)\nPWAs continue to bridge the gap between web and native applications, offering offline functionality and native-like experiences.\n\n### 3. WebAssembly Growth\nWebAssembly enables high-performance applications to run in browsers, opening new possibilities for complex web applications.\n\n### 4. Jamstack Architecture\nThe Jamstack approach provides better performance, security, and scalability for modern web applications.\n\n## Conclusion\n\nThe future of web development is bright, with new technologies enabling faster, more secure, and more user-friendly web experiences.",
			Excerpt:         "Explore the latest trends shaping web development in 2025, from AI-powered tools to progressive web apps.",
			Author:          "admin",
			Status:          "published",
			Tags:            []string{"web-development", "technology", "trends", "2025"},
			Categories:      []string{"Technology", "Development"},
			MetaTitle:       "Future of Web Development 2025 - Latest Trends & Technologies",
			MetaDescription: "Discover the key web development trends for 2025 including AI tools, PWAs, WebAssembly, and Jamstack architecture.",
			IsFeatured:      true,
			ReadingTime:     5,
			ViewCount:       0,
			PublishedAt:     time.Now().Format(time.RFC3339),
		},
		{
			Title:           "Complete Guide to Content Management Systems",
			Content:         "# Complete Guide to Content Management Systems\n\nA Content Management System (CMS) is essential for modern website management. This guide covers everything you need to know about choosing and using a CMS.\n\n## What is a CMS?\n\nA CMS allows you to create, manage, and modify content on your website without needing extensive technical knowledge.\n\n## Types of CMS\n\n### 1. Traditional CMS\n- WordPress\n- Drupal\n- Joomla\n\n### 2. Headless CMS\n- Strapi\n- Contentful\n- Ghost\n\n### 3. Static Site Generators\n- Gatsby\n- Next.js\n- Hugo\n\n## Choosing the Right CMS\n\nConsider these factors:\n- **Ease of use**\n- **Customization options**\n- **Performance**\n- **Security**\n- **Cost**\n\n## Best Practices\n\n1. Regular backups\n2. Keep software updated\n3. Use strong passwords\n4. Monitor performance\n5. Optimize for SEO\n\n## Conclusion\n\nThe right CMS can significantly improve your website management experience and boost your online presence.",
			Excerpt:         "A comprehensive guide to understanding and choosing the right Content Management System for your needs.",
			Author:          "admin",
			Status:          "published",
			Tags:            []string{"cms", "guide", "website", "management"},
			Categories:      []string{"Tutorials", "CMS"},
			MetaTitle:       "Complete CMS Guide - Choose the Right Content Management System",
			MetaDescription: "Learn about different types of CMS platforms and how to choose the best one for your website needs.",
			IsFeatured:      false,
			ReadingTime:     8,
			ViewCount:       0,
			PublishedAt:     time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
		{
			Title:           "Modern SEO Strategies That Actually Work",
			Content:         "# Modern SEO Strategies That Actually Work\n\nSearch Engine Optimization has evolved significantly. Here are the strategies that deliver real results in 2025.\n\n## Core Web Vitals\n\nGoogle prioritizes user experience metrics:\n- **Largest Contentful Paint (LCP)**\n- **First Input Delay (FID)**\n- **Cumulative Layout Shift (CLS)**\n\n## Content Quality\n\nFocus on:\n- **E-A-T (Expertise, Authoritativeness, Trustworthiness)**\n- Original research and insights\n- Comprehensive topic coverage\n- Regular content updates\n\n## Technical SEO\n\nEssential elements:\n- Mobile-first indexing\n- Page speed optimization\n- Schema markup\n- SSL certificates\n- XML sitemaps\n\n## Local SEO\n\nFor local businesses:\n- Google Business Profile optimization\n- Local citations\n- Customer reviews\n- Location-specific content\n\n## Link Building\n\nModern approaches:\n- Digital PR\n- Resource page link building\n- Broken link building\n- Guest posting (quality over quantity)\n\n## Conclusion\n\nSEO success requires a holistic approach combining technical optimization, quality content, and user experience improvements.",
			Excerpt:         "Discover the SEO strategies that actually work in 2025, from Core Web Vitals to modern link building techniques.",
			Author:          "admin",
			Status:          "published",
			Tags:            []string{"seo", "search-optimization", "google", "ranking"},
			Categories:      []string{"SEO", "Marketing"},
			MetaTitle:       "Modern SEO Strategies 2025 - What Actually Works",
			MetaDescription: "Learn effective SEO strategies for 2025 including Core Web Vitals, E-A-T, and modern link building techniques.",
			IsFeatured:      true,
			ReadingTime:     6,
			ViewCount:       0,
			PublishedAt:     time.Now().Add(-48 * time.Hour).Format(time.RFC3339),
		},
		{
			Title:           "Building Responsive Websites: Best Practices",
			Content:         "# Building Responsive Websites: Best Practices\n\nResponsive design is no longer optional. Learn how to create websites that work perfectly on all devices.\n\n## Mobile-First Design\n\nStart with mobile layouts:\n- Prioritize essential content\n- Use touch-friendly interactions\n- Optimize for thumb navigation\n- Consider one-handed use\n\n## Flexible Grid Systems\n\nCSS Grid and Flexbox:\n- Use relative units (%, em, rem)\n- Implement flexible containers\n- Create adaptive layouts\n- Maintain aspect ratios\n\n## Responsive Images\n\nOptimize images for all devices:\n- Use srcset attribute\n- Implement lazy loading\n- Choose appropriate formats\n- Compress for web\n\n## Typography\n\nReadable text across devices:\n- Scalable font sizes\n- Adequate line height\n- Proper contrast ratios\n- Limited font variations\n\n## Performance Considerations\n\nFast loading on all devices:\n- Minimize HTTP requests\n- Optimize CSS and JavaScript\n- Use CDNs\n- Enable compression\n\n## Testing\n\nEnsure quality across devices:\n- Real device testing\n- Browser developer tools\n- Responsive design mode\n- Performance testing\n\n## Conclusion\n\nResponsive design requires careful planning and attention to detail, but results in better user experiences across all devices.",
			Excerpt:         "Learn the essential best practices for building responsive websites that work perfectly on all devices.",
			Author:          "admin",
			Status:          "published",
			Tags:            []string{"responsive-design", "mobile", "css", "web-design"},
			Categories:      []string{"Design", "Development"},
			MetaTitle:       "Responsive Web Design Best Practices - Mobile-First Guide",
			MetaDescription: "Master responsive web design with mobile-first principles, flexible grids, and performance optimization techniques.",
			IsFeatured:      false,
			ReadingTime:     7,
			ViewCount:       0,
			PublishedAt:     time.Now().Add(-72 * time.Hour).Format(time.RFC3339),
		},
		{
			Title:           "JavaScript ES2025: New Features and Updates",
			Content:         "# JavaScript ES2025: New Features and Updates\n\nJavaScript continues to evolve with exciting new features. Here's what's new in ES2025.\n\n## Pattern Matching\n\nNew pattern matching syntax for cleaner conditional logic.\n\n## Temporal API\n\nBetter date and time handling with the new Temporal API.\n\n## Records and Tuples\n\nImmutable data structures coming to JavaScript.\n\n## Pipeline Operator\n\nFunctional composition made easier.\n\n## Import Assertions\n\nType-safe imports for JSON and CSS modules.\n\n## Decorators\n\nClass and method decorators for better code organization.\n\n## Conclusion\n\nES2025 brings powerful new features that make JavaScript more expressive and developer-friendly.",
			Excerpt:         "Explore the exciting new features coming to JavaScript in ES2025, from pattern matching to the Temporal API.",
			Author:          "admin",
			Status:          "published",
			Tags:            []string{"javascript", "es2025", "programming", "features"},
			Categories:      []string{"Programming", "JavaScript"},
			MetaTitle:       "JavaScript ES2025 New Features - Complete Guide",
			MetaDescription: "Discover the latest JavaScript ES2025 features including pattern matching, Temporal API, and pipeline operators.",
			IsFeatured:      false,
			ReadingTime:     4,
			ViewCount:       0,
			PublishedAt:     time.Now().Add(-96 * time.Hour).Format(time.RFC3339),
		},
		{
			Title:           "API Design: RESTful Best Practices",
			Content:         "# API Design: RESTful Best Practices\n\nWell-designed APIs are crucial for modern applications. Follow these best practices for creating robust RESTful APIs.\n\n## Resource Naming\n\nUse clear, consistent naming conventions for your API endpoints.\n\n## HTTP Methods\n\nUse appropriate HTTP verbs for different operations.\n\n## Status Codes\n\nReturn meaningful HTTP status codes for different scenarios.\n\n## Response Format\n\nMaintain consistent JSON structure across all endpoints.\n\n## Error Handling\n\nProvide detailed error information to help developers.\n\n## Authentication\n\nSecure your APIs with proper authentication mechanisms.\n\n## Documentation\n\nProvide comprehensive API documentation.\n\n## Conclusion\n\nFollowing these REST API best practices ensures your APIs are intuitive, maintainable, and secure.",
			Excerpt:         "Learn the essential best practices for designing robust and maintainable RESTful APIs.",
			Author:          "admin",
			Status:          "published",
			Tags:            []string{"api", "rest", "design", "backend"},
			Categories:      []string{"API", "Backend"},
			MetaTitle:       "RESTful API Design Best Practices - Complete Guide",
			MetaDescription: "Master RESTful API design with best practices for resource naming, HTTP methods, status codes, and security.",
			IsFeatured:      false,
			ReadingTime:     6,
			ViewCount:       0,
			PublishedAt:     time.Now().Add(-120 * time.Hour).Format(time.RFC3339),
		},
		{
			Title:           "Database Optimization Techniques",
			Content:         "# Database Optimization Techniques\n\nDatabase performance is critical for application success. Learn key optimization techniques to improve query performance.\n\n## Indexing Strategies\n\nProper indexing is crucial for database performance.\n\n## Query Optimization\n\nWrite efficient queries to minimize resource usage.\n\n## Database Design\n\nNormalize appropriately for your use case.\n\n## Connection Pooling\n\nManage database connections efficiently.\n\n## Caching Strategies\n\nReduce database load with effective caching.\n\n## Monitoring and Profiling\n\nTrack performance metrics continuously.\n\n## Partitioning\n\nScale large tables with partitioning strategies.\n\n## Conclusion\n\nDatabase optimization requires a systematic approach combining proper design, indexing, and monitoring.",
			Excerpt:         "Discover essential database optimization techniques to improve query performance and scalability.",
			Author:          "admin",
			Status:          "published",
			Tags:            []string{"database", "optimization", "performance", "sql"},
			Categories:      []string{"Database", "Performance"},
			MetaTitle:       "Database Optimization Techniques - Performance Guide",
			MetaDescription: "Learn database optimization techniques including indexing strategies, query optimization, and caching for better performance.",
			IsFeatured:      false,
			ReadingTime:     5,
			ViewCount:       0,
			PublishedAt:     time.Now().Add(-144 * time.Hour).Format(time.RFC3339),
		},
		{
			Title:           "Cybersecurity for Web Applications",
			Content:         "# Cybersecurity for Web Applications\n\nWeb application security is more important than ever. Learn how to protect your applications from common threats.\n\n## OWASP Top 10\n\nUnderstand the most critical web application security risks.\n\n## Input Validation\n\nSanitize all inputs to prevent injection attacks.\n\n## Authentication Security\n\nImplement robust authentication mechanisms.\n\n## HTTPS Implementation\n\nEncrypt data in transit with proper TLS configuration.\n\n## Data Protection\n\nSafeguard sensitive information with encryption and access controls.\n\n## Security Headers\n\nImplement protective HTTP headers.\n\n## Regular Security Testing\n\nContinuous security assessment and vulnerability management.\n\n## Incident Response\n\nPrepare for and respond to security incidents effectively.\n\n## Conclusion\n\nWeb application security requires a layered approach with continuous monitoring and regular updates.",
			Excerpt:         "Essential cybersecurity practices to protect web applications from common threats and vulnerabilities.",
			Author:          "admin",
			Status:          "published",
			Tags:            []string{"cybersecurity", "web-security", "owasp", "protection"},
			Categories:      []string{"Security", "Development"},
			MetaTitle:       "Web Application Cybersecurity - Essential Protection Guide",
			MetaDescription: "Learn essential cybersecurity practices for web applications including OWASP Top 10, authentication, and data protection.",
			IsFeatured:      true,
			ReadingTime:     7,
			ViewCount:       0,
			PublishedAt:     time.Now().Add(-168 * time.Hour).Format(time.RFC3339),
		},
		{
			Title:           "Cloud Computing: AWS vs Azure vs GCP",
			Content:         "# Cloud Computing: AWS vs Azure vs GCP\n\nChoosing the right cloud provider is crucial for your business. Compare the top three cloud platforms.\n\n## Amazon Web Services (AWS)\n\nThe market leader with extensive services and global reach.\n\n## Microsoft Azure\n\nStrong enterprise integration and hybrid cloud capabilities.\n\n## Google Cloud Platform (GCP)\n\nAdvanced AI/ML services and competitive pricing.\n\n## Key Comparison Factors\n\nPricing, services, global presence, and market position.\n\n## Decision Framework\n\nFactors to consider when choosing a cloud provider.\n\n## Multi-Cloud Strategy\n\nBenefits of using multiple cloud providers.\n\n## Conclusion\n\nEach cloud provider has unique strengths. Choose based on your specific requirements and long-term strategy.",
			Excerpt:         "Compare AWS, Azure, and Google Cloud Platform to choose the right cloud provider for your business needs.",
			Author:          "admin",
			Status:          "published",
			Tags:            []string{"cloud-computing", "aws", "azure", "gcp", "comparison"},
			Categories:      []string{"Cloud", "Technology"},
			MetaTitle:       "AWS vs Azure vs GCP - Cloud Platform Comparison 2025",
			MetaDescription: "Compare Amazon AWS, Microsoft Azure, and Google Cloud Platform features, pricing, and services to choose the right cloud provider.",
			IsFeatured:      false,
			ReadingTime:     8,
			ViewCount:       0,
			PublishedAt:     time.Now().Add(-192 * time.Hour).Format(time.RFC3339),
		},
		{
			Title:           "DevOps Best Practices for Modern Teams",
			Content:         "# DevOps Best Practices for Modern Teams\n\nDevOps transforms software development and deployment. Implement these best practices for successful DevOps adoption.\n\n## Culture and Collaboration\n\nFoster DevOps culture within your organization.\n\n## Continuous Integration (CI)\n\nAutomate code integration and testing.\n\n## Continuous Deployment (CD)\n\nStreamline deployment processes.\n\n## Infrastructure as Code (IaC)\n\nManage infrastructure programmatically.\n\n## Monitoring and Observability\n\nGain insights into system performance.\n\n## Security Integration (DevSecOps)\n\nBuild security into the development pipeline.\n\n## Tools and Technologies\n\nEssential DevOps tools and platforms.\n\n## Metrics and KPIs\n\nMeasure DevOps success with key metrics.\n\n## Common Challenges\n\nOvercome obstacles in DevOps adoption.\n\n## Getting Started\n\nSteps for successful DevOps implementation.\n\n## Conclusion\n\nSuccessful DevOps implementation requires cultural change, proper tools, and continuous improvement.",
			Excerpt:         "Learn essential DevOps best practices for modern software development teams and successful digital transformation.",
			Author:          "admin",
			Status:          "draft",
			Tags:            []string{"devops", "ci-cd", "automation", "deployment"},
			Categories:      []string{"DevOps", "Development"},
			MetaTitle:       "DevOps Best Practices - Modern Software Development Guide",
			MetaDescription: "Discover DevOps best practices including CI/CD, infrastructure as code, monitoring, and cultural transformation strategies.",
			IsFeatured:      false,
			ReadingTime:     6,
			ViewCount:       0,
		},
	}

	// API endpoint
	apiURL := "http://localhost/api/posts"

	fmt.Println("🚀 Starting to generate posts...")

	for i, post := range posts {
		fmt.Printf("📝 Creating post %d: %s\n", i+1, post.Title)

		// Convert post to JSON
		postJSON, err := json.Marshal(post)
		if err != nil {
			log.Printf("❌ Error marshaling post %d: %v", i+1, err)
			continue
		}

		// Create HTTP request
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(postJSON))
		if err != nil {
			log.Printf("❌ Error creating request for post %d: %v", i+1, err)
			continue
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzUyMTcyNDA2LCJpYXQiOjE3NTIwODYwMDZ9.XKOV7pjQQQ2xOoQLEHvcLRvEzV8B2bRYVYHEcveiLfU")

		// Make the request
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("❌ Error making request for post %d: %v", i+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			fmt.Printf("✅ Successfully created post %d: %s\n", i+1, post.Title)
		} else {
			fmt.Printf("❌ Failed to create post %d: %s (Status: %d)\n", i+1, post.Title, resp.StatusCode)
		}

		// Small delay between requests
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("🎉 Post generation completed!")
	fmt.Println("📊 Check your admin dashboard to see the new posts!")
}

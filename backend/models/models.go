package models

import "time"

type Post struct {
	ID            string     `json:"id,omitempty" db:"_id"`
	Rev           string     `json:"rev,omitempty" db:"_rev"`
	Title         string     `json:"title" validate:"required"`
	Content       string     `json:"content" validate:"required"`
	Excerpt       string     `json:"excerpt"`
	Author        string     `json:"author" validate:"required"`
	Status        string     `json:"status"` // draft, published, scheduled
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

type Category struct {
	ID          string    `json:"id,omitempty" db:"_id"`
	Rev         string    `json:"rev,omitempty" db:"_rev"`
	Name        string    `json:"name" validate:"required"`
	Slug        string    `json:"slug" validate:"required"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Icon        string    `json:"icon"`
	PostCount   int       `json:"post_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
	ID       string `json:"id,omitempty" db:"_id"`
	Rev      string `json:"rev,omitempty" db:"_rev"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
	Role     string `json:"role"` // admin, editor, author
}

type Contact struct {
	ID        string     `json:"id,omitempty" db:"_id"`
	Rev       string     `json:"rev,omitempty" db:"_rev"`
	Name      string     `json:"name" validate:"required"`
	Email     string     `json:"email" validate:"required,email"`
	Company   string     `json:"company"`
	Phone     string     `json:"phone"`
	Subject   string     `json:"subject" validate:"required"`
	Message   string     `json:"message" validate:"required"`
	Status    string     `json:"status"` // new, read, replied
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	RepliedAt *time.Time `json:"replied_at,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

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
	ID           string    `json:"id,omitempty"`
	Rev          string    `json:"_rev,omitempty"`
	Username     string    `json:"username" validate:"required,min=3,max=20"`
	Email        string    `json:"email" validate:"required,email"`
	PasswordHash string    `json:"password_hash,omitempty"`
	Role         string    `json:"role" validate:"required,oneof=admin editor author"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// UserStatsResponse represents user statistics
type UserStatsResponse struct {
	TotalUsers  int `json:"total_users"`
	AdminUsers  int `json:"admin_users"`
	EditorUsers int `json:"editor_users"`
	AuthorUsers int `json:"author_users"`
	ActiveUsers int `json:"active_users"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// PaginatedPostsResponse represents paginated posts response
type PaginatedPostsResponse struct {
	Data []Post         `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// PaginatedUsersResponse represents paginated users response
type PaginatedUsersResponse struct {
	Data []User         `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.PasswordHash = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

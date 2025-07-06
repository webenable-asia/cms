package email

import (
	"io"
)

// EmailAdapter defines the interface for email operations
type EmailAdapter interface {
	// Configuration
	Configure(config EmailConfig) error

	// Send Operations
	SendEmail(message EmailMessage) error
	SendReply(to, toName, subject, message, originalSubject, originalMessage string) error
	SendTemplatedEmail(template string, data interface{}, to []string) error

	// Health Check
	Health() error
	IsConfigured() bool
}

// EmailMessage represents an email message
type EmailMessage struct {
	To          []string          `json:"to"`
	CC          []string          `json:"cc"`
	BCC         []string          `json:"bcc"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	HTMLBody    string            `json:"html_body"`
	Attachments []EmailAttachment `json:"attachments"`
	Headers     map[string]string `json:"headers"`
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	Filename    string    `json:"filename"`
	Content     io.Reader `json:"-"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
}

// EmailConfig holds configuration for email adapters
type EmailConfig struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// EmailType constants for supported email types
const (
	EmailTypeSMTP     = "smtp"
	EmailTypeSendGrid = "sendgrid"
	EmailTypeSES      = "ses"
	EmailTypeMailgun  = "mailgun"
	EmailTypePostmark = "postmark"
)

// Common email errors
const (
	ErrInvalidRecipient = "invalid_recipient"
	ErrInvalidSender    = "invalid_sender"
	ErrSendFailed       = "send_failed"
	ErrNotConfigured    = "not_configured"
	ErrInvalidTemplate  = "invalid_template"
)
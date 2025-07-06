package email

import (
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

// SMTPAdapter implements EmailAdapter for SMTP
type SMTPAdapter struct {
	dialer *gomail.Dialer
	from   string
	config map[string]interface{}
}

// NewSMTPAdapter creates a new SMTP adapter
func NewSMTPAdapter(config map[string]interface{}) (EmailAdapter, error) {
	adapter := &SMTPAdapter{
		config: config,
	}

	if err := adapter.Configure(EmailConfig{
		Type:   EmailTypeSMTP,
		Config: config,
	}); err != nil {
		return nil, err
	}

	return adapter, nil
}

// Configure configures the SMTP adapter
func (s *SMTPAdapter) Configure(config EmailConfig) error {
	host, ok := config.Config["host"].(string)
	if !ok {
		host = "localhost"
	}

	portStr, ok := config.Config["port"].(string)
	if !ok {
		portStr = "1025"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	username, ok := config.Config["username"].(string)
	if !ok {
		username = "hello@webenable.asia"
	}

	password, _ := config.Config["password"].(string)

	from, ok := config.Config["from"].(string)
	if !ok {
		from = username
	}

	// Create dialer
	s.dialer = gomail.NewDialer(host, port, username, password)

	// For development without authentication
	if password == "" {
		s.dialer.Auth = nil
	}

	s.from = from
	return nil
}

// SendEmail sends a basic email message
func (s *SMTPAdapter) SendEmail(message EmailMessage) error {
	m := gomail.NewMessage()

	// Set headers
	m.SetHeader("From", s.from)
	m.SetHeader("To", message.To...)

	if len(message.CC) > 0 {
		m.SetHeader("Cc", message.CC...)
	}

	if len(message.BCC) > 0 {
		m.SetHeader("Bcc", message.BCC...)
	}

	m.SetHeader("Subject", message.Subject)

	// Set body
	if message.HTMLBody != "" {
		m.SetBody("text/html", message.HTMLBody)
		if message.Body != "" {
			m.AddAlternative("text/plain", message.Body)
		}
	} else {
		m.SetBody("text/plain", message.Body)
	}

	// Add custom headers
	for key, value := range message.Headers {
		m.SetHeader(key, value)
	}

	// Add attachments
	for range message.Attachments {
		// Note: gomail doesn't directly support io.Reader attachments
		// In a production system, you'd need to handle this differently
		// For now, we'll skip attachments or require file paths
	}

	// Send email
	return s.dialer.DialAndSend(m)
}

// SendReply sends an email reply using the existing reply template
func (s *SMTPAdapter) SendReply(to, toName, subject, message, originalSubject, originalMessage string) error {
	m := gomail.NewMessage()

	// Set headers
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("Re: %s", originalSubject))

	// Create HTML body using the existing template
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reply from WebEnable</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <div style="background: #2563eb; color: white; padding: 20px; text-align: center;">
            <h1 style="margin: 0;">WebEnable</h1>
            <p style="margin: 5px 0 0 0;">Digital Solutions & Development</p>
        </div>
        
        <div style="padding: 30px; background: #f9fafb; border-left: 4px solid #2563eb;">
            <h2 style="color: #2563eb; margin-top: 0;">Hi %s,</h2>
            <p>Thank you for contacting WebEnable. Here's our response to your inquiry:</p>
            
            <div style="background: white; padding: 20px; border-radius: 8px; margin: 20px 0;">
                %s
            </div>
            
            <p>If you have any additional questions, please don't hesitate to reach out to us.</p>
            
            <div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #e5e7eb;">
                <p style="margin: 0;"><strong>Best regards,</strong></p>
                <p style="margin: 5px 0;">WebEnable Team</p>
                <p style="margin: 0; color: #6b7280;">
                    Email: <a href="mailto:hello@webenable.asia" style="color: #2563eb;">hello@webenable.asia</a><br>
                    Website: <a href="https://webenable.asia" style="color: #2563eb;">webenable.asia</a>
                </p>
            </div>
        </div>
        
        <div style="background: #f3f4f6; padding: 20px; margin-top: 20px; border-radius: 8px;">
            <h3 style="margin-top: 0; color: #374151;">Your Original Message:</h3>
            <p style="margin: 0; color: #6b7280; font-style: italic;">"%s"</p>
        </div>
        
        <div style="text-align: center; margin-top: 30px; padding: 20px; color: #6b7280; font-size: 12px;">
            <p>© 2025 WebEnable. All rights reserved.</p>
            <p>Building amazing digital experiences for businesses worldwide.</p>
        </div>
    </div>
</body>
</html>`, toName, message, originalMessage)

	m.SetBody("text/html", htmlBody)

	// Set plain text alternative
	plainBody := fmt.Sprintf(`Hi %s,

Thank you for contacting WebEnable. Here's our response to your inquiry:

%s

If you have any additional questions, please don't hesitate to reach out to us.

Best regards,
WebEnable Team
Email: hello@webenable.asia
Website: webenable.asia

---
Your Original Message:
"%s"

© 2025 WebEnable. All rights reserved.
Building amazing digital experiences for businesses worldwide.`, toName, message, originalMessage)

	m.AddAlternative("text/plain", plainBody)

	// Send email
	return s.dialer.DialAndSend(m)
}

// SendTemplatedEmail sends an email using a template (basic implementation)
func (s *SMTPAdapter) SendTemplatedEmail(template string, data interface{}, to []string) error {
	// Basic template implementation
	// In a production system, you'd use a proper template engine
	return fmt.Errorf("templated email not implemented for SMTP adapter")
}

// Health checks the health of the SMTP connection
func (s *SMTPAdapter) Health() error {
	if s.dialer == nil {
		return fmt.Errorf("smtp adapter not configured")
	}
	return nil
}

// IsConfigured checks if the adapter is properly configured
func (s *SMTPAdapter) IsConfigured() bool {
	return s.dialer != nil
}
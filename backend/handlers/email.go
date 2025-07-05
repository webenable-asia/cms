package handlers

import (
	"fmt"
	"net/smtp"
	"os"
)

// EmailReply struct for email reply functionality
type EmailReply struct {
	Subject string `json:"subject" validate:"required"`
	Message string `json:"message" validate:"required"`
}

// SendEmailReply sends an email reply using SMTP
func SendEmailReply(toEmail, toName, subject, message string) error {
	// Email configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")

	// Default values for development
	if smtpHost == "" {
		smtpHost = "smtppro.zoho.com"
	}
	if smtpPort == "" {
		smtpPort = "587"
	}
	if fromEmail == "" {
		fromEmail = "hello@webenable.asia"
	}
	if fromName == "" {
		fromName = "WebEnable Team"
	}

	// Skip actual email sending in development if credentials not set
	if smtpUser == "" || smtpPass == "" {
		fmt.Printf("Email would be sent to %s (%s)\nSubject: %s\nMessage: %s\n", 
			toEmail, toName, subject, message)
		return nil // Simulate successful sending
	}

	// Create email content
	emailBody := fmt.Sprintf(`From: %s <%s>
To: %s <%s>
Subject: %s
MIME-Version: 1.0
Content-Type: text/html; charset=UTF-8

<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>%s</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <div style="background: #f8f9fa; padding: 20px; border-radius: 10px; margin-bottom: 20px;">
            <h2 style="color: #2563eb; margin: 0;">WebEnable</h2>
            <p style="margin: 5px 0 0 0; color: #666;">Digital Solutions & Development</p>
        </div>
        
        <div style="background: white; padding: 30px; border-radius: 10px; border: 1px solid #e5e7eb;">
            <h3 style="color: #374151; margin-top: 0;">Hi %s,</h3>
            
            <div style="margin: 20px 0;">
                %s
            </div>
            
            <div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #e5e7eb;">
                <p style="margin: 0; color: #666;">Best regards,<br>
                <strong>%s</strong><br>
                WebEnable Team</p>
                
                <div style="margin-top: 20px; font-size: 14px; color: #888;">
                    <p>WebEnable - Building Amazing Digital Experiences</p>
                    <p>Email: hello@webenable.asia | Website: https://webenable.asia</p>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`, fromName, fromEmail, toName, toEmail, subject, subject, toName, 
		FormatMessageForHTML(message), fromName)

	// SMTP authentication
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send email
	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		fromEmail,
		[]string{toEmail},
		[]byte(emailBody),
	)

	return err
}

// FormatMessageForHTML converts plain text to HTML format
func FormatMessageForHTML(message string) string {
	// Simple HTML formatting for line breaks
	// In a production environment, you might want more sophisticated formatting
	formatted := ""
	for _, line := range []rune(message) {
		if line == '\n' {
			formatted += "<br>"
		} else {
			formatted += string(line)
		}
	}
	return formatted
}

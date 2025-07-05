package services

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer *gomail.Dialer
	from   string
}

var Email *EmailService

func InitEmailService() {
	// Email configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	// Default values for development
	if smtpHost == "" {
		smtpHost = "localhost"
	}
	if smtpPortStr == "" {
		smtpPortStr = "1025" // MailHog default port
	}
	if smtpUser == "" {
		smtpUser = "hello@webenable.asia"
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		smtpPort = 1025
	}

	// Create dialer
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// For development without authentication
	if smtpPass == "" {
		dialer.Auth = nil
	}

	Email = &EmailService{
		dialer: dialer,
		from:   smtpUser,
	}
}

func (e *EmailService) SendReply(to, toName, subject, message, originalSubject, originalMessage string) error {
	m := gomail.NewMessage()
	
	// Set headers
	m.SetHeader("From", e.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("Re: %s", originalSubject))
	
	// Create HTML body
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
	return e.dialer.DialAndSend(m)
}

func (e *EmailService) IsConfigured() bool {
	return e != nil && e.dialer != nil
}

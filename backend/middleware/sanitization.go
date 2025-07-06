package middleware

import (
	"html"
	"net/http"
	"regexp"
	"strings"
)

// XSSProtection middleware sanitizes request data to prevent XSS attacks
func XSSProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sanitize query parameters
		for key, values := range r.URL.Query() {
			for i, value := range values {
				r.URL.Query()[key][i] = SanitizeInput(value)
			}
		}

		// Sanitize form values if present
		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			r.ParseForm()
			for key, values := range r.Form {
				for i, value := range values {
					r.Form[key][i] = SanitizeInput(value)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// SanitizeInput removes potentially dangerous HTML/JS content
func SanitizeInput(input string) string {
	// HTML escape the input
	sanitized := html.EscapeString(input)

	// Remove script tags and their content
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	sanitized = scriptRegex.ReplaceAllString(sanitized, "")

	// Remove javascript: protocol
	jsRegex := regexp.MustCompile(`(?i)javascript:`)
	sanitized = jsRegex.ReplaceAllString(sanitized, "")

	// Remove on* event handlers
	eventRegex := regexp.MustCompile(`(?i)on\w+\s*=`)
	sanitized = eventRegex.ReplaceAllString(sanitized, "")

	// Remove data: protocol (can be used for XSS)
	dataRegex := regexp.MustCompile(`(?i)data:`)
	sanitized = dataRegex.ReplaceAllString(sanitized, "")

	return strings.TrimSpace(sanitized)
}

// SanitizeHTML provides more lenient HTML sanitization for content fields
func SanitizeHTML(input string) string {
	// Remove script tags and their content
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	sanitized := scriptRegex.ReplaceAllString(input, "")

	// Remove javascript: protocol
	jsRegex := regexp.MustCompile(`(?i)javascript:`)
	sanitized = jsRegex.ReplaceAllString(sanitized, "")

	// Remove on* event handlers
	eventRegex := regexp.MustCompile(`(?i)on\w+\s*=\s*["'][^"']*["']`)
	sanitized = eventRegex.ReplaceAllString(sanitized, "")

	// Remove style attributes (can contain CSS-based XSS)
	styleRegex := regexp.MustCompile(`(?i)style\s*=\s*["'][^"']*["']`)
	sanitized = styleRegex.ReplaceAllString(sanitized, "")

	return strings.TrimSpace(sanitized)
}

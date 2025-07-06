package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})

	// Wrap with security headers middleware
	handler := SecurityHeaders(testHandler)

	// Create request
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check status
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check security headers
	headers := rr.Header()

	assert.Equal(t, "max-age=31536000; includeSubDomains; preload", headers.Get("Strict-Transport-Security"))
	assert.Contains(t, headers.Get("Content-Security-Policy"), "default-src 'self'")
	assert.Equal(t, "DENY", headers.Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", headers.Get("X-Content-Type-Options"))
	assert.Equal(t, "1; mode=block", headers.Get("X-XSS-Protection"))
	assert.Equal(t, "strict-origin-when-cross-origin", headers.Get("Referrer-Policy"))
	assert.Contains(t, headers.Get("Permissions-Policy"), "geolocation=()")
	assert.Equal(t, "", headers.Get("Server"))
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "HTML tags",
			input:    "<p>Hello</p>",
			expected: "&lt;p&gt;Hello&lt;/p&gt;",
		},
		{
			name:     "Script tag",
			input:    "<script>alert('xss')</script>",
			expected: "",
		},
		{
			name:     "JavaScript protocol",
			input:    "javascript:alert('xss')",
			expected: "alert(&#39;xss&#39;)",
		},
		{
			name:     "Event handler",
			input:    "onclick=alert('xss')",
			expected: "alert(&#39;xss&#39;)",
		},
		{
			name:     "Data protocol",
			input:    "data:text/html,<script>alert('xss')</script>",
			expected: "text/html,",
		},
		{
			name:     "Mixed content",
			input:    "<p onclick='alert()'>Hello</p><script>bad()</script>",
			expected: "&lt;p &#39;alert()&#39;&gt;Hello&lt;/p&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Script tag removal",
			input:    "<p>Hello</p><script>alert('xss')</script>",
			expected: "<p>Hello</p>",
		},
		{
			name:     "Event handler removal",
			input:    "<p onclick='alert()'>Hello</p>",
			expected: "<p >Hello</p>",
		},
		{
			name:     "Style attribute removal",
			input:    "<p style='color: red;'>Hello</p>",
			expected: "<p >Hello</p>",
		},
		{
			name:     "JavaScript protocol removal",
			input:    "<a href='javascript:alert()'>Link</a>",
			expected: "<a href='alert()'>Link</a>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeHTML(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

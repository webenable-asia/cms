package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"webenable-cms-backend/models"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    models.LoginRequest
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "Valid login",
			requestBody: models.LoginRequest{
				Username: "admin",
				Password: "/juk+vfdbNk6TICg",
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "Invalid username",
			requestBody: models.LoginRequest{
				Username: "nonexistent",
				Password: "password",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name: "Invalid password",
			requestBody: models.LoginRequest{
				Username: "admin",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name: "Empty username",
			requestBody: models.LoginRequest{
				Username: "",
				Password: "password",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name: "Empty password",
			requestBody: models.LoginRequest{
				Username: "admin",
				Password: "",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// Create request
			req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			Login(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectToken {
				// Parse response
				var response models.LoginResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Token)
				assert.Equal(t, tt.requestBody.Username, response.User.Username)
			}
		})
	}
}

func TestLoginInvalidJSON(t *testing.T) {
	// Test with invalid JSON
	req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString("invalid json"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	Login(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLogout(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/auth/logout", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	Logout(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Logged out successfully", response["message"])
}

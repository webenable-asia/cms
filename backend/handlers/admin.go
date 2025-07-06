package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"webenable-cms-backend/middleware"
	"webenable-cms-backend/models"
)

// ResetRateLimit godoc
//
//	@Summary		Reset rate limit
//	@Description	Reset rate limit for specific identifier or all limits (admin only)
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			type	query		string	false	"Type of reset: 'ip', 'user', 'all', 'api', 'auth'"
//	@Param			target	query		string	false	"Target IP address or user ID (required for ip/user types)"
//	@Success		200		{object}	models.SuccessResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		403		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/admin/rate-limit/reset [post]
func ResetRateLimit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user claims from JWT middleware
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Only admin can reset rate limits
	if claims.Role != "admin" {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	if globalRateLimiter == nil {
		http.Error(w, "Rate limiter not available", http.StatusInternalServerError)
		return
	}

	resetType := r.URL.Query().Get("type")
	target := r.URL.Query().Get("target")

	var err error
	var message string

	switch strings.ToLower(resetType) {
	case "ip":
		if target == "" {
			http.Error(w, "Target IP address is required", http.StatusBadRequest)
			return
		}
		err = globalRateLimiter.ResetRateLimitForIP(target)
		message = "Rate limit reset for IP: " + target

	case "user":
		if target == "" {
			http.Error(w, "Target user ID is required", http.StatusBadRequest)
			return
		}
		err = globalRateLimiter.ResetRateLimitForUser(target)
		message = "Rate limit reset for user: " + target

	case "api":
		err = globalRateLimiter.ResetAllAPIRateLimits()
		message = "All API rate limits reset"

	case "auth":
		err = globalRateLimiter.ResetAllAuthRateLimits()
		message = "All authentication rate limits reset"

	case "users":
		err = globalRateLimiter.ResetAllUserRateLimits()
		message = "All user rate limits reset"

	case "all":
		err = globalRateLimiter.ResetAllRateLimits()
		message = "All rate limits reset"

	default:
		http.Error(w, "Invalid reset type. Use: ip, user, api, auth, users, or all", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to reset rate limit: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.SuccessResponse{Message: message}
	json.NewEncoder(w).Encode(response)
}

// GetRateLimitStatus godoc
//
//	@Summary		Get rate limit status
//	@Description	Get current rate limit status for an IP or user (admin only)
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			type	query		string	true	"Type: 'ip' or 'user'"
//	@Param			target	query		string	true	"Target IP address or user ID"
//	@Success		200		{object}	object{api_limit=object{current=int,remaining=int,reset_time=string},auth_limit=object{current=int,remaining=int,reset_time=string},user_limit=object{current=int,remaining=int,reset_time=string}}
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		403		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/admin/rate-limit/status [get]
func GetRateLimitStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user claims from JWT middleware
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Only admin can check rate limit status
	if claims.Role != "admin" {
		http.Error(w, "Insufficient permissions", http.StatusForbidden)
		return
	}

	if globalRateLimiter == nil {
		http.Error(w, "Rate limiter not available", http.StatusInternalServerError)
		return
	}

	targetType := r.URL.Query().Get("type")
	target := r.URL.Query().Get("target")

	if target == "" {
		http.Error(w, "Target is required", http.StatusBadRequest)
		return
	}

	response := make(map[string]interface{})

	switch strings.ToLower(targetType) {
	case "ip":
		// Check API rate limit
		apiCurrent, apiRemaining, apiResetTime, err := globalRateLimiter.GetRateLimitStatus("rate_limit:api:"+target, 100)
		if err == nil {
			response["api_limit"] = map[string]interface{}{
				"current":    apiCurrent,
				"remaining":  apiRemaining,
				"reset_time": apiResetTime.Format("2006-01-02T15:04:05Z"),
			}
		}

		// Check auth rate limit
		authCurrent, authRemaining, authResetTime, err := globalRateLimiter.GetRateLimitStatus("rate_limit:auth:"+target, 10)
		if err == nil {
			response["auth_limit"] = map[string]interface{}{
				"current":    authCurrent,
				"remaining":  authRemaining,
				"reset_time": authResetTime.Format("2006-01-02T15:04:05Z"),
			}
		}

	case "user":
		// Check user rate limit
		userCurrent, userRemaining, userResetTime, err := globalRateLimiter.GetRateLimitStatus("rate_limit:user:"+target, 120)
		if err == nil {
			response["user_limit"] = map[string]interface{}{
				"current":    userCurrent,
				"remaining":  userRemaining,
				"reset_time": userResetTime.Format("2006-01-02T15:04:05Z"),
			}
		}

	default:
		http.Error(w, "Invalid type. Use: ip or user", http.StatusBadRequest)
		return
	}

	if len(response) == 0 {
		response["message"] = "No rate limit data found for target"
	}

	json.NewEncoder(w).Encode(response)
}

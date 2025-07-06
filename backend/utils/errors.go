package utils

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// ErrorCode represents standard error codes
type ErrorCode string

const (
	ErrCodeBadRequest     ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeConflict       ErrorCode = "CONFLICT"
	ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeRateLimit      ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrCodeServiceUnavail ErrorCode = "SERVICE_UNAVAILABLE"
)

// StandardError represents a standardized error response
type StandardError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// WriteErrorResponse writes a standardized error response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, code ErrorCode, message string, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := StandardError{
		Code:    code,
		Message: message,
		Details: details,
	}

	json.NewEncoder(w).Encode(errorResponse)
}

// WriteErrorResponseWithLog writes an error response and logs it
func WriteErrorResponseWithLog(w http.ResponseWriter, statusCode int, code ErrorCode, message string, err error, fields logrus.Fields) {
	details := ""
	if err != nil {
		details = err.Error()
		LogError(err, message, fields)
	} else {
		LogWarning(message, fields)
	}

	WriteErrorResponse(w, statusCode, code, message, details)
}

// Common error response helpers
func BadRequest(w http.ResponseWriter, message string, details string) {
	WriteErrorResponse(w, http.StatusBadRequest, ErrCodeBadRequest, message, details)
}

func Unauthorized(w http.ResponseWriter, message string) {
	WriteErrorResponse(w, http.StatusUnauthorized, ErrCodeUnauthorized, message, "")
}

func Forbidden(w http.ResponseWriter, message string) {
	WriteErrorResponse(w, http.StatusForbidden, ErrCodeForbidden, message, "")
}

func NotFound(w http.ResponseWriter, message string) {
	WriteErrorResponse(w, http.StatusNotFound, ErrCodeNotFound, message, "")
}

func Conflict(w http.ResponseWriter, message string, details string) {
	WriteErrorResponse(w, http.StatusConflict, ErrCodeConflict, message, details)
}

func ValidationError(w http.ResponseWriter, message string, details string) {
	WriteErrorResponse(w, http.StatusBadRequest, ErrCodeValidation, message, details)
}

func InternalError(w http.ResponseWriter, message string, err error, fields logrus.Fields) {
	WriteErrorResponseWithLog(w, http.StatusInternalServerError, ErrCodeInternal, message, err, fields)
}

func RateLimit(w http.ResponseWriter, message string) {
	WriteErrorResponse(w, http.StatusTooManyRequests, ErrCodeRateLimit, message, "")
}

func ServiceUnavailable(w http.ResponseWriter, message string, details string) {
	WriteErrorResponse(w, http.StatusServiceUnavailable, ErrCodeServiceUnavail, message, details)
}

// WriteSuccessResponse writes a standardized success response
func WriteSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// WriteCreatedResponse writes a standardized created response
func WriteCreatedResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

// WriteNoContentResponse writes a standardized no content response
func WriteNoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

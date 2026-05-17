package utils

import (
	"net/http"
)

// APIError represents a standardized API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewAPIError creates a new API error
func NewAPIError(code int, message, details string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// ErrorBadRequest creates a 400 Bad Request error
func ErrorBadRequest(details string) *APIError {
	return NewAPIError(http.StatusBadRequest, "bad request", details)
}

// ErrorNotFound creates a 404 Not Found error
func ErrorNotFound(details string) *APIError {
	return NewAPIError(http.StatusNotFound, "not found", details)
}

// ErrorInternalServer creates a 500 Internal Server Error
func ErrorInternalServer(details string) *APIError {
	return NewAPIError(http.StatusInternalServerError, "internal server error", details)
}

// ErrorUnprocessableEntity creates a 422 Unprocessable Entity error
func ErrorUnprocessableEntity(details string) *APIError {
	return NewAPIError(http.StatusUnprocessableEntity, "validation error", details)
}

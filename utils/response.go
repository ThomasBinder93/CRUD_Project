package utils

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse wraps a successful response
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// ErrorResponse wraps an error response
type ErrorResponse struct {
	Error *APIError `json:"error"`
}

// WriteJSON writes a JSON response with proper headers
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// WriteError writes a standardized error response
func WriteError(w http.ResponseWriter, apiErr *APIError) error {
	return WriteJSON(w, apiErr.Code, ErrorResponse{Error: apiErr})
}

// WriteSuccess writes a successful response with data
func WriteSuccess(w http.ResponseWriter, statusCode int, data interface{}) error {
	return WriteJSON(w, statusCode, SuccessResponse{Data: data})
}

// WriteNoContent writes a 204 No Content response
func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

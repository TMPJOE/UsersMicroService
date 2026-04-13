// Package helper provides common utility functions used across the application.
// This includes helper methods for data transformation, validation, formatting,
// and other reusable utility operations.
package helper

import (
	"encoding/json"
	"net/http"
	"os"
)

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// ErrorResponse is the standard JSON shape for all error responses.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RespondError writes a consistent JSON error response from any handler.
func RespondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: ErrorDetail{
			Code:    http.StatusText(status),
			Message: message,
		},
	})
}

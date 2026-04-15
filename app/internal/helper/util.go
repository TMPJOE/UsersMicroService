// Package helper provides common utility functions used across the application.
// This includes helper methods for data transformation, validation, formatting,
// and other reusable utility operations.
package helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
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

var (
	// General errors
	ErrDefault        = errors.New("")
	ErrInternalServer = errors.New("internal server error")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrNotFound       = errors.New("not found")
	ErrBadRequest     = errors.New("bad request")
	ErrConflict       = errors.New("resource conflict")
	ErrValidation     = errors.New("validation failed")

	// Repository/Database errors
	ErrDBConnection    = errors.New("database connection failed")
	ErrDBQuery         = errors.New("database query failed")
	ErrDBTransaction   = errors.New("database transaction failed")
	ErrRecordNotFound  = errors.New("record not found")
	ErrDuplicateEntry  = errors.New("duplicate entry")
	ErrForeignKey      = errors.New("foreign key constraint violation")
	ErrDBTimeout       = errors.New("database timeout")
	ErrDBSerialization = errors.New("database serialization error")
	ErrMigrationFailed = errors.New("database migration failed")

	// Service layer errors
	ErrServiceUnavailable  = errors.New("service unavailable")
	ErrCreateFailed        = errors.New("failed to create resource")
	ErrUpdateFailed        = errors.New("failed to update resource")
	ErrDeleteFailed        = errors.New("failed to delete resource")
	ErrFetchFailed         = errors.New("failed to fetch resource")
	ErrProcessingFailed    = errors.New("failed to process request")
	ErrDependencyFailed    = errors.New("dependency failed")
	ErrExternalServiceCall = errors.New("external service call failed")
	ErrCacheMiss           = errors.New("cache miss")
	ErrCacheWrite          = errors.New("cache write failed")

	// Authentication/Authorization errors
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token expired")
	ErrTokenGeneration      = errors.New("token generation failed")
	ErrUUIDGeneration       = errors.New("uuid generation failed")
	ErrPasswordHash         = errors.New("password hashing failed")
	ErrPasswordVerification = errors.New("password verification failed")
	ErrPermissionDenied     = errors.New("permission denied")
	ErrRoleNotFound         = errors.New("role not found")

	// Validation errors
	ErrInvalidInput    = errors.New("invalid input")
	ErrMissingField    = errors.New("required field missing")
	ErrInvalidFormat   = errors.New("invalid format")
	ErrOutOfRange      = errors.New("value out of range")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidPassword = errors.New("invalid password")
	ErrResourceExists  = errors.New("resource already exists")
)

// MapError maps a raw error to the best matching sentinel error.
// It checks for key string patterns that are common across different DB engines.
// Use this in the repo layer to wrap raw driver errors into application sentinels.
func MapError(err error) error {
	if err == nil {
		return nil
	}
	msg := strings.ToLower(err.Error())

	// ── DB-specific patterns (PostgreSQL, MySQL, SQLite, SQL Server) ──
	switch {
	// Connection errors
	case strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "connection failed") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "dial tcp") ||
		strings.Contains(msg, "connect: connection refused"):
		return ErrDBConnection

	// Timeout errors
	case strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "context deadline exceeded") ||
		strings.Contains(msg, "timed out") ||
		strings.Contains(msg, "i/o timeout"):
		return ErrDBTimeout

	// Duplicate / unique constraint violations
	case strings.Contains(msg, "duplicate") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "unique_violation") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "constraint failed") ||
		strings.Contains(msg, "unique"):
		return ErrDuplicateEntry

	// Foreign key violations
	case strings.Contains(msg, "foreign key") ||
		strings.Contains(msg, "foreign_key") ||
		strings.Contains(msg, "references") ||
		strings.Contains(msg, "constraint violation"):
		return ErrForeignKey

	// Not found — note: sql.ErrNoRows should be checked with errors.Is before calling this
	case strings.Contains(msg, "not found") ||
		strings.Contains(msg, "no rows") ||
		strings.Contains(msg, "no record") ||
		strings.Contains(msg, "does not exist") ||
		strings.Contains(msg, "unknown"):
		return ErrRecordNotFound

	// Serialization / deadlock
	case strings.Contains(msg, "deadlock") ||
		strings.Contains(msg, "serialization") ||
		strings.Contains(msg, "lock wait timeout") ||
		strings.Contains(msg, "try restarting transaction"):
		return ErrDBSerialization

	// Transaction errors
	case strings.Contains(msg, "transaction") && strings.Contains(msg, "commit") ||
		strings.Contains(msg, "transaction") && strings.Contains(msg, "rollback") ||
		strings.Contains(msg, "current transaction is aborted"):
		return ErrDBTransaction

	// Migration errors
	case strings.Contains(msg, "migration") && (strings.Contains(msg, "failed") ||
		strings.Contains(msg, "error") ||
		strings.Contains(msg, "not applied")):
		return ErrMigrationFailed

	// Default — any other DB-related error
	default:
		return ErrDBQuery
	}
}

// IsDuplicate checks if the error chain contains ErrDuplicateEntry.
func IsDuplicate(err error) bool {
	return errors.Is(err, ErrDuplicateEntry)
}

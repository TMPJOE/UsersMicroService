// Package validator provides request validation utilities using go-playground/validator.
package helper

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// NewValidator creates a configured validator instance with custom error messages.
func NewValidator() *RequestValidator {
	v := validator.New()
	return &RequestValidator{validate: v}
}

// RequestValidator wraps the validator instance and provides a clean Validate method.
type RequestValidator struct {
	validate *validator.Validate
}

// Validate checks a struct against its validation tags and returns a cleaned error message.
// Returns nil if validation passes.
func (rv *RequestValidator) Validate(s any) error {
	err := rv.validate.Struct(s)
	if err == nil {
		return nil
	}

	// Collect field-level errors into a single message
	var msgs []string
	for _, e := range err.(validator.ValidationErrors) {
		msgs = append(msgs, e.Field()+": "+e.Tag())
	}
	return &ValidationError{fields: msgs}
}

// ValidationError is a lightweight error that joins all field validation failures.
type ValidationError struct {
	fields []string
}

func (ve *ValidationError) Error() string {
	return strings.Join(ve.fields, "; ")
}

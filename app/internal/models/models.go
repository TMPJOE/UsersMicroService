// Package models defines domain data structures used across the application.
// All entities, DTOs, and shared types should be defined here
// to ensure consistency between repository, service, and handler layers.
package models

import "time"

type User struct {
	ID          string    `json:"id" db:"id"`
	Email       string    `json:"email" db:"email"`
	DisplayName string    `json:"display_name" db:"display_name"`
	UserType    string    `json:"user_type" db:"user_type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}

type UserCreation struct {
	Email       string `json:"email" db:"email"`
	DisplayName string `json:"display_name" db:"display_name"`
	Password    string `json:"password"`
}

type UserIO struct {
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserUpdate struct {
	DisplayName string `json:"display_name"`
}

type Login struct {
	ID             string    `json:"id" db:"id"`
	PswHash        string    `json:"password_hash" db:"password_hash"`
	LastLogin      time.Time `json:"last_login" db:"last_login"`
	FailedAttempts int       `json:"failed_attempts" db:"failed_attempts"`
	IsLocked       bool      `json:"is_locked" db:"is_locked"`
	UserID         string    `json:"user_id" db:"user_id"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Package repo implements the data access layer of the application.
// It handles all database queries, transactions, and data mapping,
// providing a clean interface for the service layer to interact with PostgreSQL.
package repo

import "hotel.com/app/internal/models"

type ServiceRepository interface {
	DbPing() error
	// CRUD ops
	Create(models.User, models.Login) error
	GetUserByEmail(email string) (*models.User, *models.Login, error)
	GetUserByID(id string) (*models.UserIO, error)
	UpdateLastLogin(userID string) error
	Update(displayName, userID string) error
	// Delete() error
}

//REMEMBER TRANSACTION CODE LOGIC

// Package repo implements the data access layer of the application.
// It handles all database queries, transactions, and data mapping,
// providing a clean interface for the service layer to interact with PostgreSQL.
package repo

import (
	"context"
)

type ServiceRepository interface {
	Foo(ctx context.Context) error
	DbPing() error
}

//REMEMBER TRANSACTION CODE LOGIC

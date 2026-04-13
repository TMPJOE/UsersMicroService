package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type databaseRepo struct {
	db *pgxpool.Pool
}

func NewDatabaseRepo(conn *pgxpool.Pool) ServiceRepository {
	return &databaseRepo{
		db: conn,
	}
}

func (dbr *databaseRepo) DbPing() error {
	err := dbr.db.Ping(context.Background())
	return err
}

func (dbr *databaseRepo) Foo(ctx context.Context) error {
	return nil
}

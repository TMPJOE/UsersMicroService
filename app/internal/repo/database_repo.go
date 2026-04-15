package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"hotel.com/app/internal/helper"
	"hotel.com/app/internal/models"
)

type postgreRepo struct {
	db *pgxpool.Pool
}

func NewPostgreRepo(conn *pgxpool.Pool) ServiceRepository {
	return &postgreRepo{
		db: conn,
	}
}

func (r *postgreRepo) DbPing() error {
	err := r.db.Ping(context.Background())
	return err
}

func (r *postgreRepo) Create(models.User, models.Login) error {
	stmt := "INSERT INTO users (id,email, display_name, user_type) VALUES($1,$2,$3,$4)"
	_, err := r.db.Exec(context.Background(), stmt)
	if err != nil {
		return fmt.Errorf("create user: %w", helper.MapError(err))
	}
	return nil
}

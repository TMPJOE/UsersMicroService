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

func (r *postgreRepo) Create(user models.User, login models.Login) error {
	stmt := "INSERT INTO users (id,email, display_name, user_type) VALUES($1,$2,$3,$4)"
	_, err := r.db.Exec(context.Background(), stmt)
	if err != nil {
		return fmt.Errorf("create user: %w", helper.MapError(err))
	}
	return nil
}

func (r *postgreRepo) GetUserByEmail(email string) (*models.User, *models.Login, error) {
	// Get user by email
	userStmt := "SELECT id, email, display_name, user_type, created_at, updated_at, is_active FROM users WHERE email = $1"
	var user models.User
	err := r.db.QueryRow(context.Background(), userStmt, email).Scan(
		&user.ID, &user.Email, &user.DisplayName, &user.UserType,
		&user.CreatedAt, &user.UpdatedAt, &user.IsActive,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("get user by email: %w", helper.MapError(err))
	}

	// Get login info
	loginStmt := "SELECT id, password_hash, last_login, failed_attempts, is_locked, user_id FROM logins WHERE user_id = $1"
	var login models.Login
	err = r.db.QueryRow(context.Background(), loginStmt, user.ID).Scan(
		&login.ID, &login.PswHash, &login.LastLogin, &login.FailedAttempts,
		&login.IsLocked, &login.UserID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("get login info: %w", helper.MapError(err))
	}

	return &user, &login, nil
}

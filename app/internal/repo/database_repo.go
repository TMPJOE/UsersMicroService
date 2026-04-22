package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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
	argsUsers := pgx.NamedArgs{
		"id":           user.ID,
		"email":        user.Email,
		"display_name": user.DisplayName,
		"user_type":    user.UserType,
	}
	stmtUsers := "INSERT INTO users (id,email, display_name, user_type) VALUES(@id, @email, @display_name, @user_type)"

	argsLogin := pgx.NamedArgs{
		"id_psw":        login.ID,
		"password_hash": login.PswHash,
		"user_id":       login.UserID,
	}
	stmtLogin := "INSERT INTO logins (id, password_hash, user_id) VALUES(@id_psw, @password_hash, @user_id)"

	//Transaction to ensure both user and login are created together
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("begin transaction: %w", helper.MapError(err))
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), stmtUsers, argsUsers)
	if err != nil {
		return fmt.Errorf("create user: %w", helper.MapError(err))
	}
	_, err = tx.Exec(context.Background(), stmtLogin, argsLogin)
	if err != nil {
		return fmt.Errorf("create login: %w", helper.MapError(err))
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("commit transaction: %w", helper.MapError(err))
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
	loginStmt := "SELECT id, password_hash, failed_attempts, is_locked, user_id FROM logins WHERE user_id = $1"
	var login models.Login

	err = r.db.QueryRow(context.Background(), loginStmt, user.ID).Scan(
		&login.ID, &login.PswHash, &login.FailedAttempts,
		&login.IsLocked, &login.UserID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("get login info: %w", helper.MapError(err))
	}

	return &user, &login, nil
}

func (r *postgreRepo) GetUserByID(id string) (*models.UserIO, error) {
	// Get user by ID
	userStmt := "SELECT email, display_name, updated_at, is_active FROM users WHERE id = $1"
	var user models.UserIO
	var isActive bool
	err := r.db.QueryRow(context.Background(), userStmt, id).Scan(
		&user.Email, &user.DisplayName, &user.UpdatedAt, &isActive,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by ID: %w", helper.MapError(err))
	}

	if !isActive {
		return nil, fmt.Errorf("get user by ID: %w", helper.ErrRecordNotFound)
	}

	return &user, nil
}

func (r *postgreRepo) UpdateLastLogin(userID string) error {
	stmt := "UPDATE logins SET last_login = NOW() WHERE user_id = $1"
	_, err := r.db.Exec(context.Background(), stmt, userID)
	if err != nil {
		return fmt.Errorf("update last login: %w", helper.MapError(err))
	}
	return nil
}

func (r *postgreRepo) Update(displayName, userID string) error {
	stmt := "UPDATE users SET display_name = $1, updated_at = NOW() WHERE id = $2"
	_, err := r.db.Exec(context.Background(), stmt, displayName)
	if err != nil {
		return fmt.Errorf("update user: %w", helper.MapError(err))
	}
	return nil
}

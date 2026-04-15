// Package service contains the business logic layer of the application.
// It defines service interfaces and implements use cases by orchestrating
// repositories, applying business rules, and returning results to handlers.
package service

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"hotel.com/app/internal/helper"
	"hotel.com/app/internal/models"
	"hotel.com/app/internal/repo"
)

type Service interface {
	Check() error
	CreateUser(models.UserCreation) error
}

type UserService struct {
	l *slog.Logger
	r repo.ServiceRepository
}

func New(l *slog.Logger, r repo.ServiceRepository) Service {
	return &UserService{
		l: l,
		r: r,
	}
}

func (s *UserService) Check() error {
	s.l.Info("Pinging db...")
	err := s.r.DbPing()
	s.l.Info("is service working", "err", err.Error())
	return err
}

func (s *UserService) CreateUser(usr models.UserCreation) error {

	// create uuid and hash psw
	userID, err := uuid.NewV7()
	if err != nil {
		s.l.Error("UUID generation failed", "err", err.Error())
		return err
	}

	// generate hash from password
	pswHash, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		s.l.Error(helper.ErrPasswordHash.Error(), "err", err.Error())
		return err
	}

	idPsw, err := uuid.NewV7()
	if err != nil {
		s.l.Error(helper.ErrUUIDGeneration.Error(), "err", err.Error())
		return err
	}

	// prepare data for repo
	user := models.User{
		ID:          userID.String(),
		Email:       usr.Email,
		DisplayName: usr.DisplayName,
		UserType:    "user",
	}
	lgg := models.Login{
		ID:      idPsw.String(),
		PswHash: string(pswHash),
		UserID:  userID.String(),
	}

	err = s.r.Create(user, lgg)
	if err != nil {
		if errors.Is(err, helper.ErrDuplicateEntry) {
			return helper.ErrResourceExists
		}
		return helper.ErrCreateFailed
	}

	return nil
}

// Package handler provides HTTP request handlers, routing, and middleware.
// It handles incoming HTTP requests, delegates to the service layer for
// business logic, and returns JSON responses with appropriate status codes.
package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"hotel.com/app/internal/helper"
	"hotel.com/app/internal/models"
	"hotel.com/app/internal/service"
)

type Handler struct {
	s       service.Service
	l       *slog.Logger
	jwtAuth *JWTAuthenticator
}

func New(s service.Service, l *slog.Logger, jwtAuth *JWTAuthenticator) *Handler {
	return &Handler{
		s:       s,
		l:       l,
		jwtAuth: jwtAuth,
	}
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// readinessCheck verifies if the service is ready to accept traffic
// by pinging the database and other critical dependencies.
func (h *Handler) readinessCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.s.Check(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "not ready", "reason": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ready",
		"db":     "ok",
		"from":   "users service",
	})
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var form models.UserCreation
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		helper.RespondError(w, http.StatusBadRequest, helper.ErrInvalidInput.Error())
		return
	}

	if err := h.s.CreateUser(form); err != nil {
		switch {
		case errors.Is(err, helper.ErrResourceExists):
			helper.RespondError(w, http.StatusConflict, "an account with this email already exists")
		default:
			helper.RespondError(w, http.StatusInternalServerError, helper.ErrCreateFailed.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// loginUser handles user login and returns a JWT token
func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var credentials models.LoginInput

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		helper.RespondError(w, http.StatusBadRequest, helper.ErrInvalidInput.Error())
		return
	}

	if credentials.Email == "" || credentials.Password == "" {
		helper.RespondError(w, http.StatusBadRequest, helper.ErrMissingField.Error())
		return
	}

	// Authenticate user
	user, err := h.s.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		switch err {
		case helper.ErrInvalidCredentials:
			helper.RespondError(w, http.StatusUnauthorized, helper.ErrInvalidCredentials.Error())
		case helper.ErrRecordNotFound:
			helper.RespondError(w, http.StatusUnauthorized, helper.ErrInvalidCredentials.Error())
		default:
			helper.RespondError(w, http.StatusInternalServerError, helper.ErrProcessingFailed.Error())
		}
		return
	}

	h.l.Info("User authenticated", "userID", user.ID, "email", user.Email)

	// Generate JWT token
	if h.jwtAuth == nil {
		helper.RespondError(w, http.StatusInternalServerError, "authentication not configured")
		return
	}

	token, err := h.jwtAuth.GenerateToken(user.ID, user.Email, user.UserType)
	if err != nil {
		helper.RespondError(w, http.StatusInternalServerError, helper.ErrTokenGeneration.Error())
		return
	}

	// Return token response (include token_type for client usage)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
		"token_type":   "Bearer",
	})
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromRequest(r)
	if userID == "" {
		helper.RespondError(w, http.StatusUnauthorized, helper.ErrUnauthorized.Error())
		return
	}
	profile, err := h.s.GetProfile(userID)
	if err != nil {
		switch {
		case errors.Is(err, helper.ErrRecordNotFound):
			helper.RespondError(w, http.StatusNotFound, helper.ErrNotFound.Error())
		default:
			helper.RespondError(w, http.StatusInternalServerError, helper.ErrFetchFailed.Error())
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromRequest(r)
	if userID == "" {
		helper.RespondError(w, http.StatusUnauthorized, helper.ErrUnauthorized.Error())
		return
	}
	var user models.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helper.RespondError(w, http.StatusBadRequest, helper.ErrInvalidInput.Error())
		return
	}
	if err := h.s.UpdateProfile(user.DisplayName, userID); err != nil {
		switch {
		case errors.Is(err, helper.ErrRecordNotFound):
			helper.RespondError(w, http.StatusNotFound, helper.ErrNotFound.Error())
		default:
			helper.RespondError(w, http.StatusInternalServerError, helper.ErrUpdateFailed.Error())
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "profile updated"})
}

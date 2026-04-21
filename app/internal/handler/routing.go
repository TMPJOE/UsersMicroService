package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
	"hotel.com/app/internal/helper"
)

func (h *Handler) NewServerMux(rateLimiter *RateLimiter) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(httplog.RequestLogger(h.l, &httplog.Options{
		Level:         slog.LevelDebug,
		Schema:        httplog.SchemaOTEL,
		RecoverPanics: true,
	}))
	r.Use(SecureHeaders)
	r.Use(RequestID)
	r.Use(CORS)

	// Apply rate limiting if enabled
	if rateLimiter != nil {
		r.Use(RateLimitMiddleware(rateLimiter))
	}

	// Custom error handlers (JSON instead of default HTML)
	r.NotFound(h.notFoundHandler)
	r.MethodNotAllowed(h.methodNotAllowedHandler)

	// Public routes - no authentication required
	r.Group(func(r chi.Router) {
		r.Get("/health", h.healthCheck)
		r.Get("/ready", h.readinessCheck)
		r.Post("/register", h.createUser)
		r.Post("/login", h.loginUser)
	})

	// Protected routes - require JWT authentication
	r.Group(func(r chi.Router) {
		r.Use(h.jwtAuth.Middleware()) // JWT authentication middleware

		// Add protected routes here, e.g.:
		r.Get("/profile", h.getProfile)
		r.Put("/profile", h.updateProfile)
	})

	return r
}

func (h *Handler) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	helper.RespondError(w, http.StatusNotFound, "endpoint not found")
}

func (h *Handler) methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	helper.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
}

// GetUserIDFromRequest extracts user ID from the authenticated request
func GetUserIDFromRequest(r *http.Request) string {
	return GetUserIDFromContext(r.Context())
}

// GetUserEmailFromRequest extracts user email from the authenticated request
func GetUserEmailFromRequest(r *http.Request) string {
	return GetUserEmailFromContext(r.Context())
}

// GetClaimsFromRequest extracts JWT claims from the authenticated request
func GetClaimsFromRequest(r *http.Request) *JWTClaims {
	return GetClaimsFromContext(r.Context())
}

package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
	"hotel.com/app/internal/helper"
)

func (h *Handler) NewServerMux() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(httplog.RequestLogger(h.l, &httplog.Options{
		Level:         slog.LevelDebug,
		Schema:        httplog.SchemaOTEL,
		RecoverPanics: true,
	}))
	r.Use(SecureHeaders)
	r.Use(RequestID)

	// Custom error handlers (JSON instead of default HTML)
	r.NotFound(h.notFoundHandler)
	r.MethodNotAllowed(h.methodNotAllowedHandler)

	// Routes
	r.Get("/health", h.healthCheck)
	r.Get("/ready", h.readinessCheck)

	return r
}

func (h *Handler) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	helper.RespondError(w, http.StatusNotFound, "endpoint not found")
}

func (h *Handler) methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	helper.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
}

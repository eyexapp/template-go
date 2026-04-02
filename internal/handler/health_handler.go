package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/user/myapp/internal/handler/response"
)

// RegisterHealthRoutes mounts health check routes.
func (h *Handler) RegisterHealthRoutes(r chi.Router, db *sqlx.DB) {
	r.Get("/health", h.Health)
	r.Get("/readiness", h.Readiness(db))
}

// Health returns a simple liveness check.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// Readiness checks that the database is reachable.
func (h *Handler) Readiness(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			h.Logger.ErrorContext(r.Context(), "readiness check failed", "error", err)
			response.JSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "unavailable",
				"error":  "database unreachable",
			})
			return
		}

		response.JSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	}
}

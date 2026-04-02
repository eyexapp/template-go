package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/user/myapp/internal/handler"
)

func newTestRouter(h *handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	h.RegisterItemRoutes(r)
	return r
}

func newTestHandler() *handler.Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	// Note: ItemService requires a real repo; for handler tests
	// we test the HTTP layer with a nil service to verify routing
	// and JSON parsing. Full integration needs a running DB.
	return handler.NewHandler(nil, logger)
}

func TestHealthHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(nil, logger)

	r := chi.NewRouter()
	r.Get("/health", h.Health)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	err := json.NewDecoder(w.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "ok", body["status"])
}

func TestCreateItem_InvalidBody(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	// Create a mock service for handler tests
	h := handler.NewHandler(nil, logger)

	r := chi.NewRouter()
	h.RegisterItemRoutes(r)

	// Send invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetItem_InvalidID(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewHandler(nil, logger)

	r := chi.NewRouter()
	h.RegisterItemRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/not-a-uuid", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

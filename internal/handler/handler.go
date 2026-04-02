package handler

import (
	"log/slog"

	"github.com/user/myapp/internal/service"
)

// Handler holds shared dependencies for all HTTP handlers.
type Handler struct {
	ItemService *service.ItemService
	Logger      *slog.Logger
}

// NewHandler creates a new Handler.
func NewHandler(itemService *service.ItemService, logger *slog.Logger) *Handler {
	return &Handler{
		ItemService: itemService,
		Logger:      logger,
	}
}

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/user/myapp/internal/domain"
	"github.com/user/myapp/internal/handler/response"
)

// RegisterItemRoutes mounts item CRUD routes on the given router.
func (h *Handler) RegisterItemRoutes(r chi.Router) {
	r.Route("/api/v1/items", func(r chi.Router) {
		r.Post("/", h.CreateItem)
		r.Get("/", h.ListItems)
		r.Get("/{id}", h.GetItem)
		r.Put("/{id}", h.UpdateItem)
		r.Delete("/{id}", h.DeleteItem)
	})
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, domain.NewValidationError("invalid request body"))
		return
	}

	item, err := h.ItemService.Create(r.Context(), input)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Created(w, item)
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, domain.NewValidationError("invalid item id"))
		return
	}

	item, err := h.ItemService.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, item)
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	result, err := h.ItemService.List(r.Context(), page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, domain.NewValidationError("invalid item id"))
		return
	}

	var input domain.UpdateItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, domain.NewValidationError("invalid request body"))
		return
	}

	item, err := h.ItemService.Update(r.Context(), id, input)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, domain.NewValidationError("invalid item id"))
		return
	}

	if err := h.ItemService.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}

	response.NoContent(w)
}

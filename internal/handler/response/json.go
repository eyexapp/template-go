package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/user/myapp/internal/domain"
)

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// Created writes a 201 JSON response.
func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, data)
}

// NoContent writes a 204 response with no body.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error writes an error JSON response, mapping AppError to the correct status.
func Error(w http.ResponseWriter, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		JSON(w, appErr.StatusCode, map[string]string{
			"error": appErr.Message,
			"code":  appErr.Code,
		})
		return
	}

	JSON(w, http.StatusInternalServerError, map[string]string{
		"error": "internal server error",
		"code":  "INTERNAL_ERROR",
	})
}

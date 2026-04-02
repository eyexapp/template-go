package domain

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors for common domain failures.
var (
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
	ErrValidation = errors.New("validation error")
)

// AppError is a structured error with an HTTP status code.
type AppError struct {
	Message    string `json:"message"`
	Code       string `json:"code,omitempty"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewNotFoundError creates a 404 error.
func NewNotFoundError(msg string) *AppError {
	return &AppError{
		Message:    msg,
		Code:       "NOT_FOUND",
		StatusCode: http.StatusNotFound,
		Err:        ErrNotFound,
	}
}

// NewConflictError creates a 409 error.
func NewConflictError(msg string) *AppError {
	return &AppError{
		Message:    msg,
		Code:       "CONFLICT",
		StatusCode: http.StatusConflict,
		Err:        ErrConflict,
	}
}

// NewValidationError creates a 400 error.
func NewValidationError(msg string) *AppError {
	return &AppError{
		Message:    msg,
		Code:       "VALIDATION_ERROR",
		StatusCode: http.StatusBadRequest,
		Err:        ErrValidation,
	}
}

// NewInternalError creates a 500 error.
func NewInternalError(msg string, err error) *AppError {
	return &AppError{
		Message:    msg,
		Code:       "INTERNAL_ERROR",
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

package domain_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/myapp/internal/domain"
)

func TestNewNotFoundError(t *testing.T) {
	err := domain.NewNotFoundError("item not found")

	assert.Equal(t, "item not found", err.Message)
	assert.Equal(t, "NOT_FOUND", err.Code)
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestNewConflictError(t *testing.T) {
	err := domain.NewConflictError("item already exists")

	assert.Equal(t, "item already exists", err.Message)
	assert.Equal(t, "CONFLICT", err.Code)
	assert.Equal(t, http.StatusConflict, err.StatusCode)
	assert.True(t, errors.Is(err, domain.ErrConflict))
}

func TestNewValidationError(t *testing.T) {
	err := domain.NewValidationError("title is required")

	assert.Equal(t, "title is required", err.Message)
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.True(t, errors.Is(err, domain.ErrValidation))
}

func TestNewInternalError(t *testing.T) {
	cause := errors.New("db connection lost")
	err := domain.NewInternalError("something went wrong", cause)

	assert.Equal(t, "something went wrong", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	assert.True(t, errors.Is(err, cause))
	assert.Contains(t, err.Error(), "db connection lost")
}

func TestAppError_ErrorString(t *testing.T) {
	err := domain.NewNotFoundError("item not found")
	assert.Equal(t, "item not found: not found", err.Error())

	err2 := &domain.AppError{Message: "bare error"}
	assert.Equal(t, "bare error", err2.Error())
}

func TestCreateItemInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.CreateItemInput
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   domain.CreateItemInput{Title: "My Item", Description: "desc"},
			wantErr: false,
		},
		{
			name:    "empty title",
			input:   domain.CreateItemInput{Title: "", Description: "desc"},
			wantErr: true,
		},
		{
			name:    "title too long",
			input:   domain.CreateItemInput{Title: string(make([]byte, 256)), Description: "desc"},
			wantErr: true,
		},
		{
			name:    "empty description is ok",
			input:   domain.CreateItemInput{Title: "Title"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

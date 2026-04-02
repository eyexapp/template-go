package domain

import (
	"time"

	"github.com/google/uuid"
)

// Item represents a domain entity.
type Item struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateItemInput holds the data needed to create an Item.
type CreateItemInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateItemInput holds the data needed to update an Item.
type UpdateItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

// Validate checks CreateItemInput fields.
func (c CreateItemInput) Validate() error {
	if c.Title == "" {
		return NewValidationError("title is required")
	}
	if len(c.Title) > 255 {
		return NewValidationError("title must be at most 255 characters")
	}
	return nil
}

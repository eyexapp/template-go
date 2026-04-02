package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/myapp/internal/domain"
)

// ItemRepository defines the data access contract for items.
type ItemRepository interface {
	Create(ctx context.Context, item *domain.Item) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Item, error)
	List(ctx context.Context, limit, offset int) ([]domain.Item, error)
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int, error)
}

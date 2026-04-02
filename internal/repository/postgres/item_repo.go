package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/user/myapp/internal/domain"
)

// ItemRepo implements repository.ItemRepository using sqlx.
type ItemRepo struct {
	db *sqlx.DB
}

// NewItemRepo creates a new ItemRepo.
func NewItemRepo(db *sqlx.DB) *ItemRepo {
	return &ItemRepo{db: db}
}

func (r *ItemRepo) Create(ctx context.Context, item *domain.Item) error {
	item.ID = uuid.New()
	item.CreatedAt = time.Now().UTC()
	item.UpdatedAt = item.CreatedAt

	query := `INSERT INTO items (id, title, description, created_at, updated_at)
	           VALUES (:id, :title, :description, :created_at, :updated_at)`

	_, err := r.db.NamedExecContext(ctx, query, item)
	return err
}

func (r *ItemRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	var item domain.Item
	query := `SELECT id, title, description, created_at, updated_at FROM items WHERE id = $1`

	if err := r.db.GetContext(ctx, &item, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("item not found")
		}
		return nil, err
	}

	return &item, nil
}

func (r *ItemRepo) List(ctx context.Context, limit, offset int) ([]domain.Item, error) {
	var items []domain.Item
	query := `SELECT id, title, description, created_at, updated_at
	           FROM items ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	if err := r.db.SelectContext(ctx, &items, query, limit, offset); err != nil {
		return nil, err
	}

	if items == nil {
		items = []domain.Item{}
	}

	return items, nil
}

func (r *ItemRepo) Update(ctx context.Context, item *domain.Item) error {
	item.UpdatedAt = time.Now().UTC()

	query := `UPDATE items SET title = :title, description = :description, updated_at = :updated_at
	           WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.NewNotFoundError("item not found")
	}

	return nil
}

func (r *ItemRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM items WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.NewNotFoundError("item not found")
	}

	return nil
}

func (r *ItemRepo) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM items`

	if err := r.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}

	return count, nil
}

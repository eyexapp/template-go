package service_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/user/myapp/internal/domain"
	"github.com/user/myapp/internal/service"
)

// mockItemRepo is a manual mock for repository.ItemRepository.
type mockItemRepo struct {
	items []domain.Item
}

func (m *mockItemRepo) Create(_ context.Context, item *domain.Item) error {
	item.ID = uuid.New()
	m.items = append(m.items, *item)
	return nil
}

func (m *mockItemRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Item, error) {
	for i := range m.items {
		if m.items[i].ID == id {
			return &m.items[i], nil
		}
	}
	return nil, domain.NewNotFoundError("item not found")
}

func (m *mockItemRepo) List(_ context.Context, limit, offset int) ([]domain.Item, error) {
	end := offset + limit
	if end > len(m.items) {
		end = len(m.items)
	}
	if offset >= len(m.items) {
		return []domain.Item{}, nil
	}
	return m.items[offset:end], nil
}

func (m *mockItemRepo) Update(_ context.Context, item *domain.Item) error {
	for i := range m.items {
		if m.items[i].ID == item.ID {
			m.items[i] = *item
			return nil
		}
	}
	return domain.NewNotFoundError("item not found")
}

func (m *mockItemRepo) Delete(_ context.Context, id uuid.UUID) error {
	for i := range m.items {
		if m.items[i].ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return nil
		}
	}
	return domain.NewNotFoundError("item not found")
}

func (m *mockItemRepo) Count(_ context.Context) (int, error) {
	return len(m.items), nil
}

func newTestService() (*service.ItemService, *mockItemRepo) {
	repo := &mockItemRepo{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := service.NewItemService(repo, logger)
	return svc, repo
}

func TestItemService_Create_Success(t *testing.T) {
	svc, _ := newTestService()

	item, err := svc.Create(context.Background(), domain.CreateItemInput{
		Title:       "Test Item",
		Description: "A test description",
	})

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "Test Item", item.Title)
	assert.NotEqual(t, uuid.Nil, item.ID)
}

func TestItemService_Create_ValidationError(t *testing.T) {
	svc, _ := newTestService()

	item, err := svc.Create(context.Background(), domain.CreateItemInput{
		Title: "", // empty title
	})

	assert.Error(t, err)
	assert.Nil(t, item)
}

func TestItemService_GetByID_NotFound(t *testing.T) {
	svc, _ := newTestService()

	item, err := svc.GetByID(context.Background(), uuid.New())

	assert.Error(t, err)
	assert.Nil(t, item)
}

func TestItemService_List_Pagination(t *testing.T) {
	svc, _ := newTestService()

	// Create 5 items
	for i := range 5 {
		_, err := svc.Create(context.Background(), domain.CreateItemInput{
			Title: "Item " + string(rune('A'+i)),
		})
		assert.NoError(t, err)
	}

	result, err := svc.List(context.Background(), 1, 2)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Items))
	assert.Equal(t, 5, result.Total)
	assert.Equal(t, 3, result.TotalPages)
}

func TestItemService_Update_Success(t *testing.T) {
	svc, _ := newTestService()

	created, err := svc.Create(context.Background(), domain.CreateItemInput{
		Title:       "Original",
		Description: "Original desc",
	})
	assert.NoError(t, err)

	newTitle := "Updated"
	updated, err := svc.Update(context.Background(), created.ID, domain.UpdateItemInput{
		Title: &newTitle,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Updated", updated.Title)
	assert.Equal(t, "Original desc", updated.Description)
}

func TestItemService_Delete_Success(t *testing.T) {
	svc, _ := newTestService()

	created, err := svc.Create(context.Background(), domain.CreateItemInput{
		Title: "To Delete",
	})
	assert.NoError(t, err)

	err = svc.Delete(context.Background(), created.ID)
	assert.NoError(t, err)

	_, err = svc.GetByID(context.Background(), created.ID)
	assert.Error(t, err)
}

package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/user/myapp/internal/domain"
	"github.com/user/myapp/internal/repository"
)

// ItemService handles business logic for items.
type ItemService struct {
	repo   repository.ItemRepository
	logger *slog.Logger
}

// NewItemService creates a new ItemService.
func NewItemService(repo repository.ItemRepository, logger *slog.Logger) *ItemService {
	return &ItemService{
		repo:   repo,
		logger: logger,
	}
}

// ListResult holds paginated results.
type ListResult struct {
	Items      []domain.Item `json:"items"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

func (s *ItemService) Create(ctx context.Context, input domain.CreateItemInput) (*domain.Item, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	item := &domain.Item{
		Title:       input.Title,
		Description: input.Description,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		s.logger.ErrorContext(ctx, "failed to create item", "error", err)
		return nil, domain.NewInternalError("failed to create item", err)
	}

	return item, nil
}

func (s *ItemService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ItemService) List(ctx context.Context, page, pageSize int) (*ListResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	items, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to list items", "error", err)
		return nil, domain.NewInternalError("failed to list items", err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to count items", "error", err)
		return nil, domain.NewInternalError("failed to count items", err)
	}

	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	return &ListResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *ItemService) Update(ctx context.Context, id uuid.UUID, input domain.UpdateItemInput) (*domain.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		item.Title = *input.Title
	}
	if input.Description != nil {
		item.Description = *input.Description
	}

	if err := s.repo.Update(ctx, item); err != nil {
		s.logger.ErrorContext(ctx, "failed to update item", "error", err, "id", id)
		return nil, domain.NewInternalError("failed to update item", err)
	}

	return item, nil
}

func (s *ItemService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

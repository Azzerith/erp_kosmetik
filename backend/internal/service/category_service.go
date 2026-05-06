package service

import (
	"context"

	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"

	"go.uber.org/zap"
)

type CategoryService interface {
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error)
	CreateCategory(ctx context.Context, category *models.Category) error
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, id uint64) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	logger       *zap.Logger
}

func NewCategoryService(categoryRepo repository.CategoryRepository, logger *zap.Logger) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	return s.categoryRepo.GetAll(ctx)
}

func (s *categoryService) GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error) {
	return s.categoryRepo.FindBySlug(ctx, slug)
}

func (s *categoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	return s.categoryRepo.Create(ctx, category)
}

func (s *categoryService) UpdateCategory(ctx context.Context, category *models.Category) error {
	return s.categoryRepo.Update(ctx, category)
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uint64) error {
	return s.categoryRepo.Delete(ctx, id)
}
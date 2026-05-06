package service

import (
	"context"
	"errors"

	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReviewService interface {
	CreateReview(ctx context.Context, review *models.Review) error
	UpdateReview(ctx context.Context, review *models.Review) error
	DeleteReview(ctx context.Context, id uint64) error
	GetReviewsByProduct(ctx context.Context, productID uint64, page, limit int) ([]models.Review, int64, error)
	MarkHelpful(ctx context.Context, reviewID, userID uint64) error
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
	logger     *zap.Logger
}

func NewReviewService(reviewRepo repository.ReviewRepository, logger *zap.Logger) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		logger:     logger,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, review *models.Review) error {
	// Check if user has purchased the product (verified purchase)
	// This would require checking orders table
	// For now, set based on some logic

	return s.reviewRepo.Create(ctx, review)
}

func (s *reviewService) UpdateReview(ctx context.Context, review *models.Review) error {
	existing, err := s.reviewRepo.FindByID(ctx, review.ID)
	if err != nil {
		return err
	}
	if existing.UserID != review.UserID {
		return errors.New("unauthorized to update this review")
	}

	return s.reviewRepo.Update(ctx, review)
}

func (s *reviewService) DeleteReview(ctx context.Context, id uint64) error {
	return s.reviewRepo.Delete(ctx, id)
}

func (s *reviewService) GetReviewsByProduct(ctx context.Context, productID uint64, page, limit int) ([]models.Review, int64, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.reviewRepo.GetByProductID(ctx, productID, offset, limit)
}

func (s *reviewService) MarkHelpful(ctx context.Context, reviewID, userID uint64) error {
	return s.reviewRepo.MarkHelpful(ctx, reviewID, userID)
}
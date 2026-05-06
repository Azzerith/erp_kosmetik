package repository

import (
	"context"
	"erp-cosmetics-backend/internal/models"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	FindByID(ctx context.Context, id uint64) (*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uint64) error
	GetByProductID(ctx context.Context, productID uint64, offset, limit int) ([]models.Review, int64, error)
	GetAverageRating(ctx context.Context, productID uint64) (float64, error)
	MarkHelpful(ctx context.Context, reviewID, userID uint64) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepository) FindByID(ctx context.Context, id uint64) (*models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).First(&review, id).Error
	return &review, err
}

func (r *reviewRepository) Update(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *reviewRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, id).Error
}

func (r *reviewRepository) GetByProductID(ctx context.Context, productID uint64, offset, limit int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Review{}).Where("product_id = ? AND is_approved = ?", productID, true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).
		Preload("User").
		Order("created_at DESC").
		Find(&reviews).Error

	return reviews, total, err
}

func (r *reviewRepository) GetAverageRating(ctx context.Context, productID uint64) (float64, error) {
	var avg float64
	err := r.db.WithContext(ctx).Model(&models.Review{}).
		Where("product_id = ? AND is_approved = ?", productID, true).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avg).Error
	return avg, err
}

func (r *reviewRepository) MarkHelpful(ctx context.Context, reviewID, userID uint64) error {
	// Check if already marked
	var existing models.ReviewHelpful
	err := r.db.WithContext(ctx).Where("review_id = ? AND user_id = ?", reviewID, userID).First(&existing).Error
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Create helpful record
	helpful := models.ReviewHelpful{
		ReviewID:  reviewID,
		UserID:    userID,
		IsHelpful: true,
	}
	if err := r.db.Create(&helpful).Error; err != nil {
		return err
	}

	// Increment helpful count
	return r.db.Model(&models.Review{}).Where("id = ?", reviewID).
		Update("helpful_count", gorm.Expr("helpful_count + 1")).Error
}
package repository

import (
	"context"
	"erp-cosmetics-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type TrendRepository interface {
	CreateTrendData(ctx context.Context, data *models.TrendData) error
	GetLatestTrendData(ctx context.Context, keyword string, source string) (*models.TrendData, error)
	CreateProductTrendMapping(ctx context.Context, mapping *models.ProductTrendMapping) error
	GetProductTrendMappings(ctx context.Context, productID uint64) ([]models.ProductTrendMapping, error)
	CreateTrendHistory(ctx context.Context, history *models.TrendScoreHistory) error
	GetTrendHistory(ctx context.Context, productID uint64, days int) ([]models.TrendScoreHistory, error)
}

type trendRepository struct {
	db *gorm.DB
}

func NewTrendRepository(db *gorm.DB) TrendRepository {
	return &trendRepository{db: db}
}

func (r *trendRepository) CreateTrendData(ctx context.Context, data *models.TrendData) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *trendRepository) GetLatestTrendData(ctx context.Context, keyword string, source string) (*models.TrendData, error) {
	var data models.TrendData
	err := r.db.WithContext(ctx).
		Where("keyword = ? AND source = ?", keyword, source).
		Order("recorded_at DESC").
		First(&data).Error
	return &data, err
}

func (r *trendRepository) CreateProductTrendMapping(ctx context.Context, mapping *models.ProductTrendMapping) error {
	return r.db.WithContext(ctx).Create(mapping).Error
}

func (r *trendRepository) GetProductTrendMappings(ctx context.Context, productID uint64) ([]models.ProductTrendMapping, error) {
	var mappings []models.ProductTrendMapping
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Preload("TrendData").
		Find(&mappings).Error
	return mappings, err
}

func (r *trendRepository) CreateTrendHistory(ctx context.Context, history *models.TrendScoreHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *trendRepository) GetTrendHistory(ctx context.Context, productID uint64, days int) ([]models.TrendScoreHistory, error) {
	var history []models.TrendScoreHistory
	since := time.Now().AddDate(0, 0, -days)
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND recorded_at >= ?", productID, since).
		Order("recorded_at ASC").
		Find(&history).Error
	return history, err
}
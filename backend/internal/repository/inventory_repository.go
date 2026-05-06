package repository

import (
	"context"
	"erp-cosmetics-backend/internal/models"

	"gorm.io/gorm"
)

type InventoryRepository interface {
	CreateLog(ctx context.Context, log *models.InventoryLog) error
	GetLogs(ctx context.Context, productID uint64, offset, limit int) ([]models.InventoryLog, int64, error)
	GetLowStockProducts(ctx context.Context, threshold int) ([]models.Product, error)
	CreateReservation(ctx context.Context, reservation *models.StockReservation) error
	ReleaseReservation(ctx context.Context, orderID uint64) error
	GetExpiredReservations(ctx context.Context) ([]models.StockReservation, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) CreateLog(ctx context.Context, log *models.InventoryLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *inventoryRepository) GetLogs(ctx context.Context, productID uint64, offset, limit int) ([]models.InventoryLog, int64, error) {
	var logs []models.InventoryLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.InventoryLog{}).Where("product_id = ?", productID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

func (r *inventoryRepository) GetLowStockProducts(ctx context.Context, threshold int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND stock <= min_stock_threshold", true).
		Find(&products).Error
	return products, err
}

func (r *inventoryRepository) CreateReservation(ctx context.Context, reservation *models.StockReservation) error {
	return r.db.WithContext(ctx).Create(reservation).Error
}

func (r *inventoryRepository) ReleaseReservation(ctx context.Context, orderID uint64) error {
	return r.db.WithContext(ctx).Model(&models.StockReservation{}).
		Where("order_id = ? AND is_released = ?", orderID, false).
		Updates(map[string]interface{}{
			"is_released": true,
			"released_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *inventoryRepository) GetExpiredReservations(ctx context.Context) ([]models.StockReservation, error) {
	var reservations []models.StockReservation
	err := r.db.WithContext(ctx).
		Where("expires_at <= NOW() AND is_released = ?", false).
		Find(&reservations).Error
	return reservations, err
}
package repository

import (
	"context"
	"erp-cosmetics-backend/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	CreateWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error
	FindByID(ctx context.Context, id uint64) (*models.Order, error)
	FindByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error)
	FindByUserID(ctx context.Context, userID uint64, offset, limit int) ([]models.Order, int64, error)
	Update(ctx context.Context, order *models.Order) error
	UpdateStatus(ctx context.Context, orderID uint64, status, paymentStatus, fulfillmentStatus string) error
	AddStatusHistory(ctx context.Context, history *models.OrderStatusHistory) error
	GetAll(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]models.Order, int64, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) CreateWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OrderID = order.ID
		}
		return tx.Create(&items).Error
	})
}

func (r *orderRepository) FindByID(ctx context.Context, id uint64) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Variant").
		Preload("ShippingAddress").
		Preload("Payments").
		Preload("StatusHistory").
		First(&order, id).Error
	return &order, err
}

func (r *orderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Variant").
		Preload("ShippingAddress").
		Preload("Payments").
		Where("order_number = ?", orderNumber).
		First(&order).Error
	return &order, err
}

func (r *orderRepository) FindByUserID(ctx context.Context, userID uint64, offset, limit int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Order{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Preload("Items").
		Preload("Items.Product").
		Preload("ShippingAddress").
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID uint64, status, paymentStatus, fulfillmentStatus string) error {
	updates := make(map[string]interface{})
	if status != "" {
		updates["status"] = status
	}
	if paymentStatus != "" {
		updates["payment_status"] = paymentStatus
	}
	if fulfillmentStatus != "" {
		updates["fulfillment_status"] = fulfillmentStatus
	}
	return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", orderID).Updates(updates).Error
}

func (r *orderRepository) AddStatusHistory(ctx context.Context, history *models.OrderStatusHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *orderRepository) GetAll(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Order{})

	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if dateFrom, ok := filters["date_from"]; ok {
		query = query.Where("created_at >= ?", dateFrom)
	}
	if dateTo, ok := filters["date_to"]; ok {
		query = query.Where("created_at <= ?", dateTo)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("ShippingAddress").
		Find(&orders).Error

	return orders, total, err
}
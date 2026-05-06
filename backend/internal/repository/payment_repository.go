package repository

import (
	"context"
	"erp-cosmetics-backend/internal/models"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	Update(ctx context.Context, payment *models.Payment) error
	FindByOrderID(ctx context.Context, orderID uint64) (*models.Payment, error)
	FindByTransactionID(ctx context.Context, transactionID string) (*models.Payment, error)
	CreateLog(ctx context.Context, log *models.PaymentLog) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *paymentRepository) FindByOrderID(ctx context.Context, orderID uint64) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

func (r *paymentRepository) FindByTransactionID(ctx context.Context, transactionID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&payment).Error
	return &payment, err
}

func (r *paymentRepository) CreateLog(ctx context.Context, log *models.PaymentLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
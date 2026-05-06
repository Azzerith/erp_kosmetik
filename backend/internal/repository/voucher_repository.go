package repository

import (
	"context"
	"erp-cosmetics-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type VoucherRepository interface {
	Create(ctx context.Context, voucher *models.Voucher) error
	FindByID(ctx context.Context, id uint64) (*models.Voucher, error)
	FindByCode(ctx context.Context, code string) (*models.Voucher, error)
	Update(ctx context.Context, voucher *models.Voucher) error
	Delete(ctx context.Context, id uint64) error
	GetValidVouchers(ctx context.Context) ([]models.Voucher, error)
	IncrementUsage(ctx context.Context, id uint64) error
	CreateUsage(ctx context.Context, usage *models.VoucherUsage) error
	GetUserUsageCount(ctx context.Context, voucherID, userID uint64) (int, error)
}

type voucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) VoucherRepository {
	return &voucherRepository{db: db}
}

func (r *voucherRepository) Create(ctx context.Context, voucher *models.Voucher) error {
	return r.db.WithContext(ctx).Create(voucher).Error
}

func (r *voucherRepository) FindByID(ctx context.Context, id uint64) (*models.Voucher, error) {
	var voucher models.Voucher
	err := r.db.WithContext(ctx).First(&voucher, id).Error
	return &voucher, err
}

func (r *voucherRepository) FindByCode(ctx context.Context, code string) (*models.Voucher, error) {
	var voucher models.Voucher
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&voucher).Error
	return &voucher, err
}

func (r *voucherRepository) Update(ctx context.Context, voucher *models.Voucher) error {
	return r.db.WithContext(ctx).Save(voucher).Error
}

func (r *voucherRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&models.Voucher{}, id).Error
}

func (r *voucherRepository) GetValidVouchers(ctx context.Context) ([]models.Voucher, error) {
	var vouchers []models.Voucher
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND valid_from <= ? AND valid_until >= ? AND (usage_limit IS NULL OR used_count < usage_limit)",
			true, now, now).
		Find(&vouchers).Error
	return vouchers, err
}

func (r *voucherRepository) IncrementUsage(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&models.Voucher{}).Where("id = ?", id).
		Update("used_count", gorm.Expr("used_count + 1")).Error
}

func (r *voucherRepository) CreateUsage(ctx context.Context, usage *models.VoucherUsage) error {
	return r.db.WithContext(ctx).Create(usage).Error
}

func (r *voucherRepository) GetUserUsageCount(ctx context.Context, voucherID, userID uint64) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.VoucherUsage{}).
		Where("voucher_id = ? AND user_id = ?", voucherID, userID).
		Count(&count).Error
	return int(count), err
}
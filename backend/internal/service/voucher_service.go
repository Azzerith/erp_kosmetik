package service

import (
	"context"
	"errors"
	"time"

	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"

	"go.uber.org/zap"
)

type VoucherService interface {
	ValidateVoucher(ctx context.Context, code string, orderAmount float64) (*VoucherValidationResult, error)
	ApplyVoucher(ctx context.Context, voucherID, orderID, userID uint64) error
	CreateVoucher(ctx context.Context, req *CreateVoucherRequest) (*models.Voucher, error)
	UpdateVoucher(ctx context.Context, id uint64, req *UpdateVoucherRequest) (*models.Voucher, error)
	DeleteVoucher(ctx context.Context, id uint64) error
	GetVouchers(ctx context.Context) ([]models.Voucher, error)
}

type VoucherValidationResult struct {
	Valid          bool    `json:"valid"`
	DiscountAmount float64 `json:"discount_amount"`
	Message        string  `json:"message"`
	Voucher        *models.Voucher `json:"voucher,omitempty"`
}

type CreateVoucherRequest struct {
	Code             string     `json:"code" binding:"required"`
	Name             string     `json:"name" binding:"required"`
	Description      *string    `json:"description"`
	Type             string     `json:"type" binding:"required,oneof=percentage fixed_amount free_shipping"`
	Value            float64    `json:"value" binding:"required,min=0"`
	MaxDiscountAmount *float64   `json:"max_discount_amount"`
	MinOrderAmount   float64    `json:"min_order_amount"`
	ApplicableType   string     `json:"applicable_type"`
	ApplicableIDs    []uint64   `json:"applicable_ids"`
	UsageLimit       *int       `json:"usage_limit"`
	UsagePerUser     *int       `json:"usage_per_user"`
	ValidFrom        time.Time  `json:"valid_from" binding:"required"`
	ValidUntil       time.Time  `json:"valid_until" binding:"required"`
}

type UpdateVoucherRequest struct {
	Name             *string    `json:"name"`
	Description      *string    `json:"description"`
	Value            *float64   `json:"value"`
	MaxDiscountAmount *float64   `json:"max_discount_amount"`
	MinOrderAmount   *float64   `json:"min_order_amount"`
	UsageLimit       *int       `json:"usage_limit"`
	UsagePerUser     *int       `json:"usage_per_user"`
	ValidFrom        *time.Time `json:"valid_from"`
	ValidUntil       *time.Time `json:"valid_until"`
	IsActive         *bool      `json:"is_active"`
}

type voucherService struct {
	voucherRepo repository.VoucherRepository
	logger      *zap.Logger
}

func NewVoucherService(voucherRepo repository.VoucherRepository, logger *zap.Logger) VoucherService {
	return &voucherService{
		voucherRepo: voucherRepo,
		logger:      logger,
	}
}

func (s *voucherService) ValidateVoucher(ctx context.Context, code string, orderAmount float64) (*VoucherValidationResult, error) {
	// Find voucher by code
	voucher, err := s.voucherRepo.FindByCode(ctx, code)
	if err != nil {
		return &VoucherValidationResult{
			Valid:   false,
			Message: "Voucher tidak ditemukan",
		}, nil
	}

	// Check if voucher is active
	if !voucher.IsActive {
		return &VoucherValidationResult{
			Valid:   false,
			Message: "Voucher tidak aktif",
		}, nil
	}

	// Check validity period
	now := time.Now()
	if now.Before(voucher.ValidFrom) {
		return &VoucherValidationResult{
			Valid:   false,
			Message: "Voucher belum tersedia",
		}, nil
	}
	if now.After(voucher.ValidUntil) {
		return &VoucherValidationResult{
			Valid:   false,
			Message: "Voucher sudah kadaluarsa",
		}, nil
	}

	// Check usage limit
	if voucher.UsageLimit != nil && voucher.UsedCount >= *voucher.UsageLimit {
		return &VoucherValidationResult{
			Valid:   false,
			Message: "Voucher sudah mencapai batas penggunaan",
		}, nil
	}

	// Check minimum order amount
	if orderAmount < voucher.MinOrderAmount {
		return &VoucherValidationResult{
			Valid:   false,
			Message: "Minimal belanja " + formatCurrency(voucher.MinOrderAmount) + " untuk menggunakan voucher ini",
		}, nil
	}

	// Calculate discount amount
	discountAmount := s.calculateDiscount(voucher, orderAmount)

	return &VoucherValidationResult{
		Valid:          true,
		DiscountAmount: discountAmount,
		Message:        "Voucher valid",
		Voucher:        voucher,
	}, nil
}

func (s *voucherService) ApplyVoucher(ctx context.Context, voucherID, orderID, userID uint64) error {
	// Check usage per user
	voucher, err := s.voucherRepo.FindByID(ctx, voucherID)
	if err != nil {
		return err
	}

	if voucher.UsagePerUser != nil {
		userUsage, err := s.voucherRepo.GetUserUsageCount(ctx, voucherID, userID)
		if err != nil {
			return err
		}
		if userUsage >= *voucher.UsagePerUser {
			return errors.New("voucher sudah digunakan maksimal oleh user ini")
		}
	}

	// Create usage record
	usage := &models.VoucherUsage{
		VoucherID: voucherID,
		UserID:    userID,
		OrderID:   orderID,
	}

	if err := s.voucherRepo.CreateUsage(ctx, usage); err != nil {
		return err
	}

	// Increment usage count
	return s.voucherRepo.IncrementUsage(ctx, voucherID)
}

func (s *voucherService) CreateVoucher(ctx context.Context, req *CreateVoucherRequest) (*models.Voucher, error) {
	voucher := &models.Voucher{
		Code:              req.Code,
		Name:              req.Name,
		Description:       req.Description,
		Type:              req.Type,
		Value:             req.Value,
		MaxDiscountAmount: req.MaxDiscountAmount,
		MinOrderAmount:    req.MinOrderAmount,
		ApplicableType:    req.ApplicableType,
		UsageLimit:        req.UsageLimit,
		UsagePerUser:      req.UsagePerUser,
		ValidFrom:         req.ValidFrom,
		ValidUntil:        req.ValidUntil,
		IsActive:          true,
	}

	if err := s.voucherRepo.Create(ctx, voucher); err != nil {
		return nil, err
	}

	return voucher, nil
}

func (s *voucherService) UpdateVoucher(ctx context.Context, id uint64, req *UpdateVoucherRequest) (*models.Voucher, error) {
	voucher, err := s.voucherRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		voucher.Name = *req.Name
	}
	if req.Description != nil {
		voucher.Description = req.Description
	}
	if req.Value != nil {
		voucher.Value = *req.Value
	}
	if req.MaxDiscountAmount != nil {
		voucher.MaxDiscountAmount = req.MaxDiscountAmount
	}
	if req.MinOrderAmount != nil {
		voucher.MinOrderAmount = *req.MinOrderAmount
	}
	if req.UsageLimit != nil {
		voucher.UsageLimit = req.UsageLimit
	}
	if req.UsagePerUser != nil {
		voucher.UsagePerUser = req.UsagePerUser
	}
	if req.ValidFrom != nil {
		voucher.ValidFrom = *req.ValidFrom
	}
	if req.ValidUntil != nil {
		voucher.ValidUntil = *req.ValidUntil
	}
	if req.IsActive != nil {
		voucher.IsActive = *req.IsActive
	}

	if err := s.voucherRepo.Update(ctx, voucher); err != nil {
		return nil, err
	}

	return voucher, nil
}

func (s *voucherService) DeleteVoucher(ctx context.Context, id uint64) error {
	return s.voucherRepo.Delete(ctx, id)
}

func (s *voucherService) GetVouchers(ctx context.Context) ([]models.Voucher, error) {
	return s.voucherRepo.GetValidVouchers(ctx)
}

func (s *voucherService) calculateDiscount(voucher *models.Voucher, orderAmount float64) float64 {
	switch voucher.Type {
	case "percentage":
		discount := orderAmount * (voucher.Value / 100)
		if voucher.MaxDiscountAmount != nil && discount > *voucher.MaxDiscountAmount {
			discount = *voucher.MaxDiscountAmount
		}
		return discount
	case "fixed_amount":
		if discount := voucher.Value; discount <= orderAmount {
			return discount
		}
		return orderAmount
	case "free_shipping":
		return 0 // Free shipping handled separately
	default:
		return 0
	}
}

func formatCurrency(amount float64) string {
	return "Rp " + formatNumber(amount)
}

func formatNumber(amount float64) string {
	return "0" // Simplified
}
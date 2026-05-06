package models

import (
	"time"

	"gorm.io/datatypes"
)

type Voucher struct {
	ID               uint64         `gorm:"primaryKey" json:"id"`
	Code             string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name             string         `gorm:"type:varchar(255);not null" json:"name"`
	Description      *string        `gorm:"type:text" json:"description,omitempty"`
	Type             string         `gorm:"type:enum('percentage','fixed_amount','free_shipping');not null" json:"type"`
	Value            float64        `gorm:"type:decimal(15,2);not null" json:"value"`
	MaxDiscountAmount *float64      `gorm:"type:decimal(15,2)" json:"max_discount_amount,omitempty"`
	MinOrderAmount   float64        `gorm:"type:decimal(15,2);default:0" json:"min_order_amount"`
	ApplicableType   string         `gorm:"type:enum('all','specific_products','specific_categories');default:'all'" json:"applicable_type"`
	ApplicableIDs    datatypes.JSON `gorm:"type:json" json:"applicable_ids,omitempty"`
	UsageLimit       *int           `json:"usage_limit,omitempty"`
	UsagePerUser     *int           `json:"usage_per_user,omitempty"`
	UsedCount        int            `gorm:"default:0" json:"used_count"`
	ValidFrom        time.Time      `gorm:"not null" json:"valid_from"`
	ValidUntil       time.Time      `gorm:"not null" json:"valid_until"`
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	CreatedBy        uint64         `gorm:"not null" json:"created_by"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`

	// Relationships
	Creator User `json:"creator,omitempty"`
	Usages  []VoucherUsage `json:"usages,omitempty"`
}

func (Voucher) TableName() string {
	return "vouchers"
}

type VoucherUsage struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	VoucherID      uint64    `gorm:"not null;index" json:"voucher_id"`
	UserID         uint64    `gorm:"not null;index" json:"user_id"`
	OrderID        uint64    `gorm:"not null;index" json:"order_id"`
	DiscountAmount float64   `gorm:"type:decimal(15,2);not null" json:"discount_amount"`
	UsedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"used_at"`

	// Relationships
	Voucher Voucher `json:"-"`
	User    User    `json:"-"`
	Order   Order   `json:"-"`
}

func (VoucherUsage) TableName() string {
	return "voucher_usages"
}
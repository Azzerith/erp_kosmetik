package models

import (
	"time"
)

type StockReservation struct {
	ID         uint64     `gorm:"primaryKey" json:"id"`
	ProductID  uint64     `gorm:"not null;index" json:"product_id"`
	VariantID  *uint64    `json:"variant_id,omitempty"`
	OrderID    uint64     `gorm:"not null;index" json:"order_id"`
	Quantity   int        `gorm:"not null" json:"quantity"`
	ReservedAt time.Time  `gorm:"autoCreateTime" json:"reserved_at"`
	ExpiresAt  time.Time  `gorm:"not null" json:"expires_at"`
	IsReleased bool       `gorm:"default:false" json:"is_released"`
	ReleasedAt *time.Time `json:"released_at,omitempty"`

	// Relationships
	Product Product `json:"product,omitempty"`
	Variant *ProductVariant `json:"variant,omitempty"`
	Order   Order   `json:"order,omitempty"`
}

func (StockReservation) TableName() string {
	return "stock_reservations"
}
package models

import (
	"time"
)

type CartItem struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	CartID         uint64     `gorm:"not null;index" json:"cart_id"`
	ProductID      uint64     `gorm:"not null" json:"product_id"`
	VariantID      *uint64    `json:"variant_id,omitempty"`
	Quantity       int        `gorm:"not null;default:1" json:"quantity"`
	PriceSnapshot  float64    `gorm:"type:decimal(15,2);not null" json:"price_snapshot"`
	Notes          *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Cart    Cart    `json:"-"`
	Product Product `json:"product,omitempty"`
	Variant *ProductVariant `json:"variant,omitempty"`
}

func (CartItem) TableName() string {
	return "cart_items"
}
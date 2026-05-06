package models

import (
	"time"
)

type Wishlist struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"not null;index" json:"user_id"`
	ProductID uint64    `gorm:"not null;index" json:"product_id"`
	VariantID *uint64   `json:"variant_id,omitempty"`
	Notes     *string   `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User    User    `json:"user,omitempty"`
	Product Product `json:"product,omitempty"`
	Variant *ProductVariant `json:"variant,omitempty"`
}

func (Wishlist) TableName() string {
	return "wishlists"
}
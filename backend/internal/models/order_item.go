package models

import (
	"time"
)

type OrderItem struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	OrderID        uint64     `gorm:"not null;index" json:"order_id"`
	ProductID      uint64     `gorm:"not null" json:"product_id"`
	VariantID      *uint64    `json:"variant_id,omitempty"`
	ProductName    string     `gorm:"type:varchar(255);not null" json:"product_name"`
	ProductSKU     string     `gorm:"type:varchar(100);not null" json:"product_sku"`
	VariantName    *string    `gorm:"type:varchar(255)" json:"variant_name,omitempty"`
	ProductImageURL *string   `gorm:"type:text" json:"product_image_url,omitempty"`
	Price          float64    `gorm:"type:decimal(15,2);not null" json:"price"`
	Quantity       int        `gorm:"not null" json:"quantity"`
	Subtotal       float64    `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	WeightGram     int        `gorm:"default:0" json:"weight_gram"`
	DiscountAmount float64    `gorm:"type:decimal(15,2);default:0" json:"discount_amount"`
	CreatedAt      time.Time  `json:"created_at"`

	Order   Order   `json:"-"`
	Product Product `json:"product,omitempty"`
	Variant *ProductVariant `json:"variant,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
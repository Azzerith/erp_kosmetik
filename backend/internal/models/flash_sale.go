package models

import (
	"time"
)

type FlashSale struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	StartTime   time.Time  `gorm:"not null" json:"start_time"`
	EndTime     time.Time  `gorm:"not null" json:"end_time"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedBy   uint64     `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships
	Creator User            `json:"creator,omitempty"`
	Items   []FlashSaleItem `json:"items,omitempty"`
}

func (FlashSale) TableName() string {
	return "flash_sales"
}

type FlashSaleItem struct {
	ID                uint64    `gorm:"primaryKey" json:"id"`
	FlashSaleID       uint64    `gorm:"not null;index" json:"flash_sale_id"`
	ProductID         uint64    `gorm:"not null;index" json:"product_id"`
	VariantID         *uint64   `json:"variant_id,omitempty"`
	FlashPrice        float64   `gorm:"type:decimal(15,2);not null" json:"flash_price"`
	MaxQuantityPerUser int      `gorm:"default:0" json:"max_quantity_per_user"`
	StockAllocated    int       `gorm:"not null" json:"stock_allocated"`
	StockSold         int       `gorm:"default:0" json:"stock_sold"`
	SortOrder         int       `gorm:"default:0" json:"sort_order"`

	// Relationships
	FlashSale FlashSale `json:"flash_sale,omitempty"`
	Product   Product   `json:"product,omitempty"`
	Variant   *ProductVariant `json:"variant,omitempty"`
}

func (FlashSaleItem) TableName() string {
	return "flash_sale_items"
}
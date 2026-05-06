package models

import (
	"time"
)

type InventoryLog struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	ProductID     uint64     `gorm:"not null;index" json:"product_id"`
	VariantID     *uint64    `json:"variant_id,omitempty"`
	Type          string     `gorm:"type:enum('in','out','adjustment','return','damaged');not null" json:"type"`
	Quantity      int        `gorm:"not null" json:"quantity"`
	StockBefore   int        `gorm:"not null" json:"stock_before"`
	StockAfter    int        `gorm:"not null" json:"stock_after"`
	ReferenceType *string    `gorm:"type:varchar(50)" json:"reference_type,omitempty"`
	ReferenceID   *uint64    `json:"reference_id,omitempty"`
	Note          *string    `gorm:"type:text" json:"note,omitempty"`
	CreatedBy     uint64     `gorm:"not null" json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`

	// Relationships
	Product Product `json:"product,omitempty"`
	Variant *ProductVariant `json:"variant,omitempty"`
	Creator User `json:"creator,omitempty"`
}

func (InventoryLog) TableName() string {
	return "inventory_logs"
}
package models

import (
	"time"
)

type ProductVariant struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	ProductID     uint64     `gorm:"not null;index" json:"product_id"`
	SKU           string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`
	VariantType   string     `gorm:"type:varchar(50);not null" json:"variant_type"`
	VariantValue  string     `gorm:"type:varchar(100);not null" json:"variant_value"`
	PriceModifier float64    `gorm:"type:decimal(15,2);default:0" json:"price_modifier"`
	Stock         int        `gorm:"default:0" json:"stock"`
	WeightGram    *int       `json:"weight_gram,omitempty"`
	ImageURL      *string    `gorm:"type:text" json:"image_url,omitempty"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Relationships
	Product Product `json:"product,omitempty"`
}

func (ProductVariant) TableName() string {
	return "product_variants"
}

type ProductImage struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	ProductID    uint64    `gorm:"not null;index" json:"product_id"`
	VariantID    *uint64   `json:"variant_id,omitempty"`
	URL          string    `gorm:"type:varchar(500);not null" json:"url"`
	URLThumbnail *string   `gorm:"type:varchar(500)" json:"url_thumbnail,omitempty"`
	URLMedium    *string   `gorm:"type:varchar(500)" json:"url_medium,omitempty"`
	AltText      *string   `gorm:"type:varchar(255)" json:"alt_text,omitempty"`
	IsPrimary    bool      `gorm:"default:false" json:"is_primary"`
	SortOrder    int       `gorm:"default:0" json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`

	// Relationships
	Product Product `json:"product,omitempty"`
}

func (ProductImage) TableName() string {
	return "product_images"
}

type ProductTag struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	ProductID uint64    `gorm:"not null;index" json:"product_id"`
	Tag       string    `gorm:"type:varchar(100);not null" json:"tag"`
	Weight    float64   `gorm:"type:decimal(3,2);default:1" json:"weight"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Product Product `json:"product,omitempty"`
}

func (ProductTag) TableName() string {
	return "product_tags"
}

type ProductCertification struct {
	ID                uint64     `gorm:"primaryKey" json:"id"`
	ProductID         uint64     `gorm:"not null;index" json:"product_id"`
	CertificationType string     `gorm:"type:enum('bpom','halal','vegan','organic','gmp');not null" json:"certification_type"`
	CertificateNumber *string    `gorm:"type:varchar(100)" json:"certificate_number,omitempty"`
	IssuedBy          *string    `gorm:"type:varchar(255)" json:"issued_by,omitempty"`
	IssuedDate        *time.Time `json:"issued_date,omitempty"`
	ExpiryDate        *time.Time `json:"expiry_date,omitempty"`
	DocumentURL       *string    `gorm:"type:text" json:"document_url,omitempty"`
	IsValid           bool       `gorm:"default:true" json:"is_valid"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relationships
	Product Product `json:"product,omitempty"`
}

func (ProductCertification) TableName() string {
	return "product_certifications"
}
package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	SKU         string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(280);uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"type:longtext" json:"description,omitempty"`
	ShortDesc   *string        `gorm:"type:text" json:"short_description,omitempty"`

	// Foreign keys
	CategoryID  uint64         `gorm:"not null" json:"category_id"`
	BrandID     *uint64        `json:"brand_id,omitempty"`

	// Pricing
	BasePrice     float64       `gorm:"type:decimal(15,2);not null" json:"base_price"`
	SalePrice     *float64      `gorm:"type:decimal(15,2)" json:"sale_price,omitempty"`
	WholesalePrice *float64      `gorm:"type:decimal(15,2)" json:"wholesale_price,omitempty"`

	// Physical
	WeightGram    int           `gorm:"default:0" json:"weight_gram"`

	// Certifications
	IsBPOMCertified bool          `gorm:"default:false" json:"is_bpom_certified"`
	BPOMNumber      *string       `gorm:"type:varchar(50)" json:"bpom_number,omitempty"`
	IsHalalCertified bool         `gorm:"default:false" json:"is_halal_certified"`
	HalalNumber     *string       `gorm:"type:varchar(50)" json:"halal_number,omitempty"`
	IsVegan         bool          `gorm:"default:false" json:"is_vegan"`
	IsHerbal        bool          `gorm:"default:false" json:"is_herbal"`

	// Trend
	TrendScore      float64       `gorm:"type:decimal(5,2);default:0" json:"trend_score"`
	TrendBadge      string        `gorm:"type:enum('trending','viral','best_seller','hot','none');default:'none'" json:"trend_badge"`
	TrendKeywords   datatypes.JSON `gorm:"type:json" json:"trend_keywords,omitempty"`

	// Sales & Stock
	TotalSold        int          `gorm:"default:0" json:"total_sold"`
	Stock            int          `gorm:"default:0" json:"stock"`
	MinStockThreshold int         `gorm:"default:5" json:"min_stock_threshold"`

	// Status
	IsActive        bool          `gorm:"default:true" json:"is_active"`
	IsFeatured      bool          `gorm:"default:false" json:"is_featured"`

	// SEO
	MetaTitle       *string       `gorm:"type:varchar(255)" json:"meta_title,omitempty"`
	MetaDescription *string       `gorm:"type:text" json:"meta_description,omitempty"`

	// Audit
	CreatedBy       *uint64       `json:"created_by,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Category        Category       `json:"category,omitempty"`
	Brand           *Brand         `json:"brand,omitempty"`
	Variants        []ProductVariant `json:"variants,omitempty"`
	Images          []ProductImage   `json:"images,omitempty"`
	Tags            []ProductTag     `json:"tags,omitempty"`
	Certifications  []ProductCertification `json:"certifications,omitempty"`
	Reviews         []Review         `json:"reviews,omitempty"`
	InventoryLogs   []InventoryLog   `json:"-"`
}

func (Product) TableName() string {
	return "products"
}
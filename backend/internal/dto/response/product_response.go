package response

import "time"

type ProductResponse struct {
	ID              uint64              `json:"id"`
	SKU             string              `json:"sku"`
	Name            string              `json:"name"`
	Slug            string              `json:"slug"`
	Description     string              `json:"description,omitempty"`
	ShortDesc       *string             `json:"short_description,omitempty"`
	CategoryID      uint64              `json:"category_id"`
	CategoryName    string              `json:"category_name"`
	BrandID         *uint64             `json:"brand_id,omitempty"`
	BrandName       *string             `json:"brand_name,omitempty"`
	BasePrice       float64             `json:"base_price"`
	SalePrice       *float64            `json:"sale_price,omitempty"`
	WeightGram      int                 `json:"weight_gram"`
	Stock           int                 `json:"stock"`
	TrendScore      float64             `json:"trend_score"`
	TrendBadge      string              `json:"trend_badge"`
	TotalSold       int                 `json:"total_sold"`
	IsActive        bool                `json:"is_active"`
	IsFeatured      bool                `json:"is_featured"`
	Rating          float64             `json:"rating"`
	TotalReviews    int64               `json:"total_reviews"`
	Images          []ProductImageResponse `json:"images,omitempty"`
	Variants        []ProductVariantResponse `json:"variants,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
}

type ProductImageResponse struct {
	ID        uint64  `json:"id"`
	URL       string  `json:"url"`
	IsPrimary bool    `json:"is_primary"`
	SortOrder int     `json:"sort_order"`
}

type ProductVariantResponse struct {
	ID            uint64  `json:"id"`
	VariantType   string  `json:"variant_type"`
	VariantValue  string  `json:"variant_value"`
	PriceModifier float64 `json:"price_modifier"`
	Stock         int     `json:"stock"`
	SKU           string  `json:"sku"`
}

type CategoryResponse struct {
	ID          uint64              `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Description *string             `json:"description,omitempty"`
	IconURL     *string             `json:"icon_url,omitempty"`
	Level       int                 `json:"level"`
	SortOrder   int                 `json:"sort_order"`
	ParentID    *uint64             `json:"parent_id,omitempty"`
	Children    []CategoryResponse  `json:"children,omitempty"`
}
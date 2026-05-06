package request

type CreateProductRequest struct {
	Name         string  `json:"name" binding:"required"`
	SKU          string  `json:"sku" binding:"required"`
	Description  string  `json:"description"`
	ShortDesc    *string `json:"short_description"`
	CategoryID   uint64  `json:"category_id" binding:"required"`
	BrandID      *uint64 `json:"brand_id"`
	BasePrice    float64 `json:"base_price" binding:"required,min=0"`
	SalePrice    *float64 `json:"sale_price"`
	WeightGram   int     `json:"weight_gram"`
	Stock        int     `json:"stock"`
	IsActive     bool    `json:"is_active"`
	IsFeatured   bool    `json:"is_featured"`
}

type UpdateProductRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	ShortDesc    *string  `json:"short_description"`
	CategoryID   *uint64  `json:"category_id"`
	BrandID      *uint64  `json:"brand_id"`
	BasePrice    *float64 `json:"base_price"`
	SalePrice    *float64 `json:"sale_price"`
	WeightGram   *int     `json:"weight_gram"`
	Stock        *int     `json:"stock"`
	IsActive     *bool    `json:"is_active"`
	IsFeatured   *bool    `json:"is_featured"`
}

type ListProductsRequest struct {
	Page       int     `form:"page" binding:"min=1"`
	Limit      int     `form:"limit" binding:"min=1,max=100"`
	CategoryID *uint64 `form:"category_id"`
	BrandID    *uint64 `form:"brand_id"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Search     *string `form:"search"`
	SortBy     string  `form:"sort_by"`
	SortOrder  string  `form:"sort_order"`
}
package request

type AddToCartRequest struct {
	ProductID uint64  `json:"product_id" binding:"required"`
	VariantID *uint64 `json:"variant_id"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}
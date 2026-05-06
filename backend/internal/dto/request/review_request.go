package request

type CreateReviewRequest struct {
	ProductID uint64  `json:"product_id" binding:"required"`
	OrderID   *uint64 `json:"order_id"`
	Rating    int     `json:"rating" binding:"required,min=1,max=5"`
	Title     string  `json:"title"`
	Comment   string  `json:"comment"`
}

type UpdateReviewRequest struct {
	Rating  int     `json:"rating" binding:"min=1,max=5"`
	Title   *string `json:"title"`
	Comment *string `json:"comment"`
}
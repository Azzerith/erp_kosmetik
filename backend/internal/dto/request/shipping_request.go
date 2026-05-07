package request

type CalculateCostRequest struct {
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
	Weight      int    `json:"weight" binding:"required,min=1"`
	Courier     string `json:"courier" binding:"required"`
}
package request

type CreateOrderRequest struct {
	AddressID      uint64  `json:"address_id" binding:"required"`
	Courier        string  `json:"courier" binding:"required"`
	CourierService string  `json:"courier_service" binding:"required"`
	ShippingCost   float64 `json:"shipping_cost" binding:"required,min=0"`
	VoucherCode    *string `json:"voucher_code"`
	Notes          *string `json:"notes"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending_payment paid processing shipped delivered completed cancelled refunded"`
}

type UpdateTrackingRequest struct {
	TrackingNumber string `json:"tracking_number" binding:"required"`
	Courier        string `json:"courier" binding:"required"`
}
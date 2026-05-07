package response

import "time"

type OrderResponse struct {
	ID               uint64              `json:"id"`
	OrderNumber      string              `json:"order_number"`
	Status           string              `json:"status"`
	PaymentStatus    string              `json:"payment_status"`
	FulfillmentStatus string             `json:"fulfillment_status"`
	Subtotal         float64             `json:"subtotal"`
	ShippingCost     float64             `json:"shipping_cost"`
	DiscountAmount   float64             `json:"discount_amount"`
	TotalAmount      float64             `json:"total_amount"`
	Courier          *string             `json:"courier,omitempty"`
	TrackingNumber   *string             `json:"tracking_number,omitempty"`
	PaidAt           *time.Time          `json:"paid_at,omitempty"`
	ShippedAt        *time.Time          `json:"shipped_at,omitempty"`
	DeliveredAt      *time.Time          `json:"delivered_at,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
	Items            []OrderItemResponse `json:"items"`
	ShippingAddress  AddressResponse     `json:"shipping_address"`
}

type OrderItemResponse struct {
	ID           uint64  `json:"id"`
	ProductID    uint64  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductSKU   string  `json:"product_sku"`
	VariantName  *string `json:"variant_name,omitempty"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
	Subtotal     float64 `json:"subtotal"`
}

type AddressResponse struct {
	ID             uint64  `json:"id"`
	Label          string  `json:"label"`
	RecipientName  string  `json:"recipient_name"`
	Phone          string  `json:"phone"`
	Province       string  `json:"province"`
	City           string  `json:"city"`
	District       *string `json:"district,omitempty"`
	PostalCode     string  `json:"postal_code"`
	AddressDetail  string  `json:"address_detail"`
	IsDefault      bool    `json:"is_default"`
}
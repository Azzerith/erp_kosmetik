package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	OrderNumber string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	UserID      uint64    `gorm:"not null;index" json:"user_id"`

	// Status
	Status           string `gorm:"type:enum('pending_payment','paid','processing','shipped','delivered','completed','cancelled','refunded');default:'pending_payment'" json:"status"`
	PaymentStatus    string `gorm:"type:enum('unpaid','pending','paid','failed','refunded');default:'unpaid'" json:"payment_status"`
	FulfillmentStatus string `gorm:"type:enum('unfulfilled','processing','shipped','delivered');default:'unfulfilled'" json:"fulfillment_status"`

	// Financial
	Subtotal       float64 `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	ShippingCost   float64 `gorm:"type:decimal(10,2);default:0" json:"shipping_cost"`
	DiscountAmount float64 `gorm:"type:decimal(15,2);default:0" json:"discount_amount"`
	TaxAmount      float64 `gorm:"type:decimal(15,2);default:0" json:"tax_amount"`
	TotalAmount    float64 `gorm:"type:decimal(15,2);not null" json:"total_amount"`

	// Shipping
	ShippingAddressID    uint64     `gorm:"not null" json:"shipping_address_id"`
	Courier              *string    `gorm:"type:varchar(50)" json:"courier,omitempty"`
	CourierService       *string    `gorm:"type:varchar(100)" json:"courier_service,omitempty"`
	TrackingNumber       *string    `gorm:"type:varchar(100)" json:"tracking_number,omitempty"`
	EstimatedDeliveryStart *time.Time `json:"estimated_delivery_start,omitempty"`
	EstimatedDeliveryEnd   *time.Time `json:"estimated_delivery_end,omitempty"`
	ShippingNote         *string    `gorm:"type:text" json:"shipping_note,omitempty"`

	// Payment
	PaymentMethod  *string `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
	PaymentChannel *string `gorm:"type:varchar(50)" json:"payment_channel,omitempty"`

	// Voucher
	VoucherCode     *string  `gorm:"type:varchar(50)" json:"voucher_code,omitempty"`
	VoucherDiscount *float64 `gorm:"type:decimal(15,2)" json:"voucher_discount,omitempty"`

	// Notes
	Notes         *string `gorm:"type:text" json:"notes,omitempty"`
	InternalNotes *string `gorm:"type:text" json:"internal_notes,omitempty"`

	// Timestamps
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	ShippedAt   *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	CreatedBy   *uint64    `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships
	User          User           `json:"user,omitempty"`
	ShippingAddress Address        `json:"shipping_address,omitempty"`
	Items         []OrderItem     `json:"items,omitempty"`
	StatusHistory []OrderStatusHistory `json:"status_history,omitempty"`
	Payments      []Payment       `json:"payments,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID          uint64  `gorm:"primaryKey" json:"id"`
	OrderID     uint64  `gorm:"not null;index" json:"order_id"`
	ProductID   uint64  `gorm:"not null" json:"product_id"`
	VariantID   *uint64 `json:"variant_id,omitempty"`

	// Snapshot data
	ProductName    string  `gorm:"type:varchar(255);not null" json:"product_name"`
	ProductSKU     string  `gorm:"type:varchar(100);not null" json:"product_sku"`
	VariantName    *string `gorm:"type:varchar(255)" json:"variant_name,omitempty"`
	ProductImageURL *string `gorm:"type:text" json:"product_image_url,omitempty"`

	// Pricing
	Price          float64 `gorm:"type:decimal(15,2);not null" json:"price"`
	Quantity       int     `gorm:"not null" json:"quantity"`
	Subtotal       float64 `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	WeightGram     int     `gorm:"default:0" json:"weight_gram"`
	DiscountAmount float64 `gorm:"type:decimal(15,2);default:0" json:"discount_amount"`

	CreatedAt      time.Time `json:"created_at"`

	// Relationships
	Order   Order   `json:"-"`
	Product Product `json:"product,omitempty"`
	Variant *ProductVariant `json:"variant,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

type OrderStatusHistory struct {
	ID         uint64     `gorm:"primaryKey" json:"id"`
	OrderID    uint64     `gorm:"not null;index" json:"order_id"`
	StatusFrom *string    `gorm:"type:varchar(50)" json:"status_from,omitempty"`
	StatusTo   string     `gorm:"type:varchar(50);not null" json:"status_to"`
	Notes      *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy  *uint64    `json:"created_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`

	// Relationships
	Order Order `json:"-"`
}

func (OrderStatusHistory) TableName() string {
	return "order_status_history"
}
package models

import (
	"time"

	"gorm.io/datatypes"
)

type Payment struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	OrderID       uint64    `gorm:"not null;index" json:"order_id"`
	TransactionID string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"transaction_id"`
	SnapToken     *string   `gorm:"type:varchar(255)" json:"snap_token,omitempty"`
	PaymentMethod string    `gorm:"type:varchar(50);not null" json:"payment_method"`
	PaymentChannel *string  `gorm:"type:varchar(50)" json:"payment_channel,omitempty"`
	Amount        float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	FeeAmount     float64   `gorm:"type:decimal(15,2);default:0" json:"fee_amount"`
	Status        string    `gorm:"type:enum('pending','success','failed','expired','refunded');default:'pending'" json:"status"`
	RawRequest    datatypes.JSON `gorm:"type:json" json:"raw_request,omitempty"`
	RawResponse   datatypes.JSON `gorm:"type:json" json:"raw_response,omitempty"`
	WebhookResponse datatypes.JSON `gorm:"type:json" json:"webhook_response,omitempty"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	ExpiredAt     *time.Time `json:"expired_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Relationships
	Order Order `json:"order,omitempty"`
	Logs  []PaymentLog `json:"logs,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

type PaymentLog struct {
	ID         uint64         `gorm:"primaryKey" json:"id"`
	PaymentID  uint64         `gorm:"not null;index" json:"payment_id"`
	EventType  string         `gorm:"type:varchar(50);not null" json:"event_type"`
	EventData  datatypes.JSON `gorm:"type:json" json:"event_data,omitempty"`
	Signature  *string        `gorm:"type:varchar(255)" json:"signature,omitempty"`
	IPAddress  *string        `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`

	// Relationships
	Payment Payment `json:"-"`
}

func (PaymentLog) TableName() string {
	return "payment_logs"
}
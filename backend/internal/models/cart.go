package models

import (
	"time"
)

type Cart struct {
	ID        uint64     `gorm:"primaryKey" json:"id"`
	UserID    uint64     `gorm:"not null;index" json:"user_id"`
	SessionID *string    `gorm:"type:varchar(255);index" json:"session_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relationships
	User  User       `json:"user,omitempty"`
	Items []CartItem `json:"items,omitempty"`
}

func (Cart) TableName() string {
	return "carts"
}
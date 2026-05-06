package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint64         `gorm:"primaryKey" json:"id"`
	Email           string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash    *string        `gorm:"type:varchar(255)" json:"-"`
	Name            string         `gorm:"type:varchar(100);not null" json:"name"`
	Phone           *string        `gorm:"type:varchar(20)" json:"phone,omitempty"`
	AvatarURL       *string        `gorm:"type:text" json:"avatar_url,omitempty"`
	Provider        string         `gorm:"type:enum('local','google','facebook');default:'local'" json:"provider"`
	ProviderID      *string        `gorm:"type:varchar(255)" json:"-"`
	Role            string         `gorm:"type:enum('super_admin','admin','staff','customer');default:'customer'" json:"role"`
	LoyaltyPoints   int            `gorm:"default:0" json:"loyalty_points"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty"`
	LastLoginIP     *string        `gorm:"type:varchar(45)" json:"-"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	RefreshTokens   []RefreshToken   `json:"-"`
	Addresses       []Address        `json:"addresses,omitempty"`
	Orders          []Order          `json:"orders,omitempty"`
	Reviews         []Review         `json:"reviews,omitempty"`
	Carts           []Cart           `json:"-"`
	Wishlists       []Wishlist       `json:"-"`
	Notifications   []Notification   `json:"-"`
	ActivityLogs    []ActivityLog    `json:"-"`
}

func (User) TableName() string {
	return "users"
}
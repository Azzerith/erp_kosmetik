package models

import (
	"time"
)

type RefreshToken struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsRevoked bool      `gorm:"default:false" json:"is_revoked"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User User `json:"user,omitempty"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

type PasswordReset struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(255);not null;index" json:"email"`
	Token     string    `gorm:"type:varchar(255);not null;index" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time `json:"created_at"`
}

func (PasswordReset) TableName() string {
	return "password_resets"
}
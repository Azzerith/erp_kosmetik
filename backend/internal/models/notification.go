package models

import (
	"time"

	"gorm.io/datatypes"
)

type Notification struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`
	Type      string         `gorm:"type:varchar(50);not null" json:"type"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Data      datatypes.JSON `gorm:"type:json" json:"data,omitempty"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	ReadAt    *time.Time     `json:"read_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`

	// Relationships
	User User `json:"user,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}
package models

import (
	"database/sql"
	"time"
)

type Notification struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	UserID    uint64         `gorm:"not null;index" json:"user_id"`
	Type      string         `gorm:"type:varchar(50);not null" json:"type"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Data      sql.NullJSON   `gorm:"type:json" json:"data,omitempty"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	ReadAt    *time.Time     `json:"read_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`

	// Relationships
	User User `json:"user,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}

type ActivityLog struct {
	ID         uint64         `gorm:"primaryKey" json:"id"`
	UserID     *uint64        `gorm:"index" json:"user_id,omitempty"`
	Action     string         `gorm:"type:varchar(100);not null" json:"action"`
	EntityType *string        `gorm:"type:varchar(50)" json:"entity_type,omitempty"`
	EntityID   *uint64        `json:"entity_id,omitempty"`
	OldValues  sql.NullJSON   `gorm:"type:json" json:"old_values,omitempty"`
	NewValues  sql.NullJSON   `gorm:"type:json" json:"new_values,omitempty"`
	IPAddress  *string        `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent  *string        `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`

	// Relationships
	User *User `json:"user,omitempty"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}
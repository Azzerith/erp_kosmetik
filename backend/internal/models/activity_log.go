package models

import (
	"time"

	"gorm.io/datatypes"
)

type ActivityLog struct {
	ID         uint64         `gorm:"primaryKey" json:"id"`
	UserID     *uint64        `gorm:"index" json:"user_id,omitempty"`
	Action     string         `gorm:"type:varchar(100);not null" json:"action"`
	EntityType *string        `gorm:"type:varchar(50)" json:"entity_type,omitempty"`
	EntityID   *uint64        `json:"entity_id,omitempty"`
	OldValues  datatypes.JSON `gorm:"type:json" json:"old_values,omitempty"`
	NewValues  datatypes.JSON `gorm:"type:json" json:"new_values,omitempty"`
	IPAddress  *string        `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent  *string        `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`

	User *User `json:"user,omitempty"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}
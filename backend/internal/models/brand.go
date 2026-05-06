package models

import (
	"time"

	"gorm.io/gorm"
)

type Brand struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Slug        string         `gorm:"type:varchar(120);uniqueIndex;not null" json:"slug"`
	LogoURL     *string        `gorm:"type:text" json:"logo_url,omitempty"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Products []Product `json:"products,omitempty"`
}

func (Brand) TableName() string {
	return "brands"
}
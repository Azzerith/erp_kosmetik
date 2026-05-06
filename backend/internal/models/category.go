package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(120);uniqueIndex;not null" json:"slug"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	ParentID    *uint64        `json:"parent_id,omitempty"`
	IconURL     *string        `gorm:"type:text" json:"icon_url,omitempty"`
	ImageURL    *string        `gorm:"type:text" json:"image_url,omitempty"`
	Level       int            `gorm:"default:0" json:"level"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Parent   *Category  `json:"parent,omitempty"`
	Children []Category `json:"children,omitempty"`
	Products []Product  `json:"products,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
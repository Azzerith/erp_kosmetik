package models

import (
	"time"

	"gorm.io/datatypes"
)

type TrendData struct {
	ID              uint64         `gorm:"primaryKey" json:"id"`
	Keyword         string         `gorm:"type:varchar(255);not null;index" json:"keyword"`
	Source          string         `gorm:"type:enum('google','tiktok','instagram','internal');not null" json:"source"`
	Score           float64        `gorm:"type:decimal(10,4);not null" json:"score"`
	SearchVolume    *int           `json:"search_volume,omitempty"`
	GrowthPercentage *float64      `gorm:"type:decimal(8,2)" json:"growth_percentage,omitempty"`
	RelatedKeywords datatypes.JSON `gorm:"type:json" json:"related_keywords,omitempty"`
	Region          string         `gorm:"type:varchar(10);default:'ID'" json:"region"`
	Language        string         `gorm:"type:varchar(10);default:'id'" json:"language"`
	PeriodStart     time.Time      `gorm:"not null" json:"period_start"`
	PeriodEnd       time.Time      `gorm:"not null" json:"period_end"`
	RecordedAt      time.Time      `gorm:"autoCreateTime" json:"recorded_at"`

	// Relationships
	ProductMappings []ProductTrendMapping `json:"-"`
}

func (TrendData) TableName() string {
	return "trend_data"
}
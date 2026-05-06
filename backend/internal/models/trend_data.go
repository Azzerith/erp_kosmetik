package models

import (
	"database/sql"
	"time"
)

type TrendData struct {
	ID              uint64         `gorm:"primaryKey" json:"id"`
	Keyword         string         `gorm:"type:varchar(255);not null;index" json:"keyword"`
	Source          string         `gorm:"type:enum('google','tiktok','instagram','internal');not null" json:"source"`
	Score           float64        `gorm:"type:decimal(10,4);not null" json:"score"`
	SearchVolume    *int           `json:"search_volume,omitempty"`
	GrowthPercentage *float64      `gorm:"type:decimal(8,2)" json:"growth_percentage,omitempty"`
	RelatedKeywords sql.NullJSON   `gorm:"type:json" json:"related_keywords,omitempty"`
	Region          string         `gorm:"type:varchar(10);default:'ID'" json:"region"`
	Language        string         `gorm:"type:varchar(10);default:'id'" json:"language"`
	PeriodStart     time.Time      `gorm:"not null" json:"period_start"`
	PeriodEnd       time.Time      `gorm:"not null" json:"period_end"`
	RecordedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"recorded_at"`

	// Relationships
	ProductMappings []ProductTrendMapping `json:"-"`
}

func (TrendData) TableName() string {
	return "trend_data"
}

type ProductTrendMapping struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	ProductID      uint64    `gorm:"not null;index" json:"product_id"`
	TrendDataID    uint64    `gorm:"not null;index" json:"trend_data_id"`
	RelevanceScore float64   `gorm:"type:decimal(5,2);default:1" json:"relevance_score"`
	CreatedAt      time.Time `json:"created_at"`

	// Relationships
	Product   Product   `json:"product,omitempty"`
	TrendData TrendData `json:"trend_data,omitempty"`
}

func (ProductTrendMapping) TableName() string {
	return "product_trend_mapping"
}

type TrendScoreHistory struct {
	ID                 uint64    `gorm:"primaryKey" json:"id"`
	ProductID          uint64    `gorm:"not null;index" json:"product_id"`
	TrendScore         float64   `gorm:"type:decimal(5,2);not null" json:"trend_score"`
	TrendBadge         *string   `gorm:"type:varchar(20)" json:"trend_badge,omitempty"`
	GoogleScore        *float64  `gorm:"type:decimal(5,2)" json:"google_score,omitempty"`
	TikTokScore        *float64  `gorm:"type:decimal(5,2)" json:"tiktok_score,omitempty"`
	SalesVelocityScore *float64  `gorm:"type:decimal(5,2)" json:"sales_velocity_score,omitempty"`
	RecordedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"recorded_at"`

	// Relationships
	Product Product `json:"product,omitempty"`
}

func (TrendScoreHistory) TableName() string {
	return "trend_scores_history"
}
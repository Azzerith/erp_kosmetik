package models

import (
	"time"
)

type ProductTrendMapping struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	ProductID      uint64    `gorm:"not null;index" json:"product_id"`
	TrendDataID    uint64    `gorm:"not null;index" json:"trend_data_id"`
	RelevanceScore float64   `gorm:"type:decimal(5,2);default:1" json:"relevance_score"`
	CreatedAt      time.Time `json:"created_at"`

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
	RecordedAt         time.Time `gorm:"autoCreateTime" json:"recorded_at"`

	Product Product `json:"product,omitempty"`
}

func (TrendScoreHistory) TableName() string {
	return "trend_scores_history"
}
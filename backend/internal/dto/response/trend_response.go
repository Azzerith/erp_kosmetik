package response

import "time"

type TrendingKeywordResponse struct {
	Keyword      string  `json:"keyword"`
	Score        float64 `json:"score"`
	Trend        string  `json:"trend"`
	GrowthRate   float64 `json:"growth_rate"`
}

type TrendScoreResponse struct {
	ProductID     uint64    `json:"product_id"`
	TrendScore    float64   `json:"trend_score"`
	TrendBadge    string    `json:"trend_badge"`
	GoogleScore   *float64  `json:"google_score,omitempty"`
	TikTokScore   *float64  `json:"tiktok_score,omitempty"`
	SalesVelocity *float64  `json:"sales_velocity,omitempty"`
	LastUpdated   time.Time `json:"last_updated"`
}

type TrendDashboardResponse struct {
	TopKeywords      []TrendingKeywordResponse `json:"top_keywords"`
	TopProducts      []ProductResponse         `json:"top_products"`
	ScoreHistory     []ScoreHistoryPoint       `json:"score_history"`
	Summary          TrendSummaryResponse      `json:"summary"`
}

type ScoreHistoryPoint struct {
	Date  string  `json:"date"`
	Score float64 `json:"score"`
}

type TrendSummaryResponse struct {
	TotalProductsTracked int     `json:"total_products_tracked"`
	AverageTrendScore    float64 `json:"average_trend_score"`
	HighestTrendScore    float64 `json:"highest_trend_score"`
	TrendingUpCount      int     `json:"trending_up_count"`
	TrendingDownCount    int     `json:"trending_down_count"`
}
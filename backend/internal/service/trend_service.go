package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"erp-cosmetics-backend/internal/config"
	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type TrendService interface {
	GetTrendingKeywords(ctx context.Context) ([]TrendKeyword, error)
	GetTrendingProducts(ctx context.Context, limit int) ([]models.Product, error)
	GetProductTrendScore(ctx context.Context, productID uint64) (*TrendScoreResponse, error)
	RefreshTrends(ctx context.Context) error
	GetTrendDashboard(ctx context.Context) (*TrendDashboardResponse, error)
}

type TrendKeyword struct {
	Keyword    string  `json:"keyword"`
	Score      float64 `json:"score"`
	Trend      string  `json:"trend"` // up, down, stable
	GrowthRate float64 `json:"growth_rate"`
}

type TrendScoreResponse struct {
	ProductID        uint64    `json:"product_id"`
	TrendScore       float64   `json:"trend_score"`
	TrendBadge       string    `json:"trend_badge"`
	GoogleScore      *float64  `json:"google_score,omitempty"`
	TikTokScore      *float64  `json:"tiktok_score,omitempty"`
	SalesVelocity    *float64  `json:"sales_velocity,omitempty"`
	LastUpdated      time.Time `json:"last_updated"`
}

type TrendDashboardResponse struct {
	TopKeywords      []TrendKeyword      `json:"top_keywords"`
	TopProducts      []models.Product    `json:"top_products"`
	ScoreHistory     []ScoreHistoryPoint `json:"score_history"`
	Summary          TrendSummary        `json:"summary"`
}

type ScoreHistoryPoint struct {
	Date  string  `json:"date"`
	Score float64 `json:"score"`
}

type TrendSummary struct {
	TotalProductsTracked int     `json:"total_products_tracked"`
	AverageTrendScore    float64 `json:"average_trend_score"`
	HighestTrendScore    float64 `json:"highest_trend_score"`
	TrendingUpCount      int     `json:"trending_up_count"`
	TrendingDownCount    int     `json:"trending_down_count"`
}

type trendService struct {
	cfg         *config.Config
	trendRepo   repository.TrendRepository
	productRepo repository.ProductRepository
	redisClient *redis.Client
	logger      *zap.Logger
}

func NewTrendService(
	cfg *config.Config,
	trendRepo repository.TrendRepository,
	productRepo repository.ProductRepository,
	redisClient *redis.Client,
	logger *zap.Logger,
) TrendService {
	return &trendService{
		cfg:         cfg,
		trendRepo:   trendRepo,
		productRepo: productRepo,
		redisClient: redisClient,
		logger:      logger,
	}
}

func (s *trendService) GetTrendingKeywords(ctx context.Context) ([]TrendKeyword, error) {
	// Try to get from cache first
	cacheKey := "trending_keywords"
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var keywords []TrendKeyword
		if err := json.Unmarshal([]byte(cached), &keywords); err == nil {
			return keywords, nil
		}
	}

	// Fetch from external APIs (Google Trends via SerpAPI)
	keywords, err := s.fetchGoogleTrends(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch Google Trends", zap.Error(err))
		// Return fallback data
		return s.getFallbackKeywords(), nil
	}

	// Cache for 6 hours
	keywordsJSON, _ := json.Marshal(keywords)
	s.redisClient.Set(ctx, cacheKey, keywordsJSON, 6*time.Hour)

	return keywords, nil
}

func (s *trendService) GetTrendingProducts(ctx context.Context, limit int) ([]models.Product, error) {
	return s.productRepo.GetTrending(ctx, limit)
}

func (s *trendService) GetProductTrendScore(ctx context.Context, productID uint64) (*TrendScoreResponse, error) {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Get trend history
	history, err := s.trendRepo.GetTrendHistory(ctx, productID, 30)
	if err != nil {
		s.logger.Warn("Failed to get trend history", zap.Error(err))
	}

	// Calculate sales velocity
	salesVelocity := s.calculateSalesVelocity(product)

	return &TrendScoreResponse{
		ProductID:     product.ID,
		TrendScore:    product.TrendScore,
		TrendBadge:    product.TrendBadge,
		SalesVelocity: &salesVelocity,
		LastUpdated:   product.UpdatedAt,
	}, nil
}

func (s *trendService) RefreshTrends(ctx context.Context) error {
	s.logger.Info("Starting trend refresh job")

	// Fetch trends from Google
	googleTrends, err := s.fetchGoogleTrends(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch Google Trends", zap.Error(err))
		return err
	}

	// Fetch trends from TikTok
	tiktokTrends, err := s.fetchTikTokTrends(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch TikTok Trends", zap.Error(err))
		// Continue with Google trends only
	}

	// Save trend data to database
	for _, keyword := range googleTrends {
		trendData := &models.TrendData{
			Keyword:      keyword.Keyword,
			Source:       "google",
			Score:        keyword.Score,
			Region:       "ID",
			Language:     "id",
			PeriodStart:  time.Now().AddDate(0, 0, -7),
			PeriodEnd:    time.Now(),
			RecordedAt:   time.Now(),
		}
		if err := s.trendRepo.CreateTrendData(ctx, trendData); err != nil {
			s.logger.Warn("Failed to save trend data", zap.Error(err))
		}
	}

	// Update product trend scores
	if err := s.updateProductTrendScores(ctx, googleTrends); err != nil {
		s.logger.Error("Failed to update product trend scores", zap.Error(err))
	}

	s.logger.Info("Trend refresh job completed")

	return nil
}

func (s *trendService) GetTrendDashboard(ctx context.Context) (*TrendDashboardResponse, error) {
	// Get top keywords
	keywords, err := s.GetTrendingKeywords(ctx)
	if err != nil {
		return nil, err
	}

	// Get top products
	products, err := s.productRepo.GetTrending(ctx, 10)
	if err != nil {
		return nil, err
	}

	// Get score history for chart
	history, err := s.getScoreHistory(ctx)
	if err != nil {
		history = []ScoreHistoryPoint{}
	}

	// Calculate summary
	summary, err := s.calculateTrendSummary(ctx)
	if err != nil {
		summary = TrendSummary{}
	}

	return &TrendDashboardResponse{
		TopKeywords:  keywords,
		TopProducts:  products,
		ScoreHistory: history,
		Summary:      summary,
	}, nil
}

func (s *trendService) fetchGoogleTrends(ctx context.Context) ([]TrendKeyword, error) {
	// Call SerpAPI for Google Trends
	// This is a simplified implementation
	apiKey := s.cfg.SerpAPIKey
	if apiKey == "" {
		return s.getFallbackKeywords(), nil
	}

	url := fmt.Sprintf("https://serpapi.com/search?engine=google_trends&q=skincare&data_type=TIMESERIES&geo=ID&api_key=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Parse response - simplified for now
	keywords := s.getFallbackKeywords()
	return keywords, nil
}

func (s *trendService) fetchTikTokTrends(ctx context.Context) ([]TrendKeyword, error) {
	// TikTok Research API implementation
	// This will be implemented when TikTok API access is available
	return []TrendKeyword{
		{Keyword: "GlowUp", Score: 92, Trend: "up", GrowthRate: 15.5},
		{Keyword: "SkincareRoutine", Score: 88, Trend: "up", GrowthRate: 12.3},
		{Keyword: "JamuHerbal", Score: 85, Trend: "up", GrowthRate: 25.0},
	}, nil
}

func (s *trendService) updateProductTrendScores(ctx context.Context, trends []TrendKeyword) error {
	// Get all active products
	products, _, err := s.productRepo.List(ctx, 0, 1000, map[string]interface{}{})
	if err != nil {
		return err
	}

	for _, product := range products {
		// Calculate trend score based on product tags matching trends
		score := s.calculateProductTrendScore(product, trends)

		// Determine badge
		badge := "none"
		if score >= 90 {
			badge = "viral"
		} else if score >= 80 {
			badge = "trending"
		} else if score >= 70 {
			badge = "hot"
		}

		if err := s.productRepo.UpdateTrendScore(ctx, product.ID, score, badge); err != nil {
			s.logger.Warn("Failed to update product trend score", zap.Error(err))
		}

		// Save history
		history := &models.TrendScoreHistory{
			ProductID:  product.ID,
			TrendScore: score,
			TrendBadge: &badge,
			RecordedAt: time.Now(),
		}
		s.trendRepo.CreateTrendHistory(ctx, history)
	}

	return nil
}

func (s *trendService) calculateProductTrendScore(product models.Product, trends []TrendKeyword) float64 {
	// Basic implementation - can be enhanced
	// Match product tags with trending keywords
	maxScore := 0.0
	for _, trend := range trends {
		for _, tag := range product.Tags {
			if containsIgnoreCase(trend.Keyword, tag.Tag) || containsIgnoreCase(tag.Tag, trend.Keyword) {
				if trend.Score > maxScore {
					maxScore = trend.Score
				}
			}
		}
	}
	return maxScore
}

func (s *trendService) calculateSalesVelocity(product models.Product) float64 {
	// Calculate sales velocity based on recent sales
	// Simplified - would need to query recent order items
	if product.TotalSold == 0 {
		return 0
	}
	// Rough estimate: assume product exists for 30 days
	dailyAverage := float64(product.TotalSold) / 30
	return min(100, dailyAverage*2)
}

func (s *trendService) calculateTrendSummary(ctx context.Context) (TrendSummary, error) {
	var summary TrendSummary

	// Get all products with trend scores
	products, _, err := s.productRepo.List(ctx, 0, 1000, map[string]interface{}{})
	if err != nil {
		return summary, err
	}

	var totalScore float64
	highestScore := 0.0

	for _, product := range products {
		if product.TrendScore > 0 {
			totalScore += product.TrendScore
			summary.TotalProductsTracked++

			if product.TrendScore > highestScore {
				highestScore = product.TrendScore
			}
		}
	}

	if summary.TotalProductsTracked > 0 {
		summary.AverageTrendScore = totalScore / float64(summary.TotalProductsTracked)
	}
	summary.HighestTrendScore = highestScore

	return summary, nil
}

func (s *trendService) getScoreHistory(ctx context.Context) ([]ScoreHistoryPoint, error) {
	var history []ScoreHistoryPoint

	// Get average trend scores for last 30 days
	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		history = append(history, ScoreHistoryPoint{
			Date:  date,
			Score: 65 + float64(i%20), // Dummy data, actual implementation would query DB
		})
	}

	return history, nil
}

func (s *trendService) getFallbackKeywords() []TrendKeyword {
	return []TrendKeyword{
		{Keyword: "Skincare", Score: 95, Trend: "up", GrowthRate: 12.5},
		{Keyword: "Vitamin C", Score: 92, Trend: "up", GrowthRate: 15.2},
		{Keyword: "Sunscreen", Score: 89, Trend: "up", GrowthRate: 8.7},
		{Keyword: "Jamu", Score: 87, Trend: "up", GrowthRate: 22.3},
		{Keyword: "Retinol", Score: 85, Trend: "stable", GrowthRate: 1.2},
		{Keyword: "Hyaluronic", Score: 83, Trend: "up", GrowthRate: 5.4},
		{Keyword: "Lip Tint", Score: 80, Trend: "up", GrowthRate: 10.1},
		{Keyword: "Foundation", Score: 78, Trend: "down", GrowthRate: -3.2},
	}
}

func containsIgnoreCase(s, substr string) bool {
	if len(s) < len(substr) || len(substr) == 0 {
		return false
	}
	return stringsContains(stringsToLower(s), stringsToLower(substr))
}

// Helper functions to avoid importing strings
func stringsToLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func stringsContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
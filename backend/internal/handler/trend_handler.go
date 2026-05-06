package handler

import (
	"net/http"
	"strconv"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TrendHandler struct {
	trendService service.TrendService
	logger       *zap.Logger
}

func NewTrendHandler(trendService service.TrendService, logger *zap.Logger) *TrendHandler {
	return &TrendHandler{
		trendService: trendService,
		logger:       logger,
	}
}

func (h *TrendHandler) GetTrendingKeywords(c *gin.Context) {
	keywords, err := h.trendService.GetTrendingKeywords(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch trending keywords", err)
		return
	}

	utils.SuccessResponse(c, keywords)
}

func (h *TrendHandler) GetTrendingProducts(c *gin.Context) {
	limit := 8
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	products, err := h.trendService.GetTrendingProducts(c.Request.Context(), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch trending products", err)
		return
	}

	utils.SuccessResponse(c, products)
}

func (h *TrendHandler) GetTrendScore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	score, err := h.trendService.GetProductTrendScore(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Product not found", err)
		return
	}

	utils.SuccessResponse(c, score)
}

// Admin methods
func (h *TrendHandler) GetDashboard(c *gin.Context) {
	dashboard, err := h.trendService.GetTrendDashboard(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch trend dashboard", err)
		return
	}

	utils.SuccessResponse(c, dashboard)
}

func (h *TrendHandler) RefreshTrends(c *gin.Context) {
	if err := h.trendService.RefreshTrends(c.Request.Context()); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to refresh trends", err)
		return
	}

	utils.SuccessWithMessage(c, "Trends refreshed successfully", nil)
}
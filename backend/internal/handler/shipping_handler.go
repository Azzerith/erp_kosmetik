package handler

import (
	"net/http"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ShippingHandler struct {
	shippingService service.ShippingService
	logger          *zap.Logger
}

func NewShippingHandler(shippingService service.ShippingService, logger *zap.Logger) *ShippingHandler {
	return &ShippingHandler{
		shippingService: shippingService,
		logger:          logger,
	}
}

type CalculateCostRequest struct {
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
	Weight      int    `json:"weight" binding:"required,min=1"`
	Courier     string `json:"courier" binding:"required"`
}

func (h *ShippingHandler) GetProvinces(c *gin.Context) {
	provinces, err := h.shippingService.GetProvinces(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch provinces", err)
		return
	}
	
	utils.SuccessResponse(c, provinces)
}

func (h *ShippingHandler) GetCities(c *gin.Context) {
	provinceID := c.Param("provinceId")
	if provinceID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Province ID is required", nil)
		return
	}
	
	cities, err := h.shippingService.GetCities(c.Request.Context(), provinceID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch cities", err)
		return
	}
	
	utils.SuccessResponse(c, cities)
}

func (h *ShippingHandler) CalculateCost(c *gin.Context) {
	var req CalculateCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	
	costReq := &service.CostRequest{
		Origin:      req.Origin,
		Destination: req.Destination,
		Weight:      req.Weight,
		Courier:     req.Courier,
	}
	
	result, err := h.shippingService.CalculateCost(c.Request.Context(), costReq)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to calculate cost", err)
		return
	}
	
	utils.SuccessResponse(c, result)
}
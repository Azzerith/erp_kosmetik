package handler

import (
	"net/http"
	"strconv"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InventoryHandler struct {
	inventoryService service.InventoryService
	logger           *zap.Logger
}

func NewInventoryHandler(inventoryService service.InventoryService, logger *zap.Logger) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
		logger:           logger,
	}
}

type AdjustStockRequest struct {
	ProductID  uint64  `json:"product_id" binding:"required"`
	VariantID  *uint64 `json:"variant_id"`
	Quantity   int     `json:"quantity" binding:"required"`
	Type       string  `json:"type" binding:"required,oneof=in out adjustment"`
	Note       string  `json:"note"`
}

func (h *InventoryHandler) GetInventory(c *gin.Context) {
	_ = 1 // page
	_ = 20 // limit
	
	// Get inventory list
	// This would come from inventory service
	
	utils.SuccessResponse(c, nil)
}

func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	userID := c.GetUint64("user_id")
	
	var req AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	
	if err := h.inventoryService.AdjustStock(c.Request.Context(), req.ProductID, req.VariantID, req.Quantity, req.Type, req.Note, userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to adjust stock", err)
		return
	}
	
	utils.SuccessWithMessage(c, "Stock adjusted successfully", nil)
}

func (h *InventoryHandler) GetInventoryLogs(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Query("product_id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}
	
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	
	logs, total, err := h.inventoryService.GetInventoryLogs(c.Request.Context(), productID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch inventory logs", err)
		return
	}
	
	utils.PaginatedResponse(c, logs, page, limit, total)
}
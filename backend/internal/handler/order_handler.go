package handler

import (
	"net/http"
	"strconv"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderHandler struct {
	orderService service.OrderService
	logger       *zap.Logger
}

func NewOrderHandler(orderService service.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		logger:       logger,
	}
}

type CreateOrderRequest struct {
	AddressID      uint64  `json:"address_id" binding:"required"`
	Courier        string  `json:"courier" binding:"required"`
	CourierService string  `json:"courier_service" binding:"required"`
	ShippingCost   float64 `json:"shipping_cost" binding:"required,min=0"`
	VoucherCode    *string `json:"voucher_code"`
	Notes          *string `json:"notes"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	serviceReq := &service.CreateOrderRequest{
		AddressID:      req.AddressID,
		Courier:        req.Courier,
		CourierService: req.CourierService,
		ShippingCost:   req.ShippingCost,
		VoucherCode:    req.VoucherCode,
		Notes:          req.Notes,
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), userID, serviceReq)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create order", err)
		return
	}

	utils.CreatedResponse(c, order)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := c.GetUint64("user_id")

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	orders, err := h.orderService.GetUserOrders(c.Request.Context(), userID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch orders", err)
		return
	}

	utils.SuccessResponse(c, orders)
}

func (h *OrderHandler) GetOrderByNumber(c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	if orderNumber == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Order number is required", nil)
		return
	}

	order, err := h.orderService.GetOrderByNumber(c.Request.Context(), orderNumber)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Order not found", err)
		return
	}

	utils.SuccessResponse(c, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	if err := h.orderService.CancelOrder(c.Request.Context(), id, userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to cancel order", err)
		return
	}

	utils.SuccessWithMessage(c, "Order cancelled successfully", nil)
}

// Admin methods
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	var req service.AdminListOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	resp, err := h.orderService.GetAllOrders(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch orders", err)
		return
	}

	utils.PaginatedResponse(c, resp.Orders, resp.Page, resp.Limit, resp.Total)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.orderService.UpdateOrderStatus(c.Request.Context(), id, req.Status); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update order status", err)
		return
	}

	utils.SuccessWithMessage(c, "Order status updated", nil)
}

func (h *OrderHandler) UpdateTracking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	var req struct {
		TrackingNumber string `json:"tracking_number" binding:"required"`
		Courier        string `json:"courier" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.orderService.UpdateTracking(c.Request.Context(), id, req.TrackingNumber, req.Courier); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update tracking", err)
		return
	}

	utils.SuccessWithMessage(c, "Tracking updated", nil)
}
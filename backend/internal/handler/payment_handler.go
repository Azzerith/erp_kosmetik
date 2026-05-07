package handler

import (
	"net/http"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	paymentService service.PaymentService
	logger         *zap.Logger
}

func NewPaymentHandler(paymentService service.PaymentService, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		logger:         logger,
	}
}

type InitiatePaymentRequest struct {
	OrderID uint64 `json:"order_id" binding:"required"`
}

func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	var req InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	snapToken, err := h.paymentService.InitiatePayment(c.Request.Context(), req.OrderID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to initiate payment", err)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"snap_token": snapToken,
	})
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	orderID := c.Param("orderId")

	status, err := h.paymentService.GetPaymentStatus(c.Request.Context(), orderID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Payment not found", err)
		return
	}

	utils.SuccessResponse(c, status)
}
package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WebhookHandler struct {
	paymentService service.PaymentService
	logger         *zap.Logger
}

func NewWebhookHandler(paymentService service.PaymentService, logger *zap.Logger) *WebhookHandler {
	return &WebhookHandler{
		paymentService: paymentService,
		logger:         logger,
	}
}

func (h *WebhookHandler) HandleMidtransWebhook(c *gin.Context) {
	// Read body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to read body", err)
		return
	}

	// Verify signature (implement signature verification here)
	// For now, just parse and process

	var notification map[string]interface{}
	if err := json.Unmarshal(body, &notification); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	// Process webhook
	if err := h.paymentService.ProcessWebhook(c.Request.Context(), notification); err != nil {
		h.logger.Error("Failed to process webhook", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process webhook", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
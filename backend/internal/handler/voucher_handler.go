package handler

import (
	"net/http"
	"strconv"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VoucherHandler struct {
	voucherService service.VoucherService
	logger         *zap.Logger
}

func NewVoucherHandler(voucherService service.VoucherService, logger *zap.Logger) *VoucherHandler {
	return &VoucherHandler{
		voucherService: voucherService,
		logger:         logger,
	}
}

type ValidateVoucherRequest struct {
	Code        string  `json:"code" binding:"required"`
	OrderAmount float64 `json:"order_amount" binding:"required,min=0"`
}

func (h *VoucherHandler) ValidateVoucher(c *gin.Context) {
	var req ValidateVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	
	result, err := h.voucherService.ValidateVoucher(c.Request.Context(), req.Code, req.OrderAmount)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to validate voucher", err)
		return
	}
	
	utils.SuccessResponse(c, result)
}

// Admin methods
func (h *VoucherHandler) CreateVoucher(c *gin.Context) {
	var req service.CreateVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	
	voucher, err := h.voucherService.CreateVoucher(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create voucher", err)
		return
	}
	
	utils.CreatedResponse(c, voucher)
}

func (h *VoucherHandler) UpdateVoucher(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid voucher ID", err)
		return
	}
	
	var req service.UpdateVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	
	voucher, err := h.voucherService.UpdateVoucher(c.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update voucher", err)
		return
	}
	
	utils.SuccessResponse(c, voucher)
}

func (h *VoucherHandler) DeleteVoucher(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid voucher ID", err)
		return
	}
	
	if err := h.voucherService.DeleteVoucher(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete voucher", err)
		return
	}
	
	utils.SuccessWithMessage(c, "Voucher deleted", nil)
}
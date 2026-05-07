package handler

import (
	"net/http"
	"strconv"

	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService service.UserService
	logger      *zap.Logger
}

func NewUserHandler(userService service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

type CreateAddressRequest struct {
	Label         string `json:"label" binding:"required"`
	RecipientName string `json:"recipient_name" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
	Province      string `json:"province" binding:"required"`
	City          string `json:"city" binding:"required"`
	District      string `json:"district"`
	PostalCode    string `json:"postal_code" binding:"required"`
	AddressDetail string `json:"address_detail" binding:"required"`
	IsDefault     bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	Label         *string `json:"label"`
	RecipientName *string `json:"recipient_name"`
	Phone         *string `json:"phone"`
	Province      *string `json:"province"`
	City          *string `json:"city"`
	District      *string `json:"district"`
	PostalCode    *string `json:"postal_code"`
	AddressDetail *string `json:"address_detail"`
	IsDefault     *bool   `json:"is_default"`
}

func (h *UserHandler) GetAddresses(c *gin.Context) {
	userID := c.GetUint64("user_id")
	
	addresses, err := h.userService.GetAddresses(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch addresses", err)
		return
	}
	
	utils.SuccessResponse(c, addresses)
}

func (h *UserHandler) CreateAddress(c *gin.Context) {
	userID := c.GetUint64("user_id")
	
	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	
	address := &models.Address{
		UserID:        userID,
		Label:         req.Label,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		Province:      req.Province,
		City:          req.City,
		District:      &req.District,
		PostalCode:    req.PostalCode,
		AddressDetail: req.AddressDetail,
		IsDefault:     req.IsDefault,
	}
	
	if err := h.userService.CreateAddress(c.Request.Context(), address); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create address", err)
		return
	}
	
	utils.CreatedResponse(c, address)
}

func (h *UserHandler) UpdateAddress(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid address ID", err)
		return
	}
	
	var req UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	
	// Get existing address
	addresses, err := h.userService.GetAddresses(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get address", err)
		return
	}
	
	var address *models.Address
	for _, addr := range addresses {
		if addr.ID == id {
			address = &addr
			break
		}
	}
	
	if address == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Address not found", nil)
		return
	}
	
	if req.Label != nil {
		address.Label = *req.Label
	}
	if req.RecipientName != nil {
		address.RecipientName = *req.RecipientName
	}
	if req.Phone != nil {
		address.Phone = *req.Phone
	}
	if req.Province != nil {
		address.Province = *req.Province
	}
	if req.City != nil {
		address.City = *req.City
	}
	if req.District != nil {
		address.District = req.District
	}
	if req.PostalCode != nil {
		address.PostalCode = *req.PostalCode
	}
	if req.AddressDetail != nil {
		address.AddressDetail = *req.AddressDetail
	}
	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
	}
	
	if err := h.userService.UpdateAddress(c.Request.Context(), address); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update address", err)
		return
	}
	
	utils.SuccessResponse(c, address)
}

func (h *UserHandler) DeleteAddress(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid address ID", err)
		return
	}
	
	// Verify address belongs to user
	addresses, err := h.userService.GetAddresses(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get address", err)
		return
	}
	
	found := false
	for _, addr := range addresses {
		if addr.ID == id {
			found = true
			break
		}
	}
	
	if !found {
		utils.ErrorResponse(c, http.StatusNotFound, "Address not found", nil)
		return
	}
	
	if err := h.userService.DeleteAddress(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete address", err)
		return
	}
	
	utils.SuccessWithMessage(c, "Address deleted", nil)
}
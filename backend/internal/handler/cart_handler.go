package handler

import (
	"net/http"
	"strconv"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CartHandler struct {
	cartService service.CartService
	logger      *zap.Logger
}

func NewCartHandler(cartService service.CartService, logger *zap.Logger) *CartHandler {
	return &CartHandler{
		cartService: cartService,
		logger:      logger,
	}
}

type AddToCartRequest struct {
	ProductID uint64  `json:"product_id" binding:"required"`
	VariantID *uint64 `json:"variant_id"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetUint64("user_id")

	cart, err := h.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch cart", err)
		return
	}

	utils.SuccessResponse(c, cart)
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.cartService.AddToCart(c.Request.Context(), userID, req.ProductID, req.VariantID, req.Quantity); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to add to cart", err)
		return
	}

	utils.SuccessWithMessage(c, "Item added to cart", nil)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID := c.GetUint64("user_id")
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID", err)
		return
	}

	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.cartService.UpdateCartItem(c.Request.Context(), userID, itemID, req.Quantity); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update cart item", err)
		return
	}

	utils.SuccessWithMessage(c, "Cart updated", nil)
}

func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	userID := c.GetUint64("user_id")
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID", err)
		return
	}

	if err := h.cartService.RemoveCartItem(c.Request.Context(), userID, itemID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to remove cart item", err)
		return
	}

	utils.SuccessWithMessage(c, "Item removed from cart", nil)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.GetUint64("user_id")

	if err := h.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to clear cart", err)
		return
	}

	utils.SuccessWithMessage(c, "Cart cleared", nil)
}
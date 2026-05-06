package handler

import (
	"net/http"
	"strconv"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProductHandler struct {
	productService service.ProductService
	logger         *zap.Logger
}

func NewProductHandler(productService service.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		logger:         logger,
	}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var req service.ListProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	resp, err := h.productService.ListProducts(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch products", err)
		return
	}

	utils.PaginatedResponse(c, resp.Products, resp.Page, resp.Limit, resp.Total)
}

func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Product slug is required", nil)
		return
	}

	product, err := h.productService.GetProductBySlug(c.Request.Context(), slug)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Product not found", err)
		return
	}

	utils.SuccessResponse(c, product)
}

func (h *ProductHandler) GetTrendingProducts(c *gin.Context) {
	limit := 8
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	products, err := h.productService.GetTrendingProducts(c.Request.Context(), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch trending products", err)
		return
	}

	utils.SuccessResponse(c, products)
}

func (h *ProductHandler) GetBestSellers(c *gin.Context) {
	limit := 8
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	products, err := h.productService.GetBestSellers(c.Request.Context(), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch best sellers", err)
		return
	}

	utils.SuccessResponse(c, products)
}

func (h *ProductHandler) GetFlashSale(c *gin.Context) {
	items, err := h.productService.GetFlashSale(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch flash sale", err)
		return
	}

	utils.SuccessResponse(c, items)
}

func (h *ProductHandler) GetBrands(c *gin.Context) {
	// This would typically come from a brand repository
	// For now, return empty list
	utils.SuccessResponse(c, []interface{}{})
}

func (h *ProductHandler) GetWishlist(c *gin.Context) {
	// TODO: Implement wishlist functionality
	utils.SuccessResponse(c, []interface{}{})
}

func (h *ProductHandler) AddToWishlist(c *gin.Context) {
	// TODO: Implement add to wishlist
	utils.SuccessWithMessage(c, "Added to wishlist", nil)
}

func (h *ProductHandler) RemoveFromWishlist(c *gin.Context) {
	// TODO: Implement remove from wishlist
	utils.SuccessWithMessage(c, "Removed from wishlist", nil)
}

// Admin methods
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create product", err)
		return
	}

	utils.CreatedResponse(c, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update product", err)
		return
	}

	utils.SuccessResponse(c, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete product", err)
		return
	}

	utils.SuccessWithMessage(c, "Product deleted successfully", nil)
}

func (h *ProductHandler) CreateFlashSale(c *gin.Context) {
	// TODO: Implement create flash sale
	utils.SuccessWithMessage(c, "Flash sale created", nil)
}

func (h *ProductHandler) UpdateFlashSale(c *gin.Context) {
	// TODO: Implement update flash sale
	utils.SuccessWithMessage(c, "Flash sale updated", nil)
}

func (h *ProductHandler) DeleteFlashSale(c *gin.Context) {
	// TODO: Implement delete flash sale
	utils.SuccessWithMessage(c, "Flash sale deleted", nil)
}
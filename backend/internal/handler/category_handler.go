package handler

import (
	"net/http"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	categoryService service.CategoryService
	logger          *zap.Logger
}

func NewCategoryHandler(categoryService service.CategoryService, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		logger:          logger,
	}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch categories", err)
		return
	}

	utils.SuccessResponse(c, categories)
}

func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Category slug is required", nil)
		return
	}

	category, err := h.categoryService.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Category not found", err)
		return
	}

	utils.SuccessResponse(c, category)
}
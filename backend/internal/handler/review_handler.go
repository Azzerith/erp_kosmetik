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

type ReviewHandler struct {
	reviewService service.ReviewService
	logger        *zap.Logger
}

func NewReviewHandler(reviewService service.ReviewService, logger *zap.Logger) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
		logger:        logger,
	}
}

type CreateReviewRequest struct {
	ProductID uint64  `json:"product_id" binding:"required"`
	OrderID   *uint64 `json:"order_id"`
	Rating    int     `json:"rating" binding:"required,min=1,max=5"`
	Title     string  `json:"title"`
	Comment   string  `json:"comment"`
}

type UpdateReviewRequest struct {
	Rating  int     `json:"rating" binding:"min=1,max=5"`
	Title   *string `json:"title"`
	Comment *string `json:"comment"`
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := c.GetUint64("user_id")
	
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	
	review := &models.Review{
		ProductID: req.ProductID,
		UserID:    userID,
		OrderID:   req.OrderID,
		Rating:    req.Rating,
		Title:     &req.Title,
		Comment:   &req.Comment,
		IsApproved:    false, // Admin approval required
		IsVerifiedPurchase: req.OrderID != nil,
	}
	
	if err := h.reviewService.CreateReview(c.Request.Context(), review); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create review", err)
		return
	}
	
	utils.CreatedResponse(c, review)
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID", err)
		return
	}
	
	var req UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	
	review := &models.Review{
		ID:     id,
		UserID: userID,
		Rating: req.Rating,
		Title:  req.Title,
		Comment: req.Comment,
	}
	
	if err := h.reviewService.UpdateReview(c.Request.Context(), review); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update review", err)
		return
	}
	
	utils.SuccessWithMessage(c, "Review updated", nil)
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID", err)
		return
	}
	
	// Verify review belongs to user
	// This would need a GetReviewByID method
	
	if err := h.reviewService.DeleteReview(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete review", err)
		return
	}
	
	utils.SuccessWithMessage(c, "Review deleted", nil)
}

func (h *ReviewHandler) MarkHelpful(c *gin.Context) {
	userID := c.GetUint64("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID", err)
		return
	}
	
	if err := h.reviewService.MarkHelpful(c.Request.Context(), id, userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark helpful", err)
		return
	}
	
	utils.SuccessWithMessage(c, "Marked as helpful", nil)
}
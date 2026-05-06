package handler

import (
	"net/http"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
	logger      *zap.Logger
}

func NewAuthHandler(authService service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int         `json:"expires_in"`
	User         interface{} `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UpdateProfileRequest struct {
	Name  string  `json:"name" binding:"required"`
	Phone string  `json:"phone"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.Name, req.Phone)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Registration failed", err)
		return
	}

	utils.SuccessWithMessage(c, "Registration successful", gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	accessToken, refreshToken, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err)
		return
	}

	utils.SuccessResponse(c, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Refresh token failed", err)
		return
	}

	utils.SuccessResponse(c, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetUint64("user_id")
	refreshToken := c.GetHeader("X-Refresh-Token")

	if err := h.authService.Logout(c.Request.Context(), userID, refreshToken); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Logout failed", err)
		return
	}

	utils.SuccessWithMessage(c, "Logout successful", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	resetToken, err := h.authService.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process request", err)
		return
	}

	// In production, don't return the token in response
	utils.SuccessResponse(c, gin.H{
		"message": "If your email is registered, you will receive a password reset link",
		"token":   resetToken, // Remove in production
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Reset password failed", err)
		return
	}

	utils.SuccessWithMessage(c, "Password reset successful", nil)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.GetUint64("user_id")

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"id":             user.ID,
		"email":          user.Email,
		"name":           user.Name,
		"phone":          user.Phone,
		"avatar_url":     user.AvatarURL,
		"role":           user.Role,
		"loyalty_points": user.LoyaltyPoints,
		"created_at":     user.CreatedAt,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID := c.GetUint64("user_id")

	if err := h.authService.UpdateProfile(c.Request.Context(), userID, req.Name, req.Phone, nil); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Update failed", err)
		return
	}

	utils.SuccessWithMessage(c, "Profile updated successfully", nil)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID := c.GetUint64("user_id")

	if err := h.authService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Change password failed", err)
		return
	}

	utils.SuccessWithMessage(c, "Password changed successfully", nil)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	// Redirect to Google OAuth
	// This will be implemented with actual OAuth flow
	utils.SuccessWithMessage(c, "Google login initiated", gin.H{
		"redirect_url": "/api/v1/auth/google/callback",
	})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Missing authorization code", nil)
		return
	}

	accessToken, refreshToken, user, err := h.authService.LoginWithGoogle(c.Request.Context(), code)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Google login failed", err)
		return
	}

	utils.SuccessResponse(c, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

func (h *AuthHandler) FacebookLogin(c *gin.Context) {
	utils.SuccessWithMessage(c, "Facebook login initiated", gin.H{
		"redirect_url": "/api/v1/auth/facebook/callback",
	})
}

func (h *AuthHandler) FacebookCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Missing authorization code", nil)
		return
	}

	accessToken, refreshToken, user, err := h.authService.LoginWithFacebook(c.Request.Context(), code)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Facebook login failed", err)
		return
	}

	utils.SuccessResponse(c, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}
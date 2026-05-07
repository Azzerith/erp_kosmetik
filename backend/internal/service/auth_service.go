package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"erp-cosmetics-backend/internal/config"
	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"
	"erp-cosmetics-backend/internal/utils"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, email, password, name, phone string) (*models.User, error)
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, user *models.User, err error)
	LoginWithGoogle(ctx context.Context, googleToken string) (accessToken, refreshToken string, user *models.User, err error)
	LoginWithFacebook(ctx context.Context, facebookToken string) (accessToken, refreshToken string, user *models.User, err error)
	RefreshToken(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	Logout(ctx context.Context, userID uint64, refreshToken string) error
	ForgotPassword(ctx context.Context, email string) (resetToken string, err error)
	ResetPassword(ctx context.Context, resetToken, newPassword string) error
	ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error
	GetUserByID(ctx context.Context, userID uint64) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uint64, name, phone string, avatarURL *string) error
}

type authService struct {
	cfg         *config.Config
	db          *gorm.DB
	userRepo    repository.UserRepository
	jwtManager  *utils.JWTManager
	logger      *zap.Logger
}

func NewAuthService(cfg *config.Config, db *gorm.DB, userRepo repository.UserRepository, jwtManager *utils.JWTManager, logger *zap.Logger) AuthService {
	return &authService{
		cfg:        cfg,
		db:         db,
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

func (s *authService) Register(ctx context.Context, email, password, name, phone string) (*models.User, error) {
	// Check if user already exists
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	passwordHash := string(hashedPassword)

	// Create user
	user := &models.User{
		Email:        email,
		PasswordHash: &passwordHash,
		Name:         name,
		Phone:        &phone,
		Provider:     "local",
		Role:         "customer",
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, user *models.User, err error) {
	// Find user by email
	user, err = s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil, errors.New("invalid email or password")
		}
		return "", "", nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return "", "", nil, errors.New("account is deactivated")
	}

	// Check password
	if user.PasswordHash == nil {
		return "", "", nil, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return "", "", nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, refreshToken, err = s.generateTokens(user)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *authService) LoginWithGoogle(ctx context.Context, googleToken string) (accessToken, refreshToken string, user *models.User, err error) {
	// TODO: Verify Google token and get user info
	// For now, implement basic flow
	// This will be implemented with actual Google OAuth verification

	return "", "", nil, errors.New("google login not fully implemented yet")
}

func (s *authService) LoginWithFacebook(ctx context.Context, facebookToken string) (accessToken, refreshToken string, user *models.User, err error) {
	// TODO: Verify Facebook token and get user info
	// This will be implemented with actual Facebook OAuth verification

	return "", "", nil, errors.New("facebook login not fully implemented yet")
}

func (s *authService) RefreshToken(ctx context.Context, refreshTokenStr string) (newAccessToken, newRefreshToken string, err error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// Check if refresh token exists in DB and not revoked
	var refreshTokenModel models.RefreshToken
	if err := s.db.WithContext(ctx).Where("token = ? AND is_revoked = ?", refreshTokenStr, false).First(&refreshTokenModel).Error; err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return "", "", err
	}

	// Revoke old refresh token
	s.db.WithContext(ctx).Model(&refreshTokenModel).Update("is_revoked", true)

	// Generate new tokens
	newAccessToken, newRefreshToken, err = s.generateTokens(user)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *authService) Logout(ctx context.Context, userID uint64, refreshTokenStr string) error {
	// Revoke refresh token
	return s.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("token = ? AND user_id = ?", refreshTokenStr, userID).
		Update("is_revoked", true).Error
}

func (s *authService) ForgotPassword(ctx context.Context, email string) (resetToken string, err error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if email exists or not for security
			return "", nil
		}
		return "", err
	}

	// Generate reset token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	resetToken = hex.EncodeToString(tokenBytes)

	// Save reset token to database
	reset := &models.PasswordReset{
		Email:     user.Email,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		IsUsed:    false,
	}

	if err := s.db.WithContext(ctx).Create(reset).Error; err != nil {
		return "", err
	}

	// TODO: Send email with reset link
	// emailService.SendResetPasswordEmail(user.Email, resetToken)

	return resetToken, nil
}

func (s *authService) ResetPassword(ctx context.Context, resetToken, newPassword string) error {
	// Find reset token
	var reset models.PasswordReset
	if err := s.db.WithContext(ctx).Where("token = ? AND is_used = ? AND expires_at > ?", resetToken, false, time.Now()).First(&reset).Error; err != nil {
		return errors.New("invalid or expired reset token")
	}

	// Mark token as used
	if err := s.db.WithContext(ctx).Model(&reset).Update("is_used", true).Error; err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	passwordHash := string(hashedPassword)

	// Update user password
	return s.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", reset.Email).
		Update("password_hash", passwordHash).Error
}

func (s *authService) ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if user.PasswordHash == nil {
		return errors.New("no password set for this account")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("incorrect old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	passwordHash := string(hashedPassword)

	// Update password
	return s.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).
		Update("password_hash", passwordHash).Error
}

func (s *authService) GetUserByID(ctx context.Context, userID uint64) (*models.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *authService) UpdateProfile(ctx context.Context, userID uint64, name, phone string, avatarURL *string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Name = name
	if phone != "" {
		user.Phone = &phone
	}
	if avatarURL != nil {
		user.AvatarURL = avatarURL
	}

	return s.userRepo.Update(ctx, user)
}

func (s *authService) generateTokens(user *models.User) (accessToken, refreshToken string, err error) {
	// Generate access token
	accessToken, err = s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err = s.jwtManager.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", "", err
	}

	// Store refresh token in database
	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.cfg.JWTRefreshExpiry),
		IsRevoked: false,
	}

	if err := s.db.Create(refreshTokenModel).Error; err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
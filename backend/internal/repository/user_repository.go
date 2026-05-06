package repository

import (
	"context"
	"erp-cosmetics-backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uint64) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByProvider(ctx context.Context, provider, providerID string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint64) error
	UpdateLoyaltyPoints(ctx context.Context, userID uint64, points int) error
	GetAddresses(ctx context.Context, userID uint64) ([]models.Address, error)
	CreateAddress(ctx context.Context, address *models.Address) error
	UpdateAddress(ctx context.Context, address *models.Address) error
	DeleteAddress(ctx context.Context, id uint64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByProvider(ctx context.Context, provider, providerID string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *userRepository) UpdateLoyaltyPoints(ctx context.Context, userID uint64, points int) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).
		Update("loyalty_points", gorm.Expr("loyalty_points + ?", points)).Error
}

func (r *userRepository) GetAddresses(ctx context.Context, userID uint64) ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC").Find(&addresses).Error
	return addresses, err
}

func (r *userRepository) CreateAddress(ctx context.Context, address *models.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *userRepository) UpdateAddress(ctx context.Context, address *models.Address) error {
	return r.db.WithContext(ctx).Save(address).Error
}

func (r *userRepository) DeleteAddress(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&models.Address{}, id).Error
}
package service

import (
	"context"

	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"

	"go.uber.org/zap"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uint64) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	GetAddresses(ctx context.Context, userID uint64) ([]models.Address, error)
	CreateAddress(ctx context.Context, address *models.Address) error
	UpdateAddress(ctx context.Context, address *models.Address) error
	DeleteAddress(ctx context.Context, id uint64) error
	SetDefaultAddress(ctx context.Context, userID, addressID uint64) error
}

type userService struct {
	userRepo repository.UserRepository
	logger   *zap.Logger
}

func NewUserService(userRepo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *userService) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

func (s *userService) GetAddresses(ctx context.Context, userID uint64) ([]models.Address, error) {
	return s.userRepo.GetAddresses(ctx, userID)
}

func (s *userService) CreateAddress(ctx context.Context, address *models.Address) error {
	if address.IsDefault {
		// Remove default from other addresses
		addresses, _ := s.userRepo.GetAddresses(ctx, address.UserID)
		for _, addr := range addresses {
			if addr.IsDefault {
				addr.IsDefault = false
				s.userRepo.UpdateAddress(ctx, &addr)
			}
		}
	}
	return s.userRepo.CreateAddress(ctx, address)
}

func (s *userService) UpdateAddress(ctx context.Context, address *models.Address) error {
	if address.IsDefault {
		// Remove default from other addresses
		addresses, _ := s.userRepo.GetAddresses(ctx, address.UserID)
		for _, addr := range addresses {
			if addr.ID != address.ID && addr.IsDefault {
				addr.IsDefault = false
				s.userRepo.UpdateAddress(ctx, &addr)
			}
		}
	}
	return s.userRepo.UpdateAddress(ctx, address)
}

func (s *userService) DeleteAddress(ctx context.Context, id uint64) error {
	return s.userRepo.DeleteAddress(ctx, id)
}

func (s *userService) SetDefaultAddress(ctx context.Context, userID, addressID uint64) error {
	// Get all user addresses
	addresses, err := s.userRepo.GetAddresses(ctx, userID)
	if err != nil {
		return err
	}

	for _, addr := range addresses {
		addr.IsDefault = (addr.ID == addressID)
		if err := s.userRepo.UpdateAddress(ctx, &addr); err != nil {
			return err
		}
	}

	return nil
}
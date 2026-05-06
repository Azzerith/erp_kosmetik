package service

import (
	"context"
	"errors"
	"time"

	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InventoryService interface {
	AdjustStock(ctx context.Context, productID uint64, variantID *uint64, quantity int, stockType, note string, userID uint64) error
	GetInventoryLogs(ctx context.Context, productID uint64, page, limit int) ([]models.InventoryLog, int64, error)
	GetLowStockProducts(ctx context.Context) ([]models.Product, error)
	CheckAndReleaseExpiredReservations(ctx context.Context) error
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
	productRepo   repository.ProductRepository
	logger        *zap.Logger
	db            *gorm.DB
}

func NewInventoryService(
	inventoryRepo repository.InventoryRepository,
	productRepo repository.ProductRepository,
	logger *zap.Logger,
	db *gorm.DB,
) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
		logger:        logger,
		db:            db,
	}
}

func (s *inventoryService) AdjustStock(ctx context.Context, productID uint64, variantID *uint64, quantity int, stockType, note string, userID uint64) error {
	// Get current stock
	var currentStock int
	var err error
	
	if variantID == nil {
		product, err := s.productRepo.FindByID(ctx, productID)
		if err != nil {
			return err
		}
		currentStock = product.Stock
	} else {
		// Get variant stock
		product, err := s.productRepo.FindByID(ctx, productID)
		if err != nil {
			return err
		}
		for _, variant := range product.Variants {
			if variant.ID == *variantID {
				currentStock = variant.Stock
				break
			}
		}
	}
	
	// Calculate new stock
	newStock := currentStock
	switch stockType {
	case "in":
		newStock += quantity
	case "out":
		if quantity > currentStock {
			return errors.New("insufficient stock")
		}
		newStock -= quantity
	case "adjustment":
		newStock = quantity
	default:
		return errors.New("invalid stock type")
	}
	
	// Update stock in transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update product stock
		if variantID == nil {
			if err := tx.Model(&models.Product{}).Where("id = ?", productID).Update("stock", newStock).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&models.ProductVariant{}).Where("id = ?", variantID).Update("stock", newStock).Error; err != nil {
				return err
			}
		}
		
		// Create inventory log
		log := &models.InventoryLog{
			ProductID:   productID,
			VariantID:   variantID,
			Type:        stockType,
			Quantity:    quantity,
			StockBefore: currentStock,
			StockAfter:  newStock,
			Note:        &note,
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
		}
		
		return s.inventoryRepo.CreateLog(ctx, log)
	})
	
	return err
}

func (s *inventoryService) GetInventoryLogs(ctx context.Context, productID uint64, page, limit int) ([]models.InventoryLog, int64, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 20
	}
	
	offset := (page - 1) * limit
	return s.inventoryRepo.GetLogs(ctx, productID, offset, limit)
}

func (s *inventoryService) GetLowStockProducts(ctx context.Context) ([]models.Product, error) {
	return s.inventoryRepo.GetLowStockProducts(ctx, 5)
}

func (s *inventoryService) CheckAndReleaseExpiredReservations(ctx context.Context) error {
	reservations, err := s.inventoryRepo.GetExpiredReservations(ctx)
	if err != nil {
		return err
	}
	
	for _, reservation := range reservations {
		// Release stock back to inventory
		if err := s.productRepo.UpdateStock(ctx, reservation.ProductID, reservation.VariantID, reservation.Quantity); err != nil {
			s.logger.Error("Failed to release expired reservation",
				zap.Uint64("reservation_id", reservation.ID),
				zap.Error(err))
			continue
		}
		
		// Mark as released
		reservation.IsReleased = true
		now := time.Now()
		reservation.ReleasedAt = &now
		s.db.Save(&reservation)
	}
	
	return nil
}
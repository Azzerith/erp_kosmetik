package cron

import (
	"context"
	"log"
	"time"

	"erp-cosmetics-backend/internal/repository"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderCron struct {
	db          *gorm.DB
	inventoryRepo repository.InventoryRepository
	logger      *zap.Logger
	cron        *cron.Cron
}

func NewOrderCron(db *gorm.DB, inventoryRepo repository.InventoryRepository, logger *zap.Logger) *OrderCron {
	return &OrderCron{
		db:          db,
		inventoryRepo: inventoryRepo,
		logger:      logger,
		cron:        cron.New(cron.WithLocation(time.Local)),
	}
}

func (oc *OrderCron) Start() {
	// Run every hour: 0 * * * *
	_, err := oc.cron.AddFunc("0 * * * *", func() {
		oc.cancelExpiredOrders()
		oc.releaseExpiredReservations()
	})
	if err != nil {
		log.Printf("Failed to add order cron job: %v", err)
	}

	oc.cron.Start()
	oc.logger.Info("Order cron job started")
}

func (oc *OrderCron) Stop() {
	oc.cron.Stop()
	oc.logger.Info("Order cron job stopped")
}

func (oc *OrderCron) cancelExpiredOrders() {
	ctx := context.Background()
	oc.logger.Info("Checking for expired pending orders")

	// Cancel orders pending payment for more than 24 hours
	result := oc.db.WithContext(ctx).Model(&struct{}{}).
		Exec(`UPDATE orders 
		      SET status = 'cancelled', cancelled_at = NOW() 
		      WHERE status = 'pending_payment' 
		      AND created_at < DATE_SUB(NOW(), INTERVAL 24 HOUR)`)

	if result.Error != nil {
		oc.logger.Error("Failed to cancel expired orders", zap.Error(result.Error))
	} else {
		oc.logger.Info("Expired orders cancelled", zap.Int64("rows", result.RowsAffected))
	}
}

func (oc *OrderCron) releaseExpiredReservations() {
	ctx := context.Background()
	oc.logger.Info("Releasing expired stock reservations")

	reservations, err := oc.inventoryRepo.GetExpiredReservations(ctx)
	if err != nil {
		oc.logger.Error("Failed to get expired reservations", zap.Error(err))
		return
	}

	for _, reservation := range reservations {
		// Update stock
		if reservation.VariantID == nil {
			err = oc.db.WithContext(ctx).Model(&struct{}{}).
				Exec("UPDATE products SET stock = stock + ? WHERE id = ?", reservation.Quantity, reservation.ProductID).Error
		} else {
			err = oc.db.WithContext(ctx).Model(&struct{}{}).
				Exec("UPDATE product_variants SET stock = stock + ? WHERE id = ?", reservation.Quantity, reservation.VariantID).Error
		}

		if err != nil {
			oc.logger.Error("Failed to release reservation stock",
				zap.Uint64("reservation_id", reservation.ID),
				zap.Error(err))
			continue
		}

		// Mark as released
		err = oc.db.WithContext(ctx).Model(&struct{}{}).
			Exec("UPDATE stock_reservations SET is_released = true, released_at = NOW() WHERE id = ?", reservation.ID).Error

		if err != nil {
			oc.logger.Error("Failed to mark reservation as released",
				zap.Uint64("reservation_id", reservation.ID),
				zap.Error(err))
		}
	}

	oc.logger.Info("Released expired reservations", zap.Int("count", len(reservations)))
}
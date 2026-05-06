package database

import (
	"erp-cosmetics-backend/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Auto migrate all models
	return db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.PasswordReset{},
		&models.Category{},
		&models.Brand{},
		&models.Product{},
		&models.ProductVariant{},
		&models.ProductImage{},
		&models.ProductTag{},
		&models.ProductCertification{},
		&models.InventoryLog{},
		&models.StockReservation{},
		&models.Address{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderStatusHistory{},
		&models.Payment{},
		&models.PaymentLog{},
		&models.Voucher{},
		&models.VoucherUsage{},
		&models.FlashSale{},
		&models.FlashSaleItem{},
		&models.Review{},
		&models.ReviewHelpful{},
		&models.TrendData{},
		&models.ProductTrendMapping{},
		&models.TrendScoreHistory{},
		&models.Cart{},
		&models.CartItem{},
		&models.Wishlist{},
		&models.Notification{},
		&models.ActivityLog{},
	)
}
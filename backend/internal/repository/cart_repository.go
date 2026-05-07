package repository

import (
	"context"
	"erp-cosmetics-backend/internal/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetOrCreateCart(ctx context.Context, userID uint64, sessionID *string) (*models.Cart, error)
	GetCartWithItems(ctx context.Context, cartID uint64) (*models.Cart, error)
	AddItem(ctx context.Context, item *models.CartItem) error
	UpdateItemQuantity(ctx context.Context, itemID uint64, quantity int) error
	RemoveItem(ctx context.Context, itemID uint64) error
	ClearCart(ctx context.Context, cartID uint64) error
	GetCartByUserID(ctx context.Context, userID uint64) (*models.Cart, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetOrCreateCart(ctx context.Context, userID uint64, sessionID *string) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&cart).Error
	if err == nil {
		return &cart, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	cart = models.Cart{
		UserID:    userID,
		SessionID: sessionID,
	}
	if err := r.db.Create(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetCartWithItems(ctx context.Context, cartID uint64) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Images").
		Preload("Items.Variant").
		First(&cart, cartID).Error
	return &cart, err
}

func (r *cartRepository) AddItem(ctx context.Context, item *models.CartItem) error {
	// Check if item already exists
	var existing models.CartItem
	err := r.db.WithContext(ctx).Where("cart_id = ? AND product_id = ? AND COALESCE(variant_id, 0) = COALESCE(?, 0)",
		item.CartID, item.ProductID, item.VariantID).First(&existing).Error

	if err == nil {
		// Update existing item
		return r.db.Model(&existing).Update("quantity", gorm.Expr("quantity + ?", item.Quantity)).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.Create(item).Error
}

func (r *cartRepository) UpdateItemQuantity(ctx context.Context, itemID uint64, quantity int) error {
	return r.db.WithContext(ctx).Model(&models.CartItem{}).Where("id = ?", itemID).
		Update("quantity", quantity).Error
}

func (r *cartRepository) RemoveItem(ctx context.Context, itemID uint64) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, itemID).Error
}

func (r *cartRepository) ClearCart(ctx context.Context, cartID uint64) error {
	return r.db.WithContext(ctx).Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

func (r *cartRepository) GetCartByUserID(ctx context.Context, userID uint64) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Images").
		Preload("Items.Variant").
		Where("user_id = ?", userID).
		First(&cart).Error
	return &cart, err
}
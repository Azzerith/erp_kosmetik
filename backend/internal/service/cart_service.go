package service

import (
	"context"
	"errors"

	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"

	"go.uber.org/zap"
)

type CartService interface {
	GetCart(ctx context.Context, userID uint64) (*models.Cart, error)
	AddToCart(ctx context.Context, userID uint64, productID uint64, variantID *uint64, quantity int) error
	UpdateCartItem(ctx context.Context, userID uint64, itemID uint64, quantity int) error
	RemoveCartItem(ctx context.Context, userID uint64, itemID uint64) error
	ClearCart(ctx context.Context, userID uint64) error
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	logger      *zap.Logger
}

func NewCartService(
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	logger *zap.Logger,
) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		logger:      logger,
	}
}

func (s *cartService) GetCart(ctx context.Context, userID uint64) (*models.Cart, error) {
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		// Create new cart if not exists
		return s.cartRepo.GetOrCreateCart(ctx, userID, nil)
	}
	return cart, nil
}

func (s *cartService) AddToCart(ctx context.Context, userID uint64, productID uint64, variantID *uint64, quantity int) error {
	// Get or create cart
	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID, nil)
	if err != nil {
		return err
	}

	// Get product to get price
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Check stock
	stock := product.Stock
	if variantID != nil {
		for _, variant := range product.Variants {
			if variant.ID == *variantID {
				stock = variant.Stock
				break
			}
		}
	}
	if quantity > stock {
		return errors.New("insufficient stock")
	}

	// Add item to cart
	price := product.BasePrice
	if product.SalePrice != nil && *product.SalePrice > 0 {
		price = *product.SalePrice
	}

	item := &models.CartItem{
		CartID:        cart.ID,
		ProductID:     productID,
		VariantID:     variantID,
		Quantity:      quantity,
		PriceSnapshot: price,
	}

	return s.cartRepo.AddItem(ctx, item)
}

func (s *cartService) UpdateCartItem(ctx context.Context, userID uint64, itemID uint64, quantity int) error {
	// Verify cart belongs to user
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return errors.New("cart not found")
	}

	// Check if item exists in cart
	var itemExists bool
	for _, item := range cart.Items {
		if item.ID == itemID {
			itemExists = true
			break
		}
	}
	if !itemExists {
		return errors.New("item not found in cart")
	}

	if quantity <= 0 {
		return s.cartRepo.RemoveItem(ctx, itemID)
	}

	return s.cartRepo.UpdateItemQuantity(ctx, itemID, quantity)
}

func (s *cartService) RemoveCartItem(ctx context.Context, userID uint64, itemID uint64) error {
	// Verify cart belongs to user
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return errors.New("cart not found")
	}

	// Check if item exists in cart
	var itemExists bool
	for _, item := range cart.Items {
		if item.ID == itemID {
			itemExists = true
			break
		}
	}
	if !itemExists {
		return errors.New("item not found in cart")
	}

	return s.cartRepo.RemoveItem(ctx, itemID)
}

func (s *cartService) ClearCart(ctx context.Context, userID uint64) error {
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return errors.New("cart not found")
	}
	return s.cartRepo.ClearCart(ctx, cart.ID)
}
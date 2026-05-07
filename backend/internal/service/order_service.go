package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID uint64, req *CreateOrderRequest) (*models.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID uint64, page, limit int) (*ListOrdersResponse, error)
	CancelOrder(ctx context.Context, orderID uint64, userID uint64) error
	UpdateOrderStatus(ctx context.Context, orderID uint64, status string) error
	UpdateTracking(ctx context.Context, orderID uint64, trackingNumber, courier string) error
	GetAllOrders(ctx context.Context, req *AdminListOrdersRequest) (*ListOrdersResponse, error)
}

type CreateOrderRequest struct {
	AddressID          uint64  `json:"address_id" binding:"required"`
	Courier            string  `json:"courier" binding:"required"`
	CourierService     string  `json:"courier_service" binding:"required"`
	ShippingCost       float64 `json:"shipping_cost" binding:"required,min=0"`
	VoucherCode        *string `json:"voucher_code"`
	Notes              *string `json:"notes"`
}

type ListOrdersResponse struct {
	Orders     []models.Order `json:"orders"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type AdminListOrdersRequest struct {
	Page       int     `form:"page" binding:"min=1"`
	Limit      int     `form:"limit" binding:"min=1,max=100"`
	Status     *string `form:"status"`
	UserID     *uint64 `form:"user_id"`
	DateFrom   *string `form:"date_from"`
	DateTo     *string `form:"date_to"`
}

type orderService struct {
	orderRepo     repository.OrderRepository
	cartRepo      repository.CartRepository
	paymentRepo   repository.PaymentRepository
	inventoryRepo repository.InventoryRepository
	logger        *zap.Logger
	db            *gorm.DB
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	paymentRepo repository.PaymentRepository,
	inventoryRepo repository.InventoryRepository,
	logger *zap.Logger,
	db *gorm.DB,
) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		cartRepo:      cartRepo,
		paymentRepo:   paymentRepo,
		inventoryRepo: inventoryRepo,
		logger:        logger,
		db:            db,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, userID uint64, req *CreateOrderRequest) (*models.Order, error) {
	// Get cart items
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Get shipping address
	var address models.Address
	if err := s.db.WithContext(ctx).First(&address, req.AddressID).Error; err != nil {
		return nil, errors.New("shipping address not found")
	}

	// Calculate order totals
	subtotal := 0.0
	var orderItems []models.OrderItem

	for _, item := range cart.Items {
		price := item.PriceSnapshot
		subtotal += price * float64(item.Quantity)

		orderItems = append(orderItems, models.OrderItem{
			ProductID:    item.ProductID,
			VariantID:    item.VariantID,
			ProductName:  item.Product.Name,
			ProductSKU:   item.Product.SKU,
			Price:        price,
			Quantity:     item.Quantity,
			Subtotal:     price * float64(item.Quantity),
			WeightGram:   item.Product.WeightGram,
			DiscountAmount: 0,
		})
	}

	// Apply voucher if any
	discountAmount := 0.0
	if req.VoucherCode != nil && *req.VoucherCode != "" {
		// TODO: Validate and apply voucher
		discountAmount = 0
	}

	totalAmount := subtotal + req.ShippingCost - discountAmount

	// Generate order number
	orderNumber := s.generateOrderNumber()

	// Create order
	order := &models.Order{
		OrderNumber:        orderNumber,
		UserID:             userID,
		Status:             "pending_payment",
		PaymentStatus:      "unpaid",
		FulfillmentStatus:  "unfulfilled",
		Subtotal:           subtotal,
		ShippingCost:       req.ShippingCost,
		DiscountAmount:     discountAmount,
		TaxAmount:          0,
		TotalAmount:        totalAmount,
		ShippingAddressID:  req.AddressID,
		Courier:            &req.Courier,
		CourierService:     &req.CourierService,
		Notes:              req.Notes,
	}

	// Create order with items in transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Create order items
		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		// Clear cart
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		// Create stock reservations
		for _, item := range cart.Items {
			reservation := &models.StockReservation{
				ProductID:  item.ProductID,
				VariantID:  item.VariantID,
				OrderID:    order.ID,
				Quantity:   item.Quantity,
				ExpiresAt:  time.Now().Add(1 * time.Hour), // 1 hour expiry
			}
			if err := tx.Create(reservation).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to create order", zap.Error(err))
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	return s.orderRepo.FindByOrderNumber(ctx, orderNumber)
}

func (s *orderService) GetUserOrders(ctx context.Context, userID uint64, page, limit int) (*ListOrdersResponse, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	orders, total, err := s.orderRepo.FindByUserID(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &ListOrdersResponse{
		Orders:     orders,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *orderService) CancelOrder(ctx context.Context, orderID uint64, userID uint64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errors.New("unauthorized to cancel this order")
	}

	if order.Status != "pending_payment" && order.Status != "paid" {
		return errors.New("order cannot be cancelled at this stage")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update order status
		if err := tx.Model(&models.Order{}).Where("id = ?", orderID).
			Updates(map[string]interface{}{
				"status":       "cancelled",
				"cancelled_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		// Add status history
		history := &models.OrderStatusHistory{
			OrderID:  orderID,
			StatusFrom: &order.Status,
			StatusTo: "cancelled",
			Notes:    stringPtr("Order cancelled by user"),
		}
		if err := tx.Create(history).Error; err != nil {
			return err
		}

		// Release stock reservations
		if err := tx.Model(&models.StockReservation{}).
			Where("order_id = ? AND is_released = ?", orderID, false).
			Updates(map[string]interface{}{
				"is_released": true,
				"released_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID uint64, status string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Define timestamp based on status
	var timestampField map[string]interface{}
	switch status {
	case "paid":
		timestampField = map[string]interface{}{"paid_at": time.Now()}
	case "processing":
		timestampField = map[string]interface{}{"processed_at": time.Now()}
	case "shipped":
		timestampField = map[string]interface{}{"shipped_at": time.Now()}
	case "delivered":
		timestampField = map[string]interface{}{"delivered_at": time.Now()}
	case "completed":
		timestampField = map[string]interface{}{"completed_at": time.Now()}
	default:
		timestampField = map[string]interface{}{}
	}

	// Update order status
	updates := map[string]interface{}{
		"status": status,
	}
	for k, v := range timestampField {
		updates[k] = v
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, status, "", ""); err != nil {
		return err
	}

	// Add status history
	history := &models.OrderStatusHistory{
		OrderID:    orderID,
		StatusFrom: &order.Status,
		StatusTo:   status,
	}
	return s.orderRepo.AddStatusHistory(ctx, history)
}

func (s *orderService) UpdateTracking(ctx context.Context, orderID uint64, trackingNumber, courier string) error {
	return s.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"tracking_number": trackingNumber,
			"courier":         courier,
			"status":          "shipped",
			"shipped_at":      time.Now(),
		}).Error
}

func (s *orderService) GetAllOrders(ctx context.Context, req *AdminListOrdersRequest) (*ListOrdersResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	filters := make(map[string]interface{})
	if req.Status != nil {
		filters["status"] = *req.Status
	}
	if req.UserID != nil {
		filters["user_id"] = *req.UserID
	}
	if req.DateFrom != nil {
		filters["date_from"] = *req.DateFrom
	}
	if req.DateTo != nil {
		filters["date_to"] = *req.DateTo
	}

	orders, total, err := s.orderRepo.GetAll(ctx, offset, req.Limit, filters)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &ListOrdersResponse{
		Orders:     orders,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *orderService) generateOrderNumber() string {
	return fmt.Sprintf("ORD-%s-%d", uuid.New().String()[:8], time.Now().Unix())
}

func stringPtr(s string) *string {
	return &s
}
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"erp-cosmetics-backend/internal/config"
	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PaymentService interface {
	InitiatePayment(ctx context.Context, orderID uint64) (snapToken string, err error)
	ProcessWebhook(ctx context.Context, notification map[string]interface{}) error
	GetPaymentStatus(ctx context.Context, orderID string) (*PaymentStatusResponse, error)
	RefundPayment(ctx context.Context, paymentID uint64, amount float64, reason string) error
}

type PaymentStatusResponse struct {
	OrderID     string  `json:"order_id"`
	Status      string  `json:"status"`
	Amount      float64 `json:"amount"`
	PaymentType string  `json:"payment_type"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
}

type paymentService struct {
	cfg          *config.Config
	paymentRepo  repository.PaymentRepository
	orderRepo    repository.OrderRepository
	inventoryRepo repository.InventoryRepository
	logger       *zap.Logger
	db           *gorm.DB
	snapClient   snap.Client
	coreClient   coreapi.Client
}

func NewPaymentService(
	cfg *config.Config,
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
	inventoryRepo repository.InventoryRepository,
	logger *zap.Logger,
) PaymentService {
	// Initialize Midtrans client
	var snapClient snap.Client
	var coreClient coreapi.Client
	
	if cfg.MidtransServerKey != "" {
		snapClient.New(cfg.MidtransServerKey, midtrans.Sandbox)
		if cfg.MidtransEnvironment == "production" {
			snapClient.New(cfg.MidtransServerKey, midtrans.Production)
		}
		
		coreClient.New(cfg.MidtransServerKey, midtrans.Sandbox)
		if cfg.MidtransEnvironment == "production" {
			coreClient.New(cfg.MidtransServerKey, midtrans.Production)
		}
	}

	return &paymentService{
		cfg:          cfg,
		paymentRepo:  paymentRepo,
		orderRepo:    orderRepo,
		inventoryRepo: inventoryRepo,
		logger:       logger,
		snapClient:   snapClient,
		coreClient:   coreClient,
	}
}

func (s *paymentService) InitiatePayment(ctx context.Context, orderID uint64) (string, error) {
	// Get order details
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return "", errors.New("order not found")
	}

	if order.Status != "pending_payment" {
		return "", errors.New("order cannot be paid")
	}

	// Prepare order items for Midtrans
	var items []midtrans.ItemDetails
	for _, item := range order.Items {
		items = append(items, midtrans.ItemDetails{
			ID:    fmt.Sprintf("%d", item.ProductID),
			Name:  item.ProductName,
			Price: int64(item.Price),
			Qty:   int32(item.Quantity),
		})
	}

	// Add shipping as item
	items = append(items, midtrans.ItemDetails{
		ID:    "shipping",
		Name:  "Biaya Pengiriman",
		Price: int64(order.ShippingCost),
		Qty:   1,
	})

	// Create Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.OrderNumber,
			GrossAmt: int64(order.TotalAmount),
		},
		Items: &items,
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.User.Name,
			Email: order.User.Email,
			Phone: s.getUserPhone(order.User),
		},
		EnabledPayments: snap.AllPaymentType,
	}

	// Create Snap transaction
	snapResp, err := s.snapClient.CreateTransaction(req)
	if err != nil {
		s.logger.Error("Failed to create Snap transaction", zap.Error(err))
		return "", errors.New("failed to initiate payment")
	}

	// Save payment record
	payment := &models.Payment{
		OrderID:       order.ID,
		TransactionID: order.OrderNumber,
		SnapToken:     &snapResp.Token,
		PaymentMethod: "midtrans_snap",
		Amount:        order.TotalAmount,
		Status:        "pending",
		ExpiredAt:     s.getExpiredTime(),
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		s.logger.Error("Failed to save payment record", zap.Error(err))
		return "", err
	}

	// Create payment log
	s.createPaymentLog(ctx, payment.ID, "initiate", map[string]interface{}{
		"snap_response": snapResp,
	})

	return snapResp.Token, nil
}

func (s *paymentService) ProcessWebhook(ctx context.Context, notification map[string]interface{}) error {
	orderID, ok := notification["order_id"].(string)
	if !ok {
		return errors.New("invalid order_id in notification")
	}

	transactionStatus, ok := notification["transaction_status"].(string)
	if !ok {
		return errors.New("invalid transaction_status")
	}

	paymentType, _ := notification["payment_type"].(string)

	// Get payment record
	payment, err := s.paymentRepo.FindByTransactionID(ctx, orderID)
	if err != nil {
		return errors.New("payment not found")
	}

	// Get order
	order, err := s.orderRepo.FindByID(ctx, payment.OrderID)
	if err != nil {
		return err
	}

	// Process based on transaction status
	switch transactionStatus {
	case "capture", "settlement":
		// Payment successful
		if err := s.handleSuccessfulPayment(ctx, order, payment, notification, paymentType); err != nil {
			return err
		}
	case "pending":
		// Payment pending - update status only
		payment.Status = "pending"
		s.paymentRepo.Update(ctx, payment)
	case "deny", "expire", "cancel":
		// Payment failed - release stock reservations
		if err := s.handleFailedPayment(ctx, order, payment, notification); err != nil {
			return err
		}
	case "refund":
		// Refund processed
		payment.Status = "refunded"
		s.paymentRepo.Update(ctx, payment)
	}

	// Create payment log
	s.createPaymentLog(ctx, payment.ID, "webhook", notification)

	return nil
}

func (s *paymentService) GetPaymentStatus(ctx context.Context, orderID string) (*PaymentStatusResponse, error) {
	payment, err := s.paymentRepo.FindByTransactionID(ctx, orderID)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	// Check status from Midtrans
	statusResp, err := s.coreClient.CheckTransaction(orderID)
	if err == nil {
		payment.Status = s.mapMidtransStatus(statusResp.TransactionStatus)
		s.paymentRepo.Update(ctx, payment)
	}

	return &PaymentStatusResponse{
		OrderID:     payment.TransactionID,
		Status:      payment.Status,
		Amount:      payment.Amount,
		PaymentType: payment.PaymentMethod,
		PaidAt:      payment.PaidAt,
	}, nil
}

func (s *paymentService) RefundPayment(ctx context.Context, paymentID uint64, amount float64, reason string) error {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return errors.New("payment not found")
	}

	// Process refund via Midtrans
	refundReq := &coreapi.RefundReq{
		RefundKey: fmt.Sprintf("refund-%d-%d", paymentID, time.Now().Unix()),
		Amount:    int64(amount),
		Reason:    reason,
	}

	refundResp, err := s.coreClient.Refund(payment.TransactionID, refundReq)
	if err != nil {
		s.logger.Error("Failed to process refund", zap.Error(err))
		return errors.New("refund failed")
	}

	// Update payment status
	payment.Status = "refunded"
	s.paymentRepo.Update(ctx, payment)

	// Update order status
	s.orderRepo.UpdateStatus(ctx, payment.OrderID, "refunded", "refunded", "")

	// Create refund log
	s.createPaymentLog(ctx, payment.ID, "refund", map[string]interface{}{
		"refund_response": refundResp,
		"amount":          amount,
		"reason":          reason,
	})

	return nil
}

func (s *paymentService) handleSuccessfulPayment(ctx context.Context, order *models.Order, payment *models.Payment, notification map[string]interface{}, paymentType string) error {
	now := time.Now()

	// Update payment
	payment.Status = "success"
	payment.PaymentMethod = paymentType
	payment.PaidAt = &now
	payment.RawResponse = s.toJSON(notification)

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return err
	}

	// Update order
	order.Status = "paid"
	order.PaymentStatus = "paid"
	order.PaidAt = &now

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	// Add status history
	history := &models.OrderStatusHistory{
		OrderID:    order.ID,
		StatusFrom: stringPtr("pending_payment"),
		StatusTo:   "paid",
		Notes:      stringPtr("Payment successful"),
	}
	s.orderRepo.AddStatusHistory(ctx, history)

	// Release stock (reservations are converted to actual stock deduction)
	if err := s.releaseStockReservations(ctx, order.ID); err != nil {
		s.logger.Error("Failed to release stock reservations", zap.Error(err))
	}

	// TODO: Send email notification to customer
	// s.emailService.SendOrderPaidEmail(order.User.Email, order.OrderNumber)

	return nil
}

func (s *paymentService) handleFailedPayment(ctx context.Context, order *models.Order, payment *models.Payment, notification map[string]interface{}) error {
	// Update payment
	payment.Status = "failed"
	payment.RawResponse = s.toJSON(notification)

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return err
	}

	// Update order status
	order.Status = "cancelled"
	order.CancelledAt = s.timePtr(time.Now())

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	// Add status history
	history := &models.OrderStatusHistory{
		OrderID:    order.ID,
		StatusFrom: stringPtr("pending_payment"),
		StatusTo:   "cancelled",
		Notes:      stringPtr("Payment failed: " + notification["status_message"].(string)),
	}
	s.orderRepo.AddStatusHistory(ctx, history)

	// Release stock reservations
	return s.releaseStockReservations(ctx, order.ID)
}

func (s *paymentService) releaseStockReservations(ctx context.Context, orderID uint64) error {
	return s.inventoryRepo.ReleaseReservation(ctx, orderID)
}

func (s *paymentService) createPaymentLog(ctx context.Context, paymentID uint64, eventType string, data interface{}) {
	log := &models.PaymentLog{
		PaymentID: paymentID,
		EventType: eventType,
		EventData: s.toJSON(data),
	}
	s.paymentRepo.CreateLog(ctx, log)
}

func (s *paymentService) getUserPhone(user models.User) string {
	if user.Phone != nil {
		return *user.Phone
	}
	return ""
}

func (s *paymentService) getExpiredTime() *time.Time {
	t := time.Now().Add(24 * time.Hour)
	return &t
}

func (s *paymentService) mapMidtransStatus(status string) string {
	switch status {
	case "capture", "settlement":
		return "success"
	case "pending":
		return "pending"
	case "deny", "expire", "cancel":
		return "failed"
	case "refund":
		return "refunded"
	default:
		return "pending"
	}
}

func (s *paymentService) toJSON(data interface{}) struct {
	Data interface{} `json:"data"`
} {
	return struct{ Data interface{} }{Data: data}
}

func (s *paymentService) timePtr(t time.Time) *time.Time {
	return &t
}
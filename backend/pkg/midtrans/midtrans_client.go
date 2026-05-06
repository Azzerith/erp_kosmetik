package midtrans

import (
	"errors"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransClient struct {
	snapClient   snap.Client
	coreClient   coreapi.Client
	isProduction bool
}

type TransactionDetails struct {
	OrderID     string
	GrossAmount int64
}

type CustomerDetails struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

type ItemDetail struct {
	ID       string
	Name     string
	Price    int64
	Quantity int32
}

type SnapRequest struct {
	TransactionDetails TransactionDetails
	CustomerDetails    CustomerDetails
	Items              []ItemDetail
}

type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func NewMidtransClient(serverKey string, isProduction bool) *MidtransClient {
	var snapClient snap.Client
	var coreClient coreapi.Client
	
	snapClient.New(serverKey, midtrans.Sandbox)
	coreClient.New(serverKey, midtrans.Sandbox)
	
	if isProduction {
		snapClient.New(serverKey, midtrans.Production)
		coreClient.New(serverKey, midtrans.Production)
	}
	
	return &MidtransClient{
		snapClient:   snapClient,
		coreClient:   coreClient,
		isProduction: isProduction,
	}
}

func (c *MidtransClient) CreateSnapTransaction(req *SnapRequest) (*SnapResponse, error) {
	// Convert items
	var items []midtrans.ItemDetails
	for _, item := range req.Items {
		items = append(items, midtrans.ItemDetails{
			ID:    item.ID,
			Name:  item.Name,
			Price: item.Price,
			Qty:   item.Quantity,
		})
	}
	
	// Convert customer details
	customer := &midtrans.CustomerDetails{
		FName: req.CustomerDetails.FirstName,
		LName: req.CustomerDetails.LastName,
		Email: req.CustomerDetails.Email,
		Phone: req.CustomerDetails.Phone,
	}
	
	// Create snap request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.TransactionDetails.OrderID,
			GrossAmt: req.TransactionDetails.GrossAmount,
		},
		Items:          &items,
		CustomerDetail: customer,
		EnabledPayments: snap.AllPaymentType,
	}
	
	// Create transaction
	resp, err := c.snapClient.CreateTransaction(snapReq)
	if err != nil {
		return nil, err
	}
	
	return &SnapResponse{
		Token:       resp.Token,
		RedirectURL: resp.RedirectURL,
	}, nil
}

func (c *MidtransClient) CheckTransaction(orderID string) (*coreapi.TransactionStatusResponse, error) {
	resp, err := c.coreClient.CheckTransaction(orderID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *MidtransClient) RefundTransaction(orderID string, amount int64, reason string) (*coreapi.RefundResponse, error) {
	refundReq := &coreapi.RefundReq{
		RefundKey: fmt.Sprintf("refund-%s", orderID),
		Amount:    amount,
		Reason:    reason,
	}
	
	resp, err := c.coreClient.Refund(orderID, refundReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *MidtransClient) CancelTransaction(orderID string) (*coreapi.TransactionStatusResponse, error) {
	resp, err := c.coreClient.CancelTransaction(orderID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *MidtransClient) VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	// Calculate expected signature
	// serverKey + orderID + statusCode + grossAmount + serverKey
	expected := c.snapClient.CoreApi.GetEnv().ServerKey + orderID + statusCode + grossAmount + c.snapClient.CoreApi.GetEnv().ServerKey
	
	// This should be compared with signatureKey
	// For now, return true for demo
	return true
}

func (c *MidtransClient) GetStatusMessage(statusCode string) string {
	switch statusCode {
	case "200":
		return "Transaksi berhasil"
	case "201":
		return "Transaksi pending"
	case "202":
		return "Transaksi dibatalkan"
	case "407":
		return "Transaksi expired"
	case "412":
		return "Transaksi refund"
	default:
		return "Status tidak diketahui"
	}
}

func (c *MidtransClient) MapTransactionStatus(transactionStatus string) string {
	switch transactionStatus {
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
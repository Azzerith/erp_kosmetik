package midtrans

import (
	"fmt"

	"erp-cosmetics-backend/internal/utils"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransClient struct {
	snapClient   snap.Client
	coreClient   coreapi.Client
	serverKey    string
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
	
	env := midtrans.Sandbox
	if isProduction {
		env = midtrans.Production
	}
	
	snapClient.New(serverKey, env)
	coreClient.New(serverKey, env)
	
	return &MidtransClient{
		snapClient:   snapClient,
		coreClient:   coreClient,
		serverKey:    serverKey,
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
	
	resp, err := c.coreClient.RefundTransaction(orderID, refundReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *MidtransClient) CancelTransaction(orderID string) (*coreapi.ChargeResponse, error) {
	resp, err := c.coreClient.CancelTransaction(orderID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *MidtransClient) VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	// Calculate expected signature
	// serverKey + orderID + statusCode + grossAmount + serverKey
	data := orderID + statusCode + grossAmount + c.serverKey
	expected := utils.GenerateSHA512(data)
	
	return expected == signatureKey
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
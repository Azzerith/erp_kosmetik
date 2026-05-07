package service

import (
	"context"
	"time"

	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReportService interface {
	GetSalesReport(ctx context.Context, req *SalesReportRequest) (*SalesReportResponse, error)
	GetInventoryReport(ctx context.Context) (*InventoryReportResponse, error)
	GetCustomerReport(ctx context.Context) (*CustomerReportResponse, error)
	GetDashboardSummary(ctx context.Context) (*DashboardSummary, error)
}

type SalesReportRequest struct {
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
	GroupBy  string `form:"group_by"` // day, month, year
}

type SalesReportResponse struct {
	TotalRevenue     float64                `json:"total_revenue"`
	TotalOrders      int64                  `json:"total_orders"`
	AverageOrderValue float64               `json:"average_order_value"`
	TopProducts      []TopProduct           `json:"top_products"`
	DailySales       []DailySales           `json:"daily_sales"`
	RevenueByStatus  map[string]float64     `json:"revenue_by_status"`
}

type TopProduct struct {
	ProductID   uint64  `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalSold   int     `json:"total_sold"`
	Revenue     float64 `json:"revenue"`
}

type DailySales struct {
	Date   string  `json:"date"`
	Orders int64   `json:"orders"`
	Revenue float64 `json:"revenue"`
}

type InventoryReportResponse struct {
	TotalProducts     int64                `json:"total_products"`
	TotalStock        int                  `json:"total_stock"`
	TotalValue        float64              `json:"total_value"`
	LowStockProducts  []LowStockProduct    `json:"low_stock_products"`
	OutOfStockProducts []OutOfStockProduct `json:"out_of_stock_products"`
	TopCategories     []CategoryStock      `json:"top_categories"`
}

type LowStockProduct struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	SKU      string `json:"sku"`
	Stock    int    `json:"stock"`
	Threshold int   `json:"threshold"`
}

type OutOfStockProduct struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	SKU  string `json:"sku"`
}

type CategoryStock struct {
	CategoryName string `json:"category_name"`
	TotalStock   int    `json:"total_stock"`
	ProductCount int64  `json:"product_count"`
}

type CustomerReportResponse struct {
	TotalCustomers    int64               `json:"total_customers"`
	NewCustomers      int64               `json:"new_customers"`
	ActiveCustomers   int64               `json:"active_customers"`
	TopCustomers      []TopCustomer       `json:"top_customers"`
	CustomerBySource  map[string]int64    `json:"customer_by_source"`
	CustomerByMonth   []CustomerByMonth   `json:"customer_by_month"`
}

type TopCustomer struct {
	UserID      uint64  `json:"user_id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	TotalSpent  float64 `json:"total_spent"`
	OrderCount  int64   `json:"order_count"`
}

type CustomerByMonth struct {
	Month   string `json:"month"`
	Count   int64  `json:"count"`
	Revenue float64 `json:"revenue"`
}

type DashboardSummary struct {
	TodaySales        float64   `json:"today_sales"`
	TodayOrders       int64     `json:"today_orders"`
	WeekSales         float64   `json:"week_sales"`
	MonthSales        float64   `json:"month_sales"`
	TotalCustomers    int64     `json:"total_customers"`
	TotalProducts     int64     `json:"total_products"`
	LowStockCount     int64     `json:"low_stock_count"`
	PendingOrders     int64     `json:"pending_orders"`
	AverageTrendScore float64   `json:"average_trend_score"`
	LastUpdated       time.Time `json:"last_updated"`
}

type reportService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	logger      *zap.Logger
	db          *gorm.DB
}

func NewReportService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	logger *zap.Logger,
) ReportService {
	return &reportService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

func (s *reportService) GetSalesReport(ctx context.Context, req *SalesReportRequest) (*SalesReportResponse, error) {
	// Parse dates
	dateFrom := time.Now().AddDate(0, 0, -30)
	dateTo := time.Now()
	
	if req.DateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", req.DateFrom); err == nil {
			dateFrom = parsed
		}
	}
	if req.DateTo != "" {
		if parsed, err := time.Parse("2006-01-02", req.DateTo); err == nil {
			dateTo = parsed
		}
	}
	
	// Get completed orders within date range
	var orders []models.Order
	var totalRevenue float64
	var totalOrders int64
	
	s.db.WithContext(ctx).Model(&models.Order{}).
		Where("status IN (?) AND paid_at BETWEEN ? AND ?", []string{"delivered", "completed", "paid"}, dateFrom, dateTo).
		Find(&orders)
	
	for _, order := range orders {
		totalRevenue += order.TotalAmount
	}
	totalOrders = int64(len(orders))
	
	// Calculate average order value
	avgOrderValue := 0.0
	if totalOrders > 0 {
		avgOrderValue = totalRevenue / float64(totalOrders)
	}
	
	// Get top products
	var topProducts []TopProduct
	s.db.WithContext(ctx).Table("order_items").
		Select("order_items.product_id, products.name as product_name, SUM(order_items.quantity) as total_sold, SUM(order_items.subtotal) as revenue").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("orders.status IN (?) AND orders.paid_at BETWEEN ? AND ?", []string{"delivered", "completed", "paid"}, dateFrom, dateTo).
		Group("order_items.product_id, products.name").
		Order("total_sold DESC").
		Limit(10).
		Scan(&topProducts)
	
	// Get daily sales
	var dailySales []DailySales
	for d := dateFrom; d.Before(dateTo) || d.Equal(dateTo); d = d.AddDate(0, 0, 1) {
		var dayRevenue float64
		var dayOrders int64
		
		nextDay := d.AddDate(0, 0, 1)
		s.db.WithContext(ctx).Model(&models.Order{}).
			Where("status IN (?) AND paid_at BETWEEN ? AND ?", []string{"delivered", "completed", "paid"}, d, nextDay).
			Select("COALESCE(SUM(total_amount), 0) as revenue, COUNT(*) as orders").
			Row().Scan(&dayRevenue, &dayOrders)
		
		dailySales = append(dailySales, DailySales{
			Date:    d.Format("2006-01-02"),
			Orders:  dayOrders,
			Revenue: dayRevenue,
		})
	}
	
	// Revenue by status
	var revenueByStatus map[string]float64
	s.db.WithContext(ctx).Model(&models.Order{}).
		Select("status, SUM(total_amount) as revenue").
		Where("paid_at BETWEEN ? AND ?", dateFrom, dateTo).
		Group("status").
		Scan(&revenueByStatus)
	
	return &SalesReportResponse{
		TotalRevenue:      totalRevenue,
		TotalOrders:       totalOrders,
		AverageOrderValue: avgOrderValue,
		TopProducts:       topProducts,
		DailySales:        dailySales,
		RevenueByStatus:   revenueByStatus,
	}, nil
}

func (s *reportService) GetInventoryReport(ctx context.Context) (*InventoryReportResponse, error) {
	var totalProducts int64
	var totalStock int
	var totalValue float64
	
	s.db.WithContext(ctx).Model(&models.Product{}).
		Where("is_active = ?", true).
		Count(&totalProducts)
	
	s.db.WithContext(ctx).Model(&models.Product{}).
		Select("COALESCE(SUM(stock), 0) as total_stock, COALESCE(SUM(CASE WHEN sale_price IS NOT NULL THEN sale_price * stock ELSE base_price * stock END), 0) as total_value").
		Row().Scan(&totalStock, &totalValue)
	
	// Low stock products
	var lowStockProducts []LowStockProduct
	s.db.WithContext(ctx).Model(&models.Product{}).
		Select("id, name, sku, stock, min_stock_threshold as threshold").
		Where("is_active = ? AND stock <= min_stock_threshold AND stock > 0", true).
		Find(&lowStockProducts)
	
	// Out of stock products
	var outOfStockProducts []OutOfStockProduct
	s.db.WithContext(ctx).Model(&models.Product{}).
		Select("id, name, sku").
		Where("is_active = ? AND stock = 0", true).
		Find(&outOfStockProducts)
	
	// Top categories by stock
	var topCategories []CategoryStock
	s.db.WithContext(ctx).Table("products").
		Select("categories.name as category_name, SUM(products.stock) as total_stock, COUNT(products.id) as product_count").
		Joins("JOIN categories ON categories.id = products.category_id").
		Where("products.is_active = ?", true).
		Group("categories.id, categories.name").
		Order("total_stock DESC").
		Limit(10).
		Scan(&topCategories)
	
	return &InventoryReportResponse{
		TotalProducts:      totalProducts,
		TotalStock:         totalStock,
		TotalValue:         totalValue,
		LowStockProducts:   lowStockProducts,
		OutOfStockProducts: outOfStockProducts,
		TopCategories:      topCategories,
	}, nil
}

func (s *reportService) GetCustomerReport(ctx context.Context) (*CustomerReportResponse, error) {
	var totalCustomers int64
	var newCustomers int64
	var activeCustomers int64
	
	s.db.WithContext(ctx).Model(&models.User{}).
		Where("role = ?", "customer").
		Count(&totalCustomers)
	
	// New customers in last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	s.db.WithContext(ctx).Model(&models.User{}).
		Where("role = ? AND created_at >= ?", "customer", thirtyDaysAgo).
		Count(&newCustomers)
	
	// Active customers (made purchase in last 90 days)
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)
	s.db.WithContext(ctx).Model(&models.User{}).
		Joins("JOIN orders ON orders.user_id = users.id").
		Where("users.role = ? AND orders.created_at >= ?", "customer", ninetyDaysAgo).
		Distinct("users.id").
		Count(&activeCustomers)
	
	// Top customers by spending
	var topCustomers []TopCustomer
	s.db.WithContext(ctx).Table("orders").
		Select("users.id as user_id, users.name, users.email, SUM(orders.total_amount) as total_spent, COUNT(orders.id) as order_count").
		Joins("JOIN users ON users.id = orders.user_id").
		Where("orders.status IN (?)", []string{"delivered", "completed", "paid"}).
		Group("users.id, users.name, users.email").
		Order("total_spent DESC").
		Limit(10).
		Scan(&topCustomers)
	
	// Customer by source
	var customerBySource map[string]int64
	s.db.WithContext(ctx).Model(&models.User{}).
		Where("role = ?", "customer").
		Select("provider, COUNT(*) as count").
		Group("provider").
		Scan(&customerBySource)
	
	// Customer by month
	var customerByMonth []CustomerByMonth
	for i := 11; i >= 0; i-- {
		month := time.Now().AddDate(0, -i, 0)
		monthStart := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.Local)
		monthEnd := monthStart.AddDate(0, 1, 0)
		
		var count int64
		var revenue float64
		
		s.db.WithContext(ctx).Model(&models.User{}).
			Where("role = ? AND created_at BETWEEN ? AND ?", "customer", monthStart, monthEnd).
			Count(&count)
		
		s.db.WithContext(ctx).Model(&models.Order{}).
			Where("status IN (?) AND paid_at BETWEEN ? AND ?", []string{"delivered", "completed", "paid"}, monthStart, monthEnd).
			Select("COALESCE(SUM(total_amount), 0)").
			Scan(&revenue)
		
		customerByMonth = append(customerByMonth, CustomerByMonth{
			Month:   monthStart.Format("2006-01"),
			Count:   count,
			Revenue: revenue,
		})
	}
	
	return &CustomerReportResponse{
		TotalCustomers:   totalCustomers,
		NewCustomers:     newCustomers,
		ActiveCustomers:  activeCustomers,
		TopCustomers:     topCustomers,
		CustomerBySource: customerBySource,
		CustomerByMonth:  customerByMonth,
	}, nil
}

func (s *reportService) GetDashboardSummary(ctx context.Context) (*DashboardSummary, error) {
	// Today's sales
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	todayEnd := todayStart.AddDate(0, 0, 1)
	
	var todaySales float64
	var todayOrders int64
	
	s.db.WithContext(ctx).Model(&models.Order{}).
		Where("status IN (?) AND paid_at BETWEEN ? AND ?", []string{"delivered", "completed", "paid"}, todayStart, todayEnd).
		Select("COALESCE(SUM(total_amount), 0) as sales, COUNT(*) as orders").
		Row().Scan(&todaySales, &todayOrders)
	
	// Week sales
	weekStart := todayStart.AddDate(0, 0, -int(today.Weekday()))
	var weekSales float64
	s.db.WithContext(ctx).Model(&models.Order{}).
		Where("status IN (?) AND paid_at >= ?", []string{"delivered", "completed", "paid"}, weekStart).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&weekSales)
	
	// Month sales
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.Local)
	var monthSales float64
	s.db.WithContext(ctx).Model(&models.Order{}).
		Where("status IN (?) AND paid_at >= ?", []string{"delivered", "completed", "paid"}, monthStart).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&monthSales)
	
	// Total customers
	var totalCustomers int64
	s.db.WithContext(ctx).Model(&models.User{}).
		Where("role = ?", "customer").
		Count(&totalCustomers)
	
	// Total products
	var totalProducts int64
	s.db.WithContext(ctx).Model(&models.Product{}).
		Where("is_active = ?", true).
		Count(&totalProducts)
	
	// Low stock count
	var lowStockCount int64
	s.db.WithContext(ctx).Model(&models.Product{}).
		Where("is_active = ? AND stock <= min_stock_threshold", true).
		Count(&lowStockCount)
	
	// Pending orders
	var pendingOrders int64
	s.db.WithContext(ctx).Model(&models.Order{}).
		Where("status IN (?)", []string{"pending_payment", "paid", "processing"}).
		Count(&pendingOrders)
	
	// Average trend score
	var avgTrendScore float64
	s.db.WithContext(ctx).Model(&models.Product{}).
		Where("is_active = ?", true).
		Select("COALESCE(AVG(trend_score), 0)").
		Scan(&avgTrendScore)
	
	return &DashboardSummary{
		TodaySales:        todaySales,
		TodayOrders:       todayOrders,
		WeekSales:         weekSales,
		MonthSales:        monthSales,
		TotalCustomers:    totalCustomers,
		TotalProducts:     totalProducts,
		LowStockCount:     lowStockCount,
		PendingOrders:     pendingOrders,
		AverageTrendScore: avgTrendScore,
		LastUpdated:       time.Now(),
	}, nil
}
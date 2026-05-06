package routes

import (
	"erp-cosmetics-backend/internal/config"
	"erp-cosmetics-backend/internal/handler"
	"erp-cosmetics-backend/internal/middleware"
	"erp-cosmetics-backend/internal/repository"
	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupRoutes(
	router *gin.Engine,
	cfg *config.Config,
	db *gorm.DB,
	redisClient *redis.Client,
	jwtManager *utils.JWTManager,
	logger *zap.Logger,
) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewCartRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	voucherRepo := repository.NewVoucherRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	trendRepo := repository.NewTrendRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)

	// Initialize services
	authService := service.NewAuthService(cfg, db, userRepo, jwtManager, logger)
	userService := service.NewUserService(userRepo, logger)
	productService := service.NewProductService(productRepo, categoryRepo, logger)
	categoryService := service.NewCategoryService(categoryRepo, logger)
	orderService := service.NewOrderService(orderRepo, cartRepo, paymentRepo, inventoryRepo, logger, db)
	cartService := service.NewCartService(cartRepo, productRepo, logger)
	paymentService := service.NewPaymentService(cfg, paymentRepo, orderRepo, inventoryRepo, logger, db)
	voucherService := service.NewVoucherService(voucherRepo, logger)
	reviewService := service.NewReviewService(reviewRepo, logger)
	trendService := service.NewTrendService(cfg, trendRepo, productRepo, redisClient, logger)
	shippingService := service.NewShippingService(cfg, logger)
	reportService := service.NewReportService(orderRepo, productRepo, userRepo, logger)
	inventoryService := service.NewInventoryService(inventoryRepo, productRepo, logger, db)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, logger)
	userHandler := handler.NewUserHandler(userService, logger)
	productHandler := handler.NewProductHandler(productService, logger)
	categoryHandler := handler.NewCategoryHandler(categoryService, logger)
	orderHandler := handler.NewOrderHandler(orderService, logger)
	cartHandler := handler.NewCartHandler(cartService, logger)
	paymentHandler := handler.NewPaymentHandler(paymentService, logger)
	voucherHandler := handler.NewVoucherHandler(voucherService, logger)
	reviewHandler := handler.NewReviewHandler(reviewService, logger)
	trendHandler := handler.NewTrendHandler(trendService, logger)
	shippingHandler := handler.NewShippingHandler(shippingService, logger)
	reportHandler := handler.NewReportHandler(reportService, logger)
	inventoryHandler := handler.NewInventoryHandler(inventoryService, logger)
	webhookHandler := handler.NewWebhookHandler(paymentService, logger)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("/")
		{
			// Auth
			public.POST("/auth/register", authHandler.Register)
			public.POST("/auth/login", authHandler.Login)
			public.POST("/auth/refresh", authHandler.RefreshToken)
			public.POST("/auth/forgot-password", authHandler.ForgotPassword)
			public.POST("/auth/reset-password", authHandler.ResetPassword)
			public.GET("/auth/google", authHandler.GoogleLogin)
			public.GET("/auth/google/callback", authHandler.GoogleCallback)
			public.GET("/auth/facebook", authHandler.FacebookLogin)
			public.GET("/auth/facebook/callback", authHandler.FacebookCallback)

			// Products
			public.GET("/products", productHandler.GetProducts)
			public.GET("/products/trending", productHandler.GetTrendingProducts)
			public.GET("/products/flash-sale", productHandler.GetFlashSale)
			public.GET("/products/:slug", productHandler.GetProductBySlug)
			public.GET("/categories", categoryHandler.GetCategories)
			public.GET("/categories/:slug", categoryHandler.GetCategoryBySlug)
			public.GET("/brands", productHandler.GetBrands)

			// Trends
			public.GET("/trends/keywords", trendHandler.GetTrendingKeywords)
			public.GET("/trends/products", trendHandler.GetTrendingProducts)
			public.GET("/trends/products/:id/score", trendHandler.GetTrendScore)

			// Shipping (RajaOngkir)
			public.GET("/shipping/provinces", shippingHandler.GetProvinces)
			public.GET("/shipping/cities/:provinceId", shippingHandler.GetCities)
			public.POST("/shipping/calculate", shippingHandler.CalculateCost)

			// Vouchers
			public.POST("/vouchers/validate", voucherHandler.ValidateVoucher)
		}

		// Authenticated routes
		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware(jwtManager))
		{
			// Auth
			auth.POST("/auth/logout", authHandler.Logout)
			auth.GET("/auth/me", authHandler.GetMe)
			auth.PUT("/auth/profile", authHandler.UpdateProfile)
			auth.POST("/auth/change-password", authHandler.ChangePassword)

			// Cart
			auth.GET("/cart", cartHandler.GetCart)
			auth.POST("/cart/items", cartHandler.AddToCart)
			auth.PUT("/cart/items/:id", cartHandler.UpdateCartItem)
			auth.DELETE("/cart/items/:id", cartHandler.RemoveCartItem)
			auth.DELETE("/cart/clear", cartHandler.ClearCart)

			// Orders
			auth.POST("/orders", orderHandler.CreateOrder)
			auth.GET("/orders", orderHandler.GetOrders)
			auth.GET("/orders/:orderNumber", orderHandler.GetOrderByNumber)
			auth.POST("/orders/:id/cancel", orderHandler.CancelOrder)

			// Payments
			auth.POST("/payments/initiate", paymentHandler.InitiatePayment)
			auth.GET("/payments/:orderId/status", paymentHandler.GetPaymentStatus)

			// Reviews
			auth.POST("/reviews", reviewHandler.CreateReview)
			auth.PUT("/reviews/:id", reviewHandler.UpdateReview)
			auth.DELETE("/reviews/:id", reviewHandler.DeleteReview)
			auth.POST("/reviews/:id/helpful", reviewHandler.MarkHelpful)

			// Wishlist
			auth.GET("/wishlist", productHandler.GetWishlist)
			auth.POST("/wishlist", productHandler.AddToWishlist)
			auth.DELETE("/wishlist/:id", productHandler.RemoveFromWishlist)

			// Addresses
			auth.GET("/addresses", userHandler.GetAddresses)
			auth.POST("/addresses", userHandler.CreateAddress)
			auth.PUT("/addresses/:id", userHandler.UpdateAddress)
			auth.DELETE("/addresses/:id", userHandler.DeleteAddress)
		}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(jwtManager))
		admin.Use(middleware.AdminMiddleware())
		{
			// Products
			admin.POST("/products", productHandler.CreateProduct)
			admin.PUT("/products/:id", productHandler.UpdateProduct)
			admin.DELETE("/products/:id", productHandler.DeleteProduct)

			// Orders
			admin.GET("/orders", orderHandler.GetAllOrders)
			admin.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
			admin.PUT("/orders/:id/tracking", orderHandler.UpdateTracking)

			// Inventory
			admin.GET("/inventory", inventoryHandler.GetInventory)
			admin.PUT("/inventory/adjust", inventoryHandler.AdjustStock)
			admin.GET("/inventory/logs", inventoryHandler.GetInventoryLogs)

			// Trends
			admin.GET("/trends/dashboard", trendHandler.GetDashboard)
			admin.POST("/trends/refresh", trendHandler.RefreshTrends)

			// Reports
			admin.GET("/reports/sales", reportHandler.GetSalesReport)
			admin.GET("/reports/inventory", reportHandler.GetInventoryReport)
			admin.GET("/reports/customers", reportHandler.GetCustomerReport)
			admin.GET("/dashboard/summary", reportHandler.GetDashboardSummary)

			// Vouchers
			admin.POST("/vouchers", voucherHandler.CreateVoucher)
			admin.PUT("/vouchers/:id", voucherHandler.UpdateVoucher)
			admin.DELETE("/vouchers/:id", voucherHandler.DeleteVoucher)

			// Flash Sales
			admin.POST("/flash-sales", productHandler.CreateFlashSale)
			admin.PUT("/flash-sales/:id", productHandler.UpdateFlashSale)
			admin.DELETE("/flash-sales/:id", productHandler.DeleteFlashSale)
		}

		// Webhook (no auth, signature verification inside)
		v1.POST("/payments/webhook", webhookHandler.HandleMidtransWebhook)
	}
}
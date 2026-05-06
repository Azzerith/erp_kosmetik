package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"erp-cosmetics-backend/internal/models"
	"erp-cosmetics-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*models.Product, error)
	GetProductByID(ctx context.Context, id uint64) (*models.Product, error)
	GetProductBySlug(ctx context.Context, slug string) (*models.Product, error)
	UpdateProduct(ctx context.Context, id uint64, req *UpdateProductRequest) (*models.Product, error)
	DeleteProduct(ctx context.Context, id uint64) error
	ListProducts(ctx context.Context, req *ListProductsRequest) (*ListProductsResponse, error)
	GetTrendingProducts(ctx context.Context, limit int) ([]models.Product, error)
	GetBestSellers(ctx context.Context, limit int) ([]models.Product, error)
	GetFlashSale(ctx context.Context) ([]models.FlashSaleItem, error)
	UpdateStock(ctx context.Context, productID uint64, variantID *uint64, quantity int) error
	UpdateTrendScore(ctx context.Context, productID uint64, score float64, badge string) error
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	SKU         string  `json:"sku" binding:"required"`
	Description string  `json:"description"`
	ShortDesc   *string `json:"short_description"`
	CategoryID  uint64  `json:"category_id" binding:"required"`
	BrandID     *uint64 `json:"brand_id"`
	BasePrice   float64 `json:"base_price" binding:"required,min=0"`
	SalePrice   *float64 `json:"sale_price"`
	WeightGram  int     `json:"weight_gram"`
	Stock       int     `json:"stock"`
	IsActive    bool    `json:"is_active"`
	IsFeatured  bool    `json:"is_featured"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	ShortDesc   *string  `json:"short_description"`
	CategoryID  *uint64  `json:"category_id"`
	BrandID     *uint64  `json:"brand_id"`
	BasePrice   *float64 `json:"base_price"`
	SalePrice   *float64 `json:"sale_price"`
	WeightGram  *int     `json:"weight_gram"`
	Stock       *int     `json:"stock"`
	IsActive    *bool    `json:"is_active"`
	IsFeatured  *bool    `json:"is_featured"`
}

type ListProductsRequest struct {
	Page       int     `form:"page" binding:"min=1"`
	Limit      int     `form:"limit" binding:"min=1,max=100"`
	CategoryID *uint64 `form:"category_id"`
	BrandID    *uint64 `form:"brand_id"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	Search     *string `form:"search"`
	SortBy     string  `form:"sort_by"`
	SortOrder  string  `form:"sort_order"`
}

type ListProductsResponse struct {
	Products   []models.Product `json:"products"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	logger       *zap.Logger
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, logger *zap.Logger) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*models.Product, error) {
	// Generate slug
	slug := s.generateSlug(req.Name)

	// Create product
	product := &models.Product{
		SKU:          req.SKU,
		Name:         req.Name,
		Slug:         slug,
		Description:  req.Description,
		ShortDesc:    req.ShortDesc,
		CategoryID:   req.CategoryID,
		BrandID:      req.BrandID,
		BasePrice:    req.BasePrice,
		SalePrice:    req.SalePrice,
		WeightGram:   req.WeightGram,
		Stock:        req.Stock,
		IsActive:     req.IsActive,
		IsFeatured:   req.IsFeatured,
		TrendScore:   0,
		TrendBadge:   "none",
		TotalSold:    0,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		s.logger.Error("Failed to create product", zap.Error(err))
		return nil, err
	}

	return product, nil
}

func (s *productService) GetProductByID(ctx context.Context, id uint64) (*models.Product, error) {
	return s.productRepo.FindByID(ctx, id)
}

func (s *productService) GetProductBySlug(ctx context.Context, slug string) (*models.Product, error) {
	return s.productRepo.FindBySlug(ctx, slug)
}

func (s *productService) UpdateProduct(ctx context.Context, id uint64, req *UpdateProductRequest) (*models.Product, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		product.Name = *req.Name
		product.Slug = s.generateSlug(*req.Name)
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.ShortDesc != nil {
		product.ShortDesc = req.ShortDesc
	}
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.BrandID != nil {
		product.BrandID = req.BrandID
	}
	if req.BasePrice != nil {
		product.BasePrice = *req.BasePrice
	}
	if req.SalePrice != nil {
		product.SalePrice = req.SalePrice
	}
	if req.WeightGram != nil {
		product.WeightGram = *req.WeightGram
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uint64) error {
	return s.productRepo.Delete(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context, req *ListProductsRequest) (*ListProductsResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 12
	}

	offset := (req.Page - 1) * req.Limit

	filters := make(map[string]interface{})
	if req.CategoryID != nil {
		filters["category_id"] = *req.CategoryID
	}
	if req.BrandID != nil {
		filters["brand_id"] = *req.BrandID
	}
	if req.MinPrice != nil {
		filters["min_price"] = *req.MinPrice
	}
	if req.MaxPrice != nil {
		filters["max_price"] = *req.MaxPrice
	}
	if req.Search != nil && *req.Search != "" {
		filters["search"] = *req.Search
	}
	if req.SortBy != "" {
		filters["sort_by"] = req.SortBy
	}
	if req.SortOrder != "" {
		filters["sort_order"] = req.SortOrder
	}

	products, total, err := s.productRepo.List(ctx, offset, req.Limit, filters)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &ListProductsResponse{
		Products:   products,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *productService) GetTrendingProducts(ctx context.Context, limit int) ([]models.Product, error) {
	if limit == 0 {
		limit = 8
	}
	return s.productRepo.GetTrending(ctx, limit)
}

func (s *productService) GetBestSellers(ctx context.Context, limit int) ([]models.Product, error) {
	if limit == 0 {
		limit = 8
	}
	return s.productRepo.GetBestSellers(ctx, limit)
}

func (s *productService) GetFlashSale(ctx context.Context) ([]models.FlashSaleItem, error) {
	return s.productRepo.GetFlashSale(ctx)
}

func (s *productService) UpdateStock(ctx context.Context, productID uint64, variantID *uint64, quantity int) error {
	return s.productRepo.UpdateStock(ctx, productID, variantID, quantity)
}

func (s *productService) UpdateTrendScore(ctx context.Context, productID uint64, score float64, badge string) error {
	return s.productRepo.UpdateTrendScore(ctx, productID, score, badge)
}

func (s *productService) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "&", "and")
	slug = strings.ReplaceAll(slug, "--", "-")
	slug = strings.Trim(slug, "-")
	slug = slug + "-" + uuid.New().String()[:8]
	return slug
}
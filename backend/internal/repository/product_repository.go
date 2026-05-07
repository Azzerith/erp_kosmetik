package repository

import (
	"context"
	"erp-cosmetics-backend/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	FindByID(ctx context.Context, id uint64) (*models.Product, error)
	FindBySlug(ctx context.Context, slug string) (*models.Product, error)
	FindBySKU(ctx context.Context, sku string) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]models.Product, int64, error)
	GetTrending(ctx context.Context, limit int) ([]models.Product, error)
	GetBestSellers(ctx context.Context, limit int) ([]models.Product, error)
	GetFlashSale(ctx context.Context) ([]models.FlashSaleItem, error)
	UpdateStock(ctx context.Context, productID uint64, variantID *uint64, quantity int) error
	UpdateTrendScore(ctx context.Context, productID uint64, score float64, badge string) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) FindByID(ctx context.Context, id uint64) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Variants").
		Preload("Images").
		Preload("Tags").
		First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Brand").
		Preload("Variants").
		Preload("Images").
		Preload("Tags").
		Preload("Reviews", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_approved = ?", true).Order("created_at DESC").Limit(10)
		}).
		Where("slug = ? AND is_active = ?", slug, true).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
	return &product, err
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *productRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{}).Where("is_active = ?", true)

	// Apply filters
	if categoryID, ok := filters["category_id"]; ok {
		query = query.Where("category_id = ?", categoryID)
	}
	if brandID, ok := filters["brand_id"]; ok {
		query = query.Where("brand_id = ?", brandID)
	}
	if minPrice, ok := filters["min_price"]; ok {
		query = query.Where("COALESCE(sale_price, base_price) >= ?", minPrice)
	}
	if maxPrice, ok := filters["max_price"]; ok {
		query = query.Where("COALESCE(sale_price, base_price) <= ?", maxPrice)
	}
	if search, ok := filters["search"]; ok {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search.(string)+"%", "%"+search.(string)+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "created_at"
	sortOrder := "DESC"
	if val, ok := filters["sort_by"]; ok {
		sortBy = val.(string)
	}
	if val, ok := filters["sort_order"]; ok {
		sortOrder = val.(string)
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	err := query.Offset(offset).Limit(limit).
		Preload("Category").
		Preload("Brand").
		Preload("Images").
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) GetTrending(ctx context.Context, limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND trend_score > ?", true, 70).
		Order("trend_score DESC").
		Limit(limit).
		Preload("Category").
		Preload("Images").
		Find(&products).Error
	return products, err
}

func (r *productRepository) GetBestSellers(ctx context.Context, limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("total_sold DESC").
		Limit(limit).
		Preload("Category").
		Preload("Images").
		Find(&products).Error
	return products, err
}

func (r *productRepository) GetFlashSale(ctx context.Context) ([]models.FlashSaleItem, error) {
	var items []models.FlashSaleItem
	err := r.db.WithContext(ctx).
		Joins("JOIN flash_sales ON flash_sales.id = flash_sale_items.flash_sale_id").
		Where("flash_sales.is_active = ? AND flash_sales.start_time <= NOW() AND flash_sales.end_time >= NOW()", true).
		Preload("Product").
		Preload("Product.Images").
		Find(&items).Error
	return items, err
}

func (r *productRepository) UpdateStock(ctx context.Context, productID uint64, variantID *uint64, quantity int) error {
	if variantID == nil {
		return r.db.WithContext(ctx).Model(&models.Product{}).
			Where("id = ?", productID).
			Update("stock", gorm.Expr("stock + ?", quantity)).Error
	}
	return r.db.WithContext(ctx).Model(&models.ProductVariant{}).
		Where("id = ?", variantID).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *productRepository) UpdateTrendScore(ctx context.Context, productID uint64, score float64, badge string) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"trend_score": score,
			"trend_badge": badge,
		}).Error
}
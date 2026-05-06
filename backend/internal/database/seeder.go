package database

import (
	"log"
	"time"

	"erp-cosmetics-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunSeeder(db *gorm.DB) error {
	log.Println("Running database seeder...")

	// Seed users
	if err := seedUsers(db); err != nil {
		return err
	}

	// Seed categories
	if err := seedCategories(db); err != nil {
		return err
	}

	// Seed brands
	if err := seedBrands(db); err != nil {
		return err
	}

	// Seed products
	if err := seedProducts(db); err != nil {
		return err
	}

	// Seed vouchers
	if err := seedVouchers(db); err != nil {
		return err
	}

	log.Println("Database seeder completed successfully")
	return nil
}

func seedUsers(db *gorm.DB) error {
	// Check if users already exist
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Users already seeded, skipping...")
		return nil
	}

	// Hash password for admin
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("Admin@12345"), bcrypt.DefaultCost)
	adminPasswordHash := string(adminPassword)

	// Hash password for customer
	customerPassword, _ := bcrypt.GenerateFromPassword([]byte("Customer@123"), bcrypt.DefaultCost)
	customerPasswordHash := string(customerPassword)

	users := []models.User{
		{
			Email:        "superadmin@erpcosmetics.com",
			PasswordHash: &adminPasswordHash,
			Name:         "Super Admin",
			Phone:        stringPtr("081234567890"),
			Provider:     "local",
			Role:         "super_admin",
			IsActive:     true,
		},
		{
			Email:        "admin@erpcosmetics.com",
			PasswordHash: &adminPasswordHash,
			Name:         "Admin User",
			Phone:        stringPtr("081234567891"),
			Provider:     "local",
			Role:         "admin",
			IsActive:     true,
		},
		{
			Email:        "customer@example.com",
			PasswordHash: &customerPasswordHash,
			Name:         "Customer Demo",
			Phone:        stringPtr("081234567892"),
			Provider:     "local",
			Role:         "customer",
			IsActive:     true,
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d users", len(users))
	return nil
}

func seedCategories(db *gorm.DB) error {
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count > 0 {
		log.Println("Categories already seeded, skipping...")
		return nil
	}

	categories := []models.Category{
		{Name: "Skincare", Slug: "skincare", Level: 0, SortOrder: 1, IsActive: true},
		{Name: "Makeup", Slug: "makeup", Level: 0, SortOrder: 2, IsActive: true},
		{Name: "Herbal & Jamu", Slug: "herbal-jamu", Level: 0, SortOrder: 3, IsActive: true},
		{Name: "Haircare", Slug: "haircare", Level: 0, SortOrder: 4, IsActive: true},
		{Name: "Body Care", Slug: "body-care", Level: 0, SortOrder: 5, IsActive: true},
	}

	for i := range categories {
		if err := db.Create(&categories[i]).Error; err != nil {
			return err
		}
	}

	// Sub categories
	skincareID := categories[0].ID
	subCategories := []models.Category{
		{Name: "Facial Wash", Slug: "facial-wash", ParentID: &skincareID, Level: 1, SortOrder: 1, IsActive: true},
		{Name: "Moisturizer", Slug: "moisturizer", ParentID: &skincareID, Level: 1, SortOrder: 2, IsActive: true},
		{Name: "Sunscreen", Slug: "sunscreen", ParentID: &skincareID, Level: 1, SortOrder: 3, IsActive: true},
		{Name: "Serum", Slug: "serum", ParentID: &skincareID, Level: 1, SortOrder: 4, IsActive: true},
	}

	for _, subCat := range subCategories {
		if err := db.Create(&subCat).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d categories", len(categories)+len(subCategories))
	return nil
}

func seedBrands(db *gorm.DB) error {
	var count int64
	db.Model(&models.Brand{}).Count(&count)
	if count > 0 {
		log.Println("Brands already seeded, skipping...")
		return nil
	}

	brands := []models.Brand{
		{Name: "GlowLab", Slug: "glowlab", IsActive: true},
		{Name: "HerbalIndo", Slug: "herbalindo", IsActive: true},
		{Name: "DewyLab", Slug: "dewylab", IsActive: true},
		{Name: "SunShield", Slug: "sunshield", IsActive: true},
		{Name: "Lush Beauty", Slug: "lush-beauty", IsActive: true},
		{Name: "HairGro", Slug: "hairgro", IsActive: true},
		{Name: "BodyLove", Slug: "bodylove", IsActive: true},
	}

	for _, brand := range brands {
		if err := db.Create(&brand).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d brands", len(brands))
	return nil
}

func seedProducts(db *gorm.DB) error {
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count > 0 {
		log.Println("Products already seeded, skipping...")
		return nil
	}

	// Get category IDs
	var categories []models.Category
	db.Find(&categories)
	categoryMap := make(map[string]uint64)
	for _, cat := range categories {
		categoryMap[cat.Slug] = cat.ID
	}

	// Get brand IDs
	var brands []models.Brand
	db.Find(&brands)
	brandMap := make(map[string]uint64)
	for _, brand := range brands {
		brandMap[brand.Slug] = brand.ID
	}

	salePrice1 := 159000.0
	salePrice2 := 199000.0

	products := []models.Product{
		{
			SKU:             "SKIN-001",
			Name:            "Brightening Vitamin C Serum",
			Slug:            "brightening-vitamin-c-serum",
			Description:     "Serum vitamin C dengan formula stabil untuk mencerahkan kulit dan mengurangi hiperpigmentasi.",
			ShortDesc:       stringPtr("Serum vitamin C untuk kulit cerah merata"),
			CategoryID:      categoryMap["skincare"],
			BrandID:         uint64Ptr(brandMap["glowlab"]),
			BasePrice:       189000,
			SalePrice:       &salePrice1,
			WeightGram:      50,
			IsBPOMCertified: true,
			IsHalalCertified: true,
			Stock:           89,
			TrendScore:      95.5,
			TrendBadge:      "viral",
			IsActive:        true,
			IsFeatured:      true,
		},
		{
			SKU:             "HERB-001",
			Name:            "Jamu Kunyit Asam Herbal",
			Slug:            "jamu-kunyit-asam-herbal",
			Description:     "Jamu tradisional dari kunyit dan asam jawa untuk menjaga kebugaran tubuh.",
			ShortDesc:       stringPtr("Jamu tradisional kunyit asam"),
			CategoryID:      categoryMap["herbal-jamu"],
			BrandID:         uint64Ptr(brandMap["herbalindo"]),
			BasePrice:       85000,
			WeightGram:      200,
			IsBPOMCertified: true,
			IsHalalCertified: true,
			IsHerbal:        true,
			Stock:           45,
			TrendScore:      88.2,
			TrendBadge:      "trending",
			IsActive:        true,
			IsFeatured:      true,
		},
		{
			SKU:             "SKIN-002",
			Name:            "Hydrating Moisturizer",
			Slug:            "hydrating-moisturizer",
			Description:     "Pelembab ringan dengan hyaluronic acid untuk kulit yang terhidrasi sepanjang hari.",
			ShortDesc:       stringPtr("Pelembab dengan hyaluronic acid"),
			CategoryID:      categoryMap["skincare"],
			BrandID:         uint64Ptr(brandMap["dewylab"]),
			BasePrice:       125000,
			SalePrice:       &salePrice2,
			WeightGram:      50,
			IsBPOMCertified: true,
			IsHalalCertified: true,
			Stock:           156,
			TrendScore:      92.3,
			TrendBadge:      "best_seller",
			IsActive:        true,
			IsFeatured:      true,
		},
		{
			SKU:             "SKIN-003",
			Name:            "SPF 50 PA++++ Sunscreen",
			Slug:            "spf-50-sunscreen",
			Description:     "Tabir surya dengan tekstur ringan, tidak lengket, dan perlindungan maksimal.",
			ShortDesc:       stringPtr("Sunscreen SPF 50 PA++++"),
			CategoryID:      categoryMap["skincare"],
			BrandID:         uint64Ptr(brandMap["sunshield"]),
			BasePrice:       135000,
			SalePrice:       &salePrice1,
			WeightGram:      30,
			IsBPOMCertified: true,
			IsHalalCertified: true,
			Stock:           234,
			TrendScore:      96.7,
			TrendBadge:      "viral",
			IsActive:        true,
			IsFeatured:      true,
		},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			return err
		}

		// Add product image
		image := models.ProductImage{
			ProductID: product.ID,
			URL:       "https://picsum.photos/400/500",
			IsPrimary: true,
			SortOrder: 0,
		}
		db.Create(&image)

		// Add product tags
		tags := []string{"skincare", "beauty", "trending"}
		for _, tag := range tags {
			productTag := models.ProductTag{
				ProductID: product.ID,
				Tag:       tag,
				Weight:    1.0,
			}
			db.Create(&productTag)
		}
	}

	log.Printf("Seeded %d products", len(products))
	return nil
}

func seedVouchers(db *gorm.DB) error {
	var count int64
	db.Model(&models.Voucher{}).Count(&count)
	if count > 0 {
		log.Println("Vouchers already seeded, skipping...")
		return nil
	}

	now := time.Now()
	validFrom := now
	validUntil := now.AddDate(0, 3, 0)

	createdBy := uint64(1)

	vouchers := []models.Voucher{
		{
			Code:           "WELCOME10",
			Name:           "Welcome Discount",
			Description:    stringPtr("Diskon 10% untuk member baru"),
			Type:           "percentage",
			Value:          10,
			MinOrderAmount: 50000,
			ApplicableType: "all",
			UsageLimit:     intPtr(1000),
			UsagePerUser:   intPtr(1),
			ValidFrom:      validFrom,
			ValidUntil:     validUntil,
			IsActive:       true,
			CreatedBy:      createdBy,
		},
		{
			Code:             "FLASH20",
			Name:             "Flash Sale",
			Description:      stringPtr("Diskon Rp20.000 untuk pembelian minimal Rp100.000"),
			Type:             "fixed_amount",
			Value:            20000,
			MinOrderAmount:   100000,
			ApplicableType:   "all",
			UsageLimit:       intPtr(500),
			UsagePerUser:     intPtr(1),
			ValidFrom:        validFrom,
			ValidUntil:       validUntil,
			IsActive:         true,
			CreatedBy:        createdBy,
		},
		{
			Code:             "FREESHIP",
			Name:             "Free Shipping",
			Description:      stringPtr("Gratis ongkir minimal belanja Rp150.000"),
			Type:             "free_shipping",
			Value:            0,
			MinOrderAmount:   150000,
			ApplicableType:   "all",
			UsageLimit:       intPtr(200),
			UsagePerUser:     intPtr(2),
			ValidFrom:        validFrom,
			ValidUntil:       validUntil,
			IsActive:         true,
			CreatedBy:        createdBy,
		},
	}

	for _, voucher := range vouchers {
		if err := db.Create(&voucher).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d vouchers", len(vouchers))
	return nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func uint64Ptr(u uint64) *uint64 {
	return &u
}

func intPtr(i int) *int {
	return &i
}
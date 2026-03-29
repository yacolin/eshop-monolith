package service

import (
	"context"

	"eshop-monolith/internal/domain/category"
	"eshop-monolith/internal/domain/product"
	"eshop-monolith/internal/domain/shared"
	"eshop-monolith/internal/eventbus"

	"gorm.io/gorm"
)

// ProductService 产品服务
type ProductService struct {
	repo product.Repository
	bus  *eventbus.Bus
	db   *gorm.DB
}

// NewProductService 创建产品服务
func NewProductService(repo product.Repository, bus *eventbus.Bus, db *gorm.DB) *ProductService {
	return &ProductService{
		repo: repo,
		bus:  bus,
		db:   db,
	}
}

// CreateProductRequest 创建产品请求
type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       int64   `json:"price"`
	SKU         string  `json:"sku"`
	CategoryIDs []int64 `json:"category_ids"`
}

// UpdateProductRequest 更新产品请求
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       int64   `json:"price"`
	CategoryIDs []int64 `json:"category_ids"`
}

// CreateProduct 创建产品
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*product.Product, error) {
	// 创建产品
	newProduct := &product.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
	}

	// 保存产品
	if err := s.repo.Create(ctx, newProduct); err != nil {
		return nil, err
	}

	// 关联分类（通过中间表）
	for _, categoryID := range req.CategoryIDs {
		pc := &shared.ProductCategory{
			ProductID:  newProduct.ID,
			CategoryID: categoryID,
		}
		if err := s.db.WithContext(ctx).Create(pc).Error; err != nil {
			return nil, err
		}
	}

	// 发布产品创建事件
	categoryIDValue := int64(0)
	if len(req.CategoryIDs) > 0 {
		categoryIDValue = req.CategoryIDs[0] // 选择第一个分类ID作为事件中的CategoryID
	}
	s.bus.Publish(product.ProductCreatedEvent{
		ProductID:  newProduct.ID,
		Name:       newProduct.Name,
		Price:      newProduct.Price,
		CategoryID: categoryIDValue,
	})

	return newProduct, nil
}

// GetProductByID 根据ID获取产品
func (s *ProductService) GetProductByID(ctx context.Context, id int64) (*product.Product, error) {
	return s.repo.FindByID(ctx, id)
}

// GetProductBySKU 根据SKU获取产品
func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (*product.Product, error) {
	return s.repo.FindBySKU(ctx, sku)
}

// GetProductWithCategories 获取产品及其关联的分类
func (s *ProductService) GetProductWithCategories(ctx context.Context, productID int64) (*product.Product, []category.Category, error) {
	// 获取产品
	prod, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return nil, nil, err
	}

	// 通过中间表查询关联的分类
	var categories []category.Category
	if err := s.db.WithContext(ctx).Table("categories").
		Joins("JOIN product_categories ON categories.id = product_categories.category_id").
		Where("product_categories.product_id = ?", productID).
		Find(&categories).Error; err != nil {
		return nil, nil, err
	}

	return prod, categories, nil
}

// ListProductsByCategory 根据分类列出产品
func (s *ProductService) ListProductsByCategory(ctx context.Context, categoryID int64, page, pageSize int) ([]product.Product, int64, error) {
	return s.repo.ListByCategory(ctx, categoryID, page, pageSize)
}

// ListAllProducts 列出所有产品
func (s *ProductService) ListAllProducts(ctx context.Context, page, pageSize int) ([]product.Product, int64, error) {
	return s.repo.ListAll(ctx, page, pageSize)
}

// UpdateProduct 更新产品
func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req *UpdateProductRequest) (*product.Product, error) {
	// 获取产品
	existingProduct, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新产品信息
	existingProduct.Name = req.Name
	existingProduct.Description = req.Description
	existingProduct.Price = req.Price

	// 保存产品
	if err := s.repo.Update(ctx, existingProduct); err != nil {
		return nil, err
	}

	// 删除现有的分类关联
	if err := s.db.Where("product_id = ?", id).Delete(&shared.ProductCategory{}).Error; err != nil {
		return nil, err
	}

	// 重新关联分类（通过中间表）
	for _, categoryID := range req.CategoryIDs {
		pc := &shared.ProductCategory{
			ProductID:  existingProduct.ID,
			CategoryID: categoryID,
		}
		if err := s.db.WithContext(ctx).Create(pc).Error; err != nil {
			return nil, err
		}
	}

	// 发布产品更新事件
	categoryIDValue := int64(0)
	if len(req.CategoryIDs) > 0 {
		categoryIDValue = req.CategoryIDs[0] // 选择第一个分类ID作为事件中的CategoryID
	}
	s.bus.Publish(product.ProductUpdatedEvent{
		ProductID:  existingProduct.ID,
		Name:       existingProduct.Name,
		Price:      existingProduct.Price,
		CategoryID: categoryIDValue,
	})

	return existingProduct, nil
}

// DeleteProduct 删除产品
func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	// 获取产品
	existingProduct, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 删除产品的分类关联
	if err := s.db.WithContext(ctx).Where("product_id = ?", id).Delete(&shared.ProductCategory{}).Error; err != nil {
		return err
	}

	// 删除产品
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 发布产品删除事件
	// 由于产品已删除，我们无法获取其分类关联，所以设置CategoryID为0
	s.bus.Publish(product.ProductDeletedEvent{
		ProductID:  existingProduct.ID,
		Name:       existingProduct.Name,
		CategoryID: 0,
	})

	return nil
}

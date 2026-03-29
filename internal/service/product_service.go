package service

import (
	"context"

	"eshop-monolith/internal/domain/product"
	"eshop-monolith/internal/eventbus"
)

// ProductService 产品服务
type ProductService struct {
	repo product.Repository
	bus  *eventbus.Bus
}

// NewProductService 创建产品服务
func NewProductService(repo product.Repository, bus *eventbus.Bus) *ProductService {
	return &ProductService{
		repo: repo,
		bus:  bus,
	}
}

// CreateProductRequest 创建产品请求
type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	SKU         string `json:"sku"`
	CategoryID  int64  `json:"category_id"`
}

// UpdateProductRequest 更新产品请求
type UpdateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	CategoryID  int64  `json:"category_id"`
}

// CreateProduct 创建产品
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*product.Product, error) {
	// 创建产品
	newProduct := &product.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
		CategoryID:  req.CategoryID,
	}

	// 保存产品
	if err := s.repo.Create(ctx, newProduct); err != nil {
		return nil, err
	}

	// 发布产品创建事件
	s.bus.Publish(product.ProductCreatedEvent{
		ProductID:  newProduct.ID,
		Name:       newProduct.Name,
		Price:      newProduct.Price,
		CategoryID: newProduct.CategoryID,
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
	existingProduct.CategoryID = req.CategoryID

	// 保存产品
	if err := s.repo.Update(ctx, existingProduct); err != nil {
		return nil, err
	}

	// 发布产品更新事件
	s.bus.Publish(product.ProductUpdatedEvent{
		ProductID:  existingProduct.ID,
		Name:       existingProduct.Name,
		Price:      existingProduct.Price,
		CategoryID: existingProduct.CategoryID,
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

	// 删除产品
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 发布产品删除事件
	s.bus.Publish(product.ProductDeletedEvent{
		ProductID:  existingProduct.ID,
		Name:       existingProduct.Name,
		CategoryID: existingProduct.CategoryID,
	})

	return nil
}
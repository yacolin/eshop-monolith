package service

import (
	"context"
	"strings"

	"eshop-monolith/internal/infra/domain/shared"
	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/events"
	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

// ProductService 产品服务
type ProductService struct {
	repo          repositories.IproductRepository
	inventoryRepo repositories.IinventoryRepository
	bus           *eventbus.Bus
	db            *gorm.DB
}

// NewProductService 创建产品服务
func NewProductService(
	repo repositories.IproductRepository,
	inventoryRepo repositories.IinventoryRepository,
	bus *eventbus.Bus,
	db *gorm.DB,
) *ProductService {
	return &ProductService{
		repo:          repo,
		inventoryRepo: inventoryRepo,
		bus:           bus,
		db:            db,
	}
}

// GetProductWithInventory 获取产品详情（聚合库存信息），使用 goroutine + channel 并发查询以降低延迟
func (s *ProductService) GetProductWithInventory(ctx context.Context, id int64) (*dto.ProductDetailDTO, error) {
	type prodResult struct {
		product *models.Product
		err     error
	}
	type invResult struct {
		inventory *models.Inventory
		err       error
	}

	prodCh := make(chan prodResult, 1)
	invCh := make(chan invResult, 1)

	// 并发查询产品信息
	go func() {
		p, err := s.repo.FindByID(ctx, id)
		prodCh <- prodResult{product: p, err: err}
	}()

	// 并发查询库存信息
	go func() {
		inv, err := s.inventoryRepo.FindInventoryByProductID(ctx, id)
		invCh <- invResult{inventory: inv, err: err}
	}()

	// 收集产品结果
	pr := <-prodCh
	if pr.err != nil {
		return nil, pr.err
	}

	// 收集库存结果（库存不存在时不阻塞流程，库存字段使用零值）
	ir := <-invCh
	detail := &dto.ProductDetailDTO{
		ID:          pr.product.ID,
		Name:        pr.product.Name,
		Description: pr.product.Description,
		Price:       pr.product.Price,
		SKU:         pr.product.SKU,
		CreatedAt:   pr.product.CreatedAt,
		UpdatedAt:   pr.product.UpdatedAt,
	}
	if ir.err == nil && ir.inventory != nil {
		detail.Quantity = ir.inventory.Quantity
		detail.Status = ir.inventory.Status
		detail.Reserved = ir.inventory.Reserved
		detail.Threshold = ir.inventory.Threshold
	}

	return detail, nil
}

// CreateProduct 创建产品
func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductDTO) (*models.Product, error) {
	newProduct := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 保存产品
		if err := tx.Create(newProduct).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				return errcode.ErrDuplicateSKU
			}
			return err
		}
		// 关联分类
		for _, categoryID := range req.CategoryIDs {
			pc := &shared.ProductCategory{
				ProductID:  newProduct.ID,
				CategoryID: categoryID,
			}
			if err := tx.Create(pc).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 发布产品创建事件（事务外）
	categoryIDValue := int64(0)
	if len(req.CategoryIDs) > 0 {
		categoryIDValue = req.CategoryIDs[0]
	}
	s.bus.Publish(events.ProductCreatedEvent{
		ProductID:  newProduct.ID,
		Name:       newProduct.Name,
		Price:      newProduct.Price,
		CategoryID: categoryIDValue,
	})

	return newProduct, nil
}

// GetProductByID 根据ID获取产品
func (s *ProductService) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.FindByID(ctx, id)
}

// GetProductBySKU 根据SKU获取产品
func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (*models.Product, error) {
	return s.repo.FindBySKU(ctx, sku)
}

// GetProductWithCategories 获取产品及其关联的分类
func (s *ProductService) GetProductWithCategories(ctx context.Context, productID int64) (*models.Product, []models.Category, error) {
	// 获取产品
	prod, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return nil, nil, err
	}

	// 通过中间表查询关联的分类
	var categories []models.Category
	if err := s.db.WithContext(ctx).Table("categories").
		Joins("JOIN product_categories ON categories.id = product_categories.category_id").
		Where("product_categories.product_id = ?", productID).
		Find(&categories).Error; err != nil {
		return nil, nil, err
	}

	return prod, categories, nil
}

// ListProductsByCategory 根据分类列出产品
func (s *ProductService) ListProductsByCategory(ctx context.Context, categoryID int64, page, pageSize int) ([]models.Product, int64, error) {
	return s.repo.ListByCategory(ctx, categoryID, page, pageSize)
}

// UpdateProduct 更新产品
func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req *dto.UpdateProductDTO) (*models.Product, error) {
	// 获取产品（事务外查询）
	existingProduct, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新产品信息
	if req.Name != nil {
		existingProduct.Name = *req.Name
	}
	if req.Description != nil {
		existingProduct.Description = *req.Description
	}
	if req.Price != nil {
		existingProduct.Price = *req.Price
	}

	// 事务内保存产品和分类关联
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(existingProduct).Error; err != nil {
			return err
		}
		if req.CategoryIDs != nil {
			if err := tx.Where("product_id = ?", id).Delete(&shared.ProductCategory{}).Error; err != nil {
				return err
			}
			for _, categoryID := range req.CategoryIDs {
				pc := &shared.ProductCategory{
					ProductID:  existingProduct.ID,
					CategoryID: categoryID,
				}
				if err := tx.Create(pc).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 发布产品更新事件（事务外）
	categoryIDValue := int64(0)
	if len(req.CategoryIDs) > 0 {
		categoryIDValue = req.CategoryIDs[0]
	}
	s.bus.Publish(events.ProductUpdatedEvent{
		ProductID:  existingProduct.ID,
		Name:       existingProduct.Name,
		Price:      existingProduct.Price,
		CategoryID: categoryIDValue,
	})

	return existingProduct, nil
}

// DeleteProduct 删除产品
func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	// 获取产品（事务外查询）
	existingProduct, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 事务内删除分类关联和产品
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", id).Delete(&shared.ProductCategory{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Product{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 发布产品删除事件（事务外）
	s.bus.Publish(events.ProductDeletedEvent{
		ProductID:  existingProduct.ID,
		Name:       existingProduct.Name,
		CategoryID: 0,
	})

	return nil
}

// ListAllProducts 列出所有产品
func (s *ProductService) ListAllProducts(ctx context.Context, page, pageSize int) ([]models.Product, int64, error) {
	return s.repo.ListAll(ctx, page, pageSize)
}

func (s *ProductService) ListProducts(ctx context.Context, q dto.ProductListQuery) (*dto.ProductListResult, error) {
	offset := (q.Page - 1) * q.Size
	list, err := s.repo.ListProducts(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountProducts(ctx, q)
	if err != nil {
		return nil, err
	}

	return &dto.ProductListResult{
		List:  list,
		Total: total,
	}, nil
}

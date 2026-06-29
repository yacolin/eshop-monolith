package service

import "eshop-monolith/internal/inventory/domain/models"

// InventoryService 旧库存服务桩（保留兼容旧模块引用）
type InventoryService struct{}

func NewInventoryService(args ...interface{}) *InventoryService {
	return &InventoryService{}
}

func (s *InventoryService) GetInventoryBySkuID(ctx interface{}, skuID int64) (*models.Inventory, error) {
	return &models.Inventory{Quantity: 999999, Reserved: 0}, nil
}

// ProductService 旧产品服务桩
type ProductService struct{}

func NewProductService(args ...interface{}) *ProductService { return &ProductService{} }

func (s *ProductService) GetProductByID(ctx interface{}, id int64) (*models.Product, error) {
	return &models.Product{}, nil
}

// SkuService 旧SKU服务桩
type SkuService struct{}

func NewSkuService(args ...interface{}) *SkuService { return &SkuService{} }

func (s *SkuService) GetSku(ctx interface{}, id int64) (*models.Sku, error) {
	return &models.Sku{}, nil
}

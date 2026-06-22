package service

import (
	"context"

	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/events"

	"gorm.io/gorm"
)

type SkuService struct {
	skuRepo     repositories.IskuRepository
	productRepo repositories.IproductRepository
	bus         *eventbus.Bus
	db          *gorm.DB
}

func NewSkuService(skuRepo repositories.IskuRepository, productRepo repositories.IproductRepository, bus *eventbus.Bus, db *gorm.DB) *SkuService {
	return &SkuService{skuRepo: skuRepo, productRepo: productRepo, bus: bus, db: db}
}

// syncProductMinPrice 重新计算 Product 的最低 SKU 价格
func (s *SkuService) syncProductMinPrice(ctx context.Context, productID int64) error {
	type result struct{ Price int64 }
	var r result
	if err := s.db.WithContext(ctx).Model(&models.Sku{}).
		Select("MIN(price) as price").
		Where("product_id = ?", productID).Scan(&r).Error; err != nil {
		return err
	}
	return s.db.WithContext(ctx).Model(&models.Product{}).
		Where("id = ?", productID).Update("min_price", r.Price).Error
}

func (s *SkuService) CreateSku(ctx context.Context, req *dto.CreateSkuDTO) (*models.Sku, error) {
	sku := &models.Sku{
		ProductID: req.ProductID, Name: req.Name, Price: req.Price,
		SKUCode: req.SKUCode, Image: req.Image, Spec: req.Spec,
	}
	if err := s.skuRepo.Create(ctx, sku); err != nil {
		return nil, err
	}
	s.syncProductMinPrice(ctx, req.ProductID)

	s.bus.Publish(events.SkuCreatedEvent{SkuID: sku.ID, ProductID: sku.ProductID, Price: sku.Price})
	return sku, nil
}

func (s *SkuService) GetSku(ctx context.Context, id int64) (*models.Sku, error) {
	return s.skuRepo.FindByID(ctx, id)
}

func (s *SkuService) ListByProductID(ctx context.Context, productID int64) ([]models.Sku, error) {
	return s.skuRepo.FindByProductID(ctx, productID)
}

func (s *SkuService) UpdateSku(ctx context.Context, id int64, req *dto.UpdateSkuDTO) (*models.Sku, error) {
	sku, err := s.skuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil { sku.Name = *req.Name }
	if req.Price != nil { sku.Price = *req.Price }
	if req.SKUCode != nil { sku.SKUCode = *req.SKUCode }
	if req.Image != nil { sku.Image = *req.Image }
	if req.Spec != nil { sku.Spec = req.Spec }

	if err := s.skuRepo.Update(ctx, sku); err != nil {
		return nil, err
	}
	s.syncProductMinPrice(ctx, sku.ProductID)

	s.bus.Publish(events.SkuUpdatedEvent{SkuID: sku.ID, ProductID: sku.ProductID, Price: sku.Price})
	return sku, nil
}

func (s *SkuService) DeleteSku(ctx context.Context, id int64) error {
	sku, err := s.skuRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.skuRepo.Delete(ctx, id); err != nil {
		return err
	}
	s.syncProductMinPrice(ctx, sku.ProductID)

	s.bus.Publish(events.SkuDeletedEvent{SkuID: sku.ID, ProductID: sku.ProductID})
	return nil
}

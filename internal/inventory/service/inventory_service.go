package service

import (
	"context"
	"fmt"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/events"
	"eshop-monolith/pkg/errcode"
)

// InventoryService 库存服务
type InventoryService struct {
	repo        repositories.IinventoryRepository
	skuRepo     repositories.IskuRepository
	productRepo repositories.IproductRepository
	rabbit      *rabbitmq.Client
}

// NewInventoryService 创建库存服务
func NewInventoryService(repo repositories.IinventoryRepository, skuRepo repositories.IskuRepository, productRepo repositories.IproductRepository, rabbit *rabbitmq.Client) *InventoryService {
	return &InventoryService{
		repo:        repo,
		skuRepo:     skuRepo,
		productRepo: productRepo,
		rabbit:      rabbit,
	}
}

// CreateInventoryRequest 创建库存请求
type CreateInventoryRequest struct {
	SkuID    int64 `json:"sku_id"`
	Quantity int   `json:"quantity"`
}

// UpdateInventoryRequest 更新库存请求
type UpdateInventoryRequest struct {
	Quantity int `json:"quantity"`
}

// BatchCreateInventory 批量创建库存
func (s *InventoryService) BatchCreateInventory(ctx context.Context, req *dto.BatchCreateInventoryDTO) ([]*models.Inventory, error) {
	inventories := make([]*models.Inventory, len(req.SkuIDs))
	for i, skuID := range req.SkuIDs {
		inventories[i] = &models.Inventory{
			SkuID:     skuID,
			Quantity:  req.Quantity,
			Threshold: req.Threshold,
		}
	}

	if err := s.repo.BatchCreateInventory(ctx, inventories); err != nil {
		return nil, err
	}

	for _, inv := range inventories {
		if inv.Quantity > 0 && inv.Quantity <= inv.Threshold {
			s.rabbit.Publish(ctx, events.InventoryLowEvent{
				SkuID:     fmt.Sprintf("%d", inv.SkuID),
				Quantity:  inv.Quantity,
				Threshold: inv.Threshold,
			})
		}
	}

	return inventories, nil
}

// CreateInventory 创建库存
func (s *InventoryService) CreateInventory(ctx context.Context, req *dto.CreateInventoryDTO) (*models.Inventory, error) {

	// 创建库存
	inv := &models.Inventory{
		SkuID:     req.SkuID,
		Quantity:  req.Quantity,
		Threshold: req.Threshold,
	}

	inv.UpdateStatus() // 根据数量和阈值设置初始状态

	// 保存库存
	if err := s.repo.CreateInventory(ctx, inv); err != nil {
		return nil, err
	}

	return inv, nil
}

// GetInventoryBySkuID 根据SKU ID获取库存
func (s *InventoryService) GetInventoryBySkuID(ctx context.Context, skuID int64) (*models.Inventory, error) {
	inventory, err := s.repo.FindInventoryBySkuID(ctx, skuID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return inventory, nil
}

// UpdateInventory 更新库存
func (s *InventoryService) UpdateInventory(ctx context.Context, skuID int64, req *dto.UpdateInventoryDTO) (*models.Inventory, error) {
	// 获取库存
	inv, err := s.repo.FindInventoryBySkuID(ctx, skuID)
	if err != nil {
		return nil, err
	}

	if req.Quantity != nil {
		inv.Quantity = *req.Quantity
	}
	if req.Threshold != nil {
		inv.Threshold = *req.Threshold
	}
	if req.Reserved != nil {
		inv.Reserved = *req.Reserved
	}

	// 保存库存
	if err := s.repo.UpdateInventory(ctx, inv); err != nil {
		return nil, err
	}

	// 检查低库存并发布事件
	s.checkLowInventory(ctx, skuID)

	return inv, nil
}

// ReserveInventory 预占库存
func (s *InventoryService) ReserveInventory(ctx context.Context, req *dto.ReserveInventoryDTO) error {
	if err := s.repo.ReserveInventory(ctx, req.SkuID, req.Quantity); err != nil {
		return err
	}

	// 发布库存预占事件
	s.rabbit.Publish(ctx,events.InventoryReservedEvent{
		SkuID:    fmt.Sprintf("%d", req.SkuID),
		Quantity: req.Quantity,
	})

	// 检查低库存并发布事件
	s.checkLowInventory(ctx, req.SkuID)

	return nil
}

// ReleaseInventory 释放库存
func (s *InventoryService) ReleaseInventory(ctx context.Context, req *dto.ReleaseInventoryDTO) error {
	if err := s.repo.ReleaseInventory(ctx, req.SkuID, req.Quantity); err != nil {
		return err
	}

	// 发布库存释放事件
	s.rabbit.Publish(ctx,events.InventoryReleasedEvent{
		SkuID:    fmt.Sprintf("%d", req.SkuID),
		Quantity: req.Quantity,
	})

	return nil
}

// ListInventoriesEnriched 列出所有库存（含 SKU 名称和产品信息）
func (s *InventoryService) ListInventoriesEnriched(ctx context.Context, q dto.InventoryListQuery) (*dto.InventoryEnrichedResult, error) {
	offset := (q.Page - 1) * q.Size
	list, err := s.repo.ListInventories(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountInventories(ctx, q)
	if err != nil {
		return nil, err
	}

	// 收集所有 SkuID 批量查询 SKU 信息
	skuIDs := make([]int64, len(list))
	for i, inv := range list {
		skuIDs[i] = inv.SkuID
	}

	skuMap := make(map[int64]*models.Sku, len(list))
	if len(skuIDs) > 0 {
		skus, err := s.skuRepo.FindByIDs(ctx, skuIDs)
		if err == nil {
			for _, sku := range skus {
				skuMap[sku.ID] = &sku
			}
		}
	}

	// 收集所有 ProductID 批量查询产品名称
	productIDs := make(map[int64]struct{}, len(skuMap))
	for _, sku := range skuMap {
		productIDs[sku.ProductID] = struct{}{}
	}

	productNameMap := make(map[int64]string, len(productIDs))
	if len(productIDs) > 0 {
		ids := make([]int64, 0, len(productIDs))
		for pid := range productIDs {
			ids = append(ids, pid)
		}
		products, err := s.productRepo.FindByIDs(ctx, ids)
		if err == nil {
			for _, p := range products {
				productNameMap[p.ID] = p.Name
			}
		}
	}

	enrichedList := make([]dto.InventoryEnrichedItem, len(list))
	for i, inv := range list {
		sku := skuMap[inv.SkuID]
		enrichedList[i] = dto.InventoryEnrichedItem{
			ID:          inv.ID,
			SkuID:       inv.SkuID,
			SkuName:     sku.Name,
			SkuCode:     sku.SKUCode,
			ProductID:   sku.ProductID,
			ProductName: productNameMap[sku.ProductID],
			Quantity:    inv.Quantity,
			Status:      inv.Status,
			Reserved:    inv.Reserved,
			Threshold:   inv.Threshold,
			CreatedAt:   inv.CreatedAt,
			UpdatedAt:   inv.UpdatedAt,
		}
	}

	return &dto.InventoryEnrichedResult{
		Total: total,
		List:  enrichedList,
	}, nil
}

// CheckInventory 检查库存
func (s *InventoryService) CheckInventory(ctx context.Context, skuID int64) (*models.Inventory, error) {
	return s.repo.FindInventoryBySkuID(ctx, skuID)
}

// ListInventories 列出所有库存
func (s *InventoryService) ListInventories(ctx context.Context, q dto.InventoryListQuery) (*dto.InventoryListResult, error) {
	offset := (q.Page - 1) * q.Size
	list, err := s.repo.ListInventories(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountInventories(ctx, q)
	if err != nil {
		return nil, err
	}

	return &dto.InventoryListResult{
		List:  list,
		Total: total,
	}, nil
}

// checkLowInventory 检查库存是否低于阈值，是则发布低库存事件
func (s *InventoryService) checkLowInventory(ctx context.Context, skuID int64) {
	inv, err := s.repo.FindInventoryBySkuID(ctx, skuID)
	if err != nil || inv == nil {
		return
	}
	if inv.Quantity > 0 && inv.Quantity <= inv.Threshold {
		s.rabbit.Publish(ctx, events.InventoryLowEvent{
			SkuID:     fmt.Sprintf("%d", skuID),
			Quantity:  inv.Quantity,
			Threshold: inv.Threshold,
		})
	}
}

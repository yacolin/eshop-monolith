package service

import (
	"context"
	"fmt"

	"eshop-monolith/internal/eventbus"
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/events"
	"eshop-monolith/pkg/errcode"
)

// InventoryService 库存服务
type InventoryService struct {
	repo repositories.IinventoryRepository
	bus  *eventbus.Bus
}

// NewInventoryService 创建库存服务
func NewInventoryService(repo repositories.IinventoryRepository, bus *eventbus.Bus) *InventoryService {
	return &InventoryService{
		repo: repo,
		bus:  bus,
	}
}

// CreateInventoryRequest 创建库存请求
type CreateInventoryRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// UpdateInventoryRequest 更新库存请求
type UpdateInventoryRequest struct {
	Quantity int `json:"quantity"`
}

// CreateInventory 创建库存
func (s *InventoryService) CreateInventory(ctx context.Context, req *dto.CreateInventoryDTO) (*models.Inventory, error) {

	// 创建库存
	inv := &models.Inventory{
		ProductID: req.ProductID,
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

// GetInventoryByProductID 根据产品ID获取库存
func (s *InventoryService) GetInventoryByProductID(ctx context.Context, productID int64) (*models.Inventory, error) {
	inventory, err := s.repo.FindInventoryByProductID(ctx, productID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return inventory, nil
}

// UpdateInventory 更新库存
func (s *InventoryService) UpdateInventory(ctx context.Context, productID int64, req *dto.UpdateInventoryDTO) (*models.Inventory, error) {
	// 获取库存
	inv, err := s.repo.FindInventoryByProductID(ctx, productID)
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

	return inv, nil
}

// ReserveInventory 预占库存
func (s *InventoryService) ReserveInventory(ctx context.Context, req *dto.ReserveInventoryDTO) error {
	if err := s.repo.ReserveInventory(ctx, req.ProductID, req.Quantity); err != nil {
		return err
	}

	// 发布库存预占事件
	s.bus.Publish(events.InventoryReservedEvent{
		ProductID: fmt.Sprintf("%d", req.ProductID),
		Quantity:  req.Quantity,
	})

	return nil
}

// ReleaseInventory 释放库存
func (s *InventoryService) ReleaseInventory(ctx context.Context, req *dto.ReleaseInventoryDTO) error {
	if err := s.repo.ReleaseInventory(ctx, req.ProductID, req.Quantity); err != nil {
		return err
	}

	// 发布库存释放事件
	s.bus.Publish(events.InventoryReleasedEvent{
		ProductID: fmt.Sprintf("%d", req.ProductID),
		Quantity:  req.Quantity,
	})

	return nil
}

// CheckInventory 检查库存
func (s *InventoryService) CheckInventory(ctx context.Context, productID int64) (*models.Inventory, error) {
	return s.repo.FindInventoryByProductID(ctx, productID)
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

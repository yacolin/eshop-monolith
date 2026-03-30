package service

import (
	"context"
	"fmt"

	"eshop-monolith/internal/domain/inventory"
	"eshop-monolith/internal/eventbus"
	"eshop-monolith/internal/pkg/errcode"
)

// InventoryService 库存服务
type InventoryService struct {
	repo inventory.Repository
	bus  *eventbus.Bus
}

// NewInventoryService 创建库存服务
func NewInventoryService(repo inventory.Repository, bus *eventbus.Bus) *InventoryService {
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
func (s *InventoryService) CreateInventory(ctx context.Context, req *inventory.CreateInventoryDTO) (*inventory.Inventory, error) {

	// 创建库存
	inv := &inventory.Inventory{
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
func (s *InventoryService) GetInventoryByProductID(ctx context.Context, productID int64) (*inventory.Inventory, error) {
	inventory, err := s.repo.FindInventoryByProductID(ctx, productID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return inventory, nil
}

// UpdateInventory 更新库存
func (s *InventoryService) UpdateInventory(ctx context.Context, productID int64, req *inventory.UpdateInventoryDTO) (*inventory.Inventory, error) {
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
func (s *InventoryService) ReserveInventory(ctx context.Context, req *inventory.ReserveInventoryDTO) error {
	if err := s.repo.ReserveInventory(ctx, req.ProductID, req.Quantity); err != nil {
		return err
	}

	// 发布库存预占事件
	s.bus.Publish(inventory.InventoryReservedEvent{
		ProductID: fmt.Sprintf("%d", req.ProductID),
		Quantity:  req.Quantity,
	})

	return nil
}

// ReleaseInventory 释放库存
func (s *InventoryService) ReleaseInventory(ctx context.Context, req *inventory.ReleaseInventoryDTO) error {
	if err := s.repo.ReleaseInventory(ctx, req.ProductID, req.Quantity); err != nil {
		return err
	}

	// 发布库存释放事件
	s.bus.Publish(inventory.InventoryReleasedEvent{
		ProductID: fmt.Sprintf("%d", req.ProductID),
		Quantity:  req.Quantity,
	})

	return nil
}

// CheckInventory 检查库存
func (s *InventoryService) CheckInventory(ctx context.Context, productID int64) (*inventory.Inventory, error) {
	return s.repo.FindInventoryByProductID(ctx, productID)
}

// ListInventories 列出所有库存
func (s *InventoryService) ListInventories(ctx context.Context, q inventory.InventoryListQuery) (*inventory.InventoryListResult, error) {
	offset := (q.Page - 1) * q.Size
	list, err := s.repo.ListInventories(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountInventories(ctx, q)
	if err != nil {
		return nil, err
	}

	return &inventory.InventoryListResult{
		List:  list,
		Total: total,
	}, nil
}

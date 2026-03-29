package service

import (
	"context"
	"fmt"

	"eshop-monolith/internal/domain/inventory"
	"eshop-monolith/internal/eventbus"
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
func (s *InventoryService) CreateInventory(ctx context.Context, req *CreateInventoryRequest) (*inventory.Inventory, error) {
	// 创建库存
	inv := &inventory.Inventory{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	// 保存库存
	if err := s.repo.CreateInventory(ctx, inv); err != nil {
		return nil, err
	}

	return inv, nil
}

// GetInventoryByProductID 根据产品ID获取库存
func (s *InventoryService) GetInventoryByProductID(ctx context.Context, productID int64) (*inventory.Inventory, error) {
	return s.repo.FindInventoryByProductID(ctx, productID)
}

// ListInventories 列出所有库存
func (s *InventoryService) ListInventories(ctx context.Context, page, pageSize int) ([]inventory.Inventory, int64, error) {
	return s.repo.ListInventories(ctx, page, pageSize)
}

// UpdateInventory 更新库存
func (s *InventoryService) UpdateInventory(ctx context.Context, productID int64, req *UpdateInventoryRequest) (*inventory.Inventory, error) {
	// 获取库存
	inv, err := s.repo.FindInventoryByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// 更新库存
	inv.Quantity = req.Quantity

	// 保存库存
	if err := s.repo.UpdateInventory(ctx, inv); err != nil {
		return nil, err
	}

	return inv, nil
}

// ReserveInventory 预占库存
func (s *InventoryService) ReserveInventory(ctx context.Context, productID int64, quantity int) error {
	if err := s.repo.ReserveInventory(ctx, productID, quantity); err != nil {
		return err
	}

	// 发布库存预占事件
	s.bus.Publish(inventory.InventoryReservedEvent{
		ProductID: fmt.Sprintf("%d", productID),
		Quantity:  quantity,
	})

	return nil
}

// ReleaseInventory 释放库存
func (s *InventoryService) ReleaseInventory(ctx context.Context, productID int64, quantity int) error {
	if err := s.repo.ReleaseInventory(ctx, productID, quantity); err != nil {
		return err
	}

	// 发布库存释放事件
	s.bus.Publish(inventory.InventoryReleasedEvent{
		ProductID: fmt.Sprintf("%d", productID),
		Quantity:  quantity,
	})

	return nil
}

// CheckInventory 检查库存
func (s *InventoryService) CheckInventory(ctx context.Context, productID int64) (*inventory.Inventory, error) {
	return s.repo.FindInventoryByProductID(ctx, productID)
}
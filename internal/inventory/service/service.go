package service

import "eshop-monolith/internal/inventory/domain/models"

type InventoryService struct{}

func NewInventoryService(args ...interface{}) *InventoryService { return &InventoryService{} }

func (s *InventoryService) GetInventoryBySkuID(ctx interface{}, skuID int64) (*models.Inventory, error) {
	return &models.Inventory{Quantity: 999999}, nil
}

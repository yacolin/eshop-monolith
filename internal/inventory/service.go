package inventory

import (
	"context"
	"errors"

	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

type InventoryService struct {
	repo IinventoryRepository
	db   *gorm.DB
}

func NewInventoryService(repo IinventoryRepository, db *gorm.DB) *InventoryService {
	return &InventoryService{repo: repo, db: db}
}

type InventoryLogListResult struct {
	Total int64           `json:"total"`
	List  []*InventoryLog `json:"list"`
}

// Lock 下单预占库存（FOR UPDATE 行锁防超卖）
func (s *InventoryService) Lock(ctx context.Context, req *LockStockReq) (*Inventory, error) {
	var inv *Inventory
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		inv, err = s.repo.FindOrCreateWithTx(tx, req.SkuID, 0)
		if err != nil {
			return err
		}
		if inv.Available() < req.Quantity {
			return errcode.ErrInsufficientStock
		}
		beforeQty, beforeRes := inv.Quantity, inv.Reserved
		inv.Reserved += req.Quantity

		if err := s.repo.UpdateWithTx(tx, inv); err != nil {
			return err
		}
		return s.repo.CreateLogWithTx(tx, &InventoryLog{
			SkuID:          req.SkuID,
			ChangeType:     "order_lock",
			BeforeQuantity: beforeQty,
			AfterQuantity:  inv.Quantity,
			BeforeReserved: beforeRes,
			AfterReserved:  inv.Reserved,
			ChangeAmount:   req.Quantity,
			ReferenceID:    req.ReferenceID,
			Operator:       req.Operator,
		})
	})
	return inv, err
}

// Unlock 取消释放预占
func (s *InventoryService) Unlock(ctx context.Context, req *UnlockStockReq) (*Inventory, error) {
	var inv *Inventory
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		inv, err = s.repo.FindOrCreateWithTx(tx, req.SkuID, 0)
		if err != nil {
			return err
		}
		if inv.Reserved < req.Quantity {
			return errcode.ErrInvalidStockChange
		}
		beforeQty, beforeRes := inv.Quantity, inv.Reserved
		inv.Reserved -= req.Quantity

		if err := s.repo.UpdateWithTx(tx, inv); err != nil {
			return err
		}
		return s.repo.CreateLogWithTx(tx, &InventoryLog{
			SkuID:          req.SkuID,
			ChangeType:     "order_unlock",
			BeforeQuantity: beforeQty,
			AfterQuantity:  inv.Quantity,
			BeforeReserved: beforeRes,
			AfterReserved:  inv.Reserved,
			ChangeAmount:   -req.Quantity,
			ReferenceID:    req.ReferenceID,
			Operator:       req.Operator,
		})
	})
	return inv, err
}

// Deduct 支付扣减库存
func (s *InventoryService) Deduct(ctx context.Context, req *DeductStockReq) (*Inventory, error) {
	var inv *Inventory
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		inv, err = s.repo.FindOrCreateWithTx(tx, req.SkuID, 0)
		if err != nil {
			return err
		}
		if inv.Quantity < req.Quantity || inv.Reserved < req.Quantity {
			return errcode.ErrInsufficientStock
		}
		beforeQty, beforeRes := inv.Quantity, inv.Reserved
		inv.Quantity -= req.Quantity
		inv.Reserved -= req.Quantity

		if err := s.repo.UpdateWithTx(tx, inv); err != nil {
			return err
		}
		return s.repo.CreateLogWithTx(tx, &InventoryLog{
			SkuID:          req.SkuID,
			ChangeType:     "order_deduct",
			BeforeQuantity: beforeQty,
			AfterQuantity:  inv.Quantity,
			BeforeReserved: beforeRes,
			AfterReserved:  inv.Reserved,
			ChangeAmount:   -req.Quantity,
			ReferenceID:    req.ReferenceID,
			Operator:       req.Operator,
		})
	})
	return inv, err
}

// Restock 入库/补货
func (s *InventoryService) Restock(ctx context.Context, req *RestockReq) (*Inventory, error) {
	var inv *Inventory
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		inv, err = s.repo.FindOrCreateWithTx(tx, req.SkuID, req.WarehouseID)
		if err != nil {
			return err
		}
		beforeQty := inv.Quantity
		inv.Quantity += req.Quantity

		if inv.Quantity > inv.MaxThreshold {
			inv.Quantity = inv.MaxThreshold
		}

		if err := s.repo.UpdateWithTx(tx, inv); err != nil {
			return err
		}
		return s.repo.CreateLogWithTx(tx, &InventoryLog{
			SkuID:          req.SkuID,
			WarehouseID:    req.WarehouseID,
			ChangeType:     "inbound",
			BeforeQuantity: beforeQty,
			AfterQuantity:  inv.Quantity,
			ChangeAmount:   req.Quantity,
			ReferenceID:    req.ReferenceID,
			Operator:       req.Operator,
			Note:           req.Note,
		})
	})
	return inv, err
}

// GetStock 查询库存
func (s *InventoryService) GetStock(ctx context.Context, skuID int64) (*Inventory, error) {
	inv, err := s.repo.FindBySku(ctx, skuID, 0)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrInventoryNotFound
		}
		return nil, err
	}
	return inv, nil
}

// ListLogs 查询库存变更流水
func (s *InventoryService) ListLogs(ctx context.Context, req *InventoryLogQuery) (*InventoryLogListResult, error) {
	if req.Size <= 0 {
		req.Size = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	list, total, err := s.repo.ListLogs(ctx, req.SkuID, req.ChangeType, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*InventoryLog, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return &InventoryLogListResult{Total: total, List: items}, nil
}

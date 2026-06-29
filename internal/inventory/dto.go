package inventory

// LockStockReq 下单预占库存
type LockStockReq struct {
	SkuID       int64  `json:"sku_id" binding:"required"`
	Quantity    int64  `json:"quantity" binding:"required,gt=0"`
	ReferenceID string `json:"reference_id" binding:"max=64"`
	Operator    string `json:"operator" binding:"max=50"`
}

// UnlockStockReq 取消释放库存
type UnlockStockReq struct {
	SkuID       int64  `json:"sku_id" binding:"required"`
	Quantity    int64  `json:"quantity" binding:"required,gt=0"`
	ReferenceID string `json:"reference_id" binding:"max=64"`
	Operator    string `json:"operator" binding:"max=50"`
}

// DeductStockReq 支付扣减库存
type DeductStockReq struct {
	SkuID       int64  `json:"sku_id" binding:"required"`
	Quantity    int64  `json:"quantity" binding:"required,gt=0"`
	ReferenceID string `json:"reference_id" binding:"max=64"`
	Operator    string `json:"operator" binding:"max=50"`
}

// RestockReq 入库/补货
type RestockReq struct {
	SkuID       int64  `json:"sku_id" binding:"required"`
	WarehouseID int64  `json:"warehouse_id"`
	Quantity    int64  `json:"quantity" binding:"required,gt=0"`
	ReferenceID string `json:"reference_id" binding:"max=64"`
	Operator    string `json:"operator" binding:"max=50"`
	Note        string `json:"note" binding:"max=500"`
}

type InventoryLogQuery struct {
	SkuID      int64  `form:"sku_id" binding:"required"`
	ChangeType string `form:"change_type"`
	Page       int    `form:"page,default=1" binding:"gte=1"`
	Size       int    `form:"size,default=20" binding:"gte=1,lte=100"`
}

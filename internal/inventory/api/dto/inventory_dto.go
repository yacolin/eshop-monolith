package dto

import (
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/query"
)

// InventoryListQuery 库存列表查询参数
type InventoryListQuery struct {
	query.Pagination
	ProductName string `form:"product_name"`      // 产品名称模糊搜索
	SKU         string `form:"sku"`               // SKU精确搜索
	Status      string `form:"status"`            // 库存状态
	LowStock    *bool  `form:"low_stock"`         // 是否低库存
	SortBy      string `form:"sort_by"`           // 排序字段，例如 quantity, reserved, created_at
	Order       string `form:"order,default=asc"` // asc or desc
}

// InventoryListResult 库存列表结果（使用泛型）
type InventoryListResult = query.ListResult[models.Inventory]

// CreateInventoryDTO 创建库存请求
//
// inventories 表核心字段含义:
//   quantity  — 实际物理库存。下单时不扣, 支付时扣减。
//   reserved  — 已预占库存(下单未支付)。下单时+1, 支付/取消时-1。
//   threshold — 低库存预警阈值。quantity<=threshold 时自动标为 lowstock。
//
// 可用库存 = quantity - reserved (决定能否继续下单)
type CreateInventoryDTO struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=0"`  // 初始物理库存量
	Threshold int   `json:"threshold" binding:"min=0"`          // 低库存预警阈值(默认0)
}

// UpdateInventoryDTO 更新库存请求
//
// 注意: 一般情况下不应当手动修改 reserved(预占量),
//       它由下单/支付/取消流程自动维护。
type UpdateInventoryDTO struct {
	Quantity  *int `json:"quantity" binding:"omitempty,min=0"`  // 调整后物理库存
	Threshold *int `json:"threshold" binding:"omitempty,min=0"` // 调整后预警阈值
	Reserved  *int `json:"reserved" binding:"omitempty,min=0"`  // 调整后预占量(谨慎使用)
}

// ReserveInventoryDTO 下单预占库存请求
//
// 业务流程: 用户下单时调用, 将库存从"可用"变为"预占"。
//    UPDATE: reserved += quantity
//    效果:  可用库存减少, 但实际库存不变。
//    回滚:  取消订单时由 ReleaseInventoryDTO 释放。
type ReserveInventoryDTO struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"` // 预占数量
}

// ReleaseInventoryDTO 取消订单释放库存请求
//
// 业务流程: 用户取消订单时调用, 释放之前预占的库存。
//    UPDATE: reserved -= quantity
//    效果:  可用库存恢复, 实际库存不变。
//    触发:  订单取消/退款时调用。
type ReleaseInventoryDTO struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"` // 释放数量
}

// AdjustInventoryDTO 手动调整库存请求
//
// 用途: 盘点/入库/退货时手动调整物理库存。
//   正数: 增加库存(如入库、退货)
//   负数: 减少库存(如报损、盘亏)
//   注意: 订单支付流程不经过此接口, 走 DeductInventory。
type AdjustInventoryDTO struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required"` // 正数增加, 负数减少
}

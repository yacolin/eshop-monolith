package order

import (
	"context"
)

// Repository 订单仓储接口
type Repository interface {
	// Create 创建订单
	Create(ctx context.Context, order *Order) error
	// FindByID 根据ID查询订单
	FindByID(ctx context.Context, id string) (*Order, error)
	// FindByUserID 根据用户ID查询订单列表
	FindByUserID(ctx context.Context, userID string, page, pageSize int) ([]Order, int64, error)
	// Update 更新订单
	Update(ctx context.Context, order *Order) error
	// UpdateStatus 更新订单状态
	UpdateStatus(ctx context.Context, id, status string) error
	// Delete 删除订单
	Delete(ctx context.Context, id string) error
}
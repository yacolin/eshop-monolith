package repository

import (
	"context"

	"gorm.io/gorm"

	"eshop-monolith/internal/domain/order"
)

// OrderRepository 订单仓储实现
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return OrderRepository{db: db}
}

// Create 创建订单
func (r OrderRepository) Create(ctx context.Context, order *order.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// FindByID 根据ID查询订单
func (r OrderRepository) FindByID(ctx context.Context, id string) (*order.Order, error) {
	var foundOrder order.Order
	err := r.db.WithContext(ctx).Preload("Items").First(&foundOrder, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &foundOrder, nil
}

// FindByUserID 根据用户ID查询订单列表
func (r OrderRepository) FindByUserID(ctx context.Context, userID string, page, pageSize int) ([]order.Order, int64, error) {
	var orders []order.Order
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&order.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Update 更新订单
func (r OrderRepository) Update(ctx context.Context, order *order.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// UpdateStatus 更新订单状态
func (r OrderRepository) UpdateStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&order.Order{}).Where("id = ?", id).Update("status", status).Error
}

// Delete 删除订单
func (r OrderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&order.Order{}, "id = ?", id).Error
}
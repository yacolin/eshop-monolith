package repositories

import (
	"context"

	"eshop-monolith/internal/order/api/dto"
	"eshop-monolith/internal/order/domain/models"
	"eshop-monolith/internal/pkg/query"

	"gorm.io/gorm"
)

// IorderRepository 订单仓储接口
type IorderRepository interface {
	// Create 创建订单
	Create(ctx context.Context, order *models.Order) error
	// FindByID 根据ID查询订单
	FindByID(ctx context.Context, id int64) (*models.Order, error)
	// FindByUserID 根据用户ID查询订单列表
	FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]models.Order, int64, error)
	// Update 更新订单
	Update(ctx context.Context, order *models.Order) error
	// UpdateStatus 更新订单状态
	UpdateStatus(ctx context.Context, id int64, status string) error
	// Delete 删除订单
	Delete(ctx context.Context, id int64) error

	ListByQuery(ctx context.Context, q dto.OrderListQuery, offset, limit int) ([]models.Order, error)
	CountByQuery(ctx context.Context, q dto.OrderListQuery) (int64, error)
}

type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return OrderRepository{db: db}
}

// Create 创建订单
func (r OrderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// FindByID 根据ID查询订单
func (r OrderRepository) FindByID(ctx context.Context, id int64) (*models.Order, error) {
	var foundOrder models.Order
	err := r.db.WithContext(ctx).Preload("Items").First(&foundOrder, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &foundOrder, nil
}

// FindByUserID 根据用户ID查询订单列表
func (r OrderRepository) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&models.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
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
func (r OrderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// UpdateStatus 更新订单状态
func (r OrderRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

// Delete 删除订单
func (r OrderRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Order{}, "id = ?", id).Error
}

func (r OrderRepository) ListByQuery(ctx context.Context, q dto.OrderListQuery, offset, limit int) ([]models.Order, error) {
	var list []models.Order
	db := r.applyQueryConditions(ctx, q)
	db = r.applyOrder(db, q)

	if err := db.Preload("Items").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r OrderRepository) CountByQuery(ctx context.Context, q dto.OrderListQuery) (int64, error) {
	var count int64
	db := r.applyQueryConditions(ctx, q)

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// applyQueryConditions 应用查询条件（不包含排序）
func (r OrderRepository) applyQueryConditions(ctx context.Context, q dto.OrderListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.Order{})
	if q.CustomerID != nil {
		db = db.Where("customer_id = ?", q.CustomerID)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}

	if q.MinPrice != nil {
		db = db.Where("total_amount >= ?", q.MinPrice)
	}

	if q.MaxPrice != nil {
		db = db.Where("total_amount <= ?", q.MaxPrice)
	}

	return db
}

// applyOrder 应用排序
func (r OrderRepository) applyOrder(db *gorm.DB, q dto.OrderListQuery) *gorm.DB {
	return query.ApplyOrder(db, q.SortBy, q.Order, "id asc")
}

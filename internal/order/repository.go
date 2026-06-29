package order

import (
	"context"

	"gorm.io/gorm"
)

type IorderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	CreateItem(ctx context.Context, item *OrderItem) error
	CreateLog(ctx context.Context, log *OrderLog) error
	FindByOrderNo(ctx context.Context, orderNo string) (*Order, error)
	FindByID(ctx context.Context, id int64) (*Order, error)
	List(ctx context.Context, userID int64, status string, page, size int) ([]Order, int64, error)
	ListItems(ctx context.Context, orderID int64) ([]OrderItem, error)
	UpdateStatus(ctx context.Context, orderNo string, fromStatus, toStatus string) error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) IorderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *OrderRepository) CreateItem(ctx context.Context, item *OrderItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *OrderRepository) CreateLog(ctx context.Context, log *OrderLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *OrderRepository) FindByOrderNo(ctx context.Context, orderNo string) (*Order, error) {
	var o Order
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&o).Error
	return &o, err
}

func (r *OrderRepository) FindByID(ctx context.Context, id int64) (*Order, error) {
	var o Order
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&o).Error
	return &o, err
}

func (r *OrderRepository) List(ctx context.Context, userID int64, status string, page, size int) ([]Order, int64, error) {
	db := r.db.WithContext(ctx).Model(&Order{})
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []Order
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *OrderRepository) ListItems(ctx context.Context, orderID int64) ([]OrderItem, error) {
	var items []OrderItem
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Order("id ASC").Find(&items).Error
	return items, err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderNo string, fromStatus, toStatus string) error {
	return r.db.WithContext(ctx).Model(&Order{}).
		Where("order_no = ? AND status = ?", orderNo, fromStatus).
		Update("status", toStatus).Error
}

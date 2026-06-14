package repositories

import (
	"context"

	"eshop-monolith/internal/infra/repository/models"
	"eshop-monolith/internal/order/api/dto"
	orderModels "eshop-monolith/internal/order/domain/models"
	"eshop-monolith/pkg/query"

	"gorm.io/gorm"
)

// IorderRepository 订单仓储接口
type IorderRepository interface {
	// Create 创建订单
	Create(ctx context.Context, order *orderModels.Order) error
	// FindByID 根据ID查询订单
	FindByID(ctx context.Context, id int64) (*orderModels.Order, error)
	// FindByUserID 根据用户ID查询订单列表
	FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]orderModels.Order, int64, error)
	// Update 更新订单
	Update(ctx context.Context, order *orderModels.Order) error
	// UpdateStatus 更新订单状态
	UpdateStatus(ctx context.Context, id int64, status string) error
	// Delete 删除订单
	Delete(ctx context.Context, id int64) error
	// FindItemsByOrderID 根据订单ID查询订单项（分页）
	FindItemsByOrderID(ctx context.Context, orderID int64, page, pageSize int) ([]orderModels.OrderItem, int64, error)
	// ListAllItems 查询所有订单项（分页），支持条件筛选
	ListAllItems(ctx context.Context, q dto.OrderItemListQuery, offset, limit int) ([]orderModels.OrderItem, int64, error)
	// BatchGetOrderNo 批量查询订单号
	BatchGetOrderNo(ctx context.Context, orderIDs []int64) (map[int64]string, error)

	// CreateWithTx 在已有事务内创建订单
	CreateWithTx(tx *gorm.DB, order *orderModels.Order) error
	// UpdateStatusWithTx 在已有事务内更新订单状态
	UpdateStatusWithTx(tx *gorm.DB, id int64, status string) error
	// FindByIDWithTx 在事务内查询订单
	FindByIDWithTx(tx *gorm.DB, id int64) (*orderModels.Order, error)

	ListByQuery(ctx context.Context, q dto.OrderListQuery, offset, limit int) ([]orderModels.Order, error)
	CountByQuery(ctx context.Context, q dto.OrderListQuery) (int64, error)
}

type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(db *gorm.DB) IorderRepository {
	return &OrderRepository{db: db}
}

// Create 创建订单
func (r *OrderRepository) Create(ctx context.Context, order *orderModels.Order) error {
	po := models.OrderFromDomain(order)
	if err := r.db.WithContext(ctx).Omit("Items").Create(po).Error; err != nil {
		return err
	}
	order.ID = po.ID
	for i := range order.Items {
		itemPO := models.OrderItemFromDomain(&order.Items[i])
		itemPO.OrderID = po.ID
		if err := r.db.WithContext(ctx).Create(itemPO).Error; err != nil {
			return err
		}
		order.Items[i].ID = itemPO.ID
		order.Items[i].OrderID = itemPO.OrderID
	}
	return nil
}

// FindByID 根据ID查询订单
func (r *OrderRepository) FindByID(ctx context.Context, id int64) (*orderModels.Order, error) {
	var po models.OrderPO
	err := r.db.WithContext(ctx).Preload("Items").First(&po, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindByUserID 根据用户ID查询订单列表
func (r *OrderRepository) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]orderModels.Order, int64, error) {
	var pos []models.OrderPO
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.OrderPO{}).Where("customer_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Preload("Items").Where("customer_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	orders := make([]orderModels.Order, len(pos))
	for i, po := range pos {
		orders[i] = *po.ToDomain()
	}
	return orders, total, nil
}

// Update 更新订单
func (r *OrderRepository) Update(ctx context.Context, order *orderModels.Order) error {
	po := models.OrderFromDomain(order)
	return r.db.WithContext(ctx).Save(po).Error
}

// UpdateStatus 更新订单状态
func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&models.OrderPO{}).Where("id = ?", id).Update("status", status).Error
}

// CreateWithTx 在已有事务内创建订单（显式创建订单和订单项）
func (r *OrderRepository) CreateWithTx(tx *gorm.DB, order *orderModels.Order) error {
	po := models.OrderFromDomain(order)

	// 先创建订单主记录，获取自增 ID
	if err := tx.Omit("Items").Create(po).Error; err != nil {
		return err
	}
	order.ID = po.ID

	// 显式逐条创建订单项，确保 order_id 正确关联
	for i := range order.Items {
		itemPO := models.OrderItemFromDomain(&order.Items[i])
		itemPO.OrderID = po.ID
		if err := tx.Create(itemPO).Error; err != nil {
			return err
		}
		order.Items[i].ID = itemPO.ID
		order.Items[i].OrderID = itemPO.OrderID
	}
	return nil
}

// UpdateStatusWithTx 在已有事务内更新订单状态
func (r *OrderRepository) UpdateStatusWithTx(tx *gorm.DB, id int64, status string) error {
	return tx.Model(&models.OrderPO{}).Where("id = ?", id).Update("status", status).Error
}

// Delete 删除订单
func (r *OrderRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.OrderPO{}, "id = ?", id).Error
}

// FindItemsByOrderID 根据订单ID分页查询订单项
func (r *OrderRepository) FindItemsByOrderID(ctx context.Context, orderID int64, page, pageSize int) ([]orderModels.OrderItem, int64, error) {
	var pos []models.OrderItemPO
	var total int64

	db := r.db.WithContext(ctx).Model(&models.OrderItemPO{}).Where("order_id = ?", orderID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	items := make([]orderModels.OrderItem, len(pos))
	for i, po := range pos {
		items[i] = *po.ToDomain()
	}
	return items, total, nil
}

// ListAllItems 查询所有订单项（分页）
func (r *OrderRepository) ListAllItems(ctx context.Context, q dto.OrderItemListQuery, offset, limit int) ([]orderModels.OrderItem, int64, error) {
	var pos []models.OrderItemPO
	var total int64

	db := r.db.WithContext(ctx).Model(&models.OrderItemPO{})
	if q.OrderNo != "" {
		db = db.Joins("JOIN orders ON orders.id = order_items.order_id").Where("orders.order_no = ?", q.OrderNo)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = query.ApplyOrder(db, q.SortBy, q.Order, "id asc")
	if err := db.Offset(offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	items := make([]orderModels.OrderItem, len(pos))
	for i, po := range pos {
		items[i] = *po.ToDomain()
	}
	return items, total, nil
}

// BatchGetOrderNo 批量查询订单号
func (r *OrderRepository) BatchGetOrderNo(ctx context.Context, orderIDs []int64) (map[int64]string, error) {
	if len(orderIDs) == 0 {
		return map[int64]string{}, nil
	}
	var results []struct {
		ID      int64
		OrderNo string
	}
	if err := r.db.WithContext(ctx).Model(&models.OrderPO{}).Select("id, order_no").Where("id IN ?", orderIDs).Find(&results).Error; err != nil {
		return nil, err
	}
	orderNoMap := make(map[int64]string, len(results))
	for _, r := range results {
		orderNoMap[r.ID] = r.OrderNo
	}
	return orderNoMap, nil
}

func (r *OrderRepository) ListByQuery(ctx context.Context, q dto.OrderListQuery, offset, limit int) ([]orderModels.Order, error) {
	var pos []models.OrderPO
	db := r.applyQueryConditions(ctx, q)
	db = r.applyOrder(db, q)

	if err := db.Preload("Items").Offset(offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, err
	}

	orders := make([]orderModels.Order, len(pos))
	for i, po := range pos {
		orders[i] = *po.ToDomain()
	}
	return orders, nil
}

func (r *OrderRepository) CountByQuery(ctx context.Context, q dto.OrderListQuery) (int64, error) {
	var count int64
	db := r.applyQueryConditions(ctx, q)

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// applyQueryConditions 应用查询条件（不包含排序）
func (r *OrderRepository) applyQueryConditions(ctx context.Context, q dto.OrderListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.OrderPO{})
	if q.CustomerID != nil {
		db = db.Where("customer_id = ?", q.CustomerID)
	}
	if q.OrderNo != "" {
		db = db.Where("order_no = ?", q.OrderNo)
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

// FindByIDWithTx 在事务内查询订单（含订单项）
func (r *OrderRepository) FindByIDWithTx(tx *gorm.DB, id int64) (*orderModels.Order, error) {
	var po models.OrderPO
	if err := tx.First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}

	var itemPOs []models.OrderItemPO
	if err := tx.Where("order_id = ?", id).Find(&itemPOs).Error; err != nil {
		return nil, err
	}

	order := po.ToDomain()
	items := make([]orderModels.OrderItem, len(itemPOs))
	for i, ip := range itemPOs {
		items[i] = *ip.ToDomain()
	}
	order.Items = items
	return order, nil
}

// applyOrder 应用排序
func (r *OrderRepository) applyOrder(db *gorm.DB, q dto.OrderListQuery) *gorm.DB {
	return query.ApplyOrder(db, q.SortBy, q.Order, "id asc")
}

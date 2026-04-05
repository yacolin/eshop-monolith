package repositories

import (
	"context"
	"eshop-monolith/internal/payment/api/dto"
	"eshop-monolith/internal/payment/domain/models"

	"gorm.io/gorm"
)

// IRefundRepository 退款仓储接口
type IRefundRepository interface {
	// Create 创建退款记录
	Create(ctx context.Context, refund *models.Refund) error
	// GetByID 根据ID获取退款记录
	GetByID(ctx context.Context, id int64) (*models.Refund, error)
	// GetByPaymentID 根据支付ID获取退款记录
	GetByPaymentID(ctx context.Context, paymentID int64) ([]models.Refund, error)
	// GetByOrderID 根据订单ID获取退款记录
	GetByOrderID(ctx context.Context, orderID int64) ([]models.Refund, error)
	// Update 更新退款记录
	Update(ctx context.Context, refund *models.Refund) error
	// UpdateStatus 更新退款状态
	UpdateStatus(ctx context.Context, id int64, status string) error
	// ListByQuery 根据查询条件获取退款列表
	ListByQuery(ctx context.Context, q dto.RefundListQuery, offset, limit int) ([]models.Refund, error)
	// CountByQuery 根据查询条件统计退款数量
	CountByQuery(ctx context.Context, q dto.RefundListQuery) (int64, error)
	// Delete 删除退款记录
	Delete(ctx context.Context, id int64) error
}

// refundRepository 退款仓储实现
type refundRepository struct {
	db *gorm.DB
}

// NewRefundRepository 创建退款仓储实例
func NewRefundRepository(db *gorm.DB) IRefundRepository {
	return &refundRepository{db: db}
}

// Create 创建退款记录
func (r *refundRepository) Create(ctx context.Context, refund *models.Refund) error {
	return r.db.WithContext(ctx).Create(refund).Error
}

// GetByID 根据ID获取退款记录
func (r *refundRepository) GetByID(ctx context.Context, id int64) (*models.Refund, error) {
	var refund models.Refund
	err := r.db.WithContext(ctx).First(&refund, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &refund, nil
}

// GetByPaymentID 根据支付ID获取退款记录
func (r *refundRepository) GetByPaymentID(ctx context.Context, paymentID int64) ([]models.Refund, error) {
	var refunds []models.Refund
	err := r.db.WithContext(ctx).Where("payment_id = ?", paymentID).Order("created_at DESC").Find(&refunds).Error
	if err != nil {
		return nil, err
	}
	return refunds, nil
}

// GetByOrderID 根据订单ID获取退款记录
func (r *refundRepository) GetByOrderID(ctx context.Context, orderID int64) ([]models.Refund, error) {
	var refunds []models.Refund
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Order("created_at DESC").Find(&refunds).Error
	if err != nil {
		return nil, err
	}
	return refunds, nil
}

// Update 更新退款记录
func (r *refundRepository) Update(ctx context.Context, refund *models.Refund) error {
	return r.db.WithContext(ctx).Save(refund).Error
}

// UpdateStatus 更新退款状态
func (r *refundRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&models.Refund{}).Where("id = ?", id).Update("status", status).Error
}

// ListByQuery 根据查询条件获取退款列表
func (r *refundRepository) ListByQuery(ctx context.Context, q dto.RefundListQuery, offset, limit int) ([]models.Refund, error) {
	var refunds []models.Refund
	db := r.applyQueryConditions(ctx, q)
	err := db.Offset(offset).Limit(limit).Order(q.SortBy + " " + q.Order).Find(&refunds).Error
	if err != nil {
		return nil, err
	}
	return refunds, nil
}

// CountByQuery 根据查询条件统计退款数量
func (r *refundRepository) CountByQuery(ctx context.Context, q dto.RefundListQuery) (int64, error) {
	var count int64
	db := r.applyQueryConditions(ctx, q)
	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete 删除退款记录
func (r *refundRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Refund{}, "id = ?", id).Error
}

// applyQueryConditions 应用查询条件
func (r *refundRepository) applyQueryConditions(ctx context.Context, q dto.RefundListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.Refund{})
	if q.PaymentID != 0 {
		db = db.Where("payment_id = ?", q.PaymentID)
	}
	if q.OrderID != 0 {
		db = db.Where("order_id = ?", q.OrderID)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if q.StartDate != "" {
		db = db.Where("created_at >= ?", q.StartDate)
	}
	if q.EndDate != "" {
		db = db.Where("created_at <= ?", q.EndDate)
	}
	return db
}

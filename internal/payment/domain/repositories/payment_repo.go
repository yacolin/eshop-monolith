package repositories

import (
	"context"
	"eshop-monolith/internal/payment/api/dto"
	"eshop-monolith/internal/payment/domain/models"
	"eshop-monolith/pkg/query"

	"gorm.io/gorm"
)

// IPaymentRepository 支付仓储接口
type IPaymentRepository interface {
	// Create 创建支付记录
	Create(ctx context.Context, payment *models.Payment) error
	// GetByID 根据ID获取支付记录
	GetByID(ctx context.Context, id int64) (*models.Payment, error)
	// GetByOrderID 根据订单ID获取支付记录
	GetByOrderID(ctx context.Context, orderID int64) (*models.Payment, error)
	// GetByTransactionID 根据交易ID获取支付记录
	GetByTransactionID(ctx context.Context, transactionID string) (*models.Payment, error)
	// Update 更新支付记录
	Update(ctx context.Context, payment *models.Payment) error
	// UpdateStatus 更新支付状态
	UpdateStatus(ctx context.Context, id int64, status string) error
	// ListByQuery 根据查询条件获取支付列表
	ListByQuery(ctx context.Context, q dto.PaymentListQuery, offset, limit int) ([]models.Payment, error)
	// CountByQuery 根据查询条件统计支付数量
	CountByQuery(ctx context.Context, q dto.PaymentListQuery) (int64, error)
	// Delete 删除支付记录
	Delete(ctx context.Context, id int64) error

	// CreateTransaction 创建交易记录
	CreateTransaction(ctx context.Context, transaction *models.PaymentTransaction) error
	// GetTransactionsByPaymentID 根据支付ID获取交易记录
	GetTransactionsByPaymentID(ctx context.Context, paymentID int64) ([]models.PaymentTransaction, error)
}

// paymentRepository 支付仓储实现
type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository 创建支付仓储实例
func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &paymentRepository{db: db}
}

// Create 创建支付记录
func (r *paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// GetByID 根据ID获取支付记录
func (r *paymentRepository) GetByID(ctx context.Context, id int64) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Preload("Transactions").Preload("Refunds").First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByOrderID 根据订单ID获取支付记录
func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID int64) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Preload("Transactions").Preload("Refunds").Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByTransactionID 根据交易ID获取支付记录
func (r *paymentRepository) GetByTransactionID(ctx context.Context, transactionID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Preload("Transactions").Preload("Refunds").Where("transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// Update 更新支付记录
func (r *paymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

// UpdateStatus 更新支付状态
func (r *paymentRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&models.Payment{}).Where("id = ?", id).Update("status", status).Error
}

// ListByQuery 根据查询条件获取支付列表
func (r *paymentRepository) ListByQuery(ctx context.Context, q dto.PaymentListQuery, offset, limit int) ([]models.Payment, error) {
	var payments []models.Payment
	db := r.applyQueryConditions(ctx, q)
	db = r.applyOrder(db, q)
	err := db.Preload("Transactions").Offset(offset).Limit(limit).Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

// CountByQuery 根据查询条件统计支付数量
func (r *paymentRepository) CountByQuery(ctx context.Context, q dto.PaymentListQuery) (int64, error) {
	var count int64
	db := r.applyQueryConditions(ctx, q)
	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete 删除支付记录
func (r *paymentRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Payment{}, "id = ?", id).Error
}

// CreateTransaction 创建交易记录
func (r *paymentRepository) CreateTransaction(ctx context.Context, transaction *models.PaymentTransaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

// GetTransactionsByPaymentID 根据支付ID获取交易记录
func (r *paymentRepository) GetTransactionsByPaymentID(ctx context.Context, paymentID int64) ([]models.PaymentTransaction, error) {
	var transactions []models.PaymentTransaction
	err := r.db.WithContext(ctx).Where("payment_id = ?", paymentID).Order("created_at DESC").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// applyQueryConditions 应用查询条件
func (r *paymentRepository) applyQueryConditions(ctx context.Context, q dto.PaymentListQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.Payment{})
	if q.OrderID != 0 {
		db = db.Where("order_id = ?", q.OrderID)
	}
	if q.PaymentMethod != "" {
		db = db.Where("payment_method = ?", q.PaymentMethod)
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

func (r *paymentRepository) applyOrder(db *gorm.DB, q dto.PaymentListQuery) *gorm.DB {
	return query.ApplyOrder(db, q.SortBy, q.Order, "id asc")
}

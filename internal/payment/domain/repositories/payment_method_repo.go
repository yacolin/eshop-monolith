package repositories

import (
	"context"
	"eshop-monolith/internal/payment/domain/models"

	"gorm.io/gorm"
)

// IPaymentMethodRepository 支付方式仓储接口
type IPaymentMethodRepository interface {
	// Create 创建支付方式
	Create(ctx context.Context, paymentMethod *models.PaymentMethod) error
	// GetByID 根据ID获取支付方式
	GetByID(ctx context.Context, id int64) (*models.PaymentMethod, error)
	// GetByCode 根据编码获取支付方式
	GetByCode(ctx context.Context, code string) (*models.PaymentMethod, error)
	// Update 更新支付方式
	Update(ctx context.Context, paymentMethod *models.PaymentMethod) error
	// UpdateStatus 更新支付方式状态
	UpdateStatus(ctx context.Context, id int64, status int) error
	// List 获取支付方式列表
	List(ctx context.Context, limit, offset int) ([]models.PaymentMethod, int64, error)
	// ListByStatus 根据状态获取支付方式列表
	ListByStatus(ctx context.Context, status int, limit, offset int) ([]models.PaymentMethod, int64, error)
	// Delete 删除支付方式
	Delete(ctx context.Context, id int64) error
}

// paymentMethodRepository 支付方式仓储实现
type paymentMethodRepository struct {
	db *gorm.DB
}

// NewPaymentMethodRepository 创建支付方式仓储实例
func NewPaymentMethodRepository(db *gorm.DB) IPaymentMethodRepository {
	return &paymentMethodRepository{db: db}
}

// Create 创建支付方式
func (r *paymentMethodRepository) Create(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	return r.db.WithContext(ctx).Create(paymentMethod).Error
}

// GetByID 根据ID获取支付方式
func (r *paymentMethodRepository) GetByID(ctx context.Context, id int64) (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	err := r.db.WithContext(ctx).First(&paymentMethod, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

// GetByCode 根据编码获取支付方式
func (r *paymentMethodRepository) GetByCode(ctx context.Context, code string) (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&paymentMethod).Error
	if err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

// Update 更新支付方式
func (r *paymentMethodRepository) Update(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	return r.db.WithContext(ctx).Save(paymentMethod).Error
}

// UpdateStatus 更新支付方式状态
func (r *paymentMethodRepository) UpdateStatus(ctx context.Context, id int64, status int) error {
	return r.db.WithContext(ctx).Model(&models.PaymentMethod{}).Where("id = ?", id).Update("status", status).Error
}

// List 获取支付方式列表
func (r *paymentMethodRepository) List(ctx context.Context, limit, offset int) ([]models.PaymentMethod, int64, error) {
	var paymentMethods []models.PaymentMethod
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PaymentMethod{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&paymentMethods).Error
	if err != nil {
		return nil, 0, err
	}

	return paymentMethods, total, nil
}

// ListByStatus 根据状态获取支付方式列表
func (r *paymentMethodRepository) ListByStatus(ctx context.Context, status int, limit, offset int) ([]models.PaymentMethod, int64, error) {
	var paymentMethods []models.PaymentMethod
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PaymentMethod{}).Where("status = ?", status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&paymentMethods).Error
	if err != nil {
		return nil, 0, err
	}

	return paymentMethods, total, nil
}

// Delete 删除支付方式
func (r *paymentMethodRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.PaymentMethod{}, "id = ?", id).Error
}

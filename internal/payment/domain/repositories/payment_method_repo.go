package repositories

import (
	"context"
	payModels "eshop-monolith/internal/payment/domain/models"
	"eshop-monolith/internal/infra/repository/models"

	"gorm.io/gorm"
)

// IPaymentMethodRepository 支付方式仓储接口
type IPaymentMethodRepository interface {
	// Create 创建支付方式
	Create(ctx context.Context, paymentMethod *payModels.PaymentMethod) error
	// GetByID 根据ID获取支付方式
	GetByID(ctx context.Context, id int64) (*payModels.PaymentMethod, error)
	// GetByCode 根据编码获取支付方式
	GetByCode(ctx context.Context, code string) (*payModels.PaymentMethod, error)
	// Update 更新支付方式
	Update(ctx context.Context, paymentMethod *payModels.PaymentMethod) error
	// UpdateStatus 更新支付方式状态
	UpdateStatus(ctx context.Context, id int64, status int) error
	// List 获取支付方式列表
	List(ctx context.Context, limit, offset int) ([]payModels.PaymentMethod, int64, error)
	// ListByStatus 根据状态获取支付方式列表
	ListByStatus(ctx context.Context, status int, limit, offset int) ([]payModels.PaymentMethod, int64, error)
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
func (r *paymentMethodRepository) Create(ctx context.Context, paymentMethod *payModels.PaymentMethod) error {
	po := models.PaymentMethodFromDomain(paymentMethod)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	paymentMethod.ID = po.ID
	return nil
}

// GetByID 根据ID获取支付方式
func (r *paymentMethodRepository) GetByID(ctx context.Context, id int64) (*payModels.PaymentMethod, error) {
	var po models.PaymentMethodPO
	err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// GetByCode 根据编码获取支付方式
func (r *paymentMethodRepository) GetByCode(ctx context.Context, code string) (*payModels.PaymentMethod, error) {
	var po models.PaymentMethodPO
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// Update 更新支付方式
func (r *paymentMethodRepository) Update(ctx context.Context, paymentMethod *payModels.PaymentMethod) error {
	po := models.PaymentMethodFromDomain(paymentMethod)
	return r.db.WithContext(ctx).Save(po).Error
}

// UpdateStatus 更新支付方式状态
func (r *paymentMethodRepository) UpdateStatus(ctx context.Context, id int64, status int) error {
	return r.db.WithContext(ctx).Model(&models.PaymentMethodPO{}).Where("id = ?", id).Update("status", status).Error
}

// List 获取支付方式列表
func (r *paymentMethodRepository) List(ctx context.Context, limit, offset int) ([]payModels.PaymentMethod, int64, error) {
	var pos []models.PaymentMethodPO
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PaymentMethodPO{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	paymentMethods := make([]payModels.PaymentMethod, len(pos))
	for i, po := range pos {
		paymentMethods[i] = *po.ToDomain()
	}
	return paymentMethods, total, nil
}

// ListByStatus 根据状态获取支付方式列表
func (r *paymentMethodRepository) ListByStatus(ctx context.Context, status int, limit, offset int) ([]payModels.PaymentMethod, int64, error) {
	var pos []models.PaymentMethodPO
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PaymentMethodPO{}).Where("status = ?", status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sort ASC, created_at DESC").Limit(limit).Offset(offset).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	paymentMethods := make([]payModels.PaymentMethod, len(pos))
	for i, po := range pos {
		paymentMethods[i] = *po.ToDomain()
	}
	return paymentMethods, total, nil
}

// Delete 删除支付方式
func (r *paymentMethodRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.PaymentMethodPO{}, "id = ?", id).Error
}

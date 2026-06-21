package repositories

import (
	"context"

	"eshop-monolith/internal/coupon/domain/models"
	repoModels "eshop-monolith/internal/infra/repository/models"

	"gorm.io/gorm"
)

// IcouponRepository 优惠券模板仓储接口
type IcouponRepository interface {
	Create(ctx context.Context, coupon *models.Coupon) error
	Update(ctx context.Context, coupon *models.Coupon) error
	FindByID(ctx context.Context, id int64) (*models.Coupon, error)
	FindByStatus(ctx context.Context, status models.CouponStatus, offset, limit int) ([]models.Coupon, int64, error)
	List(ctx context.Context, offset, limit int) ([]models.Coupon, int64, error)
	DecrementRemainStock(ctx context.Context, id int64, quantity int) error

	// 事务方法
	CreateWithTx(tx *gorm.DB, coupon *models.Coupon) error
	FindByIDWithTx(tx *gorm.DB, id int64) (*models.Coupon, error)
	DecrementRemainStockWithTx(tx *gorm.DB, id int64, quantity int) error
}

// CouponRepository 优惠券模板仓储实现
type CouponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) IcouponRepository {
	return &CouponRepository{db: db}
}

func (r *CouponRepository) Create(ctx context.Context, coupon *models.Coupon) error {
	po := repoModels.CouponFromDomain(coupon)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	coupon.ID = po.ID
	return nil
}

func (r *CouponRepository) Update(ctx context.Context, coupon *models.Coupon) error {
	po := repoModels.CouponFromDomain(coupon)
	return r.db.WithContext(ctx).Save(po).Error
}

func (r *CouponRepository) FindByID(ctx context.Context, id int64) (*models.Coupon, error) {
	var po repoModels.CouponPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *CouponRepository) FindByStatus(ctx context.Context, status models.CouponStatus, offset, limit int) ([]models.Coupon, int64, error) {
	var pos []repoModels.CouponPO
	var total int64

	db := r.db.WithContext(ctx).Model(&repoModels.CouponPO{}).Where("status = ?", status)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]models.Coupon, len(pos))
	for i, po := range pos {
		result[i] = *po.ToDomain()
	}
	return result, total, nil
}

func (r *CouponRepository) List(ctx context.Context, offset, limit int) ([]models.Coupon, int64, error) {
	var pos []repoModels.CouponPO
	var total int64

	if err := r.db.WithContext(ctx).Model(&repoModels.CouponPO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]models.Coupon, len(pos))
	for i, po := range pos {
		result[i] = *po.ToDomain()
	}
	return result, total, nil
}

func (r *CouponRepository) DecrementRemainStock(ctx context.Context, id int64, quantity int) error {
	return r.db.WithContext(ctx).Model(&repoModels.CouponPO{}).
		Where("id = ? AND remain_stock >= ?", id, quantity).
		UpdateColumn("remain_stock", gorm.Expr("remain_stock - ?", quantity)).Error
}

func (r *CouponRepository) CreateWithTx(tx *gorm.DB, coupon *models.Coupon) error {
	po := repoModels.CouponFromDomain(coupon)
	if err := tx.Create(po).Error; err != nil {
		return err
	}
	coupon.ID = po.ID
	return nil
}

func (r *CouponRepository) FindByIDWithTx(tx *gorm.DB, id int64) (*models.Coupon, error) {
	var po repoModels.CouponPO
	err := tx.First(&po, id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *CouponRepository) DecrementRemainStockWithTx(tx *gorm.DB, id int64, quantity int) error {
	return tx.Model(&repoModels.CouponPO{}).
		Where("id = ? AND remain_stock >= ?", id, quantity).
		UpdateColumn("remain_stock", gorm.Expr("remain_stock - ?", quantity)).Error
}

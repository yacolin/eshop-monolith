package repositories

import (
	"context"

	"eshop-monolith/internal/coupon/domain/models"
	repoModels "eshop-monolith/internal/infra/repository/models"

	"gorm.io/gorm"
)

// IuserCouponRepository 用户优惠券仓储接口
type IuserCouponRepository interface {
	Create(ctx context.Context, uc *models.UserCoupon) error
	FindByID(ctx context.Context, id int64) (*models.UserCoupon, error)
	FindByUserID(ctx context.Context, userID int64, status models.UserCouponStatus, offset, limit int) ([]models.UserCoupon, int64, error)
	FindByUserAndCoupon(ctx context.Context, userID, couponID int64) ([]models.UserCoupon, error)
	FindByCode(ctx context.Context, code string) (*models.UserCoupon, error)
	UpdateStatus(ctx context.Context, id int64, status models.UserCouponStatus, orderNo string) error
	CountByUserAndStatus(ctx context.Context, userID int64, status models.UserCouponStatus) (int64, error)
	// 事务方法
	CreateWithTx(tx *gorm.DB, uc *models.UserCoupon) error
	UpdateStatusWithTx(tx *gorm.DB, id int64, status models.UserCouponStatus, orderNo string) error
	FindByIDWithTx(tx *gorm.DB, id int64) (*models.UserCoupon, error)
	FindUsableByUserID(ctx context.Context, userID int64) ([]models.UserCoupon, error)
}

// UserCouponRepository 用户优惠券仓储实现
type UserCouponRepository struct {
	db *gorm.DB
}

func NewUserCouponRepository(db *gorm.DB) IuserCouponRepository {
	return &UserCouponRepository{db: db}
}

func (r *UserCouponRepository) Create(ctx context.Context, uc *models.UserCoupon) error {
	po := repoModels.UserCouponFromDomain(uc)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	uc.ID = po.ID
	return nil
}

func (r *UserCouponRepository) FindByID(ctx context.Context, id int64) (*models.UserCoupon, error) {
	var po repoModels.UserCouponPO
	err := r.db.WithContext(ctx).First(&po, id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *UserCouponRepository) FindByUserID(ctx context.Context, userID int64, status models.UserCouponStatus, offset, limit int) ([]models.UserCoupon, int64, error) {
	var pos []repoModels.UserCouponPO
	var total int64

	db := r.db.WithContext(ctx).Model(&repoModels.UserCouponPO{}).Where("user_id = ?", userID)
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]models.UserCoupon, len(pos))
	for i, po := range pos {
		result[i] = *po.ToDomain()
	}
	return result, total, nil
}

func (r *UserCouponRepository) FindByUserAndCoupon(ctx context.Context, userID, couponID int64) ([]models.UserCoupon, error) {
	var pos []repoModels.UserCouponPO
	if err := r.db.WithContext(ctx).Where("user_id = ? AND coupon_id = ?", userID, couponID).Find(&pos).Error; err != nil {
		return nil, err
	}
	result := make([]models.UserCoupon, len(pos))
	for i, po := range pos {
		result[i] = *po.ToDomain()
	}
	return result, nil
}

func (r *UserCouponRepository) FindByCode(ctx context.Context, code string) (*models.UserCoupon, error) {
	var po repoModels.UserCouponPO
	err := r.db.WithContext(ctx).Where("coupon_code = ?", code).First(&po).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *UserCouponRepository) UpdateStatus(ctx context.Context, id int64, status models.UserCouponStatus, orderNo string) error {
	updates := map[string]any{
		"status": status,
	}
	if status == models.UserCouponStatusUsed {
		updates["used_at"] = gorm.Expr("NOW()")
	}
	if orderNo != "" {
		updates["order_no"] = orderNo
	}
	return r.db.WithContext(ctx).Model(&repoModels.UserCouponPO{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserCouponRepository) CountByUserAndStatus(ctx context.Context, userID int64, status models.UserCouponStatus) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&repoModels.UserCouponPO{}).Where("user_id = ?", userID)
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserCouponRepository) CreateWithTx(tx *gorm.DB, uc *models.UserCoupon) error {
	po := repoModels.UserCouponFromDomain(uc)
	if err := tx.Create(po).Error; err != nil {
		return err
	}
	uc.ID = po.ID
	return nil
}

func (r *UserCouponRepository) UpdateStatusWithTx(tx *gorm.DB, id int64, status models.UserCouponStatus, orderNo string) error {
	updates := map[string]any{
		"status": status,
	}
	if status == models.UserCouponStatusUsed {
		updates["used_at"] = gorm.Expr("NOW()")
	}
	if orderNo != "" {
		updates["order_no"] = orderNo
	}
	return tx.Model(&repoModels.UserCouponPO{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserCouponRepository) FindByIDWithTx(tx *gorm.DB, id int64) (*models.UserCoupon, error) {
	var po repoModels.UserCouponPO
	err := tx.First(&po, id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *UserCouponRepository) FindUsableByUserID(ctx context.Context, userID int64) ([]models.UserCoupon, error) {
	var pos []repoModels.UserCouponPO
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND expire_at > NOW()", userID, models.UserCouponStatusUnused).
		Order("created_at DESC").
		Find(&pos).Error; err != nil {
		return nil, err
	}
	result := make([]models.UserCoupon, len(pos))
	for i, po := range pos {
		result[i] = *po.ToDomain()
	}
	return result, nil
}

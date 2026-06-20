package repositories

import (
	"context"

	"eshop-monolith/internal/coupon/domain/models"

	"gorm.io/gorm"
)

// ICouponRepository 优惠券模板仓储接口
type ICouponRepository interface {
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

// IUserCouponRepository 用户优惠券仓储接口
type IUserCouponRepository interface {
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

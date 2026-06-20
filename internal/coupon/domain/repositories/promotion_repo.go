package repositories

import (
	"context"

	"eshop-monolith/internal/coupon/domain/models"

	"gorm.io/gorm"
)

// IPromotionRepository 促销活动仓储接口
type IPromotionRepository interface {
	Create(ctx context.Context, promotion *models.Promotion) error
	Update(ctx context.Context, promotion *models.Promotion) error
	FindByID(ctx context.Context, id int64) (*models.Promotion, error)
	FindActive(ctx context.Context) ([]models.Promotion, error)
	List(ctx context.Context, offset, limit int) ([]models.Promotion, int64, error)
	UpdateStatus(ctx context.Context, id int64, status models.PromotionStatus) error
	// 事务方法
	CreateWithTx(tx *gorm.DB, promotion *models.Promotion) error
	FindByIDWithTx(tx *gorm.DB, id int64) (*models.Promotion, error)
}

// IPromotionProductRepository 促销商品仓储接口
type IPromotionProductRepository interface {
	Create(ctx context.Context, pp *models.PromotionProduct) error
	CreateBatch(ctx context.Context, items []models.PromotionProduct) error
	FindByPromotionID(ctx context.Context, promotionID int64) ([]models.PromotionProduct, error)
	FindByProductID(ctx context.Context, productID int64) ([]models.PromotionProduct, error)
	DeleteByPromotionID(ctx context.Context, promotionID int64) error
}

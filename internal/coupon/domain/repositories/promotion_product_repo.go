package repositories

import (
	"context"

	"eshop-monolith/internal/coupon/domain/models"
	repoModels "eshop-monolith/internal/infra/repository/models"

	"gorm.io/gorm"
)

// IpromotionProductRepository 促销商品仓储接口
type IpromotionProductRepository interface {
	Create(ctx context.Context, pp *models.PromotionProduct) error
	CreateBatch(ctx context.Context, items []models.PromotionProduct) error
	FindByPromotionID(ctx context.Context, promotionID int64) ([]models.PromotionProduct, error)
	FindByProductID(ctx context.Context, productID int64) ([]models.PromotionProduct, error)
	DeleteByPromotionID(ctx context.Context, promotionID int64) error
}

// PromotionProductRepository 促销商品仓储实现
type PromotionProductRepository struct {
	db *gorm.DB
}

func NewPromotionProductRepository(db *gorm.DB) IpromotionProductRepository {
	return &PromotionProductRepository{db: db}
}

func (r *PromotionProductRepository) Create(ctx context.Context, pp *models.PromotionProduct) error {
	po := repoModels.PromotionProductFromDomain(pp)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	pp.ID = po.ID
	return nil
}

func (r *PromotionProductRepository) CreateBatch(ctx context.Context, items []models.PromotionProduct) error {
	pos := make([]repoModels.PromotionProductPO, len(items))
	for i, item := range items {
		pos[i] = *repoModels.PromotionProductFromDomain(&item)
	}
	return r.db.WithContext(ctx).CreateInBatches(pos, 100).Error
}

func (r *PromotionProductRepository) FindByPromotionID(ctx context.Context, promotionID int64) ([]models.PromotionProduct, error) {
	var pos []repoModels.PromotionProductPO
	if err := r.db.WithContext(ctx).Where("promotion_id = ?", promotionID).Find(&pos).Error; err != nil {
		return nil, err
	}
	result := make([]models.PromotionProduct, len(pos))
	for i, po := range pos {
		result[i] = *po.ToDomain()
	}
	return result, nil
}

func (r *PromotionProductRepository) FindByProductID(ctx context.Context, productID int64) ([]models.PromotionProduct, error) {
	var pos []repoModels.PromotionProductPO
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&pos).Error; err != nil {
		return nil, err
	}
	result := make([]models.PromotionProduct, len(pos))
	for i, po := range pos {
		result[i] = *po.ToDomain()
	}
	return result, nil
}

func (r *PromotionProductRepository) DeleteByPromotionID(ctx context.Context, promotionID int64) error {
	return r.db.WithContext(ctx).Where("promotion_id = ?", promotionID).Delete(&repoModels.PromotionProductPO{}).Error
}

package repositories

import (
	"context"
	"errors"

	"eshop-monolith/internal/coupon/domain/models"
	repoModels "eshop-monolith/internal/infra/repository/models"

	"gorm.io/gorm"
)

// PromotionRepository 促销活动仓储实现
type PromotionRepository struct {
	db *gorm.DB
}

func NewPromotionRepository(db *gorm.DB) IPromotionRepository {
	return &PromotionRepository{db: db}
}

func (r *PromotionRepository) Create(ctx context.Context, promotion *models.Promotion) error {
	po := repoModels.PromotionFromDomain(promotion)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	promotion.ID = po.ID
	return nil
}

func (r *PromotionRepository) Update(ctx context.Context, promotion *models.Promotion) error {
	po := repoModels.PromotionFromDomain(promotion)
	return r.db.WithContext(ctx).Save(po).Error
}

func (r *PromotionRepository) FindByID(ctx context.Context, id int64) (*models.Promotion, error) {
	var po repoModels.PromotionPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *PromotionRepository) FindActive(ctx context.Context) ([]models.Promotion, error) {
	var pos []repoModels.PromotionPO
	if err := r.db.WithContext(ctx).
		Where("status = ? AND start_time <= NOW() AND end_time >= NOW()", models.PromotionStatusActive).
		Order("sort_order ASC, created_at DESC").
		Find(&pos).Error; err != nil {
		return nil, err
	}
	result := make([]models.Promotion, len(pos))
	for i, po := range pos {
		result[i] = *po.ToDomain()
	}
	return result, nil
}

func (r *PromotionRepository) List(ctx context.Context, offset, limit int) ([]models.Promotion, int64, error) {
	var pos []repoModels.PromotionPO
	var total int64

	if err := r.db.WithContext(ctx).Model(&repoModels.PromotionPO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]models.Promotion, len(pos))
	for i, po := range pos {
		result[i] = *po.ToDomain()
	}
	return result, total, nil
}

func (r *PromotionRepository) UpdateStatus(ctx context.Context, id int64, status models.PromotionStatus) error {
	return r.db.WithContext(ctx).Model(&repoModels.PromotionPO{}).Where("id = ?", id).Update("status", status).Error
}

// 事务方法
func (r *PromotionRepository) CreateWithTx(tx *gorm.DB, promotion *models.Promotion) error {
	po := repoModels.PromotionFromDomain(promotion)
	if err := tx.Create(po).Error; err != nil {
		return err
	}
	promotion.ID = po.ID
	return nil
}

func (r *PromotionRepository) FindByIDWithTx(tx *gorm.DB, id int64) (*models.Promotion, error) {
	var po repoModels.PromotionPO
	if err := tx.First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return po.ToDomain(), nil
}

// PromotionProductRepository 促销商品仓储实现
type PromotionProductRepository struct {
	db *gorm.DB
}

func NewPromotionProductRepository(db *gorm.DB) IPromotionProductRepository {
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

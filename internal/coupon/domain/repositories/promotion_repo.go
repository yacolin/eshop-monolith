package repositories

import (
	"context"

	"eshop-monolith/internal/coupon/domain/models"
	repoModels "eshop-monolith/internal/infra/repository/models"

	"gorm.io/gorm"
)

// IpromotionRepository 促销活动仓储接口
type IpromotionRepository interface {
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

// PromotionRepository 促销活动仓储实现
type PromotionRepository struct {
	db *gorm.DB
}

func NewPromotionRepository(db *gorm.DB) IpromotionRepository {
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
	err := r.db.WithContext(ctx).First(&po, id).Error
	if err != nil {
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
	err := tx.First(&po, id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

package repositories

import (
	"context"
	"fmt"
	"time"

	"eshop-monolith/internal/infra/repository/models"
	reviewModels "eshop-monolith/internal/review/domain/models"
	"eshop-monolith/pkg/utils"

	"gorm.io/gorm"
)

// ReviewRepository 评论仓储
type ReviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository 创建评论仓储
func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// AutoMigrate 自动迁移表结构
func (r *ReviewRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&models.ReviewPO{}, &models.ProductRatingSummaryPO{})
}

// Create 创建评论
func (r *ReviewRepository) Create(ctx context.Context, rv *reviewModels.Review) error {
	po := models.ReviewFromDomain(rv)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	rv.ID = po.ID
	rv.CreatedAt = utils.Timestamp(po.CreatedAt)
	rv.UpdatedAt = utils.Timestamp(po.UpdatedAt)
	return nil
}

// GetByID 根据 ID 获取评论
func (r *ReviewRepository) GetByID(ctx context.Context, id int64) (*reviewModels.Review, error) {
	var po models.ReviewPO
	if err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// ExistsByOrderItem 判断某订单项是否已存在评论（防重复评论）
func (r *ReviewRepository) ExistsByOrderItem(ctx context.Context, orderItemID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ReviewPO{}).
		Where("order_item_id = ?", orderItemID).Count(&count).Error
	return count > 0, err
}

// ListByProduct 按产品分页查询可见评论（默认仅 approved）
func (r *ReviewRepository) ListByProduct(ctx context.Context, productID int64, statuses []string, page, size int) ([]reviewModels.Review, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.ReviewPO{}).Where("product_id = ?", productID)
	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var pos []models.ReviewPO
	if err := query.Order("created_at desc").
		Offset((page - 1) * size).Limit(size).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	reviews := make([]reviewModels.Review, len(pos))
	for i, po := range pos {
		reviews[i] = *po.ToDomain()
	}
	return reviews, total, nil
}

// ListByUser 按用户分页查询评论
func (r *ReviewRepository) ListByUser(ctx context.Context, userID int64, page, size int) ([]reviewModels.Review, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.ReviewPO{}).Where("user_id = ?", userID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var pos []models.ReviewPO
	if err := query.Order("created_at desc").
		Offset((page - 1) * size).Limit(size).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	reviews := make([]reviewModels.Review, len(pos))
	for i, po := range pos {
		reviews[i] = *po.ToDomain()
	}
	return reviews, total, nil
}

// ListPending 分页查询待审核评论（管理端）
func (r *ReviewRepository) ListByStatus(ctx context.Context, status string, page, size int) ([]reviewModels.Review, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.ReviewPO{}).Where("status = ?", status)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var pos []models.ReviewPO
	if err := query.Order("created_at asc").
		Offset((page - 1) * size).Limit(size).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	reviews := make([]reviewModels.Review, len(pos))
	for i, po := range pos {
		reviews[i] = *po.ToDomain()
	}
	return reviews, total, nil
}

// UpdateStatus 更新评论审核状态（管理端审核）
func (r *ReviewRepository) UpdateStatus(ctx context.Context, id int64, status reviewModels.ReviewStatus) error {
	return r.db.WithContext(ctx).Model(&models.ReviewPO{}).
		Where("id = ?", id).Update("status", string(status)).Error
}

// UpdateReply 更新商家回复
func (r *ReviewRepository) UpdateReply(ctx context.Context, id int64, reply string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.ReviewPO{}).
		Where("id = ?", id).Updates(map[string]interface{}{
		"reply":    reply,
		"reply_at": now,
	}).Error
}

// Delete 删除评论（评论者本人或管理员）
func (r *ReviewRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.ReviewPO{}, "id = ?", id).Error
}

// RatingStats 按产品统计各评分等级的数量（仅审核通过的评论）
type RatingStats struct {
	Rating1 int64
	Rating2 int64
	Rating3 int64
	Rating4 int64
	Rating5 int64
	Total   int64
}

// CountRatingByProduct 统计产品各评分数量（仅审核通过）
func (r *ReviewRepository) CountRatingByProduct(ctx context.Context, productID int64) (*RatingStats, error) {
	type row struct {
		Rating int
		Cnt    int64
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&models.ReviewPO{}).
		Select("rating, count(*) as cnt").
		Where("product_id = ? AND status = ?", productID, string(reviewModels.ReviewStatusApproved)).
		Group("rating").Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	stats := &RatingStats{}
	for _, rw := range rows {
		switch rw.Rating {
		case 1:
			stats.Rating1 = rw.Cnt
		case 2:
			stats.Rating2 = rw.Cnt
		case 3:
			stats.Rating3 = rw.Cnt
		case 4:
			stats.Rating4 = rw.Cnt
		case 5:
			stats.Rating5 = rw.Cnt
		}
		stats.Total += rw.Cnt
	}
	return stats, nil
}

// UpsertRatingSummary 写入或更新产品评分汇总
func (r *ReviewRepository) UpsertRatingSummary(ctx context.Context, productID int64, stats *RatingStats) (*reviewModels.ProductRatingSummary, error) {
	avg := 0.0
	if stats.Total > 0 {
		sum := float64(stats.Rating1)*1 + float64(stats.Rating2)*2 + float64(stats.Rating3)*3 + float64(stats.Rating4)*4 + float64(stats.Rating5)*5
		avg = sum / float64(stats.Total)
	}

	po := &models.ProductRatingSummaryPO{
		ProductID:     productID,
		AverageRating: avg,
		ReviewCount:   stats.Total,
		Rating1Count:  stats.Rating1,
		Rating2Count:  stats.Rating2,
		Rating3Count:  stats.Rating3,
		Rating4Count:  stats.Rating4,
		Rating5Count:  stats.Rating5,
	}

	// MySQL 语义：存在则更新，不存在则插入
	if err := r.db.WithContext(ctx).Save(po).Error; err != nil {
		return nil, fmt.Errorf("upsert rating summary failed: %w", err)
	}
	return po.ToDomain(), nil
}

// GetRatingSummary 获取产品评分汇总
func (r *ReviewRepository) GetRatingSummary(ctx context.Context, productID int64) (*reviewModels.ProductRatingSummary, error) {
	var po models.ProductRatingSummaryPO
	err := r.db.WithContext(ctx).First(&po, "product_id = ?", productID).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

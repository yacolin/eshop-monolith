package review

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type IreviewRepository interface {
	Create(ctx context.Context, review *Review) error
	FindByID(ctx context.Context, id int64) (*Review, error)
	ExistsByOrderItem(ctx context.Context, orderItemID int64) (bool, error)
	ListByProduct(ctx context.Context, productID int64, statuses []string, page, size int) ([]Review, int64, error)
	ListByUser(ctx context.Context, userID int64, page, size int) ([]Review, int64, error)
	ListByStatus(ctx context.Context, status string, page, size int) ([]Review, int64, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	UpdateReply(ctx context.Context, id int64, reply string) error
	Delete(ctx context.Context, id int64) error
	CountRatingByProduct(ctx context.Context, productID int64) (*RatingStats, error)
	UpsertRatingSummary(ctx context.Context, productID int64, stats *RatingStats) (*ProductRatingSummary, error)
	GetRatingSummary(ctx context.Context, productID int64) (*ProductRatingSummary, error)
}

type RatingStats struct {
	Rating1 int64
	Rating2 int64
	Rating3 int64
	Rating4 int64
	Rating5 int64
	Total   int64
}

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) IreviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(ctx context.Context, review *Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *ReviewRepository) FindByID(ctx context.Context, id int64) (*Review, error) {
	var rv Review
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rv).Error
	return &rv, err
}

func (r *ReviewRepository) ExistsByOrderItem(ctx context.Context, orderItemID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Review{}).
		Where("order_item_id = ?", orderItemID).Count(&count).Error
	return count > 0, err
}

func (r *ReviewRepository) ListByProduct(ctx context.Context, productID int64, statuses []string, page, size int) ([]Review, int64, error) {
	q := r.db.WithContext(ctx).Model(&Review{}).Where("product_id = ?", productID)
	if len(statuses) > 0 {
		q = q.Where("status IN ?", statuses)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []Review
	err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *ReviewRepository) ListByUser(ctx context.Context, userID int64, page, size int) ([]Review, int64, error) {
	q := r.db.WithContext(ctx).Model(&Review{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []Review
	err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *ReviewRepository) ListByStatus(ctx context.Context, status string, page, size int) ([]Review, int64, error) {
	q := r.db.WithContext(ctx).Model(&Review{}).Where("status = ?", status)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []Review
	err := q.Order("created_at ASC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *ReviewRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&Review{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ReviewRepository) UpdateReply(ctx context.Context, id int64, reply string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&Review{}).Where("id = ?", id).Updates(map[string]interface{}{
		"reply":    reply,
		"reply_at": now,
	}).Error
}

func (r *ReviewRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Review{}).Error
}

func (r *ReviewRepository) CountRatingByProduct(ctx context.Context, productID int64) (*RatingStats, error) {
	type row struct {
		Rating int
		Cnt    int64
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&Review{}).
		Select("rating, count(*) as cnt").
		Where("product_id = ? AND status = ?", productID, ReviewStatusApproved).
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

func (r *ReviewRepository) UpsertRatingSummary(ctx context.Context, productID int64, stats *RatingStats) (*ProductRatingSummary, error) {
	avg := 0.0
	if stats.Total > 0 {
		sum := float64(stats.Rating1)*1 + float64(stats.Rating2)*2 + float64(stats.Rating3)*3 + float64(stats.Rating4)*4 + float64(stats.Rating5)*5
		avg = sum / float64(stats.Total)
	}
	s := &ProductRatingSummary{
		ProductID:     productID,
		AverageRating: avg,
		ReviewCount:   stats.Total,
		Rating1Count:  stats.Rating1,
		Rating2Count:  stats.Rating2,
		Rating3Count:  stats.Rating3,
		Rating4Count:  stats.Rating4,
		Rating5Count:  stats.Rating5,
	}
	if err := r.db.WithContext(ctx).Save(s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (r *ReviewRepository) GetRatingSummary(ctx context.Context, productID int64) (*ProductRatingSummary, error) {
	var s ProductRatingSummary
	err := r.db.WithContext(ctx).First(&s, "product_id = ?", productID).Error
	return &s, err
}

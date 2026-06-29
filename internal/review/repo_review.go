package review

import (
	"context"

	"gorm.io/gorm"
)

type IreviewRepository interface {
	CreateReview(ctx context.Context, review *Review) error
	FindByID(ctx context.Context, id int64) (*Review, error)
	ExistsByOrderItem(ctx context.Context, orderItemID int64) (bool, error)
	ListBySpu(ctx context.Context, spuID int64, statuses []int8, page, size int) ([]Review, int64, error)
	ListByUser(ctx context.Context, userID int64, page, size int) ([]Review, int64, error)
	ListByStatus(ctx context.Context, status int8, page, size int) ([]Review, int64, error)
	UpdateStatus(ctx context.Context, id int64, status int8, rejectReason string) error
	Delete(ctx context.Context, id int64) error

	// media
	CreateMedia(ctx context.Context, media *ReviewMedia) error
	ListMediaByReviewID(ctx context.Context, reviewID int64) ([]ReviewMedia, error)

	// replies
	CreateReply(ctx context.Context, reply *ReviewReply) (*ReviewReply, error)
	UpdateLatestReply(ctx context.Context, reviewID int64, replyID int64, replyCount int) error
	ListRepliesByReviewID(ctx context.Context, reviewID int64) ([]ReviewReply, error)

	// rating
	UpsertRatingSummary(ctx context.Context, r *ReviewRating) error
	GetRatingSummary(ctx context.Context, spuID int64) (*ReviewRating, error)
	CountRatingBySpu(ctx context.Context, spuID int64) (map[int8]int64, int64, error)

	// audit log
	CreateAuditLog(ctx context.Context, log *ReviewAuditLog) error
}

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) IreviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, review *Review) error {
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
		Where("order_item_id = ? AND status != ?", orderItemID, ReviewStatusDeleted).Count(&count).Error
	return count > 0, err
}

func (r *ReviewRepository) ListBySpu(ctx context.Context, spuID int64, statuses []int8, page, size int) ([]Review, int64, error) {
	q := r.db.WithContext(ctx).Model(&Review{}).Where("spu_id = ?", spuID)
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

func (r *ReviewRepository) ListByStatus(ctx context.Context, status int8, page, size int) ([]Review, int64, error) {
	q := r.db.WithContext(ctx).Model(&Review{}).Where("status = ?", status)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []Review
	err := q.Order("created_at ASC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *ReviewRepository) UpdateStatus(ctx context.Context, id int64, status int8, rejectReason string) error {
	updates := map[string]interface{}{"status": status}
	if rejectReason != "" {
		updates["reject_reason"] = rejectReason
	}
	return r.db.WithContext(ctx).Model(&Review{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ReviewRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Review{}).Error
}

func (r *ReviewRepository) CreateMedia(ctx context.Context, media *ReviewMedia) error {
	return r.db.WithContext(ctx).Create(media).Error
}

func (r *ReviewRepository) ListMediaByReviewID(ctx context.Context, reviewID int64) ([]ReviewMedia, error) {
	var list []ReviewMedia
	err := r.db.WithContext(ctx).Where("review_id = ?", reviewID).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ReviewRepository) CreateReply(ctx context.Context, reply *ReviewReply) (*ReviewReply, error) {
	if err := r.db.WithContext(ctx).Create(reply).Error; err != nil {
		return nil, err
	}
	return reply, nil
}

func (r *ReviewRepository) UpdateLatestReply(ctx context.Context, reviewID int64, replyID int64, replyCount int) error {
	return r.db.WithContext(ctx).Model(&Review{}).Where("id = ?", reviewID).
		Updates(map[string]interface{}{
			"latest_reply_id": replyID,
			"reply_count":     replyCount,
		}).Error
}

func (r *ReviewRepository) ListRepliesByReviewID(ctx context.Context, reviewID int64) ([]ReviewReply, error) {
	var list []ReviewReply
	err := r.db.WithContext(ctx).Where("review_id = ? AND status = ?", reviewID, ReplyStatusNormal).
		Order("created_at ASC").Find(&list).Error
	return list, err
}

func (r *ReviewRepository) UpsertRatingSummary(ctx context.Context, rating *ReviewRating) error {
	return r.db.WithContext(ctx).Save(rating).Error
}

func (r *ReviewRepository) GetRatingSummary(ctx context.Context, spuID int64) (*ReviewRating, error) {
	var rating ReviewRating
	err := r.db.WithContext(ctx).First(&rating, "spu_id = ?", spuID).Error
	return &rating, err
}

func (r *ReviewRepository) CountRatingBySpu(ctx context.Context, spuID int64) (map[int8]int64, int64, error) {
	type row struct {
		Rating int8
		Cnt    int64
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&Review{}).
		Select("overall_rating as rating, count(*) as cnt").
		Where("spu_id = ? AND status = ?", spuID, ReviewStatusApproved).
		Group("overall_rating").Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	counts := map[int8]int64{}
	var total int64
	for _, rw := range rows {
		counts[rw.Rating] = rw.Cnt
		total += rw.Cnt
	}
	return counts, total, nil
}

func (r *ReviewRepository) CreateAuditLog(ctx context.Context, log *ReviewAuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// ComputeRatingSummary calculates rating summary from raw counts.
func ComputeRatingSummary(spuID int64, counts map[int8]int64, total int64) *ReviewRating {
	avg := 0.0
	if total > 0 {
		var sum float64
		for rating, cnt := range counts {
			sum += float64(rating) * float64(cnt)
		}
		avg = sum / float64(total)
	}
	return &ReviewRating{
		SpuID:            spuID,
		AvgOverallRating: avg,
		TotalReviews:     int(total),
		Rating5Count:     int(counts[5]),
		Rating4Count:     int(counts[4]),
		Rating3Count:     int(counts[3]),
		Rating2Count:     int(counts[2]),
		Rating1Count:     int(counts[1]),
	}
}

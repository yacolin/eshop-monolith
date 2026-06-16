package models

import (
	"time"

	"github.com/bytedance/sonic"
	"gorm.io/gorm"

	domain "eshop-monolith/internal/review/domain/models"
	"eshop-monolith/pkg/utils"
)

// ReviewPO 评论持久化对象
type ReviewPO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	ProductID   int64          `gorm:"type:bigint;not null;index:idx_review_product"`
	UserID      int64          `gorm:"type:bigint;not null;index:idx_review_user"`
	OrderItemID int64          `gorm:"type:bigint;not null;uniqueIndex:uk_order_item"` // 一个订单项只能评论一次
	OrderNo     string         `gorm:"type:varchar(64);not null;default:''"`
	Rating      int            `gorm:"type:tinyint;not null"`
	Content     string         `gorm:"type:text"`
	Media       string         `gorm:"type:text"` // JSON 序列化的 []ReviewMedia
	Status      string         `gorm:"type:varchar(20);not null;default:pending;index:idx_review_status"`
	Reply       string         `gorm:"type:text"`
	ReplyAt     *time.Time     `gorm:"type:timestamp;null"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (ReviewPO) TableName() string { return "reviews" }

// ToDomain 转换为领域模型
func (po *ReviewPO) ToDomain() *domain.Review {
	r := &domain.Review{
		ID:          po.ID,
		ProductID:   po.ProductID,
		UserID:      po.UserID,
		OrderItemID: po.OrderItemID,
		OrderNo:     po.OrderNo,
		Rating:      po.Rating,
		Content:     po.Content,
		Media:       decodeReviewMedia(po.Media),
		Status:      domain.ReviewStatus(po.Status),
		Reply:       po.Reply,
		CreatedAt:   utils.Timestamp(po.CreatedAt),
		UpdatedAt:   utils.Timestamp(po.UpdatedAt),
	}
	if po.ReplyAt != nil {
		t := utils.Timestamp(*po.ReplyAt)
		r.ReplyAt = &t
	}
	return r
}

// ReviewFromDomain 从领域模型创建 PO
func ReviewFromDomain(r *domain.Review) *ReviewPO {
	po := &ReviewPO{
		ID:          r.ID,
		ProductID:   r.ProductID,
		UserID:      r.UserID,
		OrderItemID: r.OrderItemID,
		OrderNo:     r.OrderNo,
		Rating:      r.Rating,
		Content:     r.Content,
		Media:       encodeReviewMedia(r.Media),
		Status:      string(r.Status),
		Reply:       r.Reply,
		CreatedAt:   time.Time(r.CreatedAt),
		UpdatedAt:   time.Time(r.UpdatedAt),
	}
	if r.ReplyAt != nil {
		t := time.Time(*r.ReplyAt)
		po.ReplyAt = &t
	}
	return po
}

// encodeReviewMedia 将媒体列表序列化为 JSON 字符串
func encodeReviewMedia(media []domain.ReviewMedia) string {
	if len(media) == 0 {
		return ""
	}
	b, err := sonic.Marshal(media)
	if err != nil {
		return ""
	}
	return string(b)
}

// decodeReviewMedia 将 JSON 字符串反序列化为媒体列表
func decodeReviewMedia(s string) []domain.ReviewMedia {
	if s == "" {
		return nil
	}
	var media []domain.ReviewMedia
	if err := sonic.Unmarshal([]byte(s), &media); err != nil {
		return nil
	}
	return media
}

// ProductRatingSummaryPO 产品评分汇总持久化对象
type ProductRatingSummaryPO struct {
	ProductID     int64      `gorm:"primaryKey"`
	AverageRating float64    `gorm:"type:decimal(3,2);not null;default:0"`
	ReviewCount   int64      `gorm:"type:bigint;not null;default:0"`
	Rating1Count  int64      `gorm:"type:bigint;not null;default:0"`
	Rating2Count  int64      `gorm:"type:bigint;not null;default:0"`
	Rating3Count  int64      `gorm:"type:bigint;not null;default:0"`
	Rating4Count  int64      `gorm:"type:bigint;not null;default:0"`
	Rating5Count  int64      `gorm:"type:bigint;not null;default:0"`
	UpdatedAt     time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
}

func (ProductRatingSummaryPO) TableName() string { return "product_rating_summaries" }

// ToDomain 转换为领域模型
func (po *ProductRatingSummaryPO) ToDomain() *domain.ProductRatingSummary {
	return &domain.ProductRatingSummary{
		ProductID:     po.ProductID,
		AverageRating: po.AverageRating,
		ReviewCount:   po.ReviewCount,
		Rating1Count:  po.Rating1Count,
		Rating2Count:  po.Rating2Count,
		Rating3Count:  po.Rating3Count,
		Rating4Count:  po.Rating4Count,
		Rating5Count:  po.Rating5Count,
		UpdatedAt:     utils.Timestamp(po.UpdatedAt),
	}
}

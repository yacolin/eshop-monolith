package review

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
	ReviewStatusHidden   ReviewStatus = "hidden"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

type ReviewMedia struct {
	Type      MediaType `json:"type"`
	URL       string    `json:"url"`
	Thumbnail string    `json:"thumbnail"`
}

type Review struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   int64          `gorm:"not null;index:idx_review_product" json:"product_id"`
	UserID      int64          `gorm:"not null;index:idx_review_user" json:"user_id"`
	OrderItemID int64          `gorm:"not null;uniqueIndex:uk_order_item" json:"order_item_id"`
	OrderNo     string         `gorm:"type:varchar(64);not null;default:''" json:"order_no"`
	Rating      int            `gorm:"type:tinyint;not null" json:"rating"`
	Content     string         `gorm:"type:text" json:"content"`
	MediaJSON   string         `gorm:"type:text" json:"-"`
	Status      string         `gorm:"type:varchar(20);not null;default:pending;index:idx_review_status" json:"status"`
	Reply       string         `gorm:"type:text" json:"reply"`
	ReplyAt     *time.Time     `gorm:"type:datetime(3)" json:"reply_at,omitempty"`
	CreatedAt   time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`

	// 非持久化，仅用于响应组装
	UserName  string `gorm:"-" json:"user_name,omitempty"`
	UserAvatar string `gorm:"-" json:"user_avatar,omitempty"`
}

func (Review) TableName() string { return "reviews" }

func (r *Review) GetMedia() []ReviewMedia {
	if r.MediaJSON == "" {
		return nil
	}
	var media []ReviewMedia
	if err := json.Unmarshal([]byte(r.MediaJSON), &media); err != nil {
		return nil
	}
	return media
}

type ProductRatingSummary struct {
	ProductID     int64     `gorm:"primaryKey"`
	AverageRating float64   `gorm:"type:decimal(3,2);not null;default:0"`
	ReviewCount   int64     `gorm:"not null;default:0"`
	Rating1Count  int64     `gorm:"not null;default:0"`
	Rating2Count  int64     `gorm:"not null;default:0"`
	Rating3Count  int64     `gorm:"not null;default:0"`
	Rating4Count  int64     `gorm:"not null;default:0"`
	Rating5Count  int64     `gorm:"not null;default:0"`
	UpdatedAt     time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)"`
}

func (ProductRatingSummary) TableName() string { return "product_rating_summaries" }

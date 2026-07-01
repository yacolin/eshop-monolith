package review

import (
	"gorm.io/gorm"

	"eshop-monolith/pkg/utils"
)

type Review struct {
	ID              int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64           `gorm:"not null;index:idx_user" json:"user_id"`
	OrderID         int64           `gorm:"not null;index:idx_order" json:"order_id"`
	OrderItemID     *int64          `gorm:"index" json:"order_item_id,omitempty"`
	SpuID           int64           `gorm:"not null;index:idx_spu" json:"spu_id"`
	SkuID           *int64          `gorm:"index:idx_sku" json:"sku_id,omitempty"`
	OverallRating   int8            `gorm:"type:tinyint;not null" json:"overall_rating"`
	QualityRating   *int8           `gorm:"type:tinyint" json:"quality_rating,omitempty"`
	LogisticsRating *int8           `gorm:"type:tinyint" json:"logistics_rating,omitempty"`
	ServiceRating   *int8           `gorm:"type:tinyint" json:"service_rating,omitempty"`
	Content         string          `gorm:"type:text" json:"content"`
	IsAnonymous     bool            `gorm:"not null;default:false" json:"is_anonymous"`
	Status          int8            `gorm:"type:tinyint;not null;default:0;index:idx_status_created" json:"status"`
	RejectReason    string          `gorm:"type:varchar(200)" json:"reject_reason,omitempty"`
	LatestReplyID   *int64          `gorm:"" json:"latest_reply_id,omitempty"`
	ReplyCount      int             `gorm:"not null;default:0" json:"reply_count"`
	LikeCount       int             `gorm:"not null;default:0" json:"like_count"`
	HelpfulCount    int             `gorm:"not null;default:0" json:"helpful_count"`
	CreatedAt       utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt       utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"type:datetime(3);index" json:"-"`
}

func (Review) TableName() string { return "rev_reviews" }

type ReviewMedia struct {
	ID          int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	ReviewID    int64           `gorm:"not null;index:idx_review" json:"review_id"`
	MediaType   int8            `gorm:"type:tinyint;not null" json:"media_type"`
	MediaURL    string          `gorm:"type:varchar(500);not null" json:"media_url"`
	SortOrder   int             `gorm:"not null;default:0" json:"sort_order"`
	AuditStatus int8            `gorm:"type:tinyint;not null;default:0;index:idx_audit" json:"audit_status"`
	CreatedAt   utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"type:datetime(3);index" json:"-"`
}

func (ReviewMedia) TableName() string { return "rev_review_media" }

type ReviewReply struct {
	ID          int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	ReviewID    int64           `gorm:"not null;index:idx_review" json:"review_id"`
	ReplyBy     int64           `gorm:"not null;index:idx_reply_by" json:"reply_by"`
	ReplyByType int8            `gorm:"type:tinyint;not null;default:1" json:"reply_by_type"`
	Content     string          `gorm:"type:text;not null" json:"content"`
	Status      int8            `gorm:"type:tinyint;not null;default:1" json:"status"`
	CreatedAt   utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt   utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"type:datetime(3);index" json:"-"`
}

func (ReviewReply) TableName() string { return "rev_review_replies" }

type ReviewRating struct {
	ID                 int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	SpuID              int64           `gorm:"not null;uniqueIndex:uk_spu" json:"spu_id"`
	AvgOverallRating   float64         `gorm:"type:decimal(2,1);not null;default:0" json:"avg_overall_rating"`
	AvgQualityRating   float64         `gorm:"type:decimal(2,1);not null;default:0" json:"avg_quality_rating"`
	AvgLogisticsRating float64         `gorm:"type:decimal(2,1);not null;default:0" json:"avg_logistics_rating"`
	AvgServiceRating   float64         `gorm:"type:decimal(2,1);not null;default:0" json:"avg_service_rating"`
	Rating5Count       int             `gorm:"not null;default:0" json:"rating_5_count"`
	Rating4Count       int             `gorm:"not null;default:0" json:"rating_4_count"`
	Rating3Count       int             `gorm:"not null;default:0" json:"rating_3_count"`
	Rating2Count       int             `gorm:"not null;default:0" json:"rating_2_count"`
	Rating1Count       int             `gorm:"not null;default:0" json:"rating_1_count"`
	TotalReviews       int             `gorm:"not null;default:0" json:"total_reviews"`
	WithMediaCount     int             `gorm:"not null;default:0" json:"with_media_count"`
	UpdatedAt          utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
}

func (ReviewRating) TableName() string { return "rev_review_ratings" }

type ReviewAuditLog struct {
	ID             int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	ReviewID       int64           `gorm:"not null;index:idx_review" json:"review_id"`
	AuditorID      int64           `gorm:"not null;index:idx_auditor" json:"auditor_id"`
	Action         int8            `gorm:"type:tinyint;not null" json:"action"`
	Reason         string          `gorm:"type:varchar(200)" json:"reason,omitempty"`
	SensitiveWords string          `gorm:"type:json" json:"sensitive_words,omitempty"`
	CreatedAt      utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);index:idx_created" json:"created_at"`
}

func (ReviewAuditLog) TableName() string { return "rev_review_audit_logs" }

// 状态常量
const (
	ReviewStatusPending  int8 = 0
	ReviewStatusApproved int8 = 1
	ReviewStatusRejected int8 = 2
	ReviewStatusDeleted  int8 = 3
)

const (
	MediaTypeImage int8 = 1
	MediaTypeVideo int8 = 2
)

const (
	ReplyByMerchant int8 = 1
	ReplyByPlatform int8 = 2

	ReplyStatusNormal  int8 = 1
	ReplyStatusDeleted int8 = 2
)

const (
	AuditActionApprove int8 = 1
	AuditActionReject  int8 = 2
	AuditActionReview  int8 = 3
)

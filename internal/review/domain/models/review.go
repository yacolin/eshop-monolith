package models

import "eshop-monolith/pkg/utils"

// ReviewStatus 评论审核状态
type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"  // 待审核
	ReviewStatusApproved ReviewStatus = "approved" // 审核通过（对外可见）
	ReviewStatusRejected ReviewStatus = "rejected" // 审核拒绝
	ReviewStatusHidden   ReviewStatus = "hidden"   // 已隐藏（管理员下架）
)

// MediaType 评论媒体类型
type MediaType string

const (
	MediaTypeImage MediaType = "image" // 图片
	MediaTypeVideo MediaType = "video" // 视频
)

// ReviewMedia 评论媒体资源（图片/视频）
type ReviewMedia struct {
	Type     MediaType `json:"type"`      // 媒体类型
	URL      string    `json:"url"`       // 资源地址
	Thumbnail string   `json:"thumbnail"` // 缩略图地址（视频封面/图片缩略）
}

// Review 评论/评分领域模型
type Review struct {
	ID         int64            `json:"id"`
	ProductID  int64            `json:"product_id"`  // 被评论的产品 ID
	UserID     int64            `json:"user_id"`     // 评论者用户 ID
	OrderItemID int64           `json:"order_item_id"` // 关联订单项（用于校验已购买 + 防重复）
	OrderNo    string           `json:"order_no"`    // 冗余订单号，便于展示
	Rating     int              `json:"rating"`      // 评分 1-5 星
	Content    string           `json:"content"`     // 文字评论
	Media      []ReviewMedia    `json:"media"`       // 图片/视频评论
	Status     ReviewStatus     `json:"status"`      // 审核状态
	Reply      string           `json:"reply"`       // 商家回复
	ReplyAt    *utils.Timestamp `json:"reply_at,omitempty"`
	CreatedAt  utils.Timestamp  `json:"created_at"`
	UpdatedAt  utils.Timestamp  `json:"updated_at"`

	// 非持久化字段，仅用于响应组装（展示评论者昵称、头像）
	UserName  string `json:"user_name,omitempty"`
	UserAvatar string `json:"user_avatar,omitempty"`
}

// TableName 表名
func (Review) TableName() string {
	return "reviews"
}

// ProductRatingSummary 产品评分汇总（用于推荐与搜索排序）
type ProductRatingSummary struct {
	ProductID     int64           `json:"product_id"`
	AverageRating float64         `json:"average_rating"` // 平均评分
	ReviewCount   int64           `json:"review_count"`   // 评价总数
	Rating1Count  int64           `json:"rating_1_count"` // 1 星数量
	Rating2Count  int64           `json:"rating_2_count"`
	Rating3Count  int64           `json:"rating_3_count"`
	Rating4Count  int64           `json:"rating_4_count"`
	Rating5Count  int64           `json:"rating_5_count"`
	UpdatedAt     utils.Timestamp `json:"updated_at"`
}

// TableName 表名
func (ProductRatingSummary) TableName() string {
	return "product_rating_summaries"
}

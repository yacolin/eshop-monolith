package events

// ReviewCreatedEvent 评论创建事件
// 用于通知搜索/推荐系统刷新产品评分与排序
type ReviewCreatedEvent struct {
	ReviewID  int64  `json:"review_id"`
	ProductID int64  `json:"product_id"`
	UserID    int64  `json:"user_id"`
	Rating    int    `json:"rating"`
	Status    string `json:"status"` // 审核状态
}

func (e ReviewCreatedEvent) RoutingKey() string { return "review.created" }

// ReviewModeratedEvent 评论审核结果事件（通过/拒绝/隐藏）
type ReviewModeratedEvent struct {
	ReviewID  int64  `json:"review_id"`
	ProductID int64  `json:"product_id"`
	Status    string `json:"status"` // approved / rejected / hidden
}

func (e ReviewModeratedEvent) RoutingKey() string { return "review.moderated" }

// ReviewDeletedEvent 评论删除事件
type ReviewDeletedEvent struct {
	ReviewID  int64 `json:"review_id"`
	ProductID int64 `json:"product_id"`
}

func (e ReviewDeletedEvent) RoutingKey() string { return "review.deleted" }

// RatingSummaryUpdatedEvent 产品评分汇总更新事件
// 影响产品推荐与搜索排序，可由推荐系统订阅
type RatingSummaryUpdatedEvent struct {
	ProductID     int64   `json:"product_id"`
	AverageRating float64 `json:"average_rating"`
	ReviewCount   int64   `json:"review_count"`
}

func (e RatingSummaryUpdatedEvent) RoutingKey() string { return "rating-summary.updated" }

package review

import "eshop-monolith/pkg/query"

type CreateReviewReq struct {
	OrderID         int64  `json:"order_id" binding:"required,gt=0"`
	OrderItemID     *int64 `json:"order_item_id"`
	SpuID           int64  `json:"spu_id" binding:"required,gt=0"`
	SkuID           *int64 `json:"sku_id"`
	OverallRating   int8   `json:"overall_rating" binding:"required,min=1,max=5"`
	QualityRating   *int8  `json:"quality_rating" binding:"omitempty,min=1,max=5"`
	LogisticsRating *int8  `json:"logistics_rating" binding:"omitempty,min=1,max=5"`
	ServiceRating   *int8  `json:"service_rating" binding:"omitempty,min=1,max=5"`
	Content         string `json:"content" binding:"max=2000"`
	IsAnonymous     bool   `json:"is_anonymous"`
}

type ReviewResp struct {
	ID              int64  `json:"id"`
	UserID          int64  `json:"user_id"`
	OrderID         int64  `json:"order_id"`
	OrderItemID     *int64 `json:"order_item_id,omitempty"`
	SpuID           int64  `json:"spu_id"`
	SkuID           *int64 `json:"sku_id,omitempty"`
	OverallRating   int8   `json:"overall_rating"`
	QualityRating   *int8  `json:"quality_rating,omitempty"`
	LogisticsRating *int8  `json:"logistics_rating,omitempty"`
	ServiceRating   *int8  `json:"service_rating,omitempty"`
	Content         string `json:"content"`
	IsAnonymous     bool   `json:"is_anonymous"`
	Status          int8   `json:"status"`
	RejectReason    string `json:"reject_reason,omitempty"`
	ReplyCount      int    `json:"reply_count"`
	LikeCount       int    `json:"like_count"`
	HelpfulCount    int    `json:"helpful_count"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

type ReviewListReq struct {
	query.Pagination
}

type ReviewListResult struct {
	Total int64         `json:"total"`
	List  []*ReviewResp `json:"list"`
}

type CreateMediaReq struct {
	ReviewID  int64  `json:"review_id" binding:"required"`
	MediaType int8   `json:"media_type" binding:"required,oneof=1 2"`
	MediaURL  string `json:"media_url" binding:"required,max=500"`
}

type ModerateReviewReq struct {
	Status int8   `json:"status" binding:"required,oneof=1 2"`
	Reason string `json:"reason" binding:"max=200"`
}

type ReplyReviewReq struct {
	Content string `json:"content" binding:"required,max=1000"`
}

type ReviewRatingResp struct {
	SpuID              int64   `json:"spu_id"`
	AvgOverallRating   float64 `json:"avg_overall_rating"`
	AvgQualityRating   float64 `json:"avg_quality_rating"`
	AvgLogisticsRating float64 `json:"avg_logistics_rating"`
	AvgServiceRating   float64 `json:"avg_service_rating"`
	Rating5Count       int     `json:"rating_5_count"`
	Rating4Count       int     `json:"rating_4_count"`
	Rating3Count       int     `json:"rating_3_count"`
	Rating2Count       int     `json:"rating_2_count"`
	Rating1Count       int     `json:"rating_1_count"`
	TotalReviews       int     `json:"total_reviews"`
	WithMediaCount     int     `json:"with_media_count"`
}

type CreateReviewInput struct {
	UserID          int64
	OrderID         int64
	OrderItemID     *int64
	SpuID           int64
	SkuID           *int64
	OverallRating   int8
	QualityRating   *int8
	LogisticsRating *int8
	ServiceRating   *int8
	Content         string
	IsAnonymous     bool
}

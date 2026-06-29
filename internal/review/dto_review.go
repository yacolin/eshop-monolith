package review

import "eshop-monolith/pkg/query"

type CreateReviewReq struct {
	ProductID   int64        `json:"product_id" binding:"required,gt=0"`
	OrderItemID int64        `json:"order_item_id" binding:"required,gt=0"`
	Rating      int          `json:"rating" binding:"required,min=1,max=5"`
	Content     string       `json:"content" binding:"max=2000"`
	Media       []ReviewMedia `json:"media" binding:"omitempty,dive"`
}

type ReviewListReq struct {
	query.Pagination
}

type ReviewResp struct {
	ID          int64        `json:"id"`
	ProductID   int64        `json:"product_id"`
	UserID      int64        `json:"user_id"`
	UserName    string       `json:"user_name,omitempty"`
	UserAvatar  string       `json:"user_avatar,omitempty"`
	OrderItemID int64        `json:"order_item_id"`
	OrderNo     string       `json:"order_no,omitempty"`
	Rating      int          `json:"rating"`
	Content     string       `json:"content"`
	Media       []ReviewMedia `json:"media,omitempty"`
	Status      string       `json:"status"`
	Reply       string       `json:"reply,omitempty"`
	ReplyAt     int64        `json:"reply_at,omitempty"`
	CreatedAt   int64        `json:"created_at"`
	UpdatedAt   int64        `json:"updated_at"`
}

type ReviewListResult struct {
	Total int64        `json:"total"`
	List  []*ReviewResp `json:"list"`
}

type ModerateReviewReq struct {
	Status string `json:"status" binding:"required,oneof=approved rejected hidden"`
	Reason string `json:"reason"`
}

type ReplyReviewReq struct {
	Reply string `json:"reply" binding:"required,max=1000"`
}

type ProductRatingResp struct {
	ProductID     int64   `json:"product_id"`
	AverageRating float64 `json:"average_rating"`
	ReviewCount   int64   `json:"review_count"`
	Rating1Count  int64   `json:"rating_1_count"`
	Rating2Count  int64   `json:"rating_2_count"`
	Rating3Count  int64   `json:"rating_3_count"`
	Rating4Count  int64   `json:"rating_4_count"`
	Rating5Count  int64   `json:"rating_5_count"`
}

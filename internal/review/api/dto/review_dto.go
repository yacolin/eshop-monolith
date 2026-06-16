package dto

import (
	"eshop-monolith/internal/review/domain/models"
	"eshop-monolith/pkg/query"
	"eshop-monolith/pkg/utils"
)

// CreateReviewReq 创建评论请求
type CreateReviewReq struct {
	ProductID   int64                  `json:"product_id" binding:"required,gt=0"`
	OrderItemID int64                  `json:"order_item_id" binding:"required,gt=0"`
	Rating      int                    `json:"rating" binding:"required,min=1,max=5"`
	Content     string                 `json:"content" binding:"max=2000"`
	Media       []models.ReviewMedia   `json:"media" binding:"omitempty,dive"`
}

// ReviewMediaDTO 评论媒体
type ReviewMediaDTO struct {
	Type      models.MediaType `json:"type"`
	URL       string           `json:"url"`
	Thumbnail string           `json:"thumbnail,omitempty"`
}

// ReviewResp 评论响应
type ReviewResp struct {
	ID          int64                `json:"id"`
	ProductID   int64                `json:"product_id"`
	UserID      int64                `json:"user_id"`
	UserName    string               `json:"user_name,omitempty"`
	UserAvatar  string               `json:"user_avatar,omitempty"`
	OrderItemID int64                `json:"order_item_id"`
	OrderNo     string               `json:"order_no,omitempty"`
	Rating      int                  `json:"rating"`
	Content     string               `json:"content"`
	Media       []ReviewMediaDTO     `json:"media,omitempty"`
	Status      models.ReviewStatus  `json:"status"`
	Reply       string               `json:"reply,omitempty"`
	ReplyAt     *utils.Timestamp     `json:"reply_at,omitempty"`
	CreatedAt   utils.Timestamp      `json:"created_at"`
	UpdatedAt   utils.Timestamp      `json:"updated_at"`
}

// ToReviewResp 领域模型转响应
func ToReviewResp(r *models.Review) *ReviewResp {
	resp := &ReviewResp{
		ID:          r.ID,
		ProductID:   r.ProductID,
		UserID:      r.UserID,
		UserName:    r.UserName,
		UserAvatar:  r.UserAvatar,
		OrderItemID: r.OrderItemID,
		OrderNo:     r.OrderNo,
		Rating:      r.Rating,
		Content:     r.Content,
		Status:      r.Status,
		Reply:       r.Reply,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
	if r.ReplyAt != nil {
		t := *r.ReplyAt
		resp.ReplyAt = &t
	}
	if len(r.Media) > 0 {
		resp.Media = make([]ReviewMediaDTO, len(r.Media))
		for i, m := range r.Media {
			resp.Media[i] = ReviewMediaDTO{
				Type:      m.Type,
				URL:       m.URL,
				Thumbnail: m.Thumbnail,
			}
		}
	}
	return resp
}

// ToReviewRespList 领域模型列表转响应列表
func ToReviewRespList(list []*models.Review) []*ReviewResp {
	resp := make([]*ReviewResp, len(list))
	for i, r := range list {
		resp[i] = ToReviewResp(r)
	}
	return resp
}

// ReviewListResp 评论列表响应（带分页）
type ReviewListResp struct {
	Total int64          `json:"total"`
	List  []*ReviewResp  `json:"list"`
}

// ListReviewReq 评论列表请求
type ListReviewReq struct {
	query.Pagination
}

// ModerateReviewReq 审核评论请求（管理端）
type ModerateReviewReq struct {
	Status models.ReviewStatus `json:"status" binding:"required,oneof=approved rejected hidden"`
	Reason string              `json:"reason"` // 拒绝/隐藏原因（可选，记录用）
}

// ReplyReviewReq 商家回复请求
type ReplyReviewReq struct {
	Reply string `json:"reply" binding:"required,max=1000"`
}

// ProductRatingResp 产品评分汇总响应
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

// ToProductRatingResp 领域模型转响应
func ToProductRatingResp(s *models.ProductRatingSummary) *ProductRatingResp {
	return &ProductRatingResp{
		ProductID:     s.ProductID,
		AverageRating: s.AverageRating,
		ReviewCount:   s.ReviewCount,
		Rating1Count:  s.Rating1Count,
		Rating2Count:  s.Rating2Count,
		Rating3Count:  s.Rating3Count,
		Rating4Count:  s.Rating4Count,
		Rating5Count:  s.Rating5Count,
	}
}

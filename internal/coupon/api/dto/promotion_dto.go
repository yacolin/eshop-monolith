package dto

import (
	"eshop-monolith/pkg/utils"
)

// PromotionResponse 促销活动响应
type PromotionResponse struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	PromoType   string          `json:"promo_type"`
	Scope       string          `json:"scope"`
	ScopeValue  string          `json:"scope_value,omitempty"`
	Rule        string          `json:"rule"`
	StartTime   utils.Timestamp `json:"start_time"`
	EndTime     utils.Timestamp `json:"end_time"`
	Status      string          `json:"status"`
	SortOrder   int             `json:"sort_order"`
	CreatedAt   utils.Timestamp `json:"created_at"`
	UpdatedAt   utils.Timestamp `json:"updated_at"`
}

// PromotionListResult 促销活动列表
type PromotionListResult struct {
	Total int64               `json:"total"`
	List  []PromotionResponse `json:"list"`
}

// CreatePromotionDTO 创建促销活动请求
type CreatePromotionDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	PromoType   string `json:"promo_type" binding:"required,oneof=time_discount full_reduce"`
	Scope       string `json:"scope" binding:"required,oneof=all category product"`
	ScopeValue  string `json:"scope_value"`
	Rule        string `json:"rule" binding:"required"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
	SortOrder   int    `json:"sort_order"`
}

// UpdatePromotionDTO 更新促销活动请求
type UpdatePromotionDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Scope       string `json:"scope" binding:"oneof=all category product"`
	ScopeValue  string `json:"scope_value"`
	Rule        string `json:"rule"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Status      string `json:"status" binding:"oneof=pending active cancelled"`
	SortOrder   int    `json:"sort_order"`
}

// LinkProductsDTO 关联促销商品请求
type LinkProductsDTO struct {
	ProductIDs []int64 `json:"product_ids" binding:"required,min=1"`
	Discount   int64   `json:"discount" binding:"required"`
}

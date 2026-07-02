package base

import "eshop-monolith/pkg/query"

type NotificationListReq struct {
	query.Pagination
}

type NotificationResp struct {
	ID              int64  `json:"id"`
	UserID          int64  `json:"user_id"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	ContentTemplate string `json:"content_template,omitempty"`
	TemplateParams  string `json:"template_params,omitempty"`
	Channel         int8   `json:"channel"`
	Category        int8   `json:"category"`
	TargetType      string `json:"target_type,omitempty"`
	TargetID        *int64 `json:"target_id,omitempty"`
	RedirectURL     string `json:"redirect_url,omitempty"`
	IconURL         string `json:"icon_url,omitempty"`
	IsRead          bool   `json:"is_read"`
	Priority        int8   `json:"priority"`
	CreatedBy       int64  `json:"created_by"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

type NotificationListResult struct {
	Total int64               `json:"total"`
	List  []*NotificationResp `json:"list"`
}

type UnreadCountResp struct {
	Count int64 `json:"count"`
}

type SendSystemNotificationReq struct {
	UserID  int64  `json:"user_id" binding:"gte=0"`
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

package notification

import "eshop-monolith/pkg/query"

type NotificationListReq struct {
	query.Pagination
}

type NotificationResp struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	IsRead    bool   `json:"is_read"`
	ReadAt    int64  `json:"read_at,omitempty"`
	CreatedAt int64  `json:"created_at"`
}

type NotificationListResult struct {
	Total int64              `json:"total"`
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

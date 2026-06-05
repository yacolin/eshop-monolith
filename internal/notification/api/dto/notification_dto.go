package dto

import (
	"eshop-monolith/internal/notification/domain/models"
	"eshop-monolith/pkg/query"
	"eshop-monolith/pkg/utils"
)

// ListNotificationReq 通知列表请求
type ListNotificationReq struct {
	query.Pagination
}

// NotificationResp 通知响应
type NotificationResp struct {
	ID        int64                `json:"id"`
	UserID    int64                `json:"user_id"`
	Title     string               `json:"title"`
	Content   string               `json:"content"`
	Type      models.NotificationType `json:"type"`
	IsRead    bool                 `json:"is_read"`
	ReadAt    *utils.Timestamp     `json:"read_at,omitempty"`
	CreatedAt utils.Timestamp      `json:"created_at"`
}

// ToNotificationResp 领域模型转响应
func ToNotificationResp(n *models.Notification) *NotificationResp {
	resp := &NotificationResp{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Content:   n.Content,
		Type:      n.Type,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}
	if n.ReadAt != nil {
		t := *n.ReadAt
		resp.ReadAt = &t
	}
	return resp
}

// ToNotificationRespList 领域模型列表转响应列表
func ToNotificationRespList(list []*models.Notification) []*NotificationResp {
	resp := make([]*NotificationResp, len(list))
	for i, n := range list {
		resp[i] = ToNotificationResp(n)
	}
	return resp
}

// NotificationListResp 通知列表响应（带分页）
type NotificationListResp struct {
	Total int64               `json:"total"`
	List  []*NotificationResp `json:"list"`
}

// UnreadCountResp 未读计数响应
type UnreadCountResp struct {
	Count int64 `json:"count"`
}

// SendSystemNotificationReq 发送系统通知请求
type SendSystemNotificationReq struct {
	UserID int64  `json:"user_id" binding:"gte=0"` // 0 表示全体
	Title  string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

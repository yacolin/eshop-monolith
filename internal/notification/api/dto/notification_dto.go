package dto

import (
	"eshop-monolith/internal/notification/domain/models"
	"eshop-monolith/pkg/query"
	"eshop-monolith/pkg/utils"
)

// NotificationListQuery 通知列表请求
type NotificationListQuery struct {
	query.Pagination
}

// NotificationResponse 通知响应
type NotificationResponse struct {
	ID        int64                  `json:"id"`
	UserID    int64                  `json:"user_id"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Type      models.NotificationType `json:"type"`
	IsRead    bool                   `json:"is_read"`
	ReadAt    *utils.Timestamp       `json:"read_at,omitempty"`
	CreatedAt utils.Timestamp        `json:"created_at"`
}

// ToNotificationResponse 领域模型转响应
func ToNotificationResponse(n *models.Notification) *NotificationResponse {
	resp := &NotificationResponse{
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

// ToNotificationResponseList 领域模型列表转响应列表
func ToNotificationResponseList(list []*models.Notification) []*NotificationResponse {
	resp := make([]*NotificationResponse, len(list))
	for i, n := range list {
		resp[i] = ToNotificationResponse(n)
	}
	return resp
}

// NotificationListResult 通知列表响应（带分页）
type NotificationListResult struct {
	Total int64                  `json:"total"`
	List  []*NotificationResponse `json:"list"`
}

// UnreadCountResponse 未读计数响应
type UnreadCountResponse struct {
	Count int64 `json:"count"`
}

// SendSystemNotificationDTO 发送系统通知请求
type SendSystemNotificationDTO struct {
	UserID  int64  `json:"user_id" binding:"gte=0"` // 0 表示全体
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

package models

import "eshop-monolith/pkg/utils"

// NotificationType 通知类型
type NotificationType string

const (
	NotificationTypeSystem  NotificationType = "system"  // 系统通知
	NotificationTypeOrder   NotificationType = "order"   // 订单通知
	NotificationTypePayment NotificationType = "payment" // 支付通知
	NotificationTypeFlash   NotificationType = "flash"   // 秒杀通知
	NotificationTypeAdmin   NotificationType = "admin"   // 管理通知
)

// Notification 站内信/通知
type Notification struct {
	ID         int64              `json:"id"`
	UserID     int64              `json:"user_id"`     // 接收用户 ID，0 表示全体
	Title      string             `json:"title"`
	Content    string             `json:"content"`
	Type       NotificationType   `json:"type"`
	IsRead     bool               `json:"is_read"`
	ReadAt     *utils.Timestamp   `json:"read_at,omitempty"`
	CreatedAt  utils.Timestamp    `json:"created_at"`
	UpdatedAt  utils.Timestamp    `json:"updated_at"`
}

// TableName 表名
func (Notification) TableName() string {
	return "notifications"
}

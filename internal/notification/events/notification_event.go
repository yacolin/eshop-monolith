package events

// NotificationCreatedEvent 通知创建事件
type NotificationCreatedEvent struct {
	NotificationID int64  `json:"notification_id"`
	UserID         int64  `json:"user_id"`
	Title          string `json:"title"`
	Type           string `json:"type"`
}

// SystemNotificationEvent 系统通知发送事件
type SystemNotificationEvent struct {
	NotificationID int64  `json:"notification_id"`
	TargetUserID   int64  `json:"target_user_id"` // 0 表示全体
	Title          string `json:"title"`
	Content        string `json:"content"`
}

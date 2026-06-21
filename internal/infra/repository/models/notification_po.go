package models

import (
	"time"

	"gorm.io/gorm"

	notifModels "eshop-monolith/internal/notification/domain/models"
	"eshop-monolith/pkg/utils"
)

// NotificationPO 通知持久化对象
type NotificationPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"type:bigint;not null;index"`
	Title     string         `gorm:"type:varchar(200);not null"`
	Content   string         `gorm:"type:text;not null"`
	Type      string         `gorm:"type:varchar(20);not null;index"`
	IsRead    bool           `gorm:"type:tinyint(1);not null;default:0;index"`
	ReadAt    *time.Time     `gorm:"type:timestamp;null"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (NotificationPO) TableName() string { return "notifications" }

// ToDomain 转换为领域模型
func (po *NotificationPO) ToDomain() *notifModels.Notification {
	n := &notifModels.Notification{
		ID:        po.ID,
		UserID:    po.UserID,
		Title:     po.Title,
		Content:   po.Content,
		Type:      notifModels.NotificationType(po.Type),
		IsRead:    po.IsRead,
		CreatedAt: utils.Timestamp(po.CreatedAt),
		UpdatedAt: utils.Timestamp(po.UpdatedAt),
	}
	if po.ReadAt != nil {
		t := utils.Timestamp(*po.ReadAt)
		n.ReadAt = &t
	}
	return n
}

// NotificationFromDomain 从领域模型创建 PO
func NotificationFromDomain(n *notifModels.Notification) *NotificationPO {
	po := &NotificationPO{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Content:   n.Content,
		Type:      string(n.Type),
		IsRead:    n.IsRead,
		CreatedAt: time.Time(n.CreatedAt),
		UpdatedAt: time.Time(n.UpdatedAt),
	}
	if n.ReadAt != nil {
		t := time.Time(*n.ReadAt)
		po.ReadAt = &t
	}
	return po
}

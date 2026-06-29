package notification

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type InotificationRepository interface {
	Create(ctx context.Context, n *Notification) error
	FindByID(ctx context.Context, id int64) (*Notification, error)
	ListByUserID(ctx context.Context, userID int64, page, size int) ([]Notification, int64, error)
	GetUnreadCount(ctx context.Context, userID int64) (int64, error)
	MarkAsRead(ctx context.Context, id, userID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
	Delete(ctx context.Context, id, userID int64) error
}

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) InotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *NotificationRepository) FindByID(ctx context.Context, id int64) (*Notification, error) {
	var n Notification
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&n).Error
	return &n, err
}

func (r *NotificationRepository) ListByUserID(ctx context.Context, userID int64, page, size int) ([]Notification, int64, error) {
	q := r.db.WithContext(ctx).Model(&Notification{}).
		Where("user_id = ? OR user_id = 0", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []Notification
	err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Notification{}).
		Where("(user_id = ? OR user_id = 0) AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id, userID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&Notification{}).
		Where("id = ? AND (user_id = ? OR user_id = 0)", id, userID).
		Updates(map[string]any{"is_read": true, "read_at": now}).Error
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Model(&Notification{}).
		Where("(user_id = ? OR user_id = 0) AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

func (r *NotificationRepository) Delete(ctx context.Context, id, userID int64) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND (user_id = ? OR user_id = 0)", id, userID).
		Delete(&Notification{}).Error
}

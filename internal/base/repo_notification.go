package base

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type InotificationRepository interface {
	Create(ctx context.Context, n *Notification) error
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

func (r *NotificationRepository) ListByUserID(ctx context.Context, userID int64, page, size int) ([]Notification, int64, error) {
	db := r.db.WithContext(ctx)

	var total int64
	if err := db.Model(&Notification{}).
		Where("user_id IN (0, ?)", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []Notification
	err := db.Table("base_notifications n").
		Select("n.*, IF(nr.id IS NOT NULL, TRUE, FALSE) as is_read").
		Joins("LEFT JOIN base_notification_reads nr ON nr.notification_id = n.id AND nr.user_id = ?", userID).
		Where("n.user_id IN (0, ?)", userID).
		Where("n.deleted_at IS NULL").
		Order("n.priority ASC, n.created_at DESC").
		Offset((page - 1) * size).Limit(size).
		Find(&list).Error
	return list, total, err
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)
		FROM base_notifications n
		LEFT JOIN base_notification_reads nr ON nr.notification_id = n.id AND nr.user_id = ?
		WHERE n.user_id IN (0, ?)
		  AND nr.id IS NULL
		  AND n.deleted_at IS NULL
	`, userID, userID).Scan(&count).Error
	return count, err
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id, userID int64) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT IGNORE INTO base_notification_reads (notification_id, user_id, read_at) VALUES (?, ?, ?)",
		id, userID, time.Now(),
	).Error
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Exec(`
		INSERT IGNORE INTO base_notification_reads (notification_id, user_id, read_at)
		SELECT n.id, ?, ?
		FROM base_notifications n
		LEFT JOIN base_notification_reads nr ON nr.notification_id = n.id AND nr.user_id = ?
		WHERE n.user_id IN (0, ?)
		  AND nr.id IS NULL
		  AND n.deleted_at IS NULL
	`, userID, time.Now(), userID, userID).Error
}

func (r *NotificationRepository) Delete(ctx context.Context, id, userID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&Notification{}, id).Error
}

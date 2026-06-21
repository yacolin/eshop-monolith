package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"

	"eshop-monolith/internal/notification/domain/models"
	infraModels "eshop-monolith/internal/infra/repository/models"
	"eshop-monolith/pkg/utils"
)

// InotificationRepository 通知仓储接口
type InotificationRepository interface {
	Create(ctx context.Context, n *models.Notification) error
	FindByID(ctx context.Context, id int64) (*models.Notification, error)
	ListByUserID(ctx context.Context, userID int64, page, size int) ([]models.Notification, int64, error)
	GetUnreadCount(ctx context.Context, userID int64) (int64, error)
	MarkAsRead(ctx context.Context, id, userID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
	Delete(ctx context.Context, id, userID int64) error
	DeleteByUserID(ctx context.Context, userID int64) error
}

// NotificationRepository 通知仓储实现
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository 创建通知仓储
func NewNotificationRepository(db *gorm.DB) InotificationRepository {
	return &NotificationRepository{db: db}
}

// Create 创建通知
func (r *NotificationRepository) Create(ctx context.Context, n *models.Notification) error {
	po := infraModels.NotificationFromDomain(n)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	n.ID = po.ID
	n.CreatedAt = utils.Timestamp(po.CreatedAt)
	return nil
}

// FindByID 根据 ID 获取通知
func (r *NotificationRepository) FindByID(ctx context.Context, id int64) (*models.Notification, error) {
	var po infraModels.NotificationPO
	if err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// ListByUserID 分页查询用户通知列表
func (r *NotificationRepository) ListByUserID(ctx context.Context, userID int64, page, size int) ([]models.Notification, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&infraModels.NotificationPO{}).
		Where("user_id = ? OR user_id = 0", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var pos []infraModels.NotificationPO
	if err := query.Order("created_at desc").
		Offset((page - 1) * size).
		Limit(size).
		Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	notifications := make([]models.Notification, len(pos))
	for i, po := range pos {
		notifications[i] = *po.ToDomain()
	}
	return notifications, total, nil
}

// GetUnreadCount 获取用户未读通知数
func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&infraModels.NotificationPO{}).
		Where("(user_id = ? OR user_id = 0) AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// MarkAsRead 标记单条通知为已读
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id, userID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&infraModels.NotificationPO{}).
		Where("id = ? AND (user_id = ? OR user_id = 0)", id, userID).
		Updates(map[string]any{
			"is_read": true,
			"read_at": now,
		}).Error
}

// MarkAllAsRead 标记用户所有通知为已读
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&infraModels.NotificationPO{}).
		Where("(user_id = ? OR user_id = 0) AND is_read = ?", userID, false).
		Updates(map[string]any{
			"is_read": true,
			"read_at": now,
		}).Error
}

// Delete 删除通知
func (r *NotificationRepository) Delete(ctx context.Context, id, userID int64) error {
	return r.db.WithContext(ctx).Where("id = ? AND (user_id = ? OR user_id = 0)", id, userID).
		Delete(&infraModels.NotificationPO{}).Error
}

// DeleteByUserID 清空用户所有通知
func (r *NotificationRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? OR user_id = 0", userID).
		Delete(&infraModels.NotificationPO{}).Error
}

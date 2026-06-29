package user

import (
	"context"

	"gorm.io/gorm"
)

type IloginHistoryRepository interface {
	Create(ctx context.Context, history *LoginHistory) error
	GetByUserID(ctx context.Context, userID int64, page, size int) ([]LoginHistory, int64, error)
}

type LoginHistoryRepository struct {
	db *gorm.DB
}

func NewLoginHistoryRepository(db *gorm.DB) IloginHistoryRepository {
	return &LoginHistoryRepository{db: db}
}

func (r *LoginHistoryRepository) Create(ctx context.Context, history *LoginHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *LoginHistoryRepository) GetByUserID(ctx context.Context, userID int64, page, size int) ([]LoginHistory, int64, error) {
	var list []LoginHistory
	var total int64

	db := r.db.WithContext(ctx).Model(&LoginHistory{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := db.Offset(offset).Limit(size).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

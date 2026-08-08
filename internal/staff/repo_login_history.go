package staff

import (
	"context"

	"gorm.io/gorm"
)

type IstaffLoginHistoryRepository interface {
	Create(ctx context.Context, h *StaffLoginHistory) error
}

type StaffLoginHistoryRepository struct {
	db *gorm.DB
}

func NewStaffLoginHistoryRepository(db *gorm.DB) IstaffLoginHistoryRepository {
	return &StaffLoginHistoryRepository{db: db}
}

func (r *StaffLoginHistoryRepository) Create(ctx context.Context, h *StaffLoginHistory) error {
	return r.db.WithContext(ctx).Create(h).Error
}

package staff

import (
	"context"

	"gorm.io/gorm"
)

type IpermissionRepository interface {
	List(ctx context.Context, category string) ([]SysPermission, error)
}

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) IpermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) List(ctx context.Context, category string) ([]SysPermission, error) {
	db := r.db.WithContext(ctx).Model(&SysPermission{})
	if category != "" {
		db = db.Where("category = ?", category)
	}
	var perms []SysPermission
	err := db.Order("category ASC, sort_order ASC, id ASC").Find(&perms).Error
	return perms, err
}

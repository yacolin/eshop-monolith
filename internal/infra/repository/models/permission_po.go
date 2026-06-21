package models

import (
	"time"

	userDomain "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// PermissionPO 权限持久化对象
type PermissionPO struct {
	ID          int64          `gorm:"type:int;primaryKey"`
	Name        string         `gorm:"type:varchar(100);not null;uniqueIndex"`
	DisplayName string         `gorm:"type:varchar(100);not null"`
	Description string         `gorm:"type:text"`
	Resource    string         `gorm:"type:varchar(50);not null;index"`
	Action      string         `gorm:"type:varchar(50);not null;index"`
	Category    string         `gorm:"type:varchar(50)"`
	Sort        int            `gorm:"type:int;default:0;index:idx_permissions_sort"`
	Status      int            `gorm:"type:tinyint;default:1"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (PermissionPO) TableName() string { return "permissions" }

func (po *PermissionPO) ToDomain() *userDomain.Permission {
	return &userDomain.Permission{
		ID:          po.ID,
		Name:        po.Name,
		DisplayName: po.DisplayName,
		Description: po.Description,
		Resource:    po.Resource,
		Action:      po.Action,
		Category:    po.Category,
		Sort:        po.Sort,
		Status:      po.Status,
		CreatedAt:   utils.Timestamp(po.CreatedAt),
		UpdatedAt:   utils.Timestamp(po.UpdatedAt),
	}
}

func PermissionFromDomain(p *userDomain.Permission) *PermissionPO {
	return &PermissionPO{
		ID:          p.ID,
		Name:        p.Name,
		DisplayName: p.DisplayName,
		Description: p.Description,
		Resource:    p.Resource,
		Action:      p.Action,
		Category:    p.Category,
		Sort:        p.Sort,
		Status:      p.Status,
		CreatedAt:   time.Time(p.CreatedAt),
		UpdatedAt:   time.Time(p.UpdatedAt),
	}
}

package models

import (
	"time"

	userDomain "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// RolePO 角色持久化对象
type RolePO struct {
	ID          int64              `gorm:"type:int;primaryKey"`
	Name        string             `gorm:"type:varchar(50);not null;uniqueIndex"`
	DisplayName string             `gorm:"type:varchar(100);not null"`
	Description string             `gorm:"type:text"`
	Status      int                `gorm:"type:tinyint;default:1"`
	Sort        int                `gorm:"type:int;default:0"`
	IsSystem    bool               `gorm:"type:tinyint(1);default:0"`
	CreatedAt   time.Time          `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   time.Time          `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt     `gorm:"index"`
	Permissions []PermissionPO     `gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:PermissionID"`
}

func (RolePO) TableName() string { return "roles" }

func (po *RolePO) ToDomain() *userDomain.Role {
	permissions := make([]userDomain.Permission, len(po.Permissions))
	for i, p := range po.Permissions {
		permissions[i] = *p.ToDomain()
	}
	return &userDomain.Role{
		ID:          po.ID,
		Name:        po.Name,
		DisplayName: po.DisplayName,
		Description: po.Description,
		Status:      po.Status,
		Sort:        po.Sort,
		IsSystem:    po.IsSystem,
		CreatedAt:   utils.Timestamp(po.CreatedAt),
		UpdatedAt:   utils.Timestamp(po.UpdatedAt),
		Permissions: permissions,
	}
}

func RoleFromDomain(r *userDomain.Role) *RolePO {
	permissions := make([]PermissionPO, len(r.Permissions))
	for i, p := range r.Permissions {
		permissions[i] = *PermissionFromDomain(&p)
	}
	return &RolePO{
		ID:          r.ID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: r.Description,
		Status:      r.Status,
		Sort:        r.Sort,
		IsSystem:    r.IsSystem,
		CreatedAt:   time.Time(r.CreatedAt),
		UpdatedAt:   time.Time(r.UpdatedAt),
		Permissions: permissions,
	}
}

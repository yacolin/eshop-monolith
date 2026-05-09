package models

import (
	"time"

	userDomain "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// UserRolePO 用户角色关联持久化对象
type UserRolePO struct {
	ID        int64          `gorm:"type:int;primaryKey"`
	UserID    int64          `gorm:"type:int;not null;index:idx_user_role"`
	RoleID    int64          `gorm:"type:int;not null;index:idx_user_role"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Role      *RolePO        `gorm:"foreignKey:RoleID"`
}

func (UserRolePO) TableName() string { return "user_roles" }

func (po *UserRolePO) ToDomain() *userDomain.UserRole {
	var role *userDomain.Role
	if po.Role != nil {
		role = po.Role.ToDomain()
	}
	return &userDomain.UserRole{
		ID:        po.ID,
		UserID:    po.UserID,
		RoleID:    po.RoleID,
		CreatedAt: utils.Timestamp(po.CreatedAt),
		Role:      role,
	}
}

func UserRoleFromDomain(ur *userDomain.UserRole) *UserRolePO {
	return &UserRolePO{
		ID:        ur.ID,
		UserID:    ur.UserID,
		RoleID:    ur.RoleID,
		CreatedAt: time.Time(ur.CreatedAt),
	}
}

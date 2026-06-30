package user

import (
	"gorm.io/gorm"

	"eshop-monolith/pkg/utils"
)

type Permission struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex:uk_name;not null" json:"name"`
	DisplayName string         `gorm:"type:varchar(100);not null;default:''" json:"display_name"`
	Description string         `gorm:"type:varchar(255);not null;default:''" json:"description"`
	Resource    string         `gorm:"type:varchar(50);not null;index:idx_resource" json:"resource"`
	Action      string         `gorm:"type:varchar(50);not null;index:idx_action" json:"action"`
	Category    string         `gorm:"type:varchar(50);not null;default:''" json:"category"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	Status      int8           `gorm:"type:tinyint;not null;default:1;index:idx_status" json:"status"`
	CreatedAt utils.Timestamp      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt utils.Timestamp      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (Permission) TableName() string { return "usr_permissions" }

type RolePermission struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID       int64          `gorm:"not null;uniqueIndex:uk_role_permission" json:"role_id"`
	PermissionID int64          `gorm:"not null;uniqueIndex:uk_role_permission" json:"permission_id"`
	CreatedAt utils.Timestamp      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (RolePermission) TableName() string { return "usr_role_permissions" }

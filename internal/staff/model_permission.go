package staff

import (
	"time"

	"gorm.io/gorm"
)

// SysPermission B 端权限(sys_permissions)
type SysPermission struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex:uk_name;not null" json:"name"`
	DisplayName string         `gorm:"type:varchar(100);not null;default:''" json:"display_name"`
	Description string         `gorm:"type:varchar(255);not null;default:''" json:"description"`
	Resource    string         `gorm:"type:varchar(50);not null;index:idx_resource" json:"resource"`
	Action      string         `gorm:"type:varchar(50);not null;index:idx_action" json:"action"`
	ParentID    int64          `gorm:"not null;default:0" json:"parent_id"`
	Category    string         `gorm:"type:varchar(50);not null;default:''" json:"category"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	Status      int8           `gorm:"type:tinyint(1);not null;default:1;index:idx_status" json:"status"`
	CreatedAt   time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (SysPermission) TableName() string { return "sys_permissions" }

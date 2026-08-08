package staff

import (
	"time"

	"gorm.io/gorm"
)

// SysStaffRole 员工-角色关联(sys_staff_roles)
type SysStaffRole struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	StaffID   int64          `gorm:"not null;uniqueIndex:uk_staff_role" json:"staff_id"`
	RoleID    int64          `gorm:"not null;uniqueIndex:uk_staff_role" json:"role_id"`
	CreatedAt time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (SysStaffRole) TableName() string { return "sys_staff_roles" }

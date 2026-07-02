package user

import (
	"gorm.io/gorm"

	"eshop-monolith/pkg/utils"
)

type Role struct {
	ID          int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string          `gorm:"type:varchar(50);uniqueIndex:uk_name;not null" json:"name"`
	DisplayName string          `gorm:"type:varchar(100);not null;default:''" json:"display_name"`
	Description string          `gorm:"type:varchar(255);not null;default:''" json:"description"`
	RoleType    string          `gorm:"type:varchar(20);not null;default:'custom'" json:"role_type"`
	Status      int8            `gorm:"type:tinyint;not null;default:1;index:idx_status" json:"status"`
	SortOrder   int             `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt   utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt   utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (Role) TableName() string { return "usr_roles" }

type UserRole struct {
	ID        int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64           `gorm:"not null;uniqueIndex:uk_user_role" json:"user_id"`
	RoleID    int64           `gorm:"not null;uniqueIndex:uk_user_role" json:"role_id"`
	CreatedAt utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	DeletedAt gorm.DeletedAt  `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (UserRole) TableName() string { return "usr_user_roles" }

package product

import (
	"eshop-monolith/pkg/utils"

	"gorm.io/gorm"
)

type Category struct {
	ID        int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string          `gorm:"type:varchar(100);not null" json:"name"`
	ParentID  int64           `gorm:"not null;default:0;index:idx_parent_id" json:"parent_id"`
	Level     int8            `gorm:"not null;default:1;index:idx_level_status" json:"level"`
	Path      string          `gorm:"type:varchar(500);not null;default:'';index:idx_path" json:"path"`
	IconURL   string          `gorm:"type:varchar(512);default:''" json:"icon_url"`
	SortOrder int             `gorm:"not null;default:0" json:"sort_order"`
	Status    int8            `gorm:"not null;default:1;index:idx_level_status" json:"status"`
	CreatedAt utils.Timestamp `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt utils.Timestamp `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index:idx_deleted_at" json:"-"`
}

func (Category) TableName() string { return "sp_categories" }

package product

import (
	"time"

	"gorm.io/gorm"
)

type Brand struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;uniqueIndex:uk_name" json:"name"`
	EnglishName string         `gorm:"type:varchar(100);default:''" json:"english_name"`
	LogoURL     string         `gorm:"type:varchar(512);default:''" json:"logo_url"`
	FirstLetter string         `gorm:"type:char(1);default:'';index:idx_first_letter" json:"first_letter"`
	SortOrder   int            `gorm:"not null;default:0" json:"sort_order"`
	Status      int8           `gorm:"not null;default:1;index:idx_status" json:"status"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index:idx_deleted_at" json:"-"`
}

func (Brand) TableName() string { return "sp_brands" }

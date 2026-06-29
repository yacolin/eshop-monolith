package product

import "time"

type CategoryBrand struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID int64     `gorm:"not null;uniqueIndex:uk_category_brand" json:"category_id"`
	BrandID    int64     `gorm:"not null;uniqueIndex:uk_category_brand;index:idx_brand_id" json:"brand_id"`
	SortOrder  int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (CategoryBrand) TableName() string { return "sp_category_brands" }

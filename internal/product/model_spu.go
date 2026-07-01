package product

import (
	"eshop-monolith/pkg/utils"

	"gorm.io/gorm"
)

type SPU struct {
	ID             int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string          `gorm:"type:varchar(200);not null" json:"name"`
	Subtitle       string          `gorm:"type:varchar(500);default:''" json:"subtitle"`
	CategoryID     int64           `gorm:"not null;default:0;index:idx_category_status_sort" json:"category_id"`
	BrandID        int64           `gorm:"not null;default:0;index:idx_brand_status" json:"brand_id"`
	Unit           string          `gorm:"type:varchar(10);default:'件'" json:"unit"`
	MainImage      string          `gorm:"type:varchar(512);not null;default:''" json:"main_image"`
	Images         string          `gorm:"type:json" json:"images,omitempty"`
	VideoURL       string          `gorm:"type:varchar(512);default:''" json:"video_url"`
	MinPrice       int64           `gorm:"not null;default:0" json:"min_price"`
	MaxPrice       int64           `gorm:"not null;default:0" json:"max_price"`
	TotalStock     int             `gorm:"not null;default:0" json:"total_stock"`
	SalesCount     int             `gorm:"not null;default:0" json:"sales_count"`
	RatingAverage  float64         `gorm:"type:decimal(3,2);not null;default:0.00" json:"rating_average"`
	RatingCount    int             `gorm:"not null;default:0" json:"rating_count"`
	Status         int8            `gorm:"not null;default:0" json:"status"`
	SortOrder      int             `gorm:"not null;default:0" json:"sort_order"`
	HasDescription int8            `gorm:"not null;default:0" json:"has_description"`
	CreatedBy      string          `gorm:"type:varchar(50);default:''" json:"created_by"`
	UpdatedBy      string          `gorm:"type:varchar(50);default:''" json:"updated_by"`
	CreatedAt      utils.Timestamp `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      utils.Timestamp `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index:idx_deleted_at" json:"-"`
}

func (SPU) TableName() string { return "sp_products" }

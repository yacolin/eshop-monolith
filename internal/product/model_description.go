package product

import "time"

type Description struct {
	ID                int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID         int64     `gorm:"not null;uniqueIndex:uk_product_id" json:"product_id"`
	Description       string    `gorm:"type:longtext" json:"description"`
	MobileDescription string    `gorm:"type:longtext" json:"mobile_description"`
	CreatedAt         time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Description) TableName() string { return "sp_product_descriptions" }

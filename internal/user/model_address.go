package user

import (
	"gorm.io/gorm"

	"eshop-monolith/pkg/utils"
)

type Address struct {
	ID        int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64           `gorm:"not null;index:idx_user_id" json:"user_id"`
	Consignee string          `gorm:"type:varchar(64);not null;default:''" json:"consignee"`
	Phone     string          `gorm:"type:varchar(20);not null;default:''" json:"phone"`
	Country   string          `gorm:"type:varchar(32);not null;default:''" json:"country"`
	Province  string          `gorm:"type:varchar(32);not null;default:''" json:"province"`
	City      string          `gorm:"type:varchar(32);not null;default:''" json:"city"`
	District  string          `gorm:"type:varchar(32);not null;default:''" json:"district"`
	Detail    string          `gorm:"type:varchar(256);not null;default:''" json:"detail"`
	ZipCode   string          `gorm:"type:varchar(10);not null;default:''" json:"zip_code"`
	Tag       string          `gorm:"type:varchar(16);not null;default:''" json:"tag"`
	IsDefault bool            `gorm:"not null;default:false" json:"is_default"`
	CreatedAt utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"type:datetime(3);index:idx_deleted_at" json:"-"`
}

func (Address) TableName() string { return "usr_addresses" }

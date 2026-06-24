package models

import (
	"time"

	addressDomain "eshop-monolith/internal/address/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// AddressPO 地址持久化对象
type AddressPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"not null;index"`
	Consignee string         `gorm:"type:varchar(64);not null"`
	Phone     string         `gorm:"type:varchar(20);not null"`
	Province  string         `gorm:"type:varchar(32);not null"`
	City      string         `gorm:"type:varchar(32);not null"`
	District  string         `gorm:"type:varchar(32);not null"`
	Detail    string         `gorm:"type:varchar(256);not null"`
	ZipCode   string         `gorm:"type:varchar(10)"`
	IsDefault bool           `gorm:"not null;default:false"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (AddressPO) TableName() string { return "addresses" }

// ToDomain 转换为领域模型
func (po *AddressPO) ToDomain() *addressDomain.Address {
	return &addressDomain.Address{
		ID:        po.ID,
		UserID:    po.UserID,
		Consignee: po.Consignee,
		Phone:     po.Phone,
		Province:  po.Province,
		City:      po.City,
		District:  po.District,
		Detail:    po.Detail,
		ZipCode:   po.ZipCode,
		IsDefault: po.IsDefault,
		CreatedAt: utils.Timestamp(po.CreatedAt),
		UpdatedAt: utils.Timestamp(po.UpdatedAt),
	}
}

// AddressFromDomain 从领域模型创建持久化对象
func AddressFromDomain(a *addressDomain.Address) *AddressPO {
	return &AddressPO{
		ID:        a.ID,
		UserID:    a.UserID,
		Consignee: a.Consignee,
		Phone:     a.Phone,
		Province:  a.Province,
		City:      a.City,
		District:  a.District,
		Detail:    a.Detail,
		ZipCode:   a.ZipCode,
		IsDefault: a.IsDefault,
		CreatedAt: time.Time(a.CreatedAt),
		UpdatedAt: time.Time(a.UpdatedAt),
	}
}

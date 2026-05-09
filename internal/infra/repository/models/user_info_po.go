package models

import (
	"time"

	userDomain "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// UserInfoPO 用户详细信息持久化对象
type UserInfoPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"not null;uniqueIndex"`
	Avatar    string         `gorm:"type:varchar(255)"`
	Nickname  string         `gorm:"type:varchar(50)"`
	Gender    int            `gorm:"type:tinyint;default:0"`
	Birthday  *time.Time     `gorm:"type:timestamp"`
	Address   string         `gorm:"type:varchar(255)"`
	Bio       string         `gorm:"type:text"`
	Country   string         `gorm:"type:varchar(50)"`
	Province  string         `gorm:"type:varchar(50)"`
	City      string         `gorm:"type:varchar(50)"`
	ZipCode   string         `gorm:"type:varchar(20)"`
	Language  string         `gorm:"type:varchar(20);default:zh-CN"`
	Timezone  string         `gorm:"type:varchar(50);default:Asia/Shanghai"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserInfoPO) TableName() string { return "user_infos" }

func (po *UserInfoPO) ToDomain() *userDomain.UserInfo {
	return &userDomain.UserInfo{
		ID:        po.ID,
		UserID:    po.UserID,
		Avatar:    po.Avatar,
		Nickname:  po.Nickname,
		Gender:    po.Gender,
		Birthday:  po.Birthday,
		Address:   po.Address,
		Bio:       po.Bio,
		Country:   po.Country,
		Province:  po.Province,
		City:      po.City,
		ZipCode:   po.ZipCode,
		Language:  po.Language,
		Timezone:  po.Timezone,
		CreatedAt: utils.Timestamp(po.CreatedAt),
		UpdatedAt: utils.Timestamp(po.UpdatedAt),
	}
}

func UserInfoFromDomain(u *userDomain.UserInfo) *UserInfoPO {
	return &UserInfoPO{
		ID:        u.ID,
		UserID:    u.UserID,
		Avatar:    u.Avatar,
		Nickname:  u.Nickname,
		Gender:    u.Gender,
		Birthday:  u.Birthday,
		Address:   u.Address,
		Bio:       u.Bio,
		Country:   u.Country,
		Province:  u.Province,
		City:      u.City,
		ZipCode:   u.ZipCode,
		Language:  u.Language,
		Timezone:  u.Timezone,
		CreatedAt: time.Time(u.CreatedAt),
		UpdatedAt: time.Time(u.UpdatedAt),
	}
}

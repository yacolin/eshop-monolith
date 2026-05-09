package models

import (
	"time"

	userDomain "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// UserPO 用户持久化对象
type UserPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Status    int            `gorm:"type:tinyint;default:1"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserInfo  *UserInfoPO    `gorm:"foreignKey:UserID"`
}

func (UserPO) TableName() string { return "users" }

func (po *UserPO) ToDomain() *userDomain.User {
	var userInfo *userDomain.UserInfo
	if po.UserInfo != nil {
		userInfo = po.UserInfo.ToDomain()
	}
	return &userDomain.User{
		ID:        po.ID,
		Status:    po.Status,
		CreatedAt: utils.Timestamp(po.CreatedAt),
		UpdatedAt: utils.Timestamp(po.UpdatedAt),
		UserInfo:  userInfo,
	}
}

func UserFromDomain(u *userDomain.User) *UserPO {
	var userInfo *UserInfoPO
	if u.UserInfo != nil {
		userInfo = UserInfoFromDomain(u.UserInfo)
	}
	return &UserPO{
		ID:        u.ID,
		Status:    u.Status,
		CreatedAt: time.Time(u.CreatedAt),
		UpdatedAt: time.Time(u.UpdatedAt),
		UserInfo:  userInfo,
	}
}

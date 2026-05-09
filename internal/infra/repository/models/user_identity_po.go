package models

import (
	"time"

	userDomain "eshop-monolith/internal/user/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// UserIdentityPO 用户身份凭证持久化对象
type UserIdentityPO struct {
	ID         int64          `gorm:"primaryKey;autoIncrement"`
	UserID     int64          `gorm:"not null;index"`
	Provider   string         `gorm:"type:varchar(50);not null;index:idx_provider_identifier,unique"`
	Identifier string         `gorm:"type:varchar(255);not null;index:idx_provider_identifier,unique"`
	Credential string         `gorm:"type:text"`
	Verified   bool           `gorm:"default:false"`
	Meta       string         `gorm:"type:json"`
	CreatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt  time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	User       *UserPO        `gorm:"foreignKey:UserID"`
}

func (UserIdentityPO) TableName() string { return "user_identities" }

func (po *UserIdentityPO) ToDomain() *userDomain.UserIdentity {
	var user *userDomain.User
	if po.User != nil {
		user = po.User.ToDomain()
	}
	return &userDomain.UserIdentity{
		ID:         po.ID,
		UserID:     po.UserID,
		Provider:   po.Provider,
		Identifier: po.Identifier,
		Credential: po.Credential,
		Verified:   po.Verified,
		Meta:       po.Meta,
		CreatedAt:  utils.Timestamp(po.CreatedAt),
		UpdatedAt:  utils.Timestamp(po.UpdatedAt),
		User:       user,
	}
}

func UserIdentityFromDomain(u *userDomain.UserIdentity) *UserIdentityPO {
	var user *UserPO
	if u.User != nil {
		user = UserFromDomain(u.User)
	}
	return &UserIdentityPO{
		ID:         u.ID,
		UserID:     u.UserID,
		Provider:   u.Provider,
		Identifier: u.Identifier,
		Credential: u.Credential,
		Verified:   u.Verified,
		Meta:       u.Meta,
		CreatedAt:  time.Time(u.CreatedAt),
		UpdatedAt:  time.Time(u.UpdatedAt),
		User:       user,
	}
}

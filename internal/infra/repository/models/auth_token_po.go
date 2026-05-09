package models

import (
	"time"

	userDomain "eshop-monolith/internal/user/domain/models"
	"gorm.io/gorm"
)

// AuthTokenPO 认证令牌持久化对象
type AuthTokenPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"not null;index"`
	JTI       string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	TokenType string         `gorm:"type:varchar(50);not null"`
	ExpiresAt time.Time      `gorm:"type:timestamp"`
	Revoked   bool           `gorm:"default:false"`
	RevokedAt *time.Time     `gorm:"type:timestamp"`
	Meta      string         `gorm:"type:json"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      *UserPO        `gorm:"foreignKey:UserID"`
}

func (AuthTokenPO) TableName() string { return "auth_tokens" }

func (po *AuthTokenPO) ToDomain() *userDomain.AuthToken {
	var user *userDomain.User
	if po.User != nil {
		user = po.User.ToDomain()
	}
	return &userDomain.AuthToken{
		ID:        po.ID,
		UserID:    po.UserID,
		JTI:       po.JTI,
		TokenType: po.TokenType,
		ExpiresAt: po.ExpiresAt,
		Revoked:   po.Revoked,
		RevokedAt: po.RevokedAt,
		Meta:      po.Meta,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
		User:      user,
	}
}

func AuthTokenFromDomain(t *userDomain.AuthToken) *AuthTokenPO {
	return &AuthTokenPO{
		ID:        t.ID,
		UserID:    t.UserID,
		JTI:       t.JTI,
		TokenType: t.TokenType,
		ExpiresAt: t.ExpiresAt,
		Revoked:   t.Revoked,
		RevokedAt: t.RevokedAt,
		Meta:      t.Meta,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

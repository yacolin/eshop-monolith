package models

import (
	"time"

	userDomain "eshop-monolith/internal/user/domain/models"
)

// LoginHistoryPO 登录历史记录持久化对象
type LoginHistoryPO struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	UserID     int64     `gorm:"not null;index"`
	IdentityID int64     `gorm:"not null;index"`
	Provider   string    `gorm:"type:varchar(50);not null"`
	IP         string    `gorm:"type:varchar(50)"`
	UserAgent  string    `gorm:"type:varchar(500)"`
	DeviceID   string    `gorm:"type:varchar(100)"`
	Event      string    `gorm:"type:varchar(50);not null"`
	Status     string    `gorm:"type:varchar(20);not null"`
	FailReason string    `gorm:"type:varchar(255)"`
	CreatedAt  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
}

func (LoginHistoryPO) TableName() string { return "login_histories" }

func (po *LoginHistoryPO) ToDomain() *userDomain.LoginHistory {
	return &userDomain.LoginHistory{
		ID:         po.ID,
		UserID:     po.UserID,
		IdentityID: po.IdentityID,
		Provider:   po.Provider,
		IP:         po.IP,
		UserAgent:  po.UserAgent,
		DeviceID:   po.DeviceID,
		Event:      po.Event,
		Status:     po.Status,
		FailReason: po.FailReason,
		CreatedAt:  po.CreatedAt,
	}
}

func LoginHistoryFromDomain(h *userDomain.LoginHistory) *LoginHistoryPO {
	return &LoginHistoryPO{
		ID:         h.ID,
		UserID:     h.UserID,
		IdentityID: h.IdentityID,
		Provider:   h.Provider,
		IP:         h.IP,
		UserAgent:  h.UserAgent,
		DeviceID:   h.DeviceID,
		Event:      h.Event,
		Status:     h.Status,
		FailReason: h.FailReason,
		CreatedAt:  h.CreatedAt,
	}
}

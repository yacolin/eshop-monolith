package models

import (
	"time"

	"gorm.io/gorm"
)

// AuthToken 认证令牌表，用于保存 refresh token 等，便于注销、撤销
type AuthToken struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64      `gorm:"not null;index" json:"user_id"`
	JTI       string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"jti"` // JWT ID，唯一标识
	TokenType string     `gorm:"type:varchar(50);not null" json:"token_type"`       // token类型：refresh_token, access_token等
	ExpiresAt time.Time  `json:"expires_at"`                                        // 过期时间
	Revoked   bool       `gorm:"default:false" json:"revoked"`                      // 是否已撤销
	RevokedAt *time.Time `json:"revoked_at,omitempty"`                              // 撤销时间
	Meta      string     `gorm:"type:json" json:"meta,omitempty"`                   // 额外元数据（JSON格式）

	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (AuthToken) TableName() string {
	return "auth_tokens"
}

// TokenType 常量定义
const (
	TokenTypeAccessToken  = "access_token"  // 访问令牌
	TokenTypeRefreshToken = "refresh_token" // 刷新令牌
)

// AuthTokenMeta Token 元数据结构
type AuthTokenMeta struct {
	IP        string `json:"ip,omitempty"`         // 登录IP
	UserAgent string `json:"user_agent,omitempty"` // 用户代理
	DeviceID  string `json:"device_id,omitempty"`  // 设备ID
	Source    string `json:"source,omitempty"`     // 来源：web, miniapp, app等
}

// LoginHistory 登录历史记录表（可选，用于安全审计）
type LoginHistory struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"not null;index" json:"user_id"`
	IdentityID int64     `gorm:"not null;index" json:"identity_id"`         // 使用的身份凭证ID
	Provider   string    `gorm:"type:varchar(50);not null" json:"provider"` // 登录方式
	IP         string    `gorm:"type:varchar(50)" json:"ip"`
	UserAgent  string    `gorm:"type:varchar(500)" json:"user_agent"`
	DeviceID   string    `gorm:"type:varchar(100)" json:"device_id"`
	Event      string    `gorm:"type:varchar(50);not null" json:"event"`         // 事件类型：login, logout, refresh等
	Status     string    `gorm:"type:varchar(20);not null" json:"status"`        // 状态：success, failed
	FailReason string    `gorm:"type:varchar(255)" json:"fail_reason,omitempty"` // 失败原因
	CreatedAt  time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
}

func (LoginHistory) TableName() string {
	return "login_histories"
}

// LoginEvent 常量定义
const (
	LoginEventLogin    = "login"    // 登录
	LoginEventLogout   = "logout"   // 登出
	LoginEventRefresh  = "refresh"  // 刷新token
	LoginEventRevoke   = "revoke"   // 撤销token
	LoginEventPassword = "password" // 修改密码
)

// LoginStatus 常量定义
const (
	LoginStatusSuccess = "success" // 成功
	LoginStatusFailed  = "failed"  // 失败
)

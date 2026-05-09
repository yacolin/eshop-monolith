package models

import "eshop-monolith/pkg/utils"

// UserIdentity 用户身份凭证表，支持多种登录方式（密码、微信、手机、邮箱等）
type UserIdentity struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Provider   string `json:"provider"`    // 登录方式：password, wechat, phone, email, github等
	Identifier string `json:"identifier"` // 唯一标识：用户名、openid、手机号、邮箱等
	Credential string `json:"-"`                                                                // 凭证：bcrypt密码hash、加密的session_key等
	Verified   bool   `json:"verified"`      // 是否已验证（手机/邮箱）
	Meta       string `json:"meta,omitempty"` // 额外元数据：unionid、session_key_iv、source等（JSON格式）

	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`

	// 关联
	User *User `json:"user,omitempty"`
}

func (UserIdentity) TableName() string {
	return "user_identities"
}

// Provider 常量定义
type Provider string

const (
	ProviderPassword Provider = "password" // 用户名/密码登录
	ProviderWechat   Provider = "wechat"   // 微信登录（小程序/公众号）
	ProviderPhone    Provider = "phone"    // 手机号验证码登录
	ProviderEmail    Provider = "email"    // 邮箱验证码登录
	ProviderGithub   Provider = "github"   // GitHub OAuth登录
	ProviderGoogle   Provider = "google"   // Google OAuth登录
)

// String 返回 provider 字符串
func (p Provider) String() string {
	return string(p)
}

// IsValid 检查 provider 是否有效
func (p Provider) IsValid() bool {
	switch p {
	case ProviderPassword, ProviderWechat, ProviderPhone, ProviderEmail, ProviderGithub, ProviderGoogle:
		return true
	}
	return false
}

// IdentityMeta 身份元数据结构
type IdentityMeta struct {
	UnionID      string `json:"unionid,omitempty"`        // 微信unionid，用于跨应用统一登录
	SessionKey   string `json:"session_key,omitempty"`    // 微信session_key（建议加密存储）
	SessionKeyIV string `json:"session_key_iv,omitempty"` // session_key加密IV
	Source       string `json:"source,omitempty"`         // 来源：miniapp, h5, app, web等
	AppID        string `json:"appid,omitempty"`          // 微信appid
}

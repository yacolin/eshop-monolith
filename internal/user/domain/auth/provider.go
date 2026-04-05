package auth

import (
	"context"
	"eshop-monolith/internal/user/domain/models"
)

// AuthProvider 认证提供者接口，支持多种登录方式
type AuthProvider interface {
	// GetName 获取提供者名称
	GetName() string
	// Authenticate 执行认证，返回用户和身份凭证
	Authenticate(ctx context.Context, payload interface{}) (*models.User, *models.UserIdentity, error)
}

// PasswordPayload 用户名密码登录参数
type PasswordPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// WechatPayload 微信登录参数
type WechatPayload struct {
	Code          string `json:"code"`                     // 微信临时登录凭证
	AppID         string `json:"appid"`                    // 微信应用ID
	Source        string `json:"source"`                   // 来源：miniapp, h5, app等
	EncryptedData string `json:"encrypted_data,omitempty"` // 加密数据（可选）
	IV            string `json:"iv,omitempty"`             // 加密算法的初始向量（可选）
}

// WechatSession 微信session响应
type WechatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// PhonePayload 手机号验证码登录参数
type PhonePayload struct {
	Phone      string `json:"phone"`
	VerifyCode string `json:"verify_code"`
}

// EmailPayload 邮箱验证码登录参数
type EmailPayload struct {
	Email      string `json:"email"`
	VerifyCode string `json:"verify_code"`
}

// RegisterPayload 用户注册参数
type RegisterPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Provider string `json:"provider"` // 注册方式
}

// BindPayload 绑定身份凭证参数
type BindPayload struct {
	UserID     string `json:"user_id"`
	Provider   string `json:"provider"`
	Identifier string `json:"identifier"`
	Credential string `json:"credential"`
	Meta       string `json:"meta,omitempty"`
}

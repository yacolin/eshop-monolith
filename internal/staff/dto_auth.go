package staff

import (
	"time"
)

// ── Request ──

type StaffLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ── Response ──

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
}

type StaffLoginResponse struct {
	StaffID  int64    `json:"staff_id"`
	Username string   `json:"username"`
	RealName string   `json:"real_name"`
	Avatar   string   `json:"avatar"`
	Roles    []string `json:"roles"`
	TokenResponse
}

type StaffProfileResponse struct {
	StaffID   int64     `json:"staff_id"`
	Username  string    `json:"username"`
	RealName  string    `json:"real_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Avatar    string    `json:"avatar"`
	Status    int8      `json:"status"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
}

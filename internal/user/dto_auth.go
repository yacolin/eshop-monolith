package user

type PasswordLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type WechatLoginReq struct {
	Code   string `json:"code" binding:"required"`
	AppID  string `json:"appid" binding:"required"`
	Source string `json:"source"`
}

type PhoneLoginReq struct {
	Phone      string `json:"phone" binding:"required"`
	VerifyCode string `json:"verify_code" binding:"required"`
}

type EmailLoginReq struct {
	Email      string `json:"email" binding:"required,email"`
	VerifyCode string `json:"verify_code" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Provider string `json:"provider" binding:"required"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
}

type LoginResponse struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username,omitempty"`
	TokenResponse
	IsNewUser bool `json:"is_new_user"`
}

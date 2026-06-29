package user

// 用户更新个人信息请求
type UpdateUserInfoReq struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=512"`
	Gender   int8   `json:"gender" binding:"omitempty,oneof=0 1 2"`
	Birthday string `json:"birthday"`
	Bio      string `json:"bio" binding:"max=200"`
	Country  string `json:"country" binding:"max=32"`
	Province string `json:"province" binding:"max=32"`
	City     string `json:"city" binding:"max=32"`
	ZipCode  string `json:"zip_code" binding:"max=10"`
	Language string `json:"language" binding:"max=10"`
	Timezone string `json:"timezone" binding:"max=32"`
}

// 用户资料响应
type UserProfileResponse struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Phone         string `json:"phone"`
	PhoneVerified bool   `json:"phone_verified"`
	Avatar        string `json:"avatar"`
	Nickname      string `json:"nickname"`
	Status        int8   `json:"status"`
	UserInfo      *UserInfoResponse `json:"user_info,omitempty"`
}

type UserInfoResponse struct {
	Gender   int8   `json:"gender"`
	Birthday string `json:"birthday,omitempty"`
	Bio      string `json:"bio"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	ZipCode  string `json:"zip_code"`
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}

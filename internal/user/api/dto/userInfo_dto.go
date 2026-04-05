package dto

type UpdateUserInfoRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=255"`
	Gender   int    `json:"gender"`
	Birthday string `json:"birthday"`
	Address  string `json:"address" binding:"max=255"`
	Bio      string `json:"bio" binding:"max=500"`
	Country  string `json:"country" binding:"max=50"`
	Province string `json:"province" binding:"max=50"`
	City     string `json:"city" binding:"max=50"`
	ZipCode  string `json:"zip_code" binding:"max=20"`
	Language string `json:"language" binding:"max=20"`
	Timezone string `json:"timezone" binding:"max=50"`
}

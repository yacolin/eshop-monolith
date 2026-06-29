package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
)

type UserHandler struct {
	userSvc *UserService
}

func NewUserHandler(userSvc *UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GetProfile 获取当前用户资料
// @Summary 获取当前用户资料
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=UserProfileResponse}
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	user, info, err := h.userSvc.GetProfile(c, currentUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	resp := &UserProfileResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Phone:         user.Phone,
		PhoneVerified: user.PhoneVerified,
		Avatar:        user.Avatar,
		Nickname:      user.Nickname,
		Status:        user.Status,
	}
	if info != nil {
		birthday := ""
		if info.Birthday != nil {
			birthday = info.Birthday.Format("2006-01-02")
		}
		resp.UserInfo = &UserInfoResponse{
			Gender:   info.Gender,
			Birthday: birthday,
			Bio:      info.Bio,
			Country:  info.Country,
			Province: info.Province,
			City:     info.City,
			ZipCode:  info.ZipCode,
			Language: info.Language,
			Timezone: info.Timezone,
		}
	}
	response.Success(c, resp)
}

// UpdateUserInfo 更新个人信息
// @Summary 更新个人信息
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body UpdateUserInfoReq true "个人信息"
// @Success 200 {object} response.Response
// @Router /api/v1/users/info [put]
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	var req UpdateUserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.userSvc.UpdateUserInfo(c, currentUserID(c), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ── Routes ────────────────────────────────────────

func RegisterUserRoutes(v1 *gin.RouterGroup, db *gorm.DB, userRepo IuserRepository, infoRepo IuserInfoRepository) {
	userSvc := NewUserService(userRepo, infoRepo)
	h := NewUserHandler(userSvc)

	users := v1.Group("/users")
	users.Use(middleware.JWTAuth())
	{
		users.GET("/profile", h.GetProfile)
		users.PUT("/info", h.UpdateUserInfo)
	}
}

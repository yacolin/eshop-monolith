package handlers

import (
	"strconv"

	"eshop-monolith/internal/user/api/dto"
	"eshop-monolith/internal/user/service"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的完整资料（包含 User 和 UserInfo）
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=models.User}
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	user, err := h.userSvc.GetProfile(c, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, user)
}

// GetUserInfo 获取用户详细信息
// @Summary 获取用户详细信息
// @Description 获取当前登录用户的详细信息
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=models.UserInfo}
// @Router /api/v1/users/info [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	userInfo, err := h.userSvc.GetUserInfo(c, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, userInfo)
}

// UpdateUserInfo 更新用户详细信息
// @Summary 更新用户详细信息
// @Description 更新当前登录用户的详细信息（Avatar、Nickname 等）
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserInfoRequest true "用户信息"
// @Success 200 {object} response.Response{data=models.UserInfo}
// @Router /api/v1/users/info [put]
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	userInfo, err := h.userSvc.UpdateUserInfo(c, userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, userInfo)
}

// GetByID 根据ID获取用户信息
// @Summary 根据ID获取用户
// @Description 根据用户ID获取用户信息（管理员接口）
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Success 200 {object} response.Response{data=models.User}
// @Router /api/v1/users/{user_id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	userID, err := utils.ParseIntParam(c, "user_id")
	if err != nil {
		return
	}

	user, err := h.userSvc.GetByID(c, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, user)
}

func getUserIDFromContext(c *gin.Context) (int64, error) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, errcode.ErrUnauthorized
	}

	switch id := v.(type) {
	case uint:
		return int64(id), nil
	case int:
		return int64(id), nil
	case int64:
		return id, nil
	case float64:
		return int64(id), nil
	case string:
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return 0, errcode.ErrUnauthorized
		}
		return n, nil
	default:
		return 0, errcode.ErrUnauthorized
	}
}

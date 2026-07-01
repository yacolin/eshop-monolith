package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type UserHandler struct {
	userSvc *UserService
	roleSvc *RoleService
}

func NewUserHandler(userSvc *UserService, roleSvc *RoleService) *UserHandler {
	return &UserHandler{userSvc: userSvc, roleSvc: roleSvc}
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

// List 后台用户列表
// @Summary 后台用户列表
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(20)
// @Param status query int false "状态"
// @Param keyword query string false "关键字（用户名/邮箱/手机号/昵称）"
// @Success 200 {object} response.Response{data=UserListResult}
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
	var req UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	users, total, err := h.userSvc.List(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]*UserListItem, len(users))
	for i := range users {
		roles, _ := h.roleSvc.GetUserRoles(c, users[i].ID)
		roleBriefs := make([]UserRoleBrief, len(roles))
		for j := range roles {
			roleBriefs[j] = UserRoleBrief{ID: roles[j].ID, Name: roles[j].Name, DisplayName: roles[j].DisplayName}
		}
		items[i] = &UserListItem{User: &users[i], Roles: roleBriefs}
	}
	response.Success(c, &UserListResult{Total: total, List: items})
}

// AssignRole 给用户分配角色
// @Summary 给用户分配角色
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Param request body AssignRoleReq true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{user_id}/roles [post]
func (h *UserHandler) AssignRole(c *gin.Context) {
	userID, err := utils.ParseIntParam(c, "user_id")
	if err != nil {
		c.Error(err)
		return
	}
	var req AssignRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.roleSvc.AssignToUser(c, userID, req.RoleID); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// RemoveRole 移除用户的角色
// @Summary 移除用户的角色
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Param user_id path int true "用户ID"
// @Param role_id path int true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{user_id}/roles/{role_id} [delete]
func (h *UserHandler) RemoveRole(c *gin.Context) {
	userID, err := utils.ParseIntParam(c, "user_id")
	if err != nil {
		c.Error(err)
		return
	}
	roleID, err := utils.ParseIntParam(c, "role_id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.roleSvc.RemoveFromUser(c, userID, roleID); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ── Routes ────────────────────────────────────────

func RegisterUserRoutes(v1 *gin.RouterGroup, db *gorm.DB, userRepo IuserRepository, infoRepo IuserInfoRepository, roleRepo IroleRepository) {
	userSvc := NewUserService(userRepo, infoRepo)
	roleSvc := NewRoleService(roleRepo)
	h := NewUserHandler(userSvc, roleSvc)
	roleCfg := NewRequireRoleConfig(roleRepo)

	users := v1.Group("/users")
	users.Use(middleware.JWTAuth())
	{
		users.GET("/profile", h.GetProfile)
		users.PUT("/info", h.UpdateUserInfo)
	}

	admin := v1.Group("/users")
	admin.Use(middleware.JWTAuth(), RequireAdmin(roleCfg))
	{
		admin.GET("", h.List)
		admin.POST("/:user_id/roles", h.AssignRole)
		admin.DELETE("/:user_id/roles/:role_id", h.RemoveRole)
	}
}

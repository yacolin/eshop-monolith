package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type RoleHandler struct {
	roleSvc *RoleService
}

func NewRoleHandler(roleSvc *RoleService) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc}
}

// Create 创建角色
// @Summary 创建角色
// @Tags roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreateRoleReq true "角色信息"
// @Success 200 {object} response.Response{data=Role}
// @Router /api/v1/roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.roleSvc.Create(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// GetByID 获取角色
// @Summary 获取角色
// @Tags roles
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=Role}
// @Router /api/v1/roles/{id} [get]
func (h *RoleHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	result, err := h.roleSvc.GetByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Update 更新角色
// @Summary 更新角色
// @Tags roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param request body UpdateRoleReq true "角色信息"
// @Success 200 {object} response.Response{data=Role}
// @Router /api/v1/roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.roleSvc.Update(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Delete 删除角色
// @Summary 删除角色
// @Tags roles
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/v1/roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.roleSvc.Delete(c, id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// List 角色列表
// @Summary 角色列表
// @Tags roles
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(20)
// @Success 200 {object} response.Response{data=RoleListResult}
// @Router /api/v1/roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	page, size := 1, 20
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if s, err := strconv.Atoi(c.Query("size")); err == nil && s > 0 && s <= 100 {
		size = s
	}
	result, err := h.roleSvc.List(c, page, size)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ── Routes ────────────────────────────────────────

func RegisterRoleRoutes(v1 *gin.RouterGroup, db *gorm.DB, roleRepo IroleRepository) {
	roleSvc := NewRoleService(roleRepo)
	h := NewRoleHandler(roleSvc)

	roles := v1.Group("/roles")
	roles.Use(middleware.JWTAuth())
	{
		roles.GET("", h.List)
		roles.GET("/:id", h.GetByID)
	}
	admin := v1.Group("/roles")
	admin.Use(middleware.JWTAuth())
	{
		admin.POST("", h.Create)
		admin.PUT("/:id", h.Update)
		admin.DELETE("/:id", h.Delete)
	}
}

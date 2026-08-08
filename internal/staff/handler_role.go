package staff

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type RoleHandler struct {
	svc *RoleService
}

func NewRoleHandler(svc *RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

// List 角色列表(含权限数)
// @Summary 角色列表
// @Tags roles
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]RoleListItem}
// @Router /api/v1/admin/roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	items, err := h.svc.List(c)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, items)
}

// Create 创建角色
// @Summary 创建角色
// @Tags roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreateRoleReq true "角色信息"
// @Success 200 {object} response.Response{data=SysRole}
// @Router /api/v1/admin/roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	role, err := h.svc.Create(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, role)
}

// Update 更新角色(builtin 不可改)
// @Summary 更新角色
// @Tags roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "角色ID"
// @Param request body UpdateRoleReq true "更新信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/roles/{role_id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	roleID, err := utils.ParseIntParam(c, "role_id")
	if err != nil {
		c.Error(err)
		return
	}
	var req UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.Update(c, roleID, &req); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// Delete 删除角色(builtin 不可删)
// @Summary 删除角色
// @Tags roles
// @Security ApiKeyAuth
// @Produce json
// @Param role_id path int true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/roles/{role_id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	roleID, err := utils.ParseIntParam(c, "role_id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.Delete(c, roleID); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// GetPermissions 查询角色已绑权限 ID 列表
// @Summary 查询角色权限
// @Tags roles
// @Security ApiKeyAuth
// @Produce json
// @Param role_id path int true "角色ID"
// @Success 200 {object} response.Response{data=[]int64}
// @Router /api/v1/admin/roles/{role_id}/permissions [get]
func (h *RoleHandler) GetPermissions(c *gin.Context) {
	roleID, err := utils.ParseIntParam(c, "role_id")
	if err != nil {
		c.Error(err)
		return
	}
	ids, err := h.svc.GetPermissionIDs(c, roleID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, ids)
}

// AssignPermissions 全量替换角色权限绑定
// @Summary 分配角色权限
// @Tags roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param role_id path int true "角色ID"
// @Param request body AssignPermissionsReq true "权限ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/roles/{role_id}/permissions [put]
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	roleID, err := utils.ParseIntParam(c, "role_id")
	if err != nil {
		c.Error(err)
		return
	}
	var req AssignPermissionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.AssignPermissions(c, roleID, req.PermissionIDs); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ── Routes ────────────────────────────────────────

func RegisterRoleRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	roleRepo := NewRoleRepository(db)
	staffRoleRepo := NewStaffRoleRepository(db)
	svc := NewRoleService(roleRepo)
	h := NewRoleHandler(svc)

	group := v1.Group("/admin/roles")
	group.Use(middleware.JWTAuth())
	{
		group.GET("", RequirePermission(staffRoleRepo, "role:read"), h.List)
		group.POST("", RequirePermission(staffRoleRepo, "role:create"), h.Create)
		group.PUT("/:role_id", RequirePermission(staffRoleRepo, "role:update"), h.Update)
		group.DELETE("/:role_id", RequirePermission(staffRoleRepo, "role:delete"), h.Delete)
		group.GET("/:role_id/permissions", RequirePermission(staffRoleRepo, "role:read"), h.GetPermissions)
		group.PUT("/:role_id/permissions", RequirePermission(staffRoleRepo, "role:update"), h.AssignPermissions)
	}
}

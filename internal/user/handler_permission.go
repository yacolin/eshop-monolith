package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type PermissionHandler struct {
	permSvc *PermissionService
}

func NewPermissionHandler(permSvc *PermissionService) *PermissionHandler {
	return &PermissionHandler{permSvc: permSvc}
}

// Create 创建权限
// @Summary 创建权限
// @Tags permissions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreatePermissionReq true "权限信息"
// @Success 200 {object} response.Response{data=Permission}
// @Router /api/v1/permissions [post]
func (h *PermissionHandler) Create(c *gin.Context) {
	var req CreatePermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.permSvc.Create(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// GetByID 获取权限
// @Summary 获取权限
// @Tags permissions
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "权限ID"
// @Success 200 {object} response.Response{data=Permission}
// @Router /api/v1/permissions/{id} [get]
func (h *PermissionHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	result, err := h.permSvc.GetByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Update 更新权限
// @Summary 更新权限
// @Tags permissions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Param request body UpdatePermissionReq true "权限信息"
// @Success 200 {object} response.Response{data=Permission}
// @Router /api/v1/permissions/{id} [put]
func (h *PermissionHandler) Update(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req UpdatePermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.permSvc.Update(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Delete 删除权限
// @Summary 删除权限
// @Tags permissions
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "权限ID"
// @Success 200 {object} response.Response
// @Router /api/v1/permissions/{id} [delete]
func (h *PermissionHandler) Delete(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.permSvc.Delete(c, id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// List 权限列表
// @Summary 权限列表
// @Tags permissions
// @Security ApiKeyAuth
// @Produce json
// @Param request query PermissionListReq false "查询参数"
// @Success 200 {object} response.Response{data=PermissionListResult}
// @Router /api/v1/permissions [get]
func (h *PermissionHandler) List(c *gin.Context) {
	var req PermissionListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.permSvc.List(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Check 检查权限
// @Summary 检查权限
// @Tags permissions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CheckPermissionsReq true "权限检查"
// @Success 200 {object} response.Response{data=CheckPermissionsResult}
// @Router /api/v1/permissions/check [post]
func (h *PermissionHandler) Check(c *gin.Context) {
	var req CheckPermissionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.permSvc.CheckUserPermissions(c, currentUserID(c), req.PermissionNames)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, CheckPermissionsResult{Permissions: result})
}

// ── Routes ────────────────────────────────────────

func RegisterPermissionRoutes(v1 *gin.RouterGroup, db *gorm.DB, permRepo IpermissionRepository, roleRepo IroleRepository) {
	permSvc := NewPermissionService(permRepo, roleRepo)
	h := NewPermissionHandler(permSvc)

	perms := v1.Group("/permissions")
	perms.Use(middleware.JWTAuth())
	{
		perms.GET("", h.List)
		perms.GET("/:id", h.GetByID)
		perms.POST("/check", h.Check)
	}
	admin := v1.Group("/permissions")
	admin.Use(middleware.JWTAuth())
	{
		admin.POST("", h.Create)
		admin.PUT("/:id", h.Update)
		admin.DELETE("/:id", h.Delete)
	}
}

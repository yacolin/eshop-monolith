package handlers

import (
	"eshop-monolith/internal/pkg/errcode"
	"eshop-monolith/internal/pkg/response"
	"eshop-monolith/internal/pkg/utils"
	"eshop-monolith/internal/user/api/dto"
	"eshop-monolith/internal/user/service"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService service.PermissionService
	userService       *service.UserService
}

func NewPermissionHandler(
	permissionService service.PermissionService,
	userService *service.UserService,
) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
		userService:       userService,
	}
}

// CreatePermission 创建权限
// @Summary 创建权限
// @Description 创建新权限（需要管理员权限）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param request body dto.CreatePermissionRequest true "权限信息"
// @Success 200 {object} response.Response{data=models.Permission}
// @Router /api/v1/permissions [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req dto.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	permission, err := h.permissionService.CreatePermission(&dto.CreatePermissionRequest{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		Category:    req.Category,
		Sort:        req.Sort,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, permission)
}

// GetPermission 获取权限详情
// @Summary 获取权限详情
// @Description 根据ID获取权限详情
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Success 200 {object} response.Response{data=models.Permission}
// @Router /api/v1/permissions/{id} [get]
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	permission, err := h.permissionService.GetPermission(id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, permission)
}

// UpdatePermission 更新权限
// @Summary 更新权限
// @Description 更新权限信息（需要管理员权限）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Param request body dto.UpdatePermissionRequest true "权限信息"
// @Success 200 {object} response.Response{data=models.Permission}
// @Router /api/v1/permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	permission, err := h.permissionService.UpdatePermission(id, &dto.UpdatePermissionRequest{
		DisplayName: req.DisplayName,
		Description: req.Description,
		Category:    req.Category,
		Sort:        req.Sort,
		Status:      req.Status,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, permission)
}

// DeletePermission 删除权限
// @Summary 删除权限
// @Description 删除权限（需要管理员权限）
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Success 200 {object} response.Response
// @Router /api/v1/permissions/{id} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.permissionService.DeletePermission(id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Permission deleted successfully"})
}

// ListPermissions 获取权限列表
// @Summary 获取权限列表
// @Description 获取权限列表，支持分页和筛选
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Param category query string false "分类"
// @Param resource query string false "资源"
// @Param role query string false "角色"
// @Success 200 {object} response.Response{data=service.ListPermissionsResponse}
// @Router /api/v1/permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	var query dto.PermissionListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}
	query.Normalize()

	var result *dto.ListPermissionsResponse
	var err error

	switch {
	case query.Category != "":
		result, err = h.permissionService.GetPermissionsByCategory(query.Category, query.Page, query.Size)
	case query.Resource != "":
		result, err = h.permissionService.GetPermissionsByResource(query.Resource, query.Page, query.Size)
	case query.RoleID != 0:
		result, err = h.permissionService.GetPermissionsByRoleID(query.RoleID, query.Page, query.Size)
	default:
		result, err = h.permissionService.ListPermissions(query.Page, query.Size)
	}

	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

func (h *PermissionHandler) CheckPermissions(c *gin.Context) {
	var req dto.CheckPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	userID, err := utils.ParseIntParam(c, "user_id")
	if err != nil {
		c.Error(errcode.ErrUnauthorized)
		return
	}

	result, err := h.permissionService.CheckUserPermissions(userID, req.PermissionNames)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, dto.CheckPermissionsResponse{
		Permissions: result,
	})
}

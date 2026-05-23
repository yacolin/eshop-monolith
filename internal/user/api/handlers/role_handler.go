package handlers

import (
	"eshop-monolith/internal/user/api/dto"
	"eshop-monolith/internal/user/service"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	permissionSvc service.PermissionService
}

func NewRoleHandler(permissionSvc service.PermissionService) *RoleHandler {
	return &RoleHandler{
		permissionSvc: permissionSvc,
	}
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色（需要管理员权限）
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateRoleRequest true "角色信息"
// @Success 200 {object} response.Response{data=models.Role} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/roles [post]
// @Example request { "name": "editor", "display_name": "编辑员", "description": "内容编辑员", "status": 1, "sort": 5, "is_system": false }
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	role, err := h.permissionSvc.CreateRole(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, role)
}

// GetRole 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详情，包含角色的权限列表
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID" format(uuid)
// @Success 200 {object} response.Response{data=models.Role} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "角色不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/roles/{id} [get]
// @Example id "550e8400-e29b-41d4-a716-446655440000"
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	role, err := h.permissionSvc.GetRole(id)
	if err != nil {
		c.Error(errcode.ErrNotFound)
		return
	}

	response.Success(c, role)
}

// GetRoleByName 根据名称获取角色
// @Summary 根据名称获取角色
// @Description 根据角色名称获取角色详情
// @Tags roles
// @Accept json
// @Produce json
// @Param name path string true "角色名称"
// @Success 200 {object} response.Response{data=models.Role}
// @Router /api/v1/roles/name/{name} [get]
func (h *RoleHandler) GetRoleByName(c *gin.Context) {
	name := c.Param("name")
	role, err := h.permissionSvc.GetRoleByName(name)
	if err != nil {
		c.Error(errcode.ErrNotFound)
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息（需要管理员权限）。注意：系统内置角色（is_system=true）不能被修改。
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID" format(uuid)
// @Param request body service.UpdateRoleRequest true "角色信息"
// @Success 200 {object} response.Response{data=models.Role} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足或尝试修改系统角色"
// @Failure 404 {object} response.Response "角色不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/roles/{id} [put]
// @Example id "550e8400-e29b-41d4-a716-446655440000"
// @Example request { "display_name": "编辑员（更新）", "description": "更新后的描述" }
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	role, err := h.permissionSvc.UpdateRole(id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, role)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除角色（需要管理员权限）。注意：系统内置角色（is_system=true）不能被删除。
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID" format(uuid)
// @Success 200 {object} response.Response{data=map[string]string} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足或尝试删除系统角色"
// @Failure 404 {object} response.Response "角色不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/roles/{id} [delete]
// @Example id "550e8400-e29b-41d4-a716-446655440000"
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.permissionSvc.DeleteRole(id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Role deleted successfully"})
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Description 分页获取角色列表，支持按页码和每页数量查询
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" minimum(1) default(1)
// @Param page_size query int false "每页数量" minimum(1) maximum(100) default(20)
// @Success 200 {object} response.Response{data=service.ListRolesResponse} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/roles [get]
// @Example page 1
// @Example page_size 20
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.permissionSvc.ListRoles(page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// AssignRoleToUser 分配角色给用户
// @Summary 分配角色给用户
// @Description 为指定用户分配角色（需要管理员权限）
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "用户ID" format(uuid)
// @Param request body object{role_id=string} true "角色信息"
// @Success 200 {object} response.Response{data=map[string]string} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "用户或角色不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/users/{user_id}/roles [post]
// @Example user_id "550e8400-e29b-41d4-a716-446655440000"
// @Example request { "role_id": "660e8400-e29b-41d4-a716-446655440001" }
func (h *RoleHandler) AssignRoleToUser(c *gin.Context) {
	userID, err := utils.ParseIntParam(c, "user_id")
	if err != nil {
		c.Error(err)
		return
	}

	var req struct {
		RoleID int64 `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.permissionSvc.AssignRoleToUser(userID, req.RoleID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Role assigned to user successfully"})
}

// RemoveRoleFromUser 移除用户的角色
// @Summary 移除用户的角色
// @Description 移除指定用户的角色（需要管理员权限）
// @Tags roles
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param role_id path string true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/v1/users/{user_id}/roles/{role_id} [delete]
func (h *RoleHandler) RemoveRoleFromUser(c *gin.Context) {
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

	if err := h.permissionSvc.RemoveRoleFromUser(userID, roleID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Role removed from user successfully"})
}

// GetUserRoles 获取用户的角色列表
// @Summary 获取用户的角色列表
// @Description 获取指定用户的角色列表
// @Tags roles
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Success 200 {object} response.Response{data=[]models.Role}
// @Router /api/v1/users/{user_id}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID, err := utils.ParseIntParam(c, "user_id")
	if err != nil {
		c.Error(err)
		return
	}

	roles, err := h.permissionSvc.GetUserRoles(userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, roles)
}

// AssignPermissionsToRole 为角色分配权限
// @Summary 为角色分配权限
// @Description 为指定角色分配权限（需要管理员权限）
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Param request body object true "权限ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissionsToRole(c *gin.Context) {
	roleID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req struct {
		PermissionIDs []int64 `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.permissionSvc.AssignPermissionsToRole(roleID, req.PermissionIDs); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Permissions assigned to role successfully"})
}

// RemovePermissionsFromRole 移除角色的权限
// @Summary 移除角色的权限
// @Description 移除指定角色的权限（需要管理员权限）
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Param request body object true "权限ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/roles/{id}/permissions [delete]
func (h *RoleHandler) RemovePermissionsFromRole(c *gin.Context) {
	roleID, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req struct {
		PermissionIDs []int64 `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.permissionSvc.RemovePermissionsFromRole(roleID, req.PermissionIDs); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Permissions removed from role successfully"})
}

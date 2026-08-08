package staff

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type StaffHandler struct {
	svc *StaffService
}

func NewStaffHandler(svc *StaffService) *StaffHandler {
	return &StaffHandler{svc: svc}
}

// List 员工列表
// @Summary 员工列表
// @Tags staffs
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(20)
// @Param keyword query string false "关键字"
// @Param status query int false "状态"
// @Success 200 {object} response.Response{data=StaffListResult}
// @Router /api/v1/admin/staff [get]
func (h *StaffHandler) List(c *gin.Context) {
	var req StaffListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.List(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Create 创建员工
// @Summary 创建员工
// @Tags staffs
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreateStaffReq true "员工信息"
// @Success 200 {object} response.Response{data=Staff}
// @Router /api/v1/admin/staff [post]
func (h *StaffHandler) Create(c *gin.Context) {
	var req CreateStaffReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	staff, err := h.svc.Create(c, &req, currentStaffID(c))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, staff)
}

// Update 更新员工(资料/状态/重置密码)
// @Summary 更新员工
// @Tags staffs
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param staff_id path int true "员工ID"
// @Param request body UpdateStaffReq true "更新信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/staff/{staff_id} [put]
func (h *StaffHandler) Update(c *gin.Context) {
	staffID, err := utils.ParseIntParam(c, "staff_id")
	if err != nil {
		c.Error(err)
		return
	}
	var req UpdateStaffReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.Update(c, staffID, &req, currentStaffID(c)); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// Delete 删除员工(软删)
// @Summary 删除员工
// @Tags staffs
// @Security ApiKeyAuth
// @Produce json
// @Param staff_id path int true "员工ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/staff/{staff_id} [delete]
func (h *StaffHandler) Delete(c *gin.Context) {
	staffID, err := utils.ParseIntParam(c, "staff_id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.Delete(c, staffID, currentStaffID(c)); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// AssignRoles 分配角色(全量替换)
// @Summary 分配角色
// @Tags staffs
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param staff_id path int true "员工ID"
// @Param request body AssignRolesReq true "角色ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/staff/{staff_id}/roles [put]
func (h *StaffHandler) AssignRoles(c *gin.Context) {
	staffID, err := utils.ParseIntParam(c, "staff_id")
	if err != nil {
		c.Error(err)
		return
	}
	var req AssignRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.AssignRoles(c, staffID, req.RoleIDs); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ── Routes ────────────────────────────────────────

func RegisterStaffRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	staffRepo := NewStaffRepository(db)
	staffRoleRepo := NewStaffRoleRepository(db)
	svc := NewStaffService(staffRepo, staffRoleRepo)
	h := NewStaffHandler(svc)

	group := v1.Group("/admin/staff")
	group.Use(middleware.JWTAuth())
	{
		group.GET("", RequirePermission(staffRoleRepo, "staff:read"), h.List)
		group.POST("", RequirePermission(staffRoleRepo, "staff:create"), h.Create)
		group.PUT("/:staff_id", RequirePermission(staffRoleRepo, "staff:update"), h.Update)
		group.DELETE("/:staff_id", RequirePermission(staffRoleRepo, "staff:delete"), h.Delete)
		group.PUT("/:staff_id/roles", RequirePermission(staffRoleRepo, "staff:update"), h.AssignRoles)
	}
}

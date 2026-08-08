package staff

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
)

type PermissionHandler struct {
	svc *PermissionService
}

func NewPermissionHandler(svc *PermissionService) *PermissionHandler {
	return &PermissionHandler{svc: svc}
}

// List 权限列表(支持 category 过滤)
// @Summary 权限列表
// @Tags permissions
// @Security ApiKeyAuth
// @Produce json
// @Param category query string false "分类"
// @Success 200 {object} response.Response{data=[]PermissionListItem}
// @Router /api/v1/admin/permissions [get]
func (h *PermissionHandler) List(c *gin.Context) {
	var req PermissionListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	items, err := h.svc.List(c, req.Category)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, items)
}

// ── Routes ────────────────────────────────────────

func RegisterPermissionRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	permRepo := NewPermissionRepository(db)
	staffRoleRepo := NewStaffRoleRepository(db)
	svc := NewPermissionService(permRepo)
	h := NewPermissionHandler(svc)

	group := v1.Group("/admin/permissions")
	group.Use(middleware.JWTAuth(), RequirePermission(staffRoleRepo, "permission:read"))
	{
		group.GET("", h.List)
	}
}

package staff

import (
	"github.com/gin-gonic/gin"

	"eshop-monolith/pkg/errcode"
)

// RequirePermission 校验当前员工是否拥有指定权限(实时查库)。
// 需在 JWTAuth() 之后使用;admin 角色已绑定全量权限,天然放行。
func RequirePermission(staffRoleRepo IstaffRoleRepository, permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		staffID := currentStaffID(c)
		if staffID == 0 {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		has, err := staffRoleRepo.HasPermission(c, staffID, permissionName)
		if err != nil || !has {
			c.Error(errcode.ErrInsufficientPermissions)
			c.Abort()
			return
		}
		c.Next()
	}
}

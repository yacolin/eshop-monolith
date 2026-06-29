package user

import (
	"github.com/gin-gonic/gin"

	"eshop-monolith/pkg/errcode"
)

type RequireRoleConfig struct {
	RoleRepo IroleRepository
}

func NewRequireRoleConfig(roleRepo IroleRepository) *RequireRoleConfig {
	return &RequireRoleConfig{RoleRepo: roleRepo}
}

func RequireRole(config *RequireRoleConfig, roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesClaim, exists := c.Get("roles")
		if !exists {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		userRoles, ok := rolesClaim.([]string)
		if !ok {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		roleNamesFromDB := make([]string, 0, len(roleNames))
		for _, roleName := range roleNames {
			role, err := config.RoleRepo.FindByName(c, roleName)
			if err != nil {
				c.Error(errcode.ErrNotFound)
				c.Abort()
				return
			}
			roleNamesFromDB = append(roleNamesFromDB, role.Name)
		}

		for _, userRole := range userRoles {
			for _, rn := range roleNamesFromDB {
				if userRole == rn {
					c.Next()
					return
				}
			}
		}

		c.Error(errcode.ErrInsufficientPermissions)
		c.Abort()
	}
}

func RequireAdmin(config *RequireRoleConfig) gin.HandlerFunc {
	return RequireRole(config, "admin")
}

func RequireMerchant(config *RequireRoleConfig) gin.HandlerFunc {
	return RequireRole(config, "merchant", "admin")
}

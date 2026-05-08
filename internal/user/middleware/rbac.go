package middleware

import (
	"context"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/user/domain/repositories"

	"github.com/gin-gonic/gin"
)

type PermissionService interface {
	CheckUserPermissions(userID string, permissionNames []string) (map[string]bool, error)
}

type RoleService interface {
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
}

type RequireRoleConfig struct {
	RoleRepo repositories.IroleRepository
}

func NewRequireRoleConfig(roleRepo repositories.IroleRepository) *RequireRoleConfig {
	return &RequireRoleConfig{RoleRepo: roleRepo}
}

func RequirePermission(service PermissionService, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		result, err := service.CheckUserPermissions(userIDStr, permissions)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		for _, permission := range permissions {
			if has, ok := result[permission]; ok && has {
				c.Next()
				return
			}
		}

		c.Error(errcode.ErrInsufficientPermissions)
		c.Abort()
	}
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

		roles := make([]*models.Role, 0, len(roleNames))
		for _, roleName := range roleNames {
			role, err := config.RoleRepo.GetByName(c, roleName)
			if err != nil {
				c.Error(errcode.ErrNotFound)
				c.Abort()
				return
			}
			roles = append(roles, role)
		}

		roleNamesFromDB := make([]string, 0, len(roles))
		for _, role := range roles {
			roleNamesFromDB = append(roleNamesFromDB, role.Name)
		}

		for _, userRole := range userRoles {
			for _, roleName := range roleNamesFromDB {
				if userRole == roleName {
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

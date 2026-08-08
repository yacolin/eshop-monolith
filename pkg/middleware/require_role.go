package middleware

import (
	"github.com/gin-gonic/gin"

	"eshop-monolith/pkg/errcode"
)

// RequireRole 校验 token 中的角色名(纯 claims 判断,不查库)。
// 需在 JWTAuth() 之后使用;角色名单已在签发时写入 token。
func RequireRole(roleNames ...string) gin.HandlerFunc {
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
		for _, want := range roleNames {
			for _, have := range userRoles {
				if have == want {
					c.Next()
					return
				}
			}
		}
		c.Error(errcode.ErrInsufficientPermissions)
		c.Abort()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

func RequireMerchant() gin.HandlerFunc {
	return RequireRole("merchant", "admin")
}

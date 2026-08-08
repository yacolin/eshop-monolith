package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/token"
)

// JWTAuth validates the Authorization header and token, but delegates error
// responses to the centralized ErrorHandler middleware by calling `c.Error(err)`.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := token.ParseToken(parts[1])
		if err != nil {
			c.Error(errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID) // 统一 int64
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

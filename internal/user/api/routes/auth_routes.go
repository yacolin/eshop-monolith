package routes

import (
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/user/api/handlers"
	"eshop-monolith/internal/user/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB) {
	tokenService := service.NewTokenService("your-secret-key-change-in-production", repos.AuthToken, repos.Role)
	authService := service.NewAuthService(db, repos.User, repos.UserIdentity, repos.AuthToken, repos.LoginHistory, repos.Role, tokenService)
	userService := service.NewUserService(repos.User, repos.UserInfo, repos.AuthToken, repos.LoginHistory, nil)
	authHandler := handlers.NewAuthHandler(authService, tokenService, userService)

	auth := v1.Group("/auth")
	{
		auth.POST("/login/password", authHandler.LoginByPassword)
		auth.POST("/login/wechat", authHandler.LoginByWechat)
		auth.POST("/login/phone", authHandler.LoginByPhone)

		auth.POST("/register", authHandler.Register)

		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}
}

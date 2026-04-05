package routes

import (
	"eshop-monolith/internal/pkg/middleware"
	"eshop-monolith/internal/repository"
	"eshop-monolith/internal/user/api/handlers"
	"eshop-monolith/internal/user/service"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(v1 *gin.RouterGroup, repos *repository.Repositories) {
	userService := service.NewUserService(repos.User, repos.UserInfo, repos.AuthToken, repos.LoginHistory, nil)
	userHandler := handlers.NewUserHandler(userService)
	permissionService := service.NewPermissionService(repos.Permission, repos.User, repos.Role)
	roleHandler := handlers.NewRoleHandler(permissionService)

	users := v1.Group("/users")
	{
		protected := users.Group("").Use(middleware.JWTAuth())
		{
			protected.GET("/profile", userHandler.GetProfile)
			protected.GET("/info", userHandler.GetUserInfo)
			protected.PUT("/info", userHandler.UpdateUserInfo)
			protected.GET("/:user_id", userHandler.GetByID)
			protected.GET("/:user_id/roles", roleHandler.GetUserRoles)
		}

		roleConfig := middleware.NewRequireRoleConfig(repos.Role)
		admin := users.Group("").Use(middleware.JWTAuth(), middleware.RequireMerchant(roleConfig))
		{
			admin.POST("/:user_id/roles", roleHandler.AssignRoleToUser)
			admin.DELETE("/:user_id/roles/:role_id", roleHandler.RemoveRoleFromUser)
		}
	}
}

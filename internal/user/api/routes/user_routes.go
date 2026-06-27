package routes

import (
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/user/api/handlers"
	usermw "eshop-monolith/internal/user/middleware"
	"eshop-monolith/internal/user/service"
	"eshop-monolith/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, rabbit *rabbitmq.Client) {
	userService := service.NewUserService(repos.User, repos.UserInfo, repos.AuthToken, repos.LoginHistory, rabbit)
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

		roleConfig := usermw.NewRequireRoleConfig(repos.Role)
		admin := users.Group("").Use(middleware.JWTAuth(), usermw.RequireMerchant(roleConfig))
		{
			admin.POST("/:user_id/roles", roleHandler.AssignRoleToUser)
			admin.DELETE("/:user_id/roles/:role_id", roleHandler.RemoveRoleFromUser)
		}
	}
}

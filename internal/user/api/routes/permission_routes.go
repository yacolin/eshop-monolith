package routes

import (
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/user/api/handlers"
	usermw "eshop-monolith/internal/user/middleware"
	"eshop-monolith/internal/user/service"
	"eshop-monolith/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterPermissionRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB) {
	userService := service.NewUserService(repos.User, repos.UserInfo, repos.AuthToken, repos.LoginHistory, nil)
	permissionService := service.NewPermissionService(repos.Permission, repos.User, repos.Role)
	permissionHandler := handlers.NewPermissionHandler(permissionService, userService)

	permissions := v1.Group("/permissions")
	{
		permissions.Use(middleware.JWTAuth())
		{
			permissions.GET("", permissionHandler.ListPermissions)
			permissions.GET("/:id", permissionHandler.GetPermission)
			permissions.POST("/check", permissionHandler.CheckPermissions)
		}

		roleConfig := usermw.NewRequireRoleConfig(repos.Role)
		admin := permissions.Group("").Use(middleware.JWTAuth(), usermw.RequireAdmin(roleConfig))
		{
			admin.POST("", permissionHandler.CreatePermission)
			admin.PUT("/:id", permissionHandler.UpdatePermission)
			admin.DELETE("/:id", permissionHandler.DeletePermission)
		}
	}
}

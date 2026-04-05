package routes

import (
	"eshop-monolith/internal/pkg/middleware"
	"eshop-monolith/internal/repository"
	"eshop-monolith/internal/user/api/handlers"
	"eshop-monolith/internal/user/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoleRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB) {
	permissionService := service.NewPermissionService(repos.Permission, repos.User, repos.Role)
	roleHandler := handlers.NewRoleHandler(permissionService)

	roles := v1.Group("/roles")
	{
		roles.Use(middleware.JWTAuth())
		{
			roles.GET("", roleHandler.ListRoles)
			roles.GET("/:id", roleHandler.GetRole)
			roles.GET("/name/:name", roleHandler.GetRoleByName)
		}

		roleConfig := middleware.NewRequireRoleConfig(repos.Role)
		admin := roles.Group("").Use(middleware.JWTAuth(), middleware.RequireAdmin(roleConfig))
		{
			admin.POST("", roleHandler.CreateRole)
			admin.PUT("/:id", roleHandler.UpdateRole)
			admin.DELETE("/:id", roleHandler.DeleteRole)

			admin.POST("/:id/permissions", roleHandler.AssignPermissionsToRole)
			admin.DELETE("/:id/permissions", roleHandler.RemovePermissionsFromRole)
		}
	}
}

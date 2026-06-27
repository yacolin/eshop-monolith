package routes

import (
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/notification/api/handlers"
	"eshop-monolith/internal/notification/domain/repositories"
	"eshop-monolith/internal/notification/service"
	usermw "eshop-monolith/internal/user/middleware"
	"eshop-monolith/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterNotificationRoutes 注册通知相关路由
func RegisterNotificationRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, rabbit *rabbitmq.Client) *service.NotificationService {
	notifRepo := repositories.NewNotificationRepository(db)
	notifSvc := service.NewNotificationService(notifRepo)
	notifHandler := handlers.NewNotificationHandler(notifSvc)


	// 通知路由（需要认证）
	notify := v1.Group("/notifications")
	notify.Use(middleware.JWTAuth())
	{
		notify.GET("", notifHandler.ListNotifications)
		notify.GET("/unread", notifHandler.GetUnreadCount)
		notify.PUT("/:id/read", notifHandler.MarkAsRead)
		notify.PUT("/read-all", notifHandler.MarkAllAsRead)
		notify.DELETE("/:id", notifHandler.DeleteNotification)
	}

	// 系统通知发送（需要管理员权限）
	roleCfg := usermw.NewRequireRoleConfig(repos.Role)
	notify.POST("/system", usermw.RequireAdmin(roleCfg), notifHandler.SendSystemNotification)

	return notifSvc
}

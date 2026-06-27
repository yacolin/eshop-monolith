package routes

import (
	"eshop-monolith/internal/dashboard/api/handlers"
	"eshop-monolith/internal/dashboard/service"
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterDashboardRoutes 注册仪表盘相关路由
func RegisterDashboardRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, rabbit *rabbitmq.Client) *service.DashboardService {
	dashboardService := service.NewDashboardService(db, repos.Redis, repos, rabbit)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// 注册事件处理器（数据变更时自动失效缓存）
	dashboardService.RegisterEventHandlers()

	// 仪表盘路由（无需认证）
	dashboard := v1.Group("/dashboard")
	{
		dashboard.GET("/stats", dashboardHandler.GetStats)
	}

	return dashboardService
}

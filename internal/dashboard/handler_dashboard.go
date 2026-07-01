package dashboard

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/pkg/response"
)

type DashboardHandler struct {
	svc *DashboardService
}

func NewDashboardHandler(svc *DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// GetStats 获取仪表盘汇总数据
// @Summary 获取仪表盘汇总数据
// @Tags dashboard
// @Produce json
// @Success 200 {object} response.Response{data=DashboardResponse}
// @Router /api/v1/dashboard/stats [get]
func (h *DashboardHandler) GetStats(c *gin.Context) {
	resp, err := h.svc.GetStats(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, resp)
}

// ── Routes ────────────────────────────────────────

func RegisterDashboardRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, rabbit *rabbitmq.Client) *DashboardService {
	svc := NewDashboardService(db, repos.Redis, repos, rabbit)
	h := NewDashboardHandler(svc)

	dashboard := v1.Group("/dashboard")
	{
		dashboard.GET("/stats", h.GetStats)
	}

	// 后台定时刷新缓存（模块自身管理生命周期）
	svc.StartPeriodicRefresh(context.Background())

	return svc
}

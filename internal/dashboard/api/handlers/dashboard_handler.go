package handlers

import (
	"eshop-monolith/internal/dashboard/service"
	"eshop-monolith/pkg/response"

	"github.com/gin-gonic/gin"
)

// DashboardHandler 仪表盘处理器
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler 创建仪表盘处理器
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetStats 获取仪表盘汇总数据
// @Summary 获取仪表盘汇总数据
// @Description 获取仪表盘全部汇总数据，包括核心指标、订单趋势、各状态分布、热销商品等
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=dto.DashboardResponse}
// @Router /api/v1/dashboard/stats [get]
func (h *DashboardHandler) GetStats(c *gin.Context) {
	resp, err := h.dashboardService.GetStats(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, resp)
}

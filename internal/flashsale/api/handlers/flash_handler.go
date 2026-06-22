package handlers

import (
	"eshop-monolith/internal/flashsale/api/dto"
	"eshop-monolith/internal/flashsale/service"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

type FlashHandler struct {
	flashService *service.FlashService
}

func NewFlashHandler(flashService *service.FlashService) *FlashHandler {
	return &FlashHandler{flashService: flashService}
}

func (h *FlashHandler) CreateActivity(c *gin.Context) {
	var req dto.CreateActivityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	activity, err := h.flashService.CreateActivity(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, activity)
}

func (h *FlashHandler) LoadStock(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.flashService.LoadStockToRedis(c, id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "stock loaded to redis successfully"})
}

func (h *FlashHandler) FlashBuy(c *gin.Context) {
	var req dto.FlashBuyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.flashService.FlashBuy(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, resp)
}

func (h *FlashHandler) GetActivity(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	activity, err := h.flashService.GetActivity(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, activity)
}

func (h *FlashHandler) ListActivities(c *gin.Context) {
	activities, err := h.flashService.ListActivities(c)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, activities)
}

// ListActivitiesByCursor 基于游标分页查询活动列表
// @Summary 基于游标分页查询活动列表
// @Description 使用游标分页代替传统 OFFSET 分页，解决深分页性能问题
// @Tags flash_activities
// @Accept json
// @Produce json
// @Param cursor query int false "游标（上一页最后一条的 ID，首次查询传 0）" default(0)
// @Param size query int false "每页条数" default(20)
// @Param status query string false "筛选状态：pending/active/finished"
// @Success 200 {object} response.Response{data=dto.ActivityCursorResult}
// @Router /api/v1/flash/activities/cursor [get]
func (h *FlashHandler) ListActivitiesByCursor(c *gin.Context) {
	var q dto.ActivityCursorQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}
	q.Normalize()
	result, err := h.flashService.ListActivitiesByCursor(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *FlashHandler) GetOrder(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	order, err := h.flashService.GetOrder(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, order)
}

func (h *FlashHandler) GetUserOrders(c *gin.Context) {
	userID, err := utils.ParseIntParam(c, "user_id")
	if err != nil {
		c.Error(err)
		return
	}

	activityID, _ := utils.ParseIntParam(c, "activity_id")
	orders, err := h.flashService.GetUserOrders(c, userID, activityID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, orders)
}

func (h *FlashHandler) ConfirmOrder(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.flashService.ConfirmOrder(c, id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "order confirmed successfully"})
}

func (h *FlashHandler) CancelOrder(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.flashService.CancelOrder(c, id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "order cancelled successfully"})
}
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
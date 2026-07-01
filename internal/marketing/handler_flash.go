package marketing

import (
	"eshop-monolith/pkg/response"

	"github.com/gin-gonic/gin"
)

type FlashHandler struct {
	svc *FlashService
}

func NewFlashHandler(svc *FlashService) *FlashHandler {
	return &FlashHandler{svc: svc}
}

// Buy 秒杀抢购
// @Summary 秒杀抢购
// @Tags flash
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body FlashBuyReq true "抢购信息"
// @Success 200 {object} response.Response
// @Router /api/v1/flash/buy [post]
func (h *FlashHandler) Buy(c *gin.Context) {

	var req FlashBuyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Buy(c, userID(c), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Confirm 确认秒杀订单
// @Summary 确认秒杀订单
// @Tags flash
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body FlashConfirmReq true "确认信息"
// @Success 200 {object} response.Response
// @Router /api/v1/flash/confirm [post]
func (h *FlashHandler) Confirm(c *gin.Context) {

	var req FlashConfirmReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Confirm(c, userID(c), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

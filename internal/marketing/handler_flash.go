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

func (h *FlashHandler) Buy(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req FlashBuyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Buy(c, userID.(int64), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *FlashHandler) Confirm(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req FlashConfirmReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Confirm(c, userID.(int64), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

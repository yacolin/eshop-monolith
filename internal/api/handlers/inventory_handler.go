package handlers

import (
	"eshop-monolith/internal/domain/inventory"
	"eshop-monolith/internal/pkg/response"
	"eshop-monolith/internal/service"

	"github.com/gin-gonic/gin"
)

// InventoryHandler 库存处理器
type InventoryHandler struct {
	inventoryService *service.InventoryService
}

// NewInventoryHandler 创建库存处理器
func NewInventoryHandler(inventoryService *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// ListInventories 列出所有库存
// @Summary 列出所有库存
// @Description 获取所有库存的列表
// @Tags 库存管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=inventory.InventoryListResult}
// @Failure 500 {object} response.Response{error=string}
// @Router /api/inventories [get]
func (h *InventoryHandler) ListInventories(c *gin.Context) {
	var q inventory.InventoryListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	// normalize pagination values (ensure page>=1, 1<=size<=100)
	(&q).Normalize()

	result, err := h.inventoryService.ListInventories(c, q)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

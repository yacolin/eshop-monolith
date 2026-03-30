package handlers

import (
	"eshop-monolith/internal/domain/inventory"
	"eshop-monolith/internal/pkg/response"
	"eshop-monolith/internal/pkg/utils"
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

// CreateInventory 创建库存
// @Summary 创建库存
// @Description 创建一个新的库存记录
// @Tags 库存
// @Accept json
// @Produce json
// @Param inventory body dto.CreateInventoryDTO true "库存信息"
// @Success 200 {object} models.Inventory "成功"
// @Router /inventory/api/v1/inventories [post]
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req inventory.CreateInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	inventory, err := h.inventoryService.CreateInventory(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, inventory)
}

// UpdateInventory 更新库存
// @Summary 更新库存
// @Description 根据ID更新库存信息
// @Tags 库存
// @Accept json
// @Produce json
// @Param id path string true "库存ID"
// @Param inventory body dto.UpdateInventoryDTO true "库存信息"
// @Success 200 {object} models.Inventory "成功"
// @Router /inventory/api/v1/inventories/{id} [put]
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	var req inventory.UpdateInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	inventory, err := h.inventoryService.UpdateInventory(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, inventory)
}

// GetInventoryByProductID 根据产品ID获取库存
// @Summary 根据产品ID获取库存
// @Description 根据产品ID获取库存信息
// @Tags 库存
// @Produce json
// @Param productId path string true "产品ID"
// @Success 200 {object} models.Inventory "成功"
// @Router /inventory/api/v1/inventories/product/{productId} [get]
func (h *InventoryHandler) GetInventoryByProductID(c *gin.Context) {
	productId, err := utils.ParseIntParam(c, "productId")
	if err != nil {
		c.Error(err)
		return
	}

	inventory, err := h.inventoryService.GetInventoryByProductID(c, productId)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, inventory)
}

// ReserveInventory 预订库存
// @Summary 预订库存
// @Description 预订指定产品的库存
// @Tags 库存操作
// @Accept json
// @Produce json
// @Param reserve body inventory.ReserveInventoryDTO true "预订信息"
// @Success 200 {object} map[string]string "成功"
// @Router /inventory/api/v1/inventories/reserve [post]
func (h *InventoryHandler) ReserveInventory(c *gin.Context) {
	var req inventory.ReserveInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.inventoryService.ReserveInventory(c, &req); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "reserved"})
}

// ReleaseInventory 释放库存
// @Summary 释放库存
// @Description 释放之前预订的库存
// @Tags 库存操作
// @Accept json
// @Produce json
// @Param release body inventory.ReleaseInventoryDTO true "释放信息"
// @Success 200 {object} map[string]string "成功"
// @Router /inventory/api/v1/inventories/release [post]
func (h *InventoryHandler) ReleaseInventory(c *gin.Context) {
	var req inventory.ReleaseInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.inventoryService.ReleaseInventory(c, &req); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "released"})
}

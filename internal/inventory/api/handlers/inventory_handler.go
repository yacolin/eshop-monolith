package handlers

import (
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

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
// @Description 获取所有库存的列表，支持分页筛选
// @Tags inventories
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Param product_id query int false "产品ID精确搜索"
// @Param product_name query string false "产品名称模糊搜索"
// @Param sku_code query string false "SKU编码精确搜索"
// @Param sku_name query string false "SKU名称模糊搜索"
// @Param status query string false "库存状态(instock/lowstock/outofstock)"
// @Param low_stock query bool false "是否低库存"
// @Success 200 {object} response.Response{data=dto.InventoryListResult}
// @Router /api/v1/inventories [get]
func (h *InventoryHandler) ListInventories(c *gin.Context) {
	var q dto.InventoryListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	(&q).Normalize()

	result, err := h.inventoryService.ListInventories(c, q)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// ListInventoriesEnriched 列出所有库存（含 SKU 名称）
// @Summary 列出所有库存（含 SKU 名称）
// @Description 获取所有库存的列表，每条记录附带 SKU 名称，用于列表展示增强
// @Tags inventories
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Param sku_name query string false "SKU名称模糊搜索"
// @Param sku_code query string false "SKU编码精确搜索"
// @Param product_id query int false "产品ID精确搜索"
// @Param product_name query string false "产品名称模糊搜索"
// @Success 200 {object} response.Response{data=dto.InventoryEnrichedResult}
// @Router /api/v1/inventories/enriched [get]
func (h *InventoryHandler) ListInventoriesEnriched(c *gin.Context) {
	var q dto.InventoryListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	(&q).Normalize()

	result, err := h.inventoryService.ListInventoriesEnriched(c, q)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// CreateInventory 创建库存
// @Summary 创建库存
// @Description 创建一个新的库存记录
// @Tags inventories
// @Accept json
// @Produce json
// @Param inventory body dto.CreateInventoryDTO true "库存信息"
// @Success 200 {object} response.Response{data=models.Inventory}
// @Router /api/v1/inventories [post]
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req dto.CreateInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	dto, err := h.inventoryService.CreateInventory(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, dto)
}

// BatchCreateInventory 批量创建库存
// @Summary 批量创建库存
// @Description 批量创建多个SKU的库存记录
// @Tags inventories
// @Accept json
// @Produce json
// @Param inventory body dto.BatchCreateInventoryDTO true "批量库存信息"
// @Success 200 {object} response.Response{data=[]models.Inventory}
// @Router /api/v1/inventories/batch [post]
func (h *InventoryHandler) BatchCreateInventory(c *gin.Context) {
	var req dto.BatchCreateInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.inventoryService.BatchCreateInventory(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// UpdateInventory 更新库存
// @Summary 更新库存
// @Description 根据ID更新库存信息
// @Tags inventories
// @Accept json
// @Produce json
// @Param id path int true "库存ID"
// @Param inventory body dto.UpdateInventoryDTO true "库存信息"
// @Success 200 {object} response.Response{data=models.Inventory}
// @Router /api/v1/inventories/{id} [put]
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.UpdateInventoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	dto, err := h.inventoryService.UpdateInventory(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, dto)
}

// GetInventoryBySkuID 根据SKU ID获取库存
// @Summary 根据SKU ID获取库存
// @Description 根据SKU ID获取库存信息
// @Tags inventories
// @Produce json
// @Param sku_id path int true "SKU ID"
// @Success 200 {object} response.Response{data=models.Inventory}
// @Router /api/v1/inventories/sku/{sku_id} [get]
func (h *InventoryHandler) GetInventoryBySkuID(c *gin.Context) {
	skuId, err := utils.ParseIntParam(c, "sku_id")
	if err != nil {
		c.Error(err)
		return
	}

	dto, err := h.inventoryService.GetInventoryBySkuID(c, skuId)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, dto)
}

// ReserveInventory 预订库存
// @Summary 预订库存
// @Description 预订指定产品的库存
// @Tags inventories
// @Accept json
// @Produce json
// @Param reserve body dto.ReserveInventoryDTO true "预订信息"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/inventories/reserve [post]
func (h *InventoryHandler) ReserveInventory(c *gin.Context) {
	var req dto.ReserveInventoryDTO
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
// @Tags inventories
// @Accept json
// @Produce json
// @Param release body dto.ReleaseInventoryDTO true "释放信息"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/inventories/release [post]
func (h *InventoryHandler) ReleaseInventory(c *gin.Context) {
	var req dto.ReleaseInventoryDTO
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

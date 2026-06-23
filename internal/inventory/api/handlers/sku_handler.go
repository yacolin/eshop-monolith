package handlers

import (
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SkuHandler struct {
	skuService *service.SkuService
}

func NewSkuHandler(skuService *service.SkuService) *SkuHandler {
	return &SkuHandler{skuService: skuService}
}

// CreateSku 创建SKU
// @Summary 创建SKU
// @Description 创建一个新的SKU
// @Tags skus
// @Accept json
// @Produce json
// @Param sku body dto.CreateSkuDTO true "SKU信息"
// @Success 200 {object} response.Response{data=dto.SkuResponse}
// @Router /api/v1/skus [post]
func (h *SkuHandler) CreateSku(c *gin.Context) {
	var req dto.CreateSkuDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	sku, err := h.skuService.CreateSku(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, dto.SkuToResponse(sku))
}

// GetSku 获取SKU详情
// @Summary 获取SKU详情
// @Description 根据ID获取SKU详情
// @Tags skus
// @Accept json
// @Produce json
// @Param id path int true "SKU ID"
// @Success 200 {object} response.Response{data=dto.SkuResponse}
// @Router /api/v1/skus/{id} [get]
func (h *SkuHandler) GetSku(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	sku, err := h.skuService.GetSku(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, dto.SkuToResponse(sku))
}

// ListSkus SKU列表查询
// @Summary SKU列表查询
// @Description 分页查询SKU列表，支持按产品ID/名称/SKU编码/价格范围筛选
// @Tags skus
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Param product_id query int false "产品ID（可选）"
// @Param name query string false "SKU名称模糊搜索"
// @Param sku_code query string false "SKU编码精确搜索"
// @Param price_min query int false "最低价格（分）"
// @Param price_max query int false "最高价格（分）"
// @Param sort_by query string false "排序字段 (id, name, price, created_at)"
// @Param order query string false "排序方向 (asc, desc)" default(asc)
// @Success 200 {object} response.Response{data=dto.SkuListResult}
// @Router /api/v1/skus [get]
func (h *SkuHandler) ListSkus(c *gin.Context) {
	var q dto.SkuListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}
	(&q).Normalize()

	result, err := h.skuService.ListSkus(c, q)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// UpdateSku 更新SKU
// @Summary 更新SKU
// @Description 根据ID更新SKU信息
// @Tags skus
// @Accept json
// @Produce json
// @Param id path int true "SKU ID"
// @Param sku body dto.UpdateSkuDTO true "SKU信息"
// @Success 200 {object} response.Response{data=dto.SkuResponse}
// @Router /api/v1/skus/{id} [put]
func (h *SkuHandler) UpdateSku(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.UpdateSkuDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	sku, err := h.skuService.UpdateSku(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, dto.SkuToResponse(sku))
}

// DeleteSku 删除SKU
// @Summary 删除SKU
// @Description 根据ID删除SKU
// @Tags skus
// @Produce json
// @Param id path int true "SKU ID"
// @Success 200 {object} response.Response
// @Router /api/v1/skus/{id} [delete]
func (h *SkuHandler) DeleteSku(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.skuService.DeleteSku(c, id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

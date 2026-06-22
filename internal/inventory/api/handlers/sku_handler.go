package handlers

import (
	"strconv"

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

// ListByProductID 根据产品ID获取SKU列表
// @Summary 根据产品ID获取SKU列表
// @Description 根据产品ID获取该产品下所有SKU
// @Tags skus
// @Accept json
// @Produce json
// @Param product_id query int true "产品ID"
// @Success 200 {object} response.Response{data=dto.SkuListResult}
// @Router /api/v1/skus [get]
func (h *SkuHandler) ListByProductID(c *gin.Context) {
	productIDStr := c.Query("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	skus, err := h.skuService.ListByProductID(c, productID)
	if err != nil {
		c.Error(err)
		return
	}
	list := make([]dto.SkuResponse, len(skus))
	for i, sku := range skus {
		list[i] = dto.SkuToResponse(&sku)
	}
	response.Success(c, dto.SkuListResult{
		List:  list,
		Total: int64(len(skus)),
	})
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

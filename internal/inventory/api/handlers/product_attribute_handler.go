package handlers

import (
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ProductAttributeHandler struct {
	svc *service.ProductAttributeService
}

func NewProductAttributeHandler(svc *service.ProductAttributeService) *ProductAttributeHandler {
	return &ProductAttributeHandler{svc: svc}
}

// GetProductAttributes 获取产品属性定义
// @Summary 获取产品属性定义
// @Description 获取产品的所有规格属性维度及其可选值
// @Tags products
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response{data=[]dto.ProductAttributeItem}
// @Router /api/v1/products/{id}/attributes [get]
func (h *ProductAttributeHandler) GetProductAttributes(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	attrs, err := h.svc.GetProductAttributes(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, attrs)
}

// UpdateProductAttributes 更新产品属性关联（全量替换）
// @Summary 更新产品属性关联
// @Description 全量替换产品的规格属性值关联，原有关联会被清除
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Param attributes body dto.UpdateProductAttributesDTO true "属性值关联列表"
// @Success 200 {object} response.Response
// @Router /api/v1/products/{id}/attributes [put]
func (h *ProductAttributeHandler) UpdateProductAttributes(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateProductAttributesDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.svc.UpdateProductAttributes(c, id, &req); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// BatchCreateSkus 批量创建 SKU
// @Summary 批量创建 SKU
// @Description 批量创建 SKU，自动生成属性组合关联
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Param skus body dto.BatchCreateSkuDTO true "批量SKU数据"
// @Success 200 {object} response.Response{data=dto.BatchCreateSkuResult}
// @Router /api/v1/products/{id}/skus/batch [post]
func (h *ProductAttributeHandler) BatchCreateSkus(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.BatchCreateSkuDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	result, err := h.svc.BatchCreateSkus(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

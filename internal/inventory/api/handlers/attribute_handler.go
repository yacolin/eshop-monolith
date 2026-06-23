package handlers

import (
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AttributeHandler struct {
	svc *service.AttributeService
}

func NewAttributeHandler(svc *service.AttributeService) *AttributeHandler {
	return &AttributeHandler{svc: svc}
}

// ── Attribute ──────────────────────────────────────────────────

// ListAttributes 获取属性维度列表
// @Summary 获取属性维度列表
// @Description 分页查询规格属性维度列表
// @Tags attributes
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=dto.AttributeListResult}
// @Router /api/v1/attributes [get]
func (h *AttributeHandler) ListAttributes(c *gin.Context) {
	var q dto.AttributeListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}
	(&q).Normalize()

	result, err := h.svc.ListAttributes(c, q)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// GetAttribute 获取属性维度详情
// @Summary 获取属性维度详情
// @Description 根据ID获取规格属性维度详情
// @Tags attributes
// @Produce json
// @Param id path int true "属性ID"
// @Success 200 {object} response.Response{data=dto.AttributeResponse}
// @Router /api/v1/attributes/{id} [get]
func (h *AttributeHandler) GetAttribute(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	attr, err := h.svc.GetAttribute(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, dto.AttributeToResponse(attr))
}

// CreateAttribute 创建属性维度
// @Summary 创建属性维度
// @Description 创建一个新的规格属性维度，如"颜色"、"尺寸"
// @Tags attributes
// @Accept json
// @Produce json
// @Param attr body dto.CreateAttributeDTO true "属性信息"
// @Success 200 {object} response.Response{data=dto.AttributeResponse}
// @Router /api/v1/attributes [post]
func (h *AttributeHandler) CreateAttribute(c *gin.Context) {
	var req dto.CreateAttributeDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	attr, err := h.svc.CreateAttribute(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, dto.AttributeToResponse(attr))
}

// UpdateAttribute 更新属性维度
// @Summary 更新属性维度
// @Description 根据ID更新规格属性维度信息
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "属性ID"
// @Param attr body dto.UpdateAttributeDTO true "属性信息"
// @Success 200 {object} response.Response{data=dto.AttributeResponse}
// @Router /api/v1/attributes/{id} [put]
func (h *AttributeHandler) UpdateAttribute(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.UpdateAttributeDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	attr, err := h.svc.UpdateAttribute(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, dto.AttributeToResponse(attr))
}

// DeleteAttribute 删除属性维度
// @Summary 删除属性维度
// @Description 根据ID删除规格属性维度
// @Tags attributes
// @Produce json
// @Param id path int true "属性ID"
// @Success 200 {object} response.Response
// @Router /api/v1/attributes/{id} [delete]
func (h *AttributeHandler) DeleteAttribute(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.DeleteAttribute(c, id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

// ── AttributeValue ─────────────────────────────────────────────

// ListAttributeValues 获取属性值列表
// @Summary 获取属性值列表
// @Description 获取指定属性维度的所有可选值
// @Tags attribute-values
// @Produce json
// @Param id path int true "属性ID"
// @Success 200 {object} response.Response{data=[]dto.AttributeValueResponse}
// @Router /api/v1/attributes/{id}/values [get]
func (h *AttributeHandler) ListAttributeValues(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	vals, err := h.svc.ListAttributeValues(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, vals)
}

// CreateAttributeValue 创建属性值
// @Summary 创建属性值
// @Description 为指定属性维度创建一个可选值
// @Tags attribute-values
// @Accept json
// @Produce json
// @Param val body dto.CreateAttributeValueDTO true "属性值信息"
// @Success 200 {object} response.Response{data=dto.AttributeValueResponse}
// @Router /api/v1/attribute-values [post]
func (h *AttributeHandler) CreateAttributeValue(c *gin.Context) {
	var req dto.CreateAttributeValueDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	val, err := h.svc.CreateAttributeValue(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, dto.AttributeValueToResponse(val))
}

// UpdateAttributeValue 更新属性值
// @Summary 更新属性值
// @Description 根据ID更新属性可选值信息
// @Tags attribute-values
// @Accept json
// @Produce json
// @Param id path int true "属性值ID"
// @Param val body dto.UpdateAttributeValueDTO true "属性值信息"
// @Success 200 {object} response.Response{data=dto.AttributeValueResponse}
// @Router /api/v1/attribute-values/{id} [put]
func (h *AttributeHandler) UpdateAttributeValue(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.UpdateAttributeValueDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	val, err := h.svc.UpdateAttributeValue(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, dto.AttributeValueToResponse(val))
}

// DeleteAttributeValue 删除属性值
// @Summary 删除属性值
// @Description 根据ID删除属性可选值
// @Tags attribute-values
// @Produce json
// @Param id path int true "属性值ID"
// @Success 200 {object} response.Response
// @Router /api/v1/attribute-values/{id} [delete]
func (h *AttributeHandler) DeleteAttributeValue(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.DeleteAttributeValue(c, id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

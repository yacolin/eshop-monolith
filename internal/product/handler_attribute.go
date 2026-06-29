package product

import (
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttributeHandler struct {
	svc *AttributeService
}

func NewAttributeHandler(svc *AttributeService) *AttributeHandler {
	return &AttributeHandler{svc: svc}
}

// Create 创建属性
// @Summary 创建属性
// @Tags attributes
// @Accept json
// @Produce json
// @Param attr body CreateAttributeReq true "属性信息"
// @Success 200 {object} response.Response{data=Attribute}
// @Router /api/v1/attributes [post]
func (h *AttributeHandler) Create(c *gin.Context) {
	var req CreateAttributeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Create(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// GetByID 获取属性详情
// @Summary 获取属性详情
// @Tags attributes
// @Produce json
// @Param id path int true "属性ID"
// @Success 200 {object} response.Response{data=Attribute}
// @Router /api/v1/attributes/{id} [get]
func (h *AttributeHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.GetByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ListByCategory 按类目查询属性
// @Summary 按类目查询属性
// @Tags attributes
// @Produce json
// @Param category_id query int true "类目ID"
// @Success 200 {object} response.Response{data=[]Attribute}
// @Router /api/v1/attributes [get]
func (h *AttributeHandler) ListByCategory(c *gin.Context) {
	categoryID, err := utils.ParseIntParam(c, "category_id")
	if err != nil {
		// 没有 category_id 时查全部
		result, err := h.svc.ListAll(c)
		if err != nil {
			c.Error(err)
			return
		}
		response.Success(c, result)
		return
	}
	result, err := h.svc.ListByCategory(c, categoryID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ListSearchable 按类目查询可筛选项
// @Summary 按类目查询可筛选项
// @Tags attributes
// @Produce json
// @Param category_id query int true "类目ID"
// @Success 200 {object} response.Response{data=[]Attribute}
// @Router /api/v1/attributes/searchable [get]
func (h *AttributeHandler) ListSearchable(c *gin.Context) {
	categoryID, err := utils.ParseIntParam(c, "category_id")
	if err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.ListSearchable(c, categoryID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ListSkuSpec 按类目查询 SKU 规格
// @Summary 按类目查询 SKU 规格
// @Tags attributes
// @Produce json
// @Param category_id query int true "类目ID"
// @Success 200 {object} response.Response{data=[]Attribute}
// @Router /api/v1/attributes/sku-spec [get]
func (h *AttributeHandler) ListSkuSpec(c *gin.Context) {
	categoryID, err := utils.ParseIntParam(c, "category_id")
	if err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.ListSkuSpec(c, categoryID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Update 更新属性
// @Summary 更新属性
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "属性ID"
// @Param attr body UpdateAttributeReq true "属性信息"
// @Success 200 {object} response.Response{data=Attribute}
// @Router /api/v1/attributes/{id} [put]
func (h *AttributeHandler) Update(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req UpdateAttributeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Update(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Delete 删除属性
// @Summary 删除属性
// @Tags attributes
// @Produce json
// @Param id path int true "属性ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/attributes/{id} [delete]
func (h *AttributeHandler) Delete(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.Delete(c, id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

// ── Routes ────────────────────────────────────────

func RegisterAttributeRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	attrRepo := NewAttributeRepository(db)
	catRepo := NewCategoryRepository(db)
	svc := NewAttributeService(attrRepo, catRepo)
	h := NewAttributeHandler(svc)

	attrs := v1.Group("/attributes")
	{
		attrs.GET("", h.ListByCategory)
		attrs.GET("/searchable", h.ListSearchable)
		attrs.GET("/sku-spec", h.ListSkuSpec)
		attrs.GET("/:id", h.GetByID)
	}
	auth := v1.Group("/attributes")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", h.Create)
		auth.PUT("/:id", h.Update)
		auth.DELETE("/:id", h.Delete)
	}
}

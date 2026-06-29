package product

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type SKUHandler struct {
	svc *SKUService
}

func NewSKUHandler(svc *SKUService) *SKUHandler {
	return &SKUHandler{svc: svc}
}

// Create 创建 SKU
// @Summary 创建 SKU
// @Tags skus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreateSKUReq true "SKU 信息"
// @Success 200 {object} response.Response{data=SKU}
// @Router /api/v1/skus [post]
func (h *SKUHandler) Create(c *gin.Context) {
	var req CreateSKUReq
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

// GetByID 获取 SKU 详情
// @Summary 获取 SKU 详情
// @Tags skus
// @Produce json
// @Param id path int true "SKU ID"
// @Success 200 {object} response.Response{data=SKU}
// @Router /api/v1/skus/{id} [get]
func (h *SKUHandler) GetByID(c *gin.Context) {
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

// GetByCode 根据编码获取 SKU
// @Summary 根据编码获取 SKU
// @Tags skus
// @Produce json
// @Param code path string true "SKU 编码"
// @Success 200 {object} response.Response{data=SKU}
// @Router /api/v1/skus/code/:code [get]
func (h *SKUHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}
	result, err := h.svc.GetByCode(c, code)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Update 更新 SKU
// @Summary 更新 SKU
// @Tags skus
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "SKU ID"
// @Param request body UpdateSKUReq true "SKU 信息"
// @Success 200 {object} response.Response{data=SKU}
// @Router /api/v1/skus/{id} [put]
func (h *SKUHandler) Update(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req UpdateSKUReq
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

// Delete 删除 SKU
// @Summary 删除 SKU
// @Tags skus
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "SKU ID"
// @Success 200 {object} response.Response
// @Router /api/v1/skus/{id} [delete]
func (h *SKUHandler) Delete(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.Delete(c, id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ListByProduct 获取产品的 SKU 列表
// @Summary 获取产品的 SKU 列表
// @Tags skus
// @Produce json
// @Param product_id query int true "产品 ID"
// @Success 200 {object} response.Response{data=[]SKU}
// @Router /api/v1/skus [get]
func (h *SKUHandler) ListByProduct(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Query("product_id"), 10, 64)
	if err != nil || productID <= 0 {
		c.Error(errcode.ErrInvalidParams)
		return
	}
	list, err := h.svc.ListByProduct(c, productID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, list)
}

// ── Routes ────────────────────────────────────────

func RegisterSKURoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	repo := NewSpuRepository(db)
	svc := NewSKUService(repo)
	h := NewSKUHandler(svc)

	skus := v1.Group("/skus")
	{
		skus.GET("", h.ListByProduct)
		skus.GET("/code/:code", h.GetByCode)
		skus.GET("/:id", h.GetByID)
	}
	auth := v1.Group("/skus")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", h.Create)
		auth.PUT("/:id", h.Update)
		auth.DELETE("/:id", h.Delete)
	}
}

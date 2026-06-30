package product

import (

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

)

type CategoryBrandHandler struct {
	svc *CategoryBrandService
}

func NewCategoryBrandHandler(svc *CategoryBrandService) *CategoryBrandHandler {
	return &CategoryBrandHandler{svc: svc}
}

// SetBrands 设置类目关联的品牌
// @Summary 设置类目关联的品牌
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "类目ID"
// @Param req body SetCategoryBrandsReq true "品牌ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/categories/{id}/brands [put]
func (h *CategoryBrandHandler) SetBrands(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req SetCategoryBrandsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.SetCategoryBrands(c, id, &req); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ListBrands 查类目下的品牌列表（含品牌详情）
// @Summary 查类目下的品牌列表
// @Tags categories
// @Produce json
// @Param id path int true "类目ID"
// @Success 200 {object} response.Response{data=[]CategoryBrandDetail}
// @Router /api/v1/categories/{id}/brands [get]
func (h *CategoryBrandHandler) ListBrands(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.ListCategoryBrandDetails(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ── Routes ────────────────────────────────────────

func RegisterCategoryBrandRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	repo := NewCategoryBrandRepository(db)
	catRepo := NewCategoryRepository(db)
	brandRepo := NewBrandRepository(db)
	svc := NewCategoryBrandService(repo, catRepo, brandRepo, db)
	h := NewCategoryBrandHandler(svc)

	auth := v1.Group("/categories")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/:id/brands", h.ListBrands)
		auth.PUT("/:id/brands", h.SetBrands)
	}
}

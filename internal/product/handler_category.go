package product

import (

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

)

type CategoryHandler struct {
	svc *CategoryService
}

func NewCategoryHandler(svc *CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// Create 创建类目
// @Summary 创建类目
// @Tags categories
// @Accept json
// @Produce json
// @Param category body CreateCategoryReq true "类目信息"
// @Success 200 {object} response.Response{data=Category}
// @Router /api/v1/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryReq
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

// GetByID 获取类目详情
// @Summary 获取类目详情
// @Tags categories
// @Produce json
// @Param id path int true "类目ID"
// @Success 200 {object} response.Response{data=Category}
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
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

// ListRoot 根类目列表
// @Summary 根类目列表
// @Tags categories
// @Produce json
// @Success 200 {object} response.Response{data=[]Category}
// @Router /api/v1/categories/root [get]
func (h *CategoryHandler) ListRoot(c *gin.Context) {
	result, err := h.svc.ListRoot(c)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ListChildren 子类目列表
// @Summary 子类目列表
// @Tags categories
// @Produce json
// @Param id path int true "父类目ID"
// @Success 200 {object} response.Response{data=[]Category}
// @Router /api/v1/categories/{id}/children [get]
func (h *CategoryHandler) ListChildren(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.ListChildren(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ListAll 全部分类列表
// @Summary 全部分类
// @Tags categories
// @Produce json
// @Success 200 {object} response.Response{data=[]Category}
// @Router /api/v1/categories/all [get]
func (h *CategoryHandler) ListAll(c *gin.Context) {
	result, err := h.svc.ListAll(c)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ListByLevel 按层级查询
// @Summary 按层级查询
// @Tags categories
// @Produce json
// @Param level path int true "层级（1-3）"
// @Success 200 {object} response.Response{data=[]Category}
// @Router /api/v1/categories/level/{level} [get]
func (h *CategoryHandler) ListByLevel(c *gin.Context) {
	level, err := utils.ParseIntParam(c, "level")
	if err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.ListByLevel(c, int8(level))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Update 更新类目
// @Summary 更新类目
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "类目ID"
// @Param category body UpdateCategoryReq true "类目信息"
// @Success 200 {object} response.Response{data=Category}
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req UpdateCategoryReq
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

// Delete 删除类目
// @Summary 删除类目
// @Tags categories
// @Produce json
// @Param id path int true "类目ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
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

func RegisterCategoryRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	repo := NewCategoryRepository(db)
	svc := NewCategoryService(repo)
	h := NewCategoryHandler(svc)

	cats := v1.Group("/categories")
	{
		cats.GET("/root", h.ListRoot)
		cats.GET("/all", h.ListAll)
		cats.GET("/level/:level", h.ListByLevel)
		cats.GET("/:id/children", h.ListChildren)
		cats.GET("/:id", h.GetByID)
	}
	auth := v1.Group("/categories")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", h.Create)
		auth.PUT("/:id", h.Update)
		auth.DELETE("/:id", h.Delete)
	}
}

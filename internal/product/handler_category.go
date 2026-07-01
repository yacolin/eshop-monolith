package product

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
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

// List 类目列表（分页）
// @Summary 类目列表
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Param name query string false "类目名称模糊搜索"
// @Param status query int false "状态 1-启用 0-禁用"
// @Param level query int false "层级 1-3"
// @Success 200 {object} response.Response{data=CategoryListResult}
// @Router /api/v1/categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	var req CategoryListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.List(c, &req)
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
	if level < 1 || level > 3 {
		c.Error(errcode.ErrCategoryLevelInvalid)
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

func RegisterCategoryRoutes(v1 *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	repo := NewCategoryRepository(db)
	svc := NewCategoryService(repo, rdb)
	h := NewCategoryHandler(svc)

	cats := v1.Group("/categories")
	{
		cats.GET("", h.List)
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

	// 异步预热类目缓存
	if rdb != nil {
		go func() {
			n, err := svc.WarmupCache(context.Background())
			if err != nil {
				log.Printf("[warmup] category cache warmup failed: %v", err)
			} else {
				log.Printf("[warmup] category cache warmup done: %d items", n)
			}
		}()
	}
}

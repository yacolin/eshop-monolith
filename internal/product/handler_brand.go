package product

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type BrandHandler struct {
	svc *BrandService
}

func NewBrandHandler(svc *BrandService) *BrandHandler {
	return &BrandHandler{svc: svc}
}

// Create 创建品牌
// @Summary 创建品牌
// @Tags brands
// @Accept json
// @Produce json
// @Param brand body CreateBrandReq true "品牌信息"
// @Success 200 {object} response.Response{data=Brand}
// @Router /api/v1/brands [post]
func (h *BrandHandler) Create(c *gin.Context) {
	var req CreateBrandReq
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

// GetByID 获取品牌详情
// @Summary 获取品牌详情
// @Tags brands
// @Produce json
// @Param id path int true "品牌ID"
// @Success 200 {object} response.Response{data=Brand}
// @Router /api/v1/brands/{id} [get]
func (h *BrandHandler) GetByID(c *gin.Context) {
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

// List 品牌列表
// @Summary 品牌列表
// @Tags brands
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Param name query string false "品牌名称模糊搜索"
// @Param first_letter query string false "首字母筛选"
// @Param status query int false "状态 1-启用 0-禁用"
// @Success 200 {object} response.Response{data=BrandListResult}
// @Router /api/v1/brands [get]
func (h *BrandHandler) List(c *gin.Context) {
	var req BrandListReq
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

// Update 更新品牌
// @Summary 更新品牌
// @Tags brands
// @Accept json
// @Produce json
// @Param id path int true "品牌ID"
// @Param brand body UpdateBrandReq true "品牌信息"
// @Success 200 {object} response.Response{data=Brand}
// @Router /api/v1/brands/{id} [put]
func (h *BrandHandler) Update(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req UpdateBrandReq
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

// Delete 删除品牌
// @Summary 删除品牌
// @Tags brands
// @Produce json
// @Param id path int true "品牌ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/brands/{id} [delete]
func (h *BrandHandler) Delete(c *gin.Context) {
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

func RegisterBrandRoutes(v1 *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	repo := NewBrandRepository(db)
	svc := NewBrandService(repo, rdb)
	h := NewBrandHandler(svc)

	brands := v1.Group("/brands")
	{
		brands.GET("", h.List)
		brands.GET("/:id", h.GetByID)
	}
	auth := v1.Group("/brands")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", h.Create)
		auth.PUT("/:id", h.Update)
		auth.DELETE("/:id", h.Delete)
	}

	// 异步预热品牌缓存
	if rdb != nil {
		go func() {
			n, err := svc.WarmupCache(context.Background())
			if err != nil {
				log.Printf("[warmup] brand cache warmup failed: %v", err)
			} else {
				log.Printf("[warmup] brand cache warmup done: %d items", n)
			}
		}()
	}
}

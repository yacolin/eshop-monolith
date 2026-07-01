package product

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"eshop-monolith/pkg/logger"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type SpuHandler struct {
	svc *SpuService
}

func NewSpuHandler(svc *SpuService) *SpuHandler {
	return &SpuHandler{svc: svc}
}

// Create 创建商品
// @Summary 创建商品
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateSPUReq true "商品信息（含 SKU 和属性）"
// @Success 200 {object} response.Response{data=SPU}
// @Router /api/v1/products [post]
func (h *SpuHandler) Create(c *gin.Context) {
	var req CreateSPUReq
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

// GetByID 获取商品详情
// @Summary 获取商品详情
// @Tags products
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response{data=SPUDetailResponse}
// @Router /api/v1/products/{id} [get]
func (h *SpuHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.GetDetailByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
	ReleaseDetailResponse(result)
}

// List 商品列表
// @Summary 商品列表（keyset 游标分页）
// @Tags products
// @Accept json
// @Produce json
// @Param size query int false "每页条数" default(10)
// @Param cursor query string false "游标（首次请求不传，后续使用上次返回的 cursor）"
// @Param name query string false "商品名称模糊搜索"
// @Param category_id query int false "类目ID"
// @Param brand_id query int false "品牌ID"
// @Param status query int false "状态 0-草稿 1-待审 2-上架 3-下架 4-封禁"
// @Param price_min query int false "最低价格（分）"
// @Param price_max query int false "最高价格（分）"
// @Success 200 {object} response.Response{data=SPUListResult}
// @Router /api/v1/products [get]
func (h *SpuHandler) List(c *gin.Context) {
	var req SPUListReq
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

// Update 更新商品
// @Summary 更新商品
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Param product body UpdateSPUReq true "商品信息"
// @Success 200 {object} response.Response{data=SPU}
// @Router /api/v1/products/{id} [put]
func (h *SpuHandler) Update(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req UpdateSPUReq
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

// Delete 删除商品
// @Summary 删除商品
// @Tags products
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/products/{id} [delete]
func (h *SpuHandler) Delete(c *gin.Context) {
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

func RegisterProductRoutes(v1 *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	repo := NewSpuRepository(db)
	catRepo := NewCategoryRepository(db)
	brandRepo := NewBrandRepository(db)
	attrRepo := NewAttributeRepository(db)
	svc := NewSpuService(repo, catRepo, brandRepo, attrRepo, db, rdb)
	h := NewSpuHandler(svc)

	// 异步预热 SPU 缓存（L1 + Bloom Filter + L2）
	if rdb != nil {
		go func() {
			n, err := svc.WarmupCache(context.Background())
			if err != nil {
				logger.Warn("SPU cache warmup failed", "error", err)
				log.Printf("[warmup] SPU cache warmup failed: %v", err)
			} else {
				logger.Info("SPU cache warmup done", "items", n)
				log.Printf("[warmup] SPU cache warmup done: %d items", n)
			}
		}()
	}

	products := v1.Group("/products")
	{
		products.GET("", h.List)
		products.GET("/:id", h.GetByID)
	}
	auth := v1.Group("/products")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", h.Create)
		auth.PUT("/:id", h.Update)
		auth.DELETE("/:id", h.Delete)
	}
}

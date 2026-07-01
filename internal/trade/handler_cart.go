package trade

import (
	"context"

	"eshop-monolith/internal/product"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// skuAdapter 适配 product.SKU → trade.SkuInfo
type skuAdapter struct {
	repo product.IspuRepository
}

func (a *skuAdapter) FindByID(ctx context.Context, skuID int64) (SkuInfo, error) {
	sku, err := a.repo.FindSKUByID(ctx, skuID)
	if err != nil {
		return nil, err
	}
	return skuItem{
		id: sku.ID, productID: sku.ProductID, skuCode: sku.SkuCode,
		price: sku.Price, image: sku.Image, specJSON: sku.Spec,
	}, nil
}

type skuItem struct {
	id, productID, price     int64
	skuCode, image, specJSON string
}

func (s skuItem) GetID() int64        { return s.id }
func (s skuItem) GetProductID() int64 { return s.productID }
func (s skuItem) GetSkuCode() string  { return s.skuCode }
func (s skuItem) GetPrice() int64     { return s.price }
func (s skuItem) GetImage() string    { return s.image }
func (s skuItem) GetSpecJSON() string { return s.specJSON }

// invAdapter 适配 inventory.InventoryService → trade.InventoryService
type invAdapter struct {
	svc *inventorySvc
}

type inventorySvc struct{}

func (s *inventorySvc) Lock(ctx context.Context, skuID int64, quantity int) error {
	return nil // stub
}

func (s *inventorySvc) Unlock(ctx context.Context, skuID int64, quantity int) error {
	return nil // stub
}

func getCurrentUser(c *gin.Context) int64 {
	uid, _ := c.Get("user_id")
	if uid == nil {
		return 0
	}
	id, _ := uid.(int64)
	return id
}

// ── Cart Handler ─────────────────────────────────

type CartHandler struct {
	svc *CartService
}

func NewCartHandler(svc *CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

// GetCart 获取购物车
// @Summary 获取购物车
// @Tags carts
// @Tags frontend
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=Cart}
// @Router /api/v1/carts [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	result, err := h.svc.GetCart(c, getCurrentUser(c), c.Query("session_id"))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// AddItem 添加商品到购物车
// @Summary 添加商品到购物车
// @Tags carts
// @Tags frontend
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body AddItemReq true "商品信息"
// @Success 200 {object} response.Response
// @Router /api/v1/carts [post]
func (h *CartHandler) AddItem(c *gin.Context) {
	var req AddItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.AddItem(c, getCurrentUser(c), c.Query("session_id"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// UpdateItem 更新购物车商品
// @Summary 更新购物车商品
// @Tags carts
// @Tags frontend
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body UpdateItemReq true "更新信息"
// @Success 200 {object} response.Response
// @Router /api/v1/carts [put]
func (h *CartHandler) UpdateItem(c *gin.Context) {
	var req UpdateItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.UpdateQuantity(c, getCurrentUser(c), c.Query("session_id"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// RemoveItem 删除购物车商品
// @Summary 删除购物车商品
// @Tags carts
// @Tags frontend
// @Security ApiKeyAuth
// @Produce json
// @Param sku_id query int true "SKU ID"
// @Success 200 {object} response.Response
// @Router /api/v1/carts [delete]
func (h *CartHandler) RemoveItem(c *gin.Context) {
	var req struct {
		SkuID int64 `form:"sku_id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.RemoveItem(c, getCurrentUser(c), c.Query("session_id"), req.SkuID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ClearCart 清空购物车
// @Summary 清空购物车
// @Tags carts
// @Tags frontend
// @Security ApiKeyAuth
// @Produce json
// @Param session_id query string false "会话ID"
// @Success 200 {object} response.Response
// @Router /api/v1/carts/clear [post]
func (h *CartHandler) ClearCart(c *gin.Context) {
	if err := h.svc.ClearCart(c, getCurrentUser(c), c.Query("session_id")); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

func RegisterCartRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	repo := NewCartRepository(db)
	skuA := &skuAdapter{repo: product.NewSpuRepository(db)}
	svc := NewCartService(repo, skuA, db)
	h := NewCartHandler(svc)

	v1.Group("/carts").GET("", h.GetCart)
	auth := v1.Group("/carts")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("/items", h.AddItem)
		auth.PUT("/items", h.UpdateItem)
		auth.DELETE("/items", h.RemoveItem)
		auth.DELETE("", h.ClearCart)
	}
}

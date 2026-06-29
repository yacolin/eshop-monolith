package cart

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/internal/product"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
)

// skuInfo 适配 product.SKU → SkuInfo（不能对外部类型定义方法）
type skuInfo struct {
	id        int64
	productID int64
	price     int64
	image     string
	specJSON  string
}

func (s *skuInfo) GetID() int64        { return s.id }
func (s *skuInfo) GetProductID() int64  { return s.productID }
func (s *skuInfo) GetPrice() int64      { return s.price }
func (s *skuInfo) GetImage() string     { return s.image }
func (s *skuInfo) GetSpecJSON() string  { return s.specJSON }

// skuProvider 适配 product.SpuRepository → cart.SkuProvider
type skuProvider struct {
	repo product.IspuRepository
}

func (p *skuProvider) FindByID(ctx context.Context, skuID int64) (SkuInfo, error) {
	sku, err := p.repo.FindSKUByID(ctx, skuID)
	if err != nil {
		return nil, err
	}
	return &skuInfo{
		id:        sku.ID,
		productID: sku.ProductID,
		price:     sku.Price,
		image:     sku.Image,
		specJSON:  sku.Spec,
	}, nil
}

// ── Handler ──────────────────────────────────────

type CartHandler struct {
	svc *CartService
}

func NewCartHandler(svc *CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

// getCurrentUser 获取当前用户 ID（未登录返回 0）
func getCurrentUser(c *gin.Context) int64 {
	uid, _ := c.Get("user_id")
	if uid == nil {
		return 0
	}
	id, _ := uid.(int64)
	return id
}

// GetCart 获取购物车
// @Summary 获取购物车
// @Tags carts
// @Produce json
// @Success 200 {object} response.Response{data=CartResponse}
// @Router /api/v1/carts [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := getCurrentUser(c)
	sessionID := c.Query("session_id")
	result, err := h.svc.GetCart(c, userID, sessionID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// AddItem 添加商品
// @Summary 添加商品
// @Tags carts
// @Accept json
// @Produce json
// @Param req body AddItemReq true "商品信息"
// @Success 200 {object} response.Response{data=CartResponse}
// @Router /api/v1/carts/items [post]
func (h *CartHandler) AddItem(c *gin.Context) {
	userID := getCurrentUser(c)
	sessionID := c.Query("session_id")
	var req AddItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.AddItem(c, userID, sessionID, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// UpdateItem 更新商品数量
// @Summary 更新商品数量
// @Tags carts
// @Accept json
// @Produce json
// @Param req body UpdateItemReq true "数量信息"
// @Success 200 {object} response.Response{data=CartResponse}
// @Router /api/v1/carts/items [put]
func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID := getCurrentUser(c)
	sessionID := c.Query("session_id")
	var req UpdateItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.UpdateQuantity(c, userID, sessionID, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// RemoveItem 删除商品
// @Summary 删除购物车商品
// @Tags carts
// @Produce json
// @Param sku_id query int true "SKU ID"
// @Success 200 {object} response.Response{data=CartResponse}
// @Router /api/v1/carts/items [delete]
func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID := getCurrentUser(c)
	sessionID := c.Query("session_id")
	var req struct {
		SkuID int64 `form:"sku_id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.RemoveItem(c, userID, sessionID, req.SkuID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ClearCart 清空购物车
// @Summary 清空购物车
// @Tags carts
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/carts [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := getCurrentUser(c)
	sessionID := c.Query("session_id")
	if err := h.svc.ClearCart(c, userID, sessionID); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// ── Routes ────────────────────────────────────────

func RegisterCartRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	repo := NewCartRepository(db)
	skuRepo := product.NewSpuRepository(db)
	provider := &skuProvider{repo: skuRepo}
	svc := NewCartService(repo, provider, db)
	h := NewCartHandler(svc)

	carts := v1.Group("/carts")
	{
		carts.GET("", h.GetCart)
	}
	auth := v1.Group("/carts")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("/items", h.AddItem)
		auth.PUT("/items", h.UpdateItem)
		auth.DELETE("/items", h.RemoveItem)
		auth.DELETE("", h.ClearCart)
	}
}

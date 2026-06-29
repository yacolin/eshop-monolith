package order

import (
	"context"

	"eshop-monolith/internal/inventory"
	"eshop-monolith/internal/product"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ── SKU 适配 ─────────────────────────────────────

type skuInfo struct {
	id        int64
	productID int64
	skuCode   string
	price     int64
	image     string
	specJSON  string
}

func (s *skuInfo) GetID() int64       { return s.id }
func (s *skuInfo) GetProductID() int64 { return s.productID }
func (s *skuInfo) GetSkuCode() string  { return s.skuCode }
func (s *skuInfo) GetPrice() int64     { return s.price }
func (s *skuInfo) GetImage() string    { return s.image }
func (s *skuInfo) GetSpecJSON() string { return s.specJSON }

type skuService struct {
	repo product.IspuRepository
}

func (s *skuService) FindByID(ctx context.Context, skuID int64) (SkuInfo, error) {
	sku, err := s.repo.FindSKUByID(ctx, skuID)
	if err != nil {
		return nil, err
	}
	return &skuInfo{
		id:        sku.ID,
		productID: sku.ProductID,
		skuCode:   sku.SkuCode,
		price:     sku.Price,
		image:     sku.Image,
		specJSON:  sku.Spec,
	}, nil
}

// ── Inventory 适配 ───────────────────────────────

type invService struct {
	svc *inventory.InventoryService
}

func (s *invService) Lock(ctx context.Context, skuID int64, quantity int) error {
	_, err := s.svc.Lock(ctx, &inventory.LockStockReq{
		SkuID:    skuID,
		Quantity: int64(quantity),
		Operator: "order",
	})
	return err
}

func (s *invService) Unlock(ctx context.Context, skuID int64, quantity int) error {
	_, err := s.svc.Unlock(ctx, &inventory.UnlockStockReq{
		SkuID:    skuID,
		Quantity: int64(quantity),
		Operator: "order",
	})
	return err
}

// ── Handler ──────────────────────────────────────

type OrderHandler struct {
	svc *OrderService
}

func NewOrderHandler(svc *OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// Create 创建订单
// @Summary 创建订单
// @Tags orders
// @Accept json
// @Produce json
// @Param req body CreateOrderReq true "订单信息"
// @Success 200 {object} response.Response{data=Order}
// @Router /api/v1/orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	userID := int64(1) // TODO: get from JWT
	result, err := h.svc.Create(c, userID, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// GetByOrderNo 查询订单
// @Summary 查询订单
// @Tags orders
// @Produce json
// @Param order_no path string true "订单号"
// @Success 200 {object} response.Response{data=Order}
// @Router /api/v1/orders/:order_no [get]
func (h *OrderHandler) GetByOrderNo(c *gin.Context) {
	orderNo := c.Param("order_no")
	result, err := h.svc.GetByOrderNo(c, orderNo)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// List 订单列表
// @Summary 订单列表
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Param status query string false "订单状态筛选"
// @Param user_id query int false "用户ID筛选"
// @Success 200 {object} response.Response{data=OrderListResult}
// @Router /api/v1/orders [get]
func (h *OrderHandler) List(c *gin.Context) {
	var req OrderListReq
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

// UpdateStatus 更新订单状态
// @Summary 更新订单状态
// @Tags orders
// @Accept json
// @Produce json
// @Param order_no path string true "订单号"
// @Param req body UpdateOrderStatusReq true "状态信息"
// @Success 200 {object} response.Response{data=Order}
// @Router /api/v1/orders/:order_no/status [put]
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	orderNo := c.Param("order_no")
	var req UpdateOrderStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.UpdateStatus(c, orderNo, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ── Routes ────────────────────────────────────────

func RegisterOrderRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	repo := NewOrderRepository(db)
	skuSvc := &skuService{repo: product.NewSpuRepository(db)}
	invSvc := &invService{svc: inventory.NewInventoryService(inventory.NewInventoryRepository(db), db)}
	svc := NewOrderService(repo, skuSvc, invSvc, db)
	h := NewOrderHandler(svc)

	orders := v1.Group("/orders")
	orders.Use(middleware.JWTAuth())
	{
		orders.POST("", h.Create)
		orders.GET("", h.List)
		orders.GET("/:order_no", h.GetByOrderNo)
		orders.PUT("/:order_no/status", h.UpdateStatus)
	}
}

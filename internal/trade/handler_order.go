package trade

import (
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/product"
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	svc *OrderService
}

func NewOrderHandler(svc *OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// Create 创建订单
// @Summary 创建订单
// @Tags orders
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreateOrderReq true "订单信息"
// @Success 200 {object} response.Response{data=Order}
// @Router /api/v1/orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Create(c, getCurrentUser(c), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// GetByOrderNo 根据订单号查询
// @Summary 根据订单号查询
// @Tags orders
// @Security ApiKeyAuth
// @Produce json
// @Param order_no path string true "订单号"
// @Success 200 {object} response.Response{data=Order}
// @Router /api/v1/orders/{order_no} [get]
func (h *OrderHandler) GetByOrderNo(c *gin.Context) {
	result, err := h.svc.GetByOrderNo(c, c.Param("order_no"))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// List 订单列表
// @Summary 订单列表
// @Tags orders
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
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
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param order_no path string true "订单号"
// @Param request body UpdateOrderStatusReq true "状态信息"
// @Success 200 {object} response.Response{data=Order}
// @Router /api/v1/orders/{order_no}/status [put]
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	var req UpdateOrderStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.UpdateStatus(c, c.Param("order_no"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

func RegisterOrderRoutes(v1 *gin.RouterGroup, db *gorm.DB, mqClient *rabbitmq.Client) {
	repo := NewOrderRepository(db)
	skuA := &skuAdapter{repo: product.NewSpuRepository(db)}
	svc := NewOrderService(repo, skuA, &inventorySvc{}, db, mqClient)
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

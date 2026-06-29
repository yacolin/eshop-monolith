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

func (h *OrderHandler) GetByOrderNo(c *gin.Context) {
	result, err := h.svc.GetByOrderNo(c, c.Param("order_no"))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

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

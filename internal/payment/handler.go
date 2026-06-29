package payment

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
)

type PaymentHandler struct {
	svc *PaymentService
}

func NewPaymentHandler(svc *PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// CreatePayment 创建支付
// @Summary 创建支付
// @Tags payments
// @Accept json
// @Produce json
// @Param req body CreatePaymentReq true "支付请求"
// @Success 200 {object} response.Response{data=Payment}
// @Router /api/v1/payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req CreatePaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.CreatePayment(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// HandleCallback 支付回调
// @Summary 支付回调
// @Tags payments
// @Accept json
// @Produce json
// @Param req body PaymentCallbackReq true "回调请求"
// @Success 200 {object} response.Response{data=Payment}
// @Router /api/v1/payments/callback [post]
func (h *PaymentHandler) HandleCallback(c *gin.Context) {
	var req PaymentCallbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.HandleCallback(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// GetPayment 查询支付
// @Summary 查询支付
// @Tags payments
// @Produce json
// @Param payment_no query string false "支付单号"
// @Param order_no query string false "订单号"
// @Success 200 {object} response.Response{data=Payment}
// @Router /api/v1/payments [get]
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentNo := c.Query("payment_no")
	orderNo := c.Query("order_no")
	var result interface{}
	var err error
	if paymentNo != "" {
		result, err = h.svc.GetPayment(c, paymentNo)
	} else if orderNo != "" {
		result, err = h.svc.GetPaymentByOrder(c, orderNo)
	} else {
		response.Success(c, nil)
		return
	}
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// GetPaymentByOrder 按订单查支付
// @Summary 按订单查支付
// @Tags payments
// @Produce json
// @Param order_no path string true "订单号"
// @Success 200 {object} response.Response{data=Payment}
// @Router /api/v1/orders/payment/:order_no [get]
func (h *PaymentHandler) GetPaymentByOrder(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		c.Error(nil)
		return
	}
	result, err := h.svc.GetPaymentByOrder(c, orderNo)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// CreateRefund 创建退款
// @Summary 创建退款
// @Tags payments
// @Accept json
// @Produce json
// @Param req body CreateRefundReq true "退款请求"
// @Success 200 {object} response.Response{data=Refund}
// @Router /api/v1/refunds [post]
func (h *PaymentHandler) CreateRefund(c *gin.Context) {
	var req CreateRefundReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.CreateRefund(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// HandleRefundCallback 退款回调
// @Summary 退款回调
// @Tags payments
// @Accept json
// @Produce json
// @Param req body RefundCallbackReq true "退款回调"
// @Success 200 {object} response.Response{data=Refund}
// @Router /api/v1/refunds/callback [post]
func (h *PaymentHandler) HandleRefundCallback(c *gin.Context) {
	var req RefundCallbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.HandleRefundCallback(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ── Routes ────────────────────────────────────────

func RegisterPaymentRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	repo := NewPaymentRepository(db)
	svc := NewPaymentService(repo, db)
	h := NewPaymentHandler(svc)

	payments := v1.Group("/payments")
	{
		payments.GET("", h.GetPayment)
		payments.POST("/callback", h.HandleCallback)
	}
	auth := v1.Group("/payments")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("", h.CreatePayment)
	}

	refunds := v1.Group("/refunds")
	refunds.Use(middleware.JWTAuth())
	{
		refunds.POST("", h.CreateRefund)
		refunds.POST("/callback", h.HandleRefundCallback)
	}

	// 兼容旧路由（order 模块引用）
	orderPayment := v1.Group("/orders")
	orderPayment.Use(middleware.JWTAuth())
	{
		orderPayment.GET("/payment/:order_no", h.GetPaymentByOrder)
	}
}

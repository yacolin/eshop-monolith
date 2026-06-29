package trade

import (
	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	svc *PaymentService
}

func NewPaymentHandler(svc *PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

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

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentNo, orderNo := c.Query("payment_no"), c.Query("order_no")
	if paymentNo == "" && orderNo == "" {
		response.Success(c, nil)
		return
	}
	if paymentNo != "" {
		result, err := h.svc.GetPayment(c, paymentNo)
		if err != nil { c.Error(err); return }
		response.Success(c, result)
		return
	}
	result, err := h.svc.repo.FindByOrderNo(c, orderNo)
	if err != nil { c.Error(err); return }
	response.Success(c, result)
}

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
	}
}

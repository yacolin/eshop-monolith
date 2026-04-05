package routes

import (
	"eshop-monolith/internal/eventbus"
	"eshop-monolith/internal/payment/api/handlers"
	"eshop-monolith/internal/payment/domain/repositories"
	"eshop-monolith/internal/payment/service"
	"eshop-monolith/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPaymentRoutes 注册支付相关路由
func RegisterPaymentRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, bus *eventbus.Bus, db *gorm.DB) {
	// 创建仓储实例
	paymentRepo := repositories.NewPaymentRepository(db)
	refundRepo := repositories.NewRefundRepository(db)
	paymentMethodRepo := repositories.NewPaymentMethodRepository(db)

	// 创建服务实例
	paymentService := service.NewPaymentService(
		paymentRepo,
		refundRepo,
		paymentMethodRepo,
		repos.Order,
		bus,
	)

	// 创建处理器实例
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// 支付相关路由
	payments := v1.Group("/payments")
	{
		payments.POST("", paymentHandler.CreatePayment)           // 创建支付
		payments.GET("", paymentHandler.ListPayments)           // 获取支付列表
		payments.GET("/:id", paymentHandler.GetPayment)         // 获取支付详情
		payments.PATCH("/:id/status", paymentHandler.UpdatePaymentStatus) // 更新支付状态
	}

	// 订单支付路由
	orders := v1.Group("/orders")
	{
		orders.GET("/:order_id/payment", paymentHandler.GetPaymentByOrderID) // 根据订单ID获取支付
	}

	// 退款相关路由
	refunds := v1.Group("/refunds")
	{
		refunds.POST("", paymentHandler.CreateRefund)           // 创建退款
		refunds.GET("", paymentHandler.ListRefunds)           // 获取退款列表
		refunds.PATCH("/:id/status", paymentHandler.UpdateRefundStatus) // 更新退款状态
	}

	// 支付方式路由
	paymentMethods := v1.Group("/payment-methods")
	{
		paymentMethods.GET("", paymentHandler.ListPaymentMethods) // 获取支付方式列表
	}
}

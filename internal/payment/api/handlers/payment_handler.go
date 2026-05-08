package handlers

import (
	"strconv"

	"eshop-monolith/internal/payment/api/dto"
	"eshop-monolith/internal/payment/service"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

// PaymentHandler 支付处理器
type PaymentHandler struct {
	paymentService *service.PaymentService
}

// NewPaymentHandler 创建支付处理器实例
func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreatePayment 创建支付
// @Summary 创建支付
// @Description 创建新的支付记录
// @Tags 支付
// @Accept json
// @Produce json
// @Param request body dto.CreatePaymentRequest true "创建支付请求"
// @Success 200 {object} response.Response{data=dto.PaymentResponse} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	payment, paymentURL, err := h.paymentService.CreatePayment(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentMethod: payment.PaymentMethod,
		TransactionID: payment.TransactionID,
		Status:        payment.Status,
		PaymentURL:    paymentURL,
		CreatedAt:     payment.CreatedAt.ToTime().String(),
	})
}

// GetPayment 获取支付详情
// @Summary 获取支付详情
// @Description 根据ID获取支付记录详情
// @Tags 支付
// @Accept json
// @Produce json
// @Param id path int true "支付ID"
// @Success 200 {object} response.Response{data=dto.PaymentResponse} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "支付记录不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/payments/{id} [get]
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	payment, err := h.paymentService.GetPaymentByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	responseData := dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentMethod: payment.PaymentMethod,
		TransactionID: payment.TransactionID,
		Status:        payment.Status,
		CreatedAt:     payment.CreatedAt.ToTime().String(),
	}

	if payment.PaidAt != nil {
		responseData.PaidAt = payment.PaidAt.ToTime().String()
	}

	response.Success(c, responseData)
}

// GetPaymentByOrderID 根据订单ID获取支付
// @Summary 根据订单ID获取支付
// @Description 根据订单ID获取支付记录
// @Tags 支付
// @Accept json
// @Produce json
// @Param order_id path int true "订单ID"
// @Success 200 {object} response.Response{data=dto.PaymentResponse} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "支付记录不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/orders/payment/{order_id} [get]
func (h *PaymentHandler) GetPaymentByOrderID(c *gin.Context) {
	orderID, err := utils.ParseIntParam(c, "order_id")
	if err != nil {
		c.Error(err)
		return
	}

	payment, err := h.paymentService.GetPaymentByOrderID(c, orderID)
	if err != nil {
		c.Error(err)
		return
	}

	responseData := dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentMethod: payment.PaymentMethod,
		TransactionID: payment.TransactionID,
		Status:        payment.Status,
		CreatedAt:     payment.CreatedAt.ToTime().String(),
	}

	if payment.PaidAt != nil {
		responseData.PaidAt = payment.PaidAt.ToTime().String()
	}

	response.Success(c, responseData)
}

// UpdatePaymentStatus 更新支付状态
// @Summary 更新支付状态
// @Description 更新支付记录的状态
// @Tags 支付
// @Accept json
// @Produce json
// @Param id path int true "支付ID"
// @Param request body dto.UpdatePaymentStatusRequest true "更新支付状态请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "支付记录不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/payments/{id}/status [patch]
func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.paymentService.UpdatePaymentStatus(c, id, req.Status, req.TransactionID, req.FailureReason); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// ListPayments 获取支付列表
// @Summary 获取支付列表
// @Description 根据查询条件获取支付记录列表
// @Tags 支付
// @Accept json
// @Produce json
// @Param order_id query int false "订单ID"
// @Param payment_method query string false "支付方式"
// @Param status query string false "支付状态"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认20"
// @Param sort_by query string false "排序字段，默认created_at"
// @Param order query string false "排序方向，默认desc"
// @Success 200 {object} response.Response{data=dto.PaymentListResponse} "成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/payments [get]
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	var q dto.PaymentListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	// normalize pagination values (ensure page>=1, 1<=size<=100)
	(&q).Normalize()

	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.Order == "" {
		q.Order = "desc"
	}

	result, err := h.paymentService.ListPayments(c, q)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// CreateRefund 创建退款
// @Summary 创建退款
// @Description 创建新的退款记录
// @Tags 退款
// @Accept json
// @Produce json
// @Param request body dto.CreateRefundRequest true "创建退款请求"
// @Success 200 {object} response.Response{data=dto.RefundResponse} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/refunds [post]
func (h *PaymentHandler) CreateRefund(c *gin.Context) {
	var req dto.CreateRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	refund, err := h.paymentService.CreateRefund(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, dto.RefundResponse{
		ID:            refund.ID,
		PaymentID:     refund.PaymentID,
		OrderID:       refund.OrderID,
		RefundAmount:  refund.RefundAmount,
		RefundReason:  refund.RefundReason,
		TransactionID: refund.TransactionID,
		Status:        refund.Status,
		CreatedAt:     refund.CreatedAt.ToTime().String(),
	})
}

// UpdateRefundStatus 更新退款状态
// @Summary 更新退款状态
// @Description 更新退款记录的状态
// @Tags 退款
// @Accept json
// @Produce json
// @Param id path int true "退款ID"
// @Param request body dto.UpdatePaymentStatusRequest true "更新退款状态请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "退款记录不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/refunds/{id}/status [patch]
func (h *PaymentHandler) UpdateRefundStatus(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.paymentService.UpdateRefundStatus(c, id, req.Status, req.TransactionID, req.FailureReason); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// ListRefunds 获取退款列表
// @Summary 获取退款列表
// @Description 根据查询条件获取退款记录列表
// @Tags 退款
// @Accept json
// @Produce json
// @Param payment_id query int false "支付ID"
// @Param order_id query int false "订单ID"
// @Param status query string false "退款状态"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认20"
// @Param sort_by query string false "排序字段，默认created_at"
// @Param order query string false "排序方向，默认desc"
// @Success 200 {object} response.Response{data=dto.RefundListResult} "成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/refunds [get]
func (h *PaymentHandler) ListRefunds(c *gin.Context) {
	var q dto.RefundListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	// 设置默认值
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 || q.PageSize > 100 {
		q.PageSize = 20
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.Order == "" {
		q.Order = "desc"
	}

	result, err := h.paymentService.ListRefunds(c, q)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// ListPaymentMethods 获取支付方式列表
// @Summary 获取支付方式列表
// @Description 获取所有可用的支付方式
// @Tags 支付方式
// @Accept json
// @Produce json
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认20"
// @Success 200 {object} response.Response{data=dto.PaymentMethodListResponse} "成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/payment-methods [get]
func (h *PaymentHandler) ListPaymentMethods(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.paymentService.ListPaymentMethods(c, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

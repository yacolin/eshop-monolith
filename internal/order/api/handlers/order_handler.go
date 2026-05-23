package handlers

import (
	"strconv"

	"eshop-monolith/internal/order/api/dto"
	"eshop-monolith/internal/order/service"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// ListOrders 列出所有订单
// @Summary 列出所有订单
// @Description 获取所有订单的列表，支持分页筛选
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Param customer_id query int false "用户ID过滤"
// @Param status query string false "订单状态过滤"
// @Success 200 {object} response.Response{data=dto.OrderListResult}
// @Router /api/v1/orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	var q dto.OrderListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	// normalize pagination values (ensure page>=1, 1<=size<=100)
	(&q).Normalize()

	result, err := h.orderService.ListOrders(c, q)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// GetOrder 根据ID获取订单
// @Summary 获取订单详情
// @Description 根据订单ID获取订单详情
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} response.Response{data=models.Order}
// @Router /api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	order, err := h.orderService.GetOrder(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, order)
}

// CreateOrder 创建订单
// @Summary 创建订单
// @Description 创建一个新的订单
// @Tags orders
// @Accept json
// @Produce json
// @Param order body dto.CreateOrderDTO true "订单信息"
// @Success 200 {object} response.Response{data=models.Order}
// @Router /api/v1/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	order, err := h.orderService.CreateOrder(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, order)
}

// UpdateOrder 更新订单
// @Summary 更新订单
// @Description 根据ID更新订单信息
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param order body dto.UpdateOrderDTO true "订单信息"
// @Success 200 {object} response.Response{data=models.Order}
// @Router /api/v1/orders/{id} [put]
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateOrderDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	order, err := h.orderService.UpdateOrder(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, order)
}

// CancelOrder 取消订单
// @Summary 取消订单
// @Description 取消指定订单
// @Tags orders
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.orderService.CancelOrder(c, id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Order cancelled successfully"})
}

// UpdateOrderStatus 更新订单状态
// @Summary 更新订单状态
// @Description 更新指定订单的状态
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param status body dto.UpdateOrderStatusDTO true "订单状态"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/orders/{id}/status [patch]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateOrderStatusDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.orderService.UpdateOrderStatus(c, id, req.Status); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Order status updated successfully"})
}

// DeleteOrder 删除订单
// @Summary 删除订单
// @Description 根据ID删除订单
// @Tags orders
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.orderService.DeleteOrder(c, id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Order deleted successfully"})
}

// GetOrdersByUserID 根据用户ID获取订单列表
// @Summary 根据用户ID获取订单列表
// @Description 根据用户ID获取订单列表，支持分页
// @Tags orders
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=dto.OrderListResult}
// @Router /api/v1/users/{user_id}/orders [get]
func (h *OrderHandler) GetOrdersByUserID(c *gin.Context) {
	userID, err := utils.ParseIntParam(c, "user_id")
	if err != nil {
		c.Error(err)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 确保分页参数有效
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	orders, total, err := h.orderService.GetOrdersByUserID(c, userID, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"list":  orders,
		"total": total,
	})
}

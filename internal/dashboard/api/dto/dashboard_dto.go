// Package dto 定义 Dashboard 模块的请求/响应数据结构
package dto

// DashboardResponse 仪表盘汇总数据（全部聚合在一个响应中）
type DashboardResponse struct {
	Summary             SummaryDTO              `json:"summary"`
	OrderTrend          []OrderTrendDTO         `json:"order_trend"`
	OrderStatusDist     []StatusDistDTO         `json:"order_status_dist"`
	PaymentMethodDist   []MethodDistDTO         `json:"payment_method_dist"`
	CategoryDist        []CategoryDistDTO       `json:"category_dist"`
	InventoryStatusDist []StatusDistDTO         `json:"inventory_status_dist"`
	TopProducts         []TopProductDTO         `json:"top_products"`
}

// SummaryDTO 核心指标
type SummaryDTO struct {
	TotalOrders   int64 `json:"total_orders"`    // 总订单数
	TotalRevenue  int64 `json:"total_revenue"`   // 总营收（分）
	TotalProducts int64 `json:"total_products"`  // 商品总数
	LowStockCount int64 `json:"low_stock_count"` // 库存告警数
}

// OrderTrendDTO 订单趋势（按天聚合）
type OrderTrendDTO struct {
	Date   string `json:"date"`   // 日期 (MM-DD)
	Count  int64  `json:"count"`  // 订单数
	Amount int64  `json:"amount"` // 金额（分）
}

// StatusDistDTO 状态分布
type StatusDistDTO struct {
	Status string `json:"status"` // 状态编码
	Label  string `json:"label"`  // 显示标签
	Value  int64  `json:"value"`  // 数量
}

// MethodDistDTO 支付方式分布
type MethodDistDTO struct {
	Method string `json:"method"` // 支付方式编码
	Label  string `json:"label"`  // 显示标签
	Value  int64  `json:"value"`  // 数量
}

// CategoryDistDTO 分类分布
type CategoryDistDTO struct {
	Category string `json:"category"` // 分类名称
	Value    int64  `json:"value"`    // 商品数量
}

// TopProductDTO 热销商品
type TopProductDTO struct {
	Name   string `json:"name"`   // 商品名称
	Count  int64  `json:"count"`  // 销量
	Amount int64  `json:"amount"` // 销售额（分）
}

package dashboard

type DashboardResponse struct {
	Summary             SummaryDTO        `json:"summary"`
	OrderTrend          []OrderTrendDTO   `json:"order_trend"`
	OrderStatusDist     []StatusDistDTO   `json:"order_status_dist"`
	PaymentMethodDist   []MethodDistDTO   `json:"payment_method_dist"`
	CategoryDist        []CategoryDistDTO `json:"category_dist"`
	InventoryStatusDist []StatusDistDTO   `json:"inventory_status_dist"`
	TopProducts         []TopProductDTO   `json:"top_products"`
}

type SummaryDTO struct {
	TotalOrders   int64 `json:"total_orders"`
	TotalRevenue  int64 `json:"total_revenue"`
	TotalProducts int64 `json:"total_products"`
	LowStockCount int64 `json:"low_stock_count"`
}

type OrderTrendDTO struct {
	Date   string `json:"date"`
	Count  int64  `json:"count"`
	Amount int64  `json:"amount"`
}

type StatusDistDTO struct {
	Status string `json:"status"`
	Label  string `json:"label"`
	Value  int64  `json:"value"`
}

type MethodDistDTO struct {
	Method string `json:"method"`
	Label  string `json:"label"`
	Value  int64  `json:"value"`
}

type CategoryDistDTO struct {
	Category string `json:"category"`
	Value    int64  `json:"value"`
}

type TopProductDTO struct {
	Name   string `json:"name"`
	Count  int64  `json:"count"`
	Amount int64  `json:"amount"`
}

package service

import (
	"context"
	"fmt"
	"time"

	"eshop-monolith/internal/dashboard/api/dto"
	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/internal/infra/repository"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

const (
	dashboardCacheKey = "dashboard:stats"
	dashboardCacheTTL = 5 * time.Minute
)

// DashboardService 仪表盘数据聚合服务
type DashboardService struct {
	db          *gorm.DB
	rdb         *redis.Client
	repos       *repository.Repositories
	bus         *eventbus.Bus
	singleGroup singleflight.Group
}

// NewDashboardService 创建仪表盘服务
func NewDashboardService(
	db *gorm.DB,
	rdb *redis.Client,
	repos *repository.Repositories,
	bus *eventbus.Bus,
) *DashboardService {
	return &DashboardService{
		db:    db,
		rdb:   rdb,
		repos: repos,
		bus:   bus,
	}
}

// GetStats 获取仪表盘全部汇总数据（带 Redis 缓存）
func (s *DashboardService) GetStats(ctx context.Context) (*dto.DashboardResponse, error) {
	v, err, _ := s.singleGroup.Do(dashboardCacheKey, func() (interface{}, error) {
		// 1. 尝试从缓存读取
		cached, err := s.rdb.Get(ctx, dashboardCacheKey).Bytes()
		if err == nil {
			var resp dto.DashboardResponse
			if err := sonic.Unmarshal(cached, &resp); err == nil {
				return &resp, nil
			}
		}

		// 2. 从数据库计算全量汇总
		resp := &dto.DashboardResponse{}
		resp.Summary = s.computeSummary(ctx)
		resp.OrderTrend = s.computeOrderTrend(ctx)
		resp.OrderStatusDist = s.computeOrderStatusDist(ctx)
		resp.PaymentMethodDist = s.computePaymentMethodDist(ctx)
		resp.CategoryDist = s.computeCategoryDist(ctx)
		resp.InventoryStatusDist = s.computeInventoryStatusDist(ctx)
		resp.TopProducts = s.computeTopProducts(ctx)

		// 3. 写入缓存（使用背景 context 避免请求取消导致缓存失效）
		data, err := sonic.Marshal(resp)
		if err == nil {
			s.rdb.Set(context.Background(), dashboardCacheKey, data, dashboardCacheTTL)
		}

		return resp, nil
	})

	if err != nil {
		return nil, err
	}
	return v.(*dto.DashboardResponse), nil
}

// RefreshCache 强制刷新缓存（供启动预热和定时任务使用）
func (s *DashboardService) RefreshCache(ctx context.Context) error {
	resp := &dto.DashboardResponse{}
	resp.Summary = s.computeSummary(ctx)
	resp.OrderTrend = s.computeOrderTrend(ctx)
	resp.OrderStatusDist = s.computeOrderStatusDist(ctx)
	resp.PaymentMethodDist = s.computePaymentMethodDist(ctx)
	resp.CategoryDist = s.computeCategoryDist(ctx)
	resp.InventoryStatusDist = s.computeInventoryStatusDist(ctx)
	resp.TopProducts = s.computeTopProducts(ctx)

	data, err := sonic.Marshal(resp)
	if err != nil {
		return err
	}

	// 用背景 context 写入缓存，不受传入 ctx 取消影响
	return s.rdb.Set(context.Background(), dashboardCacheKey, data, dashboardCacheTTL).Err()
}

// InvalidateCache 主动失效缓存（供事件处理器调用）
func (s *DashboardService) InvalidateCache() {
	s.rdb.Del(context.Background(), dashboardCacheKey)
}

// RegisterEventHandlers 注册事件处理器，数据变更时自动失效缓存
func (s *DashboardService) RegisterEventHandlers() {
	// 订单相关事件
	s.bus.Subscribe("eshop-monolith/internal/order/events.OrderCreatedEvent", func(interface{}) {
		s.InvalidateCache()
	})
	s.bus.Subscribe("eshop-monolith/internal/order/events.OrderCancelledEvent", func(interface{}) {
		s.InvalidateCache()
	})

	// 支付相关事件
	s.bus.Subscribe("eshop-monolith/internal/payment/events.PaymentSuccessEvent", func(interface{}) {
		s.InvalidateCache()
	})
	s.bus.Subscribe("eshop-monolith/internal/payment/events.PaymentStatusUpdatedEvent", func(interface{}) {
		s.InvalidateCache()
	})
	s.bus.Subscribe("eshop-monolith/internal/payment/events.RefundCreatedEvent", func(interface{}) {
		s.InvalidateCache()
	})

	// 产品相关事件
	s.bus.Subscribe("eshop-monolith/internal/inventory/events.ProductCreatedEvent", func(interface{}) {
		s.InvalidateCache()
	})
	s.bus.Subscribe("eshop-monolith/internal/inventory/events.ProductUpdatedEvent", func(interface{}) {
		s.InvalidateCache()
	})
	s.bus.Subscribe("eshop-monolith/internal/inventory/events.ProductDeletedEvent", func(interface{}) {
		s.InvalidateCache()
	})
}

// computeSummary 计算核心指标
func (s *DashboardService) computeSummary(ctx context.Context) dto.SummaryDTO {
	var summary dto.SummaryDTO

	// 总订单数
	s.db.WithContext(ctx).Model(&struct{}{}).Table("orders").
		Select("COUNT(*)").Scan(&summary.TotalOrders)

	// 总营收（已成功支付的金额）
	s.db.WithContext(ctx).Model(&struct{}{}).Table("payments").
		Select("COALESCE(SUM(amount), 0)").
		Where("status IN ?", []string{"success", "paid"}).
		Scan(&summary.TotalRevenue)

	// 商品总数
	s.db.WithContext(ctx).Model(&struct{}{}).Table("products").
		Select("COUNT(*)").Scan(&summary.TotalProducts)

	// 库存告警数
	s.db.WithContext(ctx).Model(&struct{}{}).Table("inventories").
		Select("COUNT(*)").Where("status IN ?", []string{"lowstock", "outofstock"}).
		Scan(&summary.LowStockCount)

	return summary
}

// computeOrderTrend 近7天订单趋势
func (s *DashboardService) computeOrderTrend(ctx context.Context) []dto.OrderTrendDTO {
	type row struct {
		Date   string
		Count  int64
		Amount int64
	}

	var rows []row
	s.db.WithContext(ctx).Model(&struct{}{}).Table("orders").
		Select("DATE_FORMAT(created_at, '%m-%d') AS date, COUNT(*) AS count, COALESCE(SUM(total_amount), 0) AS amount").
		Where("created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&rows)

	// 构建完整 7 天序列（补全天数中缺失的日期）
	trend := make([]dto.OrderTrendDTO, 0, 7)
	dateMap := make(map[string]dto.OrderTrendDTO, len(rows))
	for _, r := range rows {
		dateMap[r.Date] = dto.OrderTrendDTO{
			Date:   r.Date,
			Count:  r.Count,
			Amount: r.Amount,
		}
	}

	now := time.Now()
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("01-02")
		if item, ok := dateMap[date]; ok {
			trend = append(trend, item)
		} else {
			trend = append(trend, dto.OrderTrendDTO{Date: date, Count: 0, Amount: 0})
		}
	}
	return trend
}

// computeOrderStatusDist 订单状态分布
func (s *DashboardService) computeOrderStatusDist(ctx context.Context) []dto.StatusDistDTO {
	type row struct {
		Status string
		Value  int64
	}

	var rows []row
	s.db.WithContext(ctx).Model(&struct{}{}).Table("orders").
		Select("status, COUNT(*) AS value").
		Group("status").
		Scan(&rows)

	labelMap := map[string]string{
		"pending":   "待付款",
		"paid":      "已付款",
		"shipped":   "已发货",
		"delivered": "已送达",
		"cancelled": "已取消",
		"refunded":  "已退款",
	}

	result := make([]dto.StatusDistDTO, 0, len(rows))
	for _, r := range rows {
		label := labelMap[r.Status]
		if label == "" {
			label = r.Status
		}
		result = append(result, dto.StatusDistDTO{
			Status: r.Status,
			Label:  label,
			Value:  r.Value,
		})
	}
	return result
}

// computePaymentMethodDist 支付方式分布
func (s *DashboardService) computePaymentMethodDist(ctx context.Context) []dto.MethodDistDTO {
	type row struct {
		Method string `gorm:"column:payment_method"`
		Value  int64
	}

	var rows []row
	s.db.WithContext(ctx).Model(&struct{}{}).Table("payments").
		Select("payment_method, COUNT(*) AS value").
		Where("status IN ?", []string{"success", "paid"}).
		Group("payment_method").
		Scan(&rows)

	labelMap := map[string]string{
		"alipay": "支付宝",
		"wechat": "微信支付",
		"bank":   "银行卡",
		"cash":   "现金",
		"flash":  "秒杀支付",
	}

	result := make([]dto.MethodDistDTO, 0, len(rows))
	for _, r := range rows {
		label := labelMap[r.Method]
		if label == "" {
			label = r.Method
		}
		result = append(result, dto.MethodDistDTO{
			Method: r.Method,
			Label:  label,
			Value:  r.Value,
		})
	}
	return result
}

// computeCategoryDist 商品分类分布
func (s *DashboardService) computeCategoryDist(ctx context.Context) []dto.CategoryDistDTO {
	type row struct {
		Category string
		Value    int64
	}

	var rows []row
	s.db.WithContext(ctx).Model(&struct{}{}).Table("product_categories pc").
		Select("c.name AS category, COUNT(pc.product_id) AS value").
		Joins("JOIN categories c ON c.id = pc.category_id").
		Group("pc.category_id, c.name").
		Order("value DESC").
		Limit(8).
		Scan(&rows)

	result := make([]dto.CategoryDistDTO, len(rows))
	for i, r := range rows {
		result[i] = dto.CategoryDistDTO{
			Category: r.Category,
			Value:    r.Value,
		}
	}
	return result
}

// computeInventoryStatusDist 库存状态分布
func (s *DashboardService) computeInventoryStatusDist(ctx context.Context) []dto.StatusDistDTO {
	type row struct {
		Status string
		Value  int64
	}

	var rows []row
	s.db.WithContext(ctx).Model(&struct{}{}).Table("inventories").
		Select("status, COUNT(*) AS value").
		Group("status").
		Scan(&rows)

	labelMap := map[string]string{
		"instock":    "库存充足",
		"lowstock":   "库存偏低",
		"outofstock": "缺货",
	}

	result := make([]dto.StatusDistDTO, 0, len(rows))
	for _, r := range rows {
		label := labelMap[r.Status]
		if label == "" {
			label = r.Status
		}
		result = append(result, dto.StatusDistDTO{
			Status: r.Status,
			Label:  label,
			Value:  r.Value,
		})
	}
	return result
}

// computeTopProducts 热销商品 Top 10
func (s *DashboardService) computeTopProducts(ctx context.Context) []dto.TopProductDTO {
	type row struct {
		ProductID string
		Count     int64
		Amount    int64
	}

	// 1. 从 order_items 汇总销量
	var rows []row
	s.db.WithContext(ctx).Model(&struct{}{}).Table("order_items oi").
		Select("oi.product_id, COUNT(*) AS count, COALESCE(SUM(oi.amount), 0) AS amount").
		Group("oi.product_id").
		Order("count DESC, amount DESC").
		Limit(10).
		Scan(&rows)

	if len(rows) == 0 {
		return nil
	}

	// 2. 批量查询产品名称
	productIDs := make([]string, len(rows))
	for i, r := range rows {
		productIDs[i] = r.ProductID
	}

	type nameRow struct {
		ID   string
		Name string
	}
	var nameRows []nameRow
	s.db.WithContext(ctx).Model(&struct{}{}).Table("products").
		Select("CAST(id AS CHAR) AS id, name").
		Where("CAST(id AS CHAR) IN ?", productIDs).
		Scan(&nameRows)

	nameMap := make(map[string]string, len(nameRows))
	for _, nr := range nameRows {
		nameMap[nr.ID] = nr.Name
	}

	// 3. 组装结果
	result := make([]dto.TopProductDTO, len(rows))
	for i, r := range rows {
		name, ok := nameMap[r.ProductID]
		if !ok {
			name = fmt.Sprintf("商品#%s", r.ProductID)
		}
		result[i] = dto.TopProductDTO{
			Name:   name,
			Count:  r.Count,
			Amount: r.Amount,
		}
	}
	return result
}

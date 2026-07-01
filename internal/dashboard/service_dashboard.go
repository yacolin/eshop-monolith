package dashboard

import (
	"context"
	"fmt"
	"time"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/pkg/logger"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

const (
	dashboardCacheKey = "dashboard:stats"
	dashboardCacheTTL = 20 * time.Minute
	dashboardRefreshInterval = 15 * time.Minute
)

type DashboardService struct {
	db          *gorm.DB
	rdb         *redis.Client
	repos       *repository.Repositories
	rabbit      *rabbitmq.Client
	singleGroup singleflight.Group
}

func NewDashboardService(db *gorm.DB, rdb *redis.Client, repos *repository.Repositories, rabbit *rabbitmq.Client) *DashboardService {
	return &DashboardService{db: db, rdb: rdb, repos: repos, rabbit: rabbit}
}

func (s *DashboardService) GetStats(ctx context.Context) (*DashboardResponse, error) {
	v, err, _ := s.singleGroup.Do(dashboardCacheKey, func() (interface{}, error) {
		cached, err := s.rdb.Get(ctx, dashboardCacheKey).Bytes()
		if err == nil {
			var resp DashboardResponse
			if sonic.Unmarshal(cached, &resp) == nil {
				return &resp, nil
			}
		}

		resp, err := s.rebuildStats(ctx)
		if err != nil {
			return nil, err
		}

		data, marshalErr := sonic.Marshal(resp)
		if marshalErr == nil {
			s.rdb.Set(ctx, dashboardCacheKey, data, dashboardCacheTTL)
		}
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*DashboardResponse), nil
}

// rebuildStats 并行执行 7 个统计查询
func (s *DashboardService) rebuildStats(ctx context.Context) (*DashboardResponse, error) {
	g, ctx := errgroup.WithContext(ctx)

	var summary SummaryDTO
	g.Go(func() error {
		summary = s.computeSummary(ctx)
		return nil
	})

	var orderTrend []OrderTrendDTO
	g.Go(func() error {
		orderTrend = s.computeOrderTrend(ctx)
		return nil
	})

	var orderStatusDist []StatusDistDTO
	g.Go(func() error {
		orderStatusDist = s.computeOrderStatusDist(ctx)
		return nil
	})

	var paymentMethodDist []MethodDistDTO
	g.Go(func() error {
		paymentMethodDist = s.computePaymentMethodDist(ctx)
		return nil
	})

	var categoryDist []CategoryDistDTO
	g.Go(func() error {
		categoryDist = s.computeCategoryDist(ctx)
		return nil
	})

	var inventoryStatusDist []StatusDistDTO
	g.Go(func() error {
		inventoryStatusDist = s.computeInventoryStatusDist(ctx)
		return nil
	})

	var topProducts []TopProductDTO
	g.Go(func() error {
		topProducts = s.computeTopProducts(ctx)
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &DashboardResponse{
		Summary:             summary,
		OrderTrend:          orderTrend,
		OrderStatusDist:     orderStatusDist,
		PaymentMethodDist:   paymentMethodDist,
		CategoryDist:        categoryDist,
		InventoryStatusDist: inventoryStatusDist,
		TopProducts:         topProducts,
	}, nil
}

func (s *DashboardService) RefreshCache(ctx context.Context) error {
	resp, err := s.rebuildStats(ctx)
	if err != nil {
		return err
	}
	data, marshalErr := sonic.Marshal(resp)
	if marshalErr != nil {
		return marshalErr
	}
	return s.rdb.Set(ctx, dashboardCacheKey, data, dashboardCacheTTL).Err()
}

func (s *DashboardService) InvalidateCache(ctx context.Context) {
	s.rdb.Del(ctx, dashboardCacheKey)
}

// StartPeriodicRefresh 启动后台定时刷新。应在服务初始化时启动一次。
func (s *DashboardService) StartPeriodicRefresh(ctx context.Context) {
	go func() {
		logger.Info("Starting dashboard cache warmup...")
		if err := s.RefreshCache(ctx); err != nil {
			logger.Error("Dashboard cache warmup failed", "error", err)
		} else {
			logger.Info("Dashboard cache warmup completed")
		}

		ticker := time.NewTicker(dashboardRefreshInterval)
		defer ticker.Stop()
		for range ticker.C {
			if err := s.RefreshCache(ctx); err != nil {
				logger.Error("Dashboard cache refresh failed", "error", err)
			}
		}
	}()
}

func (s *DashboardService) computeSummary(ctx context.Context) SummaryDTO {
	var summary SummaryDTO
	s.db.WithContext(ctx).Table("tx_orders").Select("COUNT(*)").Scan(&summary.TotalOrders)
	s.db.WithContext(ctx).Table("tx_payments").Select("COALESCE(SUM(amount), 0)").Where("status IN ?", []string{"success", "paid"}).Scan(&summary.TotalRevenue)
	s.db.WithContext(ctx).Table("sp_products").Select("COUNT(*)").Scan(&summary.TotalProducts)
	s.db.WithContext(ctx).Table("sp_inventories").Select("COUNT(*)").Where("status IN ?", []string{"lowstock", "outofstock"}).Scan(&summary.LowStockCount)
	return summary
}

func (s *DashboardService) computeOrderTrend(ctx context.Context) []OrderTrendDTO {
	type row struct {
		Date   string
		Count  int64
		Amount int64
	}
	var rows []row
	s.db.WithContext(ctx).Table("tx_orders").
		Select("DATE_FORMAT(created_at, '%m-%d') AS date, COUNT(*) AS count, COALESCE(SUM(total_amount), 0) AS amount").
		Where("created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)").
		Group("DATE_FORMAT(created_at, '%m-%d')").Order("date ASC").Scan(&rows)

	trend := make([]OrderTrendDTO, 0, 7)
	dateMap := make(map[string]OrderTrendDTO, len(rows))
	for _, r := range rows {
		dateMap[r.Date] = OrderTrendDTO{Date: r.Date, Count: r.Count, Amount: r.Amount}
	}
	now := time.Now()
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("01-02")
		if item, ok := dateMap[date]; ok {
			trend = append(trend, item)
		} else {
			trend = append(trend, OrderTrendDTO{Date: date, Count: 0, Amount: 0})
		}
	}
	return trend
}

func (s *DashboardService) computeOrderStatusDist(ctx context.Context) []StatusDistDTO {
	type row struct{ Status string; Value int64 }
	var rows []row
	s.db.WithContext(ctx).Table("tx_orders").Select("status, COUNT(*) AS value").Group("status").Scan(&rows)
	labelMap := map[string]string{"pending": "待付款", "paid": "已付款", "shipped": "已发货", "delivered": "已送达", "cancelled": "已取消", "refunded": "已退款"}
	result := make([]StatusDistDTO, 0, len(rows))
	for _, r := range rows {
		label := labelMap[r.Status]
		if label == "" {
			label = r.Status
		}
		result = append(result, StatusDistDTO{Status: r.Status, Label: label, Value: r.Value})
	}
	return result
}

func (s *DashboardService) computePaymentMethodDist(ctx context.Context) []MethodDistDTO {
	type row struct{ Method string `gorm:"column:payment_method"`; Value int64 }
	var rows []row
	s.db.WithContext(ctx).Table("tx_payments").Select("payment_method, COUNT(*) AS value").Where("status IN ?", []string{"success", "paid"}).Group("payment_method").Scan(&rows)
	labelMap := map[string]string{"alipay": "支付宝", "wechat": "微信支付", "wallet": "余额", "bank": "银行卡", "cash": "现金"}
	result := make([]MethodDistDTO, 0, len(rows))
	for _, r := range rows {
		label := labelMap[r.Method]
		if label == "" {
			label = r.Method
		}
		result = append(result, MethodDistDTO{Method: r.Method, Label: label, Value: r.Value})
	}
	return result
}

func (s *DashboardService) computeCategoryDist(ctx context.Context) []CategoryDistDTO {
	type row struct{ Category string; Value int64 }
	var rows []row
	s.db.WithContext(ctx).Table("sp_products p").
		Select("COALESCE(r.name, '未分类') AS category, COUNT(p.id) AS value").
		Joins("LEFT JOIN sp_categories c ON c.id = p.category_id").
		Joins("LEFT JOIN sp_categories r ON r.level = 1 AND r.id = COALESCE(NULLIF(SUBSTRING_INDEX(c.path, '/', 1), ''), c.id)").
		Where("p.deleted_at IS NULL").
		Group("r.id").Order("value DESC").Limit(8).Scan(&rows)
	result := make([]CategoryDistDTO, len(rows))
	for i, r := range rows {
		result[i] = CategoryDistDTO{Category: r.Category, Value: r.Value}
	}
	return result
}

func (s *DashboardService) computeInventoryStatusDist(ctx context.Context) []StatusDistDTO {
	type row struct{ Status string; Value int64 }
	var rows []row
	s.db.WithContext(ctx).Table("sp_inventories").Select("status, COUNT(*) AS value").Group("status").Scan(&rows)
	labelMap := map[string]string{"instock": "库存充足", "lowstock": "库存偏低", "outofstock": "缺货"}
	result := make([]StatusDistDTO, 0, len(rows))
	for _, r := range rows {
		label := labelMap[r.Status]
		if label == "" {
			label = r.Status
		}
		result = append(result, StatusDistDTO{Status: r.Status, Label: label, Value: r.Value})
	}
	return result
}

func (s *DashboardService) computeTopProducts(ctx context.Context) []TopProductDTO {
	type row struct{ ProductID string; Count int64; Amount int64 }
	var rows []row
	s.db.WithContext(ctx).Table("tx_order_items oi").
		Select("oi.product_id, COUNT(*) AS count, COALESCE(SUM(oi.subtotal), 0) AS amount").
		Group("oi.product_id").Order("count DESC, amount DESC").Limit(10).Scan(&rows)
	if len(rows) == 0 {
		return nil
	}
	productIDs := make([]string, len(rows))
	for i, r := range rows {
		productIDs[i] = r.ProductID
	}
	type nameRow struct{ ID string; Name string }
	var nameRows []nameRow
	s.db.WithContext(ctx).Table("sp_products").Select("CAST(id AS CHAR) AS id, name").Where("CAST(id AS CHAR) IN ?", productIDs).Scan(&nameRows)
	nameMap := make(map[string]string, len(nameRows))
	for _, nr := range nameRows {
		nameMap[nr.ID] = nr.Name
	}
	result := make([]TopProductDTO, len(rows))
	for i, r := range rows {
		name, ok := nameMap[r.ProductID]
		if !ok {
			name = fmt.Sprintf("商品#%s", r.ProductID)
		}
		result[i] = TopProductDTO{Name: name, Count: r.Count, Amount: r.Amount}
	}
	return result
}

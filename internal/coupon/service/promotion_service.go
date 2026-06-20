package service

import (
	"context"
	"encoding/json"
	"time"

	"eshop-monolith/internal/coupon/domain/models"
	"eshop-monolith/internal/coupon/domain/repositories"
	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

// FullReduceRule 满减活动规则
type FullReduceRule struct {
	Threshold int64 `json:"threshold"` // 满多少（分）
	Reduce    int64 `json:"reduce"`    // 减多少（分）
}

// TimeDiscountRule 限时折扣规则
type TimeDiscountRule struct {
	DiscountRate int64 `json:"discount_rate"` // 折扣率（百分比*100，如 8000 = 8折）
}

// PromotionService 促销活动服务
type PromotionService struct {
	db               *gorm.DB
	promotionRepo    repositories.IPromotionRepository
	promotionProdRepo repositories.IPromotionProductRepository
	bus              *eventbus.Bus
}

func NewPromotionService(db *gorm.DB, promotionRepo repositories.IPromotionRepository, promotionProdRepo repositories.IPromotionProductRepository, bus *eventbus.Bus) *PromotionService {
	return &PromotionService{
		db:               db,
		promotionRepo:    promotionRepo,
		promotionProdRepo: promotionProdRepo,
		bus:              bus,
	}
}

// CreatePromotion 创建促销活动
func (s *PromotionService) CreatePromotion(ctx context.Context, promotion *models.Promotion) error {
	// 默认状态为待开始
	if promotion.Status == "" {
		promotion.Status = models.PromotionStatusPending
	}
	return s.promotionRepo.Create(ctx, promotion)
}

// UpdatePromotion 更新促销活动
func (s *PromotionService) UpdatePromotion(ctx context.Context, promotion *models.Promotion) error {
	existing, err := s.promotionRepo.FindByID(ctx, promotion.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errcode.ErrNotFound
	}
	if existing.Status == models.PromotionStatusFinished || existing.Status == models.PromotionStatusCancelled {
		return errcode.ErrInvalidParams
	}
	return s.promotionRepo.Update(ctx, promotion)
}

// GetPromotion 获取促销活动
func (s *PromotionService) GetPromotion(ctx context.Context, id int64) (*models.Promotion, error) {
	promotion, err := s.promotionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if promotion == nil {
		return nil, errcode.ErrNotFound
	}
	return promotion, nil
}

// ListPromotions 促销活动列表
func (s *PromotionService) ListPromotions(ctx context.Context, page, pageSize int) ([]models.Promotion, int64, error) {
	offset := (page - 1) * pageSize
	return s.promotionRepo.List(ctx, offset, pageSize)
}

// GetActivePromotions 获取当前活动中的促销
func (s *PromotionService) GetActivePromotions(ctx context.Context) ([]models.Promotion, error) {
	return s.promotionRepo.FindActive(ctx)
}

// UpdatePromotionStatus 更新促销活动状态
func (s *PromotionService) UpdatePromotionStatus(ctx context.Context, id int64, status models.PromotionStatus) error {
	return s.promotionRepo.UpdateStatus(ctx, id, status)
}

// LinkProducts 关联促销商品
func (s *PromotionService) LinkProducts(ctx context.Context, promotionID int64, productIDs []int64, discount int64) error {
	// 先删除旧的关联
	if err := s.promotionProdRepo.DeleteByPromotionID(ctx, promotionID); err != nil {
		return err
	}
	// 创建新的关联
	items := make([]models.PromotionProduct, len(productIDs))
	for i, pid := range productIDs {
		items[i] = models.PromotionProduct{
			PromotionID: promotionID,
			ProductID:   pid,
			Discount:    discount,
		}
	}
	return s.promotionProdRepo.CreateBatch(ctx, items)
}

// GetPromotionProducts 获取促销关联商品
func (s *PromotionService) GetPromotionProducts(ctx context.Context, promotionID int64) ([]models.PromotionProduct, error) {
	return s.promotionProdRepo.FindByPromotionID(ctx, promotionID)
}

// GetProductPromotion 获取商品当前参与的促销活动（含计算后的价格）
func (s *PromotionService) GetProductPromotion(ctx context.Context, productID int64) (*models.Promotion, *models.PromotionProduct, error) {
	pp, err := s.promotionProdRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, nil, err
	}
	if len(pp) == 0 {
		return nil, nil, nil
	}

	now := time.Now()
	for _, p := range pp {
		promotion, err := s.promotionRepo.FindByID(ctx, p.PromotionID)
		if err != nil {
			continue
		}
		if promotion == nil {
			continue
		}
		// 检查活动是否在进行中
		if promotion.Status == models.PromotionStatusActive &&
			now.After(time.Time(promotion.StartTime)) &&
			now.Before(time.Time(promotion.EndTime)) {
			return promotion, &p, nil
		}
	}

	return nil, nil, nil
}

// CalculateDiscount 计算订单的促销优惠（返回总优惠金额，单位：分）
func (s *PromotionService) CalculateDiscount(ctx context.Context, items []PromotionItem) (int64, error) {
	actives, err := s.promotionRepo.FindActive(ctx)
	if err != nil {
		return 0, err
	}

	var totalDiscount int64
	for _, promo := range actives {
		switch promo.PromoType {
		case models.PromotionTypeTimeDiscount:
			discount, err := s.calcTimeDiscount(&promo, items)
			if err != nil {
				continue
			}
			totalDiscount += discount

		case models.PromotionTypeFullReduce:
			discount, err := s.calcFullReduce(&promo, items)
			if err != nil {
				continue
			}
			totalDiscount += discount
		}
	}

	return totalDiscount, nil
}

// PromotionItem 订单中参与促销的商品项
type PromotionItem struct {
	ProductID int64
	Price     int64 // 单价（分）
	Quantity  int
}

func (s *PromotionService) calcTimeDiscount(promo *models.Promotion, items []PromotionItem) (int64, error) {
	var rule TimeDiscountRule
	if err := json.Unmarshal([]byte(promo.Rule), &rule); err != nil {
		return 0, err
	}

	// 获取该促销关联的商品ID集合
	pp, err := s.promotionProdRepo.FindByPromotionID(context.TODO(), promo.ID)
	if err != nil {
		return 0, err
	}
	productMap := make(map[int64]int64)
	for _, p := range pp {
		productMap[p.ProductID] = p.Discount
	}

	var discount int64
	for _, item := range items {
		// 检查商品是否在促销范围内
		switch promo.Scope {
		case "all":
			// 全场参与
		case "product":
			if _, ok := productMap[item.ProductID]; !ok {
				continue
			}
		default:
			continue
		}

		amount := item.Price * int64(item.Quantity)
		itemDiscount := amount * rule.DiscountRate / 10000
		discount += amount - itemDiscount
	}

	return discount, nil
}

func (s *PromotionService) calcFullReduce(promo *models.Promotion, items []PromotionItem) (int64, error) {
	var rule FullReduceRule
	if err := json.Unmarshal([]byte(promo.Rule), &rule); err != nil {
		return 0, err
	}

	var totalAmount int64
	for _, item := range items {
		// 检查商品是否在促销范围内
		switch promo.Scope {
		case "all":
			totalAmount += item.Price * int64(item.Quantity)
		case "product":
			pp, err := s.promotionProdRepo.FindByProductID(context.TODO(), item.ProductID)
			if err != nil {
				continue
			}
			for _, p := range pp {
				if p.PromotionID == promo.ID {
					totalAmount += item.Price * int64(item.Quantity)
					break
				}
			}
		}
	}

	if totalAmount >= rule.Threshold {
		return rule.Reduce, nil
	}
	return 0, nil
}

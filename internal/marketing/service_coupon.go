package marketing

import (
	"context"
	"errors"
	"time"

	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

type CouponService struct {
	repo IpromotionRepository
	db   *gorm.DB
}

func NewCouponService(repo IpromotionRepository, db *gorm.DB) *CouponService {
	return &CouponService{repo: repo, db: db}
}

type UserPromotionListResult struct {
	Total int64           `json:"total"`
	List  []*UserPromotion `json:"list"`
}

// Claim 领取优惠券
func (s *CouponService) Claim(ctx context.Context, userID int64, req *ClaimCouponReq) (*UserPromotion, error) {
	promo, err := s.repo.FindByID(ctx, req.PromotionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPromotionNotFound
		}
		return nil, err
	}
	// 校验是优惠券类型
	if promo.PromoType != 1 && promo.PromoType != 2 {
		return nil, errcode.ErrPromotionRuleInvalid
	}
	// 检查状态
	if promo.Status != 2 {
		return nil, errcode.ErrPromotionRuleInvalid
	}
	// 检查时间
	now := time.Now()
	if now.Before(promo.StartTime) || now.After(promo.EndTime) {
		return nil, errcode.ErrCouponExpired
	}
	// 检查库存
	if promo.TotalQuantity > 0 && promo.UsedQuantity >= promo.TotalQuantity {
		return nil, errcode.ErrCouponSoldOut
	}
	// 检查用户已领取
	existing, err := s.repo.FindUserPromotion(ctx, userID, req.PromotionID)
	if err == nil && existing != nil {
		return nil, errcode.ErrCouponAlreadyClaimed
	}

	up := &UserPromotion{
		UserID:      userID,
		PromotionID: req.PromotionID,
		Status:      1,
	}
	if promo.PerUserLimit > 0 {
		up.ExpireTime = timePtr(time.Now().AddDate(0, 0, promo.PerUserLimit))
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(up).Error; err != nil {
			return err
		}
		// 增加已领数量
		return tx.Model(&Promotion{}).Where("id = ?", req.PromotionID).
			Update("used_quantity", gorm.Expr("used_quantity + 1")).Error
	})
	return up, err
}

// Use 使用优惠券
func (s *CouponService) Use(ctx context.Context, userID int64, req *UseCouponReq) error {
	// 直接用 db 更新
	result := s.db.Model(&UserPromotion{}).Where("id = ? AND user_id = ? AND status = 1", req.UserPromotionID, userID).
		Updates(map[string]interface{}{
			"status":   2,
			"used_time": time.Now(),
			"order_id": req.OrderID,
		})
	if result.RowsAffected == 0 {
		return errcode.ErrCouponExpired
	}
	return result.Error
}

func (s *CouponService) ListUserCoupons(ctx context.Context, userID int64, req *UserPromotionListReq) (*UserPromotionListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	list, total, err := s.repo.ListUserPromotions(ctx, userID, req.Status, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*UserPromotion, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return &UserPromotionListResult{Total: total, List: items}, nil
}

func timePtr(t time.Time) *time.Time { return &t }

package promotion

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

type FlashService struct {
	repo IpromotionRepository
	db   *gorm.DB
}

func NewFlashService(repo IpromotionRepository, db *gorm.DB) *FlashService {
	return &FlashService{repo: repo, db: db}
}

// Buy 秒杀抢购（生成排队令牌）
func (s *FlashService) Buy(ctx context.Context, userID int64, req *FlashBuyReq) (*UserPromotion, error) {
	promo, err := s.repo.FindByID(ctx, req.PromotionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPromotionNotFound
		}
		return nil, err
	}
	if promo.PromoType != 3 {
		return nil, errcode.ErrPromotionRuleInvalid
	}
	if promo.Status != 2 {
		return nil, errcode.ErrPromotionRuleInvalid
	}
	now := time.Now()
	if now.Before(promo.StartTime) || now.After(promo.EndTime) {
		return nil, errcode.ErrCouponExpired
	}
	if promo.TotalQuantity > 0 && promo.UsedQuantity >= promo.TotalQuantity {
		return nil, errcode.ErrCouponSoldOut
	}

	// 生成本地 token
	tokenBytes := make([]byte, 16)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	up := &UserPromotion{
		UserID:      userID,
		PromotionID: req.PromotionID,
		Status:      1,
		QueueToken:  token,
		AcquireTime: now,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 行锁扣减秒杀库存
		var p Promotion
		if err := tx.Where("id = ? AND promo_type = 3", req.PromotionID).
			Session(&gorm.Session{}).First(&p).Error; err != nil {
			return err
		}
		if p.TotalQuantity > 0 && p.UsedQuantity >= p.TotalQuantity {
			return errcode.ErrCouponSoldOut
		}
		if err := tx.Model(&p).Update("used_quantity", gorm.Expr("used_quantity + 1")).Error; err != nil {
			return err
		}
		return tx.Create(up).Error
	})
	return up, err
}

// Confirm 确认秒杀订单
func (s *FlashService) Confirm(ctx context.Context, userID int64, req *FlashConfirmReq) (*UserPromotion, error) {
	var up UserPromotion
	err := s.db.Where("queue_token = ? AND user_id = ? AND status = 1", req.Token, userID).First(&up).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPromotionNotFound
		}
		return nil, err
	}
	// 标记为已使用
	s.db.Model(&up).Updates(map[string]interface{}{
		"status":   2,
		"used_time": time.Now(),
	})
	return &up, nil
}

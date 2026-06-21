package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"eshop-monolith/internal/coupon/domain/models"
	"eshop-monolith/internal/coupon/domain/repositories"
	"eshop-monolith/internal/coupon/events"
	"eshop-monolith/internal/infra/eventbus"

	"gorm.io/gorm"
)

// CouponService 优惠券服务
type CouponService struct {
	db             *gorm.DB
	couponRepo     repositories.IcouponRepository
	userCouponRepo repositories.IuserCouponRepository
	bus            *eventbus.Bus
}

func NewCouponService(db *gorm.DB, couponRepo repositories.IcouponRepository, userCouponRepo repositories.IuserCouponRepository, bus *eventbus.Bus) *CouponService {
	return &CouponService{
		db:             db,
		couponRepo:     couponRepo,
		userCouponRepo: userCouponRepo,
		bus:            bus,
	}
}

// CreateCoupon 创建优惠券模板
func (s *CouponService) CreateCoupon(ctx context.Context, coupon *models.Coupon) error {
	coupon.RemainStock = coupon.TotalStock
	coupon.Status = models.CouponStatusActive
	return s.couponRepo.Create(ctx, coupon)
}

// UpdateCoupon 更新优惠券模板
func (s *CouponService) UpdateCoupon(ctx context.Context, coupon *models.Coupon) error {
	existing, err := s.couponRepo.FindByID(ctx, coupon.ID)
	if err != nil {
		return err
	}
	// 保留原有剩余数量
	coupon.RemainStock = existing.RemainStock
	return s.couponRepo.Update(ctx, coupon)
}

// GetCoupon 获取优惠券模板
func (s *CouponService) GetCoupon(ctx context.Context, id int64) (*models.Coupon, error) {
	coupon, err := s.couponRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return coupon, nil
}

// ListCoupons 优惠券模板列表
func (s *CouponService) ListCoupons(ctx context.Context, page, pageSize int) ([]models.Coupon, int64, error) {
	offset := (page - 1) * pageSize
	return s.couponRepo.List(ctx, offset, pageSize)
}

// ClaimCoupon 用户领取优惠券
func (s *CouponService) ClaimCoupon(ctx context.Context, userID int64, couponID int64) (*models.UserCoupon, error) {
	var userCoupon *models.UserCoupon

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 查询优惠券模板
		coupon, err := s.couponRepo.FindByIDWithTx(tx, couponID)
		if err != nil {
			return err
		}
		if !coupon.IsValid() {
			return errors.New("coupon is not available")
		}

		// 2. 检查领取时间
		now := time.Now()
		if now.Before(time.Time(coupon.StartTime)) || now.After(time.Time(coupon.EndTime)) {
			return errors.New("coupon claim period has not started or has ended")
		}

		// 3. 检查用户领取限制
		if coupon.UserLimit > 0 {
			claimed, err := s.userCouponRepo.FindByUserAndCoupon(ctx, userID, couponID)
			if err != nil {
				return err
			}
			if len(claimed) >= coupon.UserLimit {
				return errors.New("coupon claim limit reached")
			}
		}

		// 4. 扣减库存
		rowsAffected := tx.Model(&models.Coupon{}).Where("id = ? AND remain_stock > 0", couponID).
			UpdateColumn("remain_stock", gorm.Expr("remain_stock - 1")).RowsAffected
		if rowsAffected == 0 {
			return errors.New("coupon stock exhausted")
		}

		// 5. 计算过期时间
		var expireAt time.Time
		if coupon.ValidDays > 0 {
			expireAt = now.AddDate(0, 0, coupon.ValidDays)
		} else {
			expireAt = time.Time(coupon.EndTime)
		}

		// 6. 生成优惠券码
		code, err := generateCouponCode()
		if err != nil {
			return err
		}

		// 7. 创建用户优惠券
		userCoupon = &models.UserCoupon{
			UserID:     userID,
			CouponID:   couponID,
			CouponCode: code,
			Status:     models.UserCouponStatusUnused,
			ExpireAt:   expireAt,
		}
		return s.userCouponRepo.CreateWithTx(tx, userCoupon)
	})

	if err != nil {
		return nil, err
	}

	// 发布事件
	if s.bus != nil && userCoupon != nil {
		s.bus.Publish(events.CouponIssuedEvent{
			UserCouponID: userCoupon.ID,
			UserID:       userCoupon.UserID,
			CouponID:     userCoupon.CouponID,
			CouponCode:   userCoupon.CouponCode,
			IssuedAt:     time.Now(),
		})
	}

	return userCoupon, nil
}

// UseCoupon 使用优惠券（订单结算时调用）
func (s *CouponService) UseCoupon(ctx context.Context, userCouponID int64, userID int64, orderNo string, orderAmount int64) (int64, error) {
	var discount int64

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 查询用户优惠券
		uc, err := s.userCouponRepo.FindByIDWithTx(tx, userCouponID)
		if err != nil {
			return err
		}
		if uc.UserID != userID {
			return errors.New("coupon does not belong to this user")
		}
		if uc.Status != models.UserCouponStatusUnused {
			return errors.New("coupon has been used or expired")
		}
		if time.Now().After(uc.ExpireAt) {
			return errors.New("coupon has expired")
		}

		// 2. 查询优惠券模板
		coupon, err := s.couponRepo.FindByIDWithTx(tx, uc.CouponID)
		if err != nil {
			return err
		}

		// 3. 计算实际抵扣金额
		discount, err = calculateDiscount(coupon, orderAmount)
		if err != nil {
			return err
		}

		// 4. 更新用户优惠券状态
		return s.userCouponRepo.UpdateStatusWithTx(tx, userCouponID, models.UserCouponStatusUsed, orderNo)
	})

	if err != nil {
		return 0, err
	}

	// 发布事件
	if s.bus != nil {
		s.bus.Publish(events.CouponUsedEvent{
			UserCouponID: userCouponID,
			UserID:       userID,
			OrderNo:      orderNo,
			DiscountAmt:  discount,
			UsedAt:       time.Now(),
		})
	}

	return discount, nil
}

// GetUserCoupons 获取用户优惠券列表
func (s *CouponService) GetUserCoupons(ctx context.Context, userID int64, status models.UserCouponStatus, page, pageSize int) ([]models.UserCoupon, int64, error) {
	offset := (page - 1) * pageSize
	return s.userCouponRepo.FindByUserID(ctx, userID, status, offset, pageSize)
}

// GetUsableCoupons 获取用户可用优惠券
func (s *CouponService) GetUsableCoupons(ctx context.Context, userID int64) ([]models.UserCoupon, error) {
	return s.userCouponRepo.FindUsableByUserID(ctx, userID)
}

// ValidateCoupon 校验优惠券是否可用（结算时预校验）
func (s *CouponService) ValidateCoupon(ctx context.Context, userCouponID int64, userID int64, orderAmount int64) (int64, error) {
	uc, err := s.userCouponRepo.FindByID(ctx, userCouponID)
	if err != nil {
		return 0, err
	}
	if uc.UserID != userID {
		return 0, errors.New("coupon does not belong to this user")
	}
	if uc.Status != models.UserCouponStatusUnused {
		return 0, errors.New("coupon has been used or expired")
	}
	if time.Now().After(uc.ExpireAt) {
		return 0, errors.New("coupon has expired")
	}

	coupon, err := s.couponRepo.FindByID(ctx, uc.CouponID)
	if err != nil {
		return 0, err
	}

	return calculateDiscount(coupon, orderAmount)
}

// calculateDiscount 计算优惠金额（单位：分）
func calculateDiscount(coupon *models.Coupon, orderAmount int64) (int64, error) {
	// 检查最低消费
	if coupon.MinAmount > 0 && orderAmount < coupon.MinAmount {
		return 0, fmt.Errorf("order amount %d less than minimum %d", orderAmount, coupon.MinAmount)
	}

	switch coupon.CouponType {
	case models.CouponTypeFixed:
		// 满减券：直接减固定金额
		return coupon.Value, nil

	case models.CouponTypePercentage:
		// 折扣券：按百分比打折
		discount := orderAmount * coupon.Value / 10000
		if coupon.MaxDiscount > 0 && discount > coupon.MaxDiscount {
			discount = coupon.MaxDiscount
		}
		return discount, nil

	case models.CouponTypeVoucher:
		// 代金券：无门槛直接减
		return coupon.Value, nil

	default:
		return 0, errors.New("unknown coupon type")
	}
}

// generateCouponCode 生成唯一优惠券码
func generateCouponCode() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	n := new(big.Int).SetBytes(b)
	code := n.Mod(n, big.NewInt(10000000000)).Int64()
	return fmt.Sprintf("%010d", code), nil
}

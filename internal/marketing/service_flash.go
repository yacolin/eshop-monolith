package marketing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"eshop-monolith/pkg/errcode"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	flashStockRedisPrefix = "flash:stock:"
	flashUserRedisPrefix  = "flash:users:"
)

// flashBuyLua 原子扣库存：1=成功, 0=售罄, -1=已抢过
var flashBuyLua = redis.NewScript(`
local bought = redis.call("SISMEMBER", KEYS[2], ARGV[1])
if bought == 1 then
    return -1
end

local total = redis.call("HGET", KEYS[1], "total")
local sold = redis.call("HGET", KEYS[1], "sold")

if total and tonumber(total) > 0 then
    if not sold then sold = 0 end
    if tonumber(sold) >= tonumber(total) then
        return 0
    end
end

redis.call("HINCRBY", KEYS[1], "sold", 1)
redis.call("SADD", KEYS[2], ARGV[1])
return 1
`)

type FlashBuyJob struct {
	UserID      int64
	PromotionID int64
	QueueToken  string
	AcquireTime time.Time
}

type FlashService struct {
	repo    IpromotionRepository
	db      *gorm.DB
	rdb     *redis.Client
	jobChan chan FlashBuyJob
}

func NewFlashService(repo IpromotionRepository, db *gorm.DB, rdb *redis.Client) *FlashService {
	svc := &FlashService{
		repo:    repo,
		db:      db,
		rdb:     rdb,
		jobChan: make(chan FlashBuyJob, 1024),
	}
	go svc.asyncWorker()
	return svc
}

func (s *FlashService) asyncWorker() {
	for job := range s.jobChan {
		for i := 0; i < 3; i++ {
			up := &UserPromotion{
				UserID:      job.UserID,
				PromotionID: job.PromotionID,
				Status:      1,
				QueueToken:  job.QueueToken,
				AcquireTime: job.AcquireTime,
			}
			err := s.db.Create(up).Error
			if err == nil {
				break
			}
			if i < 2 {
				time.Sleep(time.Duration(50*(i+1)) * time.Millisecond)
			}
		}
	}
}

// Close 优雅关闭，等待异步写入完成
func (s *FlashService) Close() {
	close(s.jobChan)
}

// Buy 秒杀抢购 — Redis Lua 原子扣库存 + 异步落库
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

	// 生成排队令牌
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	stockKey := flashStockRedisPrefix + strconv.FormatInt(req.PromotionID, 10)
	userKey := flashUserRedisPrefix + strconv.FormatInt(req.PromotionID, 10)

	result, err := flashBuyLua.Run(ctx, s.rdb, []string{stockKey, userKey}, strconv.FormatInt(userID, 10)).Int()
	if err != nil {
		return nil, errcode.ErrPromotionRuleInvalid
	}

	switch result {
	case -1:
		return nil, errcode.ErrCouponAlreadyClaimed
	case 0:
		return nil, errcode.ErrCouponSoldOut
	}

	// 异步落库（channel 满时降级同步写，防止积压丢失）
	job := FlashBuyJob{
		UserID:      userID,
		PromotionID: req.PromotionID,
		QueueToken:  token,
		AcquireTime: now,
	}
	select {
	case s.jobChan <- job:
	default:
		s.db.Create(&UserPromotion{
			UserID:      userID,
			PromotionID: req.PromotionID,
			Status:      1,
			QueueToken:  token,
			AcquireTime: now,
		})
	}

	return &UserPromotion{
		UserID:      userID,
		PromotionID: req.PromotionID,
		Status:      1,
		QueueToken:  token,
		AcquireTime: now,
	}, nil
}

// Confirm 确认秒杀订单（DB 中已有异步写入的令牌记录）
func (s *FlashService) Confirm(ctx context.Context, userID int64, req *FlashConfirmReq) (*UserPromotion, error) {
	var up UserPromotion
	err := s.db.WithContext(ctx).
		Where("queue_token = ? AND user_id = ? AND status = 1", req.Token, userID).
		First(&up).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPromotionNotFound
		}
		return nil, err
	}
	s.db.Model(&up).Updates(map[string]interface{}{
		"status":    2,
		"used_time": time.Now(),
	})
	return &up, nil
}

// LoadStockToRedis 预热秒杀库存到 Redis
func (s *FlashService) LoadStockToRedis(ctx context.Context, promotionID int64) error {
	promo, err := s.repo.FindByID(ctx, promotionID)
	if err != nil {
		return err
	}
	stockKey := flashStockRedisPrefix + strconv.FormatInt(promotionID, 10)
	return s.rdb.HSet(ctx, stockKey, "total", promo.TotalQuantity, "sold", promo.UsedQuantity).Err()
}

// SyncStockToDB 活动结束后将 Redis 销量回写到 DB
func (s *FlashService) SyncStockToDB(ctx context.Context, promotionID int64) error {
	stockKey := flashStockRedisPrefix + strconv.FormatInt(promotionID, 10)
	sold, err := s.rdb.HGet(ctx, stockKey, "sold").Int()
	if err != nil {
		return err
	}
	return s.db.Model(&Promotion{}).Where("id = ?", promotionID).
		Update("used_quantity", sold).Error
}

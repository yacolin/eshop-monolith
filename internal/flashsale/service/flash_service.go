package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"eshop-monolith/internal/flashsale/api/dto"
	"eshop-monolith/internal/flashsale/domain/models"
	"eshop-monolith/internal/flashsale/domain/repositories"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var stockLuaScript = redis.NewScript(`
local stock_key = KEYS[1]
local stock = redis.call("GET", stock_key)
if not stock then
    return 0
end
stock = tonumber(stock)
if stock <= 0 then
    return -1
end
redis.call("DECR", stock_key)
return 1
`)

type FlashService struct {
	db    *gorm.DB
	rdb   *redis.Client
	repo  *repositories.FlashRepository
}

func NewFlashService(db *gorm.DB, rdb *redis.Client, repo *repositories.FlashRepository) *FlashService {
	return &FlashService{db: db, rdb: rdb, repo: repo}
}

func stockKey(activityID int64) string {
	return fmt.Sprintf("flash:stock:%d", activityID)
}

func userLimitKey(activityID, userID int64) string {
	return fmt.Sprintf("flash:limit:%d:%d", activityID, userID)
}

func (s *FlashService) CreateActivity(ctx context.Context, req *dto.CreateActivityReq) (*models.FlashActivity, error) {
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, time.Local)
	if err != nil {
		return nil, errors.New("invalid start_time format, use 2006-01-02 15:04:05")
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, time.Local)
	if err != nil {
		return nil, errors.New("invalid end_time format, use 2006-01-02 15:04:05")
	}
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return nil, errors.New("end_time must be after start_time")
	}

	now := time.Now()
	startUnix := startTime.Unix()
	endUnix := endTime.Unix()
	nowUnix := now.Unix()

	status := string(models.FlashStatusPending)
	if nowUnix >= startUnix && nowUnix < endUnix {
		status = string(models.FlashStatusActive)
	} else if nowUnix >= endUnix {
		status = string(models.FlashStatusFinished)
	}

	activity := &models.FlashActivity{
		ProductID:  req.ProductID,
		FlashPrice: req.FlashPrice,
		TotalStock: req.TotalStock,
		SoldStock:  0,
		StartTime:  utils.Timestamp(startTime),
		EndTime:    utils.Timestamp(endTime),
		Status:     status,
	}

	if err := s.repo.CreateActivity(ctx, activity); err != nil {
		return nil, err
	}
	return activity, nil
}

func (s *FlashService) LoadStockToRedis(ctx context.Context, activityID int64) error {
	activity, err := s.repo.GetActivity(ctx, activityID)
	if err != nil {
		return errcode.ErrNotFound
	}

	return s.rdb.Set(ctx, stockKey(activityID), activity.TotalStock-activity.SoldStock, 0).Err()
}

func (s *FlashService) GetActivity(ctx context.Context, id int64) (*models.FlashActivity, error) {
	return s.repo.GetActivity(ctx, id)
}

func (s *FlashService) ListActivities(ctx context.Context) ([]models.FlashActivity, error) {
	return s.repo.ListActivities(ctx)
}

func (s *FlashService) GetOrder(ctx context.Context, orderID int64) (*models.FlashOrder, error) {
	return s.repo.GetOrder(ctx, orderID)
}

func (s *FlashService) GetUserOrders(ctx context.Context, userID int64, activityID int64) ([]models.FlashOrder, error) {
	return s.repo.GetUserOrders(ctx, userID, activityID)
}

func (s *FlashService) FlashBuy(ctx context.Context, req *dto.FlashBuyReq) (*dto.FlashBuyResp, error) {
	activity, err := s.repo.GetActivity(ctx, req.ActivityID)
	if err != nil {
		return &dto.FlashBuyResp{Success: false, Message: "activity not found"}, nil
	}

	now := time.Now()
	if now.Before(time.Time(activity.StartTime)) {
		return &dto.FlashBuyResp{Success: false, Message: "flash sale not started"}, nil
	}
	if now.After(time.Time(activity.EndTime)) {
		return &dto.FlashBuyResp{Success: false, Message: "flash sale ended"}, nil
	}

	if activity.Status != string(models.FlashStatusActive) {
		return &dto.FlashBuyResp{Success: false, Message: "flash sale is not active"}, nil
	}

	limitKey := userLimitKey(req.ActivityID, req.UserID)
	exists, _ := s.rdb.Exists(ctx, limitKey).Result()
	if exists > 0 {
		return &dto.FlashBuyResp{Success: false, Message: "already purchased"}, nil
	}

	stockKey := stockKey(req.ActivityID)
	result, err := stockLuaScript.Run(ctx, s.rdb, []string{stockKey}).Result()
	if err != nil {
		return &dto.FlashBuyResp{Success: false, Message: "system error, please try again"}, nil
	}

	code := result.(int64)
	if code <= 0 {
		return &dto.FlashBuyResp{Success: false, Message: "sold out"}, nil
	}

	order := &models.FlashOrder{
		ActivityID:  req.ActivityID,
		UserID:      req.UserID,
		ProductID:   activity.ProductID,
		Quantity:    1,
		FlashPrice:  activity.FlashPrice,
		TotalAmount: activity.FlashPrice,
		Status:      string(models.FlashOrderStatusPending),
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.IncrementSoldStockWithTx(tx, req.ActivityID, 1); err != nil {
			return err
		}
		return s.repo.CreateOrderWithTx(tx, order)
	}); err != nil {
		s.rdb.Incr(ctx, stockKey)
		return &dto.FlashBuyResp{Success: false, Message: "order creation failed"}, nil
	}

	s.rdb.Set(ctx, limitKey, 1, 24*time.Hour)

	return &dto.FlashBuyResp{
		Success: true,
		OrderID: order.ID,
		Message: "order created successfully",
	}, nil
}
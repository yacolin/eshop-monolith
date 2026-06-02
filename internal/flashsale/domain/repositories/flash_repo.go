package repositories

import (
	"context"
	"time"

	"eshop-monolith/internal/flashsale/domain/models"
	infraModels "eshop-monolith/internal/infra/repository/models"

	"gorm.io/gorm"
)

type FlashRepository struct {
	db *gorm.DB
}

func NewFlashRepository(db *gorm.DB) *FlashRepository {
	return &FlashRepository{db: db}
}

func (r *FlashRepository) CreateActivity(ctx context.Context, a *models.FlashActivity) error {
	po := infraModels.FlashActivityFromDomain(a)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	a.ID = po.ID
	return nil
}

func (r *FlashRepository) GetActivity(ctx context.Context, id int64) (*models.FlashActivity, error) {
	var po infraModels.FlashActivityPO
	if err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *FlashRepository) ListActivities(ctx context.Context) ([]models.FlashActivity, error) {
	var pos []infraModels.FlashActivityPO
	if err := r.db.WithContext(ctx).Order("start_time desc").Find(&pos).Error; err != nil {
		return nil, err
	}
	activities := make([]models.FlashActivity, len(pos))
	for i, po := range pos {
		activities[i] = *po.ToDomain()
	}
	return activities, nil
}

func (r *FlashRepository) UpdateSoldStock(ctx context.Context, activityID int64, delta int) error {
	return r.db.WithContext(ctx).Model(&infraModels.FlashActivityPO{}).
		Where("id = ?", activityID).
		UpdateColumn("sold_stock", gorm.Expr("sold_stock + ?", delta)).Error
}

func (r *FlashRepository) CreateOrder(ctx context.Context, o *models.FlashOrder) error {
	po := infraModels.FlashOrderFromDomain(o)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	o.ID = po.ID
	return nil
}

func (r *FlashRepository) GetOrder(ctx context.Context, orderID int64) (*models.FlashOrder, error) {
	var po infraModels.FlashOrderPO
	if err := r.db.WithContext(ctx).First(&po, "id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *FlashRepository) GetUserOrders(ctx context.Context, userID int64, activityID int64) ([]models.FlashOrder, error) {
	var pos []infraModels.FlashOrderPO
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if activityID > 0 {
		query = query.Where("activity_id = ?", activityID)
	}
	if err := query.Order("created_at desc").Find(&pos).Error; err != nil {
		return nil, err
	}
	orders := make([]models.FlashOrder, len(pos))
	for i, po := range pos {
		orders[i] = *po.ToDomain()
	}
	return orders, nil
}

func (r *FlashRepository) UpdateActivityStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&infraModels.FlashActivityPO{}).
		Where("id = ?", id).Update("status", status).Error
}

func (r *FlashRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&infraModels.FlashActivityPO{}, &infraModels.FlashOrderPO{})
}

func (r *FlashRepository) GetActiveActivities(ctx context.Context, now time.Time) ([]models.FlashActivity, error) {
	var pos []infraModels.FlashActivityPO
	if err := r.db.WithContext(ctx).
		Where("status = ? AND start_time <= ? AND end_time > ?", models.FlashStatusActive, now, now).
		Find(&pos).Error; err != nil {
		return nil, err
	}
	activities := make([]models.FlashActivity, len(pos))
	for i, po := range pos {
		activities[i] = *po.ToDomain()
	}
	return activities, nil
}

func (r *FlashRepository) IncrementSoldStockWithTx(tx *gorm.DB, activityID int64, delta int) error {
	return tx.Model(&infraModels.FlashActivityPO{}).
		Where("id = ? AND total_stock - sold_stock >= ?", activityID, delta).
		UpdateColumn("sold_stock", gorm.Expr("sold_stock + ?", delta)).Error
}

func (r *FlashRepository) CreateOrderWithTx(tx *gorm.DB, o *models.FlashOrder) error {
	po := infraModels.FlashOrderFromDomain(o)
	if err := tx.Create(po).Error; err != nil {
		return err
	}
	o.ID = po.ID
	return nil
}

func (r *FlashRepository) UpdateOrderStatusWithTx(tx *gorm.DB, orderID int64, status string) error {
	return tx.Model(&infraModels.FlashOrderPO{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}

func (r *FlashRepository) UpdateSoldStockWithTx(tx *gorm.DB, activityID int64, delta int) error {
	if delta >= 0 {
		return tx.Model(&infraModels.FlashActivityPO{}).
			Where("id = ? AND total_stock - sold_stock >= ?", activityID, delta).
			UpdateColumn("sold_stock", gorm.Expr("sold_stock + ?", delta)).Error
	}
	return tx.Model(&infraModels.FlashActivityPO{}).
		Where("id = ? AND sold_stock >= ?", activityID, -delta).
		UpdateColumn("sold_stock", gorm.Expr("sold_stock + ?", delta)).Error
}
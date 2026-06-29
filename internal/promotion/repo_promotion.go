package promotion

import (
	"context"

	"gorm.io/gorm"
)

type IpromotionRepository interface {
	Create(ctx context.Context, p *Promotion) error
	FindByID(ctx context.Context, id int64) (*Promotion, error)
	List(ctx context.Context, status *int8, promoType *int8, page, size int) ([]Promotion, int64, error)
	Update(ctx context.Context, p *Promotion) error
	Delete(ctx context.Context, id int64) error
	// Rule
	CreateRule(ctx context.Context, r *PromotionRule) error
	FindRuleByPromotionID(ctx context.Context, promotionID int64) (*PromotionRule, error)
	UpdateRule(ctx context.Context, r *PromotionRule) error
	DeleteRuleByPromotionID(ctx context.Context, promotionID int64) error
	// Product scope
	SetProducts(tx *gorm.DB, promotionID int64, products []PromotionProduct) error
	DeleteProductsByPromotion(ctx context.Context, promotionID int64) error
	ListProducts(ctx context.Context, promotionID int64) ([]PromotionProduct, error)
	// Coupon
	ClaimCoupon(ctx context.Context, up *UserPromotion) error
	FindUserPromotion(ctx context.Context, userID, promotionID int64) (*UserPromotion, error)
	ListUserPromotions(ctx context.Context, userID int64, status *int8, page, size int) ([]UserPromotion, int64, error)
	UseCoupon(ctx context.Context, id int64, orderID int64) error
}

type PromotionRepository struct {
	db *gorm.DB
}

func NewPromotionRepository(db *gorm.DB) IpromotionRepository {
	return &PromotionRepository{db: db}
}

func (r *PromotionRepository) Create(ctx context.Context, p *Promotion) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PromotionRepository) FindByID(ctx context.Context, id int64) (*Promotion, error) {
	var p Promotion
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	return &p, err
}

func (r *PromotionRepository) List(ctx context.Context, status, promoType *int8, page, size int) ([]Promotion, int64, error) {
	db := r.db.WithContext(ctx).Model(&Promotion{})
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	if promoType != nil {
		db = db.Where("promo_type = ?", *promoType)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []Promotion
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *PromotionRepository) Update(ctx context.Context, p *Promotion) error {
	return r.db.WithContext(ctx).Model(p).Updates(p).Error
}

func (r *PromotionRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Promotion{}).Error
}

func (r *PromotionRepository) CreateRule(ctx context.Context, rule *PromotionRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *PromotionRepository) FindRuleByPromotionID(ctx context.Context, promotionID int64) (*PromotionRule, error) {
	var rule PromotionRule
	err := r.db.WithContext(ctx).Where("promotion_id = ?", promotionID).First(&rule).Error
	return &rule, err
}

func (r *PromotionRepository) UpdateRule(ctx context.Context, rule *PromotionRule) error {
	return r.db.WithContext(ctx).Model(rule).Updates(rule).Error
}

func (r *PromotionRepository) DeleteRuleByPromotionID(ctx context.Context, promotionID int64) error {
	return r.db.WithContext(ctx).Where("promotion_id = ?", promotionID).Delete(&PromotionRule{}).Error
}

func (r *PromotionRepository) SetProducts(tx *gorm.DB, promotionID int64, products []PromotionProduct) error {
	if err := tx.Where("promotion_id = ?", promotionID).Delete(&PromotionProduct{}).Error; err != nil {
		return err
	}
	if len(products) > 0 {
		return tx.Create(&products).Error
	}
	return nil
}

func (r *PromotionRepository) DeleteProductsByPromotion(ctx context.Context, promotionID int64) error {
	return r.db.WithContext(ctx).Where("promotion_id = ?", promotionID).Delete(&PromotionProduct{}).Error
}

func (r *PromotionRepository) ListProducts(ctx context.Context, promotionID int64) ([]PromotionProduct, error) {
	var list []PromotionProduct
	err := r.db.WithContext(ctx).Where("promotion_id = ?", promotionID).Order("product_type ASC, id ASC").Find(&list).Error
	return list, err
}

// Coupon / UserPromotion

func (r *PromotionRepository) ClaimCoupon(ctx context.Context, up *UserPromotion) error {
	return r.db.WithContext(ctx).Create(up).Error
}

func (r *PromotionRepository) FindUserPromotion(ctx context.Context, userID, promotionID int64) (*UserPromotion, error) {
	var up UserPromotion
	err := r.db.WithContext(ctx).Where("user_id = ? AND promotion_id = ?", userID, promotionID).First(&up).Error
	return &up, err
}

func (r *PromotionRepository) ListUserPromotions(ctx context.Context, userID int64, status *int8, page, size int) ([]UserPromotion, int64, error) {
	db := r.db.WithContext(ctx).Model(&UserPromotion{}).Where("user_id = ?", userID)
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []UserPromotion
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *PromotionRepository) UseCoupon(ctx context.Context, id int64, orderID int64) error {
	return r.db.WithContext(ctx).Model(&UserPromotion{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": 2, "used_time": nil, "order_id": orderID}).Error
}

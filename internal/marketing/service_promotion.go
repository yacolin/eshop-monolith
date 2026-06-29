package marketing

import (
	"context"
	"errors"
	"time"

	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

type PromotionService struct {
	repo IpromotionRepository
	db   *gorm.DB
}

func NewPromotionService(repo IpromotionRepository, db *gorm.DB) *PromotionService {
	return &PromotionService{repo: repo, db: db}
}

type PromotionListResult struct {
	Total int64        `json:"total"`
	List  []*Promotion `json:"list"`
}

// Create 创建促销（含规则和商品范围）
func (s *PromotionService) Create(ctx context.Context, req *CreatePromotionReq) (*Promotion, error) {
	start, _ := time.Parse("2006-01-02 15:04:05", req.StartTime)
	end, _ := time.Parse("2006-01-02 15:04:05", req.EndTime)

	p := &Promotion{
		PromoName:     req.PromoName,
		PromoType:     req.PromoType,
		PromoCode:     req.PromoCode,
		StartTime:     start,
		EndTime:       end,
		TotalQuantity: req.TotalQuantity,
		PerUserLimit:  req.PerUserLimit,
		Status:        1,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		// Create rule
		rule := &PromotionRule{
			PromotionID:    p.ID,
			RuleName:       req.RuleName,
			ConditionType:  req.ConditionType,
			ConditionValue: req.ConditionValue,
			BenefitType:    req.BenefitType,
			BenefitValue:   req.BenefitValue,
			IsStackable:    req.IsStackable,
			StackPriority:  req.StackPriority,
		}
		if err := tx.Create(rule).Error; err != nil {
			return err
		}
		// Set products
		if len(req.ProductIDs) > 0 {
			products := make([]PromotionProduct, len(req.ProductIDs))
			for i, pid := range req.ProductIDs {
				products[i] = PromotionProduct{
					PromotionID: p.ID,
					ProductType: 3,
					ProductID:   &pid,
				}
			}
			return tx.Create(&products).Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PromotionService) GetByID(ctx context.Context, id int64) (*Promotion, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPromotionNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *PromotionService) List(ctx context.Context, req *PromotionListReq) (*PromotionListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	list, total, err := s.repo.List(ctx, req.Status, req.PromoType, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*Promotion, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return &PromotionListResult{Total: total, List: items}, nil
}

func (s *PromotionService) Update(ctx context.Context, id int64, req *UpdatePromotionReq) (*Promotion, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrPromotionNotFound
		}
		return nil, err
	}
	if req.PromoName != nil {
		p.PromoName = *req.PromoName
	}
	if req.StartTime != nil {
		t, _ := time.Parse("2006-01-02 15:04:05", *req.StartTime)
		p.StartTime = t
	}
	if req.EndTime != nil {
		t, _ := time.Parse("2006-01-02 15:04:05", *req.EndTime)
		p.EndTime = t
	}
	if req.TotalQuantity != nil {
		p.TotalQuantity = *req.TotalQuantity
	}
	if req.PerUserLimit != nil {
		p.PerUserLimit = *req.PerUserLimit
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	// Update rule
	if req.RuleName != nil || req.ConditionType != nil || req.BenefitType != nil {
		rule, err := s.repo.FindRuleByPromotionID(ctx, id)
		if err == nil {
			if req.RuleName != nil {
				rule.RuleName = *req.RuleName
			}
			if req.ConditionType != nil {
				rule.ConditionType = *req.ConditionType
			}
			if req.ConditionValue != nil {
				rule.ConditionValue = *req.ConditionValue
			}
			if req.BenefitType != nil {
				rule.BenefitType = *req.BenefitType
			}
			if req.BenefitValue != nil {
				rule.BenefitValue = *req.BenefitValue
			}
			if req.IsStackable != nil {
				rule.IsStackable = *req.IsStackable
			}
			if req.StackPriority != nil {
				rule.StackPriority = *req.StackPriority
			}
			s.repo.UpdateRule(ctx, rule)
		}
	}
	return p, nil
}

func (s *PromotionService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrPromotionNotFound
		}
		return err
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("promotion_id = ?", id).Delete(&PromotionRule{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&Promotion{}).Error
	})
}

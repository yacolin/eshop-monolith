package promotion

type CreatePromotionReq struct {
	PromoName     string  `json:"promo_name" binding:"required,max=100"`
	PromoType     int8    `json:"promo_type" binding:"required,oneof=1 2 3 4 5 6"`
	PromoCode     string  `json:"promo_code" binding:"max=50"`
	StartTime     string  `json:"start_time" binding:"required"`
	EndTime       string  `json:"end_time" binding:"required"`
	TotalQuantity int     `json:"total_quantity"`
	PerUserLimit  int     `json:"per_user_limit"`
	// Rule
	RuleName      string  `json:"rule_name" binding:"max=100"`
	ConditionType int8    `json:"condition_type" binding:"oneof=1 2 3 4"`
	ConditionValue float64 `json:"condition_value"`
	BenefitType   int8    `json:"benefit_type" binding:"oneof=1 2 3 4 5"`
	BenefitValue  float64 `json:"benefit_value"`
	IsStackable   int8    `json:"is_stackable" binding:"oneof=0 1"`
	StackPriority int     `json:"stack_priority"`
	// Products
	ProductIDs []int64 `json:"product_ids"`
}

type UpdatePromotionReq struct {
	PromoName     *string  `json:"promo_name" binding:"omitempty,max=100"`
	StartTime     *string  `json:"start_time"`
	EndTime       *string  `json:"end_time"`
	TotalQuantity *int     `json:"total_quantity"`
	PerUserLimit  *int     `json:"per_user_limit"`
	Status        *int8    `json:"status" binding:"omitempty,oneof=1 2 3 4"`
	RuleName      *string  `json:"rule_name" binding:"omitempty,max=100"`
	ConditionType *int8    `json:"condition_type" binding:"omitempty,oneof=1 2 3 4"`
	ConditionValue *float64 `json:"condition_value"`
	BenefitType   *int8    `json:"benefit_type" binding:"omitempty,oneof=1 2 3 4 5"`
	BenefitValue  *float64 `json:"benefit_value"`
	IsStackable   *int8    `json:"is_stackable" binding:"omitempty,oneof=0 1"`
	StackPriority *int     `json:"stack_priority"`
}

type PromotionListReq struct {
	Page      int   `form:"page,default=1" binding:"gte=1"`
	Size      int   `form:"size,default=10" binding:"gte=1,lte=100"`
	Status    *int8 `form:"status"`
	PromoType *int8 `form:"promo_type"`
}

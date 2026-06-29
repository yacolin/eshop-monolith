package promotion

type ClaimCouponReq struct {
	PromotionID int64 `json:"promotion_id" binding:"required"`
}

type UseCouponReq struct {
	UserPromotionID int64 `json:"user_promotion_id" binding:"required"`
	OrderID         int64 `json:"order_id" binding:"required"`
}

type UserPromotionListReq struct {
	Page   int   `form:"page,default=1" binding:"gte=1"`
	Size   int   `form:"size,default=10" binding:"gte=1,lte=100"`
	Status *int8 `form:"status"`
}

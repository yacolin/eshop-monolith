package promotion

type FlashBuyReq struct {
	PromotionID int64 `json:"promotion_id" binding:"required"`
	ProductID   int64 `json:"product_id" binding:"required"`
	SkuID       int64 `json:"sku_id" binding:"required"`
	Quantity    int   `json:"quantity" binding:"required,gt=0,lte=99"`
}

type FlashConfirmReq struct {
	Token     string `json:"token" binding:"required"`
	AddressID int64  `json:"address_id"`
}

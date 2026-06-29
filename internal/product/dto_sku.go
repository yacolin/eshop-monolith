package product

import "eshop-monolith/pkg/query"

type UpdateSKUReq struct {
	Price       *int64   `json:"price" binding:"omitempty,gt=0"`
	MarketPrice *int64   `json:"market_price" binding:"omitempty,gte=0"`
	CostPrice   *int64   `json:"cost_price" binding:"omitempty,gte=0"`
	Status      *int8    `json:"status" binding:"omitempty,oneof=0 1"`
	Image       *string  `json:"image" binding:"omitempty,max=512"`
	Barcode     *string  `json:"barcode" binding:"omitempty,max=50"`
	Weight      *float64 `json:"weight"`
	Volume      *float64 `json:"volume"`
	Length      *float64 `json:"length"`
	Width       *float64 `json:"width"`
	Height      *float64 `json:"height"`
}

type SKUListReq struct {
	query.Pagination
	ProductID int64 `form:"product_id"`
	SkuCode   string `form:"sku_code"`
}

type SKUListResult struct {
	Total int64  `json:"total"`
	List  []*SKU `json:"list"`
}

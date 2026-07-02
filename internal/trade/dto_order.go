package trade

import "eshop-monolith/pkg/query"

type CreateOrderItem struct {
	SkuID    int64 `json:"sku_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"required,gt=0,lte=99"`
}

type AddressInfo struct {
	Consignee  string `json:"consignee" binding:"required,max=64"`
	Phone      string `json:"phone" binding:"required,max=20"`
	Province   string `json:"province" binding:"max=32"`
	City       string `json:"city" binding:"max=32"`
	District   string `json:"district" binding:"max=32"`
	DetailAddr string `json:"detail_addr" binding:"max=256"`
	ZipCode    string `json:"zip_code" binding:"max=10"`
}

type CreateOrderReq struct {
	Address     AddressInfo       `json:"address" binding:"required"`
	Items       []CreateOrderItem `json:"items" binding:"required,min=1,dive"`
	CouponID    *int64            `json:"coupon_id"`
	BuyerRemark string            `json:"buyer_remark" binding:"max=500"`
	Source      string            `json:"source" binding:"max=20"`
}

type OrderListReq struct {
	query.Pagination
	UserID        int64  `form:"user_id"`
	Status        string `form:"status"`
	PaymentStatus string `form:"payment_status"`
	MerchantID    int64  `form:"merchant_id"`
	OrderNo       string `form:"order_no"`
}

type OrderListResult struct {
	Total int64    `json:"total"`
	List  []*Order `json:"list"`
}

type UpdateOrderStatusReq struct {
	Status string `json:"status" binding:"required,oneof=cancelled shipped delivered completed"`
	Note   string `json:"note" binding:"max=500"`
}

// OrderDetailResponse 订单详情（含订单项）
type OrderDetailResponse struct {
	*Order
	Items []OrderItem `json:"items"`
}

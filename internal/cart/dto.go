package cart

// AddItemReq 添加商品请求
type AddItemReq struct {
	SkuID    int64 `json:"sku_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"required,gt=0,lte=99"`
}

// UpdateItemReq 更新数量请求（quantity=0 时删除）
type UpdateItemReq struct {
	SkuID    int64 `json:"sku_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"gte=0,lte=99"`
}

// CartItemResponse 购物车商品响应
type CartItemResponse struct {
	ID          int64  `json:"id"`
	SkuID       int64  `json:"sku_id"`
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	SkuSpec     string `json:"sku_spec,omitempty"`
	Image       string `json:"image"`
	Price       int64  `json:"price"`
	Quantity    int    `json:"quantity"`
	Subtotal    int64  `json:"subtotal"`
}

// CartResponse 购物车响应
type CartResponse struct {
	ID          int64              `json:"id"`
	ItemCount   int                `json:"item_count"`
	TotalAmount int64              `json:"total_amount"`
	Items       []CartItemResponse `json:"items"`
}

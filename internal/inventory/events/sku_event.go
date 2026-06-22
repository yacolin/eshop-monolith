package events

type SkuCreatedEvent struct {
	SkuID     int64  `json:"sku_id"`
	ProductID int64  `json:"product_id"`
	Price     int64  `json:"price"`
}

type SkuUpdatedEvent struct {
	SkuID     int64  `json:"sku_id"`
	ProductID int64  `json:"product_id"`
	Price     int64  `json:"price"`
}

type SkuDeletedEvent struct {
	SkuID     int64 `json:"sku_id"`
	ProductID int64 `json:"product_id"`
}

package events

type SkuCreatedEvent struct {
	SkuID     int64 `json:"sku_id"`
	ProductID int64 `json:"product_id"`
	Price     int64 `json:"price"`
}

func (e SkuCreatedEvent) RoutingKey() string { return "sku.created" }

type SkuUpdatedEvent struct {
	SkuID     int64 `json:"sku_id"`
	ProductID int64 `json:"product_id"`
	Price     int64 `json:"price"`
}

func (e SkuUpdatedEvent) RoutingKey() string { return "sku.updated" }

type SkuDeletedEvent struct {
	SkuID     int64 `json:"sku_id"`
	ProductID int64 `json:"product_id"`
}

func (e SkuDeletedEvent) RoutingKey() string { return "sku.deleted" }

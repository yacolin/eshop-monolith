package events

// ProductCreatedEvent 产品创建事件
type ProductCreatedEvent struct {
	ProductID  int64  `json:"product_id"`
	Name       string `json:"name"`
	CategoryID int64  `json:"category_id"`
}

func (e ProductCreatedEvent) RoutingKey() string { return "product.created" }

// ProductUpdatedEvent 产品更新事件
type ProductUpdatedEvent struct {
	ProductID  int64  `json:"product_id"`
	Name       string `json:"name"`
	CategoryID int64  `json:"category_id"`
}

func (e ProductUpdatedEvent) RoutingKey() string { return "product.updated" }

// ProductDeletedEvent 产品删除事件
type ProductDeletedEvent struct {
	ProductID  int64  `json:"product_id"`
	Name       string `json:"name"`
	CategoryID int64  `json:"category_id"`
}

func (e ProductDeletedEvent) RoutingKey() string { return "product.deleted" }

package events

// ProductCreatedEvent 产品创建事件
type ProductCreatedEvent struct {
	ProductID  int64  `json:"product_id"`
	Name       string `json:"name"`
	CategoryID int64  `json:"category_id"`
}

// ProductUpdatedEvent 产品更新事件
type ProductUpdatedEvent struct {
	ProductID  int64  `json:"product_id"`
	Name       string `json:"name"`
	CategoryID int64  `json:"category_id"`
}

// ProductDeletedEvent 产品删除事件
type ProductDeletedEvent struct {
	ProductID  int64  `json:"product_id"`
	Name       string `json:"name"`
	CategoryID int64  `json:"category_id"`
}

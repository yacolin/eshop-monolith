package events

// CategoryCreatedEvent 分类创建事件
type CategoryCreatedEvent struct {
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	ParentID   *int64 `json:"parent_id"`
	Level      int    `json:"level"`
	Path       string `json:"path"`
}

func (e CategoryCreatedEvent) RoutingKey() string { return "category.created" }

// CategoryUpdatedEvent 分类更新事件
type CategoryUpdatedEvent struct {
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	ParentID   *int64 `json:"parent_id"`
	Level      int    `json:"level"`
	Path       string `json:"path"`
}

func (e CategoryUpdatedEvent) RoutingKey() string { return "category.updated" }

// CategoryDeletedEvent 分类删除事件
type CategoryDeletedEvent struct {
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	ParentID   *int64 `json:"parent_id"`
	Level      int    `json:"level"`
	Path       string `json:"path"`
}

func (e CategoryDeletedEvent) RoutingKey() string { return "category.deleted" }

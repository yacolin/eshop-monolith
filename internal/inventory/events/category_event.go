package events

// CategoryCreatedEvent 分类创建事件
type CategoryCreatedEvent struct {
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	ParentID   *int64 `json:"parent_id"`
	Level      int    `json:"level"`
	Path       string `json:"path"`
}

// CategoryUpdatedEvent 分类更新事件
type CategoryUpdatedEvent struct {
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	ParentID   *int64 `json:"parent_id"`
	Level      int    `json:"level"`
	Path       string `json:"path"`
}

// CategoryDeletedEvent 分类删除事件
type CategoryDeletedEvent struct {
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	ParentID   *int64 `json:"parent_id"`
	Level      int    `json:"level"`
	Path       string `json:"path"`
}

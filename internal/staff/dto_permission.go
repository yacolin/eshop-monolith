package staff

type PermissionListReq struct {
	Category string `form:"category"`
}

type PermissionListItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	ParentID    int64  `json:"parent_id"`
	Category    string `json:"category"`
	SortOrder   int    `json:"sort_order"`
	Status      int8   `json:"status"`
}

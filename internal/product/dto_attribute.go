package product

type CreateAttributeReq struct {
	Name       string `json:"name" binding:"required,max=100"`
	CategoryID int64  `json:"category_id" binding:"required"`
	InputType  int8   `json:"input_type" binding:"oneof=1 2 3 4"`
	Values     string `json:"values"` // JSON
	Unit       string `json:"unit" binding:"max=20"`
	Required   int8   `json:"required" binding:"oneof=0 1"`
	Searchable int8   `json:"searchable" binding:"oneof=0 1"`
	IsSkuSpec  int8   `json:"is_sku_spec" binding:"oneof=0 1"`
	SortOrder  int    `json:"sort_order"`
}

type UpdateAttributeReq struct {
	Name       *string `json:"name" binding:"omitempty,max=100"`
	InputType  *int8   `json:"input_type" binding:"omitempty,oneof=1 2 3 4"`
	Values     *string `json:"values"`
	Unit       *string `json:"unit" binding:"omitempty,max=20"`
	Required   *int8   `json:"required" binding:"omitempty,oneof=0 1"`
	Searchable *int8   `json:"searchable" binding:"omitempty,oneof=0 1"`
	IsSkuSpec  *int8   `json:"is_sku_spec" binding:"omitempty,oneof=0 1"`
	SortOrder  *int    `json:"sort_order"`
	Status     *int8   `json:"status" binding:"omitempty,oneof=0 1"`
}

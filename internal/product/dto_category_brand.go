package product

type SetCategoryBrandsReq struct {
	BrandIDs  []int64 `json:"brand_ids" binding:"required"`
	SortOrder int     `json:"sort_order"`
}

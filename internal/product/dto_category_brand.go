package product

type SetCategoryBrandsReq struct {
	BrandIDs  []int64 `json:"brand_ids" binding:"required"`
	SortOrder int     `json:"sort_order"`
}

// CategoryBrandDetail 类目关联品牌的详情（含品牌信息）
type CategoryBrandDetail struct {
	BrandID     int64  `json:"brand_id"`
	BrandName   string `json:"brand_name"`
	EnglishName string `json:"english_name"`
	LogoURL     string `json:"logo_url"`
	FirstLetter string `json:"first_letter"`
	SortOrder   int    `json:"sort_order"`
}

package product

import "eshop-monolith/pkg/query"

type CreateBrandReq struct {
	Name        string `json:"name" binding:"required,max=100"`
	EnglishName string `json:"english_name" binding:"max=100"`
	LogoURL     string `json:"logo_url" binding:"max=512"`
	FirstLetter string `json:"first_letter" binding:"omitempty,len=1,alpha"`
	SortOrder   *int   `json:"sort_order"`
	Status      *int8  `json:"status" binding:"omitempty,oneof=0 1"`
	Description string `json:"description"`
}

type UpdateBrandReq struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	EnglishName *string `json:"english_name" binding:"omitempty,max=100"`
	LogoURL     *string `json:"logo_url" binding:"omitempty,max=512"`
	FirstLetter *string `json:"first_letter" binding:"omitempty,len=1,alpha"`
	SortOrder   *int    `json:"sort_order"`
	Status      *int8   `json:"status" binding:"omitempty,oneof=0 1"`
	Description *string `json:"description"`
}

type BrandListReq struct {
	query.Pagination
	Name        string `form:"name"`
	FirstLetter string `form:"first_letter"`
	Status      *int8  `form:"status"`
}

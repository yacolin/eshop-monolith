package product

import "eshop-monolith/pkg/query"

// ── Create ───────────────────────────────────────

type CreateSKUItem struct {
	SkuCode     string            `json:"sku_code" binding:"required,max=100"`
	Barcode     string            `json:"barcode" binding:"max=50"`
	Spec        map[string]string `json:"spec"`
	Price       int64             `json:"price" binding:"required"`
	MarketPrice int64             `json:"market_price"`
	CostPrice   int64             `json:"cost_price"`
	Weight      float64           `json:"weight"`
	Volume      float64           `json:"volume"`
	Length      float64           `json:"length"`
	Width       float64           `json:"width"`
	Height      float64           `json:"height"`
	MinPurchase int               `json:"min_purchase_qty"`
	MaxPurchase int               `json:"max_purchase_qty"`
	Image       string            `json:"image" binding:"max=512"`
}

type CreateProductAttrItem struct {
	AttributeID int64  `json:"attribute_id" binding:"required"`
	Value       string `json:"value" binding:"required,max=500"`
}

type CreateSPUReq struct {
	Name        string                `json:"name" binding:"required,max=200"`
	Subtitle    string                `json:"subtitle" binding:"max=500"`
	CategoryID  int64                 `json:"category_id" binding:"required"`
	BrandID     int64                 `json:"brand_id"`
	Unit        string                `json:"unit" binding:"max=10"`
	MainImage   string                `json:"main_image" binding:"required,max=512"`
	Images      []string              `json:"images"`
	VideoURL    string                `json:"video_url" binding:"max=512"`
	Description string                `json:"description"`
	MobileDesc  string                `json:"mobile_description"`
	SortOrder   int                   `json:"sort_order"`
	CreatedBy   string                `json:"created_by" binding:"max=50"`
	SKUs        []CreateSKUItem       `json:"skus" binding:"required,min=1,dive"`
	Attributes  []CreateProductAttrItem `json:"attributes" binding:"dive"`
}

// ── Update ───────────────────────────────────────

type UpdateSPUReq struct {
	Name      *string   `json:"name" binding:"omitempty,max=200"`
	Subtitle  *string   `json:"subtitle" binding:"omitempty,max=500"`
	Unit      *string   `json:"unit" binding:"omitempty,max=10"`
	MainImage *string   `json:"main_image" binding:"omitempty,max=512"`
	Images    *[]string `json:"images"`
	VideoURL  *string   `json:"video_url" binding:"omitempty,max=512"`
	SortOrder *int      `json:"sort_order"`
	Status    *int8     `json:"status" binding:"omitempty,oneof=0 1 2 3 4"`
	UpdatedBy *string   `json:"updated_by" binding:"omitempty,max=50"`
}

// ── List ─────────────────────────────────────────

type SPUListReq struct {
	query.Pagination
	Name       string `form:"name"`
	CategoryID *int64 `form:"category_id"`
	BrandID    *int64 `form:"brand_id"`
	Status     *int8  `form:"status"`
	PriceMin   int64  `form:"price_min"`
	PriceMax   int64  `form:"price_max"`
}

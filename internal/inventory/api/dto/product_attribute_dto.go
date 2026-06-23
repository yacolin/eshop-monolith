package dto

// AttributeValueItem 属性的可选值
type AttributeValueItem struct {
	ValueID int64  `json:"value_id"`
	Value   string `json:"value"`
}

// ProductAttributeItem 产品的一个属性维度及其可选值
type ProductAttributeItem struct {
	AttributeID   int64                `json:"attribute_id"`
	AttributeName string               `json:"attribute_name"`
	Values        []AttributeValueItem `json:"values"`
}

// BatchCreateSkuItem 批量创建 SKU 中的单个 SKU
type BatchCreateSkuItem struct {
	AttrValueIDs []int64 `json:"attr_value_ids" binding:"required,min=1"`
	Name         string  `json:"name" binding:"required,max=255"`
	Price        int64   `json:"price" binding:"required,gt=0"`
	SKUCode      string  `json:"sku_code" binding:"required,max=100"`
	Image        string  `json:"image" binding:"max=500"`
}

// BatchCreateSkuDTO 批量创建 SKU 请求体
type BatchCreateSkuDTO struct {
	Skus []BatchCreateSkuItem `json:"skus" binding:"required,min=1,dive"`
}

// BatchCreateSkuResult 批量创建结果
type BatchCreateSkuResult struct {
	Total   int            `json:"total"`
	Success int            `json:"success"`
	Failed  int            `json:"failed"`
	Skus    []SkuResponse  `json:"skus"`
}

// ProductAttributeUpdateItem 单个属性维度的值选择
type ProductAttributeUpdateItem struct {
	AttributeID int64   `json:"attribute_id" binding:"required"`
	ValueIDs    []int64 `json:"value_ids" binding:"required,min=1,dive,gt=0"`
}

// UpdateProductAttributesDTO 更新产品属性关联请求体
type UpdateProductAttributesDTO struct {
	Attributes []ProductAttributeUpdateItem `json:"attributes" binding:"required,min=1,dive"`
}

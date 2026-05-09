package shared

// ProductCategory 产品与分类的多对多关联表
type ProductCategory struct {
	ProductID  int64
	CategoryID int64
}

// TableName 指定表名
func (ProductCategory) TableName() string {
	return "product_categories"
}

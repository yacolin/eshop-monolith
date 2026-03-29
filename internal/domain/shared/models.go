package shared

// ProductCategory 产品与分类的多对多关联表
type ProductCategory struct {
	ProductID  int64 `gorm:"primaryKey"`
	CategoryID int64 `gorm:"primaryKey"`
}

// TableName 指定表名
func (ProductCategory) TableName() string {
	return "product_categories"
}

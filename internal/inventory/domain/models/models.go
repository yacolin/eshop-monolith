package models

// Product 旧产品模型桩
type Product struct {
	ID   int64
	Name string
}

// Sku 旧SKU模型桩
type Sku struct {
	ID        int64
	ProductID int64
	Price     int64
	SKUCode   string
	Status    int
}

// Category 旧分类模型桩
type Category struct {
	ID   int64
	Name string
}

// Inventory 旧库存模型桩
type Inventory struct {
	ID        int64
	SkuID     int64
	Quantity  int
	Reserved  int
	Threshold int
	Status    string
}

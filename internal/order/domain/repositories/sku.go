package repositories

import (
	"context"

	"gorm.io/gorm"
)

// SkuForOrder 订单模块查询 SKU 信息的防火墙接口
type SkuForOrder interface {
	// GetSkuInfo 根据 skuID 反查 productID 和 price
	GetSkuInfo(ctx context.Context, skuID int64) (productID int64, price int64, err error)
}

// SkuForOrderAdapter 轻量适配器，直查 skus 表
type SkuForOrderAdapter struct {
	db *gorm.DB
}

func NewSkuForOrderAdapter(db *gorm.DB) SkuForOrder {
	return &SkuForOrderAdapter{db: db}
}

func (a *SkuForOrderAdapter) GetSkuInfo(ctx context.Context, skuID int64) (int64, int64, error) {
	var row struct {
		ProductID int64
		Price     int64
	}
	err := a.db.WithContext(ctx).Table("skus").
		Select("product_id, price").
		Where("id = ?", skuID).Take(&row).Error
	if err != nil {
		return 0, 0, err
	}
	return row.ProductID, row.Price, nil
}

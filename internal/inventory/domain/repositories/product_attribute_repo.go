package repositories

import (
	"context"

	"eshop-monolith/internal/infra/repository/models"
	domain "eshop-monolith/internal/inventory/domain/models"

	"gorm.io/gorm"
)

type IproductAttributeRepository interface {
	// FindByProductID 查询产品的属性定义及值列表（按 sort_order 排序）
	FindByProductID(ctx context.Context, productID int64) ([]domain.Attribute, []domain.AttributeValue, error)
	// BatchCreate 批量创建产品-属性-值关联
	BatchCreate(ctx context.Context, productID int64, items []domain.ProductAttributeValue) error
	// DeleteByProductID 删除产品的所有属性关联
	DeleteByProductID(ctx context.Context, productID int64) error
}

type ProductAttributeRepository struct {
	db *gorm.DB
}

func NewProductAttributeRepository(db *gorm.DB) IproductAttributeRepository {
	return &ProductAttributeRepository{db: db}
}

// attributeRow 用于扫描多表 JOIN 结果
type attributeRow struct {
	AttrID   int64
	AttrName string
	ValID    int64
	ValValue string
}

func (r *ProductAttributeRepository) FindByProductID(ctx context.Context, productID int64) ([]domain.Attribute, []domain.AttributeValue, error) {
	var rows []attributeRow
	if err := r.db.WithContext(ctx).Table("product_attribute_values").
		Select("a.id as attr_id, a.name as attr_name, av.id as val_id, av.value as val_value").
		Joins("JOIN attribute_attributes a ON a.id = product_attribute_values.attribute_id").
		Joins("JOIN attribute_values av ON av.id = product_attribute_values.attribute_value_id").
		Where("product_attribute_values.product_id = ?", productID).
		Order("a.sort_order ASC, av.sort_order ASC").
		Scan(&rows).Error; err != nil {
		return nil, nil, err
	}

	attrs, vals := r.buildAttributeResult(rows)
	return attrs, vals, nil
}

// buildAttributeResult 将 rows 去重构造为 Attr → Values 结构
func (r *ProductAttributeRepository) buildAttributeResult(rows []attributeRow) ([]domain.Attribute, []domain.AttributeValue) {
	attrMap := make(map[int64]*domain.Attribute)
	vals := make([]domain.AttributeValue, 0, len(rows))
	for _, row := range rows {
		if _, ok := attrMap[row.AttrID]; !ok {
			attrMap[row.AttrID] = &domain.Attribute{
				ID:   row.AttrID,
				Name: row.AttrName,
			}
		}
		vals = append(vals, domain.AttributeValue{
			ID:          row.ValID,
			AttributeID: row.AttrID,
			Value:       row.ValValue,
		})
	}

	attrs := make([]domain.Attribute, 0, len(attrMap))
	for _, a := range attrMap {
		attrs = append(attrs, *a)
	}
	return attrs, vals
}

func (r *ProductAttributeRepository) BatchCreate(ctx context.Context, productID int64, items []domain.ProductAttributeValue) error {
	if len(items) == 0 {
		return nil
	}
	pos := make([]models.ProductAttributeValuePO, len(items))
	for i, item := range items {
		pos[i] = *models.ProductAttributeValueFromDomain(&item)
	}
	return r.db.WithContext(ctx).Create(&pos).Error
}

func (r *ProductAttributeRepository) DeleteByProductID(ctx context.Context, productID int64) error {
	return r.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&models.ProductAttributeValuePO{}).Error
}

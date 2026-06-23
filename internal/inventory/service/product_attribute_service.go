package service

import (
	"context"
	"fmt"
	"strings"

	"eshop-monolith/internal/infra/repository/models"
	"eshop-monolith/internal/inventory/api/dto"
	invModels "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/inventory/domain/repositories"

	"gorm.io/gorm"
)

type ProductAttributeService struct {
	attrRepo     repositories.IproductAttributeRepository
	skuRepo      repositories.IskuRepository
	productRepo  repositories.IproductRepository
	db           *gorm.DB
}

func NewProductAttributeService(
	attrRepo repositories.IproductAttributeRepository,
	skuRepo repositories.IskuRepository,
	productRepo repositories.IproductRepository,
	db *gorm.DB,
) *ProductAttributeService {
	return &ProductAttributeService{
		attrRepo:    attrRepo,
		skuRepo:     skuRepo,
		productRepo: productRepo,
		db:          db,
	}
}

// GetProductAttributes 获取产品的属性定义及可选值列表
func (s *ProductAttributeService) GetProductAttributes(ctx context.Context, productID int64) ([]dto.ProductAttributeItem, error) {
	attrs, vals, err := s.attrRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	valMap := make(map[int64][]dto.AttributeValueItem)
	for _, v := range vals {
		valMap[v.AttributeID] = append(valMap[v.AttributeID], dto.AttributeValueItem{
			ValueID: v.ID,
			Value:   v.Value,
		})
	}

	result := make([]dto.ProductAttributeItem, 0, len(attrs))
	for _, a := range attrs {
		result = append(result, dto.ProductAttributeItem{
			AttributeID:   a.ID,
			AttributeName: a.Name,
			Values:        valMap[a.ID],
		})
	}
	return result, nil
}

// UpdateProductAttributes 更新产品关联的属性值（全量替换）
func (s *ProductAttributeService) UpdateProductAttributes(ctx context.Context, productID int64, req *dto.UpdateProductAttributesDTO) error {
	if _, err := s.productRepo.FindByID(ctx, productID); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先清除原有关联
		if err := tx.Where("product_id = ?", productID).Delete(&invModels.ProductAttributeValue{}).Error; err != nil {
			return err
		}

		// 批量插入新关联
		for _, attr := range req.Attributes {
			for _, valID := range attr.ValueIDs {
				if err := tx.Create(&invModels.ProductAttributeValue{
					ProductID:       productID,
					AttributeID:     attr.AttributeID,
					AttributeValueID: valID,
				}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// BatchCreateSkus 批量创建 SKU（含 sku_attributes + product_attribute_values 关联）
func (s *ProductAttributeService) BatchCreateSkus(ctx context.Context, productID int64, req *dto.BatchCreateSkuDTO) (*dto.BatchCreateSkuResult, error) {
	if _, err := s.productRepo.FindByID(ctx, productID); err != nil {
		return nil, err
	}

	created := make([]dto.SkuResponse, 0, len(req.Skus))
	attrValSet := make(map[int64]struct{}) // 收集所有出现的属性值 ID
	var failCount int

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Skus {
			spec, err := s.buildSpec(ctx, tx, item.AttrValueIDs)
			if err != nil {
				return fmt.Errorf("构建 spec 失败: %w", err)
			}

			sku := &invModels.Sku{
				ProductID: productID,
				Name:      item.Name,
				Price:     item.Price,
				SKUCode:   item.SKUCode,
				Image:     item.Image,
				Spec:      spec,
			}
			skuPO := models.SkuFromDomain(sku)
			if err := tx.Create(skuPO).Error; err != nil {
				if strings.Contains(err.Error(), "Duplicate") {
					failCount++
					continue
				}
				return err
			}
			sku.ID = skuPO.ID

			for _, valID := range item.AttrValueIDs {
				var attrID int64
				if err := tx.Model(&invModels.AttributeValue{}).
					Select("attribute_id").Where("id = ?", valID).
					Scan(&attrID).Error; err != nil {
					return fmt.Errorf("查询属性值 %d 失败: %w", valID, err)
				}
				if err := tx.Create(&invModels.SkuAttribute{
					SkuID: sku.ID, AttributeID: attrID, AttributeValueID: valID,
				}).Error; err != nil {
					return err
				}
				attrValSet[valID] = struct{}{}
			}

			created = append(created, dto.SkuToResponse(sku))
		}

		// 批量写入 product_attribute_values（产品-属性值关联）
		if len(attrValSet) > 0 {
			type attrVal struct {
				AttrID int64
				ValID  int64
			}
			valIDs := make([]int64, 0, len(attrValSet))
			for id := range attrValSet {
				valIDs = append(valIDs, id)
			}
			var mappings []attrVal
			if err := tx.Model(&invModels.AttributeValue{}).
				Select("attribute_id as attr_id, id as val_id").
				Where("id IN ?", valIDs).
				Scan(&mappings).Error; err != nil {
				return err
			}
			for _, m := range mappings {
				if err := tx.Where("product_id = ? AND attribute_value_id = ?", productID, m.ValID).
					Delete(&invModels.ProductAttributeValue{}).Error; err != nil {
					return err
				}
				if err := tx.Create(&invModels.ProductAttributeValue{
					ProductID: productID, AttributeID: m.AttrID, AttributeValueID: m.ValID,
				}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.syncProductMinPrice(context.Background(), productID)

	return &dto.BatchCreateSkuResult{
		Total:   len(req.Skus),
		Success: len(created),
		Failed:  failCount,
		Skus:    created,
	}, nil
}

// buildSpec 根据属性值 ID 列表查询属性名→值名 映射
func (s *ProductAttributeService) buildSpec(_ context.Context, tx *gorm.DB, valIDs []int64) (map[string]string, error) {
	type row struct {
		AttrName string
		ValValue string
	}
	var rows []row
	if err := tx.Table("attribute_values av").
		Select("a.name as attr_name, av.value as val_value").
		Joins("JOIN attribute_attributes a ON a.id = av.attribute_id").
		Where("av.id IN ?", valIDs).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	spec := make(map[string]string, len(rows))
	for _, r := range rows {
		spec[r.AttrName] = r.ValValue
	}
	return spec, nil
}

// syncProductMinPrice 重新计算产品最低价
func (s *ProductAttributeService) syncProductMinPrice(ctx context.Context, productID int64) {
	type result struct{ Price int64 }
	var r result
	s.db.WithContext(ctx).Table("skus").
		Select("MIN(price) as price").
		Where("product_id = ?", productID).Scan(&r)
	s.db.WithContext(ctx).Table("products").
		Where("id = ?", productID).Update("min_price", r.Price)
}

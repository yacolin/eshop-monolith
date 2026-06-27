package service

import (
	"context"
	"encoding/json"
	"fmt"

	"eshop-monolith/internal/infra/repository/models"
	"eshop-monolith/internal/inventory/api/dto"
	invModels "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/inventory/domain/repositories"

	"gorm.io/gorm"
)

type ProductAttributeService struct {
	attrRepo    repositories.IproductAttributeRepository
	skuRepo     repositories.IskuRepository
	productRepo repositories.IproductRepository
	db          *gorm.DB
}

func NewProductAttributeService(
	attrRepo repositories.IproductAttributeRepository,
	skuRepo repositories.IskuRepository,
	productRepo repositories.IproductRepository,
	_ repositories.IcategoryAttributeRepository,
	db *gorm.DB,
) *ProductAttributeService {
	return &ProductAttributeService{
		attrRepo:    attrRepo,
		skuRepo:     skuRepo,
		productRepo: productRepo,
		db:          db,
	}
}

// GetProductAttributes 获取产品的属性定义及可选值列表（按品类过滤）
func (s *ProductAttributeService) GetProductAttributes(ctx context.Context, productID int64) ([]dto.ProductAttributeItem, error) {
	attrs, vals, err := s.attrRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// 获取产品所属品类的允许属性集合
	allowed, err := s.getAllowedAttributes(ctx, productID)
	if err != nil {
		return nil, err
	}

	valMap := make(map[int64][]dto.AttributeValueItem)
	for _, v := range vals {
		if allowed != nil {
			if _, ok := allowed[v.AttributeID]; !ok {
				continue
			}
		}
		valMap[v.AttributeID] = append(valMap[v.AttributeID], dto.AttributeValueItem{
			ValueID: v.ID,
			Value:   v.Value,
		})
	}

	result := make([]dto.ProductAttributeItem, 0, len(attrs))
	for _, a := range attrs {
		if _, ok := valMap[a.ID]; !ok {
			continue
		}
		result = append(result, dto.ProductAttributeItem{
			AttributeID:   a.ID,
			AttributeName: a.Name,
			Values:        valMap[a.ID],
		})
	}
	return result, nil
}

// getAllowedAttributes 查询产品品类的允许属性集合（单次 SQL 替代多次往返）
func (s *ProductAttributeService) getAllowedAttributes(ctx context.Context, productID int64) (map[int64]struct{}, error) {
	var attrIDs []int64
	if err := s.db.WithContext(ctx).Raw(`
		SELECT DISTINCT ca.attribute_id
		FROM product_categories pc
		JOIN categories c ON c.id = pc.category_id
		JOIN category_attributes ca ON (ca.category_id = c.id OR ca.category_id = c.parent_id)
		WHERE pc.product_id = ?
	`, productID).Scan(&attrIDs).Error; err != nil {
		return nil, err
	}
	if len(attrIDs) == 0 {
		return nil, nil
	}
	allowed := make(map[int64]struct{}, len(attrIDs))
	for _, id := range attrIDs {
		allowed[id] = struct{}{}
	}
	return allowed, nil
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
	attrValSet := make(map[int64]struct{})
	var failCount int

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 批量查重
		skuCodes := make([]string, len(req.Skus))
		for i, item := range req.Skus {
			skuCodes[i] = item.SKUCode
		}
		var existingCodes []string
		tx.Model(&models.SkuPO{}).Where("sku_code IN ?", skuCodes).Pluck("sku_code", &existingCodes)
		dupSet := make(map[string]struct{}, len(existingCodes))
		for _, c := range existingCodes {
			dupSet[c] = struct{}{}
		}

		// 2. 批量查询属性映射（一次 JOIN 查所有）
		allValIDs := make([]int64, 0, len(req.Skus)*3)
		for _, item := range req.Skus {
			allValIDs = append(allValIDs, item.AttrValueIDs...)
		}
		type keyRow struct {
			ValID    int64  `gorm:"column:val_id"`
			AttrID   int64  `gorm:"column:attr_id"`
			AttrName string `gorm:"column:attr_name"`
			ValValue string `gorm:"column:val_value"`
		}
		var krows []keyRow
		if len(allValIDs) > 0 {
			if err := tx.Table("attribute_values av").
				Select("av.id as val_id, av.attribute_id as attr_id, a.name as attr_name, av.value as val_value").
				Joins("JOIN attribute_attributes a ON a.id = av.attribute_id").
				Where("av.id IN ?", allValIDs).
				Scan(&krows).Error; err != nil {
				return err
			}
		}
		attrIDByVal := make(map[int64]int64, len(krows))
		specByVal := make(map[int64]string, len(krows))
		nameByVal := make(map[int64]string, len(krows))
		for _, r := range krows {
			attrIDByVal[r.ValID] = r.AttrID
			nameByVal[r.ValID] = r.AttrName
			specByVal[r.ValID] = r.ValValue
		}

		// 3. 构建 SKU PO
		type skuBuild struct {
			po   *models.SkuPO
			item *dto.BatchCreateSkuItem
		}
		builds := make([]skuBuild, 0, len(req.Skus))
		skuPOs := make([]*models.SkuPO, 0, len(req.Skus))
		for _, item := range req.Skus {
			if _, dup := dupSet[item.SKUCode]; dup {
				failCount++
				continue
			}
			specMap := make(map[string]string, len(item.AttrValueIDs))
			for _, valID := range item.AttrValueIDs {
				if nm, ok := nameByVal[valID]; ok {
					specMap[nm] = specByVal[valID]
				}
			}
			sku := &invModels.Sku{
				ProductID: productID,
				Name:      item.Name,
				Price:     item.Price,
				SKUCode:   item.SKUCode,
				Image:     item.Image,
				Spec:      specMap,
			}
			po := models.SkuFromDomain(sku)
			skuPOs = append(skuPOs, po)
			builds = append(builds, skuBuild{po: po, item: &item})
		}
		if len(skuPOs) == 0 {
			return nil
		}

		// 4. 批量 Insert SKU PO（50 条一批）
		if err := tx.CreateInBatches(skuPOs, 50).Error; err != nil {
			return fmt.Errorf("批量创建 SKU 失败: %w", err)
		}

		// 5. 批量构建 sku_attributes 并 Insert
		skuAttrs := make([]*invModels.SkuAttribute, 0, len(builds)*3)
		for i, b := range builds {
			po := skuPOs[i]
			sku := &invModels.Sku{
				ID:        po.ID,
				ProductID: po.ProductID,
				Name:      po.Name,
				Price:     po.Price,
				SKUCode:   po.SKUCode,
				Image:     po.Image,
				Spec:      parseSpec(po.Spec),
			}
			for _, valID := range b.item.AttrValueIDs {
				attrID, ok := attrIDByVal[valID]
				if !ok {
					continue
				}
				skuAttrs = append(skuAttrs, &invModels.SkuAttribute{
					SkuID: po.ID, AttributeID: attrID, AttributeValueID: valID,
				})
				attrValSet[valID] = struct{}{}
			}
			created = append(created, dto.SkuToResponse(sku))
		}
		if len(skuAttrs) > 0 {
			if err := tx.Create(skuAttrs).Error; err != nil {
				return fmt.Errorf("批量创建 sku_attributes 失败: %w", err)
			}
		}

		// 6. 批量写入 product_attribute_values
		if len(attrValSet) > 0 {
			valIDs := make([]int64, 0, len(attrValSet))
			for id := range attrValSet {
				valIDs = append(valIDs, id)
			}
			type attrVal struct{ AttrID int64; ValID int64 }
			var avMappings []attrVal
			if err := tx.Model(&invModels.AttributeValue{}).
				Select("attribute_id as attr_id, id as val_id").
				Where("id IN ?", valIDs).
				Scan(&avMappings).Error; err != nil {
				return err
			}
			// 批量删除旧的关联
			tx.Where("product_id = ? AND attribute_value_id IN ?", productID, valIDs).
				Delete(&invModels.ProductAttributeValue{})
			// 批量插入新的关联
			pavs := make([]*invModels.ProductAttributeValue, 0, len(avMappings))
			for _, m := range avMappings {
				pavs = append(pavs, &invModels.ProductAttributeValue{
					ProductID: productID, AttributeID: m.AttrID, AttributeValueID: m.ValID,
				})
			}
			if err := tx.Create(pavs).Error; err != nil {
				return err
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


// parseSpec 将 PO 中的 JSON 字符串 Spec 解析为 map
func parseSpec(specStr string) map[string]string {
	if specStr == "" {
		return nil
	}
	var spec map[string]string
	if err := json.Unmarshal([]byte(specStr), &spec); err != nil {
		return nil
	}
	return spec
}

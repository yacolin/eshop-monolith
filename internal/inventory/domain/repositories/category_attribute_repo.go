package repositories

import (
	"context"

	"eshop-monolith/internal/infra/repository/models"
	domain "eshop-monolith/internal/inventory/domain/models"

	"gorm.io/gorm"
)

type IcategoryAttributeRepository interface {
	// FindByCategoryID 查询品类关联的属性列表
	FindByCategoryID(ctx context.Context, categoryID int64) ([]domain.Attribute, error)
	// FindByCategoryIDs 批量查询多个品类关联的属性列表（含父品类递归）
	FindByCategoryIDs(ctx context.Context, categoryIDs []int64) ([]domain.Attribute, error)
	// SetByCategoryID 全量替换品类关联的属性
	SetByCategoryID(ctx context.Context, categoryID int64, attributeIDs []int64) error
}

type CategoryAttributeRepository struct {
	db *gorm.DB
}

func NewCategoryAttributeRepository(db *gorm.DB) IcategoryAttributeRepository {
	return &CategoryAttributeRepository{db: db}
}

func (r *CategoryAttributeRepository) FindByCategoryID(ctx context.Context, categoryID int64) ([]domain.Attribute, error) {
	type row struct {
		ID        int64
		Name      string
		SortOrder int
	}

	// 向上递归查找：先查当前品类，无绑定则查父品类，直到根
	for curID := categoryID; curID != 0; {
		var rows []row
		if err := r.db.WithContext(ctx).Table("category_attributes").
			Select("a.id, a.name, a.sort_order").
			Joins("JOIN attribute_attributes a ON a.id = category_attributes.attribute_id").
			Where("category_attributes.category_id = ?", curID).
			Order("a.sort_order asc, a.id asc").
			Scan(&rows).Error; err != nil {
			return nil, err
		}
		if len(rows) > 0 {
			attrs := make([]domain.Attribute, len(rows))
			for i, r := range rows {
				attrs[i] = domain.Attribute{ID: r.ID, Name: r.Name, SortOrder: r.SortOrder}
			}
			return attrs, nil
		}
		// 当前品类无绑定，查父品类
		var parentID *int64
		if err := r.db.WithContext(ctx).Table("categories").
			Select("parent_id").Where("id = ?", curID).Scan(&parentID).Error; err != nil {
			return nil, err
		}
		if parentID == nil || *parentID == 0 {
			break
		}
		curID = *parentID
	}
	return nil, nil
}

// FindByCategoryIDs 批量查询多个品类关联的属性（含父品类递归，逐层 IN 替代全表扫描）
func (r *CategoryAttributeRepository) FindByCategoryIDs(ctx context.Context, categoryIDs []int64) ([]domain.Attribute, error) {
	if len(categoryIDs) == 0 {
		return nil, nil
	}

	// 逐层查父品类，子→父→根，直到全部耗尽
	type catRow struct{ ID, ParentID int64 }
	allIDs := make(map[int64]struct{})
	queue := make([]int64, len(categoryIDs))
	copy(queue, categoryIDs)

	for len(queue) > 0 {
		var rows []catRow
		r.db.WithContext(ctx).Table("categories").
			Select("id, COALESCE(parent_id, 0) AS parent_id").
			Where("id IN ?", queue).Scan(&rows)

		parentOf := make(map[int64]int64, len(rows))
		for _, r := range rows {
			parentOf[r.ID] = r.ParentID
		}

		var next []int64
		for _, cid := range queue {
			if _, ok := allIDs[cid]; ok {
				continue
			}
			allIDs[cid] = struct{}{}
			pid := parentOf[cid]
			if pid != 0 {
				if _, ok := allIDs[pid]; !ok {
					next = append(next, pid)
				}
			}
		}
		queue = next
	}

	// 批量查询所有品类的属性
	ids := make([]int64, 0, len(allIDs))
	for id := range allIDs {
		ids = append(ids, id)
	}

	type row struct {
		ID        int64
		Name      string
		SortOrder int
	}
	var rows []row
	if err := r.db.WithContext(ctx).Table("category_attributes").
		Select("a.id, a.name, a.sort_order").
		Joins("JOIN attribute_attributes a ON a.id = category_attributes.attribute_id").
		Where("category_attributes.category_id IN ?", ids).
		Order("a.sort_order ASC, a.id ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	// 内存去重（同属性可能被多个品类关联）
	seen := make(map[int64]struct{}, len(rows))
	attrs := make([]domain.Attribute, 0, len(rows))
	for _, r := range rows {
		if _, ok := seen[r.ID]; ok {
			continue
		}
		seen[r.ID] = struct{}{}
		attrs = append(attrs, domain.Attribute{ID: r.ID, Name: r.Name, SortOrder: r.SortOrder})
	}
	return attrs, nil
}

func (r *CategoryAttributeRepository) SetByCategoryID(ctx context.Context, categoryID int64, attributeIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("category_id = ?", categoryID).Delete(&models.CategoryAttributePO{}).Error; err != nil {
			return err
		}
		if len(attributeIDs) == 0 {
			return nil
		}
		pos := make([]models.CategoryAttributePO, len(attributeIDs))
		for i, attrID := range attributeIDs {
			pos[i] = models.CategoryAttributePO{
				CategoryID:  categoryID,
				AttributeID: attrID,
			}
		}
		return tx.Create(&pos).Error
	})
}

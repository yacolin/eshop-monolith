package repositories

import (
	"context"

	"eshop-monolith/internal/infra/repository/models"
	"eshop-monolith/internal/inventory/api/dto"
	invModels "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/query"
	"eshop-monolith/pkg/utils"

	"gorm.io/gorm"
)

// ProductEnrichedRow 产品 enriched 查询结果行（含分类信息，子查询+LEFT JOIN 单次产出）
type ProductEnrichedRow struct {
	ID           int64           `gorm:"column:id"`
	Name         string          `gorm:"column:name"`
	Description  string          `gorm:"column:description"`
	MinPrice     int64           `gorm:"column:min_price"`
	CreatedAt    utils.Timestamp `gorm:"column:created_at"`
	UpdatedAt    utils.Timestamp `gorm:"column:updated_at"`
	CategoryID   *int64          `gorm:"column:category_id"`
	CategoryName *string         `gorm:"column:category_name"`
}

// Repository 产品仓储接口
type IproductRepository interface {
	// Create 创建产品
	Create(ctx context.Context, product *invModels.Product) error
	// FindByID 根据ID查询产品
	FindByID(ctx context.Context, id int64) (*invModels.Product, error)
	// FindByIDs 根据ID批量查询产品
	FindByIDs(ctx context.Context, ids []int64) ([]invModels.Product, error)
	// FindBySKU 根据SKU查询产品
	FindBySKU(ctx context.Context, sku string) (*invModels.Product, error)
	// ListByCategory 根据分类查询产品
	ListByCategory(ctx context.Context, categoryID int64, page, pageSize int) ([]invModels.Product, int64, error)
	// ListAll 列出所有产品
	ListAll(ctx context.Context, page, pageSize int) ([]invModels.Product, int64, error)
	// Update 更新产品
	Update(ctx context.Context, product *invModels.Product) error
	// Delete 删除产品
	Delete(ctx context.Context, id int64) error

	// ListProducts 列出产品
	ListProducts(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]invModels.Product, error)
	// ListProductsWithTotal 一次查询返回列表+总数（窗口函数）
	ListProductsWithTotal(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]invModels.Product, int64, error)
	// CountProducts 统计产品数量
	CountProducts(ctx context.Context, q dto.ProductListQuery) (int64, error)
	// FindAll 查询所有产品（不分页）
	FindAll(ctx context.Context) ([]invModels.Product, error)
	// ListProductsByCursor 基于游标查询产品列表（深分页优化）
	ListProductsByCursor(ctx context.Context, q dto.ProductCursorQuery, limit int) ([]invModels.Product, error)
	// ListProductsEnriched 列出产品（含分类信息，子查询分页+LEFT JOIN categories 单次查询）
	ListProductsEnriched(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]ProductEnrichedRow, error)
}

// ProductRepository 产品仓储实现
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建产品仓储
func NewProductRepository(db *gorm.DB) IproductRepository {
	return &ProductRepository{db: db}
}

// Create 创建产品
func (r *ProductRepository) Create(ctx context.Context, product *invModels.Product) error {
	po := models.ProductFromDomain(product)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	product.ID = po.ID
	return nil
}

// FindByID 根据ID查询产品
func (r *ProductRepository) FindByID(ctx context.Context, id int64) (*invModels.Product, error) {
	var po models.ProductPO
	err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// FindByIDs 根据ID批量查询产品
func (r *ProductRepository) FindByIDs(ctx context.Context, ids []int64) ([]invModels.Product, error) {
	var pos []models.ProductPO
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	products := make([]invModels.Product, len(pos))
	for i, po := range pos {
		products[i] = *po.ToDomain()
	}
	return products, nil
}

// FindBySKU 根据SKU查询产品
func (r *ProductRepository) FindBySKU(ctx context.Context, sku string) (*invModels.Product, error) {
	var po models.ProductPO
	err := r.db.WithContext(ctx).First(&po, "sku = ?", sku).Error
	if err != nil {
		return nil, err
	}
	return po.ToDomain(), nil
}

// ListByCategory 根据分类查询产品
func (r *ProductRepository) ListByCategory(ctx context.Context, categoryID int64, page, pageSize int) ([]invModels.Product, int64, error) {
	var pos []models.ProductPO
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&models.ProductPO{}).Where("category_id = ?", categoryID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	products := make([]invModels.Product, len(pos))
	for i, po := range pos {
		products[i] = *po.ToDomain()
	}
	return products, total, nil
}

// ListAll 列出所有产品
func (r *ProductRepository) ListAll(ctx context.Context, page, pageSize int) ([]invModels.Product, int64, error) {
	var pos []models.ProductPO
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&models.ProductPO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	products := make([]invModels.Product, len(pos))
	for i, po := range pos {
		products[i] = *po.ToDomain()
	}
	return products, total, nil
}

// Update 更新产品
func (r *ProductRepository) Update(ctx context.Context, product *invModels.Product) error {
	po := models.ProductFromDomain(product)
	return r.db.WithContext(ctx).Save(po).Error
}

// Delete 删除产品
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.ProductPO{}, "id = ?", id).Error
}

// FindAll 查询所有产品（不分页）
func (r *ProductRepository) FindAll(ctx context.Context) ([]invModels.Product, error) {
	var pos []models.ProductPO
	err := r.db.WithContext(ctx).Order("id asc").Find(&pos).Error
	if err != nil {
		return nil, err
	}

	products := make([]invModels.Product, len(pos))
	for i, po := range pos {
		products[i] = *po.ToDomain()
	}
	return products, nil
}

// ListProducts 列出产品（支持查询条件）
func (r *ProductRepository) ListProducts(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]invModels.Product, error) {
	var pos []models.ProductPO
	db := r.applyQueryConditions(ctx, q)
	db = r.applyOrder(db, q)

	// 列表查询不读取 TEXT description 字段，减少 I/O
	err := db.Select("id, name, min_price, created_at, updated_at").Offset(offset).Limit(limit).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	products := make([]invModels.Product, len(pos))
	for i, po := range pos {
		products[i] = *po.ToDomain()
	}
	return products, nil
}

// CountProducts 统计产品数量
func (r *ProductRepository) CountProducts(ctx context.Context, q dto.ProductListQuery) (int64, error) {
	var total int64
	db := r.applyQueryConditions(ctx, q)

	// 执行统计（不需要排序）
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// productWithTotalRow 标量子查询结果行
type productWithTotalRow struct {
	models.ProductPO
	TotalCount int64 `gorm:"column:total_count"`
}

// ListProductsWithTotal 单次查询：COUNT 用 GORM 子查询生成，复用过滤条件
func (r *ProductRepository) ListProductsWithTotal(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]invModels.Product, int64, error) {
	var rows []productWithTotalRow

	// 子查询：复用相同过滤条件生成 COUNT
	countDB := r.db.WithContext(ctx).Model(&models.ProductPO{})
	countDB = applyQueryConditionsOnDB(countDB, q)

	dataDB := r.applyQueryConditions(ctx, q)
	dataDB = r.applyOrder(dataDB, q)

	err := dataDB.Select("id, name, min_price, created_at, updated_at, (?) AS total_count", countDB.Select("COUNT(*)")).
		Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	total := int64(0)
	if len(rows) > 0 {
		total = rows[0].TotalCount
	}
	products := make([]invModels.Product, len(rows))
	for i := range rows {
		products[i] = *rows[i].ProductPO.ToDomain()
	}
	return products, total, nil
}

// ListProductsByCursor 基于游标查询产品列表（深分页优化）
func (r *ProductRepository) ListProductsByCursor(ctx context.Context, q dto.ProductCursorQuery, limit int) ([]invModels.Product, error) {
	var pos []models.ProductPO
	db := r.applyCursorQueryConditions(ctx, q)
	if q.Cursor > 0 {
		db = db.Where("id > ?", q.Cursor)
	}
	err := db.Order("id asc").Limit(limit).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	products := make([]invModels.Product, len(pos))
	for i, po := range pos {
		products[i] = *po.ToDomain()
	}
	return products, nil
}

// ListProductsEnriched 列出产品（含分类信息，子查询先分页再 LEFT JOIN categories 单次查询）
func (r *ProductRepository) ListProductsEnriched(ctx context.Context, q dto.ProductListQuery, offset, limit int) ([]ProductEnrichedRow, error) {
	subQuery := r.db.WithContext(ctx).Model(&models.ProductPO{})
	subQuery = r.applyQueryConditions(ctx, q)
	subQuery = r.applyOrder(subQuery, q)
	subQuery = subQuery.Offset(offset).Limit(limit)
	if offset > 0 {
		subQuery = subQuery.Select("*")
	}

	var rows []ProductEnrichedRow
	err := r.db.WithContext(ctx).
		Table("(?) AS p", subQuery).
		Select("p.id, p.name, p.description, p.min_price, p.created_at, p.updated_at, c.id AS category_id, c.name AS category_name").
		Joins("LEFT JOIN product_categories pc ON p.id = pc.product_id").
		Joins("LEFT JOIN categories c ON pc.category_id = c.id").
		Order("p.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// applyCursorQueryConditions 应用游标查询的过滤条件
func (r *ProductRepository) applyCursorQueryConditions(ctx context.Context, q dto.ProductCursorQuery) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&models.ProductPO{})
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.SKU != "" {
		db = db.Where("sku = ?", q.SKU)
	}
	if q.CategoryID != nil {
		db = db.Where("id IN (?)",
			r.db.Table("product_categories").Select("product_id").Where("category_id = ?", *q.CategoryID),
		)
	}
	return db
}

// applyQueryConditions 应用查询条件（不包含排序）
func (r *ProductRepository) applyQueryConditions(ctx context.Context, q dto.ProductListQuery) *gorm.DB {
	return applyQueryConditionsOnDB(r.db.WithContext(ctx), q)
}

// applyQueryConditionsOnDB 在指定 DB 上应用查询条件（不绑定 repo）
func applyQueryConditionsOnDB(db *gorm.DB, q dto.ProductListQuery) *gorm.DB {
	db = db.Model(&models.ProductPO{})
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.SKU != "" {
		db = db.Where("sku = ?", q.SKU)
	}
	if q.CategoryID != nil {
		db = db.Where("id IN (?)",
			db.Session(&gorm.Session{NewDB: true}).Table("product_categories").Select("product_id").Where("category_id = ?", *q.CategoryID),
		)
	}
	return db
}

// applyOrder 应用排序
func (r *ProductRepository) applyOrder(db *gorm.DB, q dto.ProductListQuery) *gorm.DB {
	return query.ApplyOrder(db, q.SortBy, q.Order, "id asc")
}

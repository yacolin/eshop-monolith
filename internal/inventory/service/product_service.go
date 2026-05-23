package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"eshop-monolith/internal/infra/domain/shared"
	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/events"
	"eshop-monolith/pkg/errcode"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	productCacheZSet     = "product:zset"
	productInfoPrefix    = "product:info:"
	cachedProductListTTL = time.Hour
)

// zrangeMGetScript 一次网络往返完成 ZRANGE + MGET
var zrangeMGetScript = redis.NewScript(`
local ids
if ARGV[3] == "desc" then
	ids = redis.call("ZREVRANGE", KEYS[1], ARGV[1], ARGV[2])
else
	ids = redis.call("ZRANGE", KEYS[1], ARGV[1], ARGV[2])
end
if #ids == 0 then return {} end
local keys = {}
for i, id in ipairs(ids) do
	keys[i] = ARGV[4] .. id
end
return redis.call("MGET", unpack(keys))
`)

// ProductService 产品服务
type ProductService struct {
	repo          repositories.IproductRepository
	inventoryRepo repositories.IinventoryRepository
	bus           *eventbus.Bus
	db            *gorm.DB
	rdb           *redis.Client
}

// NewProductService 创建产品服务
func NewProductService(
	repo repositories.IproductRepository,
	inventoryRepo repositories.IinventoryRepository,
	bus *eventbus.Bus,
	db *gorm.DB,
	rdb *redis.Client,
) *ProductService {
	return &ProductService{
		repo:          repo,
		inventoryRepo: inventoryRepo,
		bus:           bus,
		db:            db,
		rdb:           rdb,
	}
}

// WarmupProductCache 将全量商品写入 Redis zset，每个商品独立 key 存储
func (s *ProductService) WarmupProductCache(ctx context.Context) (int, error) {
	products, err := s.repo.FindAll(ctx)
	if err != nil {
		return 0, err
	}

	pipe := s.rdb.Pipeline()
	ctxBg := context.Background()

	pipe.Del(ctxBg, productCacheZSet)

	for _, p := range products {
		data, err := sonic.Marshal(p)
		if err != nil {
			return 0, err
		}
		pipe.Set(ctxBg, productInfoPrefix+strconv.FormatInt(p.ID, 10), data, cachedProductListTTL)
		pipe.ZAdd(ctxBg, productCacheZSet, redis.Z{
			Score:  float64(p.ID),
			Member: p.ID,
		})
	}

	pipe.Expire(ctxBg, productCacheZSet, cachedProductListTTL)

	_, err = pipe.Exec(ctxBg)
	if err != nil {
		return 0, err
	}

	return len(products), nil
}

// ListCachedProducts 从 Redis zset 中分页读取商品列表（Lua 脚本一次网络往返）
func (s *ProductService) ListCachedProducts(ctx context.Context, q dto.ProductListQuery) (*dto.ProductListResult, error) {
	ctxBg := context.Background()
	offset := int64((q.Page - 1) * q.Size)
	stop := offset + int64(q.Size) - 1

	order := "asc"
	if q.Order == "desc" {
		order = "desc"
	}

	dataList, err := zrangeMGetScript.Run(ctxBg, s.rdb, []string{productCacheZSet}, offset, stop, order, productInfoPrefix).Result()
	if err != nil {
		return nil, err
	}

	items, ok := dataList.([]any)
	if !ok {
		items = nil
	}

	products := make([]models.Product, 0, len(items))
	for _, data := range items {
		if data == nil {
			continue
		}
		var p models.Product
		if err := sonic.Unmarshal([]byte(data.(string)), &p); err != nil {
			continue
		}
		products = append(products, p)
	}

	total, err := s.rdb.ZCard(ctxBg, productCacheZSet).Result()
	if err != nil {
		return nil, err
	}

	return &dto.ProductListResult{List: products, Total: total}, nil
}

// GetCachedProductByID 从 Redis 缓存中查询单个商品
func (s *ProductService) GetCachedProductByID(ctx context.Context, id int64) (*models.Product, error) {
	data, err := s.rdb.Get(context.Background(), productInfoPrefix+strconv.FormatInt(id, 10)).Bytes()
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err := sonic.Unmarshal(data, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

// GetProductWithInventory 获取产品详情（聚合库存信息），使用 goroutine + channel 并发查询以降低延迟
func (s *ProductService) GetProductWithInventory(ctx context.Context, id int64) (*dto.ProductDetailDTO, error) {
	type prodResult struct {
		product *models.Product
		err     error
	}
	type invResult struct {
		inventory *models.Inventory
		err       error
	}

	prodCh := make(chan prodResult, 1)
	invCh := make(chan invResult, 1)

	go func() {
		p, err := s.repo.FindByID(ctx, id)
		prodCh <- prodResult{product: p, err: err}
	}()

	go func() {
		inv, err := s.inventoryRepo.FindInventoryByProductID(ctx, id)
		invCh <- invResult{inventory: inv, err: err}
	}()

	pr := <-prodCh
	if pr.err != nil {
		return nil, pr.err
	}

	ir := <-invCh
	detail := &dto.ProductDetailDTO{
		ID:          pr.product.ID,
		Name:        pr.product.Name,
		Description: pr.product.Description,
		Price:       pr.product.Price,
		SKU:         pr.product.SKU,
		CreatedAt:   pr.product.CreatedAt,
		UpdatedAt:   pr.product.UpdatedAt,
	}
	if ir.err == nil && ir.inventory != nil {
		detail.Quantity = ir.inventory.Quantity
		detail.Status = ir.inventory.Status
		detail.Reserved = ir.inventory.Reserved
		detail.Threshold = ir.inventory.Threshold
	}

	return detail, nil
}

// CreateProduct 创建产品
func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductDTO) (*models.Product, error) {
	newProduct := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newProduct).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				return errcode.ErrDuplicateSKU
			}
			return err
		}
		for _, categoryID := range req.CategoryIDs {
			pc := &shared.ProductCategory{
				ProductID:  newProduct.ID,
				CategoryID: categoryID,
			}
			if err := tx.Create(pc).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	categoryIDValue := int64(0)
	if len(req.CategoryIDs) > 0 {
		categoryIDValue = req.CategoryIDs[0]
	}
	s.bus.Publish(events.ProductCreatedEvent{
		ProductID:  newProduct.ID,
		Name:       newProduct.Name,
		Price:      newProduct.Price,
		CategoryID: categoryIDValue,
	})

	return newProduct, nil
}

// GetProductByID 根据ID获取产品
func (s *ProductService) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.FindByID(ctx, id)
}

// GetProductBySKU 根据SKU获取产品
func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (*models.Product, error) {
	return s.repo.FindBySKU(ctx, sku)
}

// GetProductWithCategories 获取产品及其关联的分类
func (s *ProductService) GetProductWithCategories(ctx context.Context, productID int64) (*models.Product, []models.Category, error) {
	prod, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return nil, nil, err
	}

	var categories []models.Category
	if err := s.db.WithContext(ctx).Table("categories").
		Joins("JOIN product_categories ON categories.id = product_categories.category_id").
		Where("product_categories.product_id = ?", productID).
		Find(&categories).Error; err != nil {
		return nil, nil, err
	}

	return prod, categories, nil
}

// ListProductsByCategory 根据分类列出产品
func (s *ProductService) ListProductsByCategory(ctx context.Context, categoryID int64, page, pageSize int) ([]models.Product, int64, error) {
	return s.repo.ListByCategory(ctx, categoryID, page, pageSize)
}

// UpdateProduct 更新产品
func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req *dto.UpdateProductDTO) (*models.Product, error) {
	existingProduct, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existingProduct.Name = *req.Name
	}
	if req.Description != nil {
		existingProduct.Description = *req.Description
	}
	if req.Price != nil {
		existingProduct.Price = *req.Price
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(existingProduct).Error; err != nil {
			return err
		}
		if req.CategoryIDs != nil {
			if err := tx.Where("product_id = ?", id).Delete(&shared.ProductCategory{}).Error; err != nil {
				return err
			}
			for _, categoryID := range req.CategoryIDs {
				pc := &shared.ProductCategory{
					ProductID:  existingProduct.ID,
					CategoryID: categoryID,
				}
				if err := tx.Create(pc).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	categoryIDValue := int64(0)
	if len(req.CategoryIDs) > 0 {
		categoryIDValue = req.CategoryIDs[0]
	}
	s.bus.Publish(events.ProductUpdatedEvent{
		ProductID:  existingProduct.ID,
		Name:       existingProduct.Name,
		Price:      existingProduct.Price,
		CategoryID: categoryIDValue,
	})

	return existingProduct, nil
}

// DeleteProduct 删除产品
func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	existingProduct, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", id).Delete(&shared.ProductCategory{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Product{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	s.bus.Publish(events.ProductDeletedEvent{
		ProductID:  existingProduct.ID,
		Name:       existingProduct.Name,
		CategoryID: 0,
	})

	return nil
}

// ListAllProducts 列出所有产品
func (s *ProductService) ListAllProducts(ctx context.Context, page, pageSize int) ([]models.Product, int64, error) {
	return s.repo.ListAll(ctx, page, pageSize)
}

func (s *ProductService) ListProducts(ctx context.Context, q dto.ProductListQuery) (*dto.ProductListResult, error) {
	offset := (q.Page - 1) * q.Size
	list, err := s.repo.ListProducts(ctx, q, offset, q.Size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountProducts(ctx, q)
	if err != nil {
		return nil, err
	}

	return &dto.ProductListResult{
		List:  list,
		Total: total,
	}, nil
}

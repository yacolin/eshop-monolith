package handlers

import (
	"strconv"

	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ProductHandler 产品处理器
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler 创建产品处理器
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// ListProducts 列出所有产品
// @Summary 列出所有产品
// @Description 获取所有产品的列表
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Param name query string false "产品名称模糊搜索"
// @Param sku query string false "SKU精确搜索"
// @Success 200 {object} response.Response{data=dto.ProductListResult}
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var q dto.ProductListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	(&q).Normalize()

	result, err := h.productService.ListProducts(c, q)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// ListProductsWithCategory 列出所有产品（含分类信息）
// @Summary 列出所有产品（含分类）
// @Description 获取所有产品的列表，每个产品附带其首个分类信息，通过一次批量查询补全
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Param name query string false "产品名称模糊搜索"
// @Param sku query string false "SKU精确搜索"
// @Success 200 {object} response.Response{data=dto.ProductWithCategoryListResult}
// @Router /api/v1/products/enriched [get]
func (h *ProductHandler) ListProductsWithCategory(c *gin.Context) {
	var q dto.ProductListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	(&q).Normalize()

	result, err := h.productService.ListProductsWithCategory(c, q)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// GetProduct 根据ID获取产品
// @Summary 获取产品详情
// @Description 根据产品ID获取产品详情
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response{data=models.Product}
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	product, err := h.productService.GetProductByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, product)
}

// GetProductWithCategory 获取产品详情（含分类信息）
// @Summary 获取产品详情（含分类）
// @Description 根据产品ID获取产品详情，包含产品信息和首个分类信息
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response{data=dto.ProductWithCategoryDTO}
// @Router /api/v1/products/{id}/enriched [get]
func (h *ProductHandler) GetProductWithCategory(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	product, err := h.productService.GetProductWithCategory(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, product)
}

// GetProductDetail 获取产品详情（聚合库存信息）
// @Summary 获取产品详情（含库存）
// @Description 根据产品ID获取产品详情，包含产品信息和库存信息
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response{data=dto.ProductDetailDTO}
// @Router /api/v1/products/{id}/detail [get]
func (h *ProductHandler) GetProductDetail(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	detail, err := h.productService.GetProductWithInventory(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, detail)
}

// ListProductsByCategory 根据分类获取产品
// @Summary 根据分类获取产品列表
// @Description 根据分类ID获取产品列表
// @Tags products
// @Accept json
// @Produce json
// @Param category_id path int true "分类ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} response.Response{data=[]models.Product}
// @Router /api/v1/products/category/{category_id} [get]
func (h *ProductHandler) ListProductsByCategory(c *gin.Context) {
	categoryID, err := utils.ParseIntParam(c, "category_id")
	if err != nil {
		c.Error(err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	products, total, err := h.productService.ListProductsByCategory(c, categoryID, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"list":  products,
		"total": total,
	})
}

// CreateProduct 创建产品
// @Summary 创建产品
// @Description 创建一个新的产品
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductDTO true "产品信息"
// @Success 200 {object} response.Response{data=models.Product}
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	product, err := h.productService.CreateProduct(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, product)
}

// UpdateProduct 更新产品
// @Summary 更新产品
// @Description 根据ID更新产品信息
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Param product body dto.UpdateProductDTO true "产品信息"
// @Success 200 {object} response.Response{data=models.Product}
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.UpdateProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	product, err := h.productService.UpdateProduct(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, product)
}

// DeleteProduct 删除产品
// @Summary 删除产品
// @Description 根据ID删除产品
// @Tags products
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.productService.DeleteProduct(c, id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

// WarmupCache 将全量商品预热到 Redis 缓存
// @Summary 预热商品缓存
// @Description 将全量商品数据加载到 Redis 缓存中
// @Tags products
// @Produce json
// @Success 200 {object} response.Response{data=map[string]int}
// @Router /api/v1/products/cache/warmup [post]
func (h *ProductHandler) WarmupCache(c *gin.Context) {
	total, err := h.productService.WarmupProductCache(c)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"total": total})
}

// ListCachedProducts 从缓存中获取产品列表
// @Summary 从缓存获取产品列表
// @Description 从 Redis 缓存中读取产品列表，支持分页和排序
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Param sort_by query string false "排序字段 (id, name, price)"
// @Param order query string false "排序方向 (asc, desc)" default(asc)
// @Success 200 {object} response.Response{data=dto.ProductListResult}
// @Router /api/v1/products/cache [get]
func (h *ProductHandler) ListCachedProducts(c *gin.Context) {
	var q dto.ProductListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	(&q).Normalize()

	result, err := h.productService.ListCachedProducts(c, q)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// GetCachedProduct 从缓存中查询单个商品
// @Summary 从缓存查询商品
// @Description 从 Redis 缓存中根据 ID 查询单个商品
// @Tags products
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response{data=models.Product}
// @Router /api/v1/products/cache/{id} [get]
func (h *ProductHandler) GetCachedProduct(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	product, err := h.productService.GetCachedProductByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, product)
}

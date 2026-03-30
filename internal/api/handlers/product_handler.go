package handlers

import (
	"strconv"

	"eshop-monolith/internal/domain/product"
	"eshop-monolith/internal/pkg/response"
	"eshop-monolith/internal/service"

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
// @Tags 产品管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=product.ProductListResult}
// @Failure 500 {object} response.Response{error=string}
// @Router /api/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var q product.ProductListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	// normalize pagination values (ensure page>=1, 1<=size<=100)
	(&q).Normalize()

	result, err := h.productService.ListProducts(c, q)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// GetProduct 根据ID获取产品
// @Summary 根据ID获取产品
// @Description 根据产品ID获取产品详情
// @Tags 产品管理
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response{data=product.Product}
// @Failure 400 {object} response.Response{error=string}
// @Failure 500 {object} response.Response{error=string}
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		response.BindError(c, err)
		return
	}

	product, err := h.productService.GetProductByID(c, id)
	if err != nil {
		response.SysError(c, err)
		return
	}

	response.Success(c, product)
}

// ListProductsByCategory 根据分类获取产品
// @Summary 根据分类获取产品
// @Description 根据分类ID获取产品列表
// @Tags 产品管理
// @Accept json
// @Produce json
// @Param category_id path int true "分类ID"
// @Success 200 {object} response.Response{data=[]product.Product}
// @Failure 400 {object} response.Response{error=string}
// @Failure 500 {object} response.Response{error=string}
// @Router /api/products/category/{category_id} [get]
func (h *ProductHandler) ListProductsByCategory(c *gin.Context) {
	categoryID, err := parseIntParam(c, "category_id")
	if err != nil {
		response.BindError(c, err)
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	products, total, err := h.productService.ListProductsByCategory(c, categoryID, page, pageSize)
	if err != nil {
		response.SysError(c, err)
		return
	}

	response.Success(c, gin.H{
		"list":  products,
		"total": total,
	})
}

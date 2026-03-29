package handlers

import (
	"strconv"

	"eshop-monolith/internal/pkg/response"
	"eshop-monolith/internal/service"

	"github.com/gin-gonic/gin"
)

// parseIntParam 解析路径参数为int64
func parseIntParam(c *gin.Context, param string) (int64, error) {
	idStr := c.Param(param)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// ListCategories 列出所有分类
// @Summary 列出所有分类
// @Description 获取所有分类的列表
// @Tags 分类管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]category.Category}
// @Failure 500 {object} response.Response{error=string}
// @Router /api/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.categoryService.ListAllCategories(c)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, categories)
}

// ListRootCategories 列出根分类
// @Summary 列出根分类
// @Description 获取所有根分类的列表
// @Tags 分类管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]category.Category}
// @Failure 500 {object} response.Response{error=string}
// @Router /api/categories/root [get]
func (h *CategoryHandler) ListRootCategories(c *gin.Context) {
	categories, err := h.categoryService.ListRootCategories(c)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, categories)
}

// ListSubCategories 列出子分类
// @Summary 列出子分类
// @Description 根据父分类ID获取子分类列表
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param parent_id path int true "父分类ID"
// @Success 200 {object} response.Response{data=[]category.Category}
// @Failure 400 {object} response.Response{error=string}
// @Failure 500 {object} response.Response{error=string}
// @Router /api/categories/{parent_id}/children [get]
func (h *CategoryHandler) ListSubCategories(c *gin.Context) {
	parentID, err := parseIntParam(c, "parent_id")
	if err != nil {
		c.Error(err)
		return
	}

	categories, err := h.categoryService.ListSubCategories(c, parentID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, categories)
}

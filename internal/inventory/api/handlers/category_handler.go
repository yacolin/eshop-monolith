package handlers

import (
	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/service"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

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
// @Success 200 {object} response.Response{data=dto.CategoryListResult}
// @Failure 500 {object} response.Response{error=string}
// @Router /api/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	var q dto.CategoryListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

	// normalize pagination values (ensure page>=1, 1<=size<=100)
	(&q).Normalize()

	result, err := h.categoryService.ListCategories(c, q)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// ListRootCategories 列出根分类
// @Summary 列出根分类
// @Description 获取所有根分类的列表
// @Tags 分类管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]models.Category}
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
// @Param id path int true "父分类ID"
// @Success 200 {object} response.Response{data=[]models.Category}
// @Failure 400 {object} response.Response{error=string}
// @Failure 500 {object} response.Response{error=string}
// @Router /api/categories/{id}/children [get]
func (h *CategoryHandler) ListSubCategories(c *gin.Context) {
	parentID, err := utils.ParseIntParam(c, "id")
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

// CreateCategory 创建分类
// @Summary 创建分类
// @Description 创建一个新的分类
// @Tags 分类
// @Accept json
// @Produce json
// @Param models body dto.CreateCategoryDTO true "分类信息"
// @Success 200 {object} models.Category "成功"
// @Router /inventory/api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	result, err := h.categoryService.CreateCategory(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// GetCategoryByID 获取分类详情
// @Summary 获取分类详情
// @Description 根据ID获取分类详细信息
// @Tags 分类
// @Produce json
// @Param id path string true "分类ID"
// @Success 200 {object} models.Category "成功"
// @Router /inventory/api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	models, err := h.categoryService.GetCategoryByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, models)
}

// UpdateCategory 更新分类
// @Summary 更新分类
// @Description 根据ID更新分类信息
// @Tags 分类
// @Accept json
// @Produce json
// @Param id path string true "分类ID"
// @Param models body dto.UpdateCategoryDTO true "分类信息"
// @Success 200 {object} models.Category "成功"
// @Router /inventory/api/v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateCategoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	models, err := h.categoryService.UpdateCategory(c, id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, models)
}

// DeleteCategory 删除分类
// @Summary 删除分类
// @Description 根据ID删除分类
// @Tags 分类
// @Produce json
// @Param id path string true "分类ID"
// @Success 200 {object} map[string]string "成功"
// @Router /inventory/api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.categoryService.DeleteCategory(c, id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

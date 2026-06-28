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
// @Description 获取所有分类的列表，支持分页
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Param name query string false "分类名称模糊搜索"
// @Success 200 {object} response.Response{data=dto.CategoryListResult}
// @Router /api/v1/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	var q dto.CategoryListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(err)
		return
	}

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
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=dto.CategoryListResult}
// @Router /api/v1/categories/root [get]
func (h *CategoryHandler) ListRootCategories(c *gin.Context) {
	result, err := h.categoryService.ListRootCategories(c)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// ListNonRootCategories 从缓存列出非根分类
// @Summary 从缓存列出非根分类
// @Description 从 Redis 缓存中读取所有非根分类列表（仅返回 id 和 name）
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=dto.CachedCategoryListResult}
// @Router /api/v1/categories/non-root [get]
func (h *CategoryHandler) ListNonRootCategories(c *gin.Context) {
	result, err := h.categoryService.ListNonRootCategories(c)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ListCachedCategories 从缓存列出全部分类
// @Summary 从缓存列出全部分类
// @Description 从 Redis 缓存中读取全部分类列表，仅返回 id 和 name
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.CachedCategoryItem}
// @Router /api/v1/categories/cache [get]
func (h *CategoryHandler) ListCachedCategories(c *gin.Context) {
	items, err := h.categoryService.ListCachedCategories(c)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, items)
}

// GetCachedCategory 从缓存查询单个分类
// @Summary 从缓存查询分类
// @Description 从 Redis 缓存中根据 ID 查询单个分类，仅返回 id 和 name
// @Tags categories
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response{data=dto.CachedCategoryItem}
// @Router /api/v1/categories/cache/{id} [get]
func (h *CategoryHandler) GetCachedCategory(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	item, err := h.categoryService.GetCachedCategoryByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, item)
}

// WarmupCategoryCache 预热分类缓存
// @Summary 预热分类缓存
// @Description 将全部分类数据加载到 Redis 缓存中
// @Tags categories
// @Produce json
// @Success 200 {object} response.Response{data=map[string]int}
// @Router /api/v1/categories/cache/warmup [post]
func (h *CategoryHandler) WarmupCategoryCache(c *gin.Context) {
	total, err := h.categoryService.WarmupCategoryCache(c)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, gin.H{"total": total})
}

// ListSubCategories 列出子分类
// @Summary 列出子分类
// @Description 根据父分类ID获取子分类列表
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "父分类ID"
// @Success 200 {object} response.Response{data=[]models.Category}
// @Router /api/v1/categories/{id}/children [get]
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
// @Tags categories
// @Accept json
// @Produce json
// @Param category body dto.CreateCategoryDTO true "分类信息"
// @Success 200 {object} response.Response{data=models.Category}
// @Router /api/v1/categories [post]
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
// @Tags categories
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response{data=models.Category}
// @Router /api/v1/categories/{id} [get]
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
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param category body dto.UpdateCategoryDTO true "分类信息"
// @Success 200 {object} response.Response{data=models.Category}
// @Router /api/v1/categories/{id} [put]
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
// @Tags categories
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/categories/{id} [delete]
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


// ── 品类-属性关联 ──────────────────────────────────────────────

// GetCategoryAttributes 获取品类关联的属性列表
// @Summary 获取品类关联的属性列表
// @Description 获取指定品类关联的规格属性维度列表
// @Tags categories
// @Produce json
// @Param id path int true "品类ID"
// @Success 200 {object} response.Response{data=[]dto.CategoryAttributeResponse}
// @Router /api/v1/categories/{id}/attributes [get]
func (h *CategoryHandler) GetCategoryAttributes(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	attrs, err := h.categoryService.GetCategoryAttributes(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, attrs)
}

// SetCategoryAttributes 设置品类关联的属性（全量替换）
// @Summary 设置品类关联的属性
// @Description 全量替换指定品类关联的规格属性，原有关联会被清除
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "品类ID"
// @Param attrs body dto.SetCategoryAttributesDTO true "属性ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/categories/{id}/attributes [put]
func (h *CategoryHandler) SetCategoryAttributes(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.SetCategoryAttributesDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.categoryService.SetCategoryAttributes(c, id, &req); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

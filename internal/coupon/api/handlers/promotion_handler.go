package handlers

import (
	"encoding/json"
	"strconv"
	"time"

	"eshop-monolith/internal/coupon/api/dto"
	couponModels "eshop-monolith/internal/coupon/domain/models"
	"eshop-monolith/internal/coupon/service"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PromotionHandler struct {
	promotionService *service.PromotionService
}

func NewPromotionHandler(promotionService *service.PromotionService) *PromotionHandler {
	return &PromotionHandler{promotionService: promotionService}
}

// CreatePromotion 创建促销活动
// @Summary 创建促销活动
// @Description 创建新的促销活动（限时折扣/满减活动）
// @Tags promotions
// @Accept json
// @Produce json
// @Param promotion body dto.CreatePromotionDTO true "促销活动信息"
// @Success 200 {object} response.Response{data=dto.PromotionResponse}
// @Router /api/v1/admin/promotions [post]
func (h *PromotionHandler) CreatePromotion(c *gin.Context) {
	var req dto.CreatePromotionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		c.Error(err)
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		c.Error(err)
		return
	}

	var ruleObj any
	if err := json.Unmarshal([]byte(req.Rule), &ruleObj); err != nil {
		c.Error(err)
		return
	}

	promotion := &couponModels.Promotion{
		Name:        req.Name,
		Description: req.Description,
		PromoType:   couponModels.PromotionType(req.PromoType),
		Scope:       req.Scope,
		ScopeValue:  req.ScopeValue,
		Rule:        req.Rule,
		StartTime:   utils.Timestamp(startTime),
		EndTime:     utils.Timestamp(endTime),
		SortOrder:   req.SortOrder,
	}

	if err := h.promotionService.CreatePromotion(c, promotion); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, toPromotionResponse(promotion))
}

// UpdatePromotion 更新促销活动
// @Summary 更新促销活动
// @Description 更新指定促销活动的信息
// @Tags promotions
// @Accept json
// @Produce json
// @Param id path int true "促销活动ID"
// @Param promotion body dto.UpdatePromotionDTO true "更新信息"
// @Success 200 {object} response.Response{data=dto.PromotionResponse}
// @Router /api/v1/admin/promotions/{id} [put]
func (h *PromotionHandler) UpdatePromotion(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdatePromotionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	existing, err := h.promotionService.GetPromotion(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Scope != "" {
		existing.Scope = req.Scope
	}
	if req.ScopeValue != "" {
		existing.ScopeValue = req.ScopeValue
	}
	if req.Rule != "" {
		var ruleObj any
		if parseErr := json.Unmarshal([]byte(req.Rule), &ruleObj); parseErr != nil {
			c.Error(parseErr)
			return
		}
		existing.Rule = req.Rule
	}
	if req.StartTime != "" {
		t, parseErr := time.Parse("2006-01-02 15:04:05", req.StartTime)
		if parseErr != nil {
			c.Error(parseErr)
			return
		}
		existing.StartTime = utils.Timestamp(t)
	}
	if req.EndTime != "" {
		t, parseErr := time.Parse("2006-01-02 15:04:05", req.EndTime)
		if parseErr != nil {
			c.Error(parseErr)
			return
		}
		existing.EndTime = utils.Timestamp(t)
	}
	if req.Status != "" {
		existing.Status = couponModels.PromotionStatus(req.Status)
	}
	existing.SortOrder = req.SortOrder

	if err := h.promotionService.UpdatePromotion(c, existing); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, toPromotionResponse(existing))
}

// GetPromotion 获取促销活动详情
// @Summary 获取促销活动详情
// @Description 根据ID获取促销活动详情
// @Tags promotions
// @Accept json
// @Produce json
// @Param id path int true "促销活动ID"
// @Success 200 {object} response.Response{data=dto.PromotionResponse}
// @Router /api/v1/promotions/{id} [get]
func (h *PromotionHandler) GetPromotion(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	promotion, err := h.promotionService.GetPromotion(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, toPromotionResponse(promotion))
}

// ListPromotions 促销活动列表
// @Summary 促销活动列表
// @Description 分页查询促销活动列表
// @Tags promotions
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(20)
// @Success 200 {object} response.Response{data=dto.PromotionListResult}
// @Router /api/v1/promotions [get]
func (h *PromotionHandler) ListPromotions(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 64)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	list, total, err := h.promotionService.ListPromotions(c, int(page), int(pageSize))
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]dto.PromotionResponse, len(list))
	for i, p := range list {
		items[i] = toPromotionResponse(&p)
	}

	response.Success(c, dto.PromotionListResult{
		Total: total,
		List:  items,
	})
}

// GetActivePromotions 获取当前活动
// @Summary 获取进行中的促销活动
// @Description 获取所有当前正在进行的促销活动
// @Tags promotions
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.PromotionResponse}
// @Router /api/v1/promotions/active [get]
func (h *PromotionHandler) GetActivePromotions(c *gin.Context) {
	activeList, err := h.promotionService.GetActivePromotions(c)
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]dto.PromotionResponse, len(activeList))
	for i, p := range activeList {
		items[i] = toPromotionResponse(&p)
	}

	response.Success(c, items)
}

// UpdatePromotionStatus 更新促销活动状态
// @Summary 更新促销活动状态
// @Description 更新指定促销活动的状态（pending/active/finished/cancelled）
// @Tags promotions
// @Accept json
// @Produce json
// @Param id path int true "促销活动ID"
// @Param status body object{status=string} true "目标状态"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/admin/promotions/{id}/status [put]
func (h *PromotionHandler) UpdatePromotionStatus(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending active finished cancelled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.promotionService.UpdatePromotionStatus(c, id, couponModels.PromotionStatus(req.Status)); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "status updated"})
}

// LinkProducts 关联促销商品
// @Summary 关联促销商品
// @Description 为促销活动关联指定商品，会覆盖原有商品列表
// @Tags promotions
// @Accept json
// @Produce json
// @Param id path int true "促销活动ID"
// @Param products body dto.LinkProductsDTO true "商品关联信息"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/admin/promotions/{id}/products [post]
func (h *PromotionHandler) LinkProducts(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.LinkProductsDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.promotionService.LinkProducts(c, id, req.ProductIDs, req.Discount); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "products linked"})
}

// GetPromotionProducts 获取促销关联商品
// @Summary 获取促销关联商品
// @Description 获取指定促销活动关联的所有商品
// @Tags promotions
// @Accept json
// @Produce json
// @Param id path int true "促销活动ID"
// @Success 200 {object} response.Response{data=[]object}
// @Router /api/v1/promotions/{id}/products [get]
func (h *PromotionHandler) GetPromotionProducts(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	products, err := h.promotionService.GetPromotionProducts(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, products)
}

// toPromotionResponse 转换促销活动到响应
func toPromotionResponse(p *couponModels.Promotion) dto.PromotionResponse {
	return dto.PromotionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		PromoType:   string(p.PromoType),
		Scope:       p.Scope,
		ScopeValue:  p.ScopeValue,
		Rule:        p.Rule,
		StartTime:   p.StartTime,
		EndTime:     p.EndTime,
		Status:      string(p.Status),
		SortOrder:   p.SortOrder,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

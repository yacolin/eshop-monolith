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
func (h *PromotionHandler) CreatePromotion(c *gin.Context) {
	var req dto.CreatePromotionReq
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

	// 验证规则格式
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
func (h *PromotionHandler) UpdatePromotion(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdatePromotionReq
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
func (h *PromotionHandler) LinkProducts(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.LinkProductsReq
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

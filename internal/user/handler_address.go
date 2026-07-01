package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"
)

type AddressHandler struct {
	svc *AddressService
}

func NewAddressHandler(svc *AddressService) *AddressHandler {
	return &AddressHandler{svc: svc}
}

// currentUserID 从 JWT 上下文中提取当前用户 ID
func currentUserID(c *gin.Context) int64 {
	v, _ := c.Get("user_id")
	switch id := v.(type) {
	case int64:
		return id
	case uint:
		return int64(id)
	case float64:
		return int64(id)
	case int:
		return int64(id)
	}
	return 0
}

// Create 创建地址
// @Summary 创建地址
// @Tags addresses
// @Tags frontend
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreateAddressReq true "地址信息"
// @Success 200 {object} response.Response{data=Address}
// @Router /api/v1/addresses [post]
func (h *AddressHandler) Create(c *gin.Context) {
	var req CreateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Create(c, currentUserID(c), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// List 地址列表
// @Summary 地址列表
// @Tags addresses
// @Tags frontend
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=AddressListResult}
// @Router /api/v1/addresses [get]
func (h *AddressHandler) List(c *gin.Context) {
	result, err := h.svc.ListByUser(c, currentUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// GetByID 获取地址详情
// @Summary 获取地址详情
// @Tags addresses
// @Tags frontend
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "地址ID"
// @Success 200 {object} response.Response{data=Address}
// @Router /api/v1/addresses/{id} [get]
func (h *AddressHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.GetByID(c, currentUserID(c), id)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Update 更新地址
// @Summary 更新地址
// @Tags addresses
// @Tags frontend
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "地址ID"
// @Param request body UpdateAddressReq true "更新信息"
// @Success 200 {object} response.Response{data=Address}
// @Router /api/v1/addresses/{id} [put]
func (h *AddressHandler) Update(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	var req UpdateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Update(c, currentUserID(c), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Delete 删除地址
// @Summary 删除地址
// @Tags addresses
// @Tags frontend
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "地址ID"
// @Success 200 {object} response.Response
// @Router /api/v1/addresses/{id} [delete]
func (h *AddressHandler) Delete(c *gin.Context) {
	id, err := utils.ParseIntParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.svc.Delete(c, currentUserID(c), id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// GetDefault 获取默认地址
// @Summary 获取默认地址
// @Tags addresses
// @Tags frontend
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=Address}
// @Router /api/v1/addresses/default [get]
func (h *AddressHandler) GetDefault(c *gin.Context) {
	result, err := h.svc.GetDefault(c, currentUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ── Routes ────────────────────────────────────────

func RegisterAddressRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	repo := NewAddressRepository(db)
	svc := NewAddressService(repo)
	h := NewAddressHandler(svc)

	addresses := v1.Group("/addresses")
	addresses.Use(middleware.JWTAuth())
	{
		addresses.POST("", h.Create)
		addresses.GET("", h.List)
		addresses.GET("/default", h.GetDefault)
		addresses.GET("/:id", h.GetByID)
		addresses.PUT("/:id", h.Update)
		addresses.DELETE("/:id", h.Delete)
	}
}

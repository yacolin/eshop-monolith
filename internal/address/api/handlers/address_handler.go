package handlers

import (
	"strconv"

	"eshop-monolith/internal/address/api/dto"
	"eshop-monolith/internal/address/service"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/response"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	svc *service.AddressService
}

func NewAddressHandler(svc *service.AddressService) *AddressHandler {
	return &AddressHandler{svc: svc}
}

// getCurrentUserID 从 JWT 上下文中提取当前用户 ID
func (h *AddressHandler) getCurrentUserID(c *gin.Context) int64 {
	v, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	switch id := v.(type) {
	case float64:
		return int64(id)
	case uint:
		return int64(id)
	case int64:
		return id
	case string:
		parsed, _ := strconv.ParseInt(id, 10, 64)
		return parsed
	}
	return 0
}

// Create 创建新地址
// @Tags addresses
// @Summary 创建地址
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateAddressReq true "地址信息"
// @Success 200 {object} response.Response{data=dto.AddressResp}
// @Router /api/v1/addresses [post]
func (h *AddressHandler) Create(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		_ = c.Error(errcode.ErrUnauthorized)
		return
	}
	var req dto.CreateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	addr, err := h.svc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, dto.ToAddressResp(addr))
}

// List 获取当前用户地址列表
// @Tags addresses
// @Summary 地址列表
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=dto.AddressListResp}
// @Router /api/v1/addresses [get]
func (h *AddressHandler) List(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		_ = c.Error(errcode.ErrUnauthorized)
		return
	}
	resp, err := h.svc.ListByUser(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, resp)
}

// Get 根据 ID 获取地址详情
// @Tags addresses
// @Summary 获取地址详情
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "地址ID"
// @Success 200 {object} response.Response{data=dto.AddressResp}
// @Router /api/v1/addresses/{id} [get]
func (h *AddressHandler) Get(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		_ = c.Error(errcode.ErrUnauthorized)
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParams)
		return
	}
	addr, err := h.svc.GetByID(c.Request.Context(), userID, id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, dto.ToAddressResp(addr))
}

// Update 更新地址
// @Tags addresses
// @Summary 更新地址
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "地址ID"
// @Param request body dto.UpdateAddressReq true "更新信息"
// @Success 200 {object} response.Response{data=dto.AddressResp}
// @Router /api/v1/addresses/{id} [put]
func (h *AddressHandler) Update(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		_ = c.Error(errcode.ErrUnauthorized)
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParams)
		return
	}
	var req dto.UpdateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	addr, err := h.svc.Update(c.Request.Context(), userID, id, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, dto.ToAddressResp(addr))
}

// Delete 删除地址
// @Tags addresses
// @Summary 删除地址
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "地址ID"
// @Success 200 {object} response.Response
// @Router /api/v1/addresses/{id} [delete]
func (h *AddressHandler) Delete(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		_ = c.Error(errcode.ErrUnauthorized)
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParams)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), userID, id); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, nil)
}

// GetDefault 获取默认地址
// @Tags addresses
// @Summary 获取默认地址
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=dto.AddressResp}
// @Router /api/v1/addresses/default [get]
func (h *AddressHandler) GetDefault(c *gin.Context) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		_ = c.Error(errcode.ErrUnauthorized)
		return
	}
	addr, err := h.svc.GetDefault(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, dto.ToAddressResp(addr))
}

package handlers

import (
	"strconv"

	"eshop-monolith/internal/cart/api/dto"
	"eshop-monolith/internal/cart/service"
	"eshop-monolith/pkg/errcode"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CartHandler 购物车处理器
type CartHandler struct {
	cartService *service.CartService
}

// NewCartHandler 创建购物车处理器实例
func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// GetCart 获取购物车
// @Summary 获取购物车
// @Description 根据用户ID或会话ID获取购物车详情
// @Tags 购物车
// @Accept json
// @Produce json
// @Param user_id query int64 false "用户ID"
// @Param session_id query string false "会话ID"
// @Success 200 {object} response.Response{data=dto.CartResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/carts [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	sessionID := c.Query("session_id")

	// 校验用户ID和会话ID
	if userID <= 0 && sessionID == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	cart, err := h.cartService.GetCart(c, userID, sessionID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, cart)
}

// AddToCart 添加商品到购物车
// @Summary 添加商品到购物车
// @Description 将商品添加到购物车，支持指定数量和SKU
// @Tags 购物车
// @Accept json
// @Produce json
// @Param user_id query int64 false "用户ID"
// @Param session_id query string false "会话ID"
// @Param request body dto.AddToCartDTO true "添加商品请求"
// @Success 200 {object} response.Response{data=dto.CartResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/carts/items [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	sessionID := c.Query("session_id")

	// 校验用户ID和会话ID
	if userID <= 0 && sessionID == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	var req dto.AddToCartDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	cart, err := h.cartService.AddToCart(c, userID, sessionID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, cart)
}

// UpdateCartItem 更新购物车项
// @Summary 更新购物车项
// @Description 更新购物车项的数量
// @Tags 购物车
// @Accept json
// @Produce json
// @Param user_id query int64 false "用户ID"
// @Param session_id query string false "会话ID"
// @Param item_id path int64 true "购物车项ID"
// @Param request body dto.UpdateCartItemDTO true "更新购物车项请求"
// @Success 200 {object} response.Response{data=dto.CartResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/carts/items/{item_id} [put]
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	sessionID := c.Query("session_id")

	// 校验用户ID和会话ID
	if userID <= 0 && sessionID == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	itemID, err := utils.ParseIntParam(c, "item_id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateCartItemDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	cart, err := h.cartService.UpdateCartItem(c, userID, sessionID, itemID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, cart)
}

// RemoveCartItem 删除购物车项
// @Summary 删除购物车项
// @Description 从购物车中删除指定的商品
// @Tags 购物车
// @Accept json
// @Produce json
// @Param user_id query int64 false "用户ID"
// @Param session_id query string false "会话ID"
// @Param item_id path int64 true "购物车项ID"
// @Success 200 {object} response.Response{data=dto.CartResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/carts/items/{item_id} [delete]
func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	sessionID := c.Query("session_id")

	// 校验用户ID和会话ID
	if userID <= 0 && sessionID == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	itemID, err := utils.ParseIntParam(c, "item_id")
	if err != nil {
		c.Error(err)
		return
	}

	cart, err := h.cartService.RemoveCartItem(c, userID, sessionID, itemID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, cart)
}

// ClearCart 清空购物车
// @Summary 清空购物车
// @Description 清空购物车中的所有商品
// @Tags 购物车
// @Accept json
// @Produce json
// @Param user_id query int64 false "用户ID"
// @Param session_id query string false "会话ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/carts [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	sessionID := c.Query("session_id")

	// 校验用户ID和会话ID
	if userID <= 0 && sessionID == "" {
		c.Error(errcode.ErrInvalidParams)
		return
	}

	err := h.cartService.ClearCart(c, userID, sessionID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "Cart cleared successfully"})
}

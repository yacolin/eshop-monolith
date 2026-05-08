package handlers

import (
	"errors"
	"eshop-monolith/pkg/response"
	"eshop-monolith/internal/user/api/dto"
	"eshop-monolith/internal/user/domain/auth"
	"eshop-monolith/internal/user/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService  *service.AuthService
	tokenService *service.TokenService
	userService  *service.UserService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService, tokenService *service.TokenService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		tokenService: tokenService,
		userService:  userService,
	}
}

// @Summary 密码登录
// @Description 使用用户名和密码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.PasswordLoginRequest true "登录参数"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/login/password [post]
func (h *AuthHandler) LoginByPassword(c *gin.Context) {
	var req dto.PasswordLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	payload := &auth.PasswordPayload{
		Username: req.Username,
		Password: req.Password,
	}

	user, identity, err := h.authService.LoginByPassword(c, payload)
	if err != nil {
		c.Error(err)
		return
	}

	// 生成token
	tokenPair, err := h.tokenService.GenerateTokenPair(c, user.ID, identity.ID, identity.Provider, nil)
	if err != nil {
		c.Error(err)
		return
	}

	// 记录登录历史
	h.authService.RecordLoginHistory(c, user.ID, identity.ID, identity.Provider, "login", "success", "", c.ClientIP(), c.Request.UserAgent())

	resp := dto.LoginResponse{
		UserID:       user.ID,
		Username:     identity.Identifier,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt.Unix(),
		TokenType:    tokenPair.TokenType,
		IsNewUser:    false,
	}

	response.Success(c, resp)
}

// @Summary 微信登录
// @Description 使用微信code登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.WechatLoginRequest true "登录参数"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/login/wechat [post]
func (h *AuthHandler) LoginByWechat(c *gin.Context) {
	var req dto.WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	payload := &auth.WechatPayload{
		Code:   req.Code,
		AppID:  req.AppID,
		Source: req.Source,
	}

	// TODO: 从配置获取appSecret
	appSecret := ""
	user, identity, isNew, err := h.authService.LoginByWechat(c, payload, appSecret)
	if err != nil {
		c.Error(err)
		return
	}

	// 如果是新用户，需要后续绑定或注册
	if isNew {
		response.Success(c, gin.H{
			"is_new_user": true,
			"openid":      identity.Identifier,
			"provider":    identity.Provider,
		})
		return
	}

	// 生成token
	tokenPair, err := h.tokenService.GenerateTokenPair(c, user.ID, identity.ID, identity.Provider, nil)
	if err != nil {
		c.Error(err)
		return
	}

	// 记录登录历史
	h.authService.RecordLoginHistory(c, user.ID, identity.ID, identity.Provider, "login", "success", "", c.ClientIP(), c.Request.UserAgent())

	resp := dto.LoginResponse{
		UserID:       user.ID,
		Username:     identity.Identifier,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt.Unix(),
		TokenType:    tokenPair.TokenType,
		IsNewUser:    false,
	}

	response.Success(c, resp)
}

// @Summary 手机号登录
// @Description 使用手机号和验证码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.PhoneLoginRequest true "登录参数"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/login/phone [post]
func (h *AuthHandler) LoginByPhone(c *gin.Context) {
	var req dto.PhoneLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	payload := &auth.PhonePayload{
		Phone:      req.Phone,
		VerifyCode: req.VerifyCode,
	}

	user, identity, isNew, err := h.authService.LoginByPhone(c, payload)
	if err != nil {
		c.Error(err)
		return
	}

	// 如果是新用户，需要后续绑定或注册
	if isNew {
		response.Success(c, gin.H{
			"is_new_user": true,
			"phone":       identity.Identifier,
			"provider":    identity.Provider,
		})
		return
	}

	// 生成token
	tokenPair, err := h.tokenService.GenerateTokenPair(c, user.ID, identity.ID, identity.Provider, nil)
	if err != nil {
		c.Error(err)
		return
	}

	// 记录登录历史
	h.authService.RecordLoginHistory(c, user.ID, identity.ID, identity.Provider, "login", "success", "", c.ClientIP(), c.Request.UserAgent())

	resp := dto.LoginResponse{
		UserID:       user.ID,
		Username:     identity.Identifier,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt.Unix(),
		TokenType:    tokenPair.TokenType,
		IsNewUser:    false,
	}

	response.Success(c, resp)
}

// @Summary 用户注册
// @Description 用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "注册参数"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	payload := &auth.RegisterPayload{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		Provider: req.Provider,
	}

	user, identity, err := h.authService.Register(c, payload)
	if err != nil {
		c.Error(err)
		return
	}

	// 生成token
	tokenPair, err := h.tokenService.GenerateTokenPair(c, user.ID, identity.ID, identity.Provider, nil)
	if err != nil {
		c.Error(err)
		return
	}

	resp := dto.LoginResponse{
		UserID:       user.ID,
		Username:     identity.Identifier,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt.Unix(),
		TokenType:    tokenPair.TokenType,
		IsNewUser:    true,
	}

	response.Success(c, resp)
}

// @Summary 刷新Token
// @Description 使用refresh token获取新的access token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "刷新参数"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	tokenPair, err := h.tokenService.RefreshToken(c, req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	resp := dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt.Unix(),
		TokenType:    tokenPair.TokenType,
	}

	response.Success(c, resp)
}

// @Summary 登出
// @Description 用户登出，撤销token
// @Tags 认证
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从header获取token
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	if tokenString != "" {
		claims, err := h.tokenService.ParseToken(tokenString)
		if err == nil && claims != nil {
			// 撤销token
			_ = h.tokenService.RevokeToken(c, claims.JTI)
		}
	}

	response.Success(c, gin.H{"message": "登出成功"})
}

// @Summary 获取当前用户信息
// @Description 获取当前登录用户信息
// @Tags 认证
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.Error(errors.New("未登录"))
		return
	}

	user, err := h.userService.GetByID(c, userID.(int64))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, user)
}

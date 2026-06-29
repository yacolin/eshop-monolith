package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/response"
)

type AuthHandler struct {
	authSvc *AuthService
	tokenSvc *TokenService
}

func NewAuthHandler(authSvc *AuthService, tokenSvc *TokenService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, tokenSvc: tokenSvc}
}

// LoginByPassword 密码登录
// @Summary 密码登录
// @Tags auth
// @Accept json
// @Produce json
// @Param request body PasswordLoginReq true "登录信息"
// @Success 200 {object} response.Response{data=LoginResponse}
// @Router /api/v1/auth/login/password [post]
func (h *AuthHandler) LoginByPassword(c *gin.Context) {
	var req PasswordLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	user, tokenPair, err := h.authSvc.LoginByPassword(c, req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, LoginResponse{
		UserID:        user.ID,
		Username:      user.Username,
		TokenResponse: TokenResponse{AccessToken: tokenPair.AccessToken, RefreshToken: tokenPair.RefreshToken, ExpiresAt: tokenPair.ExpiresAt, TokenType: tokenPair.TokenType},
	})
}

// Register 注册
// @Summary 注册
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterReq true "注册信息"
// @Success 200 {object} response.Response{data=LoginResponse}
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	user, tokenPair, err := h.authSvc.Register(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, LoginResponse{
		UserID:        user.ID,
		Username:      user.Username,
		TokenResponse: TokenResponse{AccessToken: tokenPair.AccessToken, RefreshToken: tokenPair.RefreshToken, ExpiresAt: tokenPair.ExpiresAt, TokenType: tokenPair.TokenType},
	})
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenReq true "刷新令牌"
// @Success 200 {object} response.Response{data=TokenResponse}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	tokenPair, err := h.tokenSvc.RefreshToken(c, req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		TokenType:    tokenPair.TokenType,
	})
}

// ── Routes ────────────────────────────────────────

func RegisterAuthRoutes(v1 *gin.RouterGroup, db *gorm.DB, tokenSvc *TokenService, roleRepo IroleRepository, userRepo IuserRepository, infoRepo IuserInfoRepository, loginHistoryRepo IloginHistoryRepository) {
	authSvc := NewAuthService(db, userRepo, infoRepo, roleRepo, loginHistoryRepo, tokenSvc)
	h := NewAuthHandler(authSvc, tokenSvc)

	auth := v1.Group("/auth")
	{
		auth.POST("/login/password", h.LoginByPassword)
		auth.POST("/register", h.Register)
		auth.POST("/refresh", h.RefreshToken)
	}
}

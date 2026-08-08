package staff

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
	"eshop-monolith/pkg/token"
)

type AuthHandler struct {
	authSvc *AuthService
}

func NewAuthHandler(authSvc *AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Login 员工登录
// @Summary 员工登录
// @Tags staffs
// @Accept json
// @Produce json
// @Param request body StaffLoginReq true "登录信息"
// @Success 200 {object} response.Response{data=StaffLoginResponse}
// @Router /api/v1/admin/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req StaffLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	staff, roleNames, tokenPair, err := h.authSvc.Login(c, req.Username, req.Password, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, StaffLoginResponse{
		StaffID:  staff.ID,
		Username: staff.Username,
		RealName: staff.RealName,
		Avatar:   staff.Avatar,
		Roles:    roleNames,
		TokenResponse: TokenResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresAt:    tokenPair.ExpiresAt,
			TokenType:    tokenPair.TokenType,
		},
	})
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Tags staffs
// @Accept json
// @Produce json
// @Param request body RefreshTokenReq true "刷新令牌"
// @Success 200 {object} response.Response{data=TokenResponse}
// @Router /api/v1/admin/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	tokenPair, err := token.RefreshToken(req.RefreshToken)
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

// Profile 当前员工信息
// @Summary 当前员工信息
// @Tags staffs
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=StaffProfileResponse}
// @Router /api/v1/admin/auth/profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	staff, roleNames, err := h.authSvc.Profile(c, currentStaffID(c))
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, StaffProfileResponse{
		StaffID:   staff.ID,
		Username:  staff.Username,
		RealName:  staff.RealName,
		Email:     staff.Email,
		Phone:     staff.Phone,
		Avatar:    staff.Avatar,
		Status:    staff.Status,
		Roles:     roleNames,
		CreatedAt: staff.CreatedAt,
	})
}

// ── Routes ────────────────────────────────────────

func RegisterAuthRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	authSvc := NewAuthService(NewStaffRepository(db), NewStaffRoleRepository(db), NewStaffLoginHistoryRepository(db))
	h := NewAuthHandler(authSvc)

	admin := v1.Group("/admin/auth")
	{
		admin.POST("/login", h.Login)
		admin.POST("/refresh", h.RefreshToken)
		admin.GET("/profile", middleware.JWTAuth(), h.Profile)
	}
}

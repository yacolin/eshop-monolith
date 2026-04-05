package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"eshop-monolith/internal/user/domain/models"
	"eshop-monolith/internal/user/domain/repositories"
	"time"

	"eshop-monolith/internal/pkg/errcode"

	"github.com/golang-jwt/jwt/v5"
)

// TokenPair Token对
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// TokenClaims JWT Claims
type TokenClaims struct {
	UserID     int64    `json:"user_id"`
	IdentityID int64    `json:"identity_id"`
	Provider   string   `json:"provider"`
	Roles      []string `json:"roles"`
	JTI        string   `json:"jti"`
	jwt.RegisteredClaims
}

// TokenService Token服务
type TokenService struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	tokenRepo     repositories.IauthTokenRepository
	roleRepo      repositories.IroleRepository
}

// generateJTI 生成JWT ID
func generateJTI() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// 如果生成失败，使用时间戳和随机数
		return hex.EncodeToString([]byte(time.Now().String()))
	}
	return hex.EncodeToString(b)
}

// TokenServiceOption Token服务配置选项
type TokenServiceOption func(*TokenService)

// WithAccessExpiry 设置访问令牌过期时间
func WithAccessExpiry(expiry time.Duration) TokenServiceOption {
	return func(s *TokenService) {
		s.accessExpiry = expiry
	}
}

// WithRefreshExpiry 设置刷新令牌过期时间
func WithRefreshExpiry(expiry time.Duration) TokenServiceOption {
	return func(s *TokenService) {
		s.refreshExpiry = expiry
	}
}

// NewTokenService 创建Token服务实例
func NewTokenService(secret string, tokenRepo repositories.IauthTokenRepository, roleRepo repositories.IroleRepository, opts ...TokenServiceOption) *TokenService {
	svc := &TokenService{
		secret:        []byte(secret),
		accessExpiry:  2 * time.Hour,      // 默认2小时
		refreshExpiry: 7 * 24 * time.Hour, // 默认7天
		tokenRepo:     tokenRepo,
		roleRepo:      roleRepo,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

// GenerateTokenPair 生成Token对
func (s *TokenService) GenerateTokenPair(ctx context.Context, userID int64, identityID int64, provider string, meta map[string]interface{}) (*TokenPair, error) {
	now := time.Now()

	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	accessJTI := generateJTI()
	refreshJTI := generateJTI()

	accessExpiresAt := now.Add(s.accessExpiry)
	accessClaims := TokenClaims{
		UserID:     userID,
		IdentityID: identityID,
		Provider:   provider,
		Roles:      roleNames,
		JTI:        accessJTI,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        accessJTI,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.secret)
	if err != nil {
		return nil, errcode.ErrGenerateAccessToken
	}

	refreshExpiresAt := now.Add(s.refreshExpiry)
	refreshClaims := TokenClaims{
		UserID: userID,
		Roles:  roleNames,
		JTI:    refreshJTI,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        refreshJTI,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.secret)
	if err != nil {
		return nil, errcode.ErrGenerateRefreshToken
	}

	metaJSON, _ := json.Marshal(meta)
	dbToken := &models.AuthToken{
		UserID:    userID,
		JTI:       refreshJTI,
		TokenType: models.TokenTypeRefreshToken,
		ExpiresAt: refreshExpiresAt,
		Revoked:   false,
		Meta:      string(metaJSON),
	}

	if err := s.tokenRepo.Create(ctx, dbToken); err != nil {
		return nil, errcode.ErrSaveRefreshToken
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiresAt,
		TokenType:    "Bearer",
	}, nil
}

// GenerateAccessToken 生成访问令牌
func (s *TokenService) GenerateAccessToken(userID int64, identityID int64, provider string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.accessExpiry)
	jti := generateJTI()

	roles, err := s.roleRepo.GetUserRoles(context.Background(), userID)
	if err != nil {
		return "", time.Time{}, err
	}

	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	claims := TokenClaims{
		UserID:     userID,
		IdentityID: identityID,
		Provider:   provider,
		Roles:      roleNames,
		JTI:        jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, errcode.ErrGenerateAccessToken
	}

	return tokenString, expiresAt, nil
}

// ParseToken 解析Token
func (s *TokenService) ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errcode.ErrUnexpectedSigningMethod
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, errcode.ErrParseToken
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errcode.ErrInvalidToken
}

// ValidateToken 验证Token
func (s *TokenService) ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt.After(time.Now().Add(s.accessExpiry)) {
		revoked, err := s.tokenRepo.IsRevoked(ctx, claims.JTI)
		if err != nil {
			return nil, err
		}
		if revoked {
			return nil, errcode.ErrTokenRevoked
		}
	}

	return claims, nil
}

// RefreshToken 刷新Token
func (s *TokenService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	revoked, err := s.tokenRepo.IsRevoked(ctx, claims.JTI)
	if err != nil {
		return nil, err
	}
	if revoked {
		return nil, errcode.ErrTokenRevoked
	}

	if err := s.tokenRepo.Revoke(ctx, claims.JTI); err != nil {
		return nil, err
	}

	return s.GenerateTokenPair(ctx, claims.UserID, claims.IdentityID, claims.Provider, nil)
}

// RevokeToken 撤销Token
func (s *TokenService) RevokeToken(ctx context.Context, jti string) error {
	return s.tokenRepo.Revoke(ctx, jti)
}

// RevokeAllUserTokens 撤销用户的所有Token
func (s *TokenService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return s.tokenRepo.RevokeAllByUserID(ctx, userID)
}

// IsTokenRevoked 检查Token是否已撤销
func (s *TokenService) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	return s.tokenRepo.IsRevoked(ctx, jti)
}

// GetTokenExpiry 获取Token过期时间
func (s *TokenService) GetTokenExpiry() time.Duration {
	return s.accessExpiry
}

// GetRefreshTokenExpiry 获取Refresh Token过期时间
func (s *TokenService) GetRefreshTokenExpiry() time.Duration {
	return s.refreshExpiry
}

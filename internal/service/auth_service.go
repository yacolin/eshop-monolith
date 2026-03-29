package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"eshop-monolith/internal/domain/shared"
	"eshop-monolith/internal/domain/user"
)

// AuthService 认证服务
type AuthService struct {
	userRepo           user.Repository
	jwtSecret          string
	expireHours        int
	refreshExpireHours int
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo user.Repository, jwtSecret string, expireHours, refreshExpireHours int) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		expireHours:        expireHours,
		refreshExpireHours: refreshExpireHours,
	}
}

// Claims JWT声明
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// HashPassword 加密密码
func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 检查密码
func (s *AuthService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateToken 生成Token
func (s *AuthService) GenerateToken(user *user.User) (string, error) {
	expireTime := time.Now().Add(time.Duration(s.expireHours) * time.Hour)
	userIDStr := fmt.Sprintf("%d", user.ID)
	claims := &Claims{
		UserID:   userIDStr,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userIDStr,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// GenerateRefreshToken 生成刷新Token
func (s *AuthService) GenerateRefreshToken(user *user.User) (string, error) {
	expireTime := time.Now().Add(time.Duration(s.refreshExpireHours) * time.Hour)
	userIDStr := fmt.Sprintf("%d", user.ID)
	claims := &Claims{
		UserID:   userIDStr,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userIDStr,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ParseToken 解析Token
func (s *AuthService) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, shared.ErrUnauthorized
	}

	return claims, nil
}

// Login 登录
func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	// 查找用户
	foundUser, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, shared.ErrInvalidCredentials
	}

	// 检查密码
	if err := s.CheckPassword(foundUser.Password, password); err != nil {
		return nil, shared.ErrInvalidCredentials
	}

	// 生成访问Token
	accessToken, err := s.GenerateToken(foundUser)
	if err != nil {
		return nil, err
	}

	// 生成刷新Token
	refreshToken, err := s.GenerateRefreshToken(foundUser)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.expireHours * 3600,
	}, nil
}

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	// 解析刷新Token
	claims, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, shared.ErrUnauthorized
	}

	// 将字符串用户ID转换为int64
	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return nil, shared.ErrUnauthorized
	}

	// 查找用户
	foundUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, shared.ErrUnauthorized
	}

	// 生成新的访问Token
	accessToken, err := s.GenerateToken(foundUser)
	if err != nil {
		return nil, err
	}

	// 生成新的刷新Token
	newRefreshToken, err := s.GenerateRefreshToken(foundUser)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.expireHours * 3600,
	}, nil
}

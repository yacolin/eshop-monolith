package token

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"eshop-monolith/pkg/config"
	"eshop-monolith/pkg/errcode"
)

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
}

// TokenClaims JWT 声明(access 与 refresh 共用结构)
type TokenClaims struct {
	UserID   int64    `json:"user_id"` // B 端=staff_id,C 端=user_id
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	JTI      string   `json:"jti"`
	jwt.RegisteredClaims
}

const issuer = "eshop"

// jwtConfig 从全局配置读取 secret 与有效期;配置未加载时回退默认值
func jwtConfig() (secret []byte, accessExpiry, refreshExpiry time.Duration) {
	secret = []byte("your-secret-key-change-in-production")
	accessExpiry = 24 * time.Hour
	refreshExpiry = 7 * 24 * time.Hour
	if cfg := config.Get(); cfg != nil {
		secret = []byte(cfg.JWT.Secret)
		if cfg.JWT.ExpireHours > 0 {
			accessExpiry = time.Duration(cfg.JWT.ExpireHours) * time.Hour
		}
		if cfg.JWT.RefreshExpireHours > 0 {
			refreshExpiry = time.Duration(cfg.JWT.RefreshExpireHours) * time.Hour
		}
	}
	return secret, accessExpiry, refreshExpiry
}

// GenerateTokenPair 签发 access + refresh 令牌
func GenerateTokenPair(userID int64, username string, roles []string) (*TokenPair, error) {
	secret, accessExpiry, refreshExpiry := jwtConfig()
	now := time.Now()

	accessClaims := TokenClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		JTI:      generateJTI(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    issuer,
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secret)
	if err != nil {
		return nil, errcode.ErrGenerateAccessToken
	}

	refreshClaims := TokenClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		JTI:      generateJTI(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    issuer,
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secret)
	if err != nil {
		return nil, errcode.ErrGenerateRefreshToken
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(accessExpiry).UnixMilli(),
		TokenType:    "Bearer",
	}, nil
}

// ParseToken 解析并校验令牌,返回结构化 claims
func ParseToken(tokenString string) (*TokenClaims, error) {
	secret, _, _ := jwtConfig()
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errcode.ErrUnexpectedSigningMethod
		}
		return secret, nil
	})
	if err != nil {
		return nil, errcode.ErrInvalidToken
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errcode.ErrInvalidToken
	}
	return claims, nil
}

// RefreshToken 用 refresh 令牌换发新的令牌对(无状态,claims 原样携带)
func RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}
	return GenerateTokenPair(claims.UserID, claims.Username, claims.Roles)
}

// signClaims 用当前 secret 重新签名 claims(仅供测试构造过期令牌)
func signClaims(claims *TokenClaims) (string, error) {
	secret, _, _ := jwtConfig()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

// jwtTime 构造 NumericDate(供测试使用)
func jwtTime(offsetSeconds int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(offsetSeconds) * time.Second))
}

func generateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

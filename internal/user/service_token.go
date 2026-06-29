package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"eshop-monolith/pkg/errcode"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
}

type TokenClaims struct {
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	JTI      string   `json:"jti"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	roleRepo      IroleRepository
}

func NewTokenService(secret string, roleRepo IroleRepository) *TokenService {
	return &TokenService{
		secret:        []byte(secret),
		accessExpiry:  2 * time.Hour,
		refreshExpiry: 7 * 24 * time.Hour,
		roleRepo:      roleRepo,
	}
}

func (s *TokenService) GenerateTokenPair(ctx context.Context, userID int64, username string) (*TokenPair, error) {
	roles, err := s.roleRepo.GetUserRoleNames(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	jti := generateJTI()

	accessClaims := TokenClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		JTI:      jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "eshop",
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.secret)
	if err != nil {
		return nil, errcode.ErrGenerateAccessToken
	}

	refreshClaims := TokenClaims{
		UserID:   userID,
		Username: username,
		JTI:      generateJTI(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "eshop",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.secret)
	if err != nil {
		return nil, errcode.ErrGenerateRefreshToken
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(s.accessExpiry).UnixMilli(),
		TokenType:    "Bearer",
	}, nil
}

func (s *TokenService) ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errcode.ErrUnexpectedSigningMethod
		}
		return s.secret, nil
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

func (s *TokenService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return s.GenerateTokenPair(ctx, claims.UserID, claims.Username)
}

func generateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

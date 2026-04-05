package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"eshop-monolith/internal/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID int64) (string, error) {
	// default to 24h access token for backwards compatibility
	return GenerateAccessToken(userID)
}

func GenerateAccessToken(userID int64) (string, error) {
	return generateTokenWithType(fmt.Sprintf("%d", userID), 15*time.Minute, "access")
}

func GenerateRefreshToken(userID int64) (string, error) {
	return generateTokenWithType(fmt.Sprintf("%d", userID), 7*24*time.Hour, "refresh")
}

func generateTokenWithType(userID string, ttl time.Duration, typ string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(ttl).Unix(),
		"typ":     typ,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	cfg := config.Get()
	secret := ""
	if cfg != nil {
		secret = cfg.JWT.Secret
	}
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	cfg := config.Get()
	secret := ""
	if cfg != nil {
		secret = cfg.JWT.Secret
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果生成失败，使用时间戳
			return fmt.Sprintf("%d", time.Now().UnixNano())
		}
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

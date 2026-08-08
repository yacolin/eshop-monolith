package token

import (
	"strings"
	"testing"
)

func TestGenerateAndParseAccess(t *testing.T) {
	pair, err := GenerateTokenPair(42, "admin", []string{"admin", "operator"})
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" || pair.TokenType != "Bearer" {
		t.Fatalf("unexpected pair: %+v", pair)
	}
	if pair.ExpiresAt <= 0 {
		t.Fatalf("expires_at must be positive")
	}

	claims, err := ParseToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.UserID != 42 || claims.Username != "admin" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if len(claims.Roles) != 2 || claims.Roles[0] != "admin" || claims.Roles[1] != "operator" {
		t.Fatalf("unexpected roles: %v", claims.Roles)
	}
	if claims.JTI == "" {
		t.Fatal("jti must be set")
	}
}

func TestParseRejectsTamperedToken(t *testing.T) {
	pair, err := GenerateTokenPair(1, "admin", nil)
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}
	tampered := pair.AccessToken[:len(pair.AccessToken)-2] + "xx"
	if _, err := ParseToken(tampered); err == nil {
		t.Fatal("expected error for tampered token")
	}
}

func TestParseRejectsExpiredToken(t *testing.T) {
	pair, err := GenerateTokenPair(1, "admin", nil)
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}
	claims, err := ParseToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	// 手动改过期时间为 1 小时前再重新签名
	claims.ExpiresAt = jwtTime(-3600)
	expired, err := signClaims(claims)
	if err != nil {
		t.Fatalf("signClaims failed: %v", err)
	}
	if _, err := ParseToken(expired); err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestRefreshReturnsNewPair(t *testing.T) {
	pair, err := GenerateTokenPair(7, "colin", []string{"user"})
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}
	newPair, err := RefreshToken(pair.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}
	if newPair.AccessToken == pair.AccessToken {
		t.Fatal("refresh must issue a new access token")
	}
	claims, err := ParseToken(newPair.AccessToken)
	if err != nil {
		t.Fatalf("parse new access token failed: %v", err)
	}
	if claims.UserID != 7 || claims.Username != "colin" {
		t.Fatalf("claims lost after refresh: %+v", claims)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "user" {
		t.Fatalf("roles lost after refresh: %v", claims.Roles)
	}
}

func TestParseRejectsNonHMACAlgorithm(t *testing.T) {
	bad := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.invalid"
	_, err := ParseToken(bad)
	if err == nil {
		t.Fatal("expected error for non-HMAC token")
	}
	if !strings.Contains(err.Error(), "invalid token") {
		t.Fatalf("unexpected error: %v", err)
	}
}

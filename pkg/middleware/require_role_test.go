package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// buildRouter 组装测试路由:rolesSetter(可选)→ require → 业务handler
func buildRouter(roles []string, require gin.HandlerFunc) (*gin.Engine, *bool) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorHandler())
	if roles != nil {
		// 模拟 JWTAuth 已写入的 context 值
		r.Use(func(c *gin.Context) {
			c.Set("roles", roles)
			c.Next()
		})
	}
	r.Use(require)
	passed := false
	r.GET("/", func(c *gin.Context) { passed = true })
	return r, &passed
}

func TestRequireRolePasses(t *testing.T) {
	r, passed := buildRouter([]string{"admin", "operator"}, RequireRole("admin"))
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !*passed {
		t.Fatal("expected request to pass with admin role")
	}
}

func TestRequireRoleDeniesWrongRole(t *testing.T) {
	r, passed := buildRouter([]string{"operator"}, RequireRole("admin"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if *passed {
		t.Fatal("expected request to be blocked")
	}
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestRequireRoleDeniesNoRolesClaim(t *testing.T) {
	r, passed := buildRouter(nil, RequireRole("admin"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if *passed {
		t.Fatal("expected request to be blocked")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireMerchantAllowsMerchantAndAdmin(t *testing.T) {
	for _, role := range []string{"merchant", "admin"} {
		r, passed := buildRouter([]string{role}, RequireMerchant())
		r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
		if !*passed {
			t.Fatalf("expected %s to pass RequireMerchant", role)
		}
	}
}

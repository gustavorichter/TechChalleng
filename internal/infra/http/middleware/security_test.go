package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInjectionGuard_BlocksSQLInjection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(InjectionGuard())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?id=1'%20OR%20'1'='1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obteve %d", w.Code)
	}
}

func TestInjectionGuard_BlocksCommandInjection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(InjectionGuard())
	r.GET("/test/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test/abc%3Brm%20-rf", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obteve %d", w.Code)
	}
}

func TestInjectionGuard_AllowsValidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(InjectionGuard())
	r.GET("/test/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperado 200, obteve %d", w.Code)
	}
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fd-api/embed"
	"fd-api/handlers"

	"github.com/gin-gonic/gin"
)

func TestAPIKeyAuthFromEnv(t *testing.T) {
	t.Setenv("FD_API_KEY", "secret")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(APIKeyAuthFromEnv())
	r.POST("/v1/embeddings", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status without auth = %d, want 401", w.Code)
	}

	w = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody)
	req.Header.Set("Authorization", "Bearer secret")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status with auth = %d, want 200; body=%s", w.Code, w.Body.String())
	}
}

func TestCORSFromEnv(t *testing.T) {
	t.Setenv("FD_CORS_ORIGINS", "https://app.example")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORSFromEnv())
	r.POST("/v1/embeddings", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody)
	req.Header.Set("Origin", "https://app.example")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
}

func TestRateLimitFromEnv(t *testing.T) {
	t.Setenv("FD_RATE_LIMIT_ENABLED", "true")
	t.Setenv("FD_RATE_LIMIT_IP_RPM", "1")
	t.Setenv("FD_RATE_LIMIT_USER_RPM", "1")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(IPRateLimitFromEnv())
	r.Use(func(c *gin.Context) {
		user := "caller"
		c.Set(handlers.ContextKeyValidatedRequest, &embed.EmbeddingsRequest{User: &user})
		c.Next()
	})
	r.Use(UserRateLimitFromEnv())
	r.POST("/v1/embeddings", func(c *gin.Context) { c.Status(http.StatusOK) })

	first := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody)
	req.RemoteAddr = "127.0.0.1:1234"
	r.ServeHTTP(first, req)
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want 200; body=%s", first.Code, first.Body.String())
	}

	second := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody)
	req.RemoteAddr = "127.0.0.1:1234"
	r.ServeHTTP(second, req)
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want 429; body=%s", second.Code, second.Body.String())
	}
}

func TestRateLimitEnvHelpersFallback(t *testing.T) {
	t.Setenv("FD_RATE_LIMIT_ENABLED", "off")
	if rateLimitEnabledFromEnv() {
		t.Fatal("rateLimitEnabledFromEnv = true, want false")
	}
	t.Setenv("FD_RATE_LIMIT_ENABLED", "true")
	if !rateLimitEnabledFromEnv() {
		t.Fatal("rateLimitEnabledFromEnv = false, want true")
	}
	t.Setenv("FD_RATE_LIMIT_IP_RPM", "not-a-number")
	if got := envInt("FD_RATE_LIMIT_IP_RPM", 42); got != 42 {
		t.Fatalf("envInt invalid = %d, want fallback", got)
	}
}

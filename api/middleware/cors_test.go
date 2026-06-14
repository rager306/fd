package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSDefaultAllowsAnyOriginAndPreflight(t *testing.T) {
	r := corsTestRouter("")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/v1/embeddings", http.NoBody)
	req.Header.Set("Origin", "https://example.test")
	req.Header.Set("Access-Control-Request-Method", "POST")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", w.Code)
	}
	assertCORSHeaders(t, w, "*")
}

func TestCORSConfiguredOriginAllowlist(t *testing.T) {
	r := corsTestRouter("https://app.example, https://admin.example")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody)
	req.Header.Set("Origin", "https://admin.example")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	assertCORSHeaders(t, w, "https://admin.example")
}

func TestCORSDisallowedOriginOmitsAllowOrigin(t *testing.T) {
	r := corsTestRouter("https://app.example")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody)
	req.Header.Set("Origin", "https://evil.example")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want empty", got)
	}
}

func corsTestRouter(origins string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORS(origins))
	r.POST("/v1/embeddings", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	return r
}

func assertCORSHeaders(t *testing.T, w *httptest.ResponseRecorder, wantOrigin string) {
	t.Helper()
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != wantOrigin {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, wantOrigin)
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); got != corsMethods {
		t.Fatalf("Access-Control-Allow-Methods = %q, want %q", got, corsMethods)
	}
	if got := w.Header().Get("Access-Control-Allow-Headers"); got != corsHeaders {
		t.Fatalf("Access-Control-Allow-Headers = %q, want %q", got, corsHeaders)
	}
}

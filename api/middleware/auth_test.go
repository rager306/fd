package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fd-api/handlers"

	"github.com/gin-gonic/gin"
)

func TestAPIKeyAuthMissingKeyRejectsProtectedRequest(t *testing.T) {
	r := authTestRouter("")
	w := performAuthRequest(r, "/v1/embeddings", "")
	assertAuthCode(t, w)
}

func TestAPIKeyAuthRequiresBearerToken(t *testing.T) {
	r := authTestRouter("test")
	w := performAuthRequest(r, "/v1/embeddings", "")
	assertAuthCode(t, w)
}

func TestAPIKeyAuthRejectsWrongBearerToken(t *testing.T) {
	r := authTestRouter("test")
	w := performAuthRequest(r, "/v1/embeddings", "Bearer wrong")
	assertAuthCode(t, w)
}

func TestAPIKeyAuthAcceptsCorrectBearerToken(t *testing.T) {
	r := authTestRouter("test")
	w := performAuthRequest(r, "/v1/embeddings", "Bearer test")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
}

func TestAPIKeyAuthSkipsPublicEndpoints(t *testing.T) {
	r := authTestRouter("test")
	for _, path := range []string{publicLivePath, publicReadyPath, publicHealthPath, publicV1Healthcheck, publicDocs, publicDocs + "/index.html", publicOpenAPI} {
		w := performAuthRequest(r, path, "")
		if w.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200; body=%s", path, w.Code, w.Body.String())
		}
	}
}

func TestAPIKeyAuthProtectsMetrics(t *testing.T) {
	r := authTestRouter("test")
	w := performAuthRequest(r, publicMetrics, "")
	assertAuthCode(t, w)
}

func authTestRouter(apiKey string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(APIKeyAuth(apiKey))
	r.GET("/*path", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	r.POST("/*path", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	return r
}

func performAuthRequest(r http.Handler, path, authorization string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, http.NoBody)
	if path == publicLivePath || path == publicReadyPath || path == publicHealthPath || path == publicV1Healthcheck || path == publicMetrics || path == publicDocs || path == publicDocs+"/index.html" || path == publicOpenAPI {
		req.Method = http.MethodGet
	}
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	r.ServeHTTP(w, req)
	return w
}

func assertAuthCode(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("HTTP = %d, want %d; body=%s", w.Code, http.StatusUnauthorized, w.Body.String())
	}
	var resp handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal envelope: %v; body=%s", err, w.Body.String())
	}
	if resp.Error.Code != handlers.CodeUnauthorized {
		t.Fatalf("error.code = %q, want %q", resp.Error.Code, handlers.CodeUnauthorized)
	}
	if resp.Error.Type != handlers.TypeAuthError {
		t.Fatalf("error.type = %q, want %q", resp.Error.Type, handlers.TypeAuthError)
	}
	if resp.Error.Param != "authorization" {
		t.Fatalf("error.param = %q, want %q", resp.Error.Param, "authorization")
	}
}

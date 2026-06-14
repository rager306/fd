package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fd-api/embed"
	"fd-api/handlers"

	"github.com/gin-gonic/gin"
)

func TestIPRateLimitRejectsRequestAfterLimit(t *testing.T) {
	r := rateLimitTestRouter(IPRateLimit(true, 100), nil)
	var last *httptest.ResponseRecorder
	for i := 0; i < 101; i++ {
		last = performRateLimitRequest(r, "", "127.0.0.1:1234")
	}
	if last.Code != http.StatusTooManyRequests {
		t.Fatalf("101st status = %d, want 429; body=%s", last.Code, last.Body.String())
	}
	assertRateLimitHeaders(t, last, "100", "0", "60")
	assertRateLimitEnvelope(t, last)
}

func TestIPRateLimitSeparatesClientIPs(t *testing.T) {
	r := rateLimitTestRouter(IPRateLimit(true, 1), nil)
	first := performRateLimitRequest(r, "", "127.0.0.1:1234")
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want 200", first.Code)
	}
	secondIP := performRateLimitRequest(r, "", "127.0.0.2:1234")
	if secondIP.Code != http.StatusOK {
		t.Fatalf("second IP status = %d, want 200", secondIP.Code)
	}
	sameIP := performRateLimitRequest(r, "", "127.0.0.1:1234")
	if sameIP.Code != http.StatusTooManyRequests {
		t.Fatalf("same IP status = %d, want 429", sameIP.Code)
	}
}

func TestUserRateLimitSeparatesUsers(t *testing.T) {
	r := rateLimitTestRouter(nil, UserRateLimit(true, 1))
	first := performRateLimitRequest(r, "alice", "127.0.0.1:1234")
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want 200", first.Code)
	}
	bob := performRateLimitRequest(r, "bob", "127.0.0.1:1234")
	if bob.Code != http.StatusOK {
		t.Fatalf("bob status = %d, want 200", bob.Code)
	}
	aliceAgain := performRateLimitRequest(r, "alice", "127.0.0.1:1234")
	if aliceAgain.Code != http.StatusTooManyRequests {
		t.Fatalf("alice again status = %d, want 429", aliceAgain.Code)
	}
}

func TestRateLimitDisabledDoesNotSetHeadersOrReject(t *testing.T) {
	r := rateLimitTestRouter(IPRateLimit(false, 1), nil)
	first := performRateLimitRequest(r, "", "127.0.0.1:1234")
	second := performRateLimitRequest(r, "", "127.0.0.1:1234")
	if first.Code != http.StatusOK || second.Code != http.StatusOK {
		t.Fatalf("disabled limiter statuses = %d/%d, want 200/200", first.Code, second.Code)
	}
	if got := second.Header().Get(headerRateLimitLimit); got != "" {
		t.Fatalf("rate limit header = %q, want empty when disabled", got)
	}
}

func rateLimitTestRouter(ipMiddleware, userMiddleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if ipMiddleware != nil {
		r.Use(ipMiddleware)
	}
	r.POST("/v1/embeddings", func(c *gin.Context) {
		user := c.GetHeader("X-Test-User")
		if user != "" {
			c.Set(handlers.ContextKeyValidatedRequest, &embed.EmbeddingsRequest{User: &user})
		}
		if userMiddleware != nil {
			userMiddleware(c)
			if c.IsAborted() {
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func performRateLimitRequest(r http.Handler, user, remoteAddr string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody)
	req.RemoteAddr = remoteAddr
	if user != "" {
		req.Header.Set("X-Test-User", user)
	}
	r.ServeHTTP(w, req)
	return w
}

func assertRateLimitHeaders(t *testing.T, w *httptest.ResponseRecorder, wantLimit, wantRemaining, wantReset string) {
	t.Helper()
	if got := w.Header().Get(headerRateLimitLimit); got != wantLimit {
		t.Fatalf("%s = %q, want %q", headerRateLimitLimit, got, wantLimit)
	}
	if got := w.Header().Get(headerRateLimitRemaining); got != wantRemaining {
		t.Fatalf("%s = %q, want %q", headerRateLimitRemaining, got, wantRemaining)
	}
	if got := w.Header().Get(headerRateLimitReset); got != wantReset {
		t.Fatalf("%s = %q, want %q", headerRateLimitReset, got, wantReset)
	}
	if got := w.Header().Get("Retry-After"); got != "60" {
		t.Fatalf("Retry-After = %q, want 60", got)
	}
}

func assertRateLimitEnvelope(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	var resp handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal envelope: %v; body=%s", err, w.Body.String())
	}
	if resp.Error.Code != handlers.CodeRateLimitExceeded {
		t.Fatalf("error.code = %q, want %q", resp.Error.Code, handlers.CodeRateLimitExceeded)
	}
	if resp.Error.Type != handlers.TypeRateLimitError {
		t.Fatalf("error.type = %q, want %q", resp.Error.Type, handlers.TypeRateLimitError)
	}
}

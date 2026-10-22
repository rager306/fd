package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCacheHeadersSetsETagAndCacheControl(t *testing.T) {
	r := cacheHeadersTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Header().Get(headerETag); got == "" {
		t.Fatalf("ETag header is empty")
	}
	if got := w.Header().Get(headerCacheControl); got != cacheControlValue {
		t.Fatalf("Cache-Control = %q, want %q", got, cacheControlValue)
	}
	if body := w.Body.String(); body != `{"ok":true}` {
		t.Fatalf("body = %q, want JSON payload", body)
	}
}

func TestCacheHeadersReturnsNotModifiedForMatchingETag(t *testing.T) {
	r := cacheHeadersTestRouter()
	first := httptest.NewRecorder()
	firstReq := httptest.NewRequest(http.MethodGet, "/info", http.NoBody)
	r.ServeHTTP(first, firstReq)
	etag := first.Header().Get(headerETag)
	if etag == "" {
		t.Fatal("first response ETag is empty")
	}

	second := httptest.NewRecorder()
	secondReq := httptest.NewRequest(http.MethodGet, "/info", http.NoBody)
	secondReq.Header.Set(headerIfNoneMatch, etag)
	r.ServeHTTP(second, secondReq)

	if second.Code != http.StatusNotModified {
		t.Fatalf("status = %d, want 304; body=%s", second.Code, second.Body.String())
	}
	if body := second.Body.String(); body != "" {
		t.Fatalf("304 body = %q, want empty", body)
	}
	if got := second.Header().Get(headerETag); got != etag {
		t.Fatalf("ETag = %q, want %q", got, etag)
	}
}

func TestCacheHeadersSkipsNonCacheablePath(t *testing.T) {
	r := cacheHeadersTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/live", http.NoBody)
	r.ServeHTTP(w, req)

	if got := w.Header().Get(headerETag); got != "" {
		t.Fatalf("ETag = %q, want empty", got)
	}
}

func TestETagMatches(t *testing.T) {
	const (
		etagABC     = `"abc"`
		etagDEF     = `"def"`
		etagWeakABC = `W/"abc"`
	)
	cases := []struct {
		name        string
		ifNoneMatch string
		etag        string
		want        bool
	}{
		{name: "empty header and empty etag", ifNoneMatch: "", etag: `""`, want: true},
		{name: "empty header", ifNoneMatch: "", etag: etagABC, want: false},
		{name: "exact quoted match", ifNoneMatch: etagABC, etag: etagABC, want: true},
		{name: "exact quoted mismatch", ifNoneMatch: etagABC, etag: etagDEF, want: false},
		{name: "list matches second", ifNoneMatch: `"abc", "def"`, etag: etagDEF, want: true},
		{name: "list with whitespace", ifNoneMatch: ` "abc" , "def" `, etag: etagDEF, want: true},
		{name: "unquoted candidate", ifNoneMatch: `abc`, etag: etagABC, want: true},
		{name: "trailing comma", ifNoneMatch: `"abc",`, etag: etagABC, want: true},
		{name: "comma-only header and empty etag", ifNoneMatch: ",,", etag: `""`, want: true},
		{name: "comma-only header", ifNoneMatch: ",,", etag: etagABC, want: false},
		{name: "weak etag exact", ifNoneMatch: etagWeakABC, etag: etagWeakABC, want: true},
		{name: "weak etag is not the strong etag", ifNoneMatch: etagWeakABC, etag: etagABC, want: false},
		{name: "unquoted etag does not match quoted candidate", ifNoneMatch: etagABC, etag: `abc`, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := etagMatches(tc.ifNoneMatch, tc.etag); got != tc.want {
				t.Fatalf("etagMatches(%q, %q) = %v, want %v", tc.ifNoneMatch, tc.etag, got, tc.want)
			}
		})
	}
}

func BenchmarkETagMatchesNoMatch(b *testing.B) {
	const (
		header = `"aaa", "bbb", "ccc", "ddd"`
		etag   = `"zzz"`
	)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if etagMatches(header, etag) {
			b.Fatal("unexpected match")
		}
	}
}

func BenchmarkETagMatchesFirstCandidate(b *testing.B) {
	const (
		header = `"zzz", "bbb", "ccc", "ddd"`
		etag   = `"zzz"`
	)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !etagMatches(header, etag) {
			b.Fatal("expected match")
		}
	}
}

func cacheHeadersTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CacheHeaders())
	r.GET("/info", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	r.POST("/v1/embeddings", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	r.GET("/live", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	return r
}

// The cacheable paths sit behind API key auth, so a shared cache must never
// store them; only the caller's own cache may.
func TestCacheHeadersMarksAuthenticatedResponsesPrivate(t *testing.T) {
	r := cacheHeadersTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody)
	r.ServeHTTP(w, req)

	got := w.Header().Get(headerCacheControl)
	if !strings.HasPrefix(got, "private") {
		t.Fatalf("Cache-Control = %q, want a private directive", got)
	}
	if strings.Contains(got, "public") {
		t.Fatalf("Cache-Control = %q, must not allow shared caching", got)
	}
}

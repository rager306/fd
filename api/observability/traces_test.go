package observability

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTraceStoreRecordsRecentRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := NewTraceStore(defaultTraceCapacity, true)
	r := gin.New()
	r.Use(store.Middleware("test-model"))
	r.GET("/request/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Header("X-Request-Id", "req-"+id)
		c.Header("X-Model-Id", "test-model")
		c.Header("X-Dimensions", "1024")
		c.JSON(http.StatusOK, gin.H{"id": id})
	})
	r.GET("/v1/traces", store.Handler())

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/request/"+strconv.Itoa(i), http.NoBody)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want 200", i, w.Code)
		}
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/traces", http.NoBody)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("traces status = %d, want 200", w.Code)
	}
	var entries []TraceEntry
	if err := json.Unmarshal(w.Body.Bytes(), &entries); err != nil {
		t.Fatalf("unmarshal traces: %v", err)
	}
	if len(entries) != 5 {
		t.Fatalf("entries = %d, want 5", len(entries))
	}
	for i, entry := range entries {
		if entry.Timestamp.IsZero() {
			t.Fatalf("entry %d timestamp is zero", i)
		}
		if entry.Status != http.StatusOK {
			t.Fatalf("entry %d status = %d, want 200", i, entry.Status)
		}
		if entry.ModelID != "test-model" {
			t.Fatalf("entry %d model_id = %q, want test-model", i, entry.ModelID)
		}
		if entry.RequestID != "req-"+strconv.Itoa(i) {
			t.Fatalf("entry %d request_id = %q", i, entry.RequestID)
		}
		if entry.Path != "/request/"+strconv.Itoa(i) {
			t.Fatalf("entry %d path = %q", i, entry.Path)
		}
		if entry.Dimensions != 1024 {
			t.Fatalf("entry %d dimensions = %d, want 1024", i, entry.Dimensions)
		}
		if entry.LatencyMS < 0 {
			t.Fatalf("entry %d latency_ms = %d, want >= 0", i, entry.LatencyMS)
		}
	}
}

func TestTraceStoreRingKeepsLastEntries(t *testing.T) {
	store := NewTraceStore(3, true)
	for i := 0; i < 5; i++ {
		store.Add(TraceEntry{Path: "/" + strconv.Itoa(i)})
	}
	entries := store.Snapshot()
	if len(entries) != 3 {
		t.Fatalf("entries = %d, want 3", len(entries))
	}
	for i, want := range []string{"/2", "/3", "/4"} {
		if entries[i].Path != want {
			t.Fatalf("entry %d path = %q, want %q", i, entries[i].Path, want)
		}
	}
}

func TestTraceStoreDisabledRecordsNothing(t *testing.T) {
	store := NewTraceStore(3, false)
	store.Add(TraceEntry{Path: "/x"})
	if entries := store.Snapshot(); len(entries) != 0 {
		t.Fatalf("disabled store entries = %d, want 0", len(entries))
	}
}

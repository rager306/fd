package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fd-api/queue"

	"github.com/gin-gonic/gin"
)

func TestQueueSubmitPublishesPendingAndCleansUpOnFull(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := queue.NewResultStore()
	t.Cleanup(func() { _ = store.Close() })
	items := make(chan queue.Item, 1)
	h := NewQueueHandler(store, items, "model", slog.New(slog.NewTextHandler(io.Discard, nil)))
	r := gin.New()
	r.POST("/v1/queue", h.Submit)

	submit := func() *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/queue", strings.NewReader(`{"input":["hello"]}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		return w
	}

	first := submit()
	if first.Code != http.StatusAccepted {
		t.Fatalf("first status = %d; body=%s", first.Code, first.Body.String())
	}
	id := first.Header().Get("X-Request-Id")
	if id == "" {
		t.Fatal("accepted request missing ID")
	}
	if result, ok := store.Get(id); !ok || result.ID != id || result.Status != queue.StatusPending {
		t.Fatalf("pending result = %#v, present %v", result, ok)
	}

	second := submit()
	if second.Code != http.StatusServiceUnavailable {
		t.Fatalf("full queue status = %d; body=%s", second.Code, second.Body.String())
	}
	if got := store.Size(); got != 1 {
		t.Fatalf("store size after rejection = %d, want 1", got)
	}
	item := <-items
	if item.ID != id {
		t.Fatalf("item ID = %q, want %q", item.ID, id)
	}
}

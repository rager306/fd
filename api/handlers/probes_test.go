package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
)

func TestLiveHandlerAlwaysReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/live", NewLiveHandler())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/live", http.NoBody))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal live response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("live status = %#v, want ok", body["status"])
	}
	if _, ok := body["time"].(string); !ok {
		t.Fatalf("live time missing or not string: %#v", body["time"])
	}
}

func TestReadyHandlerReturns503BeforeWarmup(t *testing.T) {
	state := lifecycle.NewState()
	w := serveReady(state)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusServiceUnavailable, w.Body.String())
	}
	if got := w.Header().Get("Retry-After"); got != retryAfterWarmupSeconds {
		t.Fatalf("Retry-After = %q, want %q", got, retryAfterWarmupSeconds)
	}

	var body ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal ready error response: %v", err)
	}
	if body.Error.Code != CodeModelNotLoaded {
		t.Fatalf("error.code = %q, want %q", body.Error.Code, CodeModelNotLoaded)
	}
	if body.Error.Type != TypeOverloadedError {
		t.Fatalf("error.type = %q, want %q", body.Error.Type, TypeOverloadedError)
	}
}

func TestReadyHandlerReturnsOKAfterWarmup(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	w := serveReady(state)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal ready response: %v", err)
	}
	if body["status"] != "ready" {
		t.Fatalf("ready status = %#v, want ready", body["status"])
	}
}

func TestReadyHandlerReturns503DuringShutdown(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	state.BeginShutdown()
	w := serveReady(state)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusServiceUnavailable, w.Body.String())
	}
	if got := w.Header().Get("Retry-After"); got != retryAfterWarmupSeconds {
		t.Fatalf("Retry-After = %q, want %q", got, retryAfterWarmupSeconds)
	}
}

func serveReady(state *lifecycle.State) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ready", NewReadyHandler(state))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ready", http.NoBody))
	return w
}

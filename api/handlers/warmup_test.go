package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fd-api/embed"
	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
)

type warmupTestModel struct {
	embedFunc func(context.Context, []string) ([][]float32, error)
}

func (m *warmupTestModel) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return m.embedFunc(ctx, texts)
}

func TestWarmupStatusReadyAfterWarmup(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	w := serveWarmup(t, state, successfulWarmupModel(), http.MethodGet)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	body := decodeWarmupResponse(t, w)
	if body.Status != warmupStatusReady {
		t.Fatalf("warmup status = %q, want %q", body.Status, warmupStatusReady)
	}
	if body.Progress != 1.0 {
		t.Fatalf("progress = %f, want 1", body.Progress)
	}
}

func TestWarmupStatusWarmingUpProgress(t *testing.T) {
	state := lifecycle.NewState()
	started := make(chan struct{})
	release := make(chan struct{})
	model := &warmupTestModel{embedFunc: func(_ context.Context, _ []string) ([][]float32, error) {
		close(started)
		<-release
		return [][]float32{{1}}, nil
	}}
	handler := NewWarmupHandler(state, model, time.Second)
	r := warmupRouter(handler)

	postWarmup(t, r)
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("warmup did not start")
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/warmup", http.NoBody))
	body := decodeWarmupResponse(t, w)
	if body.Status != warmupStatusWarmingUp {
		t.Fatalf("warmup status = %q, want %q", body.Status, warmupStatusWarmingUp)
	}
	if body.Progress <= 0 || body.Progress >= 1 {
		t.Fatalf("progress = %f, want fraction between 0 and 1", body.Progress)
	}
	close(release)
	waitForWarmupReady(t, state)
}

func TestWarmupTriggerReturnsOKWhenReady(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	w := serveWarmup(t, state, successfulWarmupModel(), http.MethodPost)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	body := decodeWarmupResponse(t, w)
	if body.Status != warmupStatusReady || body.Message != "already warm" {
		t.Fatalf("response = %#v, want ready/already warm", body)
	}
}

func TestWarmupTriggerStartsBackgroundWarmup(t *testing.T) {
	state := lifecycle.NewState()
	release := make(chan struct{})
	model := &warmupTestModel{embedFunc: func(_ context.Context, _ []string) ([][]float32, error) {
		<-release
		return [][]float32{{1}}, nil
	}}
	w := serveWarmup(t, state, model, http.MethodPost)

	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusAccepted)
	}
	body := decodeWarmupResponse(t, w)
	if body.Status != warmupStatusWarmingUp || body.Message != "warmup started" {
		t.Fatalf("response = %#v, want warming_up/warmup started", body)
	}
	close(release)
	waitForWarmupReady(t, state)
}

func TestWarmupTriggerStoresError(t *testing.T) {
	state := lifecycle.NewState()
	boom := errors.New("boom")
	w := serveWarmup(t, state, &warmupTestModel{embedFunc: func(_ context.Context, _ []string) ([][]float32, error) {
		return nil, boom
	}}, http.MethodPost)
	if w.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusAccepted)
	}
	waitForCondition(t, time.Second, func() bool {
		return errors.Is(state.LastError(), boom)
	})
}

func serveWarmup(t *testing.T, state *lifecycle.State, model embed.Embedder, method string) *httptest.ResponseRecorder {
	t.Helper()
	r := warmupRouter(NewWarmupHandler(state, model, time.Second))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(method, "/warmup", http.NoBody))
	return w
}

func warmupRouter(handler *WarmupHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/warmup", handler.Status)
	r.POST("/warmup", handler.Trigger)
	return r
}

func postWarmup(t *testing.T, r http.Handler) WarmupResponse {
	t.Helper()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/warmup", http.NoBody))
	if w.Code != http.StatusAccepted {
		t.Fatalf("POST /warmup status = %d, want %d", w.Code, http.StatusAccepted)
	}
	return decodeWarmupResponse(t, w)
}

func decodeWarmupResponse(t *testing.T, w *httptest.ResponseRecorder) WarmupResponse {
	t.Helper()
	var body WarmupResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal warmup response: %v", err)
	}
	return body
}

func successfulWarmupModel() *warmupTestModel {
	return &warmupTestModel{embedFunc: func(_ context.Context, _ []string) ([][]float32, error) {
		return [][]float32{{1}}, nil
	}}
}

func waitForWarmupReady(t *testing.T, state *lifecycle.State) {
	t.Helper()
	waitForCondition(t, time.Second, state.IsWarmupDone)
}

func waitForCondition(t *testing.T, timeout time.Duration, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("condition not met within %s", timeout)
}

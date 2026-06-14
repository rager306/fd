package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fd-api/buildinfo"
	"fd-api/handlers"
	"fd-api/lifecycle"
	"fd-api/middleware"
	"fd-api/observability"

	"github.com/gin-gonic/gin"
)

func TestFdV2ObservabilityEndpointsAndHeaders(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	r := newObservabilityTestRouter(state)

	checks := []struct {
		name        string
		method      string
		path        string
		body        string
		wantStatus  int
		contentType string
		wantText    string
	}{
		{name: "T-H-7 health deep", method: http.MethodGet, path: "/health", wantStatus: http.StatusOK, wantText: "model_loaded"},
		{name: "T-H-8 live", method: http.MethodGet, path: "/live", wantStatus: http.StatusOK, wantText: "ok"},
		{name: "T-H-9 ready", method: http.MethodGet, path: "/ready", wantStatus: http.StatusOK, wantText: "ready"},
		{name: "T-H-10 version", method: http.MethodGet, path: "/version", wantStatus: http.StatusOK, wantText: "version"},
		{name: "T-E-1 version exists", method: http.MethodGet, path: "/version", wantStatus: http.StatusOK, wantText: "2.0.0"},
		{name: "T-E-2 info exists", method: http.MethodGet, path: "/info", wantStatus: http.StatusOK, wantText: "models"},
		{name: "T-E-3 metrics exists", method: http.MethodGet, path: "/metrics", wantStatus: http.StatusOK, contentType: "text/plain", wantText: "fd_requests_total"},
		{name: "T-E-4 healthcheck alias", method: http.MethodGet, path: "/v1/healthcheck", wantStatus: http.StatusOK, wantText: "model_loaded"},
		{name: "warmup exists", method: http.MethodGet, path: "/warmup", wantStatus: http.StatusOK, wantText: "ready"},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			w := serveObservabilityRequest(r, check.method, check.path, check.body, "")
			if w.Code != check.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", w.Code, check.wantStatus, w.Body.String())
			}
			if check.contentType != "" && !strings.Contains(w.Header().Get("Content-Type"), check.contentType) {
				t.Fatalf("Content-Type = %q, want contains %q", w.Header().Get("Content-Type"), check.contentType)
			}
			if check.wantText != "" && !strings.Contains(w.Body.String(), check.wantText) {
				t.Fatalf("body missing %q: %s", check.wantText, w.Body.String())
			}
		})
	}
}

func TestFdV2ObservabilityHeaders(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	r := newObservabilityTestRouter(state)

	t.Run("T-HDR-1 server", func(t *testing.T) {
		w := serveObservabilityRequest(r, http.MethodGet, "/version", "", "")
		if got := w.Header().Get("Server"); got != "fd/2.0.0" {
			t.Fatalf("Server = %q, want fd/2.0.0", got)
		}
	})
	t.Run("T-HDR-2 request id echo", func(t *testing.T) {
		w := serveObservabilityRequest(r, http.MethodGet, "/version", "", "caller-id")
		if got := w.Header().Get(middleware.HeaderRequestID); got != "caller-id" {
			t.Fatalf("X-Request-Id = %q, want caller-id", got)
		}
	})
	t.Run("T-HDR-3 request id generated", func(t *testing.T) {
		w := serveObservabilityRequest(r, http.MethodGet, "/version", "", "")
		if got := w.Header().Get(middleware.HeaderRequestID); got == "" {
			t.Fatal("X-Request-Id should be generated")
		}
	})
	t.Run("T-HDR-4 model id", func(t *testing.T) {
		w := serveObservabilityRequest(r, http.MethodPost, "/v1/embeddings", `{"model":"test","input":"hello"}`, "")
		if got := w.Header().Get(middleware.HeaderModelID); got != lifecycleTestModel {
			t.Fatalf("X-Model-Id = %q, want %q", got, lifecycleTestModel)
		}
	})
	t.Run("T-HDR-5 dimensions", func(t *testing.T) {
		w := serveObservabilityRequest(r, http.MethodPost, "/v1/embeddings", `{"model":"test","input":"hello","dimensions":512}`, "")
		if got := w.Header().Get(middleware.HeaderDimensions); got != "512" {
			t.Fatalf("X-Dimensions = %q, want 512", got)
		}
	})
	t.Run("T-HDR-8 retry after", func(t *testing.T) {
		unreadyRouter := newObservabilityTestRouter(lifecycle.NewState())
		w := serveObservabilityRequest(unreadyRouter, http.MethodGet, "/ready", "", "")
		if got := w.Header().Get("Retry-After"); got != "5" {
			t.Fatalf("Retry-After = %q, want 5", got)
		}
	})
	t.Run("T-HDR-9 connection keep alive", func(t *testing.T) {
		w := serveObservabilityRequest(r, http.MethodGet, "/version", "", "")
		if got := w.Header().Get("Connection"); got != "keep-alive" {
			t.Fatalf("Connection = %q, want keep-alive", got)
		}
	})
}

func TestFdV2ObservabilityHealthPayload(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	r := newObservabilityTestRouter(state)
	post := serveObservabilityRequest(r, http.MethodPost, "/v1/embeddings", `{"model":"test","input":"hello"}`, "")
	if post.Code != http.StatusOK {
		t.Fatalf("embedding status = %d, want %d; body=%s", post.Code, http.StatusOK, post.Body.String())
	}

	health := serveObservabilityRequest(r, http.MethodGet, "/health", "", "")
	var body handlers.DeepHealthResponse
	if err := json.Unmarshal(health.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal health: %v", err)
	}
	if body.Status != "ok" || !body.ModelLoaded || !body.WarmupDone {
		t.Fatalf("unexpected health body: %#v", body)
	}
	if body.LastInferenceAt == nil || *body.LastInferenceAt == "" {
		t.Fatalf("last_inference_at missing: %#v", body.LastInferenceAt)
	}
}

func newObservabilityTestRouter(state *lifecycle.State) *gin.Engine {
	gin.SetMode(gin.TestMode)
	info := buildinfo.New(buildinfo.Info{
		Version:   "2.0.0",
		Model:     lifecycleTestModel,
		BuildHash: "testhash",
		BuildDate: "2026-06-13T00:00:00Z",
		StartedAt: time.Now().Add(-time.Second),
	})
	runtime := &handlers.RuntimeHealth{Backend: "test", Model: lifecycleTestModel, Dimensions: 1024}
	metrics := observability.NewMetrics()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	embedHandler := handlers.NewEmbeddingsHandler(staticLifecycleEmbedder(), newLifecycleTestCache(), lifecycleTestModel, logger)
	healthHandler := handlers.NewHealthHandlerWithState(runtime, state)
	warmupHandler := handlers.NewWarmupHandler(state, staticLifecycleEmbedder(), time.Second)

	r := gin.New()
	r.Use(handlers.RecoveryMiddleware(logger))
	r.Use(middleware.HeadersMiddleware(info, lifecycleTestModel))
	r.Use(metrics.Middleware())
	r.GET("/live", handlers.NewLiveHandler())
	r.GET("/ready", handlers.NewReadyHandler(state))
	r.GET("/version", handlers.NewVersionHandler(info))
	r.GET("/info", handlers.NewInfoHandler(info, runtime, state))
	r.GET("/metrics", metrics.Handler())
	r.GET("/warmup", warmupHandler.Status)
	r.POST("/warmup", warmupHandler.Trigger)
	r.GET("/health", healthHandler)
	r.GET("/v1/healthcheck", healthHandler)
	r.POST("/v1/embeddings", middleware.ValidateEmbeddingsRequest(), middleware.LifecycleGateWithCapacity(state, 0), embedHandler.CreateEmbedding)
	return r
}

func serveObservabilityRequest(r http.Handler, method, path, body, requestID string) *httptest.ResponseRecorder {
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, http.NoBody)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	if requestID != "" {
		req.Header.Set(middleware.HeaderRequestID, requestID)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

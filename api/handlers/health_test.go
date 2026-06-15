package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
)

const (
	testModelID    = "deepvk/USER-bge-m3"
	testTEIBackend = "tei"
	testVersion200 = "2.0.0"
)

func TestHealthHandlerDefaultShape(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", HealthHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal health response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status body = %#v", body["status"])
	}
	if _, ok := body["time"].(string); !ok {
		t.Fatalf("time missing or not string: %#v", body["time"])
	}
	if _, ok := body["runtime"]; ok {
		t.Fatalf("default health should not include runtime metadata: %#v", body["runtime"])
	}
}

func TestNewHealthHandlerIncludesSafeTEIRuntimeMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", NewHealthHandler(&RuntimeHealth{
		Backend:           testTEIBackend,
		Model:             testModelID,
		Dimensions:        1024,
		ProductionDefault: true,
		CacheNamespace:    "m026-tei",
	}))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal health response: %v", err)
	}
	runtime, ok := body["runtime"].(map[string]any)
	if !ok {
		t.Fatalf("runtime metadata missing: %#v", body)
	}
	if runtime["backend"] != testTEIBackend {
		t.Fatalf("backend = %#v, want tei", runtime["backend"])
	}
	if runtime["model"] != testModelID {
		t.Fatalf("model = %#v", runtime["model"])
	}
	if dims, ok := runtime["dimensions"].(float64); !ok || dims != 1024 {
		t.Fatalf("dimensions = %#v, want 1024", runtime["dimensions"])
	}
	if runtime["production_default"] != true {
		t.Fatalf("production_default = %#v, want true", runtime["production_default"])
	}
	if runtime["cache_namespace"] != "m026-tei" {
		t.Fatalf("cache_namespace = %#v", runtime["cache_namespace"])
	}
	// ONNX-only fields must not appear for TEI health
	for _, field := range []string{"artifact_id", "provider", "tokenizer_verified", "runtime_library_verified", "tei_url"} {
		if _, ok := runtime[field]; ok {
			t.Fatalf("TEI health must not expose %q (ONNX-only field)", field)
		}
	}
	// path/ssecret fields must not appear for any backend
	for _, field := range []string{"manifest_path", "runtime_library_path", "tokenizer_path", "onnx_runtime_sha256"} {
		if _, ok := runtime[field]; ok {
			t.Fatalf("runtime health must not expose %q", field)
		}
	}
}

func TestDeepHealthReportsOKWhenReady(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	w := serveDeepHealth(state)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	body := decodeDeepHealth(t, w)
	if body.Status != healthStatusOK {
		t.Fatalf("health status = %q, want %q", body.Status, healthStatusOK)
	}
	if !body.ModelLoaded || !body.WarmupDone {
		t.Fatalf("model_loaded=%v warmup_done=%v, want true/true", body.ModelLoaded, body.WarmupDone)
	}
}

func TestDeepHealthReportsDegradedBeforeWarmup(t *testing.T) {
	state := lifecycle.NewState()
	w := serveDeepHealth(state)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusServiceUnavailable, w.Body.String())
	}
	body := decodeDeepHealth(t, w)
	if body.Status != healthStatusDegraded {
		t.Fatalf("health status = %q, want %q", body.Status, healthStatusDegraded)
	}
	if body.ModelLoaded || body.WarmupDone {
		t.Fatalf("model_loaded=%v warmup_done=%v, want false/false", body.ModelLoaded, body.WarmupDone)
	}
}

func TestDeepHealthReportsDownDuringShutdown(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	state.BeginShutdown()
	w := serveDeepHealth(state)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusServiceUnavailable, w.Body.String())
	}
	body := decodeDeepHealth(t, w)
	if body.Status != healthStatusDown {
		t.Fatalf("health status = %q, want %q", body.Status, healthStatusDown)
	}
	if body.ModelLoaded {
		t.Fatal("model_loaded should be false during shutdown")
	}
	if !body.WarmupDone {
		t.Fatal("warmup_done should remain true during shutdown")
	}
}

func TestDeepHealthIncludesLastInferenceAt(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	state.MarkInferenceSuccess()
	w := serveDeepHealth(state)

	body := decodeDeepHealth(t, w)
	if body.LastInferenceAt == nil || *body.LastInferenceAt == "" {
		t.Fatalf("last_inference_at missing: %#v", body.LastInferenceAt)
	}
}

func serveDeepHealth(state *lifecycle.State) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", NewHealthHandlerWithState(&RuntimeHealth{Backend: testTEIBackend, Model: testModelID}, state))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", http.NoBody))
	return w
}

func decodeDeepHealth(t *testing.T, w *httptest.ResponseRecorder) DeepHealthResponse {
	t.Helper()
	var body DeepHealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal deep health response: %v", err)
	}
	return body
}

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fd-api/buildinfo"
	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
)

func TestVersionHandlerReturnsBuildInfoAndUptime(t *testing.T) {
	gin.SetMode(gin.TestMode)
	startedAt := time.Now().Add(-2 * time.Second)
	r := gin.New()
	r.GET("/version", NewVersionHandler(buildinfo.Info{
		Service:      "fd-api",
		Version:      testVersion200,
		Model:        testModelID,
		ModelVersion: "legal-v1",
		BuildHash:    "abc1234",
		BuildDate:    "2026-06-13T00:00:00Z",
		StartedAt:    startedAt,
	}))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/version", http.NoBody))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	var body VersionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal version response: %v", err)
	}
	if body.Version != testVersion200 {
		t.Fatalf("version = %q, want 2.0.0", body.Version)
	}
	if body.Model != testModelID {
		t.Fatalf("model = %q, want %q", body.Model, testModelID)
	}
	if body.BuildHash != "abc1234" {
		t.Fatalf("build_hash = %q, want abc1234", body.BuildHash)
	}
	if body.UptimeSeconds <= 0 {
		t.Fatalf("uptime_seconds = %d, want > 0", body.UptimeSeconds)
	}
	if body.StartedAt == "" {
		t.Fatal("started_at should be present")
	}
}

func TestInfoHandlerReportsModelLifecycleState(t *testing.T) {
	state := lifecycle.NewState()
	runtime := &RuntimeHealth{Backend: testTEIBackend, Model: testModelID, Dimensions: 1024}
	r := gin.New()
	r.GET("/info", NewInfoHandler(buildinfo.Info{Version: testVersion200, Model: testModelID}, runtime, state))

	beforeWarmup := serveInfo(t, r)
	model := onlyModel(t, beforeWarmup)
	if model.Loaded {
		t.Fatal("model should not be loaded before warmup")
	}
	if model.WarmupDone {
		t.Fatal("warmup_done should be false before warmup")
	}
	assertModelStaticInfo(t, model)

	state.MarkWarmupDone()
	afterWarmup := serveInfo(t, r)
	model = onlyModel(t, afterWarmup)
	if !model.Loaded {
		t.Fatal("model should be loaded after warmup")
	}
	if !model.WarmupDone {
		t.Fatal("warmup_done should be true after warmup")
	}
}

func TestHealthcheckAliasUsesHealthResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	healthHandler := NewHealthHandler(&RuntimeHealth{Backend: testTEIBackend, Model: testModelID})
	r.GET("/health", healthHandler)
	r.GET("/v1/healthcheck", healthHandler)

	health := httptest.NewRecorder()
	r.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", http.NoBody))
	healthcheck := httptest.NewRecorder()
	r.ServeHTTP(healthcheck, httptest.NewRequest(http.MethodGet, "/v1/healthcheck", http.NoBody))

	if health.Code != http.StatusOK || healthcheck.Code != http.StatusOK {
		t.Fatalf("status /health=%d /v1/healthcheck=%d, want both 200", health.Code, healthcheck.Code)
	}
	var healthBody map[string]any
	if err := json.Unmarshal(health.Body.Bytes(), &healthBody); err != nil {
		t.Fatalf("unmarshal health response: %v", err)
	}
	var aliasBody map[string]any
	if err := json.Unmarshal(healthcheck.Body.Bytes(), &aliasBody); err != nil {
		t.Fatalf("unmarshal healthcheck response: %v", err)
	}
	if healthBody["status"] != aliasBody["status"] {
		t.Fatalf("status mismatch: health=%#v alias=%#v", healthBody["status"], aliasBody["status"])
	}
	if _, ok := aliasBody["runtime"].(map[string]any); !ok {
		t.Fatalf("healthcheck alias missing runtime: %#v", aliasBody["runtime"])
	}
}

func serveInfo(t *testing.T, r http.Handler) InfoResponse {
	t.Helper()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/info", http.NoBody))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	var body InfoResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal info response: %v", err)
	}
	return body
}

func onlyModel(t *testing.T, body InfoResponse) ModelInfo {
	t.Helper()
	if len(body.Models) != 1 {
		t.Fatalf("models len = %d, want 1", len(body.Models))
	}
	return body.Models[0]
}

func assertModelStaticInfo(t *testing.T, model ModelInfo) {
	t.Helper()
	if model.ID != testModelID {
		t.Fatalf("model id = %q, want %q", model.ID, testModelID)
	}
	if len(model.Dimensions) != 2 || model.Dimensions[0] != 512 || model.Dimensions[1] != 1024 {
		t.Fatalf("dims = %#v, want [512 1024]", model.Dimensions)
	}
	if model.MaxInputLengthTokens != defaultModelMaxInputLengthTokens {
		t.Fatalf("max_input_length_tokens = %d, want %d", model.MaxInputLengthTokens, defaultModelMaxInputLengthTokens)
	}
	if model.MaxBatchSize != defaultModelMaxBatchSize {
		t.Fatalf("max_batch_size = %d, want %d", model.MaxBatchSize, defaultModelMaxBatchSize)
	}
	if model.Device != defaultModelDevice {
		t.Fatalf("device = %q, want %q", model.Device, defaultModelDevice)
	}
}

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthHandlerDefaultShape(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", HealthHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
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
		Backend:           "tei",
		Model:             "deepvk/USER-bge-m3",
		Dimensions:        1024,
		ProductionDefault: true,
		CacheNamespace:    "m026-tei",
	}))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
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
	if runtime["backend"] != "tei" {
		t.Fatalf("backend = %#v, want tei", runtime["backend"])
	}
	if runtime["model"] != "deepvk/USER-bge-m3" {
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

func boolPtr(v bool) *bool { return &v }

func TestNewHealthHandlerIncludesSafeRuntimeMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", NewHealthHandler(&RuntimeHealth{
		Backend:                    "onnx",
		Model:                      "deepvk/USER-bge-m3",
		ArtifactID:                 "user-bge-m3-dense-fp32",
		Dimensions:                 1024,
		MaxSequenceLength:          1024,
		ValidatedMaxSequenceLength: 1024,
		ProductionDefault:          false,
		ArtifactVerified:            boolPtr(true),
		TokenizerVerified:          boolPtr(true),
		RuntimeLibraryVerified:     boolPtr(true),
		Provider:                   "CPUExecutionProvider",
		CacheNamespace:             "m026-test",
	}))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
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
	if runtime["backend"] != "onnx" {
		t.Fatalf("backend = %#v", runtime["backend"])
	}
	if runtime["artifact_id"] != "user-bge-m3-dense-fp32" {
		t.Fatalf("artifact_id = %#v", runtime["artifact_id"])
	}
	if runtime["provider"] != "CPUExecutionProvider" || runtime["tokenizer_verified"] != true || runtime["runtime_library_verified"] != true || runtime["artifact_verified"] != true {
		t.Fatalf("verification metadata = %#v", runtime)
	}
	if _, ok := runtime["manifest_path"]; ok {
		t.Fatal("runtime health must not expose manifest_path")
	}
	if _, ok := runtime["runtime_library_path"]; ok {
		t.Fatal("runtime health must not expose runtime_library_path")
	}
	if _, ok := runtime["tokenizer_path"]; ok {
		t.Fatal("runtime health must not expose tokenizer_path")
	}
}

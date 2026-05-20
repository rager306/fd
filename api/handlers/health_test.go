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
		ArtifactVerified:           true,
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

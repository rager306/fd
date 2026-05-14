package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get("http://localhost:8000/health")
	if err != nil {
		t.Skip("API not running on :8000")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", body["status"])
	}
}

func TestEmbeddingsEndpoint_ValidRequest(t *testing.T) {
	reqBody := map[string]interface{}{
		"model": "deepvk/USER-bge-m3",
		"input": "hello world",
	}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:8000/v1/embeddings", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Skip("API not running on :8000")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)

	if body["object"] != "list" {
		t.Errorf("expected object=list, got %v", body["object"])
	}
}

func TestEmbeddingsEndpoint_ArrayInput(t *testing.T) {
	reqBody := map[string]interface{}{
		"model": "deepvk/USER-bge-m3",
		"input": []string{"hello", "world"},
	}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:8000/v1/embeddings", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Skip("API not running on :8000")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestEmbeddingsEndpoint_InvalidJSON(t *testing.T) {
	resp, err := http.Post("http://localhost:8000/v1/embeddings", "application/json", bytes.NewBuffer([]byte(`{invalid}`)))
	if err != nil {
		t.Skip("API not running on :8000")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestEmbeddingsEndpoint_EmptyInput(t *testing.T) {
	reqBody := map[string]interface{}{
		"model": "deepvk/USER-bge-m3",
		"input": []string{},
	}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:8000/v1/embeddings", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Skip("API not running on :8000")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestEmbeddingsEndpoint_MissingModel(t *testing.T) {
	reqBody := map[string]interface{}{
		"input": "hello world",
	}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:8000/v1/embeddings", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Skip("API not running on :8000")
	}
	defer resp.Body.Close()

	// Should still work - model is optional in our handler
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

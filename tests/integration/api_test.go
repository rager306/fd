package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func baseURL() string {
	return strings.TrimRight(getEnv("FD_BASE_URL", "http://localhost:8000"), "/")
}

func requireAPIKey(t *testing.T) string {
	t.Helper()
	key := os.Getenv("FD_INTEGRATION_API_KEY")
	if key == "" {
		t.Skip("FD_INTEGRATION_API_KEY is required for protected endpoint integration tests")
	}
	return key
}

func newClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

func getIntegration(t *testing.T, path string) *http.Response {
	t.Helper()
	resp, err := newClient().Get(baseURL() + path)
	if err != nil {
		t.Skipf("API not running at %s: %v", baseURL(), err)
	}
	return resp
}

func getIntegrationAuth(t *testing.T, path string, apiKey string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, baseURL()+path, nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	resp, err := newClient().Do(req)
	if err != nil {
		t.Skipf("API not running at %s: %v", baseURL(), err)
	}
	return resp
}

func postIntegrationJSON(t *testing.T, path string, body any, apiKey string) *http.Response {
	t.Helper()
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL()+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := newClient().Do(req)
	if err != nil {
		t.Skipf("API not running at %s: %v", baseURL(), err)
	}
	return resp
}

func postIntegrationRaw(t *testing.T, path string, body []byte, apiKey string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, baseURL()+path, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := newClient().Do(req)
	if err != nil {
		t.Skipf("API not running at %s: %v", baseURL(), err)
	}
	return resp
}

func decodeJSON(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode JSON response: %v", err)
	}
	return body
}

func requireStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode == want {
		return
	}
	payload, _ := io.ReadAll(resp.Body)
	t.Fatalf("expected status %d, got %d body=%s", want, resp.StatusCode, string(payload))
}

func uniqueInput(label string) string {
	return fmt.Sprintf("m050 %s %d", label, time.Now().UnixNano())
}

func intFromJSONNumber(t *testing.T, value any) int {
	t.Helper()
	number, ok := value.(float64)
	if !ok {
		t.Fatalf("expected JSON number, got %T", value)
	}
	return int(number)
}

func embeddingData(t *testing.T, body map[string]any, index int) map[string]any {
	t.Helper()
	items, ok := body["data"].([]any)
	if !ok || len(items) <= index {
		t.Fatalf("expected data[%d] in embeddings response, got %T len=%d", index, body["data"], len(items))
	}
	item, ok := items[index].(map[string]any)
	if !ok {
		t.Fatalf("expected data[%d] object, got %T", index, items[index])
	}
	return item
}

func assertEmbeddingDimensions(t *testing.T, item map[string]any, want int) {
	t.Helper()
	if got := intFromJSONNumber(t, item["dimensions"]); got != want {
		t.Fatalf("expected dimensions=%d, got %d", want, got)
	}
	embedding, ok := item["embedding"].([]any)
	if !ok {
		t.Fatalf("expected float embedding array, got %T", item["embedding"])
	}
	if len(embedding) != want {
		t.Fatalf("expected embedding length %d, got %d", want, len(embedding))
	}
}

func assertCacheHeader(t *testing.T, resp *http.Response, want string) {
	t.Helper()
	if got := resp.Header.Get("X-Cache"); !strings.EqualFold(got, want) {
		t.Fatalf("expected X-Cache=%s, got %q", want, got)
	}
}

func TestPublicRuntimeDiagnostics(t *testing.T) {
	for _, path := range []string{"/live", "/ready"} {
		resp := getIntegration(t, path)
		requireStatus(t, resp, http.StatusOK)
		resp.Body.Close()
	}

	resp := getIntegration(t, "/health")
	defer resp.Body.Close()
	requireStatus(t, resp, http.StatusOK)
	health := decodeJSON(t, resp)
	if health["status"] != "ok" {
		t.Fatalf("expected health status=ok, got %v", health["status"])
	}
	if _, ok := health["runtime"].(map[string]any); !ok {
		t.Fatalf("expected health.runtime object, got %T", health["runtime"])
	}
	if _, ok := health["dependencies"].(map[string]any); !ok {
		t.Fatalf("expected health.dependencies object, got %T", health["dependencies"])
	}
	if _, ok := health["in_flight_capacity"].(float64); !ok {
		t.Fatalf("expected health.in_flight_capacity number, got %T", health["in_flight_capacity"])
	}

}

func TestAuthenticatedMetricsDiagnostics(t *testing.T) {
	apiKey := requireAPIKey(t)
	metricsResp := getIntegrationAuth(t, "/metrics", apiKey)
	defer metricsResp.Body.Close()
	requireStatus(t, metricsResp, http.StatusOK)
	metricsBody, err := io.ReadAll(metricsResp.Body)
	if err != nil {
		t.Fatalf("read metrics body: %v", err)
	}
	metrics := string(metricsBody)
	for _, needle := range []string{"fd_in_flight_requests", "fd_in_flight_capacity", `fd_cache_entries{tier="l1"}`} {
		if !strings.Contains(metrics, needle) {
			t.Fatalf("expected metrics to contain %s", needle)
		}
	}
}

func TestProtectedEndpointsRequireAuth(t *testing.T) {
	for _, tc := range []struct {
		name string
		path string
		body any
	}{
		{name: "embeddings", path: "/v1/embeddings", body: map[string]any{"model": "deepvk/USER-bge-m3", "input": "hello world"}},
		{name: "cache flush", path: "/v1/cache/flush", body: map[string]any{}},
		{name: "cache delete", path: "/v1/cache/delete", body: map[string]any{"input": "hello world"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := postIntegrationJSON(t, tc.path, tc.body, "")
			defer resp.Body.Close()
			requireStatus(t, resp, http.StatusUnauthorized)
		})
	}
}

func TestAuthenticatedEmbeddingDimensionsAndBatch(t *testing.T) {
	apiKey := requireAPIKey(t)
	dimensions := 512
	resp := postIntegrationJSON(t, "/v1/embeddings", map[string]any{
		"model":      "deepvk/USER-bge-m3",
		"input":      []string{uniqueInput("dimensions-a"), uniqueInput("dimensions-b")},
		"dimensions": dimensions,
	}, apiKey)
	defer resp.Body.Close()
	requireStatus(t, resp, http.StatusOK)
	assertCacheHeader(t, resp, "MISS")

	body := decodeJSON(t, resp)
	if body["object"] != "list" {
		t.Fatalf("expected object=list, got %v", body["object"])
	}
	assertEmbeddingDimensions(t, embeddingData(t, body, 0), dimensions)
	assertEmbeddingDimensions(t, embeddingData(t, body, 1), dimensions)
}

func TestAuthenticatedInvalidRequestsUseCurrentValidation(t *testing.T) {
	apiKey := requireAPIKey(t)

	invalidJSON := postIntegrationRaw(t, "/v1/embeddings", []byte(`{invalid}`), apiKey)
	defer invalidJSON.Body.Close()
	requireStatus(t, invalidJSON, http.StatusBadRequest)

	emptyInput := postIntegrationJSON(t, "/v1/embeddings", map[string]any{
		"model": "deepvk/USER-bge-m3",
		"input": []string{},
	}, apiKey)
	defer emptyInput.Body.Close()
	requireStatus(t, emptyInput, http.StatusBadRequest)

	missingModel := postIntegrationJSON(t, "/v1/embeddings", map[string]any{
		"input": uniqueInput("missing-model"),
	}, apiKey)
	defer missingModel.Body.Close()
	requireStatus(t, missingModel, http.StatusOK)
}

func TestAuthenticatedEmbeddingCacheInvalidation(t *testing.T) {
	apiKey := requireAPIKey(t)
	inputForFlush := uniqueInput("flush")
	inputForDelete := uniqueInput("delete")

	flush := postIntegrationJSON(t, "/v1/cache/flush", map[string]any{}, apiKey)
	defer flush.Body.Close()
	requireStatus(t, flush, http.StatusOK)

	first := postIntegrationJSON(t, "/v1/embeddings", map[string]any{"input": inputForFlush}, apiKey)
	requireStatus(t, first, http.StatusOK)
	assertCacheHeader(t, first, "MISS")
	first.Body.Close()

	second := postIntegrationJSON(t, "/v1/embeddings", map[string]any{"input": inputForFlush}, apiKey)
	requireStatus(t, second, http.StatusOK)
	assertCacheHeader(t, second, "HIT")
	second.Body.Close()

	flush = postIntegrationJSON(t, "/v1/cache/flush", map[string]any{}, apiKey)
	requireStatus(t, flush, http.StatusOK)
	flushBody := decodeJSON(t, flush)
	flush.Body.Close()
	if flushBody["flushed"] != true {
		t.Fatalf("expected flushed=true, got %v", flushBody["flushed"])
	}

	afterFlush := postIntegrationJSON(t, "/v1/embeddings", map[string]any{"input": inputForFlush}, apiKey)
	requireStatus(t, afterFlush, http.StatusOK)
	assertCacheHeader(t, afterFlush, "MISS")
	afterFlush.Body.Close()

	deleteFirst := postIntegrationJSON(t, "/v1/embeddings", map[string]any{"input": inputForDelete}, apiKey)
	requireStatus(t, deleteFirst, http.StatusOK)
	assertCacheHeader(t, deleteFirst, "MISS")
	deleteFirst.Body.Close()

	deleteSecond := postIntegrationJSON(t, "/v1/embeddings", map[string]any{"input": inputForDelete}, apiKey)
	requireStatus(t, deleteSecond, http.StatusOK)
	assertCacheHeader(t, deleteSecond, "HIT")
	deleteSecond.Body.Close()

	deleteResp := postIntegrationJSON(t, "/v1/cache/delete", map[string]any{"input": inputForDelete, "dimensions": 1024}, apiKey)
	requireStatus(t, deleteResp, http.StatusOK)
	deleteBody := decodeJSON(t, deleteResp)
	deleteResp.Body.Close()
	if got := intFromJSONNumber(t, deleteBody["deleted"]); got != 1 {
		t.Fatalf("expected deleted=1, got %d", got)
	}

	afterDelete := postIntegrationJSON(t, "/v1/embeddings", map[string]any{"input": inputForDelete}, apiKey)
	requireStatus(t, afterDelete, http.StatusOK)
	assertCacheHeader(t, afterDelete, "MISS")
	afterDelete.Body.Close()
}

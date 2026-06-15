package embed

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestTEIClientEmbedSendsBatchAndReturnsEmbeddings(t *testing.T) {
	var captured teiRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/embeddings" {
			t.Fatalf("path = %q, want /embeddings", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("method = %q, want POST", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[{"embedding":[1,2,3]},{"embedding":[4,5,6]}]}`))
	}))
	defer server.Close()

	client := NewTEIClient(server.URL, "model-a", server.Client())
	got, err := client.Embed(context.Background(), []string{"first", "second"})
	if err != nil {
		t.Fatalf("Embed returned error: %v", err)
	}
	if captured.Model != "model-a" {
		t.Fatalf("request model = %q, want model-a", captured.Model)
	}
	if !captured.Truncate {
		t.Fatal("request truncate = false, want true")
	}
	inputs, ok := captured.Input.([]any)
	if !ok {
		t.Fatalf("request input type = %T, want []any", captured.Input)
	}
	if got, want := len(inputs), 2; got != want {
		t.Fatalf("request input length = %d, want %d", got, want)
	}
	if inputs[0] != "first" || inputs[1] != "second" {
		t.Fatalf("request input = %#v", inputs)
	}

	if gotLen, wantLen := len(got), 2; gotLen != wantLen {
		t.Fatalf("embedding count = %d, want %d", gotLen, wantLen)
	}
	if got[0][2] != 3 || got[1][0] != 4 {
		t.Fatalf("unexpected embeddings: %#v", got)
	}
}

func TestTEIClientEmbedRejectsEmptyInput(t *testing.T) {
	client := NewTEIClient("http://example.invalid", "model-a", http.DefaultClient)
	if _, err := client.Embed(context.Background(), nil); err == nil {
		t.Fatal("Embed(nil) error = nil, want error")
	}
}

func TestTEIClientEmbedReturnsStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewTEIClient(server.URL, "model-a", server.Client())
	if _, err := client.Embed(context.Background(), []string{"x"}); err == nil {
		t.Fatal("Embed status error = nil, want error")
	}
}

func TestTEIClientEmbedRetriesRetriableStatus(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempt := attempts.Add(1)
		if attempt == 1 {
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[{"embedding":[1,2,3]}]}`))
	}))
	defer server.Close()

	client := NewTEIClient(server.URL, "model-a", server.Client())
	got, err := client.Embed(context.Background(), []string{"x"})
	if err != nil {
		t.Fatalf("Embed returned error after retry: %v", err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("attempts = %d, want 2", attempts.Load())
	}
	if len(got) != 1 || len(got[0]) != 3 {
		t.Fatalf("embeddings = %#v, want one 3-dim vector", got)
	}
}

func TestTEIClientEmbedDoesNotRetryBadRequest(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewTEIClient(server.URL, "model-a", server.Client())
	if _, err := client.Embed(context.Background(), []string{"x"}); err == nil {
		t.Fatal("Embed bad request error = nil, want error")
	}
	if attempts.Load() != 1 {
		t.Fatalf("attempts = %d, want 1", attempts.Load())
	}
}

func TestTEIClientEmbedFastFailsAfterRepeatedRetriableFailures(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewTEIClient(server.URL, "model-a", server.Client())
	for i := 0; i < 3; i++ {
		if _, err := client.Embed(context.Background(), []string{"x"}); err == nil {
			t.Fatal("Embed retriable failure = nil, want error")
		}
	}
	beforeFastFail := attempts.Load()

	_, err := client.Embed(context.Background(), []string{"x"})
	if err == nil {
		t.Fatal("Embed circuit-open error = nil, want error")
	}
	if !strings.Contains(err.Error(), "TEI circuit open") {
		t.Fatalf("error = %v, want TEI circuit open", err)
	}
	if attempts.Load() != beforeFastFail {
		t.Fatalf("attempts after circuit open = %d, want unchanged %d", attempts.Load(), beforeFastFail)
	}
}

func TestTEIClientEmbedReturnsDecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	client := NewTEIClient(server.URL, "model-a", server.Client())
	if _, err := client.Embed(context.Background(), []string{"x"}); err == nil {
		t.Fatal("Embed decode error = nil, want error")
	}
}

func TestTEIClientEmbedReturnsEmptyDataError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[]}`))
	}))
	defer server.Close()

	client := NewTEIClient(server.URL, "model-a", server.Client())
	if _, err := client.Embed(context.Background(), []string{"x"}); err == nil {
		t.Fatal("Embed empty data error = nil, want error")
	}
}

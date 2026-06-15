package embed

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// TEIClient calls a Text Embeddings Inference (TEI) OpenAI-compatible
// endpoint and converts its nested response into fd's [][]float32 interface.
type TEIClient struct {
	baseURL    string
	modelID    string
	httpClient *http.Client

	retryMaxAttempts        int
	retryBaseBackoff        time.Duration
	circuitFailureThreshold int
	circuitCooldown         time.Duration
	now                     func() time.Time
	sleep                   func(context.Context, time.Duration) error

	mu                  sync.Mutex
	consecutiveFailures int
	circuitOpenUntil    time.Time
}

// NewTEIClient returns a TEIClient for baseURL/modelID using the supplied
// HTTP client. The caller owns the HTTP client's timeout and transport settings.
func NewTEIClient(baseURL, modelID string, client *http.Client) *TEIClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &TEIClient{
		baseURL:                 baseURL,
		modelID:                 modelID,
		httpClient:              client,
		retryMaxAttempts:        3,
		retryBaseBackoff:        10 * time.Millisecond,
		circuitFailureThreshold: 3,
		circuitCooldown:         5 * time.Second,
		now:                     time.Now,
		sleep:                   sleepWithContext,
	}
}

// TEI OpenAI-compatible request
type teiRequest struct {
	Input    any    `json:"input"` // string OR []string (OpenAI spec)
	Model    string `json:"model"`
	Truncate bool   `json:"truncate"`
}

// TEI OpenAI-compatible response — NESTED structure (data[0].embedding)
type teiResponse struct {
	Object string            `json:"object"`
	Data   []teiEmbeddingObj `json:"data"`
	Model  string            `json:"model"`
	Usage  Usage             `json:"usage"`
}

type teiEmbeddingObj struct {
	Embedding []float32 `json:"embedding"`
}

// Embed sends texts to TEI's OpenAI-compatible embeddings endpoint and
// returns the embedding vectors in response order.
func (c *TEIClient) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts provided")
	}
	if err := c.checkCircuitOpen(); err != nil {
		return nil, err
	}

	// Send the FULL slice as TEI input. TEI accepts both single string and
	// []string via the OpenAI-compatible endpoint (the spec's "input" field
	// is defined as either). Previously this code used texts[0] only,
	// which silently dropped batch elements when the caller passed >1
	// input — see S01-R01 for the bug pattern. Fixed 2026-06 when fd v2
	// raised maxBatchSize from 32 to 128, which made the per-item
	// chunking in the handler start sending multi-element chunks.
	reqBody := teiRequest{
		Input:    texts, // marshal as JSON []string when slice, single string when len==1
		Model:    c.modelID,
		Truncate: true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	maxAttempts := c.retryMaxAttempts
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		embeddings, retriable, err := c.doEmbedRequest(ctx, jsonBody)
		if err == nil {
			c.recordTEISuccess()
			return embeddings, nil
		}
		lastErr = err
		if !retriable || attempt == maxAttempts {
			c.recordTEIFailure(retriable)
			if retriable && attempt == maxAttempts {
				return nil, fmt.Errorf("TEI retry exhausted after %d attempts: %w", attempt, err)
			}
			return nil, err
		}
		if err := c.sleep(ctx, c.retryBackoff(attempt)); err != nil {
			c.recordTEIFailure(true)
			return nil, fmt.Errorf("TEI retry backoff: %w", err)
		}
	}
	return nil, lastErr
}

func (c *TEIClient) doEmbedRequest(ctx context.Context, jsonBody []byte) ([][]float32, bool, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/embeddings", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || ctx.Err() != nil {
			return nil, false, err
		}
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, isRetriableTEIStatus(resp.StatusCode), fmt.Errorf("TEI returned status %d", resp.StatusCode)
	}

	var result teiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, false, fmt.Errorf("decode error: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, false, fmt.Errorf("TEI returned empty data")
	}

	embeddings := make([][]float32, len(result.Data))
	for i, d := range result.Data {
		embeddings[i] = d.Embedding
	}

	return embeddings, false, nil
}

func isRetriableTEIStatus(status int) bool {
	switch status {
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func (c *TEIClient) retryBackoff(attempt int) time.Duration {
	if c.retryBaseBackoff <= 0 {
		return 0
	}
	return c.retryBaseBackoff << max(attempt-1, 0)
}

func (c *TEIClient) checkCircuitOpen() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.circuitOpenUntil.IsZero() {
		return nil
	}
	now := c.now()
	if now.Before(c.circuitOpenUntil) {
		return fmt.Errorf("TEI circuit open until %s after %d consecutive retriable failures", c.circuitOpenUntil.Format(time.RFC3339Nano), c.consecutiveFailures)
	}
	c.circuitOpenUntil = time.Time{}
	c.consecutiveFailures = 0
	return nil
}

func (c *TEIClient) recordTEISuccess() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.consecutiveFailures = 0
	c.circuitOpenUntil = time.Time{}
}

func (c *TEIClient) recordTEIFailure(retriable bool) {
	if !retriable {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.consecutiveFailures++
	if c.circuitFailureThreshold > 0 && c.consecutiveFailures >= c.circuitFailureThreshold {
		c.circuitOpenUntil = c.now().Add(c.circuitCooldown)
	}
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

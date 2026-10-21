package main

import (
	"context"
	"errors"
	"fd-api/cache"
	"fd-api/internal/envutil"
	"fd-api/lifecycle"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRouterDoesNotTrustForwardedForByDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	configureTrustedProxies(r)
	r.GET("/client-ip", func(c *gin.Context) { c.String(http.StatusOK, c.ClientIP()) })

	req := httptest.NewRequest(http.MethodGet, "/client-ip", http.NoBody)
	req.RemoteAddr = "198.51.100.10:1234"
	req.Header.Set("X-Forwarded-For", "203.0.113.99")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Body.String(); got != "198.51.100.10" {
		t.Fatalf("ClientIP = %q, want direct remote address", got)
	}
}

type warmupModelFunc func(ctx context.Context, texts []string) ([][]float32, error)

func (f warmupModelFunc) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return f(ctx, texts)
}

func TestStartModelWarmupMarksStateReady(t *testing.T) {
	state := lifecycle.NewState()
	model := warmupModelFunc(func(_ context.Context, _ []string) ([][]float32, error) {
		return [][]float32{{1}}, nil
	})
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	startModelWarmup(logger, state, model, time.Second)
	waitForCondition(t, state.IsReady)
}

func TestStartModelWarmupStoresError(t *testing.T) {
	state := lifecycle.NewState()
	boom := errors.New("boom")
	model := warmupModelFunc(func(_ context.Context, _ []string) ([][]float32, error) {
		return nil, boom
	})
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	startModelWarmup(logger, state, model, time.Second)
	waitForCondition(t, func() bool {
		return errors.Is(state.LastError(), boom)
	})
	if state.IsReady() {
		t.Fatal("state should not be ready after warmup failure")
	}
}

func TestStartModelWarmupRetriesAndMarksReady(t *testing.T) {
	state := lifecycle.NewState()
	boom := errors.New("boom")
	var attempts atomic.Int32
	model := warmupModelFunc(func(_ context.Context, _ []string) ([][]float32, error) {
		if attempts.Add(1) == 1 {
			return nil, boom
		}
		return [][]float32{{1}}, nil
	})
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	startModelWarmupWithPolicy(logger, state, model, time.Second, warmupRetryPolicy{
		maxAttempts: 3,
		backoff:     func(int) time.Duration { return 0 },
	})

	waitForCondition(t, state.IsReady)
	if attempts.Load() != 2 {
		t.Fatalf("attempts = %d, want 2", attempts.Load())
	}
	if err := state.LastError(); err != nil {
		t.Fatalf("LastError after successful retry = %v, want nil", err)
	}
}

func TestStartModelWarmupRecordsTerminalErrorAfterMaxAttempts(t *testing.T) {
	state := lifecycle.NewState()
	boom := errors.New("boom")
	var attempts atomic.Int32
	model := warmupModelFunc(func(_ context.Context, _ []string) ([][]float32, error) {
		attempts.Add(1)
		return nil, boom
	})
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	startModelWarmupWithPolicy(logger, state, model, time.Second, warmupRetryPolicy{
		maxAttempts: 3,
		backoff:     func(int) time.Duration { return 0 },
	})

	waitForCondition(t, func() bool { return attempts.Load() == 3 })
	if !errors.Is(state.LastError(), boom) {
		t.Fatalf("LastError = %v, want boom", state.LastError())
	}
	if state.IsReady() {
		t.Fatal("state should not be ready after terminal warmup failure")
	}
}

func TestReportHTTPServerErrorIgnoresWrappedServerClosed(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	shutdownCh := make(chan os.Signal, 1)

	reportHTTPServerError(logger, "127.0.0.1:0", func() error {
		return fmt.Errorf("wrapped: %w", http.ErrServerClosed)
	}, shutdownCh)

	select {
	case sig := <-shutdownCh:
		t.Fatalf("unexpected shutdown signal for ErrServerClosed: %v", sig)
	default:
	}
}

func TestReportHTTPServerErrorSignalsFatalError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	shutdownCh := make(chan os.Signal, 1)
	boom := errors.New("accept failed")

	reportHTTPServerError(logger, "127.0.0.1:0", func() error { return boom }, shutdownCh)

	select {
	case sig := <-shutdownCh:
		serverSig, ok := sig.(serverErrorSignal)
		if !ok {
			t.Fatalf("signal type = %T, want serverErrorSignal", sig)
		}
		if !errors.Is(serverSig.err, boom) {
			t.Fatalf("signal error = %v, want %v", serverSig.err, boom)
		}
		if got := sig.String(); got != "server_error" {
			t.Fatalf("signal string = %q, want server_error", got)
		}
	default:
		t.Fatal("expected fatal listener error to trigger shutdown signal")
	}
}

func TestNewHTTPServerDefaults(t *testing.T) {
	t.Setenv("FD_HTTP_READ_TIMEOUT", "")
	t.Setenv("FD_HTTP_WRITE_TIMEOUT", "")
	t.Setenv("FD_HTTP_IDLE_TIMEOUT", "")

	handler := http.NewServeMux()
	srv := newHTTPServer("127.0.0.1:8000", handler)
	if srv.Addr != "127.0.0.1:8000" {
		t.Fatalf("Addr = %q, want 127.0.0.1:8000", srv.Addr)
	}
	if srv.Handler != handler {
		t.Fatal("Handler was not attached")
	}
	if srv.ReadHeaderTimeout != defaultHTTPReadHeaderTimeout {
		t.Fatalf("ReadHeaderTimeout = %s, want %s", srv.ReadHeaderTimeout, defaultHTTPReadHeaderTimeout)
	}
	if srv.ReadTimeout != defaultHTTPReadTimeout {
		t.Fatalf("ReadTimeout = %s, want %s", srv.ReadTimeout, defaultHTTPReadTimeout)
	}
	if srv.WriteTimeout != defaultHTTPWriteTimeout {
		t.Fatalf("WriteTimeout = %s, want %s", srv.WriteTimeout, defaultHTTPWriteTimeout)
	}
	if srv.IdleTimeout != defaultHTTPIdleTimeout {
		t.Fatalf("IdleTimeout = %s, want %s", srv.IdleTimeout, defaultHTTPIdleTimeout)
	}
}

func TestNewHTTPServerTimeoutsExceedHandlerBudget(t *testing.T) {
	t.Setenv("FD_HTTP_READ_TIMEOUT", "")
	t.Setenv("FD_HTTP_WRITE_TIMEOUT", "")

	const handlerBudget = 30 * time.Second
	srv := newHTTPServer("127.0.0.1:8000", http.NewServeMux())
	if srv.ReadTimeout <= handlerBudget {
		t.Fatalf("ReadTimeout = %s, want > handler budget %s", srv.ReadTimeout, handlerBudget)
	}
	// Go starts the write deadline after the headers are read, so it must cover a
	// full body read on top of the handler budget; otherwise a slow body eats the
	// handler's time.
	if srv.WriteTimeout < srv.ReadTimeout+handlerBudget {
		t.Fatalf("WriteTimeout = %s, want >= ReadTimeout %s + handler budget %s", srv.WriteTimeout, srv.ReadTimeout, handlerBudget)
	}
}

func TestNewHTTPServerEnvOverrides(t *testing.T) {
	t.Setenv("FD_HTTP_READ_TIMEOUT", "45s")
	t.Setenv("FD_HTTP_WRITE_TIMEOUT", "50s")
	t.Setenv("FD_HTTP_IDLE_TIMEOUT", "90s")

	srv := newHTTPServer("127.0.0.1:8000", http.NewServeMux())
	if srv.ReadTimeout != 45*time.Second {
		t.Fatalf("ReadTimeout = %s, want 45s", srv.ReadTimeout)
	}
	if srv.WriteTimeout != 50*time.Second {
		t.Fatalf("WriteTimeout = %s, want 50s", srv.WriteTimeout)
	}
	if srv.IdleTimeout != 90*time.Second {
		t.Fatalf("IdleTimeout = %s, want 90s", srv.IdleTimeout)
	}
}

func TestNewHTTPServerNonPositiveEnvFallsBack(t *testing.T) {
	t.Setenv("FD_HTTP_READ_TIMEOUT", "0s")
	t.Setenv("FD_HTTP_WRITE_TIMEOUT", "-5s")
	t.Setenv("FD_HTTP_IDLE_TIMEOUT", "bogus")

	srv := newHTTPServer("127.0.0.1:8000", http.NewServeMux())
	if srv.ReadTimeout != defaultHTTPReadTimeout {
		t.Fatalf("ReadTimeout = %s, want default %s", srv.ReadTimeout, defaultHTTPReadTimeout)
	}
	if srv.WriteTimeout != defaultHTTPWriteTimeout {
		t.Fatalf("WriteTimeout = %s, want default %s", srv.WriteTimeout, defaultHTTPWriteTimeout)
	}
	if srv.IdleTimeout != defaultHTTPIdleTimeout {
		t.Fatalf("IdleTimeout = %s, want default %s", srv.IdleTimeout, defaultHTTPIdleTimeout)
	}
}

func TestHTTPServerReadTimeoutCutsSlowClient(t *testing.T) {
	t.Setenv("FD_HTTP_READ_TIMEOUT", "150ms")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})
	srv := newHTTPServer("127.0.0.1:0", handler)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	serveDone := make(chan struct{})
	go func() {
		defer close(serveDone)
		_ = srv.Serve(listener)
	}()
	t.Cleanup(func() {
		_ = srv.Close()
		<-serveDone
	})

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// Announce a 100-byte body, send one byte, and keep the write side open.
	partial := "POST /v1/embeddings HTTP/1.1\r\nHost: x\r\nContent-Length: 100\r\n\r\nx"
	if _, err := io.WriteString(conn, partial); err != nil {
		t.Fatalf("write partial request: %v", err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("set read deadline: %v", err)
	}
	start := time.Now()
	if _, err := io.Copy(io.Discard, conn); err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			t.Fatalf("server did not close the slow connection within 3s (waited %s)", time.Since(start))
		}
		// A reset is also an acceptable close.
	}
	if elapsed := time.Since(start); elapsed >= 3*time.Second {
		t.Fatalf("slow connection survived %s, want a close near the 150ms read timeout", elapsed)
	}
}

func waitForCondition(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("condition not met within 1s")
}

func TestLoadEmbeddingRuntimeConfigDefaultsToTEI(t *testing.T) {
	t.Setenv("EMBEDDING_BACKEND", "")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", "/tmp/stale-manifest.json")

	config, err := loadEmbeddingRuntimeConfig()
	if err != nil {
		t.Fatalf("loadEmbeddingRuntimeConfig default returned error: %v", err)
	}
	if config.Backend != embeddingBackendTEI {
		t.Fatalf("backend = %q, want %q", config.Backend, embeddingBackendTEI)
	}
}

func TestLoadEmbeddingRuntimeConfigRejectsInvalidBackend(t *testing.T) {
	t.Setenv("EMBEDDING_BACKEND", "candle")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected invalid backend error")
	}
}

func TestLoadEmbeddingRuntimeConfigRejectsONNXBackend(t *testing.T) {
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", "/tmp/stale-manifest.json")
	t.Setenv("ONNX_RUNTIME_LIBRARY", "/tmp/libonnxruntime.so")
	t.Setenv("ONNX_TOKENIZER_PATH", "/tmp/tokenizer.json")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected ONNX backend to be disabled")
	}
	if got := err.Error(); !strings.Contains(got, "TEI only") && !strings.Contains(got, "supports TEI only") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEmbeddingRuntimeConfigHealthReturnsSafeTEIMetadata(t *testing.T) {
	config := &embeddingRuntimeConfig{Backend: embeddingBackendTEI}
	health := config.Health("deepvk/USER-bge-m3", "v2")
	if health == nil {
		t.Fatal("TEI health metadata should not be nil")
	}
	if health.Backend != "tei" {
		t.Fatalf("backend = %q, want tei", health.Backend)
	}
	if health.Model != "deepvk/USER-bge-m3" {
		t.Fatalf("model = %q", health.Model)
	}
	if health.Dimensions != 1024 {
		t.Fatalf("dimensions = %d, want 1024", health.Dimensions)
	}
	if !health.ProductionDefault {
		t.Fatal("production_default should be true for TEI")
	}
	if health.CacheNamespace != "v2" {
		t.Fatalf("cache_namespace = %q", health.CacheNamespace)
	}
}

func TestWarmupRetryPolicyFromEnvDefaultsTo5(t *testing.T) {
	policy := warmupRetryPolicyFromEnv()
	if policy.maxAttempts != 5 {
		t.Fatalf("maxAttempts = %d, want 5", policy.maxAttempts)
	}
	if policy.backoff(1) != 5*time.Second {
		t.Fatalf("backoff(1) = %s, want 5s", policy.backoff(1))
	}
}

func TestCacheConfigFromEnv(t *testing.T) {
	t.Setenv("FD_CACHE_MAX_SIZE", "500")
	t.Setenv("FD_CACHE_LOCAL_TTL", "10s")

	size := envutil.PositiveInt("FD_CACHE_MAX_SIZE", 10000)
	ttl := envutil.DurationOrDefault("FD_CACHE_LOCAL_TTL", 30*time.Second)

	if size != 500 {
		t.Fatalf("FD_CACHE_MAX_SIZE=500, got %d", size)
	}
	if ttl != 10*time.Second {
		t.Fatalf("FD_CACHE_LOCAL_TTL=10s, got %s", ttl)
	}
}

func TestCacheConfigDefaultsWithoutEnv(t *testing.T) {
	t.Setenv("FD_CACHE_MAX_SIZE", "")
	t.Setenv("FD_CACHE_LOCAL_TTL", "")

	size := envutil.PositiveInt("FD_CACHE_MAX_SIZE", 10000)
	ttl := envutil.DurationOrDefault("FD_CACHE_LOCAL_TTL", 30*time.Second)

	if size != 10000 {
		t.Fatalf("default FD_CACHE_MAX_SIZE = %d, want 10000", size)
	}
	if ttl != 30*time.Second {
		t.Fatalf("default FD_CACHE_LOCAL_TTL = %s, want 30s", ttl)
	}
}

func TestLocalCacheWithSmallSize(t *testing.T) {
	c := cache.NewLocalCache(5, 30*time.Second)
	defer func() { _ = c.Close() }()
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		c.Set(ctx, fmt.Sprintf("k%d", i), []byte{byte(i)}, 30*time.Second)
	}
	if s := c.Size(); s > 5 {
		t.Fatalf("LocalCache size = %d, want <= 5", s)
	}
}

func TestLocalCacheWithZeroSize(t *testing.T) {
	c := cache.NewLocalCache(0, 30*time.Second)
	defer func() { _ = c.Close() }()
	ctx := context.Background()
	for i := 0; i < 20; i++ {
		c.Set(ctx, fmt.Sprintf("k%d", i), []byte{byte(i)}, 30*time.Second)
	}
	if s := c.Size(); s != 20 {
		t.Fatalf("LocalCache size=%d with maxSize=0, want 20 (unlimited)", s)
	}
}

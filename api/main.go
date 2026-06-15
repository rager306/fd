// Package main starts the fd embedding service: loads TEI runtime config and serves /v1/embeddings + observability endpoints on the configured port.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"fd-api/buildinfo"
	"fd-api/cache"
	"fd-api/embed"
	"fd-api/handlers"
	"fd-api/internal/envutil"
	"fd-api/lifecycle"
	"fd-api/middleware"
	"fd-api/observability"

	"github.com/gin-gonic/gin"
)

// Version is injected by release builds with -ldflags "-X main.Version=...".
var Version = buildinfo.DefaultVersion

// BuildHash is injected by release builds with -ldflags "-X main.BuildHash=...".
var BuildHash = buildinfo.DefaultBuildHash

// BuildDate is injected by release builds with -ldflags "-X main.BuildDate=...".
var BuildDate = buildinfo.DefaultBuildDate

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func configureTrustedProxies(r *gin.Engine) {
	// Fail closed for client identity: without an explicit trusted-proxy
	// policy, Gin must not trust spoofable X-Forwarded-For headers for
	// ClientIP(). This keeps rate limiting keyed to the direct peer by default.
	if err := r.SetTrustedProxies(nil); err != nil {
		panic("configure trusted proxies: " + err.Error())
	}
}

func getLogLevel(value string) slog.Level {
	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type embeddingBackend string

const embeddingBackendTEI embeddingBackend = "tei"

type embeddingRuntimeConfig struct {
	Backend embeddingBackend
}

// Health returns safe runtime metadata for the active TEI backend.
// TEI returns backend, model, fixed dimensions, production-default flag,
// and cache namespace — no internal URLs, paths, or secrets.
func (c *embeddingRuntimeConfig) Health(modelID, cacheNamespace string) *handlers.RuntimeHealth {
	if c == nil || c.Backend != embeddingBackendTEI {
		return nil
	}
	return &handlers.RuntimeHealth{
		Backend:           string(c.Backend),
		Model:             modelID,
		Dimensions:        1024, // deepvk/USER-bge-m3 is 1024-dimensional
		ProductionDefault: true,
		CacheNamespace:    cacheNamespace,
	}
}

func loadEmbeddingRuntimeConfig() (*embeddingRuntimeConfig, error) {
	backend := embeddingBackend(strings.ToLower(getEnv("EMBEDDING_BACKEND", string(embeddingBackendTEI))))
	if backend == "" {
		backend = embeddingBackendTEI
	}
	if backend != embeddingBackendTEI {
		return nil, fmt.Errorf("EMBEDDING_BACKEND=%q is not supported; fd currently supports TEI only", backend)
	}
	return &embeddingRuntimeConfig{Backend: embeddingBackendTEI}, nil
}

const defaultWarmupTimeout = 5 * time.Second

func closeResource(name string, resource interface{ Close() error }, logger *slog.Logger) {
	if resource == nil {
		return
	}
	if err := resource.Close(); err != nil {
		logger.Warn(name+" close failed", "error", err)
	}
}

type serverErrorSignal struct {
	err error
}

func (serverErrorSignal) Signal() {}

func (serverErrorSignal) String() string { return "server_error" }

func reportHTTPServerError(logger *slog.Logger, addr string, listen func() error, shutdownCh chan<- os.Signal) {
	logger.Info("listening", "addr", addr)
	if err := listen(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		logger.Error("server error", "error", err)
		select {
		case shutdownCh <- serverErrorSignal{err: err}:
		default:
			logger.Error("server error shutdown signal dropped", "error", err)
		}
	}
}

type warmupRetryPolicy struct {
	maxAttempts int
	backoff     func(attempt int) time.Duration
}

func defaultWarmupRetryPolicy() warmupRetryPolicy {
	return warmupRetryPolicy{
		maxAttempts: 3,
		backoff: func(attempt int) time.Duration {
			return 2 * time.Second << max(attempt-1, 0)
		},
	}
}

func startModelWarmup(logger *slog.Logger, state *lifecycle.State, model embed.Embedder, timeout time.Duration) {
	startModelWarmupWithPolicy(logger, state, model, timeout, defaultWarmupRetryPolicy())
}

func startModelWarmupWithPolicy(logger *slog.Logger, state *lifecycle.State, model embed.Embedder, timeout time.Duration, policy warmupRetryPolicy) {
	if policy.maxAttempts < 1 {
		policy.maxAttempts = 1
	}
	if policy.backoff == nil {
		policy.backoff = func(int) time.Duration { return 0 }
	}

	go func() {
		started := time.Now()
		logger.Info("model warmup started", "timeout", timeout.String(), "max_attempts", policy.maxAttempts)
		var lastErr error
		for attempt := 1; attempt <= policy.maxAttempts; attempt++ {
			attemptStarted := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := lifecycle.PreWarm(ctx, model)
			cancel()
			if err == nil {
				state.MarkWarmupDone()
				logger.Info("model warmup done", "attempt", attempt, "latency_ms", time.Since(started).Milliseconds())
				return
			}

			lastErr = err
			state.SetLastError(err)
			logger.Warn("model warmup attempt failed", "attempt", attempt, "max_attempts", policy.maxAttempts, "error", err, "latency_ms", time.Since(attemptStarted).Milliseconds())
			if attempt == policy.maxAttempts {
				break
			}
			if err := sleepWarmupBackoff(context.Background(), policy.backoff(attempt)); err != nil {
				lastErr = err
				state.SetLastError(err)
				break
			}
		}

		logger.Error("model warmup failed", "error", lastErr, "attempts", policy.maxAttempts, "latency_ms", time.Since(started).Milliseconds())
	}()
}

func sleepWarmupBackoff(ctx context.Context, d time.Duration) error {
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

func main() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(getEnv("LOG_LEVEL", "info")),
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	teiURL := getEnv("TEI_URL", "http://tei:80")
	redisHost := getEnv("REDIS_HOST", "redis:6379")
	modelID := getEnv("MODEL_ID", "deepvk/USER-bge-m3")
	bindHost := getEnv("BIND_HOST", "0.0.0.0")
	port := getEnv("PORT", "8000")
	redisPoolSize := envutil.Int("REDIS_POOL_SIZE", 50)

	runtimeConfig, err := loadEmbeddingRuntimeConfig()
	if err != nil {
		logger.Error("embedding runtime config invalid", "error", err)
		os.Exit(1)
	}
	logger.Info("embedding backend configured", "backend", runtimeConfig.Backend)

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	logger.Info("starting fd-api", "cpus", numCPU)

	// L1: Local cache — 10000 entries, 30s TTL
	localCache := cache.NewLocalCache(10000, 30*time.Second)

	// L2: Redis binary cache with pool timeouts
	redisOptions, err := cache.RedisCacheOptionsFromEnv("embed:cache:", redisPoolSize)
	if err != nil {
		logger.Error("redis cache config invalid", "error", err)
		os.Exit(1)
	}
	redisCache, err := cache.NewRedisCacheWithOptions(redisHost, redisOptions)
	if err != nil {
		logger.Error("redis cache init failed", "error", err)
		closeResource("local cache", localCache, logger)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := redisCache.Ping(ctx); err != nil {
		cancel()
		logger.Error("redis connect failed", "error", err)
		if closeErr := redisCache.Close(); closeErr != nil {
			logger.Warn("redis close failed after ping error", "error", closeErr)
		}
		closeResource("local cache", localCache, logger)
		os.Exit(1)
	}
	cancel()
	logger.Info("redis connected", "addr", redisHost, "cache_namespace", redisOptions.Namespace.String())

	// Two-tier cache
	tiered := cache.NewTieredCache(localCache, redisCache, 30*time.Second)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	lifecycleState := lifecycle.NewState()
	buildInfo := buildinfo.New(buildinfo.Info{
		Version:   Version,
		Model:     modelID,
		BuildHash: BuildHash,
		BuildDate: BuildDate,
	})
	maxInFlight := envutil.Int("FD_MAX_IN_FLIGHT", 0)
	if maxInFlight > 0 {
		logger.Info("embedding lifecycle capacity gate enabled", "max_in_flight", maxInFlight)
	}

	teiClient := embed.NewTEIClient(teiURL, modelID, httpClient)
	embeddingClient := embed.Embedder(teiClient)
	logger.Info("tei client configured", "url", teiURL, "model", modelID)

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	configureTrustedProxies(r)
	r.HandleMethodNotAllowed = true // explicit; gin's v1.12 default may differ
	// Wrap gin.Recovery so any panic returns an OpenAI-style error envelope
	// (500 internal_error) instead of gin.Recovery's default plain-text
	// 500. Without this, T-E-15 fails — panic-induced 500s would leak
	// server internals and lack the code/type envelope.
	metrics := observability.NewMetrics()
	traces := observability.NewTraceStoreFromEnv()
	r.Use(handlers.RecoveryMiddleware(logger))
	r.Use(middleware.CORSFromEnv())
	r.Use(middleware.HeadersMiddleware(buildInfo, modelID))
	r.Use(traces.Middleware(modelID))
	r.Use(metrics.Middleware())
	r.Use(middleware.APIKeyAuthFromEnv())
	r.Use(middleware.IPRateLimitFromEnv())
	r.Use(middleware.CacheHeaders())

	// 404/405 envelopes for paths/methods that don't match a registered
	// route. Without these, gin returns text/plain "404 page not found"
	// which fails the v2 spec (T-E-8, T-E-10).
	r.NoRoute(handlers.NotFoundMiddleware())
	r.NoMethod(handlers.MethodNotAllowedMiddleware())

	embedHandler := handlers.NewEmbeddingsHandler(embeddingClient, tiered, modelID, logger)
	batchHandler := handlers.NewBatchHandler(embeddingClient, tiered, modelID, logger)
	v1BatchHandler := handlers.NewV1BatchHandler(embeddingClient, tiered, logger)
	cacheHandler := handlers.NewCacheHandler(tiered)

	runtimeHealth := runtimeConfig.Health(modelID, redisOptions.Namespace.String())
	healthHandler := handlers.NewHealthHandlerWithState(runtimeHealth, lifecycleState)
	r.GET("/live", handlers.NewLiveHandler())
	r.GET("/ready", handlers.NewReadyHandler(lifecycleState))
	r.GET("/version", handlers.NewVersionHandler(buildInfo))
	r.GET("/info", handlers.NewInfoHandler(buildInfo, runtimeHealth, lifecycleState))
	r.GET("/openapi.json", handlers.NewOpenAPIHandler())
	r.GET("/docs", handlers.NewDocsHandler())
	warmupHandler := handlers.NewWarmupHandler(lifecycleState, embeddingClient, defaultWarmupTimeout)
	r.GET("/metrics", metrics.Handler())
	r.GET("/v1/traces", traces.Handler())
	r.GET("/warmup", warmupHandler.Status)
	r.POST("/warmup", warmupHandler.Trigger)
	r.POST("/v1/cache/flush", cacheHandler.Flush)
	r.POST("/v1/cache/delete", cacheHandler.Delete)
	r.GET("/health", healthHandler)
	r.GET("/v1/healthcheck", healthHandler)
	// /v1/embeddings: validation middleware runs BEFORE the handler so
	// 4xx/5xx (400 input_required, 413 input_too_long, 413 batch_too_large,
	// 413 payload_too_large) are returned without burning inference
	// capacity. The handler reads the parsed request from gin context.
	r.POST("/v1/embeddings", middleware.ValidateEmbeddingsRequest(), middleware.UserRateLimitFromEnv(), middleware.LifecycleGateWithCapacity(lifecycleState, int64(maxInFlight)), embedHandler.CreateEmbedding)
	// Batch endpoints use distinct JSON shapes, so they share the body/rate/lifecycle
	// guardrails while handlers perform shape-specific input validation before backend work.
	r.POST("/v1/batch", middleware.LimitRequestBody(), middleware.UserRateLimitFromEnv(), middleware.LifecycleGateWithCapacity(lifecycleState, int64(maxInFlight)), v1BatchHandler.CreateBatch)
	r.POST("/embeddings/batch", middleware.LimitRequestBody(), middleware.UserRateLimitFromEnv(), middleware.LifecycleGateWithCapacity(lifecycleState, int64(maxInFlight)), batchHandler.CreateBatchEmbeddings)

	addr := bindHost + ":" + port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
		// ReadHeaderTimeout: bound the time spent reading the request
		// headers to mitigate Slowloris-style attacks (gosec G112).
		// 10s is generous for /v1/embeddings callers (request bodies are
		// small; S01 caps at 10MB) and matches the Redis Ping timeout.
		ReadHeaderTimeout: 10 * time.Second,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go reportHTTPServerError(logger, addr, srv.ListenAndServe, sigCh)

	startModelWarmup(logger, lifecycleState, embeddingClient, defaultWarmupTimeout)

	if err := lifecycle.AwaitSignalAndShutdown(
		context.Background(),
		sigCh,
		srv,
		lifecycleState,
		logger,
		lifecycle.DefaultShutdownTimeout,
	); err != nil {
		logger.Error("shutdown failed", "error", err)
		closeResource("redis", redisCache, logger)
		closeResource("local cache", localCache, logger)
		os.Exit(1)
	}
	closeResource("redis", redisCache, logger)
	closeResource("local cache", localCache, logger)
}

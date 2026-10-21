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
	"fd-api/queue"

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

// HTTP server deadlines. A handler gets a 30s budget (validation, lifecycle gate,
// and the TEI call behind http.Client.Timeout = 30s), so the read and write
// deadlines leave margin for slow peers without letting one connection outlive
// the handler:
//
//	read  = handler 30s + 5s body margin  = 35s
//	write = handler 30s + 10s response margin = 40s
//	idle  = 75s keep-alive reuse
//
// ReadHeaderTimeout stays at 10s: request bodies are small (S01 caps at 10MB),
// and slow header writers must not hold a connection open (gosec G112).
const (
	defaultHTTPReadTimeout       = 35 * time.Second
	defaultHTTPWriteTimeout      = 40 * time.Second
	defaultHTTPIdleTimeout       = 75 * time.Second
	defaultHTTPReadHeaderTimeout = 10 * time.Second
)

// positiveDurationFromEnv returns a positive Go duration from the environment,
// or fallback when the value is unset, malformed, zero, or negative, so a bad
// override cannot disable a server deadline.
func positiveDurationFromEnv(key string, fallback time.Duration) time.Duration {
	if d := envutil.DurationOrDefault(key, fallback); d > 0 {
		return d
	}
	return fallback
}

// newHTTPServer builds the public HTTP server with bounded read, write, and
// idle deadlines.
func newHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: defaultHTTPReadHeaderTimeout,
		ReadTimeout:       positiveDurationFromEnv("FD_HTTP_READ_TIMEOUT", defaultHTTPReadTimeout),
		WriteTimeout:      positiveDurationFromEnv("FD_HTTP_WRITE_TIMEOUT", defaultHTTPWriteTimeout),
		IdleTimeout:       positiveDurationFromEnv("FD_HTTP_IDLE_TIMEOUT", defaultHTTPIdleTimeout),
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

// warmupRetryPolicyFromEnv builds the startup warmup retry policy from
// environment. This expands the startup window over the legacy
// defaultWarmupRetryPolicy so fd-api survives a slow TEI startup (CPU BERT
// load ~15-20s) without operator intervention. Backoff is a fixed interval
// rather than exponential: the legacy doubling backoff (2s,4s,8s) blew past
// the TEI load window and looked like a hang in logs.
//
// Env knobs (all safe to leave unset):
//
//	FD_WARMUP_START_MAX_ATTEMPTS  default 5   (>=1, clamped)
//	FD_WARMUP_START_BACKOFF_SEC   default 5   (>=0, 0 disables backoff)
func warmupRetryPolicyFromEnv() warmupRetryPolicy {
	maxAttempts := envutil.Int("FD_WARMUP_START_MAX_ATTEMPTS", 5)
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	backoffSeconds := envutil.Int("FD_WARMUP_START_BACKOFF_SEC", 5)
	backoff := time.Duration(backoffSeconds) * time.Second
	return warmupRetryPolicy{
		maxAttempts: maxAttempts,
		backoff:     func(int) time.Duration { return backoff },
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

// healthParams groups the inputs needed to build the deep /health handler.
type healthParams struct {
	runtimeConfig       *embeddingRuntimeConfig
	lifecycleState      *lifecycle.State
	modelID             string
	teiURL              string
	redisNamespace      string
	httpClient          *http.Client
	tiered              *cache.TieredCache
	maxInFlightCapacity int64
	logger              *slog.Logger
}

// newHealthHandler builds the /health handler and its runtime metadata block.
// Dependency probes are cheap reachability checks: they never run inference.
func newHealthHandler(p healthParams) (*handlers.RuntimeHealth, gin.HandlerFunc) {
	runtimeHealth := p.runtimeConfig.Health(p.modelID, p.redisNamespace)
	teiHealthURL := strings.TrimRight(p.teiURL, "/") + "/health"
	handler := handlers.NewHealthHandlerWithOptions(runtimeHealth, p.lifecycleState, handlers.HealthOptions{
		InFlightCapacity: &p.maxInFlightCapacity,
		Dependencies: &handlers.DependencyChecks{
			TEI: handlers.DependencyProbeFunc(func(ctx context.Context) handlers.DependencyStatus {
				started := time.Now()
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, teiHealthURL, http.NoBody)
				if err != nil {
					return handlers.DependencyStatus{Reachable: false, LatencyMS: elapsedMS(started), Error: err.Error()}
				}
				resp, err := p.httpClient.Do(req)
				latencyMS := elapsedMS(started)
				if err != nil {
					return handlers.DependencyStatus{Reachable: false, LatencyMS: latencyMS, Error: err.Error()}
				}
				defer func() {
					if closeErr := resp.Body.Close(); closeErr != nil {
						p.logger.Debug("tei health response close failed", "error", closeErr)
					}
				}()
				reachable := resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusInternalServerError
				status := handlers.DependencyStatus{Reachable: reachable, LatencyMS: latencyMS}
				if !reachable {
					status.Error = fmt.Sprintf("status %d", resp.StatusCode)
				}
				return status
			}),
			Redis: handlers.DependencyProbeFunc(func(ctx context.Context) handlers.DependencyStatus {
				started := time.Now()
				err := p.tiered.Ping(ctx)
				status := handlers.DependencyStatus{Reachable: err == nil, LatencyMS: elapsedMS(started), Namespace: p.redisNamespace}
				if err != nil {
					status.Error = err.Error()
				}
				return status
			}),
		},
		DependencyTimeout: 750 * time.Millisecond,
	})
	return runtimeHealth, handler
}

// elapsedMS reports the time since started in milliseconds with microsecond
// precision, used for dependency probe latency.
func elapsedMS(started time.Time) float64 {
	return float64(time.Since(started).Microseconds()) / 1000
}

func run() error {
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
		return err
	}
	logger.Info("embedding backend configured", "backend", runtimeConfig.Backend)

	logger.Info("starting fd-api", "gomaxprocs", runtime.GOMAXPROCS(0))

	// L1: Local cache — env-tunable size and TTL (M053 Phase 1a, Issue #9).
	// Defaults preserve legacy behavior (10000 entries, 30s TTL).
	l1MaxEntries := envutil.PositiveInt("FD_CACHE_MAX_SIZE", 10000)
	l1TTL := envutil.DurationOrDefault("FD_CACHE_LOCAL_TTL", 30*time.Second)
	localCache := cache.NewLocalCache(l1MaxEntries, l1TTL)
	defer closeResource("local cache", localCache, logger)

	// L2: Redis binary cache with pool timeouts
	redisOptions, err := cache.RedisCacheOptionsFromEnv("embed:cache:", redisPoolSize)
	if err != nil {
		logger.Error("redis cache config invalid", "error", err)
		return err
	}
	redisCache, err := cache.NewRedisCacheWithOptions(redisHost, redisOptions)
	if err != nil {
		logger.Error("redis cache init failed", "error", err)
		return err
	}
	defer closeResource("redis", redisCache, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := redisCache.Ping(ctx); err != nil {
		cancel()
		logger.Error("redis connect failed", "error", err)
		return err
	}
	cancel()
	logger.Info("redis connected", "addr", redisHost, "cache_namespace", redisOptions.Namespace.String())

	// Two-tier cache
	tiered := cache.NewTieredCache(localCache, redisCache, l1TTL)
	logger.Info("cache configured", "l1_max_entries", l1MaxEntries, "l1_ttl", l1TTL.String(), "l2_namespace", redisOptions.Namespace.String())

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
	maxInFlightCapacity := int64(maxInFlight)

	teiClient := embed.NewTEIClient(teiURL, modelID, httpClient)
	// M055-ucu1jl Phase 1b (Issue #9): opt-in cross-request coalescing.
	// When FD_COALESCE_ENABLED=true, wrap the TEI client in CoalescingEmbedder
	// so concurrent /v1/embeddings requests under burst share a single TEI
	// call within FD_COALESCE_WINDOW_MS. Default 0 = disabled (pass-through).
	embeddingClient := embed.Embedder(teiClient)
	if envutil.BoolOrDefault("FD_COALESCE_ENABLED", false) {
		coWindow := envutil.DurationOrDefault("FD_COALESCE_WINDOW_MS", 5*time.Millisecond)
		co := embed.NewCoalescingEmbedder(embeddingClient, coWindow)
		defer co.Close()
		embeddingClient = co
		logger.Info("embedding coalescing enabled", "window_ms", coWindow.Milliseconds())
	}
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
	// M052-mmf99p Phase 0: per-stage TEI observability hooks. Observers are
	// nil-tolerant inside TEIClient (nil fields = no-op), so we install
	// unconditionally here. Wiring after metrics declaration so the
	// Observe* methods are guaranteed available.
	teiClient.WithObservers(embed.Observers{
		ObserveDuration:  metrics.ObserveTEIRequestDuration,
		ObserveError:     metrics.IncTEIError,
		IncInFlight:      metrics.IncTEIRequestsInFlight,
		DecInFlight:      metrics.DecTEIRequestsInFlight,
		ObserveBatchFill: func(n int) { metrics.ObserveBatchFillRatio(float64(n) / 32.0) },
	})
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
	metrics.SetRuntimeObservers(lifecycleState, maxInFlightCapacity, tiered.LocalSize)
	// M052-mmf99p Phase 0: wire cache tier observer and L2 size gauge.
	// Captures per-tier hit-rate and rough L2 namespace occupancy for
	// the throughput optimization backplane (Issue #9).
	tiered.SetCacheObserver(func(tier string, hit bool) {
		result := "miss"
		if hit {
			result = "hit"
		}
		metrics.ObserveCacheResultWithTier(result, tier)
	})
	metrics.SetRedisCacheSizeObserver(func() int {
		return tiered.RedisSize(context.Background())
	})
	tiered.SetLookupDurationObserver(metrics.ObserveCacheLookupDuration)
	// L1 capacity evictions (FD_CACHE_MAX_SIZE overflow) feed
	// fd_cache_evictions_total; expired-entry cleanup is a separate event.
	localCache.SetEvictionObserver(metrics.ObserveCacheEvictions)

	runtimeHealth, healthHandler := newHealthHandler(healthParams{
		runtimeConfig:       runtimeConfig,
		lifecycleState:      lifecycleState,
		modelID:             modelID,
		teiURL:              teiURL,
		redisNamespace:      redisOptions.Namespace.String(),
		httpClient:          httpClient,
		tiered:              tiered,
		maxInFlightCapacity: maxInFlightCapacity,
		logger:              logger,
	})
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

	// M054-ybl4jr Phase 2: optional async bulk-ingestion queue. Backed by
	// an in-memory bounded channel drained by a worker goroutine. Disabled
	// by default for backward compatibility; enable via FD_QUEUE_ENABLED=true.
	queueEnabled := envutil.BoolOrDefault("FD_QUEUE_ENABLED", false)
	if queueEnabled {
		queueCap := envutil.PositiveInt("FD_QUEUE_MAX_SIZE", 1024)
		queueItems := make(chan queue.Item, queueCap)
		resultStore := queue.NewResultStore()
		queueHandler := handlers.NewQueueHandler(resultStore, queueItems, modelID, logger)
		// Same validation contract as /v1/embeddings, as the README states:
		// shape, batch size, dimensions and encoding_format are checked before
		// the item is enqueued, so a client error cannot surface as a 500 later.
		r.POST("/v1/queue", middleware.ValidateEmbeddingsRequest(), queueHandler.Submit)
		r.GET("/v1/queue/:id", queueHandler.Poll)
		logger.Info("queue enabled", "capacity", queueCap)
		// The worker shares the process lifetime: cancelling this context on the
		// way out lets it mark queued items failed and stop issuing TEI calls,
		// instead of being killed mid-request.
		queueCtx, queueCancel := context.WithCancel(context.Background())
		queueDone := queue.StartQueueWorker(queueCtx, resultStore, queueItems, embeddingClient, logger, queue.WorkerConfig{
			BatchMaxSize: envutil.PositiveInt("FD_QUEUE_BATCH_MAX_SIZE", 32),
			BatchWindow:  envutil.DurationOrDefault("FD_QUEUE_BATCH_WINDOW_MS", 10*time.Millisecond),
		})
		defer func() {
			queueCancel()
			select {
			case <-queueDone:
			case <-time.After(2 * time.Second):
				logger.Warn("queue worker did not stop within 2s")
			}
			closeResource("queue result store", resultStore, logger)
		}()
	} else {
		logger.Info("queue disabled (set FD_QUEUE_ENABLED=true to enable)")
	}
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
	srv := newHTTPServer(addr, r)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go reportHTTPServerError(logger, addr, srv.ListenAndServe, sigCh)

	warmupPolicy := warmupRetryPolicyFromEnv()
	logger.Info("model warmup policy",
		"max_attempts", warmupPolicy.maxAttempts,
		"backoff", warmupPolicy.backoff(1).String(),
	)
	startModelWarmupWithPolicy(logger, lifecycleState, embeddingClient, defaultWarmupTimeout, warmupPolicy)

	// M051-h1xr44 S01/T03: background warmup recovery. Survives the
	// startup-time race when fd_api starts before TEI finishes BERT load
	// (~15-20s on CPU) by retrying PreWarm periodically until IsWarmupDone.
	// Cancelled on signal so the goroutine exits deterministically.
	recoveryCtx, recoveryCancel := context.WithCancel(context.Background())
	defer recoveryCancel()
	recoveryInterval := time.Duration(envutil.Int("FD_WARMUP_RECOVERY_INTERVAL_SEC", 30)) * time.Second
	recoveryEnabled := envutil.BoolOrDefault("FD_WARMUP_RECOVERY_ENABLED", true)
	logger.Info("warmup recovery config",
		"enabled", recoveryEnabled,
		"interval", recoveryInterval.String(),
	)
	lifecycle.StartWarmupRecovery(recoveryCtx, logger, lifecycleState, embeddingClient, defaultWarmupTimeout, recoveryInterval, recoveryEnabled)

	if err := lifecycle.AwaitSignalAndShutdown(
		context.Background(),
		sigCh,
		srv,
		lifecycleState,
		logger,
		lifecycle.DefaultShutdownTimeout,
	); err != nil {
		logger.Error("shutdown failed", "error", err)
		return err
	}
	return nil
}

// main runs the service and exits non-zero on any fatal startup or shutdown
// error. Cleanup lives in run's defers, which run on the normal return path.
func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"fd-api/cache"
	"fd-api/embed"
	"fd-api/handlers"

	"github.com/gin-gonic/gin"
)

func getEnv(key, default_ string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return default_
}

func getEnvInt(key string, default_ int) int {
	value := os.Getenv(key)
	if value == "" {
		return default_
	}
	var n int
	for _, c := range []byte(value) {
		if c < '0' || c > '9' {
			return default_
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func sha256FileHex(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open %q for sha256: %w", path, err)
	}
	defer func() {
		_ = file.Close()
	}()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("hash %q: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func verifyFileSHA256(path, expected string) error {
	if expected == "" {
		return nil
	}
	actual, err := sha256FileHex(path)
	if err != nil {
		return err
	}
	if actual != expected {
		return fmt.Errorf("sha256 mismatch for %q: expected=%s actual=%s", path, expected, actual)
	}
	return nil
}

func verifyTokenizerJSON(path string, validation *embed.ONNXArtifactValidation) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat ONNX tokenizer JSON: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("ONNX_TOKENIZER_PATH must point to tokenizer JSON file, got directory")
	}
	if validation != nil && validation.TokenizerJSONSizeBytes > 0 && info.Size() != validation.TokenizerJSONSizeBytes {
		return fmt.Errorf("ONNX tokenizer JSON size mismatch: size=%d expected=%d", info.Size(), validation.TokenizerJSONSizeBytes)
	}
	if validation != nil && validation.TokenizerJSONSHA256 != "" {
		if err := verifyFileSHA256(path, validation.TokenizerJSONSHA256); err != nil {
			return fmt.Errorf("ONNX tokenizer JSON verification failed: %w", err)
		}
	}
	return nil
}

func boolPtr(v bool) *bool { return &v }

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

const (
	embeddingBackendTEI  embeddingBackend = "tei"
	embeddingBackendONNX embeddingBackend = "onnx"
	onnxProviderCPU      string           = "CPUExecutionProvider"
)

type embeddingRuntimeConfig struct {
	Backend                    embeddingBackend
	ONNXManifestPath           string
	ONNXRuntimeLibraryPath     string
	ONNXTokenizerPath          string
	ONNXMaxSequenceLength      int
	ONNXProvider               string
	ONNXArtifact               *embed.ONNXArtifactValidation
	ONNXTokenizerVerified      bool
	ONNXRuntimeLibraryVerified bool
}

// Health returns safe runtime metadata for the active backend.
// TEI/default returns backend, model, fixed dimensions, production-default flag,
// and cache namespace — no internal URLs, paths, or secrets.
// ONNX additionally returns artifact, tokenizer, and runtime library verification fields.
func (c *embeddingRuntimeConfig) Health(modelID, cacheNamespace string) *handlers.RuntimeHealth {
	if c == nil {
		return nil
	}
	switch c.Backend {
	case embeddingBackendTEI:
		return &handlers.RuntimeHealth{
			Backend:           string(c.Backend),
			Model:             modelID,
			Dimensions:        1024, // deepvk/USER-bge-m3 is 1024-dimensional
			ProductionDefault: true,
			CacheNamespace:    cacheNamespace,
		}
	case embeddingBackendONNX:
		if c.ONNXArtifact == nil {
			return nil
		}
		tv := c.ONNXTokenizerVerified
		rlv := c.ONNXRuntimeLibraryVerified
		return &handlers.RuntimeHealth{
			Backend:                    string(c.Backend),
			Model:                      modelID,
			ArtifactID:                 c.ONNXArtifact.ArtifactID,
			Dimensions:                 c.ONNXArtifact.Dimensions,
			MaxSequenceLength:          c.ONNXMaxSequenceLength,
			ValidatedMaxSequenceLength: c.ONNXArtifact.ValidatedMaxSequenceLength,
			ProductionDefault:          c.ONNXArtifact.ProductionDefault,
			ArtifactVerified:           boolPtr(true),
			TokenizerVerified:          &tv,
			RuntimeLibraryVerified:     &rlv,
			Provider:                   c.ONNXProvider,
			CacheNamespace:             cacheNamespace,
		}
	default:
		return nil
	}
}

func loadEmbeddingRuntimeConfig() (*embeddingRuntimeConfig, error) {
	backend := embeddingBackend(strings.ToLower(getEnv("EMBEDDING_BACKEND", string(embeddingBackendTEI))))
	config := &embeddingRuntimeConfig{Backend: backend}

	switch backend {
	case embeddingBackendTEI:
		return config, nil
	case embeddingBackendONNX:
		manifestPath := getEnv("ONNX_ARTIFACT_MANIFEST", "")
		if manifestPath == "" {
			return nil, fmt.Errorf("ONNX_ARTIFACT_MANIFEST is required when EMBEDDING_BACKEND=onnx")
		}
		runtimeLibraryPath := getEnv("ONNX_RUNTIME_LIBRARY", "")
		if runtimeLibraryPath == "" {
			return nil, fmt.Errorf("ONNX_RUNTIME_LIBRARY is required when EMBEDDING_BACKEND=onnx")
		}
		tokenizerPath := getEnv("ONNX_TOKENIZER_PATH", "")
		if tokenizerPath == "" {
			return nil, fmt.Errorf("ONNX_TOKENIZER_PATH is required when EMBEDDING_BACKEND=onnx")
		}
		validation, err := embed.ValidateONNXArtifactManifest(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("onnx artifact validation failed: %w", err)
		}
		provider := getEnv("ONNX_PROVIDER", onnxProviderCPU)
		if provider != onnxProviderCPU {
			return nil, fmt.Errorf("ONNX_PROVIDER=%q is not supported by the current Go ONNX runtime path; supported provider=%q", provider, onnxProviderCPU)
		}
		if err := verifyTokenizerJSON(tokenizerPath, validation); err != nil {
			return nil, err
		}
		runtimeLibraryVerified := false
		if expectedRuntimeSHA256 := getEnv("ONNX_RUNTIME_SHA256", ""); expectedRuntimeSHA256 != "" {
			if len(expectedRuntimeSHA256) != 64 {
				return nil, fmt.Errorf("ONNX_RUNTIME_SHA256 must be a 64-character sha256 hex digest")
			}
			if err := verifyFileSHA256(runtimeLibraryPath, expectedRuntimeSHA256); err != nil {
				return nil, fmt.Errorf("ONNX runtime library verification failed: %w", err)
			}
			runtimeLibraryVerified = true
		}
		config.ONNXManifestPath = manifestPath
		config.ONNXRuntimeLibraryPath = runtimeLibraryPath
		config.ONNXTokenizerPath = tokenizerPath
		config.ONNXProvider = provider
		config.ONNXTokenizerVerified = validation.TokenizerJSONSHA256 != "" || validation.TokenizerJSONSizeBytes > 0
		config.ONNXRuntimeLibraryVerified = runtimeLibraryVerified
		config.ONNXMaxSequenceLength = getEnvInt("ONNX_MAX_SEQUENCE_LENGTH", 512)
		if validation.ValidatedMaxSequenceLength > 0 && config.ONNXMaxSequenceLength > validation.ValidatedMaxSequenceLength {
			return nil, fmt.Errorf("ONNX_MAX_SEQUENCE_LENGTH=%d exceeds validated_max_sequence_length=%d artifact_id=%q manifest=%q", config.ONNXMaxSequenceLength, validation.ValidatedMaxSequenceLength, validation.ArtifactID, manifestPath)
		}
		config.ONNXArtifact = validation
		return config, nil
	default:
		return nil, fmt.Errorf("EMBEDDING_BACKEND must be tei or onnx, got %q", backend)
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
	redisPoolSize := getEnvInt("REDIS_POOL_SIZE", 50)

	runtimeConfig, err := loadEmbeddingRuntimeConfig()
	if err != nil {
		logger.Error("embedding runtime config invalid", "error", err)
		os.Exit(1)
	}
	logger.Info("embedding backend configured", "backend", runtimeConfig.Backend)
	if runtimeConfig.Backend == embeddingBackendONNX && runtimeConfig.ONNXArtifact != nil {
		logger.Info(
			"onnx artifact preflight verified",
			"artifact_id", runtimeConfig.ONNXArtifact.ArtifactID,
			"dimensions", runtimeConfig.ONNXArtifact.Dimensions,
			"max_sequence_length", runtimeConfig.ONNXMaxSequenceLength,
			"validated_max_sequence_length", runtimeConfig.ONNXArtifact.ValidatedMaxSequenceLength,
			"production_default", runtimeConfig.ONNXArtifact.ProductionDefault,
			"tokenizer_verified", runtimeConfig.ONNXTokenizerVerified,
			"runtime_library_verified", runtimeConfig.ONNXRuntimeLibraryVerified,
			"provider", runtimeConfig.ONNXProvider,
		)
	}

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
		os.Exit(1)
	}
	defer func() {
		if err := redisCache.Close(); err != nil {
			logger.Warn("redis close failed", "error", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := redisCache.Ping(ctx); err != nil {
		logger.Error("redis connect failed", "error", err)
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

	var embeddingClient handlers.Embedder
	if runtimeConfig.Backend == embeddingBackendONNX {
		onnxClient, err := embed.NewONNXEmbedder(embed.ONNXEmbedderOptions{
			ManifestPath:      runtimeConfig.ONNXManifestPath,
			SharedLibraryPath: runtimeConfig.ONNXRuntimeLibraryPath,
			TokenizerPath:     runtimeConfig.ONNXTokenizerPath,
			MaxSequenceLength: runtimeConfig.ONNXMaxSequenceLength,
		})
		if err != nil {
			logger.Error("onnx backend init failed", "error", err)
			os.Exit(1)
		}
		defer func() {
			if err := onnxClient.Close(); err != nil {
				logger.Warn("onnx close failed", "error", err)
			}
		}()
		embeddingClient = onnxClient
		logger.Info(
			"onnx backend ready",
			"artifact_id", runtimeConfig.ONNXArtifact.ArtifactID,
			"dimensions", runtimeConfig.ONNXArtifact.Dimensions,
			"max_sequence_length", runtimeConfig.ONNXMaxSequenceLength,
			"validated_max_sequence_length", runtimeConfig.ONNXArtifact.ValidatedMaxSequenceLength,
			"production_default", runtimeConfig.ONNXArtifact.ProductionDefault,
		)
	} else {
		teiClient := embed.NewTEIClient(teiURL, modelID, httpClient)
		embeddingClient = teiClient
		logger.Info("tei client configured", "url", teiURL, "model", modelID)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	embedHandler := handlers.NewEmbeddingsHandler(embeddingClient, tiered, modelID, logger)
	batchHandler := handlers.NewBatchHandler(embeddingClient, tiered, modelID, logger)

	r.GET("/health", handlers.NewHealthHandler(runtimeConfig.Health(modelID, redisOptions.Namespace.String())))
	r.POST("/v1/embeddings", embedHandler.CreateEmbedding)
	r.POST("/embeddings/batch", batchHandler.CreateBatchEmbeddings)

	addr := bindHost + ":" + port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		logger.Info("listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh

	logger.Info("shutting down...")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
	}
	logger.Info("stopped")
}

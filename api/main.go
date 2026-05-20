package main

import (
	"context"
	"fmt"
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
)

type embeddingRuntimeConfig struct {
	Backend                embeddingBackend
	ONNXManifestPath       string
	ONNXRuntimeLibraryPath string
	ONNXTokenizerPath      string
	ONNXMaxSequenceLength  int
	ONNXArtifact           *embed.ONNXArtifactValidation
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
		config.ONNXManifestPath = manifestPath
		config.ONNXRuntimeLibraryPath = runtimeLibraryPath
		config.ONNXTokenizerPath = tokenizerPath
		config.ONNXMaxSequenceLength = getEnvInt("ONNX_MAX_SEQUENCE_LENGTH", 512)
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
	logger.Info("redis connected", "addr", redisHost)

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

	r.GET("/health", handlers.HealthHandler)
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

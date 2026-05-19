package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
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
	var n int
	for _, c := range []byte(os.Getenv(key)) {
		if c < '0' || c > '9' {
			return default_
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func main() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	teiURL := getEnv("TEI_URL", "http://tei:80")
	redisHost := getEnv("REDIS_HOST", "redis:6379")
	modelID := getEnv("MODEL_ID", "deepvk/USER-bge-m3")
	bindHost := getEnv("BIND_HOST", "0.0.0.0")
	port := getEnv("PORT", "8000")
	redisPoolSize := getEnvInt("REDIS_POOL_SIZE", 50)

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	logger.Info("starting fd-api", "cpus", numCPU)

	// L1: Local cache — 10000 entries, 30s TTL
	localCache := cache.NewLocalCache(10000, 30*time.Second)

	// L2: Redis binary cache with pool timeouts
	redisCache := cache.NewRedisCache(redisHost, "embed:cache:", redisPoolSize)
	defer redisCache.Close()

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

	teiClient := embed.NewTEIClient(teiURL, modelID, httpClient)
	logger.Info("tei client configured", "url", teiURL, "model", modelID)

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	embedHandler := handlers.NewEmbeddingsHandler(teiClient, tiered, modelID, logger)
	batchHandler := handlers.NewBatchHandler(teiClient, tiered, modelID, logger)

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
	srv.Shutdown(ctx)
	logger.Info("stopped")
}

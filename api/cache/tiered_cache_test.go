package cache

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDiscardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestTieredCache_GetOrLoad_SeparatesDimensionsForSameText(t *testing.T) {
	ctx := context.Background()
	local := NewLocalCache(100, time.Minute)
	redis := NewRedisCache("127.0.0.1:0", "test:", 1)
	defer func() { assert.NoError(t, redis.Close()) }()
	tc := NewTieredCacheWithLogger(local, redis, time.Minute, newDiscardLogger())

	loads := 0
	loader := func(context.Context) ([]float32, error) {
		loads++
		emb := make([]float32, 1024)
		for i := range emb {
			emb[i] = float32(i)
		}
		return emb, nil
	}

	got512, err := tc.GetOrLoad(ctx, "same text", 512, loader)
	require.NoError(t, err, "512d GetOrLoad")
	assert.Len(t, got512, 512)

	got1024, err := tc.GetOrLoad(ctx, "same text", 1024, loader)
	require.NoError(t, err, "1024d GetOrLoad")
	assert.Len(t, got1024, 1024)
	assert.Equal(t, 2, loads, "loader calls for separate dimensions")
}

func TestTieredCache_GetOrLoad_ReturnsErrorForShortEmbedding(t *testing.T) {
	ctx := context.Background()
	local := NewLocalCache(100, time.Minute)
	redis := NewRedisCache("127.0.0.1:0", "test:", 1)
	defer func() { assert.NoError(t, redis.Close()) }()
	tc := NewTieredCacheWithLogger(local, redis, time.Minute, newDiscardLogger())

	_, err := tc.GetOrLoad(ctx, "short text", 512, func(context.Context) ([]float32, error) {
		return make([]float32, 128), nil
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "shorter than requested dimension")
}

func TestTieredCache_GetOrLoad_EmitsDebugCachePathWithoutRawKey(t *testing.T) {
	ctx := context.Background()
	local := NewLocalCache(100, time.Minute)
	redis := NewRedisCache("127.0.0.1:0", "test:", 1)
	defer func() { assert.NoError(t, redis.Close()) }()

	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))
	tc := NewTieredCacheWithLogger(local, redis, time.Minute, logger)

	const rawText = "sensitive benchmark text"
	loader := func(context.Context) ([]float32, error) {
		return make([]float32, 1024), nil
	}

	_, err := tc.GetOrLoad(ctx, rawText, 1024, loader)
	require.NoError(t, err, "first GetOrLoad")
	_, err = tc.GetOrLoad(ctx, rawText, 1024, loader)
	require.NoError(t, err, "second GetOrLoad")

	got := logs.String()
	for _, event := range []string{"cache_miss_load", "cache_l1_hit"} {
		assert.Contains(t, got, event)
	}
	assert.NotContains(t, got, rawText, "cache logs should not leak raw key text")
}

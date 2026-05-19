package cache

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestTieredCache_GetOrLoad_SeparatesDimensionsForSameText(t *testing.T) {
	ctx := context.Background()
	local := NewLocalCache(100, time.Minute)
	redis := NewRedisCache("127.0.0.1:0", "test:", 1)
	defer redis.Close()
	tc := NewTieredCache(local, redis, time.Minute)

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
	if err != nil {
		t.Fatalf("512d GetOrLoad error: %v", err)
	}
	if len(got512) != 512 {
		t.Fatalf("512d len=%d, want 512", len(got512))
	}

	got1024, err := tc.GetOrLoad(ctx, "same text", 1024, loader)
	if err != nil {
		t.Fatalf("1024d GetOrLoad error: %v", err)
	}
	if len(got1024) != 1024 {
		t.Fatalf("1024d len=%d, want 1024", len(got1024))
	}
	if loads != 2 {
		t.Fatalf("loader calls=%d, want 2 for separate dimensions", loads)
	}
}

func TestTieredCache_GetOrLoad_ReturnsErrorForShortEmbedding(t *testing.T) {
	ctx := context.Background()
	local := NewLocalCache(100, time.Minute)
	redis := NewRedisCache("127.0.0.1:0", "test:", 1)
	defer redis.Close()
	tc := NewTieredCache(local, redis, time.Minute)

	_, err := tc.GetOrLoad(ctx, "short text", 512, func(context.Context) ([]float32, error) {
		return make([]float32, 128), nil
	})
	if err == nil {
		t.Fatal("expected short embedding error")
	}
	if !strings.Contains(err.Error(), "shorter than requested dimension") {
		t.Fatalf("unexpected error: %v", err)
	}
}

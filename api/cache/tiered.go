package cache

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

// TieredCache provides two-tier cache-aside: L1 (local, []byte) → L2 (Redis, []float32).
// singleflight deduplicates concurrent requests for the same key.
type TieredCache struct {
	local    *LocalCache
	redis    *RedisCache
	localTTL time.Duration
	sf       singleflight.Group
}

// NewTieredCache creates a two-tier cache.
func NewTieredCache(local *LocalCache, redis *RedisCache, localTTL time.Duration) *TieredCache {
	return &TieredCache{
		local:    local,
		redis:    redis,
		localTTL: localTTL,
	}
}

// GetOrLoad checks L1 then L2, invoking loader if both miss.
func (tc *TieredCache) GetOrLoad(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error) {
	// L1 check — returns []byte
	if data, ok := tc.local.Get(ctx, key); ok {
		emb, d := unmarshalEmbedding(data)
		if d == dim {
			return emb, nil
		}
		// dim mismatch — treat as miss
	}

	// singleflight — dedup concurrent requests
	r, err, _ := tc.sf.Do(key, func() (any, error) {
		// Double-check L1
		if data, ok := tc.local.Get(ctx, key); ok {
			emb, d := unmarshalEmbedding(data)
			if d == dim {
				return emb, nil
			}
		}

		// L2 check — Redis returns []float32
		if emb, ok, err := tc.redis.Get(ctx, key, dim); err == nil && ok {
			// backfill L1 with binary
			data := marshalEmbedding(emb, dim)
			tc.local.Set(ctx, key, data, tc.localTTL)
			return emb, nil
		}

		// Both miss — invoke loader
		emb, err := loader(ctx)
		if err != nil {
			return nil, err
		}

		// Backfill L1 with binary
		data := marshalEmbedding(emb, dim)
		tc.local.Set(ctx, key, data, tc.localTTL)

		// Backfill L2
		tc.redis.Set(ctx, key, emb, dim)

		return emb, nil
	})
	if err != nil {
		return nil, err
	}
	return r.([]float32), nil
}

// Ping checks L2 (Redis) connectivity.
func (tc *TieredCache) Ping(ctx context.Context) error {
	return tc.redis.Ping(ctx)
}

// Close closes the Redis connection.
func (tc *TieredCache) Close() error {
	return tc.redis.Close()
}

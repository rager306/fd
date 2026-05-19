package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/singleflight"
)

// TieredCache provides two-tier cache-aside: L1 (local, []byte) → L2 (Redis, []float32).
// singleflight deduplicates concurrent requests for the same key.
type TieredCache struct {
	local    *LocalCache
	redis    *RedisCache
	localTTL time.Duration
	logger   *slog.Logger
	sf       singleflight.Group
}

// NewTieredCache creates a two-tier cache.
func NewTieredCache(local *LocalCache, redis *RedisCache, localTTL time.Duration) *TieredCache {
	return NewTieredCacheWithLogger(local, redis, localTTL, slog.Default().With("component", "tiered_cache"))
}

// NewTieredCacheWithLogger creates a two-tier cache with an explicit logger.
func NewTieredCacheWithLogger(local *LocalCache, redis *RedisCache, localTTL time.Duration, logger *slog.Logger) *TieredCache {
	if logger == nil {
		logger = slog.Default().With("component", "tiered_cache")
	}
	return &TieredCache{
		local:    local,
		redis:    redis,
		localTTL: localTTL,
		logger:   logger,
	}
}

// GetOrLoad checks L1 then L2, invoking loader if both miss.
func (tc *TieredCache) GetOrLoad(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error) {
	localKey := fmt.Sprintf("%s:d%d", key, dim)
	keyHash := shortCacheKeyHash(key)

	// L1 check — returns []byte
	if data, ok := tc.local.Get(ctx, localKey); ok {
		emb, d := unmarshalEmbedding(data)
		if d == dim {
			tc.logger.Debug("cache l1 hit", "event", "cache_l1_hit", "key_hash", keyHash, "dim", dim)
			return emb, nil
		}
		tc.logger.Debug("cache l1 dimension mismatch", "event", "cache_l1_dim_mismatch", "key_hash", keyHash, "requested_dim", dim, "stored_dim", d)
		// dim mismatch — treat as miss
	}

	// singleflight — dedup concurrent requests for the same text and dimension.
	r, err, shared := tc.sf.Do(localKey, func() (any, error) {
		// Double-check L1
		if data, ok := tc.local.Get(ctx, localKey); ok {
			emb, d := unmarshalEmbedding(data)
			if d == dim {
				tc.logger.Debug("cache l1 hit", "event", "cache_l1_hit", "key_hash", keyHash, "dim", dim, "source", "singleflight_double_check")
				return emb, nil
			}
			tc.logger.Debug("cache l1 dimension mismatch", "event", "cache_l1_dim_mismatch", "key_hash", keyHash, "requested_dim", dim, "stored_dim", d, "source", "singleflight_double_check")
		}

		// L2 check — Redis returns []float32
		emb, ok, err := tc.redis.Get(ctx, key, dim)
		if err != nil {
			tc.logger.Warn("cache l2 get failed", "event", "cache_l2_get_failed", "key_hash", keyHash, "dim", dim, "error", err)
		} else if ok {
			tc.logger.Debug("cache l2 hit", "event", "cache_l2_hit", "key_hash", keyHash, "dim", dim)
			// backfill L1 with binary
			data, err := marshalEmbedding(emb, dim)
			if err != nil {
				return nil, err
			}
			tc.local.Set(ctx, localKey, data, tc.localTTL)
			return emb, nil
		}

		// Both miss — invoke loader
		tc.logger.Debug("cache miss load", "event", "cache_miss_load", "key_hash", keyHash, "dim", dim)
		emb, err = loader(ctx)
		if err != nil {
			return nil, err
		}

		// Backfill L1 with binary
		data, err := marshalEmbedding(emb, dim)
		if err != nil {
			return nil, err
		}
		tc.local.Set(ctx, localKey, data, tc.localTTL)

		// Backfill L2
		if err := tc.redis.SetBytes(ctx, key, data, dim); err != nil {
			tc.logger.Warn("cache l2 set failed", "event", "cache_l2_set_failed", "key_hash", keyHash, "dim", dim, "error", err)
		}

		return emb[:dim], nil
	})
	if err != nil {
		return nil, err
	}
	if shared {
		tc.logger.Debug("cache singleflight shared", "event", "cache_singleflight_shared", "key_hash", keyHash, "dim", dim)
	}
	return r.([]float32), nil
}

func shortCacheKeyHash(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])[:12]
}

// Ping checks L2 (Redis) connectivity.
func (tc *TieredCache) Ping(ctx context.Context) error {
	return tc.redis.Ping(ctx)
}

// Close closes the Redis connection.
func (tc *TieredCache) Close() error {
	return tc.redis.Close()
}

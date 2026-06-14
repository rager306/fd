package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	localKey := localCacheKey(key, dim)
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

// GetIfPresent returns the cached embedding for (key, dim) without
// triggering any model load on miss. Used by the embeddings handler to
// peek at the cache before issuing TEI calls so a fully-cached batch
// results in zero TEI traffic. Returns (vec, true) on hit, (nil, false)
// on miss.
func (tc *TieredCache) GetIfPresent(ctx context.Context, key string, dim int) ([]float32, bool) {
	localKey := localCacheKey(key, dim)
	data, ok := tc.local.Get(ctx, localKey)
	if !ok {
		// L2 fallback (without backfilling L1)
		emb, lok, err := tc.redis.Get(ctx, key, dim)
		if err != nil || !lok {
			return nil, false
		}
		return emb[:dim], true
	}
	emb, d := unmarshalEmbedding(data)
	if d != dim {
		return nil, false
	}
	return emb[:dim], true
}

// Set stores an embedding in both L1 and L2 caches. Used by the
// embeddings handler to backfill the cache after a model call so the
// next request can use GetIfPresent.
func (tc *TieredCache) Set(ctx context.Context, key string, dim int, emb []float32) {
	if len(emb) < dim {
		return
	}
	localKey := localCacheKey(key, dim)
	data, err := marshalEmbedding(emb[:dim], dim)
	if err != nil {
		return
	}
	tc.local.Set(ctx, localKey, data, tc.localTTL)
	_ = tc.redis.SetBytes(ctx, key, data, dim)
}

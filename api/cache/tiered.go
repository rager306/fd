package cache

import (
	"context"
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
	keyHash := shortHash(key)

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

// Delete removes the cached embedding for (key, dim) from both cache tiers.
func (tc *TieredCache) Delete(ctx context.Context, key string, dim int) error {
	if tc.local != nil {
		tc.local.Delete(ctx, localCacheKey(key, dim))
	}
	if tc.redis == nil {
		return nil
	}
	return tc.redis.Delete(ctx, key, dim)
}

// Flush removes all local entries and all Redis entries in this cache namespace.
func (tc *TieredCache) Flush(ctx context.Context) (int64, error) {
	if tc.local != nil {
		tc.local.Flush(ctx)
	}
	if tc.redis == nil {
		return 0, nil
	}
	return tc.redis.FlushNamespace(ctx)
}

// LocalSize returns the number of current L1 entries.
func (tc *TieredCache) LocalSize() int {
	if tc.local == nil {
		return 0
	}
	return tc.local.Size()
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
	hits := tc.GetManyIfPresent(ctx, []string{key}, dim)
	emb, ok := hits[0]
	return emb, ok
}

// GetManyIfPresent returns cached embeddings for keys without triggering a
// model load. L1 is checked first for each key; L2 misses are fetched with a
// single Redis MGET and backfilled into L1. Hits are keyed by input index so
// callers can preserve duplicate text positions and response order.
func (tc *TieredCache) GetManyIfPresent(ctx context.Context, keys []string, dim int) map[int][]float32 {
	hits := make(map[int][]float32, len(keys))
	missIndexes := make([]int, 0, len(keys))
	missKeys := make([]string, 0, len(keys))

	for i, key := range keys {
		localKey := localCacheKey(key, dim)
		data, ok := tc.local.Get(ctx, localKey)
		if !ok {
			missIndexes = append(missIndexes, i)
			missKeys = append(missKeys, key)
			continue
		}
		emb, d := unmarshalEmbedding(data)
		if d != dim {
			missIndexes = append(missIndexes, i)
			missKeys = append(missKeys, key)
			continue
		}
		hits[i] = emb[:dim]
	}
	if len(missKeys) == 0 || tc.redis == nil {
		return hits
	}

	redisHits, err := tc.redis.GetMany(ctx, missKeys, dim)
	if err != nil {
		tc.logger.Warn("cache l2 mget failed", "event", "cache_l2_mget_failed", "dim", dim, "miss_count", len(missKeys), "error", err)
		return hits
	}
	for missOffset, emb := range redisHits {
		if len(emb) < dim {
			continue
		}
		originalIndex := missIndexes[missOffset]
		key := keys[originalIndex]
		localKey := localCacheKey(key, dim)
		data, err := marshalEmbedding(emb[:dim], dim)
		if err == nil {
			tc.local.Set(ctx, localKey, data, tc.localTTL)
		}
		hits[originalIndex] = emb[:dim]
	}
	return hits
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

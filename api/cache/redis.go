package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type cachedValue struct {
	Embedding []float32 `json:"embedding"`
	Dim      int        `json:"dim"`
}

type RedisCache struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

func NewRedisCache(addr, prefix string, poolSize int) *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			PoolSize: poolSize,
		}),
		prefix: prefix,
		ttl:    24 * time.Hour,
	}
}

func (c *RedisCache) HashText(text string) string {
	h := sha256.Sum256([]byte(text))
	return hex.EncodeToString(h[:])
}

func (c *RedisCache) Get(ctx context.Context, text string, dim int) ([]float32, bool, error) {
	key := c.prefix + "text:" + c.HashText(text)

	data, err := c.client.HGet(ctx, key, "embedding").Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var cv cachedValue
	if err := json.Unmarshal([]byte(data), &cv); err != nil {
		// Legacy data without dim — re-cache
		var embedding []float32
		if err2 := json.Unmarshal([]byte(data), &embedding); err2 == nil {
			// Re-store with dim
			c.Set(ctx, text, embedding, dim)
			return embedding[:dim], true, nil
		}
		return nil, false, err
	}

	// If stored dim != requested dim, treat as miss (will re-cache at correct dim)
	if cv.Dim != dim {
		return nil, false, nil
	}

	return cv.Embedding, true, nil
}

func (c *RedisCache) Set(ctx context.Context, text string, embedding []float32, dim int) error {
	key := c.prefix + "text:" + c.HashText(text)

	cv := cachedValue{Embedding: embedding, Dim: dim}
	embJSON, err := json.Marshal(cv)
	if err != nil {
		return fmt.Errorf("marshal embedding: %w", err)
	}

	pipe := c.client.Pipeline()
	pipe.HSet(ctx, key, "embedding", string(embJSON), "text", text)
	pipe.Expire(ctx, key, c.ttl)
	_, err = pipe.Exec(ctx)

	return err
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

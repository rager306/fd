package cache

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

func NewRedisCache(addr, prefix string, poolSize int) *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:         addr,
			PoolSize:     poolSize,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolTimeout:  4 * time.Second,
			MinIdleConns: 10,
		}),
		prefix: prefix,
		ttl:    24 * time.Hour,
	}
}

// marshalEmbedding encodes [dim:uint16][float32*dim] — 2+4*dim bytes.
// Replaces JSON (~8KB for 1024d → 4098 bytes).
func marshalEmbedding(embedding []float32, dim int) []byte {
	buf := make([]byte, 2+dim*4)
	binary.LittleEndian.PutUint16(buf[0:2], uint16(dim))
	for i := 0; i < dim; i++ {
		binary.LittleEndian.PutUint32(buf[2+i*4:2+(i+1)*4], math.Float32bits(embedding[i]))
	}
	return buf
}

// unmarshalEmbedding decodes binary format back to []float32.
func unmarshalEmbedding(data []byte) ([]float32, int) {
	if len(data) < 2 {
		return nil, 0
	}
	dim := int(binary.LittleEndian.Uint16(data[0:2]))
	if len(data) < 2+dim*4 {
		return nil, 0
	}
	emb := make([]float32, dim)
	for i := 0; i < dim; i++ {
		emb[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[2+i*4 : 2+(i+1)*4]))
	}
	return emb, dim
}

func (c *RedisCache) HashText(text string) string {
	h := sha256.Sum256([]byte(text))
	return hex.EncodeToString(h[:])
}

func (c *RedisCache) Get(ctx context.Context, text string, dim int) ([]float32, bool, error) {
	// v2 key format includes dimension suffix
	key := c.prefix + "v2:" + c.HashText(text) + ":d" + fmt.Sprintf("%d", dim)

	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	emb, gotDim := unmarshalEmbedding(data)
	if gotDim != dim {
		return nil, false, nil
	}
	return emb, true, nil
}

func (c *RedisCache) Set(ctx context.Context, text string, embedding []float32, dim int) error {
	key := c.prefix + "v2:" + c.HashText(text) + ":d" + fmt.Sprintf("%d", dim)
	data := marshalEmbedding(embedding, dim)
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

package cache

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultRedisCacheTTL = 24 * time.Hour

// RedisCacheNamespace controls correctness-affecting Redis key namespace fields.
// Optional values are hashed before they are included in keys so model/tokenizer
// identifiers do not make keys overly long or leak raw names through Redis scans.
type RedisCacheNamespace struct {
	CacheVersion     string
	ModelID          string
	ModelRevision    string
	TokenizerVersion string
	ChunkingVersion  string
}

// RedisCacheOptions controls Redis client, key namespace, and retention behavior.
type RedisCacheOptions struct {
	Prefix    string
	PoolSize  int
	TTL       time.Duration
	NoExpire  bool
	Namespace RedisCacheNamespace
}

type RedisCache struct {
	client    *redis.Client
	prefix    string
	ttl       time.Duration
	noExpire  bool
	namespace string
}

func DefaultRedisCacheOptions(prefix string, poolSize int) RedisCacheOptions {
	return RedisCacheOptions{
		Prefix:   prefix,
		PoolSize: poolSize,
		TTL:      defaultRedisCacheTTL,
		Namespace: RedisCacheNamespace{
			CacheVersion: "v2",
		},
	}
}

func RedisCacheOptionsFromEnv(prefix string, poolSize int) (RedisCacheOptions, error) {
	opts := DefaultRedisCacheOptions(prefix, poolSize)

	if v := strings.TrimSpace(os.Getenv("EMBEDDING_CACHE_VERSION")); v != "" {
		opts.Namespace.CacheVersion = v
	}
	opts.Namespace.ModelID = strings.TrimSpace(os.Getenv("EMBEDDING_MODEL_ID"))
	opts.Namespace.ModelRevision = strings.TrimSpace(os.Getenv("EMBEDDING_MODEL_REVISION"))
	opts.Namespace.TokenizerVersion = strings.TrimSpace(os.Getenv("EMBEDDING_TOKENIZER_VERSION"))
	opts.Namespace.ChunkingVersion = strings.TrimSpace(os.Getenv("EMBEDDING_CHUNKING_VERSION"))

	if v := strings.TrimSpace(os.Getenv("REDIS_CACHE_NO_EXPIRE")); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return RedisCacheOptions{}, fmt.Errorf("invalid REDIS_CACHE_NO_EXPIRE %q: %w", v, err)
		}
		opts.NoExpire = parsed
	}

	if v := strings.TrimSpace(os.Getenv("REDIS_CACHE_TTL")); v != "" {
		if opts.NoExpire {
			return RedisCacheOptions{}, fmt.Errorf("REDIS_CACHE_TTL and REDIS_CACHE_NO_EXPIRE=true are mutually exclusive")
		}
		parsed, err := time.ParseDuration(v)
		if err != nil {
			return RedisCacheOptions{}, fmt.Errorf("invalid REDIS_CACHE_TTL %q: %w", v, err)
		}
		if parsed <= 0 {
			return RedisCacheOptions{}, fmt.Errorf("REDIS_CACHE_TTL must be positive, got %s", v)
		}
		opts.TTL = parsed
	}

	return opts, nil
}

func NewRedisCache(addr, prefix string, poolSize int) *RedisCache {
	cache, err := NewRedisCacheWithOptions(addr, DefaultRedisCacheOptions(prefix, poolSize))
	if err != nil {
		panic(err)
	}
	return cache
}

func NewRedisCacheWithOptions(addr string, opts RedisCacheOptions) (*RedisCache, error) {
	if opts.Prefix == "" {
		return nil, fmt.Errorf("redis cache prefix must not be empty")
	}
	if opts.PoolSize <= 0 {
		return nil, fmt.Errorf("redis pool size must be positive, got %d", opts.PoolSize)
	}
	if !opts.NoExpire && opts.TTL <= 0 {
		return nil, fmt.Errorf("redis cache TTL must be positive unless no-expire mode is enabled")
	}

	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:         addr,
			PoolSize:     opts.PoolSize,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolTimeout:  4 * time.Second,
			MinIdleConns: 10,
		}),
		prefix:    opts.Prefix,
		ttl:       opts.TTL,
		noExpire:  opts.NoExpire,
		namespace: opts.Namespace.String(),
	}, nil
}

func (n RedisCacheNamespace) String() string {
	cacheVersion := strings.TrimSpace(n.CacheVersion)
	if cacheVersion == "" {
		cacheVersion = "v2"
	}
	parts := []string{cacheVersion}
	if n.ModelID != "" {
		parts = append(parts, "m"+shortNamespaceHash(n.ModelID))
	}
	if n.ModelRevision != "" {
		parts = append(parts, "r"+shortNamespaceHash(n.ModelRevision))
	}
	if n.TokenizerVersion != "" {
		parts = append(parts, "t"+shortNamespaceHash(n.TokenizerVersion))
	}
	if n.ChunkingVersion != "" {
		parts = append(parts, "c"+shortNamespaceHash(n.ChunkingVersion))
	}
	return strings.Join(parts, ":")
}

func shortNamespaceHash(value string) string {
	h := sha256.Sum256([]byte(value))
	return hex.EncodeToString(h[:])[:12]
}

func (c *RedisCache) expiration() time.Duration {
	if c.noExpire {
		return 0
	}
	return c.ttl
}

func (c *RedisCache) key(text string, dim int) string {
	return c.prefix + c.namespace + ":" + c.HashText(text) + ":d" + fmt.Sprintf("%d", dim)
}

// marshalEmbedding encodes [dim:uint16][float32*dim] — 2+4*dim bytes.
// Replaces JSON (~8KB for 1024d → 4098 bytes).
func marshalEmbedding(embedding []float32, dim int) ([]byte, error) {
	if dim <= 0 {
		return nil, fmt.Errorf("dimension must be positive, got %d", dim)
	}
	if len(embedding) < dim {
		return nil, fmt.Errorf("embedding length %d shorter than requested dimension %d", len(embedding), dim)
	}

	buf := make([]byte, 2+dim*4)
	binary.LittleEndian.PutUint16(buf[0:2], uint16(dim))
	for i := 0; i < dim; i++ {
		binary.LittleEndian.PutUint32(buf[2+i*4:2+(i+1)*4], math.Float32bits(embedding[i]))
	}
	return buf, nil
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
	data, err := c.client.Get(ctx, c.key(text, dim)).Bytes()
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
	data, err := marshalEmbedding(embedding, dim)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.key(text, dim), data, c.expiration()).Err()
}

// SetBytes stores pre-marshaled binary embedding.
func (c *RedisCache) SetBytes(ctx context.Context, text string, data []byte, dim int) error {
	return c.client.Set(ctx, c.key(text, dim), data, c.expiration()).Err()
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

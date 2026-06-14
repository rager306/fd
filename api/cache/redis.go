package cache

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
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

// RedisCache is the L2 cache backed by a Redis cluster. It speaks the
// OpenAI-compatible embedding-cache protocol (sha256-keyed binary blobs
// with a dim prefix) so that L1 (LocalCache) and L2 (Redis) can
// exchange values.
type RedisCache struct {
	client    *redis.Client
	prefix    string
	ttl       time.Duration
	noExpire  bool
	namespace string
}

// DefaultRedisCacheOptions returns sane defaults for a production fd
// deployment: 24h TTL, no no-expire flag, prefix=prefix, poolSize=poolSize.
// Caller can override individual fields before passing to NewRedisCacheWithOptions.
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

// RedisCacheOptionsFromEnv reads EMBEDDING_CACHE_VERSION, EMBEDDING_MODEL_ID,
// EMBEDDING_MODEL_REVISION, EMBEDDING_TOKENIZER_VERSION, EMBEDDING_CHUNKING_VERSION,
// REDIS_CACHE_NO_EXPIRE, REDIS_CACHE_TTL from the process environment and
// overlays them onto DefaultRedisCacheOptions. Returns an error if
// REDIS_CACHE_NO_EXPIRE or REDIS_CACHE_TTL are set but unparseable.
//
// The namespace fields are critical for isolating TEI vs ONNX cache keys
// (per M040/M041 gotcha: cross-backend cache pollution silently produces
// false-positive equivalence).
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

// NewRedisCache connects to Redis at addr with default options. Panics
// if the connection fails — startup is a programmer error, not a runtime
// degradation. For env-driven configuration use NewRedisCacheWithOptions
// + RedisCacheOptionsFromEnv.
func NewRedisCache(addr, prefix string, poolSize int) *RedisCache {
	cache, err := NewRedisCacheWithOptions(addr, DefaultRedisCacheOptions(prefix, poolSize))
	if err != nil {
		panic(err)
	}
	return cache
}

// NewRedisCacheWithOptions connects to Redis at addr using the supplied
// opts. Returns an error if Prefix is empty, PoolSize is non-positive,
// or TTL is non-positive when no-expire is disabled. Caller is responsible
// for calling Close to release the connection pool.
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
	return c.prefix + c.namespace + ":" + c.HashText(text) + ":d" + strconv.Itoa(dim)
}

// marshalEmbedding encodes [dim:uint16][float32*dim] — 2+4*dim bytes.
// Replaces JSON (~8KB for 1024d → 4098 bytes).
func marshalEmbedding(embedding []float32, dim int) ([]byte, error) {
	if dim <= 0 {
		return nil, fmt.Errorf("dimension must be positive, got %d", dim)
	}
	if dim > maxUint16 {
		// Stored as uint16 in the cache blob (2-byte dim prefix + dim*4 bytes
		// of float32 LE). fd's model only supports 512/1024 today, but the
		// limit exists so future larger models fail loudly instead of
		// silently truncating.
		return nil, fmt.Errorf("dimension %d exceeds cache uint16 capacity %d", dim, maxUint16)
	}
	if len(embedding) < dim {
		return nil, fmt.Errorf("embedding length %d shorter than requested dimension %d", len(embedding), dim)
	}

	buf := make([]byte, 2+dim*4)
	binary.LittleEndian.PutUint16(buf[0:2], uint16(dim)) //nolint:gosec // G115: bounds-checked above
	for i := 0; i < dim; i++ {
		binary.LittleEndian.PutUint32(buf[2+i*4:2+(i+1)*4], math.Float32bits(embedding[i]))
	}
	return buf, nil
}

// maxUint16 is the upper bound for the cached dim prefix (uint16). Set as
// a constant so the bound is explicit at the call site and the gosec
// G115 suppression carries meaning.
const maxUint16 = 0xFFFF

// unmarshalEmbedding decodes binary format back to []float32.
func unmarshalEmbedding(data []byte) (embedding []float32, dim int) {
	if len(data) < 2 {
		return nil, 0
	}
	dim = int(binary.LittleEndian.Uint16(data[0:2]))
	if len(data) < 2+dim*4 {
		return nil, 0
	}
	emb := make([]float32, dim)
	for i := 0; i < dim; i++ {
		emb[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[2+i*4 : 2+(i+1)*4]))
	}
	return emb, dim
}

// HashText returns the sha256-hex digest of text. Used as the cache key
// component for the embedding text (the dim and prefix are added by
// the key() method to form the full Redis key).
func (c *RedisCache) HashText(text string) string {
	h := sha256.Sum256([]byte(text))
	return hex.EncodeToString(h[:])
}

// Get retrieves the cached embedding vector for (text, dim). Returns
// (nil, false, nil) on miss, the error from Redis on transport failure.
// The dim is appended to the cache key; a dim mismatch (e.g., the
// value was stored for a different dim) is treated as a miss.
func (c *RedisCache) Get(ctx context.Context, text string, dim int) (embedding []float32, found bool, err error) {
	data, err := c.client.Get(ctx, c.key(text, dim)).Bytes()
	if errors.Is(err, redis.Nil) {
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

// Set stores the embedding vector for (text, dim). If no-expire is
// configured the value is written without TTL; otherwise the configured
// TTL is applied. The embedding is encoded via marshalEmbedding before
// transport.
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

// Ping checks Redis liveness. Used by fd's startup preflight (after the
// constructor) and by TieredCache.Ping.
func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close releases the Redis connection pool. Safe to call multiple times.
func (c *RedisCache) Close() error {
	return c.client.Close()
}

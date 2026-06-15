package cache

import (
	"container/list"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	defaultLRUCacheSize     = 10000
	defaultLRUCacheTTLHours = 24
)

// LRUCacheMetrics is the metrics surface used by LRUCache.
type LRUCacheMetrics interface {
	ObserveCacheResult(result string)
	ObserveCacheEviction()
}

type lruCacheEntry struct {
	key       string
	value     []float32
	expiresAt time.Time
}

// LRUCache is a goroutine-safe in-memory LRU cache for embedding vectors.
type LRUCache struct {
	mu      sync.Mutex
	entries map[string]*list.Element
	order   *list.List
	maxSize int
	ttl     time.Duration
	metrics LRUCacheMetrics
}

// NewLRUCache creates a TTL-aware in-memory LRU cache. maxSize <= 0 disables
// size eviction; ttl <= 0 disables time expiration.
func NewLRUCache(maxSize int, ttl time.Duration, metrics LRUCacheMetrics) *LRUCache {
	return &LRUCache{
		entries: make(map[string]*list.Element),
		order:   list.New(),
		maxSize: maxSize,
		ttl:     ttl,
		metrics: metrics,
	}
}

// NewLRUCacheFromEnv creates an LRU cache configured by FD_CACHE_SIZE and
// FD_CACHE_TTL_HOURS, defaulting to 10000 entries and 24 hours.
func NewLRUCacheFromEnv(metrics LRUCacheMetrics) *LRUCache {
	return NewLRUCache(envInt("FD_CACHE_SIZE", defaultLRUCacheSize), time.Duration(envInt("FD_CACHE_TTL_HOURS", defaultLRUCacheTTLHours))*time.Hour, metrics)
}

// EmbeddingCacheKey returns the SHA256 key for input text and dimensions.
func EmbeddingCacheKey(input string, dimensions int) string {
	sum := sha256.Sum256([]byte(input + "|" + strconv.Itoa(dimensions)))
	return hex.EncodeToString(sum[:])
}

// Get returns a cached embedding copy for input/dimensions.
func (c *LRUCache) Get(input string, dimensions int) ([]float32, bool) {
	key := EmbeddingCacheKey(input, dimensions)
	c.mu.Lock()
	defer c.mu.Unlock()

	element, ok := c.entries[key]
	if !ok {
		c.observeResult("miss")
		return nil, false
	}
	entry := element.Value.(*lruCacheEntry)
	if c.expired(entry) {
		c.removeElement(element)
		c.observeResult("miss")
		return nil, false
	}
	c.order.MoveToFront(element)
	c.observeResult("hit")
	return append([]float32(nil), entry.value...), true
}

// Put stores an embedding copy for input/dimensions.
func (c *LRUCache) Put(input string, dimensions int, value []float32) {
	key := EmbeddingCacheKey(input, dimensions)
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.entries[key]; ok {
		entry := element.Value.(*lruCacheEntry)
		entry.value = append([]float32(nil), value...)
		entry.expiresAt = c.expiry()
		c.order.MoveToFront(element)
		return
	}

	entry := &lruCacheEntry{key: key, value: append([]float32(nil), value...), expiresAt: c.expiry()}
	c.entries[key] = c.order.PushFront(entry)
	c.evictOverflow()
}

// GetIfPresent returns a cached embedding copy without invoking a loader.
func (c *LRUCache) GetIfPresent(_ context.Context, key string, dimensions int) ([]float32, bool) {
	return c.Get(key, dimensions)
}

// GetManyIfPresent returns cached embedding copies keyed by input index.
func (c *LRUCache) GetManyIfPresent(ctx context.Context, keys []string, dimensions int) map[int][]float32 {
	hits := make(map[int][]float32, len(keys))
	for i, key := range keys {
		if emb, ok := c.GetIfPresent(ctx, key, dimensions); ok {
			hits[i] = emb
		}
	}
	return hits
}

// Set stores an embedding copy under key/dimensions.
func (c *LRUCache) Set(_ context.Context, key string, dimensions int, value []float32) {
	c.Put(key, dimensions, value)
}

// GetOrLoad returns a cached embedding or stores the loader result on miss.
func (c *LRUCache) GetOrLoad(ctx context.Context, key string, dimensions int, loader func(context.Context) ([]float32, error)) ([]float32, error) {
	if value, ok := c.Get(key, dimensions); ok {
		return value, nil
	}
	value, err := loader(ctx)
	if err != nil {
		return nil, err
	}
	c.Put(key, dimensions, value)
	return append([]float32(nil), value...), nil
}

// Len returns the current number of stored entries, including entries that may
// be expired but have not yet been accessed.
func (c *LRUCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}

func (c *LRUCache) evictOverflow() {
	for c.maxSize > 0 && len(c.entries) > c.maxSize {
		oldest := c.order.Back()
		if oldest == nil {
			return
		}
		c.removeElement(oldest)
		if c.metrics != nil {
			c.metrics.ObserveCacheEviction()
		}
	}
}

func (c *LRUCache) removeElement(element *list.Element) {
	c.order.Remove(element)
	delete(c.entries, element.Value.(*lruCacheEntry).key)
}

func (c *LRUCache) expiry() time.Time {
	if c.ttl <= 0 {
		return time.Time{}
	}
	return time.Now().Add(c.ttl)
}

func (c *LRUCache) expired(entry *lruCacheEntry) bool {
	return !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt)
}

func (c *LRUCache) observeResult(result string) {
	if c.metrics != nil {
		c.metrics.ObserveCacheResult(result)
	}
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}

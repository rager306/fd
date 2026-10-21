// Package cache provides an in-memory LRU cache used as the L1 layer of fd's two-tier embedding cache (L1 memory + L2 Redis).
package cache

import (
	"context"
	"math/rand/v2"
	"sync"
	"time"
)

type l1Entry struct {
	value     []byte
	expiresAt time.Time
}

// LocalCache is an in-memory TTL-aware cache used as fd's L1 layer
// (Redis is the L2). It is goroutine-safe.
type LocalCache struct {
	mu       sync.Mutex
	data     map[string]l1Entry
	maxSize  int
	evictTTL time.Duration

	evictionObserver func(int)

	stopCh    chan struct{}
	doneCh    chan struct{}
	closeOnce sync.Once
}

// NewLocalCache returns a LocalCache with capacity maxSize entries and a
// background eviction loop running every evictTTL. If maxSize is <= 0,
// size eviction is disabled. If evictTTL is <= 0, background TTL eviction
// is disabled, but Get still expires entries lazily.
func NewLocalCache(maxSize int, evictTTL time.Duration) *LocalCache {
	c := &LocalCache{
		data:     make(map[string]l1Entry),
		maxSize:  maxSize,
		evictTTL: evictTTL,
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}
	if evictTTL > 0 {
		go c.evictLoop()
	} else {
		close(c.doneCh)
	}
	return c
}

// Close stops the background eviction loop. It is safe to call multiple times.
func (c *LocalCache) Close() error {
	c.closeOnce.Do(func() {
		close(c.stopCh)
		<-c.doneCh
	})
	return nil
}

// Get returns the value stored under key, or (nil, false) on miss
// or expired entry. The context is accepted for interface symmetry
// with the cache interfaces; it is not currently used by the in-memory
// implementation.
func (c *LocalCache) Get(_ context.Context, key string) ([]byte, bool) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.data[key]
	if !ok {
		return nil, false
	}
	if now.After(e.expiresAt) {
		delete(c.data, key)
		return nil, false
	}
	return e.value, true
}

// SetEvictionObserver registers a callback that receives the number of entries
// dropped by one capacity sweep. Pass nil to detach. The observer runs outside
// the cache lock, after the write completes; it must be non-blocking.
func (c *LocalCache) SetEvictionObserver(observer func(int)) {
	c.mu.Lock()
	c.evictionObserver = observer
	c.mu.Unlock()
}

// Set stores value under key with a per-entry TTL. If key already exists
// the value and TTL are refreshed in place; the entry is not duplicated. When
// the write pushes the cache past maxSize, the entries dropped to make room are
// reported to the eviction observer.
func (c *LocalCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	c.data[key] = l1Entry{
		value:     value,
		expiresAt: time.Now().Add(jitteredTTL(ttl)),
	}
	evicted := c.enforceMaxSizeLocked(key)
	observer := c.evictionObserver
	c.mu.Unlock()

	if evicted > 0 && observer != nil {
		observer(evicted)
	}
}

// ttlJitterFraction is the relative spread applied to positive TTLs, so entries
// written together do not expire in the same tick.
const ttlJitterFraction = 0.1

// jitteredTTL spreads a positive TTL uniformly over ±10%. Non-positive TTLs are
// returned unchanged: they mean "no expiry" or "expire immediately", and
// jittering either would change the caller's intent.
func jitteredTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return ttl
	}
	factor := 1 + (rand.Float64()*2-1)*ttlJitterFraction //nolint:gosec // G404: TTL jitter only; unpredictability is not a security requirement
	return time.Duration(float64(ttl) * factor)
}

// Delete removes the entry under key. No-op if the key is absent.
func (c *LocalCache) Delete(_ context.Context, key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

// Flush removes all entries from the local cache.
func (c *LocalCache) Flush(_ context.Context) {
	c.mu.Lock()
	c.data = make(map[string]l1Entry)
	c.mu.Unlock()
}

// Size returns the current number of local cache entries.
func (c *LocalCache) Size() int {
	return c.currentSize()
}

// enforceMaxSizeLocked drops entries until the cache is back within maxSize and
// returns how many it dropped. protectedKey — the entry just written — is never
// dropped. Callers must hold c.mu.
func (c *LocalCache) enforceMaxSizeLocked(protectedKey string) int {
	if c.maxSize <= 0 {
		return 0
	}
	evicted := 0
	for len(c.data) > c.maxSize {
		deleted := false
		for key := range c.data {
			if key == protectedKey {
				continue
			}
			delete(c.data, key)
			evicted++
			deleted = true
			break
		}
		if !deleted {
			return evicted
		}
	}
	return evicted
}

func (c *LocalCache) currentSize() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data)
}

func (c *LocalCache) evictLoop() {
	defer close(c.doneCh)
	ticker := time.NewTicker(c.evictTTL)
	defer ticker.Stop()
	for {
		select {
		case now := <-ticker.C:
			c.evictExpired(now)
		case <-c.stopCh:
			return
		}
	}
}

func (c *LocalCache) evictExpired(now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, value := range c.data {
		if now.After(value.expiresAt) {
			delete(c.data, key)
		}
	}
}

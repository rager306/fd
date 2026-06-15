// Package cache provides an in-memory LRU cache used as the L1 layer of fd's two-tier embedding cache (L1 memory + L2 Redis).
package cache

import (
	"context"
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

// Set stores value under key with a per-entry TTL. If key already exists
// the value and TTL are refreshed in place; the entry is not duplicated.
func (c *LocalCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = l1Entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	c.enforceMaxSizeLocked(key)
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

func (c *LocalCache) enforceMaxSizeLocked(protectedKey string) {
	if c.maxSize <= 0 {
		return
	}
	for len(c.data) > c.maxSize {
		deleted := false
		for key := range c.data {
			if key == protectedKey {
				continue
			}
			delete(c.data, key)
			deleted = true
			break
		}
		if !deleted {
			return
		}
	}
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

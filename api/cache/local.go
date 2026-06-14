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
	data     sync.Map
	maxSize  int
	mu       sync.Mutex
	size     int
	evictTTL time.Duration
}

// NewLocalCache returns a LocalCache with capacity maxSize entries and a
// background eviction loop running every evictTTL. If maxSize is <= 0,
// eviction is disabled (the cache grows unbounded until OOM).
func NewLocalCache(maxSize int, evictTTL time.Duration) *LocalCache {
	c := &LocalCache{maxSize: maxSize, evictTTL: evictTTL}
	go c.evictLoop()
	return c
}

// Get returns the value stored under key, or (nil, false) on miss
// or expired entry. The context is accepted for interface symmetry
// with the cache interface; it is not used.
func (c *LocalCache) Get(_ context.Context, key string) ([]byte, bool) {
	raw, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}
	e := raw.(*l1Entry)
	if time.Now().After(e.expiresAt) {
		c.data.Delete(key)
		c.decrSize()
		return nil, false
	}
	return e.value, true
}

// Set stores value under key with a per-entry TTL. If key already exists
// the value and TTL are refreshed in place; the entry is not duplicated.
func (c *LocalCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) {
	_, loaded := c.data.Load(key)
	c.data.Store(key, &l1Entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	})
	if loaded {
		return
	}

	c.incrSize()
	c.enforceMaxSize(key)
}

// Delete removes the entry under key. No-op if the key is absent.
func (c *LocalCache) Delete(_ context.Context, key string) {
	if _, loaded := c.data.LoadAndDelete(key); loaded {
		c.decrSize()
	}
}

func (c *LocalCache) incrSize() {
	c.mu.Lock()
	c.size++
	c.mu.Unlock()
}

func (c *LocalCache) decrSize() {
	c.mu.Lock()
	if c.size > 0 {
		c.size--
	}
	c.mu.Unlock()
}

func (c *LocalCache) enforceMaxSize(protectedKey string) {
	if c.maxSize <= 0 {
		return
	}

	for c.currentSize() > c.maxSize {
		deleted := false
		c.data.Range(func(key, _ any) bool {
			if key == protectedKey {
				return true
			}
			if _, loaded := c.data.LoadAndDelete(key); loaded {
				c.decrSize()
				deleted = true
			}
			return false
		})
		if !deleted {
			return
		}
	}
}

func (c *LocalCache) currentSize() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.size
}

func (c *LocalCache) evictLoop() {
	ticker := time.NewTicker(c.evictTTL)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		c.data.Range(func(key, value any) bool {
			if now.After(value.(*l1Entry).expiresAt) {
				c.data.Delete(key)
				c.decrSize()
			}
			return true
		})
	}
}

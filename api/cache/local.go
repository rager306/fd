// Package cache provides an in-memory LRU cache used as the L1 layer of fd's two-tier embedding cache (L1 memory + L2 Redis).
package cache

import (
	"container/list"
	"context"
	"math/rand/v2"
	"sync"
	"time"
)

type l1Entry struct {
	key       string
	value     []byte
	expiresAt time.Time
	bytes     int64
}

// LocalCache is an in-memory TTL-aware cache used as fd's L1 layer
// (Redis is the L2). It is goroutine-safe.
//
// Entries are kept in least-recently-used order: Get and Set move an entry to
// the front, and capacity eviction drops entries from the back. Two limits
// bound the cache, an entry count and a byte budget; a limit of 0 disables
// that dimension. The byte budget counts payload bytes as len(key)+len(value)
// per entry. That is payload accounting, not Go heap usage: map and list
// overhead and slice capacity slack are not counted.
type LocalCache struct {
	mu sync.Mutex
	// data maps each key to its element in lru; lru.Front is the most
	// recently used entry and lru.Back the least recently used.
	data       map[string]*list.Element
	lru        list.List
	maxEntries int
	maxBytes   int64
	bytes      int64
	evictTTL   time.Duration

	evictionObserver func(int)

	stopCh    chan struct{}
	doneCh    chan struct{}
	closeOnce sync.Once
}

// NewLocalCache returns a LocalCache with capacity maxSize entries and a
// background eviction loop running every evictTTL. If maxSize is <= 0, entry
// count eviction is disabled; the byte budget is disabled. If evictTTL is <= 0,
// background TTL eviction is disabled, but Get still expires entries lazily.
func NewLocalCache(maxSize int, evictTTL time.Duration) *LocalCache {
	return NewLocalCacheWithBudget(maxSize, 0, evictTTL)
}

// NewLocalCacheWithBudget returns a LocalCache bounded by maxEntries entries
// and maxBytes payload bytes. A limit of 0 disables it. See LocalCache for what
// the byte budget counts. Background TTL eviction follows the same rule as
// NewLocalCache.
func NewLocalCacheWithBudget(maxEntries int, maxBytes int64, evictTTL time.Duration) *LocalCache {
	c := &LocalCache{
		data:       make(map[string]*list.Element),
		maxEntries: maxEntries,
		maxBytes:   maxBytes,
		evictTTL:   evictTTL,
		stopCh:     make(chan struct{}),
		doneCh:     make(chan struct{}),
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

	elem, ok := c.data[key]
	if !ok {
		return nil, false
	}
	e := elem.Value.(*l1Entry)
	if now.After(e.expiresAt) {
		c.removeLocked(elem)
		return nil, false
	}
	c.lru.MoveToFront(elem)
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

// Set stores value under key with a per-entry TTL and marks the key as most
// recently used. If key already exists the value and TTL are refreshed in
// place; the entry is not duplicated. When the write pushes the cache past
// maxEntries or maxBytes, the entries dropped to make room are reported to the
// eviction observer. The entry just written is kept whenever another entry can
// be dropped instead.
func (c *LocalCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) {
	size := int64(len(key) + len(value))
	c.mu.Lock()
	if elem, ok := c.data[key]; ok {
		e := elem.Value.(*l1Entry)
		c.bytes -= e.bytes
		e.value = value
		e.expiresAt = time.Now().Add(jitteredTTL(ttl))
		e.bytes = size
		c.bytes += size
		c.lru.MoveToFront(elem)
	} else {
		e := &l1Entry{
			key:       key,
			value:     value,
			expiresAt: time.Now().Add(jitteredTTL(ttl)),
			bytes:     size,
		}
		c.data[key] = c.lru.PushFront(e)
		c.bytes += size
	}
	evicted := c.enforceLimitsLocked(key)
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
	if elem, ok := c.data[key]; ok {
		c.removeLocked(elem)
	}
	c.mu.Unlock()
}

// Flush removes all entries from the local cache.
func (c *LocalCache) Flush(_ context.Context) {
	c.mu.Lock()
	c.data = make(map[string]*list.Element)
	c.lru.Init()
	c.bytes = 0
	c.mu.Unlock()
}

// Size returns the current number of local cache entries.
func (c *LocalCache) Size() int {
	return c.currentSize()
}

// enforceLimitsLocked drops least recently used entries until the cache is back
// within maxEntries and maxBytes, and returns how many it dropped.
// protectedKey — the entry just written, which sits at the front — is never
// dropped while another entry remains. Callers must hold c.mu.
func (c *LocalCache) enforceLimitsLocked(protectedKey string) int {
	evicted := 0
	for c.overLimitLocked() {
		elem := c.lru.Back()
		if elem == nil || elem.Value.(*l1Entry).key == protectedKey {
			break
		}
		c.removeLocked(elem)
		evicted++
	}
	return evicted
}

// overLimitLocked reports whether either capacity limit is exceeded. Callers
// must hold c.mu.
func (c *LocalCache) overLimitLocked() bool {
	if c.maxEntries > 0 && len(c.data) > c.maxEntries {
		return true
	}
	return c.maxBytes > 0 && c.bytes > c.maxBytes
}

// removeLocked unlinks one entry and updates the byte counter. Callers must
// hold c.mu.
func (c *LocalCache) removeLocked(elem *list.Element) {
	e := elem.Value.(*l1Entry)
	delete(c.data, e.key)
	c.lru.Remove(elem)
	c.bytes -= e.bytes
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
	for _, elem := range c.data {
		if now.After(elem.Value.(*l1Entry).expiresAt) {
			c.removeLocked(elem)
		}
	}
}

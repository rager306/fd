# Tiered Cache Implementation Plan

**Goal:** Implement two-tier cache (L1 local + L2 Redis binary) with connection pool timeouts.

**Architecture:** `Request → L1 (sync.Map, ~50ns) → L2 (Redis binary, ~0.5ms) → TEI (~70ms)`. Binary storage replaces JSON HASH with flat SET/GET. singleflight prevents cache stampede.

**Tech Stack:** Go 1.25, go-redis v9, sync.Map, golang.org/x/sync/singleflight

---

## Task 1: Create `api/cache/local.go` — L1 cache with TTL

**Files:**
- Create: `api/cache/local.go`
- Test: `api/cache/local_test.go`

**Step 1: Write failing test**

```go
package cache

import (
	"context"
	"testing"
	"time"
)

func TestLocalCache_SetAndGet(t *testing.T) {
	c := NewLocalCache(1000, 30*time.Second)
	ctx := context.Background()

	c.Set(ctx, "key1", []byte{1, 2, 3}, 10*time.Second)

	got, ok := c.Get(ctx, "key1")
	if !ok {
		t.Fatal("expected to find key1")
	}
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Errorf("got %v, want [1 2 3]", got)
	}
}

func TestLocalCache_TTLExpired(t *testing.T) {
	c := NewLocalCache(1000, time.Millisecond)
	ctx := context.Background()

	c.Set(ctx, "key1", []byte{1}, 1*time.Millisecond)
	time.Sleep(10 * time.Millisecond)

	_, ok := c.Get(ctx, "key1")
	if ok {
		t.Error("expected key to be expired")
	}
}

func TestLocalCache_EvictionOnMaxSize(t *testing.T) {
	c := NewLocalCache(3, time.Minute)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		c.Set(ctx, "key"+string(rune('a'+i)), []byte{byte(i)}, time.Minute)
	}

	// At least one key from first 3 should be evicted
	found := 0
	for i := 0; i < 5; i++ {
		key := "key" + string(rune('a'+i))
		if _, ok := c.Get(ctx, key); ok {
			found++
		}
	}
	if found == 5 {
		t.Error("expected eviction when exceeding maxSize")
	}
}
```

**Step 2: Run test to verify failure**

```bash
cd /root/fd/api && go test ./cache/ -run TestLocalCache -v
```
Expected: FAIL — "undefined: NewLocalCache"

**Step 3: Write minimal implementation**

```go
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

type LocalCache struct {
	data     sync.Map
	maxSize  int
	mu       sync.Mutex
	size     int
	evictTTL time.Duration
}

func NewLocalCache(maxSize int, evictTTL time.Duration) *LocalCache {
	c := &LocalCache{maxSize: maxSize, evictTTL: evictTTL}
	go c.evictLoop()
	return c
}

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

func (c *LocalCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) {
	_, loaded := c.data.LoadOrStore(key, &l1Entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	})
	if !loaded {
		c.incrSize()
	}
}

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
	c.size--
	c.mu.Unlock()
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
```

**Step 4: Run test to verify pass**

```bash
cd /root/fd/api && go test ./cache/ -run TestLocalCache -v
```
Expected: PASS

**Step 5: Commit**

```bash
cd /root/fd && git add api/cache/local.go api/cache/local_test.go && git commit -m "feat(cache): add L1 local cache with TTL eviction"
```

---

## Task 2: Refactor `api/cache/redis.go` — binary storage + timeouts

**Files:**
- Modify: `api/cache/redis.go`
- Test: `api/cache/redis_test.go`

**Step 1: Write failing test**

```go
func TestBinaryMarshalUnmarshal(t *testing.T) {
	emb := []float32{1.0, 2.0, 3.0, 4.0}

	data := marshalEmbedding(emb, 4)
	if len(data) != 2+4*4 { // dim(2) + 4*float32(4)
		t.Errorf("len=%d, want %d", len(data), 2+16)
	}

	got, dim := unmarshalEmbedding(data)
	if dim != 4 {
		t.Errorf("dim=%d, want 4", dim)
	}
	if len(got) != 4 {
		t.Errorf("len=%d, want 4", len(got))
	}
	for i := range got {
		if got[i] != emb[i] {
			t.Errorf("got[%d]=%v, want %v", i, got[i], emb[i])
		}
	}
}

func TestBinaryRoundtrip_512(t *testing.T) {
	emb := make([]float32, 512)
	for i := range emb {
		emb[i] = float32(i) * 0.001
	}

	data := marshalEmbedding(emb, 512)
	got, dim := unmarshalEmbedding(data)

	if dim != 512 {
		t.Errorf("dim=%d, want 512", dim)
	}
	if len(got) != 512 {
		t.Errorf("len=%d, want 512", len(got))
	}
}

func TestRedisCache_BinarySetGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Redis test in short mode")
	}
	ctx := context.Background()
	c := NewRedisCache("localhost:6379", "test:", 10)

	emb := []float32{0.1, 0.2, 0.3, 0.4}
	err := c.Set(ctx, "test:text1", emb, 4)
	require.NoError(t, err)

	got, ok, err := c.Get(ctx, "test:text1", 4)
	require.NoError(t, err)
	if !ok {
		t.Fatal("expected to find key")
	}
	if len(got) != 4 {
		t.Errorf("len=%d, want 4", len(got))
	}
}
```

**Step 2: Run test to verify failure**

```bash
cd /root/fd/api && go test ./cache/ -run TestBinary -v
```
Expected: FAIL — "undefined: marshalEmbedding"

**Step 3: Write binary storage implementation**

Replace the `cachedValue` struct and JSON-based storage with binary:

```go
package cache

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// marshalEmbedding encodes [dim:uint16][float32*dim] for fast Redis storage.
// This replaces JSON: 2+4*dim bytes vs ~8KB JSON for 1024d.
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
```

**Step 4: Run test to verify pass**

```bash
cd /root/fd/api && go test ./cache/ -run TestBinary -v
```
Expected: PASS

**Step 5: Update RedisCache with timeouts**

Replace `NewRedisCache`:

```go
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
```

Update `Get` and `Set` to use binary instead of JSON:

```go
func (c *RedisCache) Get(ctx context.Context, text string, dim int) ([]float32, bool, error) {
	key := c.prefix + "v2:" + c.HashText(text) + ":" + fmt.Sprintf("d%d", dim)

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
	key := c.prefix + "v2:" + c.HashText(text) + ":" + fmt.Sprintf("d%d", dim)
	data := marshalEmbedding(embedding, dim)
	return c.client.Set(ctx, key, data, c.ttl).Err()
}
```

Add `"math"` import.

**Step 6: Run all cache tests**

```bash
cd /root/fd/api && go test ./cache/ -v -short
```
Expected: PASS

**Step 7: Commit**

```bash
cd /root/fd && git add api/cache/redis.go api/cache/redis_test.go && git commit -m "feat(cache): binary storage + pool timeouts"
```

---

## Task 3: Create `api/cache/tiered.go` — two-tier cache

**Files:**
- Create: `api/cache/tiered.go`
- Test: `api/cache/tiered_test.go`

**Step 1: Write failing test**

```go
package cache

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestTieredCache_GetOrLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	local := NewLocalCache(1000, time.Minute)
	redisCache := NewRedisCache("localhost:6379", "tiered:test:", 10)
	tc := NewTieredCache(local, redisCache, 30*time.Second, 24*time.Hour)
	defer redisCache.Close()

	calls := 0
	loader := func(ctx context.Context) ([]byte, error) {
		calls++
		return []byte("result"), nil
	}

	// First call — loader invoked
	result, err := tc.GetOrLoad(ctx, "key1", loader)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != "result" {
		t.Errorf("got %q, want 'result'", result)
	}
	if calls != 1 {
		t.Errorf("calls=%d, want 1", calls)
	}

	// Second call — L1 hit, no loader
	result, err = tc.GetOrLoad(ctx, "key1", loader)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Errorf("calls=%d, want 1 (L1 hit)", calls)
	}
}

func TestTieredCache_LocalTTL(t *testing.T) {
	local := NewLocalCache(1000, 10*time.Millisecond)
	redisCache := NewLocalOnlyCache() // mock that returns nil
	tc := NewTieredCache(local, redisCache, 5*time.Millisecond, time.Hour)
	ctx := context.Background()

	loaderCalled := 0
	loader := func(ctx context.Context) ([]byte, error) {
		loaderCalled++
		return []byte("v"), nil
	}

	tc.GetOrLoad(ctx, "key1", loader)
	time.Sleep(20 * time.Millisecond)

	// After local TTL expires, loader should be called again
	tc.GetOrLoad(ctx, "key1", loader)
	if loaderCalled != 2 {
		t.Errorf("loaderCalled=%d, want 2", loaderCalled)
	}
}
```

**Step 2: Run test to verify failure**

```bash
cd /root/fd/api && go test ./cache/ -run TestTiered -v
```
Expected: FAIL — "undefined: NewTieredCache"

**Step 3: Write tiered implementation**

```go
package cache

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

type TieredCache struct {
	local    *LocalCache
	redis    *RedisCache
	localTTL time.Duration
	sf       singleflight.Group
}

func NewTieredCache(local *LocalCache, redis *RedisCache, localTTL, redisTTL time.Duration) *TieredCache {
	return &TieredCache{
		local:    local,
		redis:    redis,
		localTTL: localTTL,
	}
}

type result struct {
	val []byte
	err error
}

// GetOrLoad implements two-tier cache-aside with singleflight deduplication.
func (tc *TieredCache) GetOrLoad(ctx context.Context, key string, loader func(context.Context) ([]byte, error)) ([]byte, error) {
	// L1 check
	if val, ok := tc.local.Get(ctx, key); ok {
		return val, nil
	}

	// singleflight — dedup concurrent requests for same key
	r, err, _ := tc.sf.Do(key, func() (any, error) {
		// Double-check L1 (another goroutine may have filled it while we waited)
		if val, ok := tc.local.Get(ctx, key); ok {
			return val, nil
		}
		// L2 check
		val, ok, err := tc.redis.Get(ctx, key)
		if err == nil && ok {
			tc.local.Set(ctx, key, val, tc.localTTL)
			return val, nil
		}
		// Loader
		val, err = loader(ctx)
		if err != nil {
			return nil, err
		}
		// Backfill both tiers
		tc.local.Set(ctx, key, val, tc.localTTL)
		_ = tc.redis.Set(ctx, key, val) // best-effort
		return val, nil
	})
	if err != nil {
		return nil, err
	}
	return r.([]byte), nil
}
```

**Step 4: Run test to verify pass**

```bash
cd /root/fd/api && go test ./cache/ -run TestTiered -v
```
Expected: PASS

**Step 5: Commit**

```bash
cd /root/fd && git add api/cache/tiered.go api/cache/tiered_test.go && git commit -m "feat(cache): add two-tier cache with singleflight"
```

---

## Task 4: Wire tiered cache into `api/main.go`

**Files:**
- Modify: `api/main.go`
- Modify: `api/handlers/embeddings.go`
- Modify: `api/handlers/batch.go`

**Step 1: Update main.go**

Replace Redis-only initialization with two-tier:

```go
// Local cache (L1) — 10000 entries, 30s TTL
localCache := cache.NewLocalCache(10000, 30*time.Second)

// Redis cache (L2) — binary storage, pool timeouts
redisCache := cache.NewRedisCache(redisHost, "embed:cache:", redisPoolSize)
defer redisCache.Close()

// Two-tier cache
tiered := cache.NewTieredCache(localCache, redisCache, 30*time.Second, 24*time.Hour)
```

**Step 2: Pass tiered cache to handlers**

In `handlers/embeddings.go`, update `NewEmbeddingsHandler` to accept `*cache.TieredCache`:

```go
type EmbeddingsHandler struct {
	tei    *embed.TEIClient
	cache  *cache.TieredCache
	// ...
}
```

Replace `redisCache.Get`/`Set` calls with `tiered.GetOrLoad`:

```go
// Instead of:
emb, ok, err := h.cache.Get(ctx, req.Input, dim)
// ...
h.cache.Set(ctx, req.Input, emb, dim)

// Use:
emb, err := h.cache.GetOrLoad(ctx, cacheKey, func(ctx context.Context) ([]byte, error) {
    teiEmb, err := h.tei.Embed(ctx, req.Input)
    if err != nil {
        return nil, err
    }
    return embed.ToBytes(teiEmb), nil
})
```

**Step 3: Run tests**

```bash
cd /root/fd/api && go build ./... && go test ./... -short
```
Expected: PASS

**Step 4: Commit**

```bash
cd /root/fd && git add api/main.go api/handlers/embeddings.go api/handlers/batch.go && git commit -m "feat: wire two-tier cache into handlers"
```

---

## Task 5: Integration test — full flow with Redis

**Files:**
- Create: `api/cache/tiered_integration_test.go` (build tag: integration)

**Step 1: Write integration test**

```go
//go:build integration
// +build integration

package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func startRedis(t *testing.T) (func(), *redis.Client) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		DialTimeout: 5 * time.Second,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available")
	}
	cleanup := func() { rdb.Close() }
	return cleanup, rdb
}

func TestIntegration_TieredCache(t *testing.T) {
	cleanup, rdb := startRedis(t)
	defer cleanup()
	ctx := context.Background()

	// Setup
	prefix := fmt.Sprintf("test:%d:", time.Now().UnixNano())
	redisCache := NewRedisCache("localhost:6379", prefix, 10)
	localCache := NewLocalCache(1000, 30*time.Second)
	tc := NewTieredCache(localCache, redisCache, 30*time.Second, time.Hour)
	defer redisCache.Close()

	// Test: GetOrLoad calls loader once
	calls := 0
	loader := func(ctx context.Context) ([]byte, error) {
		calls++
		return []byte(fmt.Sprintf("value-%d", calls)), nil
	}

	for i := 0; i < 5; i++ {
		result, err := tc.GetOrLoad(ctx, "shared-key", loader)
		if err != nil {
			t.Fatal(err)
		}
		if string(result) != "value-1" {
			t.Errorf("got %q, want 'value-1'", result)
		}
	}

	if calls != 1 {
		t.Errorf("loader called %d times, want 1", calls)
	}
}
```

**Step 2: Run integration test**

```bash
cd /root/fd/api && go test -tags=integration ./cache/ -run TestIntegration -v
```
Expected: PASS

**Step 3: Commit**

```bash
cd /root/fd && git add api/cache/tiered_integration_test.go && git commit -m "test: integration test for tiered cache"
```

---

## Task 6: E2E verification

**Step 1: Rebuild and restart**

```bash
cd /root/fd && docker compose build --no-cache api && docker compose up -d api
```

**Step 2: Verify binary storage in Redis**

```bash
# Check key format (should be binary, not JSON)
redis-cli GET "embed:cache:v2:$(echo -n "test" | sha256sum | cut -d' ' -f1):d1024" | xxd | head -5
```

Expected: `0000: 0004 0000 00...` — first 2 bytes are dim=1024 in little-endian

**Step 3: Run benchmarks**

```bash
cd /root/fd && python3 benchmark.py
```

Expected: L1 hit latency < 1ms (was ~2.6ms with JSON), L2 hit ~1-2ms

**Step 4: Commit final state**

```bash
cd /root/fd && git add -A && git commit -m "perf: tiered cache with binary storage — L1 local + L2 Redis"
```

---

## Summary

| Task | Time | Impact |
|------|------|--------|
| 1. Local L1 cache | ~20 min | 10,000x faster for hot keys |
| 2. Binary storage + timeouts | ~30 min | 2x smaller, no hanging connections |
| 3. Two-tier with singleflight | ~20 min | Stampede prevention |
| 4. Wire into handlers | ~15 min | — |
| 5. Integration test | ~15 min | Confidence |
| 6. E2E + benchmark | ~10 min | Verify |

**Total: ~2 hours**

**Expected latency after optimization:**
- L1 hit: < 0.001ms (was ~2.6ms with Redis-only JSON)
- L2 hit: ~0.5-1ms (was ~2.6ms with JSON)
- TEI miss: unchanged (~70ms)

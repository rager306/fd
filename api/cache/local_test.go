package cache

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLocalCache_CloseIsIdempotent(t *testing.T) {
	c := NewLocalCache(1000, time.Millisecond)
	if err := c.Close(); err != nil {
		t.Fatalf("first Close returned error: %v", err)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("second Close returned error: %v", err)
	}
}

func TestLocalCacheFlushRemovesAllEntries(t *testing.T) {
	c := NewLocalCache(1000, time.Minute)
	defer func() {
		if err := c.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()
	ctx := context.Background()
	c.Set(ctx, "a", []byte{1}, time.Minute)
	c.Set(ctx, "b", []byte{2}, time.Minute)
	if c.currentSize() != 2 {
		t.Fatalf("size before flush = %d, want 2", c.currentSize())
	}
	c.Flush(ctx)
	if c.currentSize() != 0 {
		t.Fatalf("size after flush = %d, want 0", c.currentSize())
	}
	if _, ok := c.Get(ctx, "a"); ok {
		t.Fatal("key a survived flush")
	}
}

func TestLocalCacheSizeReportsCurrentEntries(t *testing.T) {
	c := NewLocalCache(1000, time.Minute)
	defer func() {
		if err := c.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()
	ctx := context.Background()
	c.Set(ctx, "a", []byte{1}, time.Minute)
	c.Set(ctx, "b", []byte{2}, time.Minute)
	if got := c.Size(); got != 2 {
		t.Fatalf("Size = %d, want 2", got)
	}
}

func TestLocalCache_ConcurrentOverwriteKeepsSingleEntry(t *testing.T) {
	c := NewLocalCache(1000, time.Minute)
	defer func() {
		if err := c.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()
	ctx := context.Background()
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			c.Set(ctx, "same-key", []byte{1}, time.Minute)
		}()
	}
	close(start)
	wg.Wait()

	if c.currentSize() != 1 {
		t.Fatalf("size=%d, want 1 after concurrent overwrites", c.currentSize())
	}
	if _, ok := c.Get(ctx, "same-key"); !ok {
		t.Fatal("expected same-key to be retained")
	}
}

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

func TestLocalCache_SetRefreshesExistingValueAndTTL(t *testing.T) {
	c := NewLocalCache(1000, time.Minute)
	ctx := context.Background()

	c.Set(ctx, "key1", []byte{1}, 1*time.Millisecond)
	c.Set(ctx, "key1", []byte{2, 3}, time.Minute)
	time.Sleep(10 * time.Millisecond)

	got, ok := c.Get(ctx, "key1")
	if !ok {
		t.Fatal("expected refreshed key to still exist")
	}
	if len(got) != 2 || got[0] != 2 || got[1] != 3 {
		t.Fatalf("got %v, want [2 3]", got)
	}
	if c.currentSize() != 1 {
		t.Fatalf("size=%d, want 1 after overwrite", c.currentSize())
	}
}

func TestLocalCache_EnforcesMaxSize(t *testing.T) {
	c := NewLocalCache(2, time.Minute)
	ctx := context.Background()

	c.Set(ctx, "key1", []byte{1}, time.Minute)
	c.Set(ctx, "key2", []byte{2}, time.Minute)
	c.Set(ctx, "key3", []byte{3}, time.Minute)

	if c.currentSize() > 2 {
		t.Fatalf("size=%d, want <= 2", c.currentSize())
	}
	if _, ok := c.Get(ctx, "key3"); !ok {
		t.Fatal("expected newest key to be retained")
	}
}

func TestLocalCache_Delete(t *testing.T) {
	c := NewLocalCache(1000, time.Minute)
	ctx := context.Background()

	c.Set(ctx, "key1", []byte{1}, time.Minute)
	c.Delete(ctx, "key1")

	_, ok := c.Get(ctx, "key1")
	if ok {
		t.Error("expected cache miss after delete")
	}
}

func TestLocalCache_NotFound(t *testing.T) {
	c := NewLocalCache(1000, time.Minute)
	ctx := context.Background()

	_, ok := c.Get(ctx, "nonexistent")
	if ok {
		t.Error("expected cache miss for nonexistent key")
	}
}

func TestLocalCacheReportsCapacityEvictions(t *testing.T) {
	const maxSize = 2
	c := NewLocalCache(maxSize, 0)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})

	var sweeps []int
	c.SetEvictionObserver(func(n int) { sweeps = append(sweeps, n) })

	ctx := context.Background()
	for _, key := range []string{"k1", "k2", "k3", "k4", "k5"} {
		c.Set(ctx, key, []byte("value"), time.Minute)
	}

	evicted := 0
	for _, n := range sweeps {
		evicted += n
	}
	if evicted != 3 {
		t.Fatalf("evicted %d entries across %d sweeps, want 3", evicted, len(sweeps))
	}
	if got := c.Size(); got != maxSize {
		t.Fatalf("Size() = %d, want %d", got, maxSize)
	}
}

func TestJitteredTTLIsWithinTenPercent(t *testing.T) {
	const ttl = 30 * time.Second
	low, high := ttl*9/10, ttl*11/10
	seen := make(map[time.Duration]struct{})
	for i := 0; i < 1000; i++ {
		got := jitteredTTL(ttl)
		if got < low || got > high {
			t.Fatalf("jitteredTTL(%s) = %s, want within [%s, %s]", ttl, got, low, high)
		}
		seen[got] = struct{}{}
	}
	if len(seen) < 2 {
		t.Fatalf("jitteredTTL produced %d distinct values, want at least 2", len(seen))
	}
}

func TestJitteredTTLLeavesNonPositiveTTLAlone(t *testing.T) {
	for _, ttl := range []time.Duration{0, -time.Second} {
		if got := jitteredTTL(ttl); got != ttl {
			t.Fatalf("jitteredTTL(%s) = %s, want unchanged", ttl, got)
		}
	}
}

func TestLocalCacheSetJittersExpiry(t *testing.T) {
	c := NewLocalCache(1000, 0)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})

	const (
		ttl     = time.Minute
		entries = 100
	)
	ctx := context.Background()
	before := time.Now()
	for i := 0; i < entries; i++ {
		c.Set(ctx, fmt.Sprintf("k%d", i), []byte{1}, ttl)
	}
	after := time.Now()

	low, high := before.Add(ttl*9/10), after.Add(ttl*11/10)
	c.mu.Lock()
	defer c.mu.Unlock()
	expiries := make(map[time.Time]struct{}, entries)
	for key, elem := range c.data {
		e := elem.Value.(*l1Entry)
		if e.expiresAt.Before(low) || e.expiresAt.After(high) {
			t.Fatalf("%s expiresAt = %s, want within [%s, %s]", key, e.expiresAt, low, high)
		}
		expiries[e.expiresAt] = struct{}{}
	}
	if len(expiries) < 2 {
		t.Fatalf("Set produced %d distinct expiresAt values, want at least 2", len(expiries))
	}
}

func TestLocalCacheWithoutObserverStillEvicts(t *testing.T) {
	c := NewLocalCache(1, 0)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})

	ctx := context.Background()
	c.Set(ctx, "k1", []byte("value"), time.Minute)
	c.Set(ctx, "k2", []byte("value"), time.Minute)

	if got := c.Size(); got != 1 {
		t.Fatalf("Size() = %d, want 1", got)
	}
	if _, ok := c.Get(ctx, "k2"); !ok {
		t.Fatal("most recently written entry was dropped")
	}
}

func TestLocalCache_LRUOrderEvictsLeastRecentlyUsed(t *testing.T) {
	c := NewLocalCache(3, 0)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})
	ctx := context.Background()

	c.Set(ctx, "a", []byte("1"), time.Minute)
	c.Set(ctx, "b", []byte("1"), time.Minute)
	c.Set(ctx, "c", []byte("1"), time.Minute)
	// Recency order is now c, b, a. Touching a makes b the oldest entry.
	if _, ok := c.Get(ctx, "a"); !ok {
		t.Fatal("expected a to be present")
	}
	c.Set(ctx, "d", []byte("1"), time.Minute)

	if _, ok := c.Get(ctx, "b"); ok {
		t.Fatal("b should have been evicted as least recently used")
	}
	for _, key := range []string{"a", "c", "d"} {
		if _, ok := c.Get(ctx, key); !ok {
			t.Fatalf("%s should have survived eviction", key)
		}
	}
	assertLocalCacheAccounting(t, c)
}

func TestLocalCache_KeepsHotKeyUnderInsertPressure(t *testing.T) {
	c := NewLocalCache(4, 0)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})
	ctx := context.Background()

	c.Set(ctx, "hot", []byte("hot"), time.Minute)
	for i := 0; i < 200; i++ {
		if _, ok := c.Get(ctx, "hot"); !ok {
			t.Fatalf("hot key evicted after %d cold inserts", i)
		}
		c.Set(ctx, fmt.Sprintf("cold-%d", i), []byte("cold"), time.Minute)
	}
	if _, ok := c.Get(ctx, "hot"); !ok {
		t.Fatal("hot key evicted at the end of the run")
	}
	assertLocalCacheAccounting(t, c)
}

func TestLocalCache_MaxBytesEvictsLeastRecentlyUsed(t *testing.T) {
	// Every key is one byte and every value is four bytes, so one entry
	// accounts for five payload bytes and the budget holds three entries.
	c := NewLocalCacheWithBudget(0, 15, 0)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})
	ctx := context.Background()

	c.Set(ctx, "a", []byte("1111"), time.Minute)
	c.Set(ctx, "b", []byte("1111"), time.Minute)
	c.Set(ctx, "c", []byte("1111"), time.Minute)
	if c.bytes != 15 {
		t.Fatalf("bytes = %d, want 15", c.bytes)
	}

	c.Set(ctx, "d", []byte("2222"), time.Minute)
	if _, ok := c.Get(ctx, "a"); ok {
		t.Fatal("a should have been evicted to make room for d")
	}
	if _, ok := c.Get(ctx, "d"); !ok {
		t.Fatal("the entry just written was evicted")
	}

	// Touching b makes c the oldest entry before e is inserted.
	if _, ok := c.Get(ctx, "b"); !ok {
		t.Fatal("expected b to be present")
	}
	c.Set(ctx, "e", []byte("3333"), time.Minute)
	if _, ok := c.Get(ctx, "c"); ok {
		t.Fatal("c should have been evicted as least recently used")
	}
	if _, ok := c.Get(ctx, "e"); !ok {
		t.Fatal("the entry just written was evicted")
	}
	if c.bytes > 15 {
		t.Fatalf("bytes = %d, want <= 15", c.bytes)
	}
	assertLocalCacheAccounting(t, c)
}

func TestLocalCache_MaxBytesKeepsOversizedWriteAlone(t *testing.T) {
	c := NewLocalCacheWithBudget(0, 4, 0)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})
	ctx := context.Background()

	// 3 key bytes + 5 value bytes exceed the budget on their own, so the
	// write must not evict itself.
	c.Set(ctx, "big", []byte("12345"), time.Minute)
	if _, ok := c.Get(ctx, "big"); !ok {
		t.Fatal("oversized entry evicted itself")
	}
	if c.Size() != 1 {
		t.Fatalf("Size() = %d, want 1", c.Size())
	}
	assertLocalCacheAccounting(t, c)
}

func TestLocalCache_TracksPayloadBytesAcrossMutations(t *testing.T) {
	c := NewLocalCache(0, 0)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})
	ctx := context.Background()

	c.Set(ctx, "a", []byte("123"), time.Minute) // 4 payload bytes
	c.Set(ctx, "b", []byte("12"), time.Minute)  // 3 payload bytes
	if c.bytes != 7 {
		t.Fatalf("bytes after inserts = %d, want 7", c.bytes)
	}

	c.Set(ctx, "a", []byte("12345"), time.Minute) // a grows to 6
	if c.bytes != 9 {
		t.Fatalf("bytes after grow = %d, want 9", c.bytes)
	}

	c.Set(ctx, "a", []byte("1"), time.Minute) // a shrinks to 2
	if c.bytes != 5 {
		t.Fatalf("bytes after shrink = %d, want 5", c.bytes)
	}

	c.Delete(ctx, "a")
	c.Delete(ctx, "missing")
	if c.bytes != 3 {
		t.Fatalf("bytes after delete = %d, want 3", c.bytes)
	}

	c.Flush(ctx)
	if c.bytes != 0 {
		t.Fatalf("bytes after flush = %d, want 0", c.bytes)
	}
	assertLocalCacheAccounting(t, c)
}

func TestLocalCache_TTLExpiryUpdatesByteAccounting(t *testing.T) {
	c := NewLocalCache(0, 0)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})
	ctx := context.Background()

	c.Set(ctx, "short", []byte("v"), time.Millisecond)
	c.Set(ctx, "long", []byte("v"), time.Minute)
	time.Sleep(10 * time.Millisecond)

	if _, ok := c.Get(ctx, "short"); ok {
		t.Fatal("expected short entry to be expired")
	}
	if c.bytes != int64(len("long")+1) {
		t.Fatalf("bytes after lazy expiry = %d, want %d", c.bytes, len("long")+1)
	}
	assertLocalCacheAccounting(t, c)
}

func TestLocalCache_BackgroundEvictionRemovesExpiredEntries(t *testing.T) {
	c := NewLocalCache(0, time.Millisecond)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})
	ctx := context.Background()

	c.Set(ctx, "k", []byte("v"), time.Millisecond)
	deadline := time.Now().Add(2 * time.Second)
	for c.Size() != 0 {
		if time.Now().After(deadline) {
			t.Fatal("background eviction did not remove the expired entry")
		}
		time.Sleep(time.Millisecond)
	}
	if c.bytes != 0 {
		t.Fatalf("bytes = %d, want 0", c.bytes)
	}
	assertLocalCacheAccounting(t, c)
}

func TestLocalCacheReportsByteBudgetEvictions(t *testing.T) {
	c := NewLocalCacheWithBudget(0, 10, 0)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})

	evicted := 0
	c.SetEvictionObserver(func(n int) { evicted += n })

	ctx := context.Background()
	for _, key := range []string{"a", "b", "c"} {
		c.Set(ctx, key, make([]byte, 4), time.Minute) // 5 payload bytes each
	}
	if evicted != 1 {
		t.Fatalf("evicted %d entries, want 1", evicted)
	}
	assertLocalCacheAccounting(t, c)
}

func TestLocalCache_ConcurrentSetGetDelete(t *testing.T) {
	c := NewLocalCacheWithBudget(64, 4096, 0)
	t.Cleanup(func() {
		if err := c.Close(); err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	})
	ctx := context.Background()

	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				key := fmt.Sprintf("k%d", (g*7+i)%100)
				switch i % 3 {
				case 0:
					c.Set(ctx, key, []byte(strings.Repeat("v", i%8)), time.Minute)
				case 1:
					c.Get(ctx, key)
				default:
					c.Delete(ctx, key)
				}
			}
		}(g)
	}
	wg.Wait()

	if c.Size() > 64 {
		t.Fatalf("Size() = %d, want <= 64", c.Size())
	}
	if c.bytes > 4096 {
		t.Fatalf("bytes = %d, want <= 4096", c.bytes)
	}
	assertLocalCacheAccounting(t, c)
}

// assertLocalCacheAccounting verifies that the byte counter matches the live
// entries and that the map and LRU list hold the same entries.
func assertLocalCacheAccounting(t *testing.T, c *LocalCache) {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()

	var sum int64
	for _, elem := range c.data {
		e := elem.Value.(*l1Entry)
		sum += e.bytes
		if e.bytes != int64(len(e.key)+len(e.value)) {
			t.Fatalf("entry %q accounts %d bytes, want %d", e.key, e.bytes, len(e.key)+len(e.value))
		}
	}
	if sum != c.bytes {
		t.Fatalf("byte counter = %d, want %d", c.bytes, sum)
	}
	if len(c.data) != c.lru.Len() {
		t.Fatalf("map holds %d entries, LRU list holds %d", len(c.data), c.lru.Len())
	}
}

package cache

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type lruMetricsSpy struct {
	hits      atomic.Int64
	misses    atomic.Int64
	evictions atomic.Int64
}

func (m *lruMetricsSpy) ObserveCacheResult(result string) {
	switch result {
	case "hit":
		m.hits.Add(1)
	case "miss":
		m.misses.Add(1)
	}
}

func (m *lruMetricsSpy) ObserveCacheEviction() {
	m.evictions.Add(1)
}

func TestLRUCacheGetPutCopiesValues(t *testing.T) {
	metrics := &lruMetricsSpy{}
	cache := NewLRUCache(10, time.Hour, metrics)
	value := []float32{1, 2, 3}
	cache.Put("hello", 1024, value)
	value[0] = 99

	got, ok := cache.Get("hello", 1024)
	if !ok {
		t.Fatal("Get returned miss, want hit")
	}
	if got[0] != 1 {
		t.Fatalf("cached value was mutated: %#v", got)
	}
	got[0] = 42
	gotAgain, ok := cache.Get("hello", 1024)
	if !ok || gotAgain[0] != 1 {
		t.Fatalf("Get should return a copy, got ok=%v value=%#v", ok, gotAgain)
	}
	if metrics.hits.Load() != 2 {
		t.Fatalf("hits = %d, want 2", metrics.hits.Load())
	}
}

func TestLRUCacheMissAndDimensionKey(t *testing.T) {
	metrics := &lruMetricsSpy{}
	cache := NewLRUCache(10, time.Hour, metrics)
	cache.Put("hello", 1024, []float32{1})

	if _, ok := cache.Get("hello", 512); ok {
		t.Fatal("Get with different dimensions should miss")
	}
	if metrics.misses.Load() != 1 {
		t.Fatalf("misses = %d, want 1", metrics.misses.Load())
	}
}

func TestLRUCacheEvictsLeastRecentlyUsed(t *testing.T) {
	metrics := &lruMetricsSpy{}
	cache := NewLRUCache(2, time.Hour, metrics)
	cache.Put("one", 1024, []float32{1})
	cache.Put("two", 1024, []float32{2})
	if _, ok := cache.Get("one", 1024); !ok {
		t.Fatal("one should be present before eviction")
	}
	cache.Put("three", 1024, []float32{3})

	if _, ok := cache.Get("two", 1024); ok {
		t.Fatal("two should be evicted as least recently used")
	}
	if _, ok := cache.Get("one", 1024); !ok {
		t.Fatal("one should remain after recent access")
	}
	if got := metrics.evictions.Load(); got != 1 {
		t.Fatalf("evictions = %d, want 1", got)
	}
}

func TestLRUCacheExpiresEntries(t *testing.T) {
	cache := NewLRUCache(10, time.Millisecond, nil)
	cache.Put("hello", 1024, []float32{1})
	time.Sleep(2 * time.Millisecond)

	if _, ok := cache.Get("hello", 1024); ok {
		t.Fatal("expired entry should miss")
	}
	if got := cache.Len(); got != 0 {
		t.Fatalf("Len after expired get = %d, want 0", got)
	}
}

func TestLRUCacheFromEnv(t *testing.T) {
	t.Setenv("FD_CACHE_SIZE", "7")
	t.Setenv("FD_CACHE_TTL_HOURS", "2")
	cache := NewLRUCacheFromEnv(nil)

	if cache.maxSize != 7 {
		t.Fatalf("maxSize = %d, want 7", cache.maxSize)
	}
	if cache.ttl != 2*time.Hour {
		t.Fatalf("ttl = %s, want 2h", cache.ttl)
	}
}

func TestEmbeddingCacheKeyStableAndDimensionAware(t *testing.T) {
	one := EmbeddingCacheKey("hello", 1024)
	two := EmbeddingCacheKey("hello", 1024)
	three := EmbeddingCacheKey("hello", 512)
	if one != two {
		t.Fatalf("keys should be stable: %q != %q", one, two)
	}
	if one == three {
		t.Fatal("keys should include dimensions")
	}
}

func TestLRUCacheConcurrentAccess(t *testing.T) {
	cache := NewLRUCache(100, time.Hour, nil)
	var wg sync.WaitGroup
	for worker := range 8 {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := range 100 {
				input := string(rune('a' + (worker+i)%26))
				cache.Put(input, 1024, []float32{float32(i)})
				_, _ = cache.Get(input, 1024)
			}
		}(worker)
	}
	wg.Wait()
}

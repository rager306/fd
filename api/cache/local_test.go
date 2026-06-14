package cache

import (
	"context"
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

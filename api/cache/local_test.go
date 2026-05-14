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

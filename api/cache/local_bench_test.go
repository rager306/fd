package cache

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func BenchmarkLocalCacheGetHit(b *testing.B) {
	c := NewLocalCache(1000, 0)
	b.Cleanup(func() { _ = c.Close() })
	ctx := context.Background()
	c.Set(ctx, "hot", make([]byte, 512), time.Minute)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := c.Get(ctx, "hot"); !ok {
			b.Fatal("unexpected cache miss")
		}
	}
}

func BenchmarkLocalCacheSet(b *testing.B) {
	c := NewLocalCache(1000, 0)
	b.Cleanup(func() { _ = c.Close() })
	ctx := context.Background()
	keys := make([]string, 256)
	for i := range keys {
		keys[i] = fmt.Sprintf("key-%d", i)
	}
	value := make([]byte, 512)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(ctx, keys[i%len(keys)], value, time.Minute)
	}
}

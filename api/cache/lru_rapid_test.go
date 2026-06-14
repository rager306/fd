package cache

import (
	"context"
	"testing"
	"time"

	"pgregory.net/rapid"
)

func TestLRUCacheWrapperMethodsCopyValues_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		key := rapid.StringMatching(`[a-zA-Z0-9_-]{1,32}`).Draw(t, "key")
		dimensions := rapid.IntRange(1, 64).Draw(t, "dimensions")
		value := rapid.SliceOfN(rapid.Float32Range(-10_000, 10_000), dimensions, dimensions).Draw(t, "value")

		cache := NewLRUCache(128, time.Hour, nil)
		cache.Set(context.Background(), key, dimensions, value)
		value[0]++

		got, ok := cache.GetIfPresent(context.Background(), key, dimensions)
		if !ok {
			t.Fatal("GetIfPresent returned miss after Set")
		}
		if got[0] == value[0] {
			t.Fatalf("cache retained caller slice alias: got[0]=%v mutated[0]=%v", got[0], value[0])
		}
		got[0]++
		gotAgain, ok := cache.GetIfPresent(context.Background(), key, dimensions)
		if !ok {
			t.Fatal("GetIfPresent returned miss after read")
		}
		if gotAgain[0] == got[0] {
			t.Fatalf("cache returned mutable internal slice: gotAgain[0]=%v mutatedRead[0]=%v", gotAgain[0], got[0])
		}
	})
}

func TestLRUCacheGetOrLoadCallsLoaderOnlyOnMiss_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		key := rapid.StringMatching(`[a-z]{1,24}`).Draw(t, "key")
		dimensions := rapid.IntRange(1, 32).Draw(t, "dimensions")
		value := rapid.SliceOfN(rapid.Float32Range(-1_000, 1_000), dimensions, dimensions).Draw(t, "value")

		cache := NewLRUCache(128, time.Hour, nil)
		loads := 0
		loader := func(context.Context) ([]float32, error) {
			loads++
			return append([]float32(nil), value...), nil
		}

		first, err := cache.GetOrLoad(context.Background(), key, dimensions, loader)
		if err != nil {
			t.Fatalf("first GetOrLoad error: %v", err)
		}
		second, err := cache.GetOrLoad(context.Background(), key, dimensions, loader)
		if err != nil {
			t.Fatalf("second GetOrLoad error: %v", err)
		}
		if loads != 1 {
			t.Fatalf("loader calls = %d, want 1", loads)
		}
		if len(first) != dimensions || len(second) != dimensions {
			t.Fatalf("lengths = %d/%d, want %d", len(first), len(second), dimensions)
		}
		first[0]++
		third, ok := cache.GetIfPresent(context.Background(), key, dimensions)
		if !ok {
			t.Fatal("GetIfPresent returned miss after GetOrLoad")
		}
		if third[0] == first[0] {
			t.Fatalf("cached value aliased returned slice: third[0]=%v first[0]=%v", third[0], first[0])
		}
	})
}

func TestEmbeddingCacheKeySeparatesDimensions_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		input := rapid.String().Draw(t, "input")
		dimA := rapid.IntRange(1, 4096).Draw(t, "dimA")
		dimB := rapid.IntRange(1, 4096).Draw(t, "dimB")

		keyA1 := EmbeddingCacheKey(input, dimA)
		keyA2 := EmbeddingCacheKey(input, dimA)
		if keyA1 != keyA2 {
			t.Fatalf("EmbeddingCacheKey is not deterministic: %q != %q", keyA1, keyA2)
		}
		if dimA != dimB && keyA1 == EmbeddingCacheKey(input, dimB) {
			t.Fatalf("EmbeddingCacheKey collision across dimensions %d and %d for input %q", dimA, dimB, input)
		}
	})
}

package cache

import "testing"

func BenchmarkShortHashOld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = shortHash("this is a very long string that we are hashing repeatedly to test performance")
	}
}

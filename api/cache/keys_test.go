package cache

import (
	"testing"
)

func TestLocalCacheKey(t *testing.T) {
	tests := []struct {
		key  string
		dim  int
		want string
	}{
		{"abc", 1024, "abc:d1024"},
		{"abc", 512, "abc:d512"},
		{"xyz", 768, "xyz:d768"},
		{"", 128, ":d128"},
	}

	for _, tt := range tests {
		got := localCacheKey(tt.key, tt.dim)
		if got != tt.want {
			t.Errorf("localCacheKey(%q, %d) = %q, want %q", tt.key, tt.dim, got, tt.want)
		}
	}
}

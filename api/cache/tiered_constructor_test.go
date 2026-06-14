package cache

import (
	"log/slog"
	"testing"
	"time"
)

func TestNewTieredCacheUsesDefaultLogger(t *testing.T) {
	local := NewLocalCache(1024, time.Minute)
	cache := NewTieredCache(local, nil, time.Minute)
	if cache == nil {
		t.Fatal("NewTieredCache returned nil")
	}
	if cache.local != local {
		t.Fatal("NewTieredCache did not store local cache")
	}
	if cache.localTTL != time.Minute {
		t.Fatalf("localTTL = %s, want 1m", cache.localTTL)
	}
	if cache.logger == nil {
		t.Fatal("logger is nil")
	}
}

func TestNewTieredCacheWithLoggerNilFallsBack(t *testing.T) {
	cache := NewTieredCacheWithLogger(nil, nil, time.Second, nil)
	if cache == nil || cache.logger == nil {
		t.Fatal("NewTieredCacheWithLogger nil logger did not fall back")
	}
}

func TestNewTieredCacheWithLoggerStoresLogger(t *testing.T) {
	logger := slog.Default()
	cache := NewTieredCacheWithLogger(nil, nil, time.Second, logger)
	if cache.logger != logger {
		t.Fatal("NewTieredCacheWithLogger did not store explicit logger")
	}
}

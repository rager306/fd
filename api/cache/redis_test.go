package cache

import (
	"strings"
	"testing"
	"time"
)

func clearRedisCacheEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"REDIS_CACHE_TTL",
		"REDIS_CACHE_NO_EXPIRE",
		"EMBEDDING_CACHE_VERSION",
		"EMBEDDING_MODEL_ID",
		"EMBEDDING_MODEL_REVISION",
		"EMBEDDING_TOKENIZER_VERSION",
		"EMBEDDING_CHUNKING_VERSION",
	} {
		t.Setenv(key, "")
	}
}

func closeRedisCache(t *testing.T, c *RedisCache) {
	t.Helper()
	if err := c.Close(); err != nil {
		t.Fatalf("close cache: %v", err)
	}
}

func TestHashText(t *testing.T) {
	c := &RedisCache{prefix: "test:", namespace: "v2"}

	hash1 := c.HashText("hello")
	hash2 := c.HashText("hello")
	hash3 := c.HashText("world")

	if hash1 != hash2 {
		t.Errorf("same text should produce same hash")
	}
	if hash1 == hash3 {
		t.Errorf("different text should produce different hash")
	}
	if len(hash1) != 64 { // SHA256 hex length
		t.Errorf("expected 64 char hash, got %d", len(hash1))
	}
}

func TestHashText_Deterministic(t *testing.T) {
	c := &RedisCache{prefix: "embed:cache:", namespace: "v2"}
	text := "test text for hashing"

	hash1 := c.HashText(text)
	hash2 := c.HashText(text)

	if hash1 != hash2 {
		t.Errorf("hash not deterministic: %s != %s", hash1, hash2)
	}
}

func TestRedisCacheDefaultOptionsPreserveLegacyKey(t *testing.T) {
	opts := DefaultRedisCacheOptions("embed:cache:", 50)
	c, err := NewRedisCacheWithOptions("127.0.0.1:0", opts)
	if err != nil {
		t.Fatalf("new cache: %v", err)
	}
	defer closeRedisCache(t, c)

	key := c.key("hello", 1024)
	wantPrefix := "embed:cache:v2:"
	if !strings.HasPrefix(key, wantPrefix) {
		t.Fatalf("key %q should preserve prefix %q", key, wantPrefix)
	}
	if !strings.HasSuffix(key, ":d1024") {
		t.Fatalf("key %q should include dimension suffix", key)
	}
	if c.expiration() != 24*time.Hour {
		t.Fatalf("expiration=%s, want 24h", c.expiration())
	}
}

func TestRedisCacheNamespaceIncludesConfiguredFields(t *testing.T) {
	opts := DefaultRedisCacheOptions("embed:cache:", 50)
	opts.Namespace.ModelID = "deepvk/USER-bge-m3"
	opts.Namespace.ModelRevision = "rev-1"
	opts.Namespace.TokenizerVersion = "tok-1"
	opts.Namespace.ChunkingVersion = "chunk-1"

	c, err := NewRedisCacheWithOptions("127.0.0.1:0", opts)
	if err != nil {
		t.Fatalf("new cache: %v", err)
	}
	defer closeRedisCache(t, c)

	key := c.key("hello", 1024)
	if !strings.HasPrefix(key, "embed:cache:v2:m") {
		t.Fatalf("key %q should include model namespace hash", key)
	}
	for _, raw := range []string{"deepvk/USER-bge-m3", "rev-1", "tok-1", "chunk-1"} {
		if strings.Contains(key, raw) {
			t.Fatalf("key %q should not contain raw namespace value %q", key, raw)
		}
	}
	if key == (&RedisCache{prefix: "embed:cache:", namespace: "v2"}).key("hello", 1024) {
		t.Fatalf("configured namespace should change key")
	}
}

func TestRedisCacheOptionsFromEnvDefaults(t *testing.T) {
	clearRedisCacheEnv(t)

	opts, err := RedisCacheOptionsFromEnv("embed:cache:", 50)
	if err != nil {
		t.Fatalf("options: %v", err)
	}
	if opts.TTL != 24*time.Hour {
		t.Fatalf("TTL=%s, want 24h", opts.TTL)
	}
	if opts.NoExpire {
		t.Fatalf("NoExpire should default false")
	}
	if opts.Namespace.CacheVersion != "v2" {
		t.Fatalf("CacheVersion=%q, want v2", opts.Namespace.CacheVersion)
	}
}

func TestRedisCacheOptionsFromEnvTTL(t *testing.T) {
	clearRedisCacheEnv(t)
	t.Setenv("REDIS_CACHE_TTL", "168h")

	opts, err := RedisCacheOptionsFromEnv("embed:cache:", 50)
	if err != nil {
		t.Fatalf("options: %v", err)
	}
	if opts.TTL != 168*time.Hour {
		t.Fatalf("TTL=%s, want 168h", opts.TTL)
	}
}

func TestRedisCacheOptionsFromEnvNoExpire(t *testing.T) {
	clearRedisCacheEnv(t)
	t.Setenv("REDIS_CACHE_NO_EXPIRE", "true")

	opts, err := RedisCacheOptionsFromEnv("embed:cache:", 50)
	if err != nil {
		t.Fatalf("options: %v", err)
	}
	if !opts.NoExpire {
		t.Fatalf("NoExpire should be true")
	}
	c, err := NewRedisCacheWithOptions("127.0.0.1:0", opts)
	if err != nil {
		t.Fatalf("new cache: %v", err)
	}
	defer closeRedisCache(t, c)
	if c.expiration() != 0 {
		t.Fatalf("expiration=%s, want no expiration", c.expiration())
	}
}

func TestRedisCacheOptionsFromEnvRejectsTTLWithNoExpire(t *testing.T) {
	clearRedisCacheEnv(t)
	t.Setenv("REDIS_CACHE_NO_EXPIRE", "true")
	t.Setenv("REDIS_CACHE_TTL", "168h")

	_, err := RedisCacheOptionsFromEnv("embed:cache:", 50)
	if err == nil {
		t.Fatal("expected conflict error")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRedisCacheOptionsFromEnvRejectsInvalidTTL(t *testing.T) {
	clearRedisCacheEnv(t)
	t.Setenv("REDIS_CACHE_TTL", "not-a-duration")

	_, err := RedisCacheOptionsFromEnv("embed:cache:", 50)
	if err == nil {
		t.Fatal("expected invalid TTL error")
	}
}

func TestRedisCacheOptionsFromEnvNamespace(t *testing.T) {
	clearRedisCacheEnv(t)
	t.Setenv("EMBEDDING_CACHE_VERSION", "v3")
	t.Setenv("EMBEDDING_MODEL_ID", "deepvk/USER-bge-m3")
	t.Setenv("EMBEDDING_MODEL_REVISION", "main")
	t.Setenv("EMBEDDING_TOKENIZER_VERSION", "tokenizer-sha")
	t.Setenv("EMBEDDING_CHUNKING_VERSION", "chunk-v1")

	opts, err := RedisCacheOptionsFromEnv("embed:cache:", 50)
	if err != nil {
		t.Fatalf("options: %v", err)
	}
	if opts.Namespace.CacheVersion != "v3" {
		t.Fatalf("CacheVersion=%q, want v3", opts.Namespace.CacheVersion)
	}
	if opts.Namespace.ModelID == "" || opts.Namespace.TokenizerVersion == "" || opts.Namespace.ChunkingVersion == "" {
		t.Fatalf("namespace fields not populated: %+v", opts.Namespace)
	}
}

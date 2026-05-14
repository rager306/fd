package cache

import (
	"testing"
)

func TestHashText(t *testing.T) {
	c := &RedisCache{prefix: "test:"}

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
	c := &RedisCache{prefix: "embed:cache:"}
	text := "test text for hashing"

	hash1 := c.HashText(text)
	hash2 := c.HashText(text)

	if hash1 != hash2 {
		t.Errorf("hash not deterministic: %s != %s", hash1, hash2)
	}
}
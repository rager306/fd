package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// shortHash feeds cache keys and the Redis namespace string, so its output
// format (12 lowercase hex characters) is part of fd's cache identity. A
// change here silently invalidates every cached embedding.
func TestShortHashFormatAndStability(t *testing.T) {
	const want = "8bd574fdb05c" // sha256("fd")[:12]
	got := shortHash("fd")
	if got != want {
		t.Fatalf("shortHash(\"fd\") = %q, want %q", got, want)
	}
	if len(got) != 12 {
		t.Fatalf("shortHash length = %d, want 12", len(got))
	}
	for _, value := range []string{"", "договор поставки", "44-ФЗ статья 33"} {
		h := sha256.Sum256([]byte(value))
		if expanded := hex.EncodeToString(h[:])[:12]; shortHash(value) != expanded {
			t.Errorf("shortHash(%q) = %q, want the first 12 hex chars of its sha256", value, shortHash(value))
		}
	}
}

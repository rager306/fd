package envutil

import (
	"testing"
	"time"
)

func TestDurationOrDefaultFallbackWhenUnset(t *testing.T) {
	t.Setenv("FD_TEST_DUR_UNSET", "")
	d := DurationOrDefault("FD_TEST_DUR_UNSET", 30*time.Second)
	if d != 30*time.Second {
		t.Fatalf("unset fallback: got %s, want 30s", d)
	}
}

func TestDurationOrDefaultValidDurations(t *testing.T) {
	cases := []struct {
		env  string
		want time.Duration
	}{
		{"60s", 60 * time.Second},
		{"5m", 5 * time.Minute},
		{"100ms", 100 * time.Millisecond},
		{"0s", 0},
	}
	for _, c := range cases {
		t.Run(c.env, func(t *testing.T) {
			t.Setenv("FD_TEST_DUR", c.env)
			if got := DurationOrDefault("FD_TEST_DUR", 30*time.Second); got != c.want {
				t.Fatalf("env=%q: got %s, want %s", c.env, got, c.want)
			}
		})
	}
}

func TestDurationOrDefaultInvalidFallsBack(t *testing.T) {
	invalid := []string{"30", "notaduration", "-5s", "abc"}
	for _, s := range invalid {
		t.Run(s, func(t *testing.T) {
			t.Setenv("FD_TEST_DUR", s)
			if got := DurationOrDefault("FD_TEST_DUR", 10*time.Second); got != 10*time.Second {
				t.Fatalf("env=%q: got %s, want fallback 10s", s, got)
			}
		})
	}
}

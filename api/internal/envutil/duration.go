package envutil

import (
	"os"
	"strings"
	"time"
)

// DurationOrDefault returns a Go duration from the environment or fallback.
// Accepts Go duration strings (s, ms, us, ns, m, h). Returns the fallback
// if the env var is unset, empty, or unparseable.
func DurationOrDefault(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}

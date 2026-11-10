// Package envutil centralizes safe environment variable parsing helpers.
package envutil

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Int returns a non-negative integer from the environment or fallback.
func Int(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}

// PositiveInt returns a positive integer from the environment or fallback.
func PositiveInt(key string, fallback int) int {
	parsed := Int(key, fallback)
	if parsed <= 0 {
		return fallback
	}
	return parsed
}

// RequiredPositiveInt returns the fallback only when the variable is unset.
// A configured value must be a positive integer.
func RequiredPositiveInt(key string, fallback int) (int, error) {
	value, present := os.LookupEnv(key)
	if !present {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return parsed, nil
}

// Int64 returns a non-negative 64-bit integer from the environment or fallback.
func Int64(key string, fallback int64) int64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}

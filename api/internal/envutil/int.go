// Package envutil centralizes safe environment variable parsing helpers.
package envutil

import (
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

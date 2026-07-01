package envutil

import (
	"os"
	"strings"
)

// BoolOrDefault returns a bool from the environment or fallback. The value is
// truthy for: "1", "true", "yes", "on", "y", "t" (case-insensitive); falsy for
// "0", "false", "no", "off", "n", "f", "". An unrecognized value returns the
// fallback so misconfiguration cannot silently flip a feature flag.
func BoolOrDefault(key string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch value {
	case "":
		return fallback
	case "1", "true", "yes", "on", "y", "t":
		return true
	case "0", "false", "no", "off", "n", "f":
		return false
	default:
		return fallback
	}
}

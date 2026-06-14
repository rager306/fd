package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	corsMethods = "GET,POST,OPTIONS"
	corsHeaders = "Content-Type,Authorization,X-Request-Id"
)

// CORSFromEnv returns CORS middleware configured by FD_CORS_ORIGINS.
// FD_CORS_ORIGINS accepts a comma-separated allowlist and defaults to "*".
func CORSFromEnv() gin.HandlerFunc {
	return CORS(os.Getenv("FD_CORS_ORIGINS"))
}

// CORS adds fd v2 CORS headers and terminates OPTIONS preflight with 204.
func CORS(origins string) gin.HandlerFunc {
	allowed := parseCORSOrigins(origins)
	return func(c *gin.Context) {
		origin := selectAllowedOrigin(c.GetHeader("Origin"), allowed)
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", corsMethods)
		c.Header("Access-Control-Allow-Headers", corsHeaders)
		c.Header("Vary", "Origin")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func parseCORSOrigins(origins string) []string {
	if strings.TrimSpace(origins) == "" {
		return []string{"*"}
	}
	parts := strings.Split(origins, ",")
	allowed := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin != "" {
			allowed = append(allowed, origin)
		}
	}
	if len(allowed) == 0 {
		return []string{"*"}
	}
	return allowed
}

func selectAllowedOrigin(requestOrigin string, allowed []string) string {
	for _, origin := range allowed {
		if origin == "*" {
			return "*"
		}
		if requestOrigin != "" && origin == requestOrigin {
			return requestOrigin
		}
	}
	return ""
}

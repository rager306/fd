package middleware

import (
	"crypto/subtle"
	"os"
	"strings"

	"fd-api/handlers"

	"github.com/gin-gonic/gin"
)

const (
	bearerPrefix        = "Bearer "
	publicLivePath      = "/live"
	publicReadyPath     = "/ready"
	publicHealthPath    = "/health"
	publicV1Healthcheck = "/v1/healthcheck"
	publicMetrics       = "/metrics"
	publicDocs          = "/docs"
	publicOpenAPI       = "/openapi.json"
)

// APIKeyAuthFromEnv returns fail-closed auth middleware configured from FD_API_KEY.
// When FD_API_KEY is empty, protected endpoints reject requests instead of
// silently disabling authentication. Public probe/docs endpoints stay open.
func APIKeyAuthFromEnv() gin.HandlerFunc {
	return APIKeyAuth(os.Getenv("FD_API_KEY"))
}

// APIKeyAuth requires Authorization: Bearer <apiKey> on protected endpoints.
// Public endpoints are limited to cheap liveness/metadata/docs surfaces.
func APIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isAuthPublicPath(c.Request.URL.Path) || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		if strings.TrimSpace(apiKey) == "" {
			handlers.WriteError(c, handlers.CodeUnauthorized, "authorization", "api key is not configured")
			c.Abort()
			return
		}

		authorization := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, bearerPrefix) {
			handlers.WriteError(c, handlers.CodeUnauthorized, "authorization", "missing bearer token")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authorization, bearerPrefix)
		tokenBytes := []byte(token)
		apiKeyBytes := []byte(apiKey)

		// Mask length mismatch to prevent length-based timing attacks.
		// If lengths differ, we compare the expected key to itself to ensure
		// the timing is always relative to the expected key's length.
		match := 1
		if len(tokenBytes) != len(apiKeyBytes) {
			tokenBytes = apiKeyBytes
			match = 0
		}
		match &= subtle.ConstantTimeCompare(tokenBytes, apiKeyBytes)

		if match != 1 {
			handlers.WriteError(c, handlers.CodeUnauthorized, "authorization", "invalid bearer token")
			c.Abort()
			return
		}
		c.Next()
	}
}

func isAuthPublicPath(path string) bool {
	return path == publicLivePath ||
		path == publicReadyPath ||
		path == publicHealthPath ||
		path == publicV1Healthcheck ||
		path == publicOpenAPI ||
		path == publicDocs ||
		strings.HasPrefix(path, publicDocs+"/")
}

package middleware

import (
	"crypto/subtle"
	"os"
	"strings"

	"fd-api/handlers"

	"github.com/gin-gonic/gin"
)

const (
	bearerPrefix   = "Bearer "
	publicLivePath = "/live"
	publicMetrics  = "/metrics"
	publicDocs     = "/docs"
	publicOpenAPI  = "/openapi.json"
)

// APIKeyAuthFromEnv returns auth middleware configured from FD_API_KEY.
// When FD_API_KEY is empty, authentication is disabled for local/dev parity.
func APIKeyAuthFromEnv() gin.HandlerFunc {
	return APIKeyAuth(os.Getenv("FD_API_KEY"))
}

// APIKeyAuth requires Authorization: Bearer <apiKey> on protected endpoints.
// Public endpoints are limited to cheap liveness/metadata/docs surfaces.
func APIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiKey == "" || isAuthPublicPath(c.Request.URL.Path) || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		authorization := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, bearerPrefix) {
			handlers.WriteError(c, handlers.CodeUnauthorized, "authorization", "missing bearer token")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authorization, bearerPrefix)
		if subtle.ConstantTimeCompare([]byte(token), []byte(apiKey)) != 1 {
			handlers.WriteError(c, handlers.CodeUnauthorized, "authorization", "invalid bearer token")
			c.Abort()
			return
		}
		c.Next()
	}
}

func isAuthPublicPath(path string) bool {
	return path == publicLivePath ||
		path == publicMetrics ||
		path == publicOpenAPI ||
		path == publicDocs ||
		strings.HasPrefix(path, publicDocs+"/")
}

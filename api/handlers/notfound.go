package handlers

import "github.com/gin-gonic/gin"

// NotFoundMiddleware returns a 404 envelope for unknown paths (NoRoute).
// Per docs/fd-v2.md Section 3, "path {path} not found" → code=not_found.
// Replaces gin's default text/plain "404 page not found" body.
func NotFoundMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		WriteError(c, CodeNotFound, "",
			"path "+c.Request.URL.Path+" not found")
	}
}

// MethodNotAllowedMiddleware returns a 405 envelope when the path is
// known but the method isn't (e.g. GET /v1/embeddings which is POST-only).
// Spec T-E-8: GET /v1/embeddings → 405 (NOT 404).
func MethodNotAllowedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		WriteError(c, CodeMethodNotAllowed, "method", "method "+c.Request.Method+" not allowed on "+c.Request.URL.Path)
	}
}

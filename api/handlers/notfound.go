package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
//
// We use 405 with a not_found envelope shape (not its own code) because
// OpenAI's catalog in docs/fd-v2.md Section 3 does not define a
// method_not_allowed code. The HTTP status is the machine-readable
// signal; the body code stays not_found.
func MethodNotAllowedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Status is set explicitly because WriteError maps CodeNotFound → 404.
		// We need 405 here, so we override.
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, ErrorResponse{
			Error: ErrorDetail{
				Code:    "method_not_allowed",
				Type:    "invalid_request_error",
				Param:   "method",
				Message: "method " + c.Request.Method + " not allowed on " + c.Request.URL.Path,
			},
		})
	}
}

package handlers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"log/slog"
)

// RecoveryMiddleware wraps gin.Recovery so unhandled panics return the
// OpenAI-style error envelope (500 internal_error) instead of gin's
// default plain-text "Internal Server Error" response.
//
// Without this, T-E-15 (forced internal error → 500 internal_error with
// X-Request-Id in message) fails. Plain-text 500s also leak stack-trace
// fragments when gin.RecoveryWithWriter is misconfigured.
//
// Usage: r.Use(handlers.RecoveryMiddleware(logger)) BEFORE all other
// middleware so even validation/header errors that panic are caught.
func RecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := debug.Stack()
				logger.Error("panic recovered",
					"panic", fmt.Sprintf("%v", rec),
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"stack", string(stack),
				)
				// Build a stable request_id reference for the response
				// message. If headers middleware (S03) is wired, it will
				// have already set X-Request-Id; we re-use that value so
				// the body's request_id matches the response header.
				requestID := c.Writer.Header().Get("X-Request-Id")
				msg := "internal server error"
				if requestID != "" {
					msg = "internal server error; request_id=" + requestID
				}
				// If WriteError was already called (handler aborted
				// before panic), do not double-write. IsAborted is set
				// by AbortWithStatusJSON.
				if !c.Writer.Written() {
					c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
						Error: ErrorDetail{
							Code:    CodeInternalError,
							Type:    TypeInternalError,
							Message: msg,
						},
					})
				}
			}
		}()
		c.Next()
	}
}

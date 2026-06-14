package handlers

import (
	"net/http"
	"time"

	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
)

const retryAfterWarmupSeconds = "5"

// NewLiveHandler returns a cheap liveness probe. It intentionally does not
// depend on model warmup or downstream services; it only proves the process can
// answer HTTP.
func NewLiveHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			responseStatusKey: "ok",
			responseTimeKey:   time.Now().Format(time.RFC3339),
		})
	}
}

// NewReadyHandler returns a readiness probe backed by lifecycle state.
// It is 200 only after warmup succeeds and shutdown has not begun.
func NewReadyHandler(state *lifecycle.State) gin.HandlerFunc {
	return func(c *gin.Context) {
		if state != nil && state.IsReady() {
			c.JSON(http.StatusOK, gin.H{
				responseStatusKey: "ready",
				responseTimeKey:   time.Now().Format(time.RFC3339),
			})
			return
		}

		c.Header("Retry-After", retryAfterWarmupSeconds)
		WriteError(c, CodeModelNotLoaded, "", "model warmup is not complete")
	}
}

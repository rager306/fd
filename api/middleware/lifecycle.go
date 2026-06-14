package middleware

import (
	"fd-api/handlers"
	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
)

const (
	warmupRetryAfterSeconds   = "5"
	shutdownRetryAfterSeconds = "30"
)

// LifecycleGate rejects embedding requests while the process is warming up or
// shutting down, and tracks accepted requests for graceful drain.
func LifecycleGate(state *lifecycle.State) gin.HandlerFunc {
	return LifecycleGateWithCapacity(state, 0)
}

// LifecycleGateWithCapacity behaves like LifecycleGate and additionally
// rejects requests with model_overloaded when maxInFlight is reached. A
// maxInFlight value <= 0 disables capacity limiting.
func LifecycleGateWithCapacity(state *lifecycle.State, maxInFlight int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if state != nil && state.IsShuttingDown() {
			c.Header("Retry-After", shutdownRetryAfterSeconds)
			handlers.WriteError(c, handlers.CodeShuttingDown, "", "server is shutting down")
			return
		}
		if state == nil || !state.IsReady() {
			c.Header("Retry-After", warmupRetryAfterSeconds)
			handlers.WriteError(c, handlers.CodeModelNotLoaded, "", "model warmup is not complete")
			return
		}

		done, ok := state.TryTrackRequest(maxInFlight)
		if !ok {
			c.Header("Retry-After", warmupRetryAfterSeconds)
			handlers.WriteError(c, handlers.CodeModelOverloaded, "", "model capacity is exhausted")
			return
		}
		defer done()
		c.Request = c.Request.WithContext(lifecycle.WithState(c.Request.Context(), state))
		c.Next()
	}
}

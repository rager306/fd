package handlers

import (
	"net/http"
	"time"

	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
)

const (
	responseStatusKey = "status"
	responseTimeKey   = "time"

	healthStatusOK       = "ok"
	healthStatusDegraded = "degraded"
	healthStatusDown     = "down"
)

// RuntimeHealth describes the active embedding runtime reported by /health.
// It is metadata only; readiness still requires a smoke embedding request.
type RuntimeHealth struct {
	Backend           string `json:"backend"`
	Model             string `json:"model,omitempty"`
	Dimensions        int    `json:"dimensions,omitempty"`
	ProductionDefault bool   `json:"production_default"`
	CacheNamespace    string `json:"cache_namespace,omitempty"`
}

// DeepHealthResponse is the wire shape returned by GET /health.
type DeepHealthResponse struct {
	Status           string         `json:"status"`
	Time             string         `json:"time"`
	ModelLoaded      bool           `json:"model_loaded"`
	WarmupDone       bool           `json:"warmup_done"`
	Device           string         `json:"device"`
	LastInferenceAt  *string        `json:"last_inference_at"`
	InFlightRequests int64          `json:"in_flight_requests"`
	Runtime          *RuntimeHealth `json:"runtime,omitempty"`
}

// HealthHandler serves the default /health response without runtime metadata.
func HealthHandler(c *gin.Context) {
	writeHealth(c, nil, nil)
}

// NewHealthHandler returns a /health handler that includes runtime metadata
// when the embedding backend exposes it.
func NewHealthHandler(runtime *RuntimeHealth) gin.HandlerFunc {
	return NewHealthHandlerWithState(runtime, nil)
}

// NewHealthHandlerWithState returns a deep /health handler backed by lifecycle
// state. It reports ok for ready processes, degraded while warming up, and down
// once shutdown begins.
func NewHealthHandlerWithState(runtime *RuntimeHealth, state *lifecycle.State) gin.HandlerFunc {
	return func(c *gin.Context) {
		writeHealth(c, runtime, state)
	}
}

func writeHealth(c *gin.Context, runtime *RuntimeHealth, state *lifecycle.State) {
	status, httpStatus := healthStatus(state)
	response := DeepHealthResponse{
		Status:           status,
		Time:             time.Now().Format(time.RFC3339),
		Device:           defaultModelDevice,
		Runtime:          runtime,
		InFlightRequests: 0,
	}
	if state != nil {
		response.ModelLoaded = state.IsReady()
		response.WarmupDone = state.IsWarmupDone()
		response.InFlightRequests = state.InFlightCount()
		if lastInferenceAt, ok := state.LastInferenceAt(); ok {
			formatted := lastInferenceAt.Format(time.RFC3339)
			response.LastInferenceAt = &formatted
		}
	}
	c.JSON(httpStatus, response)
}

func healthStatus(state *lifecycle.State) (status string, httpStatus int) {
	if state == nil || state.IsReady() {
		return healthStatusOK, http.StatusOK
	}
	if state.IsShuttingDown() {
		return healthStatusDown, http.StatusServiceUnavailable
	}
	return healthStatusDegraded, http.StatusServiceUnavailable
}

package handlers

import (
	"context"
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

// HealthError describes the latest lifecycle error visible in /health.
type HealthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	At      string `json:"at"`
}

// DependencyStatus describes lightweight reachability for a runtime dependency.
type DependencyStatus struct {
	Reachable bool    `json:"reachable"`
	LatencyMS float64 `json:"latency_ms"`
	Namespace string  `json:"namespace,omitempty"`
	Error     string  `json:"error,omitempty"`
}

// DependencyProbe checks one dependency without performing embedding inference.
type DependencyProbe interface {
	Check(context.Context) DependencyStatus
}

// DependencyProbeFunc adapts a function into a DependencyProbe.
type DependencyProbeFunc func(context.Context) DependencyStatus

// Check implements DependencyProbe.
func (f DependencyProbeFunc) Check(ctx context.Context) DependencyStatus {
	return f(ctx)
}

// DependencyChecks groups optional dependency probes reported by /health.
type DependencyChecks struct {
	TEI   DependencyProbe
	Redis DependencyProbe
}

// HealthDependencies is the JSON dependency block returned by /health.
type HealthDependencies struct {
	TEI   *DependencyStatus `json:"tei,omitempty"`
	Redis *DependencyStatus `json:"redis,omitempty"`
}

// HealthOptions controls optional diagnostic fields for /health.
type HealthOptions struct {
	InFlightCapacity  *int64
	Dependencies      *DependencyChecks
	DependencyTimeout time.Duration
}

// DeepHealthResponse is the wire shape returned by GET /health.
type DeepHealthResponse struct {
	Status           string              `json:"status"`
	Time             string              `json:"time"`
	ModelLoaded      bool                `json:"model_loaded"`
	WarmupDone       bool                `json:"warmup_done"`
	Device           string              `json:"device"`
	LastInferenceAt  *string             `json:"last_inference_at"`
	LastError        *HealthError        `json:"last_error,omitempty"`
	InFlightRequests int64               `json:"in_flight_requests"`
	InFlightCapacity *int64              `json:"in_flight_capacity,omitempty"`
	Runtime          *RuntimeHealth      `json:"runtime,omitempty"`
	Dependencies     *HealthDependencies `json:"dependencies,omitempty"`
}

// HealthHandler serves the default /health response without runtime metadata.
func HealthHandler(c *gin.Context) {
	writeHealth(c, nil, nil, HealthOptions{})
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
	return NewHealthHandlerWithOptions(runtime, state, HealthOptions{})
}

// NewHealthHandlerWithOptions returns a deep /health handler with optional
// agent-facing diagnostics such as capacity and dependency reachability.
func NewHealthHandlerWithOptions(runtime *RuntimeHealth, state *lifecycle.State, opts HealthOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		writeHealth(c, runtime, state, opts)
	}
}

func writeHealth(c *gin.Context, runtime *RuntimeHealth, state *lifecycle.State, opts HealthOptions) {
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
		if err := state.LastError(); err != nil {
			at := time.Now()
			if recordedAt, ok := state.LastErrorAt(); ok {
				at = recordedAt
			}
			response.LastError = &HealthError{Code: "lifecycle_error", Message: err.Error(), At: at.Format(time.RFC3339)}
		}
		if lastInferenceAt, ok := state.LastInferenceAt(); ok {
			formatted := lastInferenceAt.Format(time.RFC3339)
			response.LastInferenceAt = &formatted
		}
	}
	if opts.InFlightCapacity != nil {
		capacity := *opts.InFlightCapacity
		response.InFlightCapacity = &capacity
	}
	if opts.Dependencies != nil {
		response.Dependencies = opts.Dependencies.Check(c.Request.Context(), opts.DependencyTimeout)
	}
	c.JSON(httpStatus, response)
}

// Check runs configured dependency probes and returns the dependency block.
func (checks *DependencyChecks) Check(ctx context.Context, timeout time.Duration) *HealthDependencies {
	if checks == nil {
		return nil
	}
	if timeout <= 0 {
		timeout = 750 * time.Millisecond
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	deps := &HealthDependencies{}
	if checks.TEI != nil {
		status := checks.TEI.Check(ctx)
		deps.TEI = &status
	}
	if checks.Redis != nil {
		status := checks.Redis.Check(ctx)
		deps.Redis = &status
	}
	if deps.TEI == nil && deps.Redis == nil {
		return nil
	}
	return deps
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

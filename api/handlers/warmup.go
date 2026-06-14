package handlers

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
)

const (
	warmupStatusReady     = "ready"
	warmupStatusWarmingUp = "warming_up"
)

// WarmupHandler serves manual warmup status and trigger endpoints.
type WarmupHandler struct {
	state      *lifecycle.State
	model      lifecycle.WarmupModel
	timeout    time.Duration
	inProgress atomic.Bool
}

// WarmupResponse is the wire shape returned by /warmup endpoints.
type WarmupResponse struct {
	Status   string  `json:"status"`
	Progress float64 `json:"progress"`
	Message  string  `json:"message,omitempty"`
}

// NewWarmupHandler wires lifecycle state and model warmup behavior.
func NewWarmupHandler(state *lifecycle.State, model lifecycle.WarmupModel, timeout time.Duration) *WarmupHandler {
	return &WarmupHandler{state: state, model: model, timeout: timeout}
}

// Status returns current warmup readiness state.
func (h *WarmupHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, h.response(""))
}

// Trigger starts background warmup when the model is not ready yet.
func (h *WarmupHandler) Trigger(c *gin.Context) {
	if h.ready() {
		c.JSON(http.StatusOK, h.response("already warm"))
		return
	}
	if h.inProgress.CompareAndSwap(false, true) {
		go h.runWarmup()
	}
	c.JSON(http.StatusAccepted, h.response("warmup started"))
}

func (h *WarmupHandler) runWarmup() {
	defer h.inProgress.Store(false)
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()
	if err := lifecycle.PreWarm(ctx, h.model); err != nil {
		h.state.SetLastError(err)
		return
	}
	h.state.MarkWarmupDone()
}

func (h *WarmupHandler) response(message string) WarmupResponse {
	status := warmupStatusWarmingUp
	progress := 0.0
	if h.ready() {
		status = warmupStatusReady
		progress = 1.0
	} else if h.inProgress.Load() {
		progress = 0.5
	}
	return WarmupResponse{Status: status, Progress: progress, Message: message}
}

func (h *WarmupHandler) ready() bool {
	return h.state != nil && h.state.IsWarmupDone() && h.state.LastError() == nil
}

package handlers

import (
	"net/http"
	"time"

	"fd-api/buildinfo"
	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
)

const (
	defaultModelMaxInputLengthTokens = 512
	defaultModelMaxBatchSize         = 32
	defaultModelDevice               = "cpu"
)

// VersionResponse is the wire shape returned by GET /version.
type VersionResponse struct {
	Service       string `json:"service"`
	Version       string `json:"version"`
	Model         string `json:"model,omitempty"`
	ModelVersion  string `json:"model_version,omitempty"`
	BuildHash     string `json:"build_hash"`
	BuildDate     string `json:"build_date"`
	StartedAt     string `json:"started_at"`
	Uptime        string `json:"uptime"`
	UptimeSeconds int64  `json:"uptime_seconds"`
}

// ModelInfo describes one embedding model exposed by GET /info.
type ModelInfo struct {
	ID                   string `json:"id"`
	Dimensions           []int  `json:"dims"`
	MaxInputLengthTokens int    `json:"max_input_length_tokens"`
	MaxBatchSize         int    `json:"max_batch_size"`
	Loaded               bool   `json:"loaded"`
	WarmupDone           bool   `json:"warmup_done"`
	Device               string `json:"device"`
}

// InfoResponse is the wire shape returned by GET /info.
type InfoResponse struct {
	Service string      `json:"service"`
	Version string      `json:"version"`
	Models  []ModelInfo `json:"models"`
}

// NewVersionHandler returns build and process uptime metadata.
func NewVersionHandler(info buildinfo.Info) gin.HandlerFunc {
	info = buildinfo.New(info)
	return func(c *gin.Context) {
		uptime := info.Uptime()
		c.JSON(http.StatusOK, VersionResponse{
			Service:       info.Service,
			Version:       info.Version,
			Model:         info.Model,
			ModelVersion:  info.ModelVersion,
			BuildHash:     info.BuildHash,
			BuildDate:     info.BuildDate,
			StartedAt:     info.StartedAt.Format(time.RFC3339),
			Uptime:        uptime.String(),
			UptimeSeconds: int64(uptime.Seconds()),
		})
	}
}

// NewInfoHandler returns model metadata and lifecycle loading state.
func NewInfoHandler(info buildinfo.Info, runtime *RuntimeHealth, state *lifecycle.State) gin.HandlerFunc {
	info = buildinfo.New(info)
	return func(c *gin.Context) {
		model := modelInfo(runtime, state, info.Model)
		c.JSON(http.StatusOK, InfoResponse{
			Service: info.Service,
			Version: info.Version,
			Models:  []ModelInfo{model},
		})
	}
}

func modelInfo(runtime *RuntimeHealth, state *lifecycle.State, fallbackModel string) ModelInfo {
	modelID := fallbackModel
	if runtime != nil && runtime.Model != "" {
		modelID = runtime.Model
	}
	warmupDone := false
	loaded := false
	if state != nil {
		warmupDone = state.IsWarmupDone()
		loaded = state.IsReady()
	}
	return ModelInfo{
		ID:                   modelID,
		Dimensions:           []int{512, 1024},
		MaxInputLengthTokens: defaultModelMaxInputLengthTokens,
		MaxBatchSize:         defaultModelMaxBatchSize,
		Loaded:               loaded,
		WarmupDone:           warmupDone,
		Device:               defaultModelDevice,
	}
}

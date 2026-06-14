package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RuntimeHealth describes the active embedding runtime reported by /health.
// It is metadata only; readiness still requires a smoke embedding request.
type RuntimeHealth struct {
	Backend                    string `json:"backend"`
	Model                      string `json:"model,omitempty"`
	ArtifactID                 string `json:"artifact_id,omitempty"`
	Dimensions                 int    `json:"dimensions,omitempty"`
	MaxSequenceLength          int    `json:"max_sequence_length,omitempty"`
	ValidatedMaxSequenceLength int    `json:"validated_max_sequence_length,omitempty"`
	ProductionDefault          bool   `json:"production_default"`
	// Pointer bools omitted from JSON when nil (TEI path); set for ONNX only.
	ArtifactVerified       *bool  `json:"artifact_verified,omitempty"`
	TokenizerVerified      *bool  `json:"tokenizer_verified,omitempty"`
	RuntimeLibraryVerified *bool  `json:"runtime_library_verified,omitempty"`
	Provider               string `json:"provider,omitempty"`
	CacheNamespace         string `json:"cache_namespace,omitempty"`
}

// HealthHandler serves the basic /health response without runtime metadata.
func HealthHandler(c *gin.Context) {
	writeHealth(c, nil)
}

// NewHealthHandler returns a /health handler that includes runtime metadata
// when the embedding backend exposes it.
func NewHealthHandler(runtime *RuntimeHealth) gin.HandlerFunc {
	return func(c *gin.Context) {
		writeHealth(c, runtime)
	}
}

func writeHealth(c *gin.Context, runtime *RuntimeHealth) {
	body := gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	}
	if runtime != nil {
		body["runtime"] = runtime
	}
	c.JSON(http.StatusOK, body)
}

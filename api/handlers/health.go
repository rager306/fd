package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RuntimeHealth struct {
	Backend                    string `json:"backend"`
	Model                      string `json:"model,omitempty"`
	ArtifactID                 string `json:"artifact_id,omitempty"`
	Dimensions                 int    `json:"dimensions,omitempty"`
	MaxSequenceLength          int    `json:"max_sequence_length,omitempty"`
	ValidatedMaxSequenceLength int    `json:"validated_max_sequence_length,omitempty"`
	ProductionDefault          bool   `json:"production_default"`
	ArtifactVerified           bool   `json:"artifact_verified"`
	TokenizerVerified          bool   `json:"tokenizer_verified"`
	RuntimeLibraryVerified     bool   `json:"runtime_library_verified"`
	Provider                   string `json:"provider,omitempty"`
	CacheNamespace             string `json:"cache_namespace,omitempty"`
}

func HealthHandler(c *gin.Context) {
	writeHealth(c, nil)
}

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

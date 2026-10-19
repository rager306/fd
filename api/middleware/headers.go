package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"

	"fd-api/buildinfo"
	"fd-api/embed"
	"fd-api/handlers"

	"github.com/gin-gonic/gin"
)

const (
	// HeaderRequestID is the fd correlation-id response/request header.
	HeaderRequestID = "X-Request-Id"
	// HeaderModelID reports the authoritative embedding model id on responses.
	HeaderModelID = "X-Model-Id"
	// HeaderDimensions reports the embedding dimensionality used for a response.
	HeaderDimensions = "X-Dimensions"
)

type headerWriter struct {
	gin.ResponseWriter
	context *gin.Context
	modelID string
}

// HeadersMiddleware sets stable fd response headers and injects embedding
// metadata headers before the first response write.
func HeadersMiddleware(info buildinfo.Info, modelID string) gin.HandlerFunc {
	info = buildinfo.New(info)
	serverValue := "fd/" + info.Version
	return func(c *gin.Context) {
		requestID := c.GetHeader(HeaderRequestID)
		if requestID == "" {
			requestID = newRequestID()
		}
		header := c.Writer.Header()
		header.Set("Server", serverValue)
		header.Set(HeaderRequestID, requestID)
		header.Set("Connection", "keep-alive")
		c.Writer = &headerWriter{ResponseWriter: c.Writer, context: c, modelID: modelID}
		c.Next()
	}
}

func (w *headerWriter) WriteHeader(code int) {
	w.setEmbeddingHeaders()
	w.ResponseWriter.WriteHeader(code)
}

func (w *headerWriter) WriteHeaderNow() {
	w.setEmbeddingHeaders()
	w.ResponseWriter.WriteHeaderNow()
}

func (w *headerWriter) Write(data []byte) (int, error) {
	w.setEmbeddingHeaders()
	return w.ResponseWriter.Write(data)
}

func (w *headerWriter) WriteString(data string) (int, error) {
	w.setEmbeddingHeaders()
	return w.ResponseWriter.WriteString(data)
}

func (w *headerWriter) setEmbeddingHeaders() {
	if w.context.FullPath() != "/v1/embeddings" {
		return
	}
	value, ok := w.context.Get(handlers.ContextKeyValidatedRequest)
	if !ok {
		return
	}
	req, ok := value.(*embed.EmbeddingsRequest)
	if !ok {
		return
	}
	dimensions := 1024
	if req.Dimensions != nil {
		dimensions = *req.Dimensions
	}
	header := w.Header()
	header.Set(HeaderModelID, w.modelID)

	// Fast path: bypassing strconv.Itoa for common dimension sizes
	// saves heap allocations for values outside the 0-99 cache.
	switch dimensions {
	case 1024:
		header.Set(HeaderDimensions, "1024")
	case 512:
		header.Set(HeaderDimensions, "512")
	default:
		header.Set(HeaderDimensions, strconv.Itoa(dimensions))
	}
}

func newRequestID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		panic("headers: crypto/rand unavailable: " + err.Error())
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	encoded := hex.EncodeToString(bytes[:])
	return encoded[0:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:32]
}

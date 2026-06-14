package observability

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const defaultTraceCapacity = 100

// TraceEntry is one request trace exposed by GET /v1/traces.
type TraceEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	LatencyMS  int64     `json:"latency_ms"`
	Status     int       `json:"status"`
	ModelID    string    `json:"model_id"`
	RequestID  string    `json:"request_id"`
	Path       string    `json:"path"`
	Dimensions int       `json:"dimensions"`
}

// TraceStore keeps a concurrency-safe ring buffer of recent request traces.
type TraceStore struct {
	mu       sync.Mutex
	entries  []TraceEntry
	next     int
	filled   bool
	capacity int
	enabled  bool
}

// NewTraceStore creates a trace ring buffer. Non-positive capacity uses 100.
func NewTraceStore(capacity int, enabled bool) *TraceStore {
	if capacity <= 0 {
		capacity = defaultTraceCapacity
	}
	return &TraceStore{entries: make([]TraceEntry, capacity), capacity: capacity, enabled: enabled}
}

// NewTraceStoreFromEnv creates the default store. FD_TRACES_ENABLED defaults true.
func NewTraceStoreFromEnv() *TraceStore {
	return NewTraceStore(defaultTraceCapacity, tracesEnabledFromEnv())
}

// Middleware records completed requests into the trace ring.
func (s *TraceStore) Middleware(modelID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.enabled || c.Request.URL.Path == "/v1/traces" {
			c.Next()
			return
		}
		started := time.Now()
		c.Next()
		s.Add(TraceEntry{
			Timestamp:  time.Now().UTC(),
			LatencyMS:  time.Since(started).Milliseconds(),
			Status:     c.Writer.Status(),
			ModelID:    firstNonEmpty(c.Writer.Header().Get("X-Model-Id"), modelID),
			RequestID:  c.Writer.Header().Get("X-Request-Id"),
			Path:       c.Request.URL.Path,
			Dimensions: parseDimensionHeader(c.Writer.Header().Get("X-Dimensions")),
		})
	}
}

// Handler returns the current trace entries in chronological order.
func (s *TraceStore) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, s.Snapshot())
	}
}

// Add appends one trace to the ring buffer.
func (s *TraceStore) Add(entry TraceEntry) {
	if !s.enabled {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[s.next] = entry
	s.next = (s.next + 1) % s.capacity
	if s.next == 0 {
		s.filled = true
	}
}

// Snapshot returns trace entries oldest-first.
func (s *TraceStore) Snapshot() []TraceEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.filled {
		out := make([]TraceEntry, s.next)
		copy(out, s.entries[:s.next])
		return out
	}
	out := make([]TraceEntry, 0, s.capacity)
	out = append(out, s.entries[s.next:]...)
	out = append(out, s.entries[:s.next]...)
	return out
}

func tracesEnabledFromEnv() bool {
	value := strings.TrimSpace(os.Getenv("FD_TRACES_ENABLED"))
	return value == "" || strings.EqualFold(value, "true")
}

func parseDimensionHeader(value string) int {
	dimensions, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return dimensions
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

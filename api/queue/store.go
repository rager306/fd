package queue

import (
	"errors"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"fd-api/embed"
)

// defaultResultTTL is how long a completed result stays in the store
// before the eviction goroutine prunes it. Operators can override this
// through FD_QUEUE_RESULT_TTL (not wired in this version; default applied).
const (
	defaultResultTTL        = 5 * time.Minute
	DefaultResultMaxBytes   = 128 << 20
	DefaultResultMaxEntries = 4096
	maxResultErrorBytes     = 1024
)

// ResultStore provides in-memory storage for queue results, indexed by
// request ID. A background eviction loop runs at most once per tick to
// prune expired entries. The store is goroutine-safe.
type ResultStore struct {
	mu                  sync.Mutex
	data                map[string]*storedResult
	maxBytes, usedBytes int
	maxEntries          int
	closed              bool
	done                chan struct{}
}

type storedResult struct {
	res                *Result
	ts                 time.Time
	bytes, count, dims int
}

// NewResultStore creates a store with TTL-based eviction. Call Close
// to stop the eviction goroutine before discarding the store.
func NewResultStore() *ResultStore {
	return NewResultStoreWithLimits(DefaultResultMaxBytes, DefaultResultMaxEntries)
}

// NewResultStoreWithLimits sets positive limits for reserved vector payload and entries.
func NewResultStoreWithLimits(maxBytes, maxEntries int) *ResultStore {
	if maxBytes <= 0 || maxEntries <= 0 {
		panic("queue result limits must be positive")
	}
	s := &ResultStore{
		data:       make(map[string]*storedResult),
		done:       make(chan struct{}),
		maxBytes:   maxBytes,
		maxEntries: maxEntries,
	}
	go s.evictLoop()
	return s
}

// ReservePending atomically reserves the maximum vector payload before enqueue.
// Pending entries remain reserved until completed or explicitly removed.
func (s *ResultStore) ReservePending(id string, createdAt int64, count, dims int) bool {
	if len(id) != 16 || count < 1 || count > 128 || (dims != 512 && dims != 1024) {
		return false
	}
	bytes := count * dims * 4
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return false
	}
	if len(s.data) >= s.maxEntries || bytes > s.maxBytes-s.usedBytes {
		s.evictLocked(time.Now())
		if len(s.data) >= s.maxEntries || bytes > s.maxBytes-s.usedBytes {
			return false
		}
	}
	if _, exists := s.data[id]; exists {
		return false
	}
	id = strings.Clone(id)
	s.data[id] = &storedResult{res: &Result{ID: id, Status: StatusPending, CreatedAt: createdAt}, bytes: bytes, count: count, dims: dims}
	s.usedBytes += bytes
	return true
}

// DeletePending releases a reservation after an unsuccessful enqueue.
func (s *ResultStore) DeletePending(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.data[id]; ok && entry.res.Status == StatusPending {
		s.removeLocked(id, entry)
	}
}

// CompleteExisting publishes a compact, bounded terminal result for a reservation.
func (s *ResultStore) CompleteExisting(id string, res *Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data[id]
	if s.closed || !ok || entry.res.Status != StatusPending || res == nil {
		return
	}
	result := Result{ID: entry.res.ID, Status: res.Status, CreatedAt: entry.res.CreatedAt, CompletedAt: res.CompletedAt}
	if res.Status == StatusCompleted {
		if validResult(res, entry) {
			result.Dims = entry.dims
			result.EncodingFormat = strings.Clone(res.EncodingFormat)
			result.PromptTokens = res.PromptTokens
			result.Embeddings = make([][]float32, entry.count)
			for i, vector := range res.Embeddings {
				result.Embeddings[i] = make([]float32, entry.dims)
				copy(result.Embeddings[i], vector)
			}
		} else {
			result.Status = StatusFailed
			result.Err = errors.New("invalid embedding dimensions or count")
		}
	} else {
		result.Status = StatusFailed
		message := "unknown"
		if res.Err != nil {
			message = res.Err.Error()
		}
		message = strings.ToValidUTF8(message, "?")
		if len(message) > maxResultErrorBytes {
			message = message[:maxResultErrorBytes]
			for !utf8.ValidString(message) {
				message = message[:len(message)-1]
			}
		}
		result.Err = errors.New(strings.Clone(message))
	}
	entry.res = &result
	entry.ts = time.Now()
}

func validResult(res *Result, entry *storedResult) bool {
	if len(res.Embeddings) != entry.count || res.Dims != entry.dims ||
		(res.EncodingFormat != "" && res.EncodingFormat != embed.EncodingFormatFloat && res.EncodingFormat != embed.EncodingFormatBase64) {
		return false
	}
	for _, vector := range res.Embeddings {
		if len(vector) != entry.dims {
			return false
		}
	}
	return true
}

// Get retrieves a result by id. Returns nil, false if absent or expired.
func (s *ResultStore) Get(id string) (*Result, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sr, ok := s.data[id]
	if !ok {
		return nil, false
	}
	if sr.res.Status != StatusPending && time.Since(sr.ts) > defaultResultTTL {
		s.removeLocked(id, sr)
		return nil, false
	}
	return sr.res, true
}

// Size returns the current number of entries in the store. Exported for
// observability.
func (s *ResultStore) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.data)
}

// Close stops the background eviction loop. Safe to call multiple times.
func (s *ResultStore) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	s.mu.Unlock()
	close(s.done)
	return nil
}

func (s *ResultStore) evictLoop() {
	ticker := time.NewTicker(defaultResultTTL / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.evict()
		case <-s.done:
			return
		}
	}
}

func (s *ResultStore) evict() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evictLocked(time.Now())
}

func (s *ResultStore) evictLocked(now time.Time) {
	for id, sr := range s.data {
		if sr.res.Status != StatusPending && now.Sub(sr.ts) > defaultResultTTL {
			s.removeLocked(id, sr)
		}
	}
}

func (s *ResultStore) removeLocked(id string, sr *storedResult) {
	delete(s.data, id)
	s.usedBytes -= sr.bytes
}

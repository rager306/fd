package queue

import (
	"sync"
	"time"
)

// defaultResultTTL is how long a completed result stays in the store
// before the eviction goroutine prunes it. Operators can override this
// through FD_QUEUE_RESULT_TTL (not wired in this version; default applied).
const defaultResultTTL = 5 * time.Minute

// ResultStore provides in-memory storage for queue results, indexed by
// request ID. A background eviction loop runs at most once per tick to
// prune expired entries. The store is goroutine-safe.
type ResultStore struct {
	mu     sync.Mutex
	data   map[string]*storedResult
	closed bool
	done   chan struct{}
}

type storedResult struct {
	res *Result
	ts  time.Time
}

// NewResultStore creates a store with TTL-based eviction. Call Close
// to stop the eviction goroutine before discarding the store.
func NewResultStore() *ResultStore {
	s := &ResultStore{
		data: make(map[string]*storedResult),
		done: make(chan struct{}),
	}
	go s.evictLoop()
	return s
}

// Save puts a queue result into the store under id. If id already exists
// the previous entry is overwritten; callers should ensure unique IDs.
func (s *ResultStore) Save(id string, res *Result) {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	s.data[id] = &storedResult{res: res, ts: time.Now()}
	s.mu.Unlock()
}

// Get retrieves a result by id. Returns nil, false if absent or expired.
func (s *ResultStore) Get(id string) (*Result, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sr, ok := s.data[id]
	if !ok {
		return nil, false
	}
	if time.Since(sr.ts) > defaultResultTTL {
		delete(s.data, id)
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
	now := time.Now()
	for id, sr := range s.data {
		if now.Sub(sr.ts) > defaultResultTTL {
			delete(s.data, id)
		}
	}
}

// Package lifecycle provides process lifecycle state used by readiness,
// shutdown, and request-gating middleware.
package lifecycle

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// ErrDrainTimeout is returned when in-flight requests do not drain before
// the caller-provided timeout.
var ErrDrainTimeout = errors.New("lifecycle drain timeout")

// State tracks fd process lifecycle: warmup completion, shutdown mode,
// in-flight request count, and the last warmup/runtime error.
type State struct {
	warmupDone      atomic.Bool
	shuttingDown    atomic.Bool
	inflight        sync.WaitGroup
	inflightCount   atomic.Int64
	lastInferenceAt atomic.Value // stores time.Time
	lastError       atomic.Value // stores errorSnapshot
}

type errorSnapshot struct {
	err error
}

type contextKey struct{}

// NewState returns a zero-value lifecycle state. A new state starts unready,
// not shutting down, with no in-flight requests and no recorded error.
func NewState() *State {
	return &State{}
}

// WithState stores state in ctx for handlers and middleware that need lifecycle
// decisions without depending on globals. A nil state leaves ctx unchanged.
func WithState(ctx context.Context, state *State) context.Context {
	if state == nil {
		return ctx
	}
	return context.WithValue(ctx, contextKey{}, state)
}

// FromContext returns the lifecycle state previously stored by WithState.
func FromContext(ctx context.Context) (*State, bool) {
	state, ok := ctx.Value(contextKey{}).(*State)
	return state, ok
}

// MarkWarmupDone marks model warmup as completed and clears any previous error.
func (s *State) MarkWarmupDone() {
	s.SetLastError(nil)
	s.warmupDone.Store(true)
}

// IsWarmupDone reports whether model warmup has completed successfully.
func (s *State) IsWarmupDone() bool {
	return s.warmupDone.Load()
}

// IsReady reports whether fd should accept embedding requests.
func (s *State) IsReady() bool {
	return s.IsWarmupDone() && !s.shuttingDown.Load() && s.LastError() == nil
}

// BeginShutdown marks fd as shutting down. New embedding requests should be
// rejected while existing tracked requests drain.
func (s *State) BeginShutdown() {
	s.shuttingDown.Store(true)
}

// IsShuttingDown reports whether shutdown has begun.
func (s *State) IsShuttingDown() bool {
	return s.shuttingDown.Load()
}

// TrackRequest increments the in-flight request counter and returns an
// idempotent done function that decrements it. Callers should defer the
// returned function after accepting a request.
func (s *State) TrackRequest() func() {
	s.inflight.Add(1)
	s.inflightCount.Add(1)
	return s.doneRequestOnce()
}

// TryTrackRequest tracks a request only when maxInFlight is not already
// reached. A maxInFlight value <= 0 disables capacity limiting.
func (s *State) TryTrackRequest(maxInFlight int64) (func(), bool) {
	if maxInFlight <= 0 {
		return s.TrackRequest(), true
	}
	for {
		current := s.inflightCount.Load()
		if current >= maxInFlight {
			return nil, false
		}
		if s.inflightCount.CompareAndSwap(current, current+1) {
			s.inflight.Add(1)
			return s.doneRequestOnce(), true
		}
	}
}

func (s *State) doneRequestOnce() func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			s.inflightCount.Add(-1)
			s.inflight.Done()
		})
	}
}

// InFlightCount returns the current number of tracked in-flight requests.
func (s *State) InFlightCount() int64 {
	return s.inflightCount.Load()
}

// MarkInferenceSuccess records the timestamp of the latest successful
// embedding response.
func (s *State) MarkInferenceSuccess() {
	s.lastInferenceAt.Store(time.Now())
}

// LastInferenceAt returns the latest successful embedding response timestamp.
func (s *State) LastInferenceAt() (time.Time, bool) {
	value := s.lastInferenceAt.Load()
	if value == nil {
		return time.Time{}, false
	}
	return value.(time.Time), true
}

// WaitDrain waits until all tracked requests finish or timeout elapses.
// A timeout <= 0 performs a non-blocking check.
func (s *State) WaitDrain(timeout time.Duration) error {
	drained := make(chan struct{})
	go func() {
		s.inflight.Wait()
		close(drained)
	}()

	if timeout <= 0 {
		if s.inflightCount.Load() == 0 {
			return nil
		}
		return ErrDrainTimeout
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-drained:
		return nil
	case <-timer.C:
		return ErrDrainTimeout
	}
}

// SetLastError records the last lifecycle error. Passing nil clears it.
func (s *State) SetLastError(err error) {
	s.lastError.Store(errorSnapshot{err: err})
}

// LastError returns the last lifecycle error recorded by SetLastError.
func (s *State) LastError() error {
	value := s.lastError.Load()
	if value == nil {
		return nil
	}
	return value.(errorSnapshot).err
}

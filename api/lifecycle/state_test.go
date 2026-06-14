package lifecycle

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestStateReadinessAndShutdown(t *testing.T) {
	state := NewState()

	if state.IsReady() {
		t.Fatal("new state should not be ready before warmup")
	}
	if state.IsShuttingDown() {
		t.Fatal("new state should not be shutting down")
	}

	state.MarkWarmupDone()
	if !state.IsReady() {
		t.Fatal("state should be ready after warmup")
	}

	state.BeginShutdown()
	if !state.IsShuttingDown() {
		t.Fatal("state should report shutdown after BeginShutdown")
	}
	if state.IsReady() {
		t.Fatal("state should not be ready while shutting down")
	}
}

func TestStateLastErrorAffectsReadiness(t *testing.T) {
	state := NewState()
	state.MarkWarmupDone()

	boom := errors.New("boom")
	state.SetLastError(boom)
	if !errors.Is(state.LastError(), boom) {
		t.Fatalf("LastError() = %v, want %v", state.LastError(), boom)
	}
	if state.IsReady() {
		t.Fatal("state with last error should not be ready")
	}

	state.SetLastError(nil)
	if state.LastError() != nil {
		t.Fatalf("LastError() after clear = %v, want nil", state.LastError())
	}
	if !state.IsReady() {
		t.Fatal("state should be ready after clearing last error")
	}
}

func TestWaitDrainEmptyReturnsImmediately(t *testing.T) {
	state := NewState()

	start := time.Now()
	if err := state.WaitDrain(0); err != nil {
		t.Fatalf("WaitDrain(0) empty = %v, want nil", err)
	}
	if elapsed := time.Since(start); elapsed > 20*time.Millisecond {
		t.Fatalf("WaitDrain(0) empty took %s, want immediate", elapsed)
	}
}

func TestTrackRequestAndWaitDrainBlocksUntilDone(t *testing.T) {
	state := NewState()
	doneRequest := state.TrackRequest()

	waitDone := make(chan error, 1)
	go func() {
		waitDone <- state.WaitDrain(200 * time.Millisecond)
	}()

	select {
	case err := <-waitDone:
		t.Fatalf("WaitDrain returned before request finished: %v", err)
	case <-time.After(25 * time.Millisecond):
		// Expected: still waiting while request is in flight.
	}

	doneRequest()
	select {
	case err := <-waitDone:
		if err != nil {
			t.Fatalf("WaitDrain after done = %v, want nil", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("WaitDrain did not return after request finished")
	}

	// The returned done function is idempotent.
	doneRequest()
	if err := state.WaitDrain(0); err != nil {
		t.Fatalf("WaitDrain after idempotent done = %v, want nil", err)
	}
}

func TestWaitDrainTimesOutWhenInflight(t *testing.T) {
	state := NewState()
	doneRequest := state.TrackRequest()
	defer doneRequest()

	err := state.WaitDrain(10 * time.Millisecond)
	if !errors.Is(err, ErrDrainTimeout) {
		t.Fatalf("WaitDrain timeout = %v, want ErrDrainTimeout", err)
	}
}

func TestConcurrentShutdownWhileRequestsInflight(t *testing.T) {
	state := NewState()
	const requestCount = 8

	doneRequests := make([]func(), 0, requestCount)
	for range requestCount {
		doneRequests = append(doneRequests, state.TrackRequest())
	}

	state.BeginShutdown()
	if !state.IsShuttingDown() {
		t.Fatal("state should report shutdown during in-flight drain")
	}

	waitDone := make(chan error, 1)
	go func() {
		waitDone <- state.WaitDrain(200 * time.Millisecond)
	}()

	select {
	case err := <-waitDone:
		t.Fatalf("WaitDrain returned while requests still in flight: %v", err)
	case <-time.After(25 * time.Millisecond):
		// Expected: blocked until all done functions run.
	}

	for _, done := range doneRequests {
		done()
	}

	select {
	case err := <-waitDone:
		if err != nil {
			t.Fatalf("WaitDrain after concurrent drain = %v, want nil", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("WaitDrain did not return after all in-flight requests finished")
	}
}

func TestStateContextHelpers(t *testing.T) {
	state := NewState()
	ctx := WithState(context.Background(), state)

	got, ok := FromContext(ctx)
	if !ok {
		t.Fatal("FromContext did not find lifecycle state")
	}
	if got != state {
		t.Fatalf("FromContext returned %p, want %p", got, state)
	}

	if _, ok := FromContext(context.Background()); ok {
		t.Fatal("FromContext should not find state in a plain context")
	}

	plain := context.Background()
	if WithState(plain, nil) != plain {
		t.Fatal("WithState(nil) should leave context unchanged")
	}
}

package lifecycle

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// warmupModelFunc and discardLogger are shared helpers already defined in
// sibling lifecycle test files (warmup_test.go, shutdown_test.go). Recovery
// tests reuse them to stay consistent with the existing package test
// conventions. waitForCondition is defined here because no sibling test
// currently needs it; mirroring api/main_test.go's helper with a longer
// 2s deadline for robustness on loaded CI.
func waitForCondition(t *testing.T, pred func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if pred() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("condition not met within 2s")
}

func TestStartWarmupRecoveryDisabledIsNoop(t *testing.T) {
	state := NewState()
	var calls atomic.Int32
	model := warmupModelFunc(func(context.Context, []string) ([][]float32, error) {
		calls.Add(1)
		return [][]float32{{1}}, nil
	})

	StartWarmupRecovery(context.Background(), discardLogger(), state, model, time.Second, 10*time.Millisecond, false)

	// Give the goroutine a generous window to prove it never runs.
	time.Sleep(50 * time.Millisecond)
	if got := calls.Load(); got != 0 {
		t.Fatalf("disabled recovery called model %d times, want 0", got)
	}
	if state.IsWarmupDone() {
		t.Fatal("disabled recovery should not mark warmup done")
	}
}

func TestStartWarmupRecoveryAlreadyWarmNoop(t *testing.T) {
	state := NewState()
	state.MarkWarmupDone() // already warm before recovery starts
	var calls atomic.Int32
	model := warmupModelFunc(func(context.Context, []string) ([][]float32, error) {
		calls.Add(1)
		return [][]float32{{1}}, nil
	})

	StartWarmupRecovery(context.Background(), discardLogger(), state, model, time.Second, 10*time.Millisecond, true)

	// The goroutine should observe IsWarmupDone on its first check and exit
	// without invoking the model.
	time.Sleep(50 * time.Millisecond)
	if got := calls.Load(); got != 0 {
		t.Fatalf("already-warm recovery called model %d times, want 0", got)
	}
}

func TestStartWarmupRecoveryRetriesAndMarksReady(t *testing.T) {
	state := NewState()
	boom := errors.New("boom")
	var attempts atomic.Int32
	model := warmupModelFunc(func(context.Context, []string) ([][]float32, error) {
		if attempts.Add(1) < 3 {
			return nil, boom
		}
		return [][]float32{{1}}, nil
	})

	StartWarmupRecovery(context.Background(), discardLogger(), state, model, time.Second, 10*time.Millisecond, true)

	waitForCondition(t, state.IsReady)
	if got := attempts.Load(); got != 3 {
		t.Fatalf("recovery attempts = %d, want 3", got)
	}
	if err := state.LastError(); err != nil {
		t.Fatalf("LastError after recovery success = %v, want nil", err)
	}
}

func TestStartWarmupRecoveryStopsAtSuccess(t *testing.T) {
	state := NewState()
	var attempts atomic.Int32
	model := warmupModelFunc(func(context.Context, []string) ([][]float32, error) {
		attempts.Add(1)
		return [][]float32{{1}}, nil
	})

	StartWarmupRecovery(context.Background(), discardLogger(), state, model, time.Second, 10*time.Millisecond, true)

	waitForCondition(t, state.IsReady)
	// Wait past at least one more tick to confirm the goroutine exited.
	time.Sleep(40 * time.Millisecond)
	first := attempts.Load()
	time.Sleep(40 * time.Millisecond)
	second := attempts.Load()
	if second != first {
		t.Fatalf("recovery continued after success: attempts grew %d -> %d", first, second)
	}
	if first != 1 {
		t.Fatalf("recovery attempts = %d, want exactly 1", first)
	}
}

func TestStartWarmupRecoveryStopsAtShutdown(t *testing.T) {
	state := NewState()
	boom := errors.New("boom")
	var attempts atomic.Int32
	model := warmupModelFunc(func(context.Context, []string) ([][]float32, error) {
		attempts.Add(1)
		return nil, boom
	})

	StartWarmupRecovery(context.Background(), discardLogger(), state, model, time.Second, 15*time.Millisecond, true)

	// Let a couple of attempts run, then begin shutdown.
	time.Sleep(40 * time.Millisecond)
	state.BeginShutdown()
	// Snapshot attempts right after shutdown signal; then wait multiple
	// intervals and confirm no growth.
	time.Sleep(20 * time.Millisecond)
	beforeWait := attempts.Load()
	time.Sleep(80 * time.Millisecond)
	afterWait := attempts.Load()
	if afterWait != beforeWait {
		t.Fatalf("recovery continued after shutdown: attempts grew %d -> %d", beforeWait, afterWait)
	}
	if state.IsWarmupDone() {
		t.Fatal("recovery should not mark warmup done on shutdown")
	}
}

func TestStartWarmupRecoveryStopsAtContextCancel(t *testing.T) {
	state := NewState()
	boom := errors.New("boom")
	var attempts atomic.Int32
	model := warmupModelFunc(func(context.Context, []string) ([][]float32, error) {
		attempts.Add(1)
		return nil, boom
	})

	ctx, cancel := context.WithCancel(context.Background())
	StartWarmupRecovery(ctx, discardLogger(), state, model, time.Second, 15*time.Millisecond, true)

	time.Sleep(40 * time.Millisecond)
	cancel()
	// Snapshot and confirm growth stops after cancel.
	time.Sleep(15 * time.Millisecond)
	beforeWait := attempts.Load()
	time.Sleep(80 * time.Millisecond)
	afterWait := attempts.Load()
	if afterWait != beforeWait {
		t.Fatalf("recovery continued after ctx cancel: attempts grew %d -> %d", beforeWait, afterWait)
	}
}

func TestStartWarmupRecoveryNilArgsDontPanic(t *testing.T) {
	// These branches must not spawn a goroutine and must not panic.
	StartWarmupRecovery(context.Background(), discardLogger(), nil, nil, time.Second, 10*time.Millisecond, true)
	StartWarmupRecovery(context.Background(), discardLogger(), NewState(), nil, time.Second, 10*time.Millisecond, true)
}

func TestStartWarmupRecoveryNonPositiveIntervalNoop(t *testing.T) {
	state := NewState()
	var calls atomic.Int32
	model := warmupModelFunc(func(context.Context, []string) ([][]float32, error) {
		calls.Add(1)
		return [][]float32{{1}}, nil
	})

	StartWarmupRecovery(context.Background(), discardLogger(), state, model, time.Second, 0, true)
	time.Sleep(30 * time.Millisecond)
	if got := calls.Load(); got != 0 {
		t.Fatalf("zero-interval recovery called model %d times, want 0", got)
	}
}

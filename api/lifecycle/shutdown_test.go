package lifecycle

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

type fakeGracefulServer struct {
	shutdownFunc func(context.Context) error
	closeCount   atomic.Int64
}

func (s *fakeGracefulServer) Shutdown(ctx context.Context) error {
	if s.shutdownFunc != nil {
		return s.shutdownFunc(ctx)
	}
	return nil
}

func (s *fakeGracefulServer) Close() error {
	s.closeCount.Add(1)
	return nil
}

func TestAwaitSignalAndShutdownIdleCleanExit(t *testing.T) {
	state := NewState()
	logger := discardLogger()
	signals := make(chan os.Signal, 1)
	shutdownCalled := make(chan struct{})
	server := &fakeGracefulServer{shutdownFunc: func(_ context.Context) error {
		if !state.IsShuttingDown() {
			t.Error("state should be shutting down before server shutdown")
		}
		close(shutdownCalled)
		return nil
	}}
	signals <- syscall.SIGTERM
	started := time.Now()

	if err := AwaitSignalAndShutdown(context.Background(), signals, server, state, logger, time.Second); err != nil {
		t.Fatalf("AwaitSignalAndShutdown returned error: %v", err)
	}
	if elapsed := time.Since(started); elapsed >= time.Second {
		t.Fatalf("idle shutdown took %s, want < 1s", elapsed)
	}
	select {
	case <-shutdownCalled:
	default:
		t.Fatal("server shutdown was not called")
	}
	if got := server.closeCount.Load(); got != 0 {
		t.Fatalf("Close calls = %d, want 0", got)
	}
}

func TestGracefulShutdownWaitsForInflightThenCleanExit(t *testing.T) {
	state := NewState()
	doneRequest := state.TrackRequest()
	server := &fakeGracefulServer{}
	go func() {
		time.Sleep(10 * time.Millisecond)
		doneRequest()
	}()

	if err := GracefulShutdown(context.Background(), "SIGTERM", server, state, discardLogger(), time.Second); err != nil {
		t.Fatalf("GracefulShutdown returned error: %v", err)
	}
	if err := state.WaitDrain(0); err != nil {
		t.Fatalf("WaitDrain after shutdown = %v, want nil", err)
	}
	if got := server.closeCount.Load(); got != 0 {
		t.Fatalf("Close calls = %d, want 0", got)
	}
}

func TestGracefulShutdownForceClosesOnTimeout(t *testing.T) {
	state := NewState()
	_ = state.TrackRequest()
	server := &fakeGracefulServer{shutdownFunc: func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}}
	started := time.Now()

	err := GracefulShutdown(context.Background(), "SIGTERM", server, state, discardLogger(), 10*time.Millisecond)
	if !errors.Is(err, ErrShutdownTimeout) && !errors.Is(err, ErrDrainTimeout) {
		t.Fatalf("GracefulShutdown error = %v, want shutdown/drain timeout", err)
	}
	if elapsed := time.Since(started); elapsed > 500*time.Millisecond {
		t.Fatalf("forced shutdown took %s, want <= 500ms", elapsed)
	}
	if got := server.closeCount.Load(); got != 1 {
		t.Fatalf("Close calls = %d, want 1", got)
	}
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

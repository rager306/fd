package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// DefaultShutdownTimeout is the maximum time fd gives active handlers to drain
// after SIGTERM or SIGINT before force-closing the HTTP server.
const DefaultShutdownTimeout = 30 * time.Second

// ErrShutdownTimeout is returned when graceful shutdown exceeds its deadline.
var ErrShutdownTimeout = errors.New("lifecycle shutdown timeout")

// GracefulServer is the http.Server shutdown surface used by lifecycle.
type GracefulServer interface {
	Shutdown(ctx context.Context) error
	Close() error
}

// AwaitSignalAndShutdown waits for one process signal, marks lifecycle shutdown,
// drains the HTTP server and in-flight lifecycle requests, then returns nil for
// a clean exit or an error when force close was required.
func AwaitSignalAndShutdown(
	ctx context.Context,
	signals <-chan os.Signal,
	server GracefulServer,
	state *State,
	logger *slog.Logger,
	timeout time.Duration,
) error {
	select {
	case sig := <-signals:
		return GracefulShutdown(ctx, sig.String(), server, state, logger, timeout)
	case <-ctx.Done():
		return ctx.Err()
	}
}

// GracefulShutdown marks shutdown, waits for the HTTP server and lifecycle
// in-flight tracker to drain under one deadline, and force-closes on timeout.
func GracefulShutdown(
	ctx context.Context,
	signalName string,
	server GracefulServer,
	state *State,
	logger *slog.Logger,
	timeout time.Duration,
) error {
	if timeout <= 0 {
		timeout = DefaultShutdownTimeout
	}
	if state != nil {
		state.BeginShutdown()
	}
	logger.Info("shutdown signal received", "signal", signalName, "timeout", timeout.String())

	shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	serverDone := make(chan error, 1)
	go func() {
		serverDone <- server.Shutdown(shutdownCtx)
	}()
	drainDone := make(chan error, 1)
	go func() {
		if state == nil {
			drainDone <- nil
			return
		}
		drainDone <- state.WaitDrain(timeout)
	}()

	var serverErr error
	var drainErr error
	for serverDone != nil || drainDone != nil {
		select {
		case err := <-serverDone:
			serverErr = err
			serverDone = nil
		case err := <-drainDone:
			drainErr = err
			drainDone = nil
		case <-shutdownCtx.Done():
			return forceClose(server, logger, fmt.Errorf("%w: %w", ErrShutdownTimeout, shutdownCtx.Err()))
		}
	}

	if serverErr != nil || drainErr != nil {
		return forceClose(server, logger, errors.Join(serverErr, drainErr))
	}
	logger.Info("shutdown complete")
	return nil
}

func forceClose(server GracefulServer, logger *slog.Logger, shutdownErr error) error {
	if closeErr := server.Close(); closeErr != nil {
		logger.Error("server force close failed", "error", closeErr)
		return errors.Join(shutdownErr, closeErr)
	}
	logger.Error("server force closed during shutdown", "error", shutdownErr)
	return shutdownErr
}

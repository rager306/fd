package lifecycle

import (
	"context"
	"log/slog"
	"time"

	"fd-api/embed"
)

// StartWarmupRecovery runs a background goroutine that periodically retries
// model warmup after the startup warmup loop has exhausted its attempts.
// This eliminates the fd-api "degraded forever" failure mode where a race
// with TEI startup (CPU BERT load ~15-20s) leaves model_loaded=false and
// requires a manual `docker restart fd_api`.
//
// Behaviour:
//   - If enabled is false, the function is a no-op and returns immediately.
//   - The goroutine exits on any of: ctx cancelled, state.IsWarmupDone(),
//     state.IsShuttingDown(), or a successful PreWarm.
//   - Each attempt uses a per-call timeout (timeout) derived from the parent
//     context, so a slow TEI does not pin the goroutine.
//   - Errors are recorded via state.SetLastError so /health.last_error reflects
//     the current recovery state. Successful recovery clears it via
//     state.MarkWarmupDone().
//
// StartWarmupRecovery does not coordinate with the manual POST /warmup
// handler or the startup warmup goroutine. PreWarm is idempotent on a given
// model; concurrent attempts are safe but may double-count in logs. The
// goroutine is the single long-lived recovery loop for the process.
func StartWarmupRecovery(
	ctx context.Context,
	logger *slog.Logger,
	state *State,
	model embed.Embedder,
	timeout, interval time.Duration,
	enabled bool,
) {
	if !enabled {
		logger.Info("warmup recovery disabled")
		return
	}
	if interval <= 0 {
		logger.Warn("warmup recovery enabled with non-positive interval; disabling", "interval", interval.String())
		return
	}
	if state == nil || model == nil {
		logger.Warn("warmup recovery disabled; state or model is nil")
		return
	}

	go runWarmupRecovery(ctx, logger, state, model, timeout, interval)
}

func runWarmupRecovery(
	ctx context.Context,
	logger *slog.Logger,
	state *State,
	model embed.Embedder,
	timeout, interval time.Duration,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	attempt := 0
	started := time.Now()
	logger.Info("warmup recovery started", "interval", interval.String(), "timeout", timeout.String())

	for {
		// Fast exit on terminal conditions before waiting on the ticker.
		if state.IsWarmupDone() {
			logger.Info("warmup recovery stopped", "reason", "already_warm", "attempts", attempt, "elapsed_ms", time.Since(started).Milliseconds())
			return
		}
		if state.IsShuttingDown() {
			logger.Info("warmup recovery stopped", "reason", "shutdown", "attempts", attempt, "elapsed_ms", time.Since(started).Milliseconds())
			return
		}
		if ctx.Err() != nil {
			logger.Info("warmup recovery stopped", "reason", "context_cancelled", "attempts", attempt, "elapsed_ms", time.Since(started).Milliseconds())
			return
		}

		select {
		case <-ctx.Done():
			logger.Info("warmup recovery stopped", "reason", "context_cancelled", "attempts", attempt, "elapsed_ms", time.Since(started).Milliseconds())
			return
		case <-ticker.C:
		}

		// Re-check terminal conditions after waking: state may have changed
		// while we were waiting (e.g. manual POST /warmup succeeded, or
		// shutdown began).
		if state.IsWarmupDone() {
			logger.Info("warmup recovery stopped", "reason", "warmed_outside_loop", "attempts", attempt, "elapsed_ms", time.Since(started).Milliseconds())
			return
		}
		if state.IsShuttingDown() {
			logger.Info("warmup recovery stopped", "reason", "shutdown", "attempts", attempt, "elapsed_ms", time.Since(started).Milliseconds())
			return
		}

		attempt++
		attemptStarted := time.Now()
		callCtx, cancel := context.WithTimeout(ctx, timeout)
		err := PreWarm(callCtx, model)
		cancel()
		latencyMs := time.Since(attemptStarted).Milliseconds()

		if err == nil {
			state.MarkWarmupDone()
			logger.Info("warmup recovery succeeded",
				"attempt", attempt,
				"latency_ms", latencyMs,
				"elapsed_ms", time.Since(started).Milliseconds(),
			)
			return
		}

		state.SetLastError(err)
		logger.Warn("warmup recovery attempt failed",
			"attempt", attempt,
			"error", err,
			"latency_ms", latencyMs,
			"next_interval", interval.String(),
		)
	}
}
